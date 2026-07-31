"""Tests for generate_env.py.

GenerateEnvs.__init__ parses argv and immediately runs generation, so most tests
build an instance with ``GenerateEnvs.__new__(GenerateEnvs)`` and set only the
attributes the method under test needs, instead of going through ``__init__``.
Name-mangled private methods are invoked via their mangled name
(``_GenerateEnvs__method``), same as the class itself would.
"""

import os

import pytest

from generate_env import GenerateEnvs


def make_instance(**attrs) -> GenerateEnvs:
    instance = GenerateEnvs.__new__(GenerateEnvs)
    instance.variables = {}
    instance.path_to_env_folder = ""
    instance.path_to_toml = ""
    instance.path_to_variables_toml = ""
    for name, value in attrs.items():
        setattr(instance, name, value)
    return instance


# ---------------------------------------------------------------------------
# __resolve_variable_refs / __lookup_variable
# ---------------------------------------------------------------------------


class TestResolveVariableRefs:
    def test_full_match_preserves_int_type(self):
        instance = make_instance(variables={"REDIS_PORT": 6379})
        result = instance._GenerateEnvs__resolve_variable_refs("${REDIS_PORT}")
        assert result == 6379
        assert isinstance(result, int)

    def test_full_match_preserves_bool_type(self):
        instance = make_instance(variables={"FEATURE_ON": True})
        result = instance._GenerateEnvs__resolve_variable_refs("${FEATURE_ON}")
        assert result is True

    def test_partial_match_interpolates_into_string(self):
        instance = make_instance(
            variables={
                "DB_USER": "postgres_user",
                "DB_PASSWORD": "pass",
                "DB_HOST": "localhost",
                "DB_PORT": 5432,
                "DB_NAME": "dcat",
            }
        )
        dsn = "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
        result = instance._GenerateEnvs__resolve_variable_refs(dsn)
        assert result == "postgres://postgres_user:pass@localhost:5432/dcat"

    def test_string_without_placeholder_is_unchanged(self):
        instance = make_instance(variables={})
        assert instance._GenerateEnvs__resolve_variable_refs("localhost") == "localhost"

    def test_non_string_values_pass_through_unchanged(self):
        instance = make_instance(variables={})
        assert instance._GenerateEnvs__resolve_variable_refs(5432) == 5432
        assert instance._GenerateEnvs__resolve_variable_refs(True) is True

    def test_missing_variable_raises_keyerror(self):
        instance = make_instance(variables={})
        with pytest.raises(KeyError, match="UNKNOWN_VAR"):
            instance._GenerateEnvs__resolve_variable_refs("${UNKNOWN_VAR}")

    def test_missing_variable_inside_larger_string_raises_keyerror(self):
        instance = make_instance(variables={"DB_USER": "postgres_user"})
        with pytest.raises(KeyError, match="DB_PASSWORD"):
            instance._GenerateEnvs__resolve_variable_refs("${DB_USER}:${DB_PASSWORD}")


# ---------------------------------------------------------------------------
# generate_one_enf_file
# ---------------------------------------------------------------------------


class TestGenerateOneEnfFile:
    def test_renders_service_and_section_comments(self):
        instance = make_instance(variables={})
        params = {
            "comment": "Service banner",
            "section": {"comment": "Section banner", "KEY": "value"},
        }
        text = instance.generate_one_enf_file(params)
        assert "Service banner" in text
        assert "Section banner" in text
        assert "KEY=value" in text

    def test_quotes_values_containing_spaces(self):
        instance = make_instance(variables={})
        params = {"section": {"ORIGINS": "http://a, http://b"}}
        text = instance.generate_one_enf_file(params)
        assert 'ORIGINS="http://a, http://b"\n' in text

    def test_does_not_quote_values_without_spaces(self):
        instance = make_instance(variables={})
        params = {"section": {"PORT": 8080}}
        text = instance.generate_one_enf_file(params)
        assert "PORT=8080\n" in text
        assert '"' not in text

    def test_comment_keys_are_removed_from_input_dict(self):
        instance = make_instance(variables={})
        params = {"comment": "banner", "section": {"comment": "inner", "KEY": "v"}}
        instance.generate_one_enf_file(params)
        assert "comment" not in params
        assert "comment" not in params["section"]

    def test_missing_service_level_comment_is_fine(self):
        instance = make_instance(variables={})
        params = {"section": {"KEY": "value"}}
        text = instance.generate_one_enf_file(params)
        assert text.strip() == "KEY=value"

    def test_substitutes_variable_references(self):
        instance = make_instance(variables={"REDIS_PORT": 6379})
        params = {"section": {"REDIS_PORT": "${REDIS_PORT}"}}
        text = instance.generate_one_enf_file(params)
        assert "REDIS_PORT=6379\n" in text


