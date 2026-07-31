"""Generates per-service and aggregate ``.env`` files from ``params.toml``.

Reads ``params.toml`` in the given deploy env folder, converts nested TOML sections
into ``KEY=value`` lines (quoting values that contain spaces), and writes:

* one ``.env`` per microservice under ``../../<service_name>/.env`` when that
  directory exists;
* a combined ``.env`` in the env folder with all services concatenated.

The CLI entry point is ``python generate_env.py <path_to_env_folder>``.
"""

import argparse
import logging
import os
import re
import typing

import tomllib

logger = logging.getLogger(__name__)

logging.basicConfig(
    format="%(asctime)s %(levelname)-8s %(message)s",
    level=logging.INFO,
    datefmt="%Y-%m-%d %H:%M:%S",
)


class GenerateEnvs:
    """Orchestrates ``.env`` generation from a shared TOML parameter manifest.

    Each top-level TOML key names a microservice (for example ``datacatalogue``).
    Nested tables become grouped sections in the output; optional ``comment`` keys
    at the service or section level render as banner comments in the ``.env`` file.

    Values (or parts of string values) may reference shared constants defined in
    the optional ``params.variables.toml`` manifest using ``${VARIABLE_NAME}``.
    This avoids retyping the same value (e.g. a shared Redis port, or DB
    credentials reused inside a connection string) in multiple places.

    Attributes:
        ALL_SETTINGS_TOML: Basename of the TOML manifest inside ``path_to_env_folder``.
        VARIABLES_TOML: Basename of the optional shared-variables TOML manifest
            inside ``path_to_env_folder``.
        VARIABLE_REF_PATTERN: Regex matching ``${VARIABLE_NAME}`` references inside
            string values.
        ENV_COMMENT: Multi-line banner template for section headers (one ``{}`` placeholder).
        COMMENT: Single-line comment template (one ``{}`` placeholder).
        PARAMETER: ``KEY=value`` line template without quoting.
        PARAMETER_WTH_STRING: ``KEY="value"`` line template for values containing spaces.
        path_to_env_folder: Directory that holds ``params.toml`` and the aggregate
            ``.env`` output.
        path_to_toml: Resolved filesystem path to ``params.toml`` after validation.
        path_to_variables_toml: Resolved filesystem path to ``params.variables.toml``,
            or ``""`` if that optional file is not present.
        variables: Flat name -> value map loaded from ``params.variables.toml``.
    """

    ALL_SETTINGS_TOML: str = "params.toml"
    VARIABLES_TOML: str = "params.variables.toml"
    VARIABLE_REF_PATTERN: re.Pattern = re.compile(r"\$\{(?P<name>[A-Za-z_][A-Za-z0-9_]*)\}")
    ENV_COMMENT: str = "# ----------------------------------------\n# {}\n# ----------------------------------------\n"
    COMMENT: str = "# {}"
    PARAMETER: str = "{}={}\n"
    PARAMETER_WTH_STRING: str = """{}="{}"\n"""

    path_to_env_folder: str = ""
    path_to_toml: str = ""
    path_to_variables_toml: str = ""
    variables: dict[str, typing.Any]

    def __init__(self) -> None:
        """Parses CLI arguments, validates ``params.toml``, and runs generation.

        Raises:
            ValueError: If the required ``path_to_env_folder`` argument is missing.
            FileExistsError: If ``params.toml`` is not found under the env folder.
        """
        self.variables = {}
        self.__check_arguments_and_parse()
        self.__check_if_setting_file_exists(
            folder_path_to_env_settings=self.path_to_env_folder
        )
        self.__check_if_variables_file_exists(
            folder_path_to_env_settings=self.path_to_env_folder
        )

        self.run_all()

    def run_all(self) -> None:
        """Reads settings, rebuilds the aggregate ``.env``, and writes per-service files.

        Removes any existing aggregate ``.env`` in ``path_to_env_folder``, then for each
        top-level key in the TOML manifest generates service content and appends it to
        the combined file while also attempting to write ``../../<service>/.env``.
        """
        self.variables = self.__read_toml_variables()
        setting = self.__read_toml_settings()
        self.__remove_main_env_file()
        for lib in setting:
            one_file_text = self.generate_one_enf_file(setting[lib])

            self.__write_env_to_lib(lib, one_file_text)

            with open(os.path.join(self.path_to_env_folder, ".env"), "w", encoding="utf-8") as text_file:
                text_file.write(self.ENV_COMMENT.format(lib) + one_file_text + "\n\n")

    def __write_env_to_lib(self, lib_name: str, env_file_text: str) -> None:
        """Writes the same ``.env`` body to every known layout for one microservice.

        Delegates to ``__write_file`` for the service source tree
        (``path_to_env_folder/../../<lib_name>``) and for the compose deploy tree
        (``path_to_env_folder/../compose/<lib_name>``). Missing directories are
        skipped with an info log; existing ones receive an overwritten ``.env``.

        Args:
            lib_name: Top-level TOML key and microservice folder name.
            env_file_text: Full ``.env`` body for that service.
        """
        logger.info(f"Начинаем обработку {lib_name}")
        library_path = os.path.join(self.path_to_env_folder, "../..", lib_name)
        self.__write_file(dir_path=library_path, env_file_text=env_file_text)
        deploy_path = os.path.join(self.path_to_env_folder, "../compose/", lib_name)
        self.__write_file(dir_path=deploy_path, env_file_text=env_file_text)

    def __write_file(self, dir_path: str, env_file_text: str) -> None:
        """Writes ``.env`` into ``dir_path`` when that directory exists.

        Args:
            dir_path: Target directory; the file is created as ``<dir_path>/.env``.
            env_file_text: Full ``.env`` body to write.
        """
        if os.path.isdir(dir_path):
            with open(os.path.join(dir_path, ".env"), "w", encoding="utf-8") as text_file:
                text_file.write(env_file_text)
            logger.info(f"""Сохранен файл .env в {dir_path}""")
        else:
            logger.info(f"""Папка {dir_path} отстутствует""")

    def generate_one_enf_file(self, params: dict[str, typing.Any]) -> str:
        """Builds one service ``.env`` body from a nested TOML subtree.

        Walks two nesting levels under the service key. Optional ``comment`` entries at
        the service root or under a section name are consumed and emitted as banner
        comments. Leaf values are rendered as ``KEY=value``; values whose string form
        contains a space are wrapped in double quotes.

        Args:
            params: Nested dict for one microservice (mutated in place: ``comment``
                keys are removed while building output).

        Returns:
            The assembled ``.env`` text for that service, including section banners.
        """
        env_file_text: str = ""

        if "comment" in params:
            env_file_text += self.ENV_COMMENT.format(params["comment"])
            del params["comment"]

        for param_level_2 in params:
            if "comment" in params[param_level_2]:
                env_file_text += "\n" + self.ENV_COMMENT.format(
                    params[param_level_2]["comment"]
                )
                del params[param_level_2]["comment"]

            for param_level_3 in params[param_level_2]:
                inserting_param = self.__resolve_variable_refs(
                    params[param_level_2][param_level_3]
                )
                if " " in str(inserting_param):
                    env_file_text += self.PARAMETER_WTH_STRING.format(
                        param_level_3, inserting_param
                    )
                else:
                    env_file_text += self.PARAMETER.format(
                        param_level_3, inserting_param
                    )

        return env_file_text

    def __resolve_variable_refs(self, value: typing.Any) -> typing.Any:
        """Substitutes ``${VARIABLE_NAME}`` references with values from ``variables``.

        A value that is *exactly* one reference (e.g. ``"${REDIS_PORT}"``) is replaced
        by the referenced value as-is, preserving its original TOML type (int, bool, ...).
        A reference embedded inside a larger string (e.g. a connection string) is
        substituted as text.

        Args:
            value: Raw leaf value from ``params.toml``. Non-string values are returned
                unchanged, since only strings can contain ``${...}`` references.

        Returns:
            The value with all variable references resolved.

        Raises:
            KeyError: If a referenced variable is not defined in ``params.variables.toml``.
        """
        if not isinstance(value, str):
            return value

        full_match = self.VARIABLE_REF_PATTERN.fullmatch(value)
        if full_match:
            return self.__lookup_variable(full_match.group("name"))

        return self.VARIABLE_REF_PATTERN.sub(
            lambda match: str(self.__lookup_variable(match.group("name"))), value
        )

    def __lookup_variable(self, name: str) -> typing.Any:
        """Looks up a variable by name, raising a descriptive error if it is missing.

        Args:
            name: Variable name referenced as ``${name}`` in ``params.toml``.

        Returns:
            The value defined for ``name`` in ``params.variables.toml``.

        Raises:
            KeyError: If ``name`` is not defined in ``params.variables.toml``.
        """
        if name not in self.variables:
            error_message = (
                f'Переменная "${{{name}}}" не найдена в {self.VARIABLES_TOML}'
            )
            logger.error(error_message)
            raise KeyError(error_message)

        return self.variables[name]

    def __check_arguments_and_parse(self) -> None:
        """Reads ``path_to_env_folder`` from the first positional CLI argument.

        Raises:
            ValueError: If ``path_to_env_folder`` is not provided on the command line.
        """
        folder_path_to_env_settings: str = ""
        parser = argparse.ArgumentParser()
        parser.add_argument(
            "path_to_env_folder",
            help="Путь к папке, где лежит файл с настройками и куда будут сохраняться env файлы",
            type=str,
        )
        args = parser.parse_args()

        if args.path_to_env_folder:
            folder_path_to_env_settings = args.path_to_env_folder
        else:
            error_message: str = """Аргумент "path_to_env_folder" не передан """
            logger.error(error_message)
            raise ValueError(error_message)

        logger.info(f"Путь к папке с параметрами env: '{folder_path_to_env_settings}'")

        self.path_to_env_folder = folder_path_to_env_settings

    def __remove_main_env_file(self) -> None:
        """Deletes the aggregate ``.env`` in ``path_to_env_folder`` if it already exists.

        Ensures a clean rebuild of the combined env file on each run.
        """
        main_env_file = os.path.join(self.path_to_env_folder, ".env")
        if os.path.isfile(main_env_file):
            logger.info(f"""Найден {main_env_file} файл""")
            os.remove(main_env_file)
            logger.info(f"""Удален {main_env_file} файл""")

    def __check_if_setting_file_exists(self, folder_path_to_env_settings: str) -> None:
        """Resolves ``path_to_toml`` when ``params.toml`` is present in the env folder.

        Args:
            folder_path_to_env_settings: Directory expected to contain ``params.toml``.

        Raises:
            FileExistsError: If ``params.toml`` is missing under that directory.
        """
        if os.path.isfile(
            os.path.join(folder_path_to_env_settings, self.ALL_SETTINGS_TOML)
        ):
            self.path_to_toml = os.path.join(
                folder_path_to_env_settings, self.ALL_SETTINGS_TOML
            )
        else:
            error_message: str = f"""File "{os.path.join(folder_path_to_env_settings, self.ALL_SETTINGS_TOML)}" does not exist"""
            logger.error(error_message)
            raise FileExistsError(error_message)

    def __check_if_variables_file_exists(self, folder_path_to_env_settings: str) -> None:
        """Resolves ``path_to_variables_toml`` when ``params.variables.toml`` is present.

        The file is optional: if it is absent, ``path_to_variables_toml`` stays empty
        and no variable substitution is available in ``params.toml``.

        Args:
            folder_path_to_env_settings: Directory that may contain
                ``params.variables.toml`` alongside ``params.toml``.
        """
        candidate_path = os.path.join(
            folder_path_to_env_settings, self.VARIABLES_TOML
        )
        if os.path.isfile(candidate_path):
            self.path_to_variables_toml = candidate_path
            logger.info(f"Найден файл с переменными: '{candidate_path}'")
        else:
            logger.info(
                f"Файл с переменными '{candidate_path}' отсутствует, подстановка ${{VAR}} недоступна"
            )

    def __read_toml_variables(self) -> dict[str, typing.Any]:
        """Loads the flat shared-variable manifest from ``path_to_variables_toml``.

        Returns:
            Parsed TOML as a flat dict, or an empty dict if ``params.variables.toml``
            is not present.
        """
        if not self.path_to_variables_toml:
            return {}

        with open(self.path_to_variables_toml, mode="rb") as fp:
            return tomllib.load(fp)

    def __read_toml_settings(self) -> dict[str, typing.Any]:
        """Loads the full parameter manifest from ``path_to_toml``.

        Returns:
            Parsed TOML as a nested dict keyed by microservice name.
        """
        with open(self.path_to_toml, mode="rb") as fp:
            setting = tomllib.load(fp)

        return setting


if __name__ == "__main__":
    GenerateEnvs()
