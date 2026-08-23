# -*- coding: utf-8 -*-
"""Собирает REST-контракт Shepherd из api/openapi/v1/spec/*.yaml в один openapi.yaml.

Контракт разрезан по пакетам микросервиса datacatalogue (см. root.yaml:
x-spec-parts). Каждая часть — самостоятельный валидный OpenAPI-документ, поэтому
её можно линтовать и открывать в редакторе отдельно; ссылки между частями пишутся
как ./common.yaml#/components/schemas/Error.

Собирать нужно потому, что oapi-codegen (v2) генерирует модели только из
components корневого документа: внешние $ref он либо отказывается резолвить,
либо молча выкидывает типы. Поэтому источник истины — spec/*.yaml, а
api/openapi/v1/openapi.yaml — сборка, которую нельзя править руками.

Запуск:  python deploy/python-scripts/bundle_openapi.py [spec-dir] [out-file]
"""
import os
import sys

import yaml

DEFAULT_SPEC = os.path.join("shepherd", "api", "openapi", "v1", "spec")
DEFAULT_OUT = os.path.join("shepherd", "api", "openapi", "v1", "openapi.yaml")

HEADER = (
    "# СГЕНЕРИРОВАННЫЙ ФАЙЛ — НЕ РЕДАКТИРОВАТЬ.\n"
    "# Собран из api/openapi/v1/spec/*.yaml командой `task openapi:bundle`.\n"
    "# Правки вносить в spec/<пакет>.yaml, затем `task openapi:gen`.\n"
)

# Секции components, которые склеиваем. Порядок — как в собранном файле.
SECTIONS = ("securitySchemes", "parameters", "responses", "schemas")


class Dumper(yaml.SafeDumper):
    """Без якорей/алиасов: одинаковые куски должны разворачиваться в текст."""

    def ignore_aliases(self, data):
        return True

    def increase_indent(self, flow=False, indentless=False):
        return super(Dumper, self).increase_indent(flow, False)


def _str(dumper, data):
    """Длинные описания — свёрнутым блоком `>`, многострочные — литеральным `|`."""
    body = data[:-1] if data.endswith("\n") else data
    if "\n" in body:
        style = "|"
    elif data.endswith("\n") and len(data) > 80:
        style = ">"
    else:
        style = None
    return dumper.represent_scalar("tag:yaml.org,2002:str", data, style=style)


Dumper.add_representer(str, _str)


def read(path):
    with open(path, encoding="utf-8") as f:
        return yaml.safe_load(f)


def flatten_refs(node, parts):
    """./tables.yaml#/components/schemas/Alias -> #/components/schemas/Alias."""
    if isinstance(node, dict):
        for key, value in node.items():
            if key == "$ref" and isinstance(value, str) and "#" in value and not value.startswith("#"):
                head, pointer = value.split("#", 1)
                name = os.path.basename(head)
                if name not in parts:
                    raise SystemExit("неизвестная часть контракта в $ref: %s" % value)
                node[key] = "#" + pointer
            else:
                flatten_refs(value, parts)
    elif isinstance(node, list):
        for value in node:
            flatten_refs(value, parts)
    return node


def merge(dst, src, section, origin):
    for name, value in src.items():
        if name in dst:
            raise SystemExit("дубль %s/%s: %s и %s" % (section, name, dst[name + "@from"], origin))
        dst[name] = value
    return dst


def bundle(spec_dir, out_file):
    root = read(os.path.join(spec_dir, "root.yaml"))
    parts = root.pop("x-spec-parts", None)
    if not parts:
        raise SystemExit("в root.yaml нет x-spec-parts со списком частей контракта")

    paths = dict(root.get("paths") or {})
    components = {s: dict((root.get("components") or {}).get(s) or {}) for s in SECTIONS}
    owner = {}

    for part_name in parts:
        part = read(os.path.join(spec_dir, part_name))
        for path, item in (part.get("paths") or {}).items():
            if path in owner:
                raise SystemExit("путь %s объявлен и в %s, и в %s" % (path, owner[path], part_name))
            owner[path] = part_name
            paths[path] = item
        for section in SECTIONS:
            for name, value in ((part.get("components") or {}).get(section) or {}).items():
                key = section + "/" + name
                if key in owner:
                    raise SystemExit("компонент %s объявлен и в %s, и в %s"
                                     % (key, owner[key], part_name))
                owner[key] = part_name
                components[section][name] = value

    doc = {}
    for key in ("openapi", "info", "servers", "security"):
        if key in root:
            doc[key] = root[key]
    doc["paths"] = paths
    doc["components"] = {s: components[s] for s in SECTIONS if components[s]}
    flatten_refs(doc, set(parts) | {"root.yaml"})

    with open(out_file, "w", encoding="utf-8", newline="\n") as f:
        f.write(HEADER)
        yaml.dump(doc, f, Dumper=Dumper, allow_unicode=True, sort_keys=False,
                  default_flow_style=False, width=110)

    print("[bundle] %s <- %s (%d путей, %d схем)"
          % (out_file, ", ".join(parts), len(paths), len(components["schemas"])))


def main():
    spec_dir = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_SPEC
    out_file = sys.argv[2] if len(sys.argv) > 2 else DEFAULT_OUT
    bundle(spec_dir, out_file)


if __name__ == "__main__":
    main()