# ---------------------------------------------------------------------------
# __check_if_setting_file_exists
# ---------------------------------------------------------------------------


class TestCheckIfSettingFileExists:
    def test_sets_path_when_params_toml_present(self, tmp_path):
        (tmp_path / "params.toml").write_text("a=1\n", encoding="utf-8")
        instance = make_instance()
        instance._GenerateEnvs__check_if_setting_file_exists(str(tmp_path))
        assert instance.path_to_toml == os.path.join(str(tmp_path), "params.toml")

    def test_raises_when_params_toml_missing(self, tmp_path):
        instance = make_instance()
        with pytest.raises(FileExistsError):
            instance._GenerateEnvs__check_if_setting_file_exists(str(tmp_path))


# ---------------------------------------------------------------------------
# __check_if_variables_file_exists / __read_toml_variables
# ---------------------------------------------------------------------------


class TestVariablesFile:
    def test_path_set_when_variables_toml_present(self, tmp_path):
        (tmp_path / "params.variables.toml").write_text('X="y"\n', encoding="utf-8")
        instance = make_instance()
        instance._GenerateEnvs__check_if_variables_file_exists(str(tmp_path))
        assert instance.path_to_variables_toml == os.path.join(
            str(tmp_path), "params.variables.toml"
        )

    def test_path_stays_empty_when_variables_toml_absent(self, tmp_path):
        instance = make_instance()
        instance._GenerateEnvs__check_if_variables_file_exists(str(tmp_path))
        assert instance.path_to_variables_toml == ""

    def test_read_toml_variables_returns_empty_dict_without_file(self):
        instance = make_instance(path_to_variables_toml="")
        assert instance._GenerateEnvs__read_toml_variables() == {}

    def test_read_toml_variables_parses_existing_file(self, tmp_path):
        variables_file = tmp_path / "params.variables.toml"
        variables_file.write_text('REDIS_PORT=6379\nHOST="localhost"\n', encoding="utf-8")
        instance = make_instance(path_to_variables_toml=str(variables_file))
        assert instance._GenerateEnvs__read_toml_variables() == {
            "REDIS_PORT": 6379,
            "HOST": "localhost",
        }


# ---------------------------------------------------------------------------
# __write_file / __write_env_to_lib
# ---------------------------------------------------------------------------


class TestWriteFile:
    def test_writes_env_file_when_directory_exists(self, tmp_path):
        target_dir = tmp_path / "service"
        target_dir.mkdir()
        instance = make_instance()
        instance._GenerateEnvs__write_file(str(target_dir), "KEY=value\n")
        assert (target_dir / ".env").read_text(encoding="utf-8") == "KEY=value\n"

    def test_skips_silently_when_directory_missing(self, tmp_path):
        missing_dir = tmp_path / "nope"
        instance = make_instance()
        instance._GenerateEnvs__write_file(str(missing_dir), "KEY=value\n")
        assert not missing_dir.exists()

    def test_overwrites_existing_env_file(self, tmp_path):
        target_dir = tmp_path / "service"
        target_dir.mkdir()
        (target_dir / ".env").write_text("OLD=1\n", encoding="utf-8")
        instance = make_instance()
        instance._GenerateEnvs__write_file(str(target_dir), "NEW=2\n")
        assert (target_dir / ".env").read_text(encoding="utf-8") == "NEW=2\n"


# ---------------------------------------------------------------------------
# __remove_main_env_file
# ---------------------------------------------------------------------------


class TestRemoveMainEnvFile:
    def test_deletes_existing_env_file(self, tmp_path):
        (tmp_path / ".env").write_text("OLD=1\n", encoding="utf-8")
        instance = make_instance(path_to_env_folder=str(tmp_path))
        instance._GenerateEnvs__remove_main_env_file()
        assert not (tmp_path / ".env").exists()

    def test_noop_when_env_file_absent(self, tmp_path):
        instance = make_instance(path_to_env_folder=str(tmp_path))
        instance._GenerateEnvs__remove_main_env_file()  # must not raise


# ---------------------------------------------------------------------------
# __check_arguments_and_parse
# ---------------------------------------------------------------------------


