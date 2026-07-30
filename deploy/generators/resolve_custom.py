# -*- coding: utf-8 -*-
"""Шаг 2 конвейера для доменов без табличного CRUD: собирает модель по
build/parsed_custom.json. Пара к resolve.py, но единица сборки — один запрос
query.sql, а не одна таблица: нет Create<E>, нет строки-сущности, есть только
аргументы запроса и его возвращаемое значение.

Любое расхождение между query.sql, sqlc и proto — ошибка, а не повод
догадываться, как и в resolve.py. Результат — build/model_custom.json.

Запуск:  python deploy/generators/resolve_custom.py [--config ...] [--work-dir ...]
"""
from genconfig import parse_args, read_json, write_json


def norm(name):
    return name.replace("_", "").lower()


class Fail(Exception):
    pass


# Пары типов sqlc -> proto, которые конвертер умеет переносить.
# None — типы совпадают, переносится напрямую; иначе — %-шаблон выражения.
TYPE_RULES = {
    ("bool", "bool"): None,
    ("int64", "int64"): None,
    ("string", "string"): None,
    ("int16", "int32"): "int32(%s)",
    ("time.Time", "*timestamppb.Timestamp"): "converter.TimeToProto(%s)",
    ("sql.NullTime", "*timestamppb.Timestamp"): "converter.NullTimeToProto(%s)",
    ("uuid.UUID", "string"): "converter.UUIDToProto(%s)",
}


def find_column(schema, col, table_hint=None):
    """Ищет колонку col по таблицам домена (запросы здесь — джойны нескольких
    таблиц, колонка не привязана к одной). Если одноимённая колонка встречается
    в разных таблицах с разной границей varchar, выбор конкретной таблицы важен
    и его нельзя делать молча — тогда обязателен table_hint из column_hints
    в конфиге (какую таблицу и колонку реально фильтрует запрос)."""
    if table_hint is not None:
        if table_hint not in schema:
            raise Fail("column_hints: таблицы %s нет в schema.sql" % table_hint)
        if col not in schema[table_hint]:
            raise Fail("column_hints: колонки %s нет в таблице %s" % (col, table_hint))
        return schema[table_hint][col]

    found = [(t, c[col]) for t, c in schema.items() if col in c]
    if not found:
        return None
    varchars = {info["varchar"] for _, info in found}
    if len(varchars) > 1:
        raise Fail(
            "колонка %s встречается в нескольких таблицах с разной длиной varchar: %s — "
            "нужен column_hints.<query>.%s в конфиге, указывающий нужную таблицу"
            % (col, found, col)
        )
    return found[0][1]


