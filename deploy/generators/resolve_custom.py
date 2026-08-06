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
    ("int32", "int32"): None,
    ("int64", "int64"): None,
    ("string", "string"): None,
    ("int16", "int32"): "int32(%s)",
    ("time.Time", "*timestamppb.Timestamp"): "converter.TimeToProto(%s)",
    ("sql.NullTime", "*timestamppb.Timestamp"): "converter.NullTimeToProto(%s)",
    ("uuid.UUID", "string"): "converter.UUIDToProto(%s)",
    # internal id как optional-поле: внешние клиенты полагаются на external_id,
    # см. datacatalogue/internal/converter/convert.go:Int64ToProto.
    ("int64", "*int64"): "converter.Int64ToProto(%s)",
}


def find_column(schema, col, table_hint=None, strict=True):
    """strict=False — подсказка необязательная: колонки нет в этой таблице,
    ищем по всей схеме. Так ведёт себя таблица блока query.sql: параметр
    запроса может прийти из подзапроса по соседней таблице (external_id
    приезжает из dc."user", а блок — dc.alias)."""
    if table_hint is not None and not strict:
        if table_hint in schema and col in schema[table_hint]:
            return schema[table_hint][col]
        table_hint = None

    return _find_column_strict(schema, col, table_hint)


def _find_column_strict(schema, col, table_hint=None):
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


# Постраничный ответ: строки плюс два счётчика. Имена полей задаёт SG Buddy,
# и опознаём мы страницу именно по ним, а не по имени запроса.
PAGE_ROWS = "rows"
PAGE_TOTAL_ROWS = "total_rows"
PAGE_TOTAL_PAGES = "total_pages"
COUNT_PREFIX = "Count"


def counter_name(query):
    """Имя парного счётчика для постраничной выборки."""
    return COUNT_PREFIX + query


def resolve_args(sqlc, data, name, func):
    """Форма аргументов запроса: без них, один скаляр или Params-структура."""
    args = func["args"]
    if not args:
        return {"arg_form": "none", "params_type": None, "arg_name": None, "arg_type": None}, []

    if len(args) != 1:
        raise Fail("%s: %s — ожидался один аргумент, найдено %s" % (sqlc, name, args))

    arg = args[0]
    if arg["type"] in data["param_structs"]:
        return (
            {"arg_form": "params", "params_type": arg["type"], "arg_name": None, "arg_type": None},
            data["param_structs"][arg["type"]],
        )

    # Скаляр: sqlc так делает, когда параметр ровно один. Поле запроса всё
    # равно одно, поэтому дальше он обрабатывается тем же кодом.
    return (
        {"arg_form": "scalar", "params_type": None, "arg_name": arg["name"], "arg_type": arg["type"]},
        [{"field": arg["name"][:1].upper() + arg["name"][1:], "type": arg["type"], "proto_name": None}],
    )


def pair_fields(sqlc, data, name, go_fields, proto_msg, column_hints, what, table=None):
    """Сопоставляет поля Go с полями proto-сообщения по нормализованному имени."""
    proto = data["proto_structs"]
    pfields = {norm(f["field"]): f for f in proto[proto_msg]}

    if len(pfields) != len(go_fields):
        raise Fail(
            "%s: число полей %s (%d) не совпадает с числом полей %s (%d)"
            % (sqlc, what, len(go_fields), proto_msg, len(pfields))
        )

    paired = []
    for f in go_fields:
        key = norm(f["field"])
        if key not in pfields:
            raise Fail("%s: поле %s.%s не найдено в %s" % (sqlc, what, f["field"], proto_msg))
        pf = pfields[key]
        col = pf["proto_name"]
        if col is None:
            raise Fail("%s: у поля %s.%s нет имени колонки" % (sqlc, proto_msg, pf["field"]))

        # Параметры страницы колонками не являются — проверять их по схеме нечем.
        colinfo = None
        if col not in (PAGE_LIMIT, PAGE_OFFSET):
            # Таблица блока отвечает на вопрос «чья это колонка» в обычном
            # случае; column_hints из конфига сильнее и обязан быть точным.
            explicit = column_hints.get(col)
            colinfo = find_column(
                data["schema"], col, explicit or table, strict=explicit is not None
            )
            if colinfo is None:
                raise Fail("%s: колонки %s нет в schema.sql (запрос %s)" % (sqlc, col, name))

        pair = (f["type"], pf["type"])
        if pair not in TYPE_RULES:
            raise Fail(
                "%s: нет правила для %s -> %s (поле %s, запрос %s)"
                % (sqlc, f["type"], pf["type"], col, name)
            )

        paired.append(
            {
                "go": f["field"],
                "go_type": f["type"],
                "proto": pf["field"],
                "proto_type": pf["type"],
                "col": col,
                "varchar": colinfo["varchar"] if colinfo else None,
                "page_param": col in (PAGE_LIMIT, PAGE_OFFSET),
            }
        )
    return paired


PAGE_LIMIT = "page_limit"
PAGE_OFFSET = "page_offset"


