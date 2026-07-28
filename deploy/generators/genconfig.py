# -*- coding: utf-8 -*-
"""Общая настройка для parse.py, resolve.py и gen.py.

Корень репозитория вычисляется от расположения скрипта (deploy/generators/../..),
поэтому абсолютных путей в коде нет и репозиторий можно перенести куда угодно.
Все пути в crud_config.json задаются от корня репозитория через прямой слэш.
"""
import argparse
import io
import json
import os

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.normpath(os.path.join(HERE, "..", ".."))

DEFAULT_CONFIG = os.path.join(HERE, "crud_config.json")
DEFAULT_WORK_DIR = os.path.join(HERE, "build")


def repo_path(*parts):
    """Абсолютный путь от корня репозитория."""
    joined = "/".join(p for p in parts if p)
    return os.path.join(ROOT, joined.replace("/", os.sep))


def load_config(path):
    with io.open(path, encoding="utf-8") as f:
        config = json.loads(f.read())

    required = ["module", "internal", "sqlc_root", "proto_go_root", "shared_imports", "domains"]
    for key in required:
        if key not in config:
            raise SystemExit("в конфиге нет обязательного ключа %r: %s" % (key, path))

    if not config["domains"]:
        raise SystemExit("в конфиге пустой список domains: %s" % path)

    return config


def parse_args(description):
    parser = argparse.ArgumentParser(description=description)
    parser.add_argument(
        "--config",
        default=DEFAULT_CONFIG,
        help="путь к crud_config.json (по умолчанию рядом со скриптом)",
    )
    parser.add_argument(
        "--work-dir",
        default=DEFAULT_WORK_DIR,
        help="каталог для промежуточных parsed.json и model.json",
    )
    args = parser.parse_args()

    if not os.path.isfile(args.config):
        raise SystemExit("конфиг не найден: %s" % args.config)

    if not os.path.isdir(args.work_dir):
        os.makedirs(args.work_dir)

    return args, load_config(args.config)


def read_text(path):
    if not os.path.isfile(path):
        raise SystemExit("не найден файл источника: %s" % path)

    return io.open(path, encoding="utf-8").read()


def write_json(work_dir, name, data):
    path = os.path.join(work_dir, name)
    with io.open(path, "w", encoding="utf-8") as f:
        f.write(json.dumps(data, ensure_ascii=False, indent=1))

    return path


def read_json(work_dir, name):
    path = os.path.join(work_dir, name)
    if not os.path.isfile(path):
        raise SystemExit("не найден %s — сначала запустите предыдущий шаг конвейера" % path)

    return json.loads(io.open(path, encoding="utf-8").read())
