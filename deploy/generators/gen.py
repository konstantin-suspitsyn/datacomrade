# -*- coding: utf-8 -*-
"""Шаг 3 конвейера: генерация слоёв converter / validation / service / api.

Читает build/model.json и переписывает по файлу на таблицу в каждом слое.
ВНИМАНИЕ: перезаписывает файлы без спроса. Список того, что генерируется,
а что писано руками, — в documentation/dev_instructions/crud/generate_gprpc_go.md

Таблица не обязана иметь весь стандартный набор из восьми запросов —
генерируется только то, что реально есть в model.json (см. resolve.py).
Три операции по id — Update/Delete/Undelete<E>ById — могут как возвращать
строку (:one, RETURNING), так и быть :exec с пустым ответом; в последнем
случае аргумент может быть голым id или структурой параметров (когда
запрос заодно резолвит user_id через external_id). Списочные запросы
(поле "lists") бывают простыми (без аргументов, один ответ-слайс) или
постраничными (Params{Order?,Page,PageLimit}, ответ — список плюс общее
сообщение Pagination, парный Count<имя>). Строка ответа — либо сущность
(<Entity>ToProto), либо собственная строка запроса, если тот джойнит соседние
таблицы (<Query>RowToProto) — см.
documentation/dev_instructions/crud/optional_standard_ops.md

Запуск:  python deploy/generators/gen.py [--config ...] [--work-dir ...]
"""
import io
import os

from genconfig import parse_args, read_json, repo_path

# Заполняются в main() из конфига.
INTERNAL = None
CONV_IMPORT = None
VALID_IMPORT = None
ERRORS_IMPORT = None
VALIDATOR_IMPORT = None
APIERR_IMPORT = None

written = []


def write(rel, text):
    path = os.path.join(INTERNAL, rel.replace("/", os.sep))
    d = os.path.dirname(path)
    if not os.path.isdir(d):
        os.makedirs(d)
    with io.open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(text)
    written.append(rel)


def lower_first(s):
    return s[0].lower() + s[1:] if s else s


def limit_const(entity, proto_field):
    return "%s%sMaxLen" % (lower_first(entity), proto_field)


# ---------- преобразование значений ----------
def row_to_proto_expr(f, var):
    """Выражение для поля proto-сущности из поля строки таблицы."""
    src = "%s.%s" % (var, f["go"])
    pair = (f["go_type"], f["proto_type"])
    if pair in (("bool", "bool"), ("int64", "int64"), ("string", "string")):
        return src
    if pair == ("int16", "int32"):
        return "int32(%s)" % src
    if pair == ("time.Time", "*timestamppb.Timestamp"):
        return "converter.TimeToProto(%s)" % src
    if pair == ("sql.NullTime", "*timestamppb.Timestamp"):
        return "converter.NullTimeToProto(%s)" % src
    if pair == ("sql.NullInt64", "*int64"):
        # Nullable-колонка (dc.column_cat.alias_id): NULL обязан доехать до
        # клиента как отсутствие значения, а не как 0 — для FK это разное.
        return "converter.NullInt64ToProto(%s)" % src
    if pair == ("sql.NullString", "*string"):
        return "converter.NullStringToProto(%s)" % src
    if pair == ("string", "*string"):
        return "converter.StringToProto(%s)" % src
    if pair == ("int64", "*int64"):
        # Колонка NOT NULL, а поле в .proto помечено optional.
        return "converter.Int64ToProto(%s)" % src
    if pair == ("uuid.UUID", "string"):
        return "converter.UUIDToProto(%s)" % src
    raise Exception("нет правила для %s -> %s" % pair)


# Необязательный параметр (sqlc.narg) ↔ optional-поле proto: "не передан" и
# "передан" различаются, поэтому значение берётся прямо из поля-указателя, а
# не из геттера — тот отдаёт 0/false и на nil, и на явно переданный ноль.
# Из-за прямого доступа к полю такие функции нуждаются в защите от nil-req
# (см. req_guard), которую геттеры дают бесплатно.
NULLABLE_PAIRS = {
    ("sql.NullInt64", "*int64"): "converter.ProtoToNullInt64",
    ("sql.NullBool", "*bool"): "converter.ProtoToNullBool",
    ("sql.NullString", "*string"): "converter.ProtoToNullString",
    ("uuid.NullUUID", "*string"): "converter.ProtoToNullUUID",
}


def req_to_param_expr(f, var):
    """Выражение для поля sqlc-параметров из поля запроса gRPC."""
    getter = "%s.Get%s()" % (var, f["proto"])
    pair = (f["go_type"], f["proto_type"])
    if pair in (("bool", "bool"), ("int64", "int64"), ("string", "string")):
        return getter
    if pair == ("int16", "int32"):
        return "int16(%s)" % getter
    if pair == ("string", "*string"):
        return getter
    if pair == ("uuid.UUID", "string"):
        return "converter.ProtoToUUID(%s)" % getter
    if pair == ("time.Time", "*timestamppb.Timestamp"):
        # Клиент явно передаёт значение (например created_at/updated_at в
        # Update/Delete) — не nullable, поле не помечено optional в .proto.
        return "converter.ProtoToTime(%s)" % getter
    if pair in NULLABLE_PAIRS:
        return "%s(%s.%s)" % (NULLABLE_PAIRS[pair], var, f["proto"])
    raise Exception("нет правила для %s -> %s" % pair)


def needs_req_guard(fields):
    """Есть ли среди полей такие, что читаются напрямую, а не через геттер."""
    return any((f["go_type"], f["proto_type"]) in NULLABLE_PAIRS for f in fields)


def req_guard(proto, req_type):
    """Замена nil-запроса пустым: дальше поля читаются напрямую (см. NULLABLE_PAIRS),
    а геттеры protobuf безопасны на nil — сервер не должен падать ни там, ни там."""
    return [
        "\tif req == nil {",
        "\t\treq = &%s.%s{}" % (proto, req_type),
        "\t}",
        "",
    ]


def lookup_arg_expr(lk, var):
    getter = "%s.Get%s()" % (var, lk["proto_field"])
    pair = (lk["arg_type"], lk["proto_type"])
    if pair in (("int64", "int64"), ("string", "string")):
        return getter
    if pair == ("uuid.UUID", "string"):
        return "converter.ProtoToUUID(%s)" % getter
    raise Exception("нет правила для аргумента выборки %s -> %s" % pair)


def op_fields(e, op_key):
    op = e["ops"].get(op_key)
    if not op:
        return []
    return op.get("fields") or []


def all_param_fields(e):
    """Все поля, приходящие через структуры параметров (для needs_*/limits)."""
    fields = list(e["create_fields"])
    for op_key in ("update", "delete", "undelete"):
        fields += op_fields(e, op_key)
    return fields


def response_rows(e):
    """Все строки, которые сущность отдаёт наружу, по одной записи на запрос."""
    rows = [e["ops"][k]["row"] for k in ("get_by_id", "get_deleted_by_id") if k in e["ops"]]
    rows += [l["row"] for l in e["lists"]]
    rows += [fk["row"] for fk in e["fks"]]
    rows += [lk["row"] for lk in e["lookups"]]
    return rows


def query_rows(e):
    """Собственные строки запросов (join отдаёт колонки соседних таблиц) —
    по одной на структуру sqlc, с пометкой, нужен ли слайс-конвертер."""
    seen = {}
    slices = set()
    for l in e["lists"]:
        slices.add(l["row"]["go"])
    for fk in e["fks"]:
        slices.add(fk["row"]["go"])
    for r in response_rows(e):
        if not r["own"] and r["go"] not in seen:
            seen[r["go"]] = dict(r)
    for go, r in seen.items():
        r["slice"] = go in slices
    return list(seen.values())


def query_row_fields(e):
    fields = []
    for r in query_rows(e):
        fields += r["fields"]
    return fields


def all_filter_fields(e):
    """Поля фильтра всех постраничных списков сущности (см. gen_converter)."""
    fields = []
    for l in e["lists"]:
        if l["paginated"]:
            fields += l.get("filter_fields") or []
    return fields


def needs_converter_pkg(e):
    for f in e["row_fields"] + all_param_fields(e) + all_filter_fields(e) + query_row_fields(e):
        if f["proto_type"] in ("*timestamppb.Timestamp", "*string"):
            return True
        if f["go_type"] == "uuid.UUID":
            return True
        if (f["go_type"], f["proto_type"]) in NULLABLE_PAIRS:
            return True
    for lk in e["lookups"]:
        if lk["arg_type"] == "uuid.UUID":
            return True
    return False


def needs_uuid_pkg(e, extra_fields=()):
    for f in list(e["row_fields"]) + list(extra_fields) + all_filter_fields(e) + query_row_fields(e):
        if f["go_type"] in ("uuid.UUID", "uuid.NullUUID"):
            return True
    for lk in e["lookups"]:
        if lk["arg_type"] == "uuid.UUID":
            return True
    return False


