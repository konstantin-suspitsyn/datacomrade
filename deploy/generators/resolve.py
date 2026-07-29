# -*- coding: utf-8 -*-
"""Шаг 2 конвейера: сборка модели сущностей из build/parsed.json.

Любое расхождение между query.sql, sqlc и proto — ошибка, а не повод
догадываться: конвейер обязан упасть здесь, а не выдать кривой Go.
Результат — build/model.json.

Запуск:  python deploy/generators/resolve.py [--config ...] [--work-dir ...]
"""
import re

from genconfig import parse_args, read_json, write_json

def norm(name):
    return name.replace("_", "").lower()


def snake(name):
    """PascalCase -> snake_case для имён файлов."""
    s = re.sub(r"(.)([A-Z][a-z]+)", r"\1_\2", name)
    return re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", s).lower()


class Fail(Exception):
    pass


def resolve_entity(sqlc, data, block):
    table = block["table"]  # dc.host
    short = table.split(".")[-1]
    names = [q[0] for q in block["queries"]]
    kinds = dict(block["queries"])

    # Сущность — из Create<E>, множественное число — из Get<P>.
    creates = [n for n in names if n.startswith("Create")]
    if len(creates) != 1:
        raise Fail("%s: ожидался ровно один Create, найдено %s" % (table, creates))
    entity = creates[0][len("Create"):]

    plurals = [
        n[len("Get"):]
        for n in names
        if n.startswith("Get")
        and not n.startswith("GetDeleted")
        and not n.endswith("ById")
        and "By" not in n[len("Get"):]
    ]
    if len(plurals) != 1:
        raise Fail("%s: не удалось определить множественное число: %s" % (table, plurals))
    plural = plurals[0]

    fk_queries = []
    for n in names:
        m = re.match(r"^Get(%s)By(\w+)$" % re.escape(plural), n)
        if m and not n.startswith("GetDeleted"):
            fk_queries.append({"query": n, "by": m.group(2)})

    # Выборка одной строки по уникальной колонке: Get<E>By<Column> с :one.
    # От FK-выборки отличается тем, что колонка уникальна, поэтому запрос
    # возвращает строку, а не список.
    lookup_queries = []
    for n in names:
        m = re.match(r"^Get(%s)By(\w+)$" % re.escape(entity), n)
        if not m or n == "Get%sById" % entity or n.startswith("GetDeleted"):
            continue
        if kinds[n] != "one":
            raise Fail(
                "%s: %s отбирает по уникальной колонке и обязан быть :one, а не :%s"
                % (table, n, kinds[n])
            )
        lookup_queries.append({"query": n, "by": m.group(2)})

    expected = (
        [
            "Get%sById" % entity,
            "Get%s" % plural,
            "GetDeleted%sById" % entity,
            "GetDeleted%s" % plural,
        ]
        + [f["query"] for f in fk_queries]
        + [l["query"] for l in lookup_queries]
        + [
            "Create%s" % entity,
            "Update%sById" % entity,
            "Delete%sById" % entity,
            "Undelete%sById" % entity,
        ]
    )
    if sorted(expected) != sorted(names):
        raise Fail("%s: набор запросов расходится\n  ожидали %s\n  нашли   %s" % (table, sorted(expected), sorted(names)))

    funcs = data["funcs"]
    for n in names:
        if n not in funcs:
            raise Fail("%s: нет метода репозитория %s" % (table, n))

    # Структура строки — из возвращаемого типа Get<E>ById.
    ret = funcs["Get%sById" % entity]["ret"]
    row_struct = ret.split(",")[0].strip()
    if row_struct not in data["model_structs"]:
        raise Fail("%s: нет модели %s" % (table, row_struct))

    create_params = "Create%sParams" % entity
    update_params = "Update%sByIdParams" % entity
    for p in (create_params, update_params):
        if p not in data["param_structs"]:
            raise Fail("%s: нет структуры %s" % (table, p))

    # Сообщения proto.
    proto = data["proto_structs"]
    for msg in [entity, "Create%sRequest" % entity, "Update%sByIdRequest" % entity]:
        if msg not in proto:
            raise Fail("%s: нет proto-сообщения %s" % (table, msg))

    schema_cols = data["schema"][table]

    def pair_fields(go_fields, proto_msg, what):
        """Сопоставляет поля Go-структуры с полями proto-сообщения."""
        pfields = {norm(f["field"]): f for f in proto[proto_msg]}
        paired = []
        for f in go_fields:
            key = norm(f["field"])
            if key not in pfields:
                raise Fail("%s: поле %s.%s не найдено в %s" % (table, what, f["field"], proto_msg))
            pf = pfields[key]
            col = pf["proto_name"]
            if col is None:
                raise Fail("%s: у поля %s.%s нет имени колонки" % (table, proto_msg, pf["field"]))
            if col not in schema_cols:
                raise Fail("%s: колонки %s нет в schema.sql" % (table, col))
            paired.append(
                {
                    "go": f["field"],
                    "go_type": f["type"],
                    "proto": pf["field"],
                    "proto_type": pf["type"],
                    "col": col,
                    "varchar": schema_cols[col]["varchar"],
                }
            )
        return paired

    row_fields = pair_fields(data["model_structs"][row_struct], entity, row_struct)
    create_fields = pair_fields(data["param_structs"][create_params], "Create%sRequest" % entity, create_params)
    update_fields = pair_fields(
        data["param_structs"][update_params], "Update%sByIdRequest" % entity, update_params
    )

    # Поля ответов: Get<E>ByIdResponse{<E> *<E>}, Get<P>Response{<P> []*<E>}.
    def response_field(msg):
        fields = proto.get(msg)
        if not fields:
            raise Fail("%s: нет proto-сообщения %s" % (table, msg))
        if len(fields) != 1:
            raise Fail("%s: в %s ожидалось одно поле, найдено %s" % (table, msg, fields))
        return fields[0]["field"]

    resp_fields = {
        "Get%sById" % entity: response_field("Get%sByIdResponse" % entity),
        "Get%s" % plural: response_field("Get%sResponse" % plural),
        "GetDeleted%sById" % entity: response_field("GetDeleted%sByIdResponse" % entity),
        "GetDeleted%s" % plural: response_field("GetDeleted%sResponse" % plural),
        "Create%s" % entity: response_field("Create%sResponse" % entity),
        "Update%sById" % entity: response_field("Update%sByIdResponse" % entity),
        "Delete%sById" % entity: response_field("Delete%sByIdResponse" % entity),
        "Undelete%sById" % entity: response_field("Undelete%sByIdResponse" % entity),
    }

    fks = []
    for fk in fk_queries:
        args = funcs[fk["query"]]["args"]
        if len(args) != 1:
            raise Fail("%s: у %s ожидался один аргумент, найдено %s" % (table, fk["query"], args))
        arg = args[0]
        req_msg = fk["query"] + "Request"
        pf = proto.get(req_msg)
        if not pf or len(pf) != 1:
            raise Fail("%s: в %s ожидалось одно поле" % (table, req_msg))
        col = pf[0]["proto_name"]
        fks.append(
            {
                "query": fk["query"],
                "arg_name": arg["name"],
                "arg_type": arg["type"],
                "proto_field": pf[0]["field"],
                "col": col,
                "resp_field": response_field(fk["query"] + "Response"),
            }
        )

    lookups = []
    for lk in lookup_queries:
        args = funcs[lk["query"]]["args"]
        if len(args) != 1:
            raise Fail("%s: у %s ожидался один аргумент, найдено %s" % (table, lk["query"], args))
        arg = args[0]
        req_msg = lk["query"] + "Request"
        pf = proto.get(req_msg)
        if not pf or len(pf) != 1:
            raise Fail("%s: в %s ожидалось одно поле" % (table, req_msg))
        col = pf[0]["proto_name"]
        if col not in schema_cols:
            raise Fail("%s: колонки %s нет в schema.sql" % (table, col))
        lookups.append(
            {
                "query": lk["query"],
                "arg_name": arg["name"],
                "arg_type": arg["type"],
                "proto_field": pf[0]["field"],
                "proto_type": pf[0]["type"],
                "col": col,
                "varchar": schema_cols[col]["varchar"],
                "resp_field": response_field(lk["query"] + "Response"),
            }
        )

    return {
        "sqlc": sqlc,
        "table": table,
        "short": short,
        "file": snake(entity),
        "entity": entity,
        "plural": plural,
        "row_struct": row_struct,
        "create_params": create_params,
        "update_params": update_params,
        "row_fields": row_fields,
        "create_fields": create_fields,
        "update_fields": update_fields,
        "resp_fields": resp_fields,
        "fks": fks,
        "lookups": lookups,
        "kinds": kinds,
    }