def resolve_response(sqlc, data, name, func, resp_msg, counters):
    """Форма ответа: пустой (:exec), одна величина или страница со счётчиками."""
    proto = data["proto_structs"]
    fields = proto[resp_msg]
    ret_type = func["ret"].split(",")[0].strip()

    if not fields:
        if ret_type != "error":
            raise Fail("%s: %s ничего не возвращает в proto, а репозиторий отдаёт %s" % (sqlc, name, ret_type))
        return {"resp_form": "empty"}

    is_slice = ret_type.startswith("[]")
    go_elem = ret_type[2:] if is_slice else ret_type

    names = [f["proto_name"] for f in fields]
    if len(fields) == 3 and names == [PAGE_ROWS, PAGE_TOTAL_ROWS, PAGE_TOTAL_PAGES]:
        if not is_slice:
            raise Fail("%s: %s отдаёт страницу, а репозиторий возвращает %s" % (sqlc, name, ret_type))
        counter = counter_name(name)
        if counter not in data["funcs"]:
            raise Fail(
                "%s: у постраничной выборки %s нет парного счётчика %s" % (sqlc, name, counter)
            )
        counters.add(counter)
        return {
            "resp_form": "page",
            "is_slice": True,
            "go_elem": go_elem,
            "row_msg": fields[0]["type"].lstrip("[]*"),
            "resp_field": fields[0]["field"],
            "counter": counter,
            "counter_row": data["funcs"][counter]["ret"].split(",")[0].strip(),
        }

    if len(fields) != 1:
        raise Fail("%s: в %s ожидалось одно поле или страница, найдено %s" % (sqlc, resp_msg, names))

    resp_field = fields[0]
    proto_type = resp_field["type"]
    proto_is_slice = proto_type.startswith("[]")
    proto_elem = proto_type[2:] if proto_is_slice else proto_type

    if is_slice != proto_is_slice:
        raise Fail(
            "%s: %s возвращает %s, а %s.%s — %s (repeated не совпадает)"
            % (sqlc, name, ret_type, resp_msg, resp_field["field"], proto_type)
        )

    # Строка таблицы: тип sqlc есть в models.go, а поле ответа — сообщение.
    if go_elem in data["model_structs"]:
        return {
            "resp_form": "row",
            "is_slice": is_slice,
            "go_elem": go_elem,
            "row_msg": proto_elem.lstrip("*"),
            "resp_field": resp_field["field"],
        }

    elem_pair = (go_elem, proto_elem.lstrip("*"))
    if elem_pair not in TYPE_RULES:
        raise Fail("%s: нет правила для %s -> %s (ответ %s)" % (sqlc, go_elem, proto_elem, resp_msg))

    return {
        "resp_form": "scalar",
        "is_slice": is_slice,
        "go_elem": go_elem,
        "proto_elem": proto_elem,
        "elem_conv": TYPE_RULES[elem_pair],
        "resp_field": resp_field["field"],
    }


def resolve_rows(sqlc, data, row_structs):
    """Пары «строка sqlc -> сообщение proto» для запросов, отдающих строки."""
    rows = []
    for go_struct, msg, table in sorted(row_structs):
        if go_struct not in data["model_structs"]:
            raise Fail("%s: нет модели %s" % (sqlc, go_struct))
        if msg not in data["proto_structs"]:
            raise Fail("%s: нет proto-сообщения %s" % (sqlc, msg))
        rows.append(
            {
                "go_struct": go_struct,
                "msg": msg,
                "fields": pair_fields(
                    sqlc,
                    data,
                    go_struct,
                    data["model_structs"][go_struct],
                    msg,
                    {},
                    go_struct,
                    table,
                ),
            }
        )
    return rows


def resolve_query(sqlc, data, name, kind, table, column_hints, counters):
    funcs = data["funcs"]
    if name not in funcs:
        raise Fail("%s: нет метода репозитория %s" % (sqlc, name))
    func = funcs[name]

    # Разделитель блока может нести не таблицу, а произвольный заголовок —
    # тогда привязки нет и колонки ищутся по всей схеме, как раньше.
    table = table if table in data["schema"] else None

    proto = data["proto_structs"]
    req_msg = name + "Request"
    resp_msg = name + "Response"
    for msg in (req_msg, resp_msg):
        if msg not in proto:
            raise Fail("%s: нет proto-сообщения %s" % (sqlc, msg))

    arg_info, go_fields = resolve_args(sqlc, data, name, func)
    req_fields = pair_fields(
        sqlc, data, name, go_fields, req_msg, column_hints, arg_info["params_type"] or name, table
    )
    response = resolve_response(sqlc, data, name, func, resp_msg, counters)

    call = {"sqlc": sqlc, "query": name, "kind": kind, "table": table, "req_fields": req_fields}
    call.update(arg_info)
    call.update(response)
    return call


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

        # Счётчики страниц собственными вызовами не становятся: их дёргает
        # сервис внутри постраничного метода. Собираем их на первом проходе.
        counters = set()
        calls = [
            resolve_query(sqlc, data, name, kind, table, hints.get(name) or {}, counters)
            for name, kind, table in data["queries"]
            if not name.startswith(COUNT_PREFIX)
        ]

        declared = {n for n, _, _ in data["queries"] if n.startswith(COUNT_PREFIX)}
        orphans = sorted(declared - counters)
        if orphans:
            raise Fail(
                "%s: счётчики без постраничной выборки: %s" % (sqlc, ", ".join(orphans))
            )

        rows = resolve_rows(
            sqlc,
            data,
            {(c["go_elem"], c["row_msg"], c["table"]) for c in calls if c.get("row_msg")},
        )

        model[sqlc] = {"meta": meta_by_sqlc[sqlc], "calls": calls, "rows": rows}
        total += len(calls)

    write_json(args.work_dir, "model_custom.json", model)

    print("всего запросов:", total)
    for sqlc, d in model.items():
        print("\n%s (%d запросов)" % (sqlc, len(d["calls"])))
        for c in d["calls"]:
            print("  %-45s :%-6s args=%d" % (c["query"], c["kind"], len(c["req_fields"])))


if __name__ == "__main__":
    main()