# ---------- converter ----------
def gen_converter(meta, e):
    proto = meta["proto_alias"]
    repo = meta["repo_alias"]
    E, P = e["entity"], e["plural"]

    paginated_lists = [l for l in e["lists"] if l["paginated"]]

    imports = []
    if any(lk["arg_type"] == "uuid.UUID" for lk in e["lookups"]):
        imports.append('\t"github.com/google/uuid"')
        imports.append("")
    if needs_converter_pkg(e):
        imports.append('\t"%s"' % CONV_IMPORT)
    if paginated_lists:
        imports.append('\t"%s"' % VALID_IMPORT)
    imports.append('\t"%s"' % meta["repo_import"])
    imports.append("\t" + meta["proto_import"])

    lines = []
    lines.append("package %s" % meta["conv_pkg"])
    lines.append("")
    lines.append("import (")
    lines.extend(imports)
    lines.append(")")
    lines.append("")

    # Сущность -> proto
    lines.append("// %sToProto переводит строку %s в сущность gRPC." % (E, e["table"]))
    lines.append("func %sToProto(row %s.%s) *%s.%s {" % (E, repo, e["row_struct"], proto, E))
    lines.append("\treturn &%s.%s{" % (proto, E))
    for f in e["row_fields"]:
        lines.append("\t\t%s: %s," % (f["proto"], row_to_proto_expr(f, "row")))
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    lines.append("// %sToProto переводит список строк %s в список сущностей gRPC." % (P, e["table"]))
    lines.append("// Для пустого входа возвращается пустой, а не nil-слайс.")
    lines.append("func %sToProto(rows []%s.%s) []*%s.%s {" % (P, repo, e["row_struct"], proto, E))
    lines.append("\titems := make([]*%s.%s, 0, len(rows))" % (proto, E))
    lines.append("")
    lines.append("\tfor _, row := range rows {")
    lines.append("\t\titems = append(items, %sToProto(row))" % E)
    lines.append("\t}")
    lines.append("")
    lines.append("\treturn items")
    lines.append("}")
    lines.append("")

    # Собственные строки запросов: join отдаёт колонки соседних таблиц, у такой
    # строки свой набор полей и своё proto-сообщение (см. resolve_row).
    for r in query_rows(e):
        lines.append("// %sToProto переводит строку выборки в сообщение gRPC." % r["conv"])
        lines.append("func %sToProto(row %s.%s) *%s.%s {" % (r["conv"], repo, r["go"], proto, r["proto"]))
        lines.append("\treturn &%s.%s{" % (proto, r["proto"]))
        for f in r["fields"]:
            lines.append("\t\t%s: %s," % (f["proto"], row_to_proto_expr(f, "row")))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        if not r["slice"]:
            continue

        lines.append("// %sToProto переводит список строк выборки в список сообщений gRPC." % r["conv_plural"])
        lines.append("// Для пустого входа возвращается пустой, а не nil-слайс.")
        lines.append(
            "func %sToProto(rows []%s.%s) []*%s.%s {" % (r["conv_plural"], repo, r["go"], proto, r["proto"])
        )
        lines.append("\titems := make([]*%s.%s, 0, len(rows))" % (proto, r["proto"]))
        lines.append("")
        lines.append("\tfor _, row := range rows {")
        lines.append("\t\titems = append(items, %sToProto(row))" % r["conv"])
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn items")
        lines.append("}")
        lines.append("")

    if e["has_create"]:
        lines.append("// ToCreate%sParams собирает параметры вставки %s из запроса gRPC." % (E, e["table"]))
        lines.append("// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.")
        lines.append(
            "func ToCreate%sParams(req *%s.Create%sRequest) %s.%s {" % (E, proto, E, repo, e["create_params"])
        )
        if needs_req_guard(e["create_fields"]):
            lines.extend(req_guard(proto, "Create%sRequest" % E))
        lines.append("\treturn %s.%s{" % (repo, e["create_params"]))
        for f in e["create_fields"]:
            lines.append("\t\t%s: %s," % (f["go"], req_to_param_expr(f, "req")))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    for op_key, doc in (
        ("update", "обновления"),
        ("delete", "мягкого удаления"),
        ("undelete", "обратного удаления"),
    ):
        op = e["ops"].get(op_key)
        if not op or op["arg_form"] != "params":
            continue
        name = op["query"]
        lines.append("// To%sParams собирает параметры %s %s из запроса gRPC." % (name, doc, e["table"]))
        lines.append(
            "func To%sParams(req *%s.%sRequest) %s.%s {" % (name, proto, name, repo, op["params_type"])
        )
        if needs_req_guard(op["fields"]):
            lines.extend(req_guard(proto, name + "Request"))
        lines.append("\treturn %s.%s{" % (repo, op["params_type"]))
        for f in op["fields"]:
            lines.append("\t\t%s: %s," % (f["go"], req_to_param_expr(f, "req")))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    for lk in e["lookups"]:
        lines.append(
            "// To%sArg достаёт из запроса gRPC значение %s для выборки %s."
            % (lk["query"], lk["col"], e["table"])
        )
        lines.append(
            "func To%sArg(req *%s.%sRequest) %s {" % (lk["query"], proto, lk["query"], lk["arg_type"])
        )
        lines.append("\treturn %s" % lookup_arg_expr(lk, "req"))
        lines.append("}")
        lines.append("")

    for l in paginated_lists:
        name = l["query"]
        lines.append("// To%sParams собирает параметры страницы %s из запроса gRPC." % (name, e["table"]))
        lines.append(
            "func To%sParams(req *%s.%sRequest) %s.%s {" % (name, proto, name, repo, l["params_type"])
        )
        if needs_req_guard(l.get("filter_fields") or []):
            lines.extend(req_guard(proto, name + "Request"))
        lines.append("\tlimit := req.GetPageLimit()")
        lines.append("\tif limit == 0 {")
        lines.append("\t\tlimit = validation.DefaultPageSize")
        lines.append("\t}")
        lines.append("")
        lines.append("\tpage := req.GetPage()")
        lines.append("\tif page == 0 {")
        lines.append("\t\tpage = 1")
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn %s.%s{" % (repo, l["params_type"]))
        if l["order_field"]:
            lines.append("\t\tOrder:     req.GetOrder(),")
        if l.get("order_by_field"):
            lines.append("\t\tOrderBy:   req.GetOrderBy(),")
        lines.append("\t\tPage:      page,")
        lines.append("\t\tPageLimit: limit,")
        for f in l.get("filter_fields") or []:
            lines.append("\t\t%s: %s," % (f["go"], req_to_param_expr(f, "req")))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


# ---------- значения для тестов ----------
def uuid_literal(idx, variant=0):
    return "00000000-0000-4000-8000-%012d" % (idx + variant * 1000 + 1)


def sample_value(f, idx, variant=0, ptr_helper=None):
    t = f["go_type"]
    tag = f.get("col") or f["go"]
    if t == "string":
        # Поля фильтра постраничного списка не привязаны к колонке (нет "col").
        return '"%s-%d"' % (tag.replace("_", "-"), variant)
    if t == "int64":
        return str(100 + idx + variant * 1000)
    if t == "int16":
        return str(10 + idx + variant)
    if t == "bool":
        return "true"
    if t == "uuid.UUID":
        return '"%s"' % uuid_literal(idx, variant)
    if t == "sql.NullInt64":
        # optional int64 в proto — литерал обязан быть указателем.
        if ptr_helper is None:
            raise Exception("для %s (%s) нужен ptr-хелпер" % (tag, t))
        return "%s(int64(%d))" % (ptr_helper, 100 + idx + variant * 1000)
    # Молча вернуть None значило бы записать в .go строку "None" — генератор
    # обязан падать на неизвестном типе, а не выдавать несобираемый код.
    raise Exception("нет тестового значения для %s (%s)" % (tag, t))


def repo_sample_value(f, idx, variant=0):
    if f["go_type"] == "sql.NullInt64":
        return "sql.NullInt64{Int64: %d, Valid: true}" % (100 + idx + variant * 1000)
    value = sample_value(f, idx, variant)
    if f["go_type"] == "uuid.UUID":
        return "uuid.MustParse(%s)" % value
    return value


def filter_proto_literal(f, idx, ptr_helper):
    """Значение поля фильтра в proto-запросе теста. Необязательный фильтр —
    указатель, поэтому значение заворачивается локальным хелпером."""
    pair = (f["go_type"], f["proto_type"])
    if pair == ("sql.NullInt64", "*int64"):
        return "%s(int64(%d))" % (ptr_helper, 100 + idx)
    if pair == ("sql.NullBool", "*bool"):
        return "%s(true)" % ptr_helper
    if pair == ("sql.NullString", "*string"):
        return '%s("%s-%d")' % (ptr_helper, f["go"].lower(), idx)
    if pair == ("uuid.NullUUID", "*string"):
        return '%s("%s")' % (ptr_helper, uuid_literal(idx))
    return sample_value(f, idx)


def filter_repo_literal(f, idx):
    """То же значение на стороне sqlc-параметров: у необязательного фильтра
    переданное значение всегда Valid — отсутствие проверяется отдельно."""
    # Составной литерал в условии `if` обязан быть в скобках — иначе `{`
    # разбирается как начало тела блока.
    pair = (f["go_type"], f["proto_type"])
    if pair == ("sql.NullInt64", "*int64"):
        return "(sql.NullInt64{Int64: %d, Valid: true})" % (100 + idx)
    if pair == ("sql.NullBool", "*bool"):
        return "(sql.NullBool{Bool: true, Valid: true})"
    if pair == ("sql.NullString", "*string"):
        return '(sql.NullString{String: "%s-%d", Valid: true})' % (f["go"].lower(), idx)
    if pair == ("uuid.NullUUID", "*string"):
        return '(uuid.NullUUID{UUID: uuid.MustParse("%s"), Valid: true})' % uuid_literal(idx)
    return repo_sample_value(f, idx)


def gen_row_sample_block(repo, struct, func_name, fields, time_prefix, doc):
    """test<Что>Row() — строка со значениями, различимыми между полями."""
    lines = []
    lines.append("// %s — %s со значениями, различимыми между полями." % (func_name, doc))
    lines.append("func %s() %s.%s {" % (func_name, repo, struct))
    lines.append('\treturn %s.%s{' % (repo, struct))
    for i, f in enumerate(fields):
        if f["go_type"] == "time.Time":
            v = "%sCreatedAt" % time_prefix if f["col"] == "created_at" else "%sUpdatedAt" % time_prefix
        elif f["go_type"] == "sql.NullTime":
            v = "sql.NullTime{Time: %sUpdatedAt, Valid: true}" % time_prefix
        elif f["col"] == "is_deleted":
            v = "false"
        else:
            v = repo_sample_value(f, i)
        lines.append('\t\t%s: %s,' % (f["go"], v))
    lines.append('\t}')
    lines.append("}")
    lines.append("")
    return lines


def gen_row_to_proto_test(conv, sample_func, fields):
    """Test<Конвертер>ToProto: каждое поле строки доехало без потери значения."""
    lines = []
    lines.append("func Test%sToProto(t *testing.T) {" % conv)
    lines.append('\trow := %s()' % sample_func)
    lines.append('\tgot := %sToProto(row)' % conv)
    lines.append("")
    lines.append('\tif got == nil {')
    lines.append('\t\tt.Fatal("%sToProto() = nil, want value")' % conv)
    lines.append('\t}')
    lines.append("")
    for f in fields:
        g = "got.Get%s()" % f["proto"]
        pair = (f["go_type"], f["proto_type"])
        if pair == ("time.Time", "*timestamppb.Timestamp"):
            lines.append('\tif !%s.AsTime().Equal(row.%s) {' % (g, f["go"]))
            lines.append('\t\tt.Errorf("%s = %%v, want %%v", %s.AsTime(), row.%s)' % (f["proto"], g, f["go"]))
        elif pair == ("sql.NullTime", "*timestamppb.Timestamp"):
            lines.append('\tif !%s.AsTime().Equal(row.%s.Time) {' % (g, f["go"]))
            lines.append(
                '\t\tt.Errorf("%s = %%v, want %%v", %s.AsTime(), row.%s.Time)' % (f["proto"], g, f["go"])
            )
        elif pair == ("int16", "int32"):
            lines.append('\tif %s != int32(row.%s) {' % (g, f["go"]))
            lines.append('\t\tt.Errorf("%s = %%d, want %%d", %s, row.%s)' % (f["proto"], g, f["go"]))
        elif pair == ("uuid.UUID", "string"):
            lines.append('\tif %s != row.%s.String() {' % (g, f["go"]))
            lines.append('\t\tt.Errorf("%s = %%q, want %%q", %s, row.%s.String())' % (f["proto"], g, f["go"]))
        elif pair == ("sql.NullInt64", "*int64"):
            # Сверяется поле, а не геттер: тот отдаёт 0 и на nil, и на явный
            # ноль, и потеря значения выглядела бы как успешный перенос.
            lines.append('\tif got.%s == nil || *got.%s != row.%s.Int64 {' % (f["proto"], f["proto"], f["go"]))
            lines.append(
                '\t\tt.Errorf("%s = %%v, want %%d", got.%s, row.%s.Int64)' % (f["proto"], f["proto"], f["go"])
            )
        else:
            verb = "%q" if f["go_type"] == "string" else ("%d" if f["go_type"] == "int64" else "%v")
            lines.append('\tif %s != row.%s {' % (g, f["go"]))
            lines.append('\t\tt.Errorf("%s = %s, want %s", %s, row.%s)' % (f["proto"], verb, verb, g, f["go"]))
        lines.append('\t}')
        lines.append("")

    # NULL в колонке обязан доехать до клиента как отсутствие поля. Без этой
    # проверки конвертер, подставляющий 0 вместо nil, прошёл бы тесты выше.
    for f in fields:
        if (f["go_type"], f["proto_type"]) != ("sql.NullInt64", "*int64"):
            continue
        lines.append('\trow.%s = sql.NullInt64{}' % f["go"])
        lines.append('\tif got := %sToProto(row); got.%s != nil {' % (conv, f["proto"]))
        lines.append(
            '\t\tt.Errorf("%s = %%v, want nil для SQL NULL", got.%s)' % (f["proto"], f["proto"])
        )
        lines.append('\t}')
        lines.append("")

    lines.append("}")
    lines.append("")
    return lines


def gen_rows_to_proto_test(plural_conv, sample_func, repo, struct, fields):
    """Test<Конвертер>sToProto: длина слайса и пустой (не nil) ответ на пустом входе."""
    first_str = next((f for f in fields if f["go_type"] == "string"), None)
    id_field = next((f for f in fields if f["col"] == "id"), None)
    lines = []
    lines.append("func Test%sToProto(t *testing.T) {" % plural_conv)
    lines.append('\tfirst := %s()' % sample_func)
    lines.append("")
    lines.append('\tsecond := %s()' % sample_func)
    if id_field:
        lines.append('\tsecond.%s = 999' % id_field["go"])
    if first_str:
        lines.append('\tsecond.%s = "second-value"' % first_str["go"])
    lines.append("")
    lines.append('\ttests := []struct {')
    lines.append('\t\tname    string')
    lines.append('\t\tinput   []%s.%s' % (repo, struct))
    lines.append('\t\twantLen int')
    lines.append('\t}{')
    lines.append('\t\t{name: "two rows", input: []%s.%s{first, second}, wantLen: 2},' % (repo, struct))
    lines.append('\t\t{name: "empty slice", input: []%s.%s{}, wantLen: 0},' % (repo, struct))
    lines.append('\t\t{name: "nil slice", input: nil, wantLen: 0},')
    lines.append('\t}')
    lines.append("")
    lines.append('\tfor _, tt := range tests {')
    lines.append('\t\tt.Run(tt.name, func(t *testing.T) {')
    lines.append('\t\t\tgot := %sToProto(tt.input)' % plural_conv)
    lines.append("")
    lines.append('\t\t\tif got == nil {')
    lines.append('\t\t\t\tt.Fatal("%sToProto() = nil, want empty slice")' % plural_conv)
    lines.append('\t\t\t}')
    lines.append("")
    lines.append('\t\t\tif len(got) != tt.wantLen {')
    lines.append('\t\t\t\tt.Fatalf("len = %d, want %d", len(got), tt.wantLen)')
    lines.append('\t\t\t}')
    lines.append('\t\t})')
    lines.append('\t}')
    lines.append("}")
    lines.append("")
    return lines


def gen_params_test_block(proto, repo, req_type, params_type, func_name, fields, ptr_helper=None, time_prefix=None):
    """TestTo<Func>Params + TestTo<Func>ParamsNil — общий шаблон для Create и
    для update/delete/undelete в форме параметров.

    time.Time сравнивается через Equal(), а не через общий "got != want":
    ProtoToTime гоняет значение через timestamppb (Unix-секунды и наносекунды),
    и хотя для дат в разумном диапазоне это даёт побитово тот же time.Time,
    что и литерал time.Date(...), полагаться на такое совпадение — забег
    с граблями. После проверки оба временных поля обнуляются, чтобы остаток
    структуры остался сравним одним "!=", как для всех прочих типов полей.
    """
    time_fields = [f for f in fields if f["go_type"] == "time.Time"]

    def time_var(f):
        return "%sCreatedAt" % time_prefix if f["col"] == "created_at" else "%sUpdatedAt" % time_prefix

    lines = []
    lines.append("func Test%s(t *testing.T) {" % func_name)
    lines.append("\treq := &%s.%s{" % (proto, req_type))
    for i, f in enumerate(fields):
        if f["go_type"] == "time.Time":
            v = "timestamppb.New(%s)" % time_var(f)
        else:
            v = sample_value(f, i, ptr_helper=ptr_helper)
        lines.append("\t\t%s: %s," % (f["proto"], v))
    lines.append("\t}")
    lines.append("")
    lines.append("\twant := %s.%s{" % (repo, params_type))
    for i, f in enumerate(fields):
        if f["go_type"] == "time.Time":
            v = time_var(f)
        else:
            v = repo_sample_value(f, i)
            if f["go_type"] == "int16":
                v = str(10 + i)
        lines.append("\t\t%s: %s," % (f["go"], v))
    lines.append("\t}")
    lines.append("")
    if time_fields:
        lines.append("\tgot := %s(req)" % func_name)
        for f in time_fields:
            lines.append("\tif !got.%s.Equal(want.%s) {" % (f["go"], f["go"]))
            lines.append(
                '\t\tt.Errorf("%s.%s = %%v, want %%v", got.%s, want.%s)' % (func_name, f["go"], f["go"], f["go"])
            )
            lines.append("\t}")
            lines.append("\tgot.%s, want.%s = time.Time{}, time.Time{}" % (f["go"], f["go"]))
        lines.append("")
        lines.append("\tif got != want {")
        lines.append('\t\tt.Errorf("%s() = %%+v, want %%+v", got, want)' % func_name)
        lines.append("\t}")
    else:
        lines.append("\tif got := %s(req); got != want {" % func_name)
        lines.append('\t\tt.Errorf("%s() = %%+v, want %%+v", got, want)' % func_name)
        lines.append("\t}")
    lines.append("}")
    lines.append("")

    lines.append("func Test%sNil(t *testing.T) {" % func_name)
    lines.append("\t// Геттеры protobuf безопасны на nil: сервер не должен падать.")
    lines.append("\tif got := %s(nil); got != (%s.%s{}) {" % (func_name, repo, params_type))
    lines.append('\t\tt.Errorf("%s(nil) = %%+v, want zero value", got)' % func_name)
    lines.append("\t}")
    lines.append("}")
    lines.append("")
    return lines


def gen_converter_test(meta, e):
    proto = meta["proto_alias"]
    repo = meta["repo_alias"]
    E, P = e["entity"], e["plural"]

    paginated_lists = [l for l in e["lists"] if l["paginated"]]
    param_ops = [
        (op_key, e["ops"][op_key])
        for op_key in ("update", "delete", "undelete")
        if op_key in e["ops"] and e["ops"][op_key]["arg_form"] == "params"
    ]

    row_fields = e["row_fields"] + query_row_fields(e)
    extra_fields = list(e["create_fields"])
    for _, op in param_ops:
        extra_fields += op["fields"]

    # time.Time в extra_fields — клиент явно передаёт created_at/updated_at в
    # Update/Delete (см. ProtoToTime); своих CreatedAt/UpdatedAt var'ов у
    # такой таблицы могло бы не быть, даже если в строке колонки есть, но
    # они всегда есть — created_at/updated_at стандартны для всех таблиц.
    has_time = any(f["go_type"] in ("time.Time", "sql.NullTime") for f in row_fields + extra_fields)
    has_nulltime = any(f["go_type"] == "sql.NullTime" for f in row_fields)
    # Необязательные фильтры сравниваются с sql.Null*-литералами и требуют
    # адреса значения для optional-полей proto — отсюда импорт и хелпер.
    optional_filters = [f for f in all_filter_fields(e) if (f["go_type"], f["proto_type"]) in NULLABLE_PAIRS]
    ptr_helper = "%sPtr" % lower_first(E)

    # Указатель нужен не только фильтрам: nullable-колонка (dc.column_cat.alias_id)
    # приезжает в строке и в параметрах Create/Update тем же optional-полем proto.
    needs_ptr = bool(optional_filters) or any(
        (f["go_type"], f["proto_type"]) in NULLABLE_PAIRS for f in row_fields + extra_fields
    )

    # Шапка собирается в конце: у таблицы-справочника (без Create и без
    # запросов с параметрами) proto-пакет в тестах не упоминается вовсе,
    # и неиспользованный импорт не дал бы собрать пакет.
    lines = []

    if needs_ptr:
        lines.append("// %s — адрес значения для optional-полей proto." % ptr_helper)
        lines.append("func %s[T any](v T) *T {" % ptr_helper)
        lines.append("\treturn &v")
        lines.append("}")
        lines.append("")

    if has_time:
        lines.append("var (")
        lines.append(
            "\t%sCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)" % lower_first(E)
        )
        lines.append(
            "\t%sUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)" % lower_first(E)
        )
        lines.append(")")
        lines.append("")

    lines.extend(
        gen_row_sample_block(
            repo, e["row_struct"], "test%sRow" % E, e["row_fields"], lower_first(E), "строка %s" % e["table"]
        )
    )
    lines.extend(gen_row_to_proto_test(E, "test%sRow" % E, e["row_fields"]))
    lines.extend(gen_rows_to_proto_test(P, "test%sRow" % E, repo, e["row_struct"], e["row_fields"]))

    # Собственные строки запросов — свой набор полей, свой конвертер.
    for r in query_rows(e):
        sample = "test%s" % r["go"]
        lines.extend(
            gen_row_sample_block(
                repo, r["go"], sample, r["fields"], lower_first(E), "строка выборки %s" % r["go"]
            )
        )
        lines.extend(gen_row_to_proto_test(r["conv"], sample, r["fields"]))
        if r["slice"]:
            lines.extend(gen_rows_to_proto_test(r["conv_plural"], sample, repo, r["go"], r["fields"]))

    # Create params
    if e["has_create"]:
        lines.extend(
            gen_params_test_block(
                proto, repo, "Create%sRequest" % E, e["create_params"], "ToCreate%sParams" % E,
                e["create_fields"], ptr_helper, lower_first(E),
            )
        )

    # update/delete/undelete в форме параметров
    for op_key, op in param_ops:
        name = op["query"]
        lines.extend(
            gen_params_test_block(
                proto, repo, name + "Request", op["params_type"], "To%sParams" % name, op["fields"],
                ptr_helper, lower_first(E),
            )
        )

    for lk in e["lookups"]:
        proto_lit = '"%s"' % uuid_literal(0, 7) if lk["arg_type"] == "uuid.UUID" else (
            "77" if lk["arg_type"] == "int64" else '"%s-lookup"' % lk["col"].replace("_", "-")
        )
        want = "uuid.MustParse(%s)" % proto_lit if lk["arg_type"] == "uuid.UUID" else proto_lit
        zero = "uuid.Nil" if lk["arg_type"] == "uuid.UUID" else ("0" if lk["arg_type"] == "int64" else '""')
        lines.append("func TestTo%sArg(t *testing.T) {" % lk["query"])
        lines.append("\treq := &%s.%sRequest{%s: %s}" % (proto, lk["query"], lk["proto_field"], proto_lit))
        lines.append("")
        lines.append("\tif got := To%sArg(req); got != %s {" % (lk["query"], want))
        lines.append('\t\tt.Errorf("To%sArg() = %%v, want %%v", got, %s)' % (lk["query"], want))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func TestTo%sArgNil(t *testing.T) {" % lk["query"])
        lines.append("\tif got := To%sArg(nil); got != %s {" % (lk["query"], zero))
        lines.append('\t\tt.Errorf("To%sArg(nil) = %%v, want zero value", got)' % lk["query"])
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    for l in paginated_lists:
        name = l["query"]
        lines.append("func Test%sDefaultsPageLimit(t *testing.T) {" % name)
        lines.append("\tgot := To%sParams(&%s.%sRequest{Page: 3})" % (name, proto, name))
        lines.append("")
        lines.append("\tif got.PageLimit != validation.DefaultPageSize {")
        lines.append('\t\tt.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)')
        lines.append("\t}")
        lines.append("")
        lines.append("\tif got.Page != 3 {")
        lines.append('\t\tt.Errorf("Page = %d, want 3", got.Page)')
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func Test%sDefaultsPage(t *testing.T) {" % name)
        lines.append("\tgot := To%sParams(&%s.%sRequest{PageLimit: 10})" % (name, proto, name))
        lines.append("")
        lines.append("\tif got.Page != 1 {")
        lines.append('\t\tt.Errorf("Page = %d, want 1", got.Page)')
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func Test%sKeepsExplicitPageLimit(t *testing.T) {" % name)
        lines.append("\tgot := To%sParams(&%s.%sRequest{PageLimit: 10, Page: 5})" % (name, proto, name))
        lines.append("")
        lines.append("\tif got.PageLimit != 10 {")
        lines.append('\t\tt.Errorf("PageLimit = %d, want 10", got.PageLimit)')
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        if l.get("filter_fields"):
            lines.append("func Test%sPassesFilterFields(t *testing.T) {" % name)
            lines.append("\treq := &%s.%sRequest{PageLimit: 10, Page: 1," % (proto, name))
            for i, f in enumerate(l["filter_fields"]):
                lines.append("\t\t%s: %s," % (f["proto"], filter_proto_literal(f, i, ptr_helper)))
            lines.append("\t}")
            lines.append("")
            lines.append("\tgot := To%sParams(req)" % name)
            lines.append("")
            for i, f in enumerate(l["filter_fields"]):
                v = filter_repo_literal(f, i)
                verb = "%q" if f["go_type"] == "string" else ("%d" if f["go_type"] in ("int64", "int16") else "%v")
                lines.append("\tif got.%s != %s {" % (f["go"], v))
                lines.append('\t\tt.Errorf("%s = %s, want %s", got.%s, %s)' % (f["go"], verb, verb, f["go"], v))
                lines.append("\t}")
            lines.append("}")
            lines.append("")

        list_optional_filters = [
            f for f in (l.get("filter_fields") or []) if (f["go_type"], f["proto_type"]) in NULLABLE_PAIRS
        ]
        if list_optional_filters:
            lines.append("// Непереданный необязательный фильтр обязан остаться SQL NULL, а не")
            lines.append("// превратиться в 0/false: иначе выборка молча сузилась бы.")
            lines.append("func Test%sOmitsAbsentFilters(t *testing.T) {" % name)
            lines.append("\tgot := To%sParams(&%s.%sRequest{PageLimit: 10, Page: 1})" % (name, proto, name))
            lines.append("")
            for f in list_optional_filters:
                lines.append("\tif got.%s.Valid {" % f["go"])
                lines.append('\t\tt.Errorf("%s = %%+v, want invalid", got.%s)' % (f["go"], f["go"]))
                lines.append("\t}")
            lines.append("}")
            lines.append("")

    body = "\n".join(lines)

    head = ["package %s" % meta["conv_pkg"], "", "import ("]
    # По телу файла, а не по флагам: sql.Null*-литералы приезжают из строки,
    # из параметров и из фильтров, и перечислять источники — лишний шанс
    # выписать неиспользованный импорт, а это ошибка сборки Go.
    if "sql.Null" in body:
        head.append('\t"database/sql"')
    head.append('\t"testing"')
    if has_time:
        head.append('\t"time"')
    head.append("")
    if needs_uuid_pkg(e, extra_fields):
        head.append('\t"github.com/google/uuid"')
    head.append('\t"%s"' % meta["repo_import"])
    if "%s." % proto in body:
        head.append("\t" + meta["proto_import"])
    if paginated_lists:
        head.append('\t"%s"' % VALID_IMPORT)
    if "timestamppb." in body:
        head.append('\t"google.golang.org/protobuf/types/known/timestamppb"')
    head.append(")")
    head.append("")
    head.append("")

    return "\n".join(head) + body


# ---------- validation ----------
def field_validation_line(f, limit_entity):
    col, getter = f["col"], "req.Get%s()" % f["proto"]
    if f["go_type"] == "string":
        if f["varchar"] is None:
            raise Exception("нет границы varchar для %s" % col)
        return '\tv.StringVarchar("%s", %s, %s)' % (col, getter, limit_const(limit_entity, f["proto"]))
    if f["go_type"] == "int64":
        return '\tv.Int64ID("%s", %s)' % (col, getter)
    if f["go_type"] == "sql.NullInt64":
        # Nullable-колонка: не передали — проходит, передали — обязан быть
        # больше нуля. Значение берётся из поля, а не из геттера: тот отдаёт 0
        # и на "не передан", и на явный ноль, и presence потерялась бы.
        return '\tv.OptionalInt64ID("%s", req.%s)' % (col, f["proto"])
    if f["go_type"] == "int16":
        return '\tv.Int32Between("%s", %s, math.MinInt16, math.MaxInt16)' % (col, getter)
    if f["go_type"] == "uuid.UUID":
        return '\tv.StringUUID("%s", %s)' % (col, getter)
    if f["go_type"] == "time.Time":
        # Клиент явно передаёт значение (created_at/updated_at в Update/Delete).
        # AsTime() на nil-Timestamp дал бы unix-эпоху (не нулевое time.Time),
        # и TimeRequired не заметила бы "поле не передано" — поэтому через
        # converter.ProtoToTime, как и в самом конвертере параметров.
        return '\tv.TimeRequired("%s", converter.ProtoToTime(%s))' % (col, getter)
    if f["go_type"] == "bool":
        return None
    raise Exception("нет правила валидации для %s (%s)" % (col, f["go_type"]))


def gen_validate_func(proto, func_name, req_type, fields, limit_entity):
    lines = []
    lines.append("func %s(req *%s.%s) error {" % (func_name, proto, req_type))
    lines.append("\tv := validator.New()")
    lines.append("")
    lines.append("\tif req == nil {")
    lines.append('\t\tv.AddError("request", validator.MsgRequired)')
    lines.append("\t\treturn v.Err()")
    lines.append("\t}")
    lines.append("")
    for f in fields:
        call = field_validation_line(f, limit_entity)
        if call:
            lines.append(call)
    lines.append("")
    lines.append("\treturn v.Err()")
    lines.append("}")
    lines.append("")
    return lines


def gen_validation(meta, e):
    proto = meta["proto_alias"]
    E = e["entity"]

    param_ops = [
        (op_key, e["ops"][op_key])
        for op_key in ("update", "delete", "undelete")
        if op_key in e["ops"] and e["ops"][op_key]["arg_form"] == "params"
    ]
    paginated_lists = [l for l in e["lists"] if l["paginated"]]

    needs_math = any(f["go_type"] == "int16" for f in e["create_fields"]) or any(
        f["go_type"] == "int16" for _, op in param_ops for f in op["fields"]
    )
    needs_conv = any(
        f["go_type"] == "time.Time" for _, op in param_ops for f in op["fields"]
    )

    lines = []
    lines.append("package %s" % meta["conv_pkg"])
    lines.append("")
    lines.append("import (")
    if needs_math:
        lines.append('\t"math"')
        lines.append("")
    lines.append("\t" + meta["proto_import"])
    if needs_conv:
        lines.append('\t"%s"' % CONV_IMPORT)
    if paginated_lists:
        lines.append('\t"%s"' % VALID_IMPORT)
    lines.append('\t"%s"' % VALIDATOR_IMPORT)
    lines.append(")")
    lines.append("")

    if e["has_create"]:
        lines.append("// ValidateCreate%s проверяет запрос на вставку строки %s." % (E, e["table"]))
        lines.extend(
            gen_validate_func(proto, "ValidateCreate%s" % E, "Create%sRequest" % E, e["create_fields"], E)
        )

    op_doc = {"update": "обновление", "delete": "мягкое удаление", "undelete": "обратное удаление"}
    for op_key, op in param_ops:
        name = op["query"]
        lines.append(
            "// Validate%s проверяет запрос на %s строки %s." % (name, op_doc[op_key], e["table"])
        )
        lines.extend(gen_validate_func(proto, "Validate%s" % name, name + "Request", op["fields"], E))

    for lk in e["lookups"]:
        lines.append(
            "// Validate%s проверяет запрос на выборку строки %s по %s."
            % (lk["query"], e["table"], lk["col"])
        )
        lines.append("func Validate%s(req *%s.%sRequest) error {" % (lk["query"], proto, lk["query"]))
        lines.append("\tv := validator.New()")
        lines.append("")
        lines.append("\tif req == nil {")
        lines.append('\t\tv.AddError("request", validator.MsgRequired)')
        lines.append("\t\treturn v.Err()")
        lines.append("\t}")
        lines.append("")
        getter = "req.Get%s()" % lk["proto_field"]
        if lk["arg_type"] == "uuid.UUID":
            call = '\tv.StringUUID("%s", %s)' % (lk["col"], getter)
        elif lk["arg_type"] == "int64":
            call = '\tv.Int64ID("%s", %s)' % (lk["col"], getter)
        else:
            call = '\tv.StringVarchar("%s", %s, %s)' % (lk["col"], getter, limit_const(E, lk["proto_field"]))
        lines.append(call)
        lines.append("")
        lines.append("\treturn v.Err()")
        lines.append("}")
        lines.append("")

    for l in e["lists"]:
        if l["paginated"] or l["arg_form"] != "order":
            continue
        # Список целиком, но с выбором направления сортировки: проверять
        # нечего, кроме самого order — страницы у такого запроса нет.
        name = l["query"]
        lines.append("// Validate%s проверяет запрос списка %s." % (name, e["table"]))
        lines.append("func Validate%s(req *%s.%sRequest) error {" % (name, proto, name))
        lines.append("\tv := validator.New()")
        lines.append("")
        lines.append("\tif req == nil {")
        lines.append('\t\tv.AddError("request", validator.MsgRequired)')
        lines.append("\t\treturn v.Err()")
        lines.append("\t}")
        lines.append("")
        lines.append('\tv.StringIn("order", req.GetOrder(), "", "ASC", "DESC")')
        lines.append("")
        lines.append("\treturn v.Err()")
        lines.append("}")
        lines.append("")

    for l in paginated_lists:
        name = l["query"]
        lines.append("// Validate%s проверяет запрос страницы %s." % (name, e["table"]))
        lines.append("func Validate%s(req *%s.%sRequest) error {" % (name, proto, name))
        lines.append("\tv := validator.New()")
        lines.append("")
        lines.append("\tif req == nil {")
        lines.append('\t\tv.AddError("request", validator.MsgRequired)')
        lines.append("\t\treturn v.Err()")
        lines.append("\t}")
        lines.append("")
        lines.append('\tv.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)')
        lines.append('\tv.Int32Min("page", req.GetPage(), 0)')
        if l["order_field"]:
            lines.append('\tv.StringIn("order", req.GetOrder(), "", "ASC", "DESC")')
        if l.get("order_by_field"):
            # Допустимые значения вычитаны из самого ORDER BY (см. parse.py):
            # неизвестное имя колонки в SQL молча даёт сортировку по умолчанию,
            # поэтому опечатку обязан поймать api-слой.
            allowed = ", ".join('"%s"' % v for v in l["order_by_values"])
            lines.append('\tv.StringIn("order_by", req.GetOrderBy(), "", %s)' % allowed)
        lines.append("")
        lines.append("\treturn v.Err()")
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


def gen_validation_test(meta, e):
    proto = meta["proto_alias"]
    E = e["entity"]

    param_ops = [
        (op_key, e["ops"][op_key])
        for op_key in ("update", "delete", "undelete")
        if op_key in e["ops"] and e["ops"][op_key]["arg_form"] == "params"
    ]
    paginated_lists = [l for l in e["lists"] if l["paginated"]]

    ptr_helper = "%sPtr" % lower_first(E)
    needs_ptr = False

    def literal(f, i):
        nonlocal needs_ptr
        if f["go_type"] == "string":
            return '"%s-%d"' % (f["col"].replace("_", "-"), i)
        if f["go_type"] == "int64":
            return str(100 + i)
        if f["go_type"] == "int16":
            return str(10 + i)
        if f["go_type"] == "bool":
            return "true"
        if f["go_type"] == "uuid.UUID":
            return '"%s"' % uuid_literal(i)
        if f["go_type"] == "sql.NullInt64":
            needs_ptr = True
            return "%s(int64(%d))" % (ptr_helper, 100 + i)
        if f["go_type"] == "time.Time":
            return "timestamppb.New(time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC))"
        # Прежняя заглушка "0" молча подставляла неверный литерал вместо того,
        # чтобы упасть: для optional-поля он к тому же не компилировался.
        raise Exception("нет тестового литерала для %s (%s)" % (f["col"], f["go_type"]))

    lines = []
    lines.append("// %sFieldErrors достаёт из ошибки список полей с претензиями." % lower_first(E))
    lines.append("func %sFieldErrors(t *testing.T, err error) map[string][]string {" % lower_first(E))
    lines.append("\tt.Helper()")
    lines.append("")
    lines.append("\tvar validationErr *validator.ValidationError")
    lines.append("\tif !errors.As(err, &validationErr) {")
    lines.append('\t\tt.Fatalf("error = %v, want *validator.ValidationError", err)')
    lines.append("\t}")
    lines.append("")
    lines.append("\treturn validationErr.Errors")
    lines.append("}")
    lines.append("")

    def gen_cases(fields, reqtype, validate_fn):
        out = []
        out.append("func valid%sRequest() *%s.%s {" % (validate_fn, proto, reqtype))
        out.append("\treturn &%s.%s{" % (proto, reqtype))
        for i, f in enumerate(fields):
            out.append("\t\t%s: %s," % (f["proto"], literal(f, i)))
        out.append("\t}")
        out.append("}")
        out.append("")

        out.append("func Test%s(t *testing.T) {" % validate_fn)
        out.append("\ttests := []struct {")
        out.append("\t\tname      string")
        out.append("\t\tmutate    func(*%s.%s)" % (proto, reqtype))
        out.append("\t\twantField string")
        out.append("\t}{")
        out.append('\t\t{name: "valid", mutate: func(*%s.%s) {}},' % (proto, reqtype))
        for f in fields:
            col = f["col"]
            if f["go_type"] == "string":
                const = limit_const(E, f["proto"])
                out.append(
                    '\t\t{name: "empty %s", mutate: func(r *%s.%s) { r.%s = "" }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
                out.append(
                    '\t\t{name: "%s too long", mutate: func(r *%s.%s) { r.%s = strings.Repeat("a", %s+1) }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], const, col)
                )
            elif f["go_type"] == "int64":
                out.append(
                    '\t\t{name: "zero %s", mutate: func(r *%s.%s) { r.%s = 0 }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
            elif f["go_type"] == "int16":
                out.append(
                    '\t\t{name: "%s over smallint", mutate: func(r *%s.%s) { r.%s = math.MaxInt16 + 1 }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
            elif f["go_type"] == "uuid.UUID":
                out.append(
                    '\t\t{name: "empty %s", mutate: func(r *%s.%s) { r.%s = "" }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
                out.append(
                    '\t\t{name: "malformed %s", mutate: func(r *%s.%s) { r.%s = "not-a-uuid" }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
            elif f["go_type"] == "sql.NullInt64":
                # Два случая, и оба важны: непереданное значение обязано
                # проходить (иначе колонка не необязательная), а переданный
                # ноль — отбиваться, как у обычного FK.
                out.append(
                    '\t\t{name: "absent %s", mutate: func(r *%s.%s) { r.%s = nil }},'
                    % (col, proto, reqtype, f["proto"])
                )
                out.append(
                    '\t\t{name: "zero %s", mutate: func(r *%s.%s) { r.%s = %s(int64(0)) }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], ptr_helper, col)
                )
            elif f["go_type"] == "time.Time":
                out.append(
                    '\t\t{name: "nil %s", mutate: func(r *%s.%s) { r.%s = nil }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
        out.append("\t}")
        out.append("")
        out.append("\tfor _, tt := range tests {")
        out.append("\t\tt.Run(tt.name, func(t *testing.T) {")
        out.append("\t\t\treq := valid%sRequest()" % validate_fn)
        out.append("\t\t\ttt.mutate(req)")
        out.append("")
        out.append("\t\t\terr := %s(req)" % validate_fn)
        out.append("")
        out.append('\t\t\tif tt.wantField == "" {')
        out.append("\t\t\t\tif err != nil {")
        out.append('\t\t\t\t\tt.Fatalf("%s() = %%v, want nil", err)' % validate_fn)
        out.append("\t\t\t\t}")
        out.append("\t\t\t\treturn")
        out.append("\t\t\t}")
        out.append("")
        out.append("\t\t\tif err == nil {")
        out.append('\t\t\t\tt.Fatalf("%s() = nil, want error on %%q", tt.wantField)' % validate_fn)
        out.append("\t\t\t}")
        out.append("")
        out.append("\t\t\tfields := %sFieldErrors(t, err)" % lower_first(E))
        out.append("\t\t\tif len(fields[tt.wantField]) == 0 {")
        out.append('\t\t\t\tt.Errorf("no error on %q, got %v", tt.wantField, fields)')
        out.append("\t\t\t}")
        out.append("\t\t})")
        out.append("\t}")
        out.append("}")
        out.append("")

        out.append("func Test%sNil(t *testing.T) {" % validate_fn)
        out.append("\tif err := %s(nil); err == nil {" % validate_fn)
        out.append('\t\tt.Error("%s(nil) = nil, want error")' % validate_fn)
        out.append("\t}")
        out.append("}")
        out.append("")
        return out

    if e["has_create"]:
        lines.extend(gen_cases(e["create_fields"], "Create%sRequest" % E, "ValidateCreate%s" % E))

    for op_key, op in param_ops:
        name = op["query"]
        lines.extend(gen_cases(op["fields"], name + "Request", "Validate%s" % name))

    for lk in e["lookups"]:
        good = (
            '"%s"' % uuid_literal(0, 7)
            if lk["arg_type"] == "uuid.UUID"
            else ("77" if lk["arg_type"] == "int64" else '"%s-lookup"' % lk["col"].replace("_", "-"))
        )
        if lk["arg_type"] == "uuid.UUID":
            bad = [("empty", '""'), ("malformed", '"not-a-uuid"')]
        elif lk["arg_type"] == "int64":
            bad = [("zero", "0")]
        else:
            bad = [("empty", '""')]

        lines.append("func TestValidate%s(t *testing.T) {" % lk["query"])
        lines.append("\ttests := []struct {")
        lines.append("\t\tname    string")
        lines.append("\t\tvalue   %s" % lk["proto_type"])
        lines.append("\t\twantErr bool")
        lines.append("\t}{")
        lines.append('\t\t{name: "valid", value: %s},' % good)
        for name_, value in bad:
            lines.append('\t\t{name: "%s", value: %s, wantErr: true},' % (name_, value))
        lines.append("\t}")
        lines.append("")
        lines.append("\tfor _, tt := range tests {")
        lines.append("\t\tt.Run(tt.name, func(t *testing.T) {")
        lines.append(
            "\t\t\terr := Validate%s(&%s.%sRequest{%s: tt.value})"
            % (lk["query"], proto, lk["query"], lk["proto_field"])
        )
        lines.append("")
        lines.append("\t\t\tif (err != nil) != tt.wantErr {")
        lines.append(
            '\t\t\t\tt.Errorf("Validate%s() error = %%v, wantErr %%v", err, tt.wantErr)' % lk["query"]
        )
        lines.append("\t\t\t}")
        lines.append("\t\t})")
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func TestValidate%sNil(t *testing.T) {" % lk["query"])
        lines.append("\tif err := Validate%s(nil); err == nil {" % lk["query"])
        lines.append('\t\tt.Error("Validate%s(nil) = nil, want error")' % lk["query"])
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    for l in e["lists"]:
        if l["paginated"] or l["arg_form"] != "order":
            continue
        name = l["query"]
        lines.append("func TestValidate%s(t *testing.T) {" % name)
        lines.append("\ttests := []struct {")
        lines.append("\t\tname    string")
        lines.append("\t\torder   string")
        lines.append("\t\twantErr bool")
        lines.append("\t}{")
        lines.append('\t\t{name: "empty order", order: ""},')
        lines.append('\t\t{name: "ASC", order: "ASC"},')
        lines.append('\t\t{name: "DESC", order: "DESC"},')
        lines.append('\t\t{name: "invalid order", order: "sideways", wantErr: true},')
        lines.append("\t}")
        lines.append("")
        lines.append("\tfor _, tt := range tests {")
        lines.append("\t\tt.Run(tt.name, func(t *testing.T) {")
        lines.append(
            "\t\t\terr := Validate%s(&%s.%sRequest{Order: tt.order})" % (name, proto, name)
        )
        lines.append("\t\t\tif (err != nil) != tt.wantErr {")
        lines.append(
            '\t\t\t\tt.Errorf("Validate%s() error = %%v, wantErr %%v", err, tt.wantErr)' % name
        )
        lines.append("\t\t\t}")
        lines.append("\t\t})")
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func TestValidate%sNil(t *testing.T) {" % name)
        lines.append("\tif err := Validate%s(nil); err == nil {" % name)
        lines.append('\t\tt.Error("Validate%s(nil) = nil, want error")' % name)
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    for l in paginated_lists:
        name = l["query"]
        lines.append("func TestValidate%s(t *testing.T) {" % name)
        lines.append("\ttests := []struct {")
        lines.append("\t\tname      string")
        lines.append("\t\treq       *%s.%sRequest" % (proto, name))
        lines.append("\t\twantErr   bool")
        lines.append("\t}{")
        lines.append('\t\t{name: "valid", req: &%s.%sRequest{PageLimit: 50, Page: 1}},' % (proto, name))
        lines.append(
            '\t\t{name: "zero page_limit and page ok", req: &%s.%sRequest{PageLimit: 0, Page: 0}},' % (proto, name)
        )
        lines.append(
            '\t\t{name: "negative page_limit", req: &%s.%sRequest{PageLimit: -1}, wantErr: true},' % (proto, name)
        )
        lines.append(
            '\t\t{name: "negative page", req: &%s.%sRequest{Page: -1}, wantErr: true},' % (proto, name)
        )
        if l["order_field"]:
            lines.append(
                '\t\t{name: "invalid order", req: &%s.%sRequest{Order: "sideways"}, wantErr: true},' % (proto, name)
            )
        if l.get("order_by_field"):
            lines.append(
                '\t\t{name: "valid order_by", req: &%s.%sRequest{OrderBy: "%s"}},'
                % (proto, name, l["order_by_values"][0])
            )
            lines.append(
                '\t\t{name: "unknown order_by", req: &%s.%sRequest{OrderBy: "nonexistent"}, wantErr: true},'
                % (proto, name)
            )
        lines.append("\t}")
        lines.append("")
        lines.append("\tfor _, tt := range tests {")
        lines.append("\t\tt.Run(tt.name, func(t *testing.T) {")
        lines.append("\t\t\terr := Validate%s(tt.req)" % name)
        lines.append("\t\t\tif (err != nil) != tt.wantErr {")
        lines.append(
            '\t\t\t\tt.Errorf("Validate%s() error = %%v, wantErr %%v", err, tt.wantErr)' % name
        )
        lines.append("\t\t\t}")
        lines.append("\t\t})")
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func TestValidate%sNil(t *testing.T) {" % name)
        lines.append("\tif err := Validate%s(nil); err == nil {" % name)
        lines.append('\t\tt.Error("Validate%s(nil) = nil, want error")' % name)
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    body = "\n".join(lines)

    if needs_ptr:
        body = (
            "// %s — адрес значения для optional-полей proto.\n"
            "func %s[T any](v T) *T {\n"
            "\treturn &v\n"
            "}\n\n" % (ptr_helper, ptr_helper)
        ) + body

    head = ["package %s" % meta["conv_pkg"], "", "import (", '\t"errors"']
    if "math.MaxInt16" in body:
        head.append('\t"math"')
    if "strings.Repeat(" in body:
        head.append('\t"strings"')
    head.append('\t"testing"')
    if "time.Date(" in body:
        head.append('\t"time"')
    head.append("")
    head.append("\t" + meta["proto_import"])
    if "timestamppb." in body:
        head.append('\t"google.golang.org/protobuf/types/known/timestamppb"')
    head.append('\t"%s"' % VALIDATOR_IMPORT)
    head.append(")")
    head.append("")
    head.append("")

    return "\n".join(head) + body


# ---------- limits ----------
def gen_limits(meta, entities):
    lines = []
    lines.append("package %s" % meta["conv_pkg"])
    lines.append("")
    lines.append("// Ограничения длины строковых колонок. Значения взяты из")
    lines.append("// datacatalogue/db/sqlc/%s/schema.sql и должны меняться вместе" % meta["sqlc"])
    lines.append("// с ней: валидация обязана отсекать слишком длинную строку раньше,")
    lines.append("// чем её отвергнет Postgres.")
    lines.append("const (")
    first = True
    for e in entities:
        seen = {}
        for f in all_param_fields(e):
            if f["go_type"] == "string" and f["col"] not in seen:
                seen[f["col"]] = f
        for lk in e["lookups"]:
            if lk["arg_type"] == "string" and lk["col"] not in seen:
                seen[lk["col"]] = {"proto": lk["proto_field"], "varchar": lk["varchar"], "col": lk["col"]}
        if not seen:
            continue
        if not first:
            lines.append("")
        first = False
        lines.append("\t// %s" % e["table"])
        for col, f in seen.items():
            lines.append(
                "\t%s = %d // %s character varying(%d)"
                % (limit_const(e["entity"], f["proto"]), f["varchar"], col, f["varchar"])
            )
    lines.append(")")
    lines.append("")
    return "\n".join(lines)


# ---------- service ----------
def gen_service(meta, e):
    repo = meta["repo_alias"]
    field = meta["repo_field"]
    recv = "s"
    st = meta["service_type"]
    E, P = e["entity"], e["plural"]
    row = "%s.%s" % (repo, e["row_struct"])
    low = e["table"]
    ops = e["ops"]

    # Шапка собирается в конце: у таблицы без Get<E>ById в файле не остаётся ни
    # одной проверки sql.ErrNoRows, и неиспользованный импорт не дал бы собрать пакет.
    lines = []

    def row_type(r):
        return "%s.%s" % (repo, r["go"])

    def one(name, doc, notfound_msg, row):
        out = []
        out.append("// %s %s" % (name, doc))
        out.append("func (%s *%s) %s(ctx context.Context, id int64) (%s, error) {" % (recv, st, name, row))
        out.append("\trow, err := %s.%s.%s(ctx, id)" % (recv, field, name))
        out.append("")
        out.append("\tif err != nil {")
        out.append("\t\tif errors.Is(err, sql.ErrNoRows) {")
        out.append(
            '\t\t\treturn %s{}, fmt.Errorf("%s id = %%d: %%w", id, customerrors.ErrNotFound)'
            % (row, notfound_msg)
        )
        out.append("\t\t}")
        out.append("")
        out.append('\t\treturn %s{}, fmt.Errorf("get %s id = %%d: %%w", id, err)' % (row, notfound_msg))
        out.append("\t}")
        out.append("")
        out.append("\treturn row, nil")
        out.append("}")
        out.append("")
        return out

    if "get_by_id" in ops:
        lines.extend(
            one(
                "Get%sById" % E,
                "возвращает активную строку %s по id." % low,
                low,
                row_type(ops["get_by_id"]["row"]),
            )
        )
    if "get_deleted_by_id" in ops:
        lines.extend(
            one(
                "GetDeleted%sById" % E,
                "возвращает мягко удалённую строку %s по id." % low,
                "deleted " + low,
                row_type(ops["get_deleted_by_id"]["row"]),
            )
        )

    for l in e["lists"]:
        name = l["query"]
        row = row_type(l["row"])
        if not l["paginated"]:
            ordered = l["arg_form"] == "order"
            lines.append("// %s возвращает строки %s." % (name, low))
            lines.append(
                "func (%s *%s) %s(ctx context.Context%s) ([]%s, error) {"
                % (recv, st, name, ", order string" if ordered else "", row)
            )
            lines.append("\trows, err := %s.%s.%s(ctx%s)" % (recv, field, name, ", order" if ordered else ""))
            lines.append("")
            lines.append("\tif err != nil {")
            lines.append('\t\treturn nil, fmt.Errorf("%s: %%w", err)' % name)
            lines.append("\t}")
            lines.append("")
            lines.append("\treturn rows, nil")
            lines.append("}")
            lines.append("")
            continue

        counter_row = "%s.%s" % (repo, "%sRow" % l["counter"])
        lines.append("// %s возвращает страницу строк %s и её счётчики." % (name, low))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, params %s.%s) ([]%s, %s, error) {"
            % (recv, st, name, repo, l["params_type"], row, counter_row)
        )
        if l["counter_arg_form"] == "int32":
            lines.append("\tcount, err := %s.%s.%s(ctx, params.PageLimit)" % (recv, field, l["counter"]))
        else:
            lines.append("\tcountParams := %s.%s{" % (repo, l["counter_params_type"]))
            lines.append("\t\tPageLimit: params.PageLimit,")
            for f in l.get("filter_fields") or []:
                lines.append("\t\t%s: params.%s," % (f["go"], f["go"]))
            lines.append("\t}")
            lines.append("\tcount, err := %s.%s.%s(ctx, countParams)" % (recv, field, l["counter"]))
        lines.append("\tif err != nil {")
        lines.append('\t\treturn nil, %s{}, fmt.Errorf("count %s: %%w", err)' % (counter_row, low))
        lines.append("\t}")
        lines.append("")
        lines.append("\tif count.%s == 0 {" % l["counter_total_field"])
        lines.append("\t\treturn []%s{}, count, nil" % row)
        lines.append("\t}")
        lines.append("")
        lines.append("\trows, err := %s.%s.%s(ctx, params)" % (recv, field, name))
        lines.append("\tif err != nil {")
        lines.append('\t\treturn nil, %s{}, fmt.Errorf("get %s page: %%w", err)' % (counter_row, low))
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn rows, count, nil")
        lines.append("}")
        lines.append("")

    for fk in e["fks"]:
        row = row_type(fk["row"])
        lines.append("// %s возвращает активные строки %s, отобранные по %s." % (fk["query"], low, fk["col"]))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, %s %s) ([]%s, error) {"
            % (recv, st, fk["query"], fk["arg_name"], fk["arg_type"], row)
        )
        lines.append("\trows, err := %s.%s.%s(ctx, %s)" % (recv, field, fk["query"], fk["arg_name"]))
        lines.append("")
        lines.append("\tif err != nil {")
        lines.append(
            '\t\treturn nil, fmt.Errorf("get %s by %s = %%v: %%w", %s, err)' % (low, fk["col"], fk["arg_name"])
        )
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn rows, nil")
        lines.append("}")
        lines.append("")

    for lk in e["lookups"]:
        row = row_type(lk["row"])
        verb = "%d" if lk["arg_type"] == "int64" else ("%q" if lk["arg_type"] == "string" else "%v")
        lines.append(
            "// %s возвращает активную строку %s по уникальной колонке %s." % (lk["query"], low, lk["col"])
        )
        lines.append(
            "func (%s *%s) %s(ctx context.Context, %s %s) (%s, error) {"
            % (recv, st, lk["query"], lk["arg_name"], lk["arg_type"], row)
        )
        lines.append("\trow, err := %s.%s.%s(ctx, %s)" % (recv, field, lk["query"], lk["arg_name"]))
        lines.append("")
        lines.append("\tif err != nil {")
        lines.append("\t\tif errors.Is(err, sql.ErrNoRows) {")
        lines.append(
            '\t\t\treturn %s{}, fmt.Errorf("%s %s = %s: %%w", %s, customerrors.ErrNotFound)'
            % (row, low, lk["col"], verb, lk["arg_name"])
        )
        lines.append("\t\t}")
        lines.append("")
        lines.append(
            '\t\treturn %s{}, fmt.Errorf("get %s %s = %s: %%w", %s, err)' % (row, low, lk["col"], verb, lk["arg_name"])
        )
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn row, nil")
        lines.append("}")
        lines.append("")

    # Create и Update/Delete/Undelete с RETURNING отдают строку своей таблицы.
    row = "%s.%s" % (repo, e["row_struct"])

    # Create — у таблицы-справочника его нет
    if e["has_create"]:
        lines.append("// Create%s вставляет строку %s и возвращает её целиком." % (E, low))
        lines.append(
            "func (%s *%s) Create%s(ctx context.Context, params %s.%s) (%s, error) {"
            % (recv, st, E, repo, e["create_params"], row)
        )
        lines.append("\trow, err := %s.%s.Create%s(ctx, params)" % (recv, field, E))
        lines.append("")
        lines.append("\tif err != nil {")
        lines.append('\t\treturn %s{}, fmt.Errorf("%%w: %s: %%w", customerrors.ErrCreate, err)' % (row, low))
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn row, nil")
        lines.append("}")
        lines.append("")

    def id_expr(op):
        if op["arg_form"] == "id":
            return "id"
        return "params.ID"

    def arg_decl(op, params_type):
        if op["arg_form"] == "id":
            return "id int64"
        return "params %s.%s" % (repo, params_type)

    if "update" in ops:
        op = ops["update"]
        name = op["query"]
        if op["resp_form"] == "row":
            lines.append("// %s обновляет активную строку %s и возвращает её целиком." % (name, low))
            lines.append("//")
            lines.append("// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую")
            lines.append("// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,")
            lines.append("// чтобы api-слой ответил NotFound, а не Internal.")
            lines.append(
                "func (%s *%s) %s(ctx context.Context, %s) (%s, error) {"
                % (recv, st, name, arg_decl(op, op["params_type"]), row)
            )
            lines.append("\trow, err := %s.%s.%s(ctx, params)" % (recv, field, name))
            lines.append("")
            lines.append("\tif err != nil {")
            lines.append("\t\tif errors.Is(err, sql.ErrNoRows) {")
            lines.append(
                '\t\t\treturn %s{}, fmt.Errorf("%s id = %%d: %%w", %s, customerrors.ErrNotFound)'
                % (row, low, id_expr(op))
            )
            lines.append("\t\t}")
            lines.append("")
            lines.append(
                '\t\treturn %s{}, fmt.Errorf("%%w: %s id = %%d: %%w", customerrors.ErrUpdate, %s, err)'
                % (row, low, id_expr(op))
            )
            lines.append("\t}")
            lines.append("")
            lines.append("\treturn row, nil")
            lines.append("}")
            lines.append("")
        else:
            lines.append("// %s обновляет строку %s." % (name, low))
            if "get_by_id" in ops:
                lines.append("//")
                lines.append("// Запрос — :exec и не сообщает число затронутых строк, поэтому")
                lines.append("// существование активной записи проверяется заранее.")
            lines.append(
                "func (%s *%s) %s(ctx context.Context, %s) error {"
                % (recv, st, name, arg_decl(op, op["params_type"]))
            )
            if "get_by_id" in ops:
                lines.append("\tif _, err := %s.Get%sById(ctx, %s); err != nil {" % (recv, E, id_expr(op)))
                lines.append("\t\treturn fmt.Errorf(\"%%w: %s: %%w\", customerrors.ErrUpdate, err)" % low)
                lines.append("\t}")
                lines.append("")
            lines.append("\tif err := %s.%s.%s(ctx, params); err != nil {" % (recv, field, name))
            lines.append(
                '\t\treturn fmt.Errorf("%%w: %s id = %%d: %%w", customerrors.ErrUpdate, %s, err)'
                % (low, id_expr(op))
            )
            lines.append("\t}")
            lines.append("")
            lines.append("\treturn nil")
            lines.append("}")
            lines.append("")

    for op_key, prefix, err_name, precheck_slot in (
        ("delete", "Delete", "ErrDelete", "get_by_id"),
        ("undelete", "Undelete", "ErrUndelete", "get_deleted_by_id"),
    ):
        if op_key not in ops:
            continue
        op = ops[op_key]
        name = op["query"]
        lines.append("// %s %s строку %s." % (name, "мягко удаляет" if op_key == "delete" else "восстанавливает мягко удалённую", low))
        if precheck_slot in ops:
            lines.append("//")
            lines.append(
                "// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута"
                if op_key == "delete"
                else "// Существование удалённой записи проверяется заранее по той же причине,"
            )
            if op_key == "delete":
                lines.append("// строка, поэтому существование активной записи проверяем заранее —")
                lines.append("// иначе удаление несуществующего id молча возвращало бы успех.")
            else:
                lines.append("// что и в Delete%sById." % E)
        lines.append(
            "func (%s *%s) %s(ctx context.Context, %s) error {"
            % (recv, st, name, arg_decl(op, op.get("params_type")))
        )
        if precheck_slot in ops:
            precheck_call = "Get%sById" % E if precheck_slot == "get_by_id" else "GetDeleted%sById" % E
            lines.append("\tif _, err := %s.%s(ctx, %s); err != nil {" % (recv, precheck_call, id_expr(op)))
            lines.append("\t\treturn errors.Join(customerrors.%s, err)" % err_name)
            lines.append("\t}")
            lines.append("")
        call_arg = "params" if op["arg_form"] == "params" else "id"
        lines.append("\tif err := %s.%s.%s(ctx, %s); err != nil {" % (recv, field, name, call_arg))
        lines.append(
            '\t\treturn fmt.Errorf("%%w: %s id = %%d: %%w", customerrors.%s, %s, err)'
            % (low, err_name, id_expr(op))
        )
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn nil")
        lines.append("}")
        lines.append("")

    body = "\n".join(lines)

    head = ["package %s" % meta["service_pkg"], "", "import (", '\t"context"']
    if "sql.ErrNoRows" in body:
        head.append('\t"database/sql"')
    if "errors.Is(" in body or "errors.Join(" in body:
        head.append('\t"errors"')
    head.append('\t"fmt"')
    head.append("")
    if any(lk["arg_type"] == "uuid.UUID" for lk in e["lookups"]):
        head.append('\t"github.com/google/uuid"')
    head.append('\t"%s"' % meta["repo_import"])
    if "customerrors." in body:
        head.append('\tcustomerrors "%s"' % ERRORS_IMPORT)
    head.append(")")
    head.append("")
    head.append("")

    return "\n".join(head) + body


# ---------- api ----------
def gen_api(meta, e):
    proto = meta["proto_alias"]
    recv = meta["api_recv"]
    at = meta["api_type"]
    svc = meta["service_field"]
    E, P = e["entity"], e["plural"]
    conv = meta["conv_pkg"] + "conv"
    valid = meta["conv_pkg"] + "validation"
    repo = meta["repo_alias"]
    ops = e["ops"]

    # validation.ValidateID(bare пакет) нужен только там, где аргумент — голый id.
    needs_bare_validation = (
        "get_by_id" in ops
        or "get_deleted_by_id" in ops
        or bool(e["fks"])
        or any(ops[k]["arg_form"] == "id" for k in ("delete", "undelete") if k in ops)
    )
    # emptypb.Empty{} нужен только для ответов-обёрток (не для буквально пустых).
    needs_emptypb = any(
        ops[k]["resp_form"] == "field" for k in ("update", "delete", "undelete") if k in ops
    )

    lines = []
    lines.append("package %s" % meta["api_pkg"])
    lines.append("")
    lines.append("import (")
    lines.append('\t"context"')
    lines.append("")
    lines.append('\t"%s"' % APIERR_IMPORT)
    lines.append('\t%s "%s/%s"' % (conv, CONV_IMPORT, meta["conv_pkg"]))
    if needs_bare_validation:
        lines.append('\t"%s"' % VALID_IMPORT)
    lines.append('\t%s "%s/%s"' % (valid, VALID_IMPORT, meta["conv_pkg"]))
    lines.append("\t" + meta["proto_import"])
    if needs_emptypb:
        lines.append('\t"google.golang.org/protobuf/types/known/emptypb"')
    lines.append(")")
    lines.append("")

    def by_id(rpc, doc, svc_call, resp_field, row_conv):
        out = []
        out.append("// %s %s" % (rpc, doc))
        out.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, rpc, proto, rpc, proto, rpc)
        )
        out.append("\tif err := validation.ValidateID(req.GetId()); err != nil {")
        out.append("\t\treturn nil, apierror.Wrap(err)")
        out.append("\t}")
        out.append("")
        out.append("\trow, err := %s.services.%s.%s(ctx, req.GetId())" % (recv, svc, svc_call))
        out.append("\tif err != nil {")
        out.append("\t\treturn nil, apierror.Wrap(err)")
        out.append("\t}")
        out.append("")
        out.append(
            "\treturn &%s.%sResponse{%s: %s.%sToProto(row)}, nil" % (proto, rpc, resp_field, conv, row_conv)
        )
        out.append("}")
        out.append("")
        return out

    if "get_by_id" in ops:
        lines.extend(
            by_id(
                "Get%sById" % E,
                "отдаёт активную строку %s по id." % e["table"],
                "Get%sById" % E,
                ops["get_by_id"]["resp_field"],
                ops["get_by_id"]["row"]["conv"],
            )
        )
    if "get_deleted_by_id" in ops:
        lines.extend(
            by_id(
                "GetDeleted%sById" % E,
                "отдаёт мягко удалённую строку %s по id." % e["table"],
                "GetDeleted%sById" % E,
                ops["get_deleted_by_id"]["resp_field"],
                ops["get_deleted_by_id"]["row"]["conv"],
            )
        )

    for l in e["lists"]:
        name = l["query"]
        if not l["paginated"]:
            ordered = l["arg_form"] == "order"
            lines.append("// %s отдаёт строки %s." % (name, e["table"]))
            lines.append(
                "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
                % (recv, at, name, proto, name, proto, name)
            )
            if ordered:
                lines.append("\tif err := %s.Validate%s(req); err != nil {" % (valid, name))
                lines.append("\t\treturn nil, apierror.Wrap(err)")
                lines.append("\t}")
                lines.append("")
            lines.append(
                "\trows, err := %s.services.%s.%s(ctx%s)"
                % (recv, svc, name, ", req.GetOrder()" if ordered else "")
            )
            lines.append("\tif err != nil {")
            lines.append("\t\treturn nil, apierror.Wrap(err)")
            lines.append("\t}")
            lines.append("")
            lines.append(
                "\treturn &%s.%sResponse{%s: %s.%sToProto(rows)}, nil"
                % (proto, name, l["resp_field"], conv, l["row"]["conv_plural"])
            )
            lines.append("}")
            lines.append("")
            continue

        lines.append("// %s отдаёт страницу строк %s." % (name, e["table"]))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, name, proto, name, proto, name)
        )
        lines.append("\tif err := %s.Validate%s(req); err != nil {" % (valid, name))
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append("\tparams := %s.To%sParams(req)" % (conv, name))
        lines.append("")
        lines.append(
            "\trows, count, err := %s.services.%s.%s(ctx, params)" % (recv, svc, name)
        )
        lines.append("\tif err != nil {")
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        pm = l["pagination_msg_fields"]
        lines.append("\treturn &%s.%sResponse{" % (proto, name))
        lines.append("\t\t%s: %s.%sToProto(rows)," % (l["resp_field"], conv, l["row"]["conv_plural"]))
        lines.append("\t\t%s: &%s.Pagination{" % (l["pagination_field"], proto))
        lines.append("\t\t\t%s: params.Page," % pm["page"])
        lines.append("\t\t\t%s: params.PageLimit," % pm["per_page"])
        lines.append("\t\t\t%s: count.%s," % (pm["total_items"], l["counter_total_field"]))
        lines.append("\t\t\t%s: count.%s," % (pm["total_pages"], l["counter_pages_field"]))
        lines.append("\t\t},")
        lines.append("\t}, nil")
        lines.append("}")
        lines.append("")

    for fk in e["fks"]:
        rpc = fk["query"]
        lines.append("// %s отдаёт активные строки %s, отобранные по %s." % (rpc, e["table"], fk["col"]))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, rpc, proto, rpc, proto, rpc)
        )
        lines.append("\tif err := validation.ValidateID(req.Get%s()); err != nil {" % fk["proto_field"])
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\trows, err := %s.services.%s.%s(ctx, req.Get%s())" % (recv, svc, rpc, fk["proto_field"])
        )
        lines.append("\tif err != nil {")
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\treturn &%s.%sResponse{%s: %s.%sToProto(rows)}, nil"
            % (proto, rpc, fk["resp_field"], conv, fk["row"]["conv_plural"])
        )
        lines.append("}")
        lines.append("")

    for lk in e["lookups"]:
        rpc = lk["query"]
        lines.append(
            "// %s отдаёт активную строку %s по уникальной колонке %s." % (rpc, e["table"], lk["col"])
        )
        lines.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, rpc, proto, rpc, proto, rpc)
        )
        lines.append("\tif err := %s.Validate%s(req); err != nil {" % (valid, rpc))
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\trow, err := %s.services.%s.%s(ctx, %s.To%sArg(req))" % (recv, svc, rpc, conv, rpc)
        )
        lines.append("\tif err != nil {")
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\treturn &%s.%sResponse{%s: %s.%sToProto(row)}, nil"
            % (proto, rpc, lk["resp_field"], conv, lk["row"]["conv"])
        )
        lines.append("}")
        lines.append("")

    # Create — у таблицы-справочника его нет
    if e["has_create"]:
        lines.append("// Create%s вставляет строку %s и отдаёт её целиком." % (E, e["table"]))
        lines.append(
            "func (%s *%s) Create%s(ctx context.Context, req *%s.Create%sRequest) (*%s.Create%sResponse, error) {"
            % (recv, at, E, proto, E, proto, E)
        )
        lines.append("\tif err := %s.ValidateCreate%s(req); err != nil {" % (valid, E))
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\trow, err := %s.services.%s.Create%s(ctx, %s.ToCreate%sParams(req))" % (recv, svc, E, conv, E)
        )
        lines.append("\tif err != nil {")
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\treturn &%s.Create%sResponse{%s: %s.%sToProto(row)}, nil"
            % (proto, E, e["create_resp_field"], conv, E)
        )
        lines.append("}")
        lines.append("")

    def empty_return(resp_type, resp_field, resp_form):
        if resp_form == "empty":
            return "&%s.%s{}" % (proto, resp_type)
        return "&%s.%s{%s: &emptypb.Empty{}}" % (proto, resp_type, resp_field)

    if "update" in ops:
        op = ops["update"]
        name = op["query"]
        lines.append("// %s обновляет строку %s." % (name, e["table"]))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, name, proto, name, proto, name)
        )
        lines.append("\tif err := %s.Validate%s(req); err != nil {" % (valid, name))
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        if op["resp_form"] == "row":
            lines.append(
                "\trow, err := %s.services.%s.%s(ctx, %s.To%sParams(req))" % (recv, svc, name, conv, name)
            )
            lines.append("\tif err != nil {")
            lines.append("\t\treturn nil, apierror.Wrap(err)")
            lines.append("\t}")
            lines.append("")
            lines.append(
                "\treturn &%s.%sResponse{%s: %s.%sToProto(row)}, nil"
                % (proto, name, op["resp_field"], conv, E)
            )
        else:
            lines.append(
                "\tif err := %s.services.%s.%s(ctx, %s.To%sParams(req)); err != nil {"
                % (recv, svc, name, conv, name)
            )
            lines.append("\t\treturn nil, apierror.Wrap(err)")
            lines.append("\t}")
            lines.append("")
            lines.append("\treturn %s, nil" % empty_return(name + "Response", op["resp_field"], op["resp_form"]))
        lines.append("}")
        lines.append("")

    for op_key, prefix, doc in (
        ("delete", "Delete", "мягко удаляет строку"),
        ("undelete", "Undelete", "восстанавливает мягко удалённую строку"),
    ):
        if op_key not in ops:
            continue
        op = ops[op_key]
        rpc = op["query"]
        lines.append("// %s %s %s." % (rpc, doc, e["table"]))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, rpc, proto, rpc, proto, rpc)
        )
        if op["arg_form"] == "id":
            lines.append("\tif err := validation.ValidateID(req.GetId()); err != nil {")
            lines.append("\t\treturn nil, apierror.Wrap(err)")
            lines.append("\t}")
            lines.append("")
            lines.append("\tif err := %s.services.%s.%s(ctx, req.GetId()); err != nil {" % (recv, svc, rpc))
        else:
            lines.append("\tif err := %s.Validate%s(req); err != nil {" % (valid, rpc))
            lines.append("\t\treturn nil, apierror.Wrap(err)")
            lines.append("\t}")
            lines.append("")
            lines.append(
                "\tif err := %s.services.%s.%s(ctx, %s.To%sParams(req)); err != nil {"
                % (recv, svc, rpc, conv, rpc)
            )
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn %s, nil" % empty_return(rpc + "Response", op["resp_field"], op["resp_form"]))
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


