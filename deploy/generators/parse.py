# -*- coding: utf-8 -*-
"""Шаг 1 конвейера: разбор источников истины.

Читает schema.sql, query.sql, сгенерированные sqlc-модели и proto pb.go,
складывает всё в build/parsed.json. Ничего не додумывает — только разбирает.

Запуск:  python deploy/generators/parse.py [--config ...] [--work-dir ...]
"""
import re

from genconfig import parse_args, read_text, repo_path, write_json


def read(path):
    return read_text(path)


def norm(name):
    """Ключ для сопоставления полей sqlc и protoc: UserID и UserId -> userid."""
    return name.replace("_", "").lower()


# ---------- schema.sql: колонки и границы varchar ----------
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
            vm = re.search(r"character varying\((\d+)\)", coltype)
            cols[col] = {
                "type": coltype,
                "varchar": int(vm.group(1)) if vm else None,
                "nullable": "not null" not in coltype.lower(),
            }
        tables[table] = cols
    return tables


# ---------- query.sql: блоки по таблицам, имена запросов ----------
def parse_queries(path):
    text = read(path)
    blocks = []
    # Разделитель вида: -- ====\n-- dc.host\n-- ====
    parts = re.split(r"-- =+\n-- (\S+)\n-- =+", text)
    for i in range(1, len(parts), 2):
        table = parts[i]
        body = parts[i + 1]
        names = re.findall(r"^-- name: (\w+) :(\w+)", body, re.M)
        blocks.append({"table": table, "queries": names})
    return blocks


# ---------- Go: структуры и сигнатуры ----------
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
    args, config = parse_args("Разбор schema.sql, query.sql, sqlc-моделей и proto pb.go")

    sqlc_root = config["sqlc_root"]
    internal = config["internal"]
    proto_go_root = config["proto_go_root"]

    out = {}
    for d in config["domains"]:
        sqlc = d["sqlc"]
        repo_dir = "%s/repository/%s" % (internal, sqlc)

        schema = parse_schema(repo_path(sqlc_root, sqlc, "schema.sql"))
        blocks = parse_queries(repo_path(sqlc_root, sqlc, "query.sql"))
        models = read(repo_path(repo_dir, "models.go"))
        queries_go = read(repo_path(repo_dir, "query.sql.go"))
        pb = read(repo_path(proto_go_root, d["proto_pkg"], d["proto_file"]))

        out[sqlc] = {
            "schema": schema,
            "blocks": blocks,
            "model_structs": parse_go_structs(models),
            "param_structs": parse_go_structs(queries_go),
            "funcs": parse_go_funcs(queries_go),
            "proto_structs": parse_go_structs(pb),
        }

    path = write_json(args.work_dir, "parsed.json", out)

    for sqlc, data in out.items():
        print(sqlc)
        print("  таблиц в схеме:      ", len(data["schema"]))
        print("  блоков запросов:     ", len(data["blocks"]))
        print("  запросов всего:      ", sum(len(b["queries"]) for b in data["blocks"]))
        print("  структур моделей:    ", len(data["model_structs"]))
        print("  структур параметров: ", len([k for k in data["param_structs"] if k.endswith("Params")]))
        print("  методов репозитория: ", len(data["funcs"]))
        print("  сообщений proto:     ", len(data["proto_structs"]))

    print("\nзаписано:", path)


main()
