# -*- coding: utf-8 -*-
"""Шаг 1 конвейера для доменов без табличного CRUD (crud_config.json/custom_domains).

Обычный parse.py (см. parse.py) ожидает, что query.sql разбит на блоки по
таблицам и что у каждой таблицы есть Create/Update/Delete. Домены вроде
auth_logic — сквозные выборки поверх нескольких таблиц, без единой сущности —
под эту форму не подходят, поэтому у них отдельный конвейер:
parse_custom.py -> resolve_custom.py -> gen_custom.py.

Читает schema.sql (все таблицы сразу, без привязки к одной), query.sql
(плоский список запросов, без разбора по блокам), сгенерированные sqlc-модели
и proto pb.go. Складывает всё в build/parsed_custom.json. Ничего не додумывает —
только разбирает, как и parse.py.

Запуск:  python deploy/generators/parse_custom.py [--config ...] [--work-dir ...]
"""
import re

from genconfig import parse_args, read_text, repo_path, write_json


def read(path):
    return read_text(path)


# ---------- schema.sql: колонки и границы varchar по всем таблицам ----------
def parse_schema(path):
    text = read(path)
    tables = {}
    for m in re.finditer(r"CREATE TABLE (\S+)\s*\((.*?)\n\);", text, re.S | re.I):
        table = m.group(1).replace('"', "")
        cols = {}
        for line in m.group(2).split("\n"):
            line = line.strip().rstrip(",")
            if not line:
                continue
            cm = re.match(r'"?([a-z_]+)"?\s+(.*)', line)
            if not cm:
                continue
            col, coltype = cm.group(1), cm.group(2)
            # auth_logic/schema.sql пишет "varchar(n)", а не "character varying(n)"
            # (как в остальных schema.sql) — оба означают один тип в Postgres.
            vm = re.search(r"(?:character varying|varchar)\((\d+)\)", coltype)
            cols[col] = {
                "type": coltype,
                "varchar": int(vm.group(1)) if vm else None,
                "nullable": "not null" not in coltype.lower(),
            }
        tables[table] = cols
    return tables


# ---------- query.sql: плоский список запросов, порядок как в файле ----------
def parse_queries_flat(path):
    text = read(path)
    return [list(m) for m in re.findall(r"^-- name: (\w+) :(\w+)", text, re.M)]


# ---------- Go: структуры и сигнатуры (те же правила, что в parse.py) ----------
def parse_go_structs(text):
    """Возвращает {StructName: [(FieldName, GoType, protoName|None)]}."""
    structs = {}
    for m in re.finditer(r"^type (\w+) struct \{\n(.*?)^\}", text, re.S | re.M):
        name = m.group(1)
        fields = []
        for line in m.group(2).split("\n"):
            line = line.rstrip()
            fm = re.match(r"\t([A-Z]\w*)\s+([\w\.\*\[\]]+)(?:\s+`(.*)`)?$", line)
            if not fm:
                continue
            field, gotype, tag = fm.group(1), fm.group(2), fm.group(3)
            proto_name = None
            if tag:
                pm = re.search(r"name=([a-z_0-9]+)", tag)
                if pm:
                    proto_name = pm.group(1)
            if field in ("state", "unknownFields", "sizeCache"):
                continue
            fields.append({"field": field, "type": gotype, "proto_name": proto_name})
        structs[name] = fields
    return structs


def parse_go_funcs(text):
    """Сигнатуры методов Queries: {Name: {args, ret}}."""
    funcs = {}
    for m in re.finditer(
        r"^func \(q \*Queries\) (\w+)\(ctx context\.Context(.*?)\) \((.*?)\)|"
        r"^func \(q \*Queries\) (\w+)\(ctx context\.Context(.*?)\) (error)",
        text,
        re.M,
    ):
        if m.group(1):
            name, args, ret = m.group(1), m.group(2), m.group(3)
        else:
            name, args, ret = m.group(4), m.group(5), m.group(6)
        arg_list = []
        for a in args.split(",")[1:] if args.startswith(",") else []:
            a = a.strip()
            if not a:
                continue
            an, at = a.split(" ", 1)
            arg_list.append({"name": an, "type": at})
        funcs[name] = {"args": arg_list, "ret": ret.strip()}
    return funcs


def main():
    args, config = parse_args("Разбор источников для доменов без табличного CRUD (custom_domains)")

    custom_domains = config.get("custom_domains") or []
    if not custom_domains:
        write_json(args.work_dir, "parsed_custom.json", {})
        print("custom_domains пуст в конфиге — писать нечего")
        return

    sqlc_root = config["sqlc_root"]
    internal = config["internal"]
    proto_go_root = config["proto_go_root"]

    out = {}
    for d in custom_domains:
        sqlc = d["sqlc"]
        repo_dir = "%s/repository/%s" % (internal, sqlc)

        schema = parse_schema(repo_path(sqlc_root, sqlc, "schema.sql"))
        queries = parse_queries_flat(repo_path(sqlc_root, sqlc, "query.sql"))
        queries_go = read(repo_path(repo_dir, "query.sql.go"))
        pb = read(repo_path(proto_go_root, d["proto_pkg"], d["proto_file"]))

        out[sqlc] = {
            "schema": schema,
            "queries": queries,
            "param_structs": parse_go_structs(queries_go),
            "funcs": parse_go_funcs(queries_go),
            "proto_structs": parse_go_structs(pb),
        }

    path = write_json(args.work_dir, "parsed_custom.json", out)

    for sqlc, data in out.items():
        print(sqlc)
        print("  таблиц в схеме:      ", len(data["schema"]))
        print("  запросов:            ", len(data["queries"]))
        print("  методов репозитория: ", len(data["funcs"]))
        print("  сообщений proto:     ", len(data["proto_structs"]))

    print("\nзаписано:", path)


if __name__ == "__main__":
    main()