def main():
    global INTERNAL, CONV_IMPORT, VALID_IMPORT, ERRORS_IMPORT, VALIDATOR_IMPORT, APIERR_IMPORT

    args, config = parse_args("Генерация Go-слоёв по model.json")

    INTERNAL = repo_path(config["internal"])
    shared = config["shared_imports"]
    CONV_IMPORT = shared["converter"]
    VALID_IMPORT = shared["validation"]
    ERRORS_IMPORT = shared["custom_errors"]
    VALIDATOR_IMPORT = shared["validator"]
    APIERR_IMPORT = shared["apierror"]

    model = read_json(args.work_dir, "model.json")

    for sqlc, d in model.items():
        meta = dict(d["meta"])
        meta["sqlc"] = sqlc
        entities = d["entities"]

        for e in entities:
            write("%s/%s.go" % (meta["conv_dir"], e["file"]), gen_converter(meta, e))
            write("%s/%s_test.go" % (meta["conv_dir"], e["file"]), gen_converter_test(meta, e))

            write("%s/%s.go" % (meta["valid_dir"], e["file"]), gen_validation(meta, e))
            write("%s/%s_test.go" % (meta["valid_dir"], e["file"]), gen_validation_test(meta, e))

            write("%s/%s.go" % (meta["service_dir"], e["file"]), gen_service(meta, e))
            write("%s/%s.go" % (meta["api_dir"], e["file"]), gen_api(meta, e))

        write("%s/limits.go" % meta["valid_dir"], gen_limits(meta, entities))

    print("файлов записано:", len(written))
    print("каталог:", INTERNAL)
    print("\nдальше: gofmt -w <internal>, затем go build ./... и go test ./internal/...")


main()
