# -*- coding: utf-8 -*-
"""Шаг 2 конвейера: сборка модели сущностей из build/parsed.json.

Любое расхождение между query.sql, sqlc и proto — ошибка, а не повод
догадываться: конвейер обязан упасть здесь, а не выдать кривой Go.
Результат — build/model.json.

Таблица не обязана иметь все 8 стандартных запросов — присутствует только
то, что реально есть в query.sql (см. documentation/dev_instructions/crud/
generate_gprpc_go.md). Обязателен только Create<Entity> — из него берутся
имя сущности и структура строки. Остальные пять именованных операций
(Get<E>ById, GetDeleted<E>ById, Update<E>ById, Delete<E>ById,
Undelete<E>ById) — не обязаны присутствовать.

Списочные запросы (`:many`, без обязательного аргумента-фильтра) собираются
в один список "lists" независимо от имени: генератору не важно, значит ли
имя "активные", "удалённые" или что-то ещё — важна форма (простой список
или постраничный с парным Count<имя>).

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


def pluralize(word):
    """Простая английская плюрализация — только для имени вспомогательного
    конвертера списка (%sToProto), не для сопоставления запросов: конкретные
    запросы находятся по точному имени (Get<Entity>ById и т.д.), а не
    угадыванием множественного числа."""
    if re.search(r"[sxz]$", word) or re.search(r"[cs]h$", word):
        return word + "es"
    if re.search(r"[^aeiouAEIOU]y$", word):
        return word[:-1] + "ies"
    return word + "s"


class Fail(Exception):
    pass


# Соглашение: proto-поле "<col>_external_id" резолвит внешний id через
# dc.user.external_id в internal id колонки "<col>_id" той же таблицы,
# вместо прямого bigint-параметра. В query.sql это выглядит как
# `<col>_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = $N)`.
# sqlc называет такой параметр по колонке сравнения — всегда "ExternalID".
EXTERNAL_ID_SUFFIX = "_external_id"
USER_TABLE = "dc.user"

# Постраничный ответ: список строк плюс общее для всего файла сообщение
# Pagination (см. pagination_proto.md) — не разбираем по имени запроса,
# только по форме: два поля, второе — *Pagination.
PAGINATION_MSG = "Pagination"
COUNT_PREFIX = "Count"

# Поля структуры параметров, из которых конвейер строит страницу. Тип и
# необходимость присутствия — фиксированы, состав запроса — нет: Order можно
# не запрашивать (тогда сортировки на выбор нет), Page и PageLimit нужны
# всегда — без них не собрать LIMIT/OFFSET и не заполнить Pagination.
PAGE_FIELD_TYPES = {"PageLimit": "int32", "Page": "int32", "Order": "string"}
REQUIRED_PAGE_FIELDS = {"PageLimit", "Page"}

# Типы полей фильтра постраничного списка (см. "Списочные запросы с фильтром"
# в optional_standard_ops.md) — те же, для которых у req_to_param_expr в
# gen.py есть правило прямого переноса значения из запроса в параметры.
FILTER_FIELD_TYPES = {"bool", "int64", "int16", "string", "uuid.UUID"}


def pair_fields(data, table, go_fields, proto_msg, what):
    """Сопоставляет поля Go-структуры с полями proto-сообщения по имени."""
    schema_cols = data["schema"][table]
    proto = data["proto_structs"]
    if proto_msg not in proto:
        raise Fail("%s: нет proto-сообщения %s" % (table, proto_msg))
    pfields = {norm(f["field"]): f for f in proto[proto_msg]}
    used = set()

    def find_proto_field(f):
        # ExternalID проверяется первым и отдельно от прямого совпадения по
        # имени: колонка с единственной внешней ссылкой в таблице называется
        # в proto просто "external_id" (без префикса column_), и тогда прямое
        # совпадение имени нашло бы его само — но это была бы ложная удача:
        # "external_id" не колонка этой таблицы, а конвенция резолва через
        # dc.user, и должна разбираться через base_col, а не как обычная
        # колонка. Исключение — таблица, у которой "external_id" — её
        # собственная колонка (dc.user): там это прямое поле, без резолва.
        if f["field"] == "ExternalID" and f["type"] == "uuid.UUID" and "external_id" not in schema_cols:
            candidates = [
                p
                for p in proto[proto_msg]
                if (p["proto_name"] or "") == "external_id"
                or (p["proto_name"] or "").endswith(EXTERNAL_ID_SUFFIX)
            ]
            candidates = [p for p in candidates if norm(p["field"]) not in used]
            if len(candidates) == 1:
                return candidates[0]
            if len(candidates) > 1:
                raise Fail(
                    "%s: несколько кандидатов на внешний id в %s: %s"
                    % (table, proto_msg, [p["field"] for p in candidates])
                )

        key = norm(f["field"])
        if key in pfields:
            return pfields[key]
        return None

    paired = []
    for f in go_fields:
        pf = find_proto_field(f)
        if pf is None:
            raise Fail("%s: поле %s.%s не найдено в %s" % (table, what, f["field"], proto_msg))
        used.add(norm(pf["field"]))
        col = pf["proto_name"]
        if col is None:
            raise Fail("%s: у поля %s.%s нет имени колонки" % (table, proto_msg, pf["field"]))

        if (col == "external_id" or col.endswith(EXTERNAL_ID_SUFFIX)) and col not in schema_cols:
            # "external_id" без префикса — единственная внешняя ссылка в
            # таблице, разрешаемая в user_id (см. standard_crud.md); с
            # префиксом ("<col>_external_id") база — "<col>_id". Если
            # "external_id" — собственная колонка этой таблицы (dc.user),
            # сюда не попадаем: ветка ниже находит её напрямую в schema_cols.
            base_col = "user_id" if col == "external_id" else col[: -len(EXTERNAL_ID_SUFFIX)] + "_id"
            if base_col not in schema_cols:
                raise Fail(
                    "%s: %s.%s (%s) подразумевает колонку %s, которой нет в schema.sql"
                    % (table, proto_msg, pf["field"], col, base_col)
                )
            user_cols = data["schema"].get(USER_TABLE)
            if not user_cols or "external_id" not in user_cols:
                raise Fail("%s: нет %s.external_id для разрешения %s" % (table, USER_TABLE, col))
            paired.append(
                {
                    "go": f["field"],
                    "go_type": f["type"],
                    "proto": pf["field"],
                    "proto_type": pf["type"],
                    "col": base_col,
                    "varchar": schema_cols[base_col]["varchar"],
                }
            )
            continue

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


def pair_filter_fields(data, table, go_fields, proto_msg, what):
    """Сопоставляет поля фильтра постраничного списка с полями proto-запроса
    по имени. В отличие от pair_fields, не требует существования колонки в
    schema.sql и не резолвит external_id: фильтр — это условие в WHERE
    (например LIKE), а не запись в конкретную колонку, и может не совпадать
    с ней 1:1 (см. "Списочные запросы с фильтром" в optional_standard_ops.md)."""
    proto = data["proto_structs"]
    if proto_msg not in proto:
        raise Fail("%s: нет proto-сообщения %s" % (table, proto_msg))
    pfields = {norm(f["field"]): f for f in proto[proto_msg]}
    paired = []
    for f in go_fields:
        pf = pfields.get(norm(f["field"]))
        if pf is None:
            raise Fail("%s: поле фильтра %s.%s не найдено в %s" % (table, what, f["field"], proto_msg))
        paired.append({"go": f["field"], "go_type": f["type"], "proto": pf["field"], "proto_type": pf["type"]})
    return paired


def row_or_empty_response(data, msg, entity):
    """Форма ответа операций по id (get/create/update/delete/undelete):

    - 0 полей   -> "empty" — ответ буквально пустой, вставлять нечего;
    - 1 поле, тип которого — сама сущность -> "row";
    - 1 поле любого другого типа (старое соглашение: поле-обёртка над
      google.protobuf.Empty) -> "field" — конвертер молча кладёт туда
      &emptypb.Empty{}, как и раньше.
    """
    fields = data["proto_structs"].get(msg)
    if fields is None:
        raise Fail("нет proto-сообщения %s" % msg)
    if not fields:
        return {"form": "empty", "field": None}
    if len(fields) == 1:
        f = fields[0]
        if f["type"] in (entity, "*" + entity):
            return {"form": "row", "field": f["field"]}
        return {"form": "field", "field": f["field"]}
    raise Fail("%s: ожидалось 0 или 1 поле, найдено %s" % (msg, [f["field"] for f in fields]))


def list_response(data, msg):
    """Форма ответа списочного запроса: простой список или страница.

    Страница — это ровно два поля: список строк и *Pagination (общее для
    всего файла сообщение, см. resolve_pagination_fields). Опознаём по типу
    второго поля, а не по имени первого — оно может быть "rows", "data" или
    как угодно ещё."""
    fields = data["proto_structs"].get(msg)
    if not fields:
        raise Fail("нет proto-сообщения %s или оно пустое: %s" % (msg, msg))
    names = [f["proto_name"] for f in fields]
    if len(fields) == 2 and fields[1]["type"] == "*%s" % PAGINATION_MSG:
        if not fields[0]["type"].startswith("[]"):
            raise Fail("%s: первое поле страницы обязано быть списком, найдено %s" % (msg, fields[0]["type"]))
        return {
            "form": "page",
            "resp_field": fields[0]["field"],
            "row_type": fields[0]["type"],
            "pagination_field": fields[1]["field"],
        }
    if len(fields) == 1:
        return {"form": "plain", "resp_field": fields[0]["field"], "row_type": fields[0]["type"]}
    raise Fail("%s: неожиданная форма ответа списка: %s" % (msg, names))


def resolve_pagination_fields(data):
    """Поля общего сообщения Pagination в порядке объявления: [page, per_page,
    total_items, total_pages]. Порядок фиксирован (см. pagination_proto.md —
    "одна на файл, поля одни и те же везде"), конкретные имена — нет, поэтому
    берём позиционно, а не по имени, как и в остальном конвейере."""
    fields = data["proto_structs"].get(PAGINATION_MSG)
    if not fields or len(fields) != 4:
        raise Fail(
            "сообщение %s обязано иметь ровно 4 поля (page, per_page, total_items, total_pages), найдено %s"
            % (PAGINATION_MSG, [f["field"] for f in (fields or [])])
        )
    return {
        "page": fields[0]["field"],
        "per_page": fields[1]["field"],
        "total_items": fields[2]["field"],
        "total_pages": fields[3]["field"],
    }


def resolve_exec_args(data, name):
    """Форма аргумента операций по id: голый id или структура параметров."""
    func = data["funcs"][name]
    args = func["args"]
    if len(args) == 1 and args[0]["name"] == "id" and args[0]["type"] == "int64":
        return {"arg_form": "id", "params_type": None}, None
    if len(args) == 1 and args[0]["type"] in data["param_structs"]:
        ptype = args[0]["type"]
        return {"arg_form": "params", "params_type": ptype}, data["param_structs"][ptype]
    raise Fail("%s: неожиданная форма аргументов: %s" % (name, args))


def resolve_id_op(data, table, entity, name):
    """Общая логика для Get<E>ById / GetDeleted<E>ById: строго один id-аргумент,
    ответ обязан быть строкой сущности."""
    if name not in data["funcs"]:
        raise Fail("%s: нет метода репозитория %s" % (table, name))
    resp = row_or_empty_response(data, name + "Response", entity)
    if resp["form"] != "row":
        raise Fail("%s: %s обязан возвращать строку %s" % (table, name, entity))
    return {"query": name, "resp_field": resp["field"]}


def resolve_exec_op(data, table, entity, name):
    """Общая логика для Update<E>ById / Delete<E>ById / Undelete<E>ById:
    аргумент — id или структура параметров; ответ — строка, пустое поле-обёртка
    или буквально пустое сообщение."""
    if name not in data["funcs"]:
        raise Fail("%s: нет метода репозитория %s" % (table, name))
    arg_info, go_fields = resolve_exec_args(data, name)
    resp = row_or_empty_response(data, name + "Response", entity)
    entry = {
        "query": name,
        "arg_form": arg_info["arg_form"],
        "params_type": arg_info["params_type"],
        "resp_form": resp["form"],
        "resp_field": resp["field"],
    }
    if arg_info["arg_form"] == "params":
        entry["fields"] = pair_fields(data, table, go_fields, name + "Request", arg_info["params_type"])
        if not any(f["col"] == "id" for f in entry["fields"]):
            raise Fail("%s: %s принимает структуру параметров без колонки id" % (table, name))
    return entry


def resolve_entity(sqlc, data, block):
    table = block["table"].replace('"', "")  # dc.host — parse_queries не снимает кавычки
    short = table.split(".")[-1]
    names = [q[0] for q in block["queries"]]
    kinds = dict(block["queries"])
    funcs = data["funcs"]
    proto = data["proto_structs"]

    creates = [n for n in names if n.startswith("Create")]
    if len(creates) != 1:
        raise Fail("%s: ожидался ровно один Create, найдено %s" % (table, creates))
    entity = creates[0][len("Create"):]
    E = entity

    claimed = set()

    def take(name):
        if name in names and name not in claimed:
            claimed.add(name)
            return True
        return False

    # ---------- Create<E>: обязателен, отсюда берётся сущность и строка ----------
    create_name = "Create%s" % E
    if not take(create_name):
        raise Fail("%s: нет запроса %s" % (table, create_name))
    if create_name not in funcs:
        raise Fail("%s: нет метода репозитория %s" % (table, create_name))
    create_params = "Create%sParams" % E
    if create_params not in data["param_structs"]:
        raise Fail("%s: нет структуры %s" % (table, create_params))

    row_struct = funcs[create_name]["ret"].split(",")[0].strip()
    if row_struct not in data["model_structs"]:
        raise Fail("%s: нет модели %s" % (table, row_struct))
    if entity not in proto:
        raise Fail("%s: нет proto-сообщения %s" % (table, entity))

    create_fields = pair_fields(
        data, table, data["param_structs"][create_params], "Create%sRequest" % E, create_params
    )
    create_resp = row_or_empty_response(data, "Create%sResponse" % E, entity)
    if create_resp["form"] != "row":
        raise Fail("%s: %s обязан возвращать строку %s" % (table, create_name, entity))

    row_fields = pair_fields(data, table, data["model_structs"][row_struct], entity, row_struct)

    # ---------- Пять именованных операций: у таблицы может не быть любой из них ----------
    ops = {}

    name = "Get%sById" % E
    if take(name):
        ops["get_by_id"] = resolve_id_op(data, table, entity, name)

    name = "GetDeleted%sById" % E
    if take(name):
        ops["get_deleted_by_id"] = resolve_id_op(data, table, entity, name)

    name = "Update%sById" % E
    if take(name):
        op = resolve_exec_op(data, table, entity, name)
        if op["arg_form"] != "params":
            raise Fail("%s: %s обязан принимать структуру параметров" % (table, name))
        ops["update"] = op

    for op_key, prefix in (("delete", "Delete"), ("undelete", "Undelete")):
        name = "%s%sById" % (prefix, E)
        if take(name):
            ops[op_key] = resolve_exec_op(data, table, entity, name)

    # ---------- Count<...> размечаются, но разбираются вместе с их выборкой ----------
    declared_counters = {n for n in names if n.startswith(COUNT_PREFIX) and n not in claimed}
    consumed_counters = set()

    remaining = [n for n in names if n not in claimed and n not in declared_counters]

    fks = []
    lookups = []
    lists = []

    for name in remaining:
        kind = kinds[name]
        if name not in funcs:
            raise Fail("%s: нет метода репозитория %s" % (table, name))
        func = funcs[name]
        args = func["args"]

        if kind == "one":
            m = re.match(r"^Get%sBy(\w+)$" % re.escape(E), name)
            if not m:
                raise Fail("%s: нераспознанный запрос %s (:one)" % (table, name))
            if len(args) != 1:
                raise Fail("%s: у %s ожидался один аргумент, найдено %s" % (table, name, args))
            arg = args[0]
            req_msg = name + "Request"
            pf = proto.get(req_msg)
            if not pf or len(pf) != 1:
                raise Fail("%s: в %s ожидалось одно поле" % (table, req_msg))
            col = pf[0]["proto_name"]
            schema_cols = data["schema"][table]
            if col not in schema_cols:
                raise Fail("%s: колонки %s нет в schema.sql" % (table, col))
            resp = row_or_empty_response(data, name + "Response", entity)
            if resp["form"] != "row":
                raise Fail("%s: %s обязан возвращать строку %s" % (table, name, entity))
            lookups.append(
                {
                    "query": name,
                    "arg_name": arg["name"],
                    "arg_type": arg["type"],
                    "proto_field": pf[0]["field"],
                    "proto_type": pf[0]["type"],
                    "col": col,
                    "varchar": schema_cols[col]["varchar"],
                    "resp_field": resp["field"],
                }
            )
            continue

        if kind != "many":
            raise Fail("%s: нераспознанный запрос %s (:%s)" % (table, name, kind))

        if len(args) == 1 and args[0]["type"] not in data["param_structs"]:
            # Выборка по неуникальной колонке (FK): один голый скалярный аргумент.
            arg = args[0]
            req_msg = name + "Request"
            pf = proto.get(req_msg)
            if not pf or len(pf) != 1:
                raise Fail("%s: в %s ожидалось одно поле" % (table, req_msg))
            col = pf[0]["proto_name"]
            resp = list_response(data, name + "Response")
            if resp["form"] != "plain":
                raise Fail("%s: %s — выборка по колонке пока не поддерживает пагинацию" % (table, name))
            fks.append(
                {
                    "query": name,
                    "arg_name": arg["name"],
                    "arg_type": arg["type"],
                    "proto_field": pf[0]["field"],
                    "col": col,
                    "resp_field": resp["resp_field"],
                }
            )
            continue

        if len(args) == 0:
            resp = list_response(data, name + "Response")
            if resp["form"] != "plain":
                raise Fail("%s: %s без аргументов не может отдавать страницу" % (table, name))
            lists.append({"query": name, "arg_form": "none", "params_type": None, "paginated": False, "resp_field": resp["resp_field"]})
            continue

        if len(args) == 1 and args[0]["type"] in data["param_structs"]:
            ptype = args[0]["type"]
            pfields = {f["field"]: f["type"] for f in data["param_structs"][ptype]}
            pnames = set(pfields)
            missing = REQUIRED_PAGE_FIELDS - pnames
            if missing:
                raise Fail(
                    "%s: %s — в параметрах не хватает обязательных полей страницы %s"
                    % (table, name, sorted(missing))
                )
            for fname in pnames & set(PAGE_FIELD_TYPES):
                if PAGE_FIELD_TYPES[fname] != pfields[fname]:
                    raise Fail(
                        "%s: %s.%s имеет тип %s, ожидался %s"
                        % (table, ptype, fname, pfields[fname], PAGE_FIELD_TYPES[fname])
                    )

            # Поля сверх пагинации — фильтры постраничного списка (см.
            # "Списочные запросы с фильтром" в optional_standard_ops.md):
            # проверенный по типу набор полей, парный с одноимёнными полями
            # в Request, без привязки к конкретной колонке schema.sql.
            filter_names = pnames - set(PAGE_FIELD_TYPES)
            unsupported = {n: pfields[n] for n in filter_names if pfields[n] not in FILTER_FIELD_TYPES}
            if unsupported:
                raise Fail(
                    "%s: %s — не поддерживаются типы полей фильтра %s (ожидались %s)"
                    % (table, name, unsupported, sorted(FILTER_FIELD_TYPES))
                )
            req_msg = name + "Request"
            filter_go_fields = [f for f in data["param_structs"][ptype] if f["field"] in filter_names]
            filter_fields = (
                pair_filter_fields(data, table, filter_go_fields, req_msg, ptype) if filter_go_fields else []
            )

            resp = list_response(data, name + "Response")
            if resp["form"] != "page":
                raise Fail("%s: %s принимает страницу, но ответ не постраничный" % (table, name))

            counter = COUNT_PREFIX + name
            if counter not in funcs:
                raise Fail("%s: у постраничной выборки %s нет счётчика %s" % (table, name, counter))
            cargs = funcs[counter]["args"]
            if not filter_fields and len(cargs) == 1 and cargs[0]["type"] == "int32":
                # Без фильтра: как раньше, счётчик принимает голый page_limit.
                counter_arg_form = "int32"
                counter_params_type = None
                counter_arg = cargs[0]["name"]
            elif len(cargs) == 1 and cargs[0]["type"] in data["param_structs"]:
                # С фильтром: точный подсчёт требует и его значение, поэтому
                # счётчик берёт структуру параметров — PageLimit плюс те же
                # поля фильтра, что и у самой выборки, не больше и не меньше.
                counter_arg_form = "params"
                counter_params_type = cargs[0]["type"]
                counter_arg = None
                counter_pfields = {f["field"] for f in data["param_structs"][counter_params_type]}
                expected = {"PageLimit"} | {f["go"] for f in filter_fields}
                if counter_pfields != expected:
                    raise Fail(
                        "%s: счётчик %s принимает поля %s, ожидались %s"
                        % (table, counter, sorted(counter_pfields), sorted(expected))
                    )
            else:
                raise Fail(
                    "%s: счётчик %s обязан принимать один аргумент — либо page_limit (int32), "
                    "либо структуру параметров из PageLimit и полей фильтра, найдено %s"
                    % (table, counter, cargs)
                )
            counter_row = funcs[counter]["ret"].split(",")[0].strip()
            counter_fields = data["param_structs"].get(counter_row)
            if not counter_fields or len(counter_fields) != 2:
                raise Fail(
                    "%s: %s обязан возвращать ровно два поля (число строк и число страниц), найдено %s"
                    % (table, counter_row, [f["field"] for f in (counter_fields or [])])
                )
            consumed_counters.add(counter)

            pagination = resolve_pagination_fields(data)
            lists.append(
                {
                    "query": name,
                    "arg_form": "params",
                    "params_type": ptype,
                    "paginated": True,
                    "resp_field": resp["resp_field"],
                    "pagination_field": resp["pagination_field"],
                    "order_field": "Order" if "Order" in pnames else None,
                    "filter_fields": filter_fields,
                    "counter": counter,
                    "counter_arg_form": counter_arg_form,
                    "counter_params_type": counter_params_type,
                    "counter_arg": counter_arg,
                    "counter_total_field": counter_fields[0]["field"],
                    "counter_pages_field": counter_fields[1]["field"],
                    "pagination_msg_fields": pagination,
                }
            )
            continue

        raise Fail("%s: нераспознанный запрос %s (:many, args=%s)" % (table, name, args))

    orphans = sorted(declared_counters - consumed_counters)
    if orphans:
        raise Fail("%s: счётчики без постраничной выборки: %s" % (table, ", ".join(orphans)))

    return {
        "sqlc": sqlc,
        "table": table,
        "short": short,
        "file": snake(entity),
        "entity": entity,
        "plural": pluralize(entity),
        "row_struct": row_struct,
        "create_params": create_params,
        "create_fields": create_fields,
        "create_resp_field": create_resp["field"],
        "row_fields": row_fields,
        "ops": ops,
        "fks": fks,
        "lookups": lookups,
        "lists": lists,
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
            total_rpc += 1 + len(e["ops"]) + len(e["fks"]) + len(e["lookups"]) + len(e["lists"])

    write_json(args.work_dir, "model.json", model)

    print("всего RPC:", total_rpc)
    for sqlc, d in model.items():
        print("\n%s (%d таблиц)" % (sqlc, len(d["entities"])))
        for e in d["entities"]:
            ops = ",".join(sorted(e["ops"]))
            fk = " fk=%s" % ",".join(f["query"] for f in e["fks"]) if e["fks"] else ""
            lk = " lookup=%s" % ",".join(l["query"] for l in e["lookups"]) if e["lookups"] else ""
            ls = " lists=%s" % ",".join(l["query"] for l in e["lists"]) if e["lists"] else ""
            print(
                "  %-24s %-16s file=%-20s ops=%s%s%s%s"
                % (e["table"], e["entity"], e["file"] + ".go", ops or "-", fk, lk, ls)
            )

    # Типы, которые встречаются в полях — список должен быть закрытым.
    types = set()
    for d in model.values():
        for e in d["entities"]:
            for f in e["row_fields"] + e["create_fields"]:
                types.add((f["go_type"], f["proto_type"]))
            for l in e["lists"]:
                for f in l.get("filter_fields") or []:
                    types.add((f["go_type"], f["proto_type"]))
            if "update" in e["ops"]:
                for f in e["ops"]["update"].get("fields") or []:
                    types.add((f["go_type"], f["proto_type"]))
            for op_key in ("delete", "undelete"):
                if op_key in e["ops"]:
                    for f in e["ops"][op_key].get("fields") or []:
                        types.add((f["go_type"], f["proto_type"]))
    print("\nпары типов sqlc -> proto:")
    for t in sorted(types):
        print("  %-16s -> %s" % t)


main()