def main():
    args, config = parse_args("Сборка модели сущностей из parsed.json")

    parsed = read_json(args.work_dir, "parsed.json")
    meta_by_sqlc = {d["sqlc"]: d for d in config["domains"]}

    model = {}
    total_rpc = 0
    for sqlc, data in parsed.items():
        if sqlc not in meta_by_sqlc:
            raise SystemExit("домена %r нет в конфиге — перезапустите parse.py" % sqlc)

        entities = []
        for block in data["blocks"]:
            entities.append(resolve_entity(sqlc, data, block))

        model[sqlc] = {"meta": meta_by_sqlc[sqlc], "entities": entities}
        for e in entities:
            total_rpc += 8 + len(e["fks"]) + len(e["lookups"])

    write_json(args.work_dir, "model.json", model)

    print("всего RPC:", total_rpc)
    for sqlc, d in model.items():
        print("\n%s (%d таблиц)" % (sqlc, len(d["entities"])))
        for e in d["entities"]:
            fk = " fk=%s" % ",".join(f["query"] for f in e["fks"]) if e["fks"] else ""
            lk = " lookup=%s" % ",".join(l["query"] for l in e["lookups"]) if e["lookups"] else ""
            print(
                "  %-24s %-22s file=%-22s%s%s"
                % (e["table"], e["entity"] + "/" + e["plural"], e["file"] + ".go", fk, lk)
            )

    # Типы, которые встречаются в полях — список должен быть закрытым.
    types = set()
    for d in model.values():
        for e in d["entities"]:
            for group in ("row_fields", "create_fields", "update_fields"):
                for f in e[group]:
                    types.add((f["go_type"], f["proto_type"]))
    print("\nпары типов sqlc -> proto:")
    for t in sorted(types):
        print("  %-16s -> %s" % t)


main()