def resolve_query(sqlc, data, name, kind, column_hints):
    funcs = data["funcs"]
    if name not in funcs:
        raise Fail("%s: нет метода репозитория %s" % (sqlc, name))
    func = funcs[name]

    # Пока поддержан только вид "один аргумент — Params-структура sqlc":
    # это всё, что реально встречается в custom_domains сейчас. Запрос с
    # другой формой аргументов должен упасть здесь, а не сгенерировать
    # что-то наугад.
    args = func["args"]
    if len(args) != 1 or args[0]["type"] not in data["param_structs"]:
        raise Fail(
            "%s: %s — аргументы %s не в виде одной Params-структуры, для такой формы правил ещё нет"
            % (sqlc, name, args)
        )
    params_type = args[0]["type"]
    param_fields = data["param_structs"][params_type]

    proto = data["proto_structs"]
    req_msg = name + "Request"
    resp_msg = name + "Response"
    if req_msg not in proto:
        raise Fail("%s: нет proto-сообщения %s" % (sqlc, req_msg))
    if resp_msg not in proto:
        raise Fail("%s: нет proto-сообщения %s" % (sqlc, resp_msg))

    pfields = {norm(f["field"]): f for f in proto[req_msg]}
    if len(pfields) != len(param_fields):
        raise Fail(
            "%s: число полей %s (%d) не совпадает с числом полей %s (%d)"
            % (sqlc, params_type, len(param_fields), req_msg, len(pfields))
        )

    req_fields = []
    for f in param_fields:
        key = norm(f["field"])
        if key not in pfields:
            raise Fail("%s: поле %s.%s не найдено в %s" % (sqlc, params_type, f["field"], req_msg))
        pf = pfields[key]
        col = pf["proto_name"]
        if col is None:
            raise Fail("%s: у поля %s.%s нет имени колонки" % (sqlc, req_msg, pf["field"]))
        colinfo = find_column(data["schema"], col, column_hints.get(col))
        if colinfo is None:
            raise Fail("%s: колонки %s нет в schema.sql" % (sqlc, col))
        pair = (f["type"], pf["type"])
        if pair not in TYPE_RULES:
            raise Fail("%s: нет правила для %s -> %s (поле %s)" % (sqlc, f["type"], pf["type"], col))
        req_fields.append(
            {
                "go": f["field"],
                "go_type": f["type"],
                "proto": pf["field"],
                "proto_type": pf["type"],
                "col": col,
                "varchar": colinfo["varchar"],
            }
        )

    if len(proto[resp_msg]) != 1:
        raise Fail("%s: в %s ожидалось одно поле, найдено %s" % (sqlc, resp_msg, proto[resp_msg]))
    resp_field = proto[resp_msg][0]

    ret_type = func["ret"].split(",")[0].strip()
    is_slice = ret_type.startswith("[]")
    go_elem = ret_type[2:] if is_slice else ret_type

    proto_type = resp_field["type"]
    proto_is_slice = proto_type.startswith("[]")
    proto_elem = proto_type[2:] if proto_is_slice else proto_type

    if is_slice != proto_is_slice:
        raise Fail(
            "%s: %s возвращает %s, а %s.%s — %s (repeated не совпадает)"
            % (sqlc, name, ret_type, resp_msg, resp_field["field"], proto_type)
        )

    elem_pair = (go_elem, proto_elem)
    if elem_pair not in TYPE_RULES:
        raise Fail("%s: нет правила для %s -> %s (ответ %s)" % (sqlc, go_elem, proto_elem, resp_msg))

    return {
        "sqlc": sqlc,
        "query": name,
        "kind": kind,
        "params_type": params_type,
        "req_fields": req_fields,
        "is_slice": is_slice,
        "go_elem": go_elem,
        "proto_elem": proto_elem,
        "elem_conv": TYPE_RULES[elem_pair],
        "resp_field": resp_field["field"],
    }


def main():
    args, config = parse_args("Сборка модели запросов для custom_domains из parsed_custom.json")

    parsed = read_json(args.work_dir, "parsed_custom.json")
    if not parsed:
        write_json(args.work_dir, "model_custom.json", {})
        print("нечего собирать (custom_domains пуст)")
        return

    meta_by_sqlc = {d["sqlc"]: d for d in config.get("custom_domains") or []}

    model = {}
    total = 0
    for sqlc, data in parsed.items():
        if sqlc not in meta_by_sqlc:
            raise SystemExit("домена %r нет в custom_domains конфига — перезапустите parse_custom.py" % sqlc)

        hints = meta_by_sqlc[sqlc].get("column_hints") or {}
        calls = [resolve_query(sqlc, data, name, kind, hints.get(name) or {}) for name, kind in data["queries"]]

        model[sqlc] = {"meta": meta_by_sqlc[sqlc], "calls": calls}
        total += len(calls)

    write_json(args.work_dir, "model_custom.json", model)

    print("всего запросов:", total)
    for sqlc, d in model.items():
        print("\n%s (%d запросов)" % (sqlc, len(d["calls"])))
        for c in d["calls"]:
            print("  %-45s :%-6s args=%d" % (c["query"], c["kind"], len(c["req_fields"])))


if __name__ == "__main__":
    main()
