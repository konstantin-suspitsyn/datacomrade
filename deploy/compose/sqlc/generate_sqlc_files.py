"""Generates per-service sqlc layouts from YAML and PostgreSQL DDL dumps.

Loads ``schema_to_create.yml`` to discover microservices and models, refreshes
``schema.sql`` via the local Docker-backed dump script, then writes ``sqlc.yml`` and
per-model ``schema.sql`` / ``query.sql`` paths under each service's ``db/sqlc`` tree.

The CLI entry point is ``python generate_sqlc_files.py <path_to_sqlc_config>``; the
sqlc config path is stored on ``SqlCreator`` for tooling that invokes this module.
"""

import logging
import os
import subprocess
import sys
from pathlib import Path

import pydantic
import yaml

logger = logging.getLogger(__name__)


class SqlcModel(pydantic.BaseModel):
    """One sqlc codegen package and the PostgreSQL tables it includes.

    Corresponds to a single key under ``models`` in ``schema_to_create.yml``: a model
    name and the schema-qualified relations that belong in that package's
    ``schema.sql``.

    Attributes:
        name: Model directory and sqlc package name (YAML key under ``models``).
        tables: Table names as ``schema.table``; each entry must include a schema
            prefix (for example ``public.users`` or ``dc.database_type``).
    """

    name: str
    tables: list[str]


class SqlcMicroservice(pydantic.BaseModel):
    """One deployable backend and its sqlc output layout.

    Corresponds to one top-level key in ``schema_to_create.yml`` (for example
    ``datacomrade``).

    Attributes:
        name: Microservice key from the YAML root (often the repository folder name).
        path: Relative path segment for generated Go (or other) output; from the
            ``save_path`` field in YAML.
        env_path: Path to the env file passed to the schema dump script for this
            service (credentials and connection for Docker PostgreSQL).
        sqlc_models: All ``SqlcModel`` instances declared under this microservice's
            ``models`` block.
    """

    name: str
    path: str
    env_path: str
    sqlc_models: list[SqlcModel]