class TestCheckArgumentsAndParse:
    def test_sets_path_from_first_positional_argument(self, monkeypatch, tmp_path):
        monkeypatch.setattr("sys.argv", ["generate_env.py", str(tmp_path)])
        instance = make_instance()
        instance._GenerateEnvs__check_arguments_and_parse()
        assert instance.path_to_env_folder == str(tmp_path)

    def test_missing_argument_exits_via_argparse(self, monkeypatch):
        monkeypatch.setattr("sys.argv", ["generate_env.py"])
        instance = make_instance()
        with pytest.raises(SystemExit):
            instance._GenerateEnvs__check_arguments_and_parse()

    def test_empty_string_argument_raises_value_error(self, monkeypatch):
        monkeypatch.setattr("sys.argv", ["generate_env.py", ""])
        instance = make_instance()
        with pytest.raises(ValueError):
            instance._GenerateEnvs__check_arguments_and_parse()


# ---------------------------------------------------------------------------
# End-to-end: run_all()
# ---------------------------------------------------------------------------


class TestRunAllIntegration:
    def _build_env_folder(self, tmp_path, params_toml: str, variables_toml: str | None = None):
        repo_root = tmp_path / "repo"
        env_folder = repo_root / "deploy" / "env"
        env_folder.mkdir(parents=True)
        (env_folder / "params.toml").write_text(params_toml, encoding="utf-8")
        if variables_toml is not None:
            (env_folder / "params.variables.toml").write_text(
                variables_toml, encoding="utf-8"
            )
        return repo_root, env_folder

    def test_generates_per_service_env_with_variable_substitution(self, tmp_path):
        params_toml = """
datacatalogue.host_and_port.DATACATALOGUE_GRPC_HOST="${GRPC_HOST}"
datacatalogue.host_and_port.DATACATALOGUE_GRPC_PORT=5052

redis.db.REDIS_PORT="${REDIS_PORT}"
redis.db.REDIS_PASSWORD="${REDIS_PASSWORD}"
"""
        variables_toml = """
GRPC_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD="secret"
"""
        repo_root, env_folder = self._build_env_folder(
            tmp_path, params_toml, variables_toml
        )
        # Service directories that __write_file expects two levels above env_folder.
        (repo_root / "datacatalogue").mkdir()
        (repo_root / "deploy" / "compose" / "redis").mkdir(parents=True)

        instance = make_instance(path_to_env_folder=str(env_folder))
        instance._GenerateEnvs__check_if_setting_file_exists(str(env_folder))
        instance._GenerateEnvs__check_if_variables_file_exists(str(env_folder))
        instance.run_all()

        datacatalogue_env = (repo_root / "datacatalogue" / ".env").read_text(
            encoding="utf-8"
        )
        assert "DATACATALOGUE_GRPC_HOST=localhost\n" in datacatalogue_env
        assert "DATACATALOGUE_GRPC_PORT=5052\n" in datacatalogue_env

        redis_env = (
            repo_root / "deploy" / "compose" / "redis" / ".env"
        ).read_text(encoding="utf-8")
        assert "REDIS_PORT=6379\n" in redis_env
        assert "REDIS_PASSWORD=secret\n" in redis_env

    def test_missing_variable_reference_raises_keyerror(self, tmp_path):
        params_toml = 'lib.section.HOST="${UNDEFINED}"\n'
        repo_root, env_folder = self._build_env_folder(tmp_path, params_toml)

        instance = make_instance(path_to_env_folder=str(env_folder))
        instance._GenerateEnvs__check_if_setting_file_exists(str(env_folder))
        instance._GenerateEnvs__check_if_variables_file_exists(str(env_folder))

        with pytest.raises(KeyError):
            instance.run_all()

    def test_works_without_variables_file(self, tmp_path):
        params_toml = 'lib.section.HOST="localhost"\n'
        repo_root, env_folder = self._build_env_folder(tmp_path, params_toml)

        instance = make_instance(path_to_env_folder=str(env_folder))
        instance._GenerateEnvs__check_if_setting_file_exists(str(env_folder))
        instance._GenerateEnvs__check_if_variables_file_exists(str(env_folder))
        instance.run_all()  # must not raise despite no params.variables.toml

        aggregate_env = (env_folder / ".env").read_text(encoding="utf-8")
        assert "HOST=localhost" in aggregate_env

    def test_removes_stale_aggregate_env_before_rebuilding(self, tmp_path):
        params_toml = 'lib.section.HOST="localhost"\n'
        repo_root, env_folder = self._build_env_folder(tmp_path, params_toml)
        (env_folder / ".env").write_text("STALE=1\n", encoding="utf-8")

        instance = make_instance(path_to_env_folder=str(env_folder))
        instance._GenerateEnvs__check_if_setting_file_exists(str(env_folder))
        instance._GenerateEnvs__check_if_variables_file_exists(str(env_folder))
        instance.run_all()

        aggregate_env = (env_folder / ".env").read_text(encoding="utf-8")
        assert "STALE" not in aggregate_env