class SqlCreator:
    """Orchestrates DDL dumps and sqlc file generation for configured microservices.

    Reads ``schema_to_create.yml``, runs ``generate_schema_local.sh`` per service to
    refresh ``schema.sql``, slices ``CREATE SCHEMA`` / ``CREATE TABLE`` fragments from
    that dump, and writes ``sqlc.yml`` plus per-model ``schema.sql`` beside each
    service tree.

    Attributes:
        shell_script_relative_path: Bash script under ``current_folder_path`` that
            dumps PostgreSQL DDL into ``SQL_FILE_PATH`` using the env file passed at
            runtime.
        path_to_sqlc_config: Path to the sqlc config file from the first CLI argument
            (``sys.argv[1]``); stored for callers that need it when invoking this
            module.
        current_folder_path: Directory containing this script, templates, and the
            transient ``schema.sql`` dump.
        SQL_FILE_PATH: Basename of the dump file read and removed during generation.
        POSTGRES_SQLC_TEMPLATE_FOLDER: Subdirectory with shared PostgreSQL sqlc YAML
            fragments (``header.yml``, ``each_schema.yml``).
        POSTGRES_TEMPLATE_HEADER: Relative path to the sqlc YAML header template.
        POSTGRES_TEMPLATE_SCHEMA: Relative path to the per-schema sqlc block template.
        SCHEMA_TO_CREATE: Basename of the microservice/model manifest YAML.
    """

    shell_script_relative_path: str = "generate_schema_local.sh"
    path_to_sqlc_config: str = ""
    current_folder_path: Path = Path(__file__).resolve().parent
    SQL_FILE_PATH: str = "schema.sql"
    POSTGRES_SQLC_TEMPLATE_FOLDER: str = r"template/postgres"
    POSTGRES_TEMPLATE_HEADER: str = os.path.join(
        POSTGRES_SQLC_TEMPLATE_FOLDER, "header.yml"
    )
    POSTGRES_TEMPLATE_SCHEMA: str = os.path.join(
        POSTGRES_SQLC_TEMPLATE_FOLDER, "each_schema.yml"
    )
    SCHEMA_TO_CREATE: str = r"schema_to_create.yml"

    def __init__(self) -> None:
        """Initializes paths from CLI arguments.

        Expects ``sys.argv[1]`` to be the path to the sqlc config file consumed by
        tooling that runs this script.

        Raises:
            ValueError: If no argument is present after the script name.
        """
        cmd_variables: list[str] = sys.argv[1:]
        if len(cmd_variables) < 1:
            raise ValueError(
                "Not enough variables from command line were provided. "
                "Expected: script.py path_to_sqlc_config path_to_env"
            )

        self.path_to_sqlc_config = cmd_variables[0]

    def get_table_dump_from_postgres_docker(self, path_to_env: str) -> None:
        """Refreshes ``schema.sql`` by running the local Docker-backed PostgreSQL dump.

        Args:
            path_to_env: Filesystem path to the env file passed to
                ``generate_schema_local.sh`` (database connection and credentials).

        Raises:
            subprocess.CalledProcessError: If the shell script exits with a non-zero
                status.
            OSError: If ``bash`` cannot be started (for example executable not found).
        """
        command: str = f"bash, ./{self.shell_script_relative_path}, {path_to_env}, {self.SQL_FILE_PATH}"
        logger.info(f"Running {command}")

        try:
            subprocess.run(
                [
                    "bash",
                    f"./{self.shell_script_relative_path}",
                    f"{path_to_env}",
                    f"./{self.SQL_FILE_PATH}",
                ],
                cwd=self.current_folder_path,
                check=True,
            )
        except subprocess.CalledProcessError as e:
            logger.error(
                f"Schema dump script failed (exit {e.returncode}). Schema was not copied from Docker.",
            )
            raise
        except OSError as e:
            logger.error(f"Could not run schema dump script: {e}")
            raise

    def __get_create_schema_lines(self) -> list[str]:
        """Collects lines from the DDL dump that mention ``CREATE SCHEMA``.

        Scans ``schema.sql`` sequentially and keeps every line where the lowercased
        text contains the substring ``create schema``.

        Returns:
            All matching lines in file order, each including its trailing newline if
            present in the source file.

        Raises:
            ValueError: If the dump file is missing or cannot be read; the chained
                exception preserves the original ``FileNotFoundError`` or ``OSError``.
        """
        create_schemas: list[str] = []

        schema_path = self.current_folder_path / self.SQL_FILE_PATH

        try:
            with open(schema_path, encoding="utf-8") as file:
                for line in file:
                    if "create schema" in line.lower():
                        create_schemas.append(line)
        except FileNotFoundError as e:
            raise ValueError(f"Schema dump file not found: {schema_path}") from e
        except OSError as e:
            raise ValueError(f"Cannot read schema dump file {schema_path}: {e}") from e

        return create_schemas

    def __copy_create_table_lines(self, table_name: str) -> list[str]:
        """Extracts the contiguous ``CREATE TABLE`` statement for one qualified table.

        Scans ``schema.sql`` from the first matching ``CREATE TABLE`` through the line
        that contains the terminating semicolon.

        Args:
            table_name: Qualified name ``schema.table`` as in the dump; matching is
                case-insensitive and ignores double quotes around identifiers.

        Returns:
            Lines of the statement in file order, preserving newlines from the dump.

        Raises:
            ValueError: If ``table_name`` contains no ``.``, if the table is not found,
                or if the dump file is missing or unreadable (wrapping
                ``FileNotFoundError`` or ``OSError``).
        """
        if "." not in table_name:
            raise ValueError(
                f"В таблице {table_name} не найдена схема. Вид названия таблицы должен быть схема.таблица"
            )

        sql_create: list[str] = []
        starter: bool = False
        schema_path = self.current_folder_path / self.SQL_FILE_PATH

        try:
            with open(schema_path, encoding="utf-8") as file:
                for line in file:
                    if (
                        ("create table" in line.lower())
                        and (f"{table_name.lower()} " in line.replace('"', "").lower())
                    ) or starter:
                        starter = True
                        sql_create.append(line)
                    if starter and ";" in line:
                        break
        except FileNotFoundError as e:
            raise ValueError(f"Schema dump file not found: {schema_path}") from e
        except OSError as e:
            raise ValueError(f"Cannot read schema dump file {schema_path}: {e}") from e

        if len(sql_create) == 0:
            raise ValueError(f"Table {table_name} was not found")

        return sql_create

    def create_overall_schema(self, path_to_env: str) -> None:
        """Refreshes the local ``schema.sql`` using the given env file.

        Args:
            path_to_env: Path forwarded to ``get_table_dump_from_postgres_docker``.
        """
        self.get_table_dump_from_postgres_docker(path_to_env=path_to_env)

    def clean_up(self) -> None:
        """Deletes ``schema.sql`` under ``current_folder_path``.

        Raises:
            OSError: If the file cannot be removed (for example permission denied).
        """
        os.remove(self.current_folder_path / self.SQL_FILE_PATH)

    def __read_header_postgres_yml(self) -> list[str]:
        """Loads the shared PostgreSQL sqlc YAML header fragment.

        Returns:
            Lines read from ``POSTGRES_TEMPLATE_HEADER`` (including newlines).

        Raises:
            FileNotFoundError: If the header template is missing.
            OSError: If the file cannot be read.
        """
        yml_text_buffer: list[str] = []
        header_path = self.current_folder_path / self.POSTGRES_TEMPLATE_HEADER
        try:
            with open(header_path, encoding="utf-8") as p_header:
                for line in p_header:
                    yml_text_buffer.append(line)
        except FileNotFoundError:
            logger.error(f"sqlc header template not found: {header_path}")
            raise
        except OSError as e:
            logger.error(f"Cannot read sqlc header template {header_path}: {e}")
            raise

        return yml_text_buffer

    def __read_schema_to_create_postgres(
        self, placeholder_replacement: dict[str, str]
    ) -> list[str]:
        """Loads the per-model sqlc YAML block and substitutes placeholders.

        Each template line is passed through ``str.format`` with
        ``placeholder_replacement`` plus ``end=""`` so literal braces in YAML remain
        valid.

        Args:
            placeholder_replacement: Keys such as ``schema``, ``queries``, ``package``,
                and ``out`` for ``template/postgres/each_schema.yml``.

        Returns:
            Template lines with placeholders expanded.

        Raises:
            FileNotFoundError: If the per-schema template is missing.
            OSError: If the file cannot be read.
            KeyError: If a required placeholder is missing from the mapping.
        """
        yml_text_buffer: list[str] = []
        template_path = self.current_folder_path / self.POSTGRES_TEMPLATE_SCHEMA
        try:
            with open(template_path, encoding="utf-8") as schema_to_create:
                for line in schema_to_create:
                    yml_text_buffer.append(
                        line.format(**placeholder_replacement, end="")
                    )
        except FileNotFoundError:
            logger.error(f"sqlc schema template not found: {template_path}")
            raise
        except OSError as e:
            logger.error(f"Cannot read sqlc schema template {template_path}: {e}")
            raise

        return yml_text_buffer

    def save_yml_sql_postgres(self, sqlc_microservice: SqlcMicroservice) -> None:
        """Writes ``sqlc.yml`` and each model's ``schema.sql`` for one microservice.

        Resolves the output root to ``<repo>/<name>/db/sqlc/`` (three levels above
        this script). For every model, appends the shared header, one templated sqlc
        schema block, then all ``CREATE SCHEMA`` lines from the dump followed by
        ``CREATE TABLE`` slices for each listed table.

        Args:
            sqlc_microservice: Service name, relative codegen path (``save_path`` from
                YAML), and the models with their table lists.

        Raises:
            OSError: If directories cannot be created or files cannot be written.
            ValueError: If reading or slicing ``schema.sql`` fails (see
                ``__get_create_schema_lines`` and ``__copy_create_table_lines``).
        """
        microservice_name: str = sqlc_microservice.name

        # Каждый микросервис будет иметь следующую схему хранения файлов:
        # Имя микросервиса
        # └─── db/
        #      └─── sqlс/
        #           ├─── sqlc.yml
        #           └─── имя репозитория (для каждого репозитория)/
        #                ├─── schema.sql
        #                └─── query.sql

        save_dir = str(
            (
                self.current_folder_path
                / ".."
                / ".."
                / ".."
                / microservice_name
                / "db"
                / "sqlc"
            ).resolve()
        )
        sqlc_yml_path = os.path.join(save_dir, "sqlc.yml")
        repository_internal_path = os.path.join("..", "..", sqlc_microservice.path)

        sqlc_yml: list[str] = self.__read_header_postgres_yml()

        for repository_model in sqlc_microservice.sqlc_models:
            model_name: str = repository_model.name

            # С помощью этого словаря будут заменены плэйсхолдеры в файле each_schema
            text_placeholders: dict[str, str] = {}
            text_placeholders["schema"] = f"./{model_name}/schema.sql"
            text_placeholders["queries"] = f"./{model_name}/query.sql"
            text_placeholders["package"] = f"{model_name}"
            text_placeholders["out"] = os.path.join(
                repository_internal_path, model_name
            )

            # Добавили схему в sqlc.yml
            sqlc_yml.extend(self.__read_schema_to_create_postgres(text_placeholders))

            all_tables_creation: list[str] = []

            all_tables_creation.extend(self.__get_create_schema_lines())

            for table_name in repository_model.tables:
                all_tables_creation.extend(
                    self.__copy_create_table_lines(table_name=table_name)
                )

            # Сохраняем схему
            schema_sql_path: str = os.path.join(save_dir, model_name, "schema.sql")
            os.makedirs(os.path.dirname(schema_sql_path), exist_ok=True)
            with open(schema_sql_path, "w") as file:
                for line in all_tables_creation:
                    file.write(line)

        os.makedirs(os.path.dirname(sqlc_yml_path), exist_ok=True)
        with open(sqlc_yml_path, "w", encoding="utf-8") as file:
            for line in sqlc_yml:
                file.write(line)

    def __create_sql_microservices_structure(
        self,
        yaml_payload: dict,  # pyright: ignore[reportMissingTypeArgument]
    ) -> list[SqlcMicroservice]:
        """Builds ``SqlcMicroservice`` instances from the root YAML mapping.

        Each top-level value must provide ``models`` (name to model definition),
        ``save_path`` (string), and ``env_path`` (string for the dump script).

        Args:
            yaml_payload: Mapping from microservice name to its configuration dict.

        Returns:
            One ``SqlcMicroservice`` per top-level key, in YAML iteration order.

        Raises:
            ValueError: If ``models`` or ``save_path`` is missing or has the wrong
                type.
            KeyError: If ``env_path``, ``tables``, or other keys required to build
                ``SqlcModel`` are absent.
        """
        # Result container; one element per microservice key in yaml_payload.
        sql_microservices_data: list[SqlcMicroservice] = []

        for microservice_name, microservice_structure in yaml_payload.items():
            # Models declared under this microservice (`models` -> name -> fields).
            sqlc_models: list[SqlcModel] = []

            if "models" not in microservice_structure:
                raise ValueError(f"Microservice {microservice_name} has no models")
            if not isinstance(microservice_structure["models"], dict):
                raise ValueError(
                    f"Microservice {microservice_name} models must be a dictionary"
                )

            for model_name, model_value in microservice_structure["models"].items():
                sqlc_model: SqlcModel = SqlcModel(
                    name=model_name, tables=model_value["tables"]
                )
                sqlc_models.append(sqlc_model)

            if "save_path" not in microservice_structure:
                raise ValueError(f"Microservice {microservice_name} has no save_path")
            if not isinstance(microservice_structure["save_path"], str):
                raise ValueError(
                    f"Microservice {microservice_name} save_path must be a string"
                )

            sqlc_microservice = SqlcMicroservice(
                name=microservice_name,
                sqlc_models=sqlc_models,
                path=microservice_structure["save_path"],
                env_path=microservice_structure["env_path"],
            )

            sql_microservices_data.append(sqlc_microservice)

        return sql_microservices_data

    def __get_sql_models_from_yml(self, path: str) -> list[SqlcMicroservice]:
        """Loads ``schema_to_create.yml`` (or another manifest) and parses services.

        Relative paths are resolved against ``current_folder_path``.

        Args:
            path: Path to the YAML file; may be absolute or relative to this script.

        Returns:
            Parsed ``SqlcMicroservice`` list from ``__create_sql_microservices_structure``.

        Raises:
            ValueError: If the YAML root is not a mapping, or if the structure fails
                validation in ``__create_sql_microservices_structure``.
            yaml.YAMLError: If parsing fails (after logging the error).
            OSError: If the file cannot be opened or read.
            KeyError: Propagated when required keys are missing in the YAML tree.
        """
        # Parsed root object; must be a dict for `__create_sql_microservices_structure`.
        structure_to_generate: dict = {}  # pyright: ignore[reportMissingTypeArgument, reportUnknownVariableType]
        yml_path = Path(path)
        if not yml_path.is_absolute():
            yml_path = self.current_folder_path / yml_path

        with open(yml_path, encoding="utf-8") as stream:
            try:
                structure_to_generate = yaml.safe_load(stream)  # pyright: ignore[reportAny]
            except yaml.YAMLError as exc:
                logger.error("%s", exc)
                raise

        if not isinstance(structure_to_generate, dict):
            raise ValueError(
                f"YAML root must be a mapping, got {type(structure_to_generate).__name__}"
            )

        return self.__create_sql_microservices_structure(
            yaml_payload=structure_to_generate
        )

    def save_all_schemas_postgres(self) -> None:
        """Generates sqlc YAML and schema SQL for each microservice in the manifest.

        For each entry from ``SCHEMA_TO_CREATE``: dumps PostgreSQL DDL with that
        service's ``env_path``, writes ``sqlc.yml`` and model ``schema.sql`` files,
        then deletes the local ``schema.sql`` copy.

        Raises:
            ValueError: From the dump step, YAML parsing, or DDL slicing helpers.
            subprocess.CalledProcessError: If a schema dump script invocation fails.
            OSError: If files cannot be read or written.
            yaml.YAMLError: If the manifest YAML is invalid.
            KeyError: If the manifest omits required keys.
        """
        sql_service_models: list[SqlcMicroservice] = self.__get_sql_models_from_yml(
            self.SCHEMA_TO_CREATE
        )
        for service in sql_service_models:
            self.get_table_dump_from_postgres_docker(service.env_path)
            self.save_yml_sql_postgres(service)
            self.clean_up()
            logger.info(f"SQLC files for {service.name} generated successfully")


if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s:%(message)s",
    )
    logger.info("%s", sys.argv)

    sql_creator = SqlCreator()
    sql_creator.save_all_schemas_postgres()
