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
сообщение Pagination, парный Count<имя>) — см.
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
    if pair == ("string", "*string"):
        return "converter.StringToProto(%s)" % src
    if pair == ("int64", "*int64"):
        # Колонка NOT NULL, а поле в .proto помечено optional.
        return "converter.Int64ToProto(%s)" % src
    if pair == ("uuid.UUID", "string"):
        return "converter.UUIDToProto(%s)" % src
    raise Exception("нет правила для %s -> %s" % pair)


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
    raise Exception("нет правила для %s -> %s" % pair)


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


def all_filter_fields(e):
    """Поля фильтра всех постраничных списков сущности (см. gen_converter)."""
    fields = []
    for l in e["lists"]:
        if l["paginated"]:
            fields += l.get("filter_fields") or []
    return fields


def needs_converter_pkg(e):
    for f in e["row_fields"] + all_param_fields(e) + all_filter_fields(e):
        if f["proto_type"] in ("*timestamppb.Timestamp", "*string"):
            return True
        if f["go_type"] == "uuid.UUID":
            return True
    for lk in e["lookups"]:
        if lk["arg_type"] == "uuid.UUID":
            return True
    return False


def needs_uuid_pkg(e, extra_fields=()):
    for f in list(e["row_fields"]) + list(extra_fields) + all_filter_fields(e):
        if f["go_type"] == "uuid.UUID":
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

    lines.append("// ToCreate%sParams собирает параметры вставки %s из запроса gRPC." % (E, e["table"]))
    lines.append("// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.")
    lines.append(
        "func ToCreate%sParams(req *%s.Create%sRequest) %s.%s {" % (E, proto, E, repo, e["create_params"])
    )
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


def sample_value(f, idx, variant=0):
    t = f["go_type"]
    if t == "string":
        # Поля фильтра постраничного списка не привязаны к колонке (нет "col").
        tag = f.get("col") or f["go"]
        return '"%s-%d"' % (tag.replace("_", "-"), variant)
    if t == "int64":
        return str(100 + idx + variant * 1000)
    if t == "int16":
        return str(10 + idx + variant)
    if t == "bool":
        return "true"
    if t == "uuid.UUID":
        return '"%s"' % uuid_literal(idx, variant)
    return None


def repo_sample_value(f, idx, variant=0):
    value = sample_value(f, idx, variant)
    if f["go_type"] == "uuid.UUID":
        return "uuid.MustParse(%s)" % value
    return value


def gen_params_test_block(proto, repo, req_type, params_type, func_name, fields):
    """TestTo<Func>Params + TestTo<Func>ParamsNil — общий шаблон для Create и
    для update/delete/undelete в форме параметров."""
    lines = []
    lines.append("func Test%s(t *testing.T) {" % func_name)
    lines.append("\treq := &%s.%s{" % (proto, req_type))
    for i, f in enumerate(fields):
        v = sample_value(f, i)
        lines.append("\t\t%s: %s," % (f["proto"], v))
    lines.append("\t}")
    lines.append("")
    lines.append("\twant := %s.%s{" % (repo, params_type))
    for i, f in enumerate(fields):
        v = repo_sample_value(f, i)
        if f["go_type"] == "int16":
            v = str(10 + i)
        lines.append("\t\t%s: %s," % (f["go"], v))
    lines.append("\t}")
    lines.append("")
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

    has_time = any(f["go_type"] in ("time.Time", "sql.NullTime") for f in e["row_fields"])
    has_nulltime = any(f["go_type"] == "sql.NullTime" for f in e["row_fields"])

    extra_fields = list(e["create_fields"])
    for _, op in param_ops:
        extra_fields += op["fields"]

    lines = []
    lines.append("package %s" % meta["conv_pkg"])
    lines.append("")
    lines.append("import (")
    if has_nulltime:
        lines.append('\t"database/sql"')
    lines.append('\t"testing"')
    if has_time:
        lines.append('\t"time"')
    lines.append("")
    if needs_uuid_pkg(e, extra_fields):
        lines.append('\t"github.com/google/uuid"')
    lines.append('\t"%s"' % meta["repo_import"])
    lines.append("\t" + meta["proto_import"])
    if paginated_lists:
        lines.append('\t"%s"' % VALID_IMPORT)
    lines.append(")")
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

    lines.append("// test%sRow — строка %s со значениями, различимыми между полями." % (E, e["table"]))
    lines.append("func test%sRow() %s.%s {" % (E, repo, e["row_struct"]))
    lines.append("\treturn %s.%s{" % (repo, e["row_struct"]))
    for i, f in enumerate(e["row_fields"]):
        if f["go_type"] == "time.Time":
            v = "%sCreatedAt" % lower_first(E) if f["col"] == "created_at" else "%sUpdatedAt" % lower_first(E)
        elif f["go_type"] == "sql.NullTime":
            v = "sql.NullTime{Time: %sUpdatedAt, Valid: true}" % lower_first(E)
        elif f["col"] == "is_deleted":
            v = "false"
        else:
            v = repo_sample_value(f, i)
        lines.append("\t\t%s: %s," % (f["go"], v))
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    lines.append("func Test%sToProto(t *testing.T) {" % E)
    lines.append("\trow := test%sRow()" % E)
    lines.append("\tgot := %sToProto(row)" % E)
    lines.append("")
    lines.append("\tif got == nil {")
    lines.append('\t\tt.Fatal("%sToProto() = nil, want value")' % E)
    lines.append("\t}")
    lines.append("")
    for f in e["row_fields"]:
        g = "got.Get%s()" % f["proto"]
        pair = (f["go_type"], f["proto_type"])
        if pair == ("time.Time", "*timestamppb.Timestamp"):
            lines.append("\tif !%s.AsTime().Equal(row.%s) {" % (g, f["go"]))
            lines.append('\t\tt.Errorf("%s = %%v, want %%v", %s.AsTime(), row.%s)' % (f["proto"], g, f["go"]))
        elif pair == ("sql.NullTime", "*timestamppb.Timestamp"):
            lines.append("\tif !%s.AsTime().Equal(row.%s.Time) {" % (g, f["go"]))
            lines.append(
                '\t\tt.Errorf("%s = %%v, want %%v", %s.AsTime(), row.%s.Time)' % (f["proto"], g, f["go"])
            )
        elif pair == ("int16", "int32"):
            lines.append("\tif %s != int32(row.%s) {" % (g, f["go"]))
            lines.append('\t\tt.Errorf("%s = %%d, want %%d", %s, row.%s)' % (f["proto"], g, f["go"]))
        elif pair == ("uuid.UUID", "string"):
            lines.append("\tif %s != row.%s.String() {" % (g, f["go"]))
            lines.append(
                '\t\tt.Errorf("%s = %%q, want %%q", %s, row.%s.String())' % (f["proto"], g, f["go"])
            )
        else:
            verb = "%q" if f["go_type"] == "string" else ("%d" if f["go_type"] == "int64" else "%v")
            lines.append("\tif %s != row.%s {" % (g, f["go"]))
            lines.append(
                '\t\tt.Errorf("%s = %s, want %s", %s, row.%s)' % (f["proto"], verb, verb, g, f["go"])
            )
        lines.append("\t}")
        lines.append("")
    lines.append("}")
    lines.append("")

    first_str = next((f for f in e["row_fields"] if f["go_type"] == "string"), None)
    id_field = next((f for f in e["row_fields"] if f["col"] == "id"), None)
    mark = first_str or id_field
    lines.append("func Test%sToProto(t *testing.T) {" % P)
    lines.append("\tfirst := test%sRow()" % E)
    lines.append("")
    lines.append("\tsecond := test%sRow()" % E)
    if id_field:
        lines.append("\tsecond.%s = 999" % id_field["go"])
    if first_str:
        lines.append('\tsecond.%s = "second-value"' % first_str["go"])
    lines.append("")
    lines.append("\ttests := []struct {")
    lines.append("\t\tname    string")
    lines.append("\t\tinput   []%s.%s" % (repo, e["row_struct"]))
    lines.append("\t\twantLen int")
    lines.append("\t}{")
    lines.append(
        '\t\t{name: "two rows", input: []%s.%s{first, second}, wantLen: 2},' % (repo, e["row_struct"])
    )
    lines.append('\t\t{name: "empty slice", input: []%s.%s{}, wantLen: 0},' % (repo, e["row_struct"]))
    lines.append('\t\t{name: "nil slice", input: nil, wantLen: 0},')
    lines.append("\t}")
    lines.append("")
    lines.append("\tfor _, tt := range tests {")
    lines.append("\t\tt.Run(tt.name, func(t *testing.T) {")
    lines.append("\t\t\tgot := %sToProto(tt.input)" % P)
    lines.append("")
    lines.append("\t\t\tif got == nil {")
    lines.append('\t\t\t\tt.Fatal("%sToProto() = nil, want empty slice")' % P)
    lines.append("\t\t\t}")
    lines.append("")
    lines.append("\t\t\tif len(got) != tt.wantLen {")
    lines.append('\t\t\t\tt.Fatalf("len = %d, want %d", len(got), tt.wantLen)')
    lines.append("\t\t\t}")
    lines.append("\t\t})")
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    # Create params
    lines.extend(
        gen_params_test_block(
            proto, repo, "Create%sRequest" % E, e["create_params"], "ToCreate%sParams" % E, e["create_fields"]
        )
    )

    # update/delete/undelete в форме параметров
    for op_key, op in param_ops:
        name = op["query"]
        lines.extend(
            gen_params_test_block(proto, repo, name + "Request", op["params_type"], "To%sParams" % name, op["fields"])
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
                lines.append("\t\t%s: %s," % (f["proto"], sample_value(f, i)))
            lines.append("\t}")
            lines.append("")
            lines.append("\tgot := To%sParams(req)" % name)
            lines.append("")
            for i, f in enumerate(l["filter_fields"]):
                v = repo_sample_value(f, i)
                verb = "%q" if f["go_type"] == "string" else ("%d" if f["go_type"] in ("int64", "int16") else "%v")
                lines.append("\tif got.%s != %s {" % (f["go"], v))
                lines.append('\t\tt.Errorf("%s = %s, want %s", got.%s, %s)' % (f["go"], verb, verb, f["go"], v))
                lines.append("\t}")
            lines.append("}")
            lines.append("")

    return "\n".join(lines)


# ---------- validation ----------
def field_validation_line(f, limit_entity):
    col, getter = f["col"], "req.Get%s()" % f["proto"]
    if f["go_type"] == "string":
        if f["varchar"] is None:
            raise Exception("нет границы varchar для %s" % col)
        return '\tv.StringVarchar("%s", %s, %s)' % (col, getter, limit_const(limit_entity, f["proto"]))
    if f["go_type"] == "int64":
        return '\tv.Int64ID("%s", %s)' % (col, getter)
    if f["go_type"] == "int16":
        return '\tv.Int32Between("%s", %s, math.MinInt16, math.MaxInt16)' % (col, getter)
    if f["go_type"] == "uuid.UUID":
        return '\tv.StringUUID("%s", %s)' % (col, getter)
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

    lines = []
    lines.append("package %s" % meta["conv_pkg"])
    lines.append("")
    lines.append("import (")
    if needs_math:
        lines.append('\t"math"')
        lines.append("")
    lines.append("\t" + meta["proto_import"])
    if paginated_lists:
        lines.append('\t"%s"' % VALID_IMPORT)
    lines.append('\t"%s"' % VALIDATOR_IMPORT)
    lines.append(")")
    lines.append("")

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

    def literal(f, i):
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
        return "0"

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

    head = ["package %s" % meta["conv_pkg"], "", "import (", '\t"errors"']
    if "math.MaxInt16" in body:
        head.append('\t"math"')
    if "strings.Repeat(" in body:
        head.append('\t"strings"')
    head.append('\t"testing"')
    head.append("")
    head.append("\t" + meta["proto_import"])
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

    lines = []
    lines.append("package %s" % meta["service_pkg"])
    lines.append("")
    lines.append("import (")
    lines.append('\t"context"')
    lines.append('\t"database/sql"')
    lines.append('\t"errors"')
    lines.append('\t"fmt"')
    lines.append("")
    if any(lk["arg_type"] == "uuid.UUID" for lk in e["lookups"]):
        lines.append('\t"github.com/google/uuid"')
    lines.append('\t"%s"' % meta["repo_import"])
    lines.append('\tcustomerrors "%s"' % ERRORS_IMPORT)
    lines.append(")")
    lines.append("")

    def one(name, doc, notfound_msg):
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
        lines.extend(one("Get%sById" % E, "возвращает активную строку %s по id." % low, low))
    if "get_deleted_by_id" in ops:
        lines.extend(
            one("GetDeleted%sById" % E, "возвращает мягко удалённую строку %s по id." % low, "deleted " + low)
        )

    for l in e["lists"]:
        name = l["query"]
        if not l["paginated"]:
            lines.append("// %s возвращает строки %s." % (name, low))
            lines.append("func (%s *%s) %s(ctx context.Context) ([]%s, error) {" % (recv, st, name, row))
            lines.append("\trows, err := %s.%s.%s(ctx)" % (recv, field, name))
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

    # Create — всегда есть
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
            lines.append("//")
            lines.append("// Запрос — :exec и не сообщает число затронутых строк, поэтому")
            if "get_by_id" in ops:
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

    return "\n".join(lines)


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

    def by_id(rpc, doc, svc_call, resp_field):
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
            "\treturn &%s.%sResponse{%s: %s.%sToProto(row)}, nil" % (proto, rpc, resp_field, conv, E)
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
            )
        )
    if "get_deleted_by_id" in ops:
        lines.extend(
            by_id(
                "GetDeleted%sById" % E,
                "отдаёт мягко удалённую строку %s по id." % e["table"],
                "GetDeleted%sById" % E,
                ops["get_deleted_by_id"]["resp_field"],
            )
        )

    for l in e["lists"]:
        name = l["query"]
        if not l["paginated"]:
            lines.append("// %s отдаёт строки %s." % (name, e["table"]))
            lines.append(
                "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
                % (recv, at, name, proto, name, proto, name)
            )
            lines.append("\trows, err := %s.services.%s.%s(ctx)" % (recv, svc, name))
            lines.append("\tif err != nil {")
            lines.append("\t\treturn nil, apierror.Wrap(err)")
            lines.append("\t}")
            lines.append("")
            lines.append(
                "\treturn &%s.%sResponse{%s: %s.%sToProto(rows)}, nil" % (proto, name, l["resp_field"], conv, P)
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
        lines.append("\t\t%s: %s.%sToProto(rows)," % (l["resp_field"], conv, P))
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
            % (proto, rpc, fk["resp_field"], conv, P)
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
            % (proto, rpc, lk["resp_field"], conv, E)
        )
        lines.append("}")
        lines.append("")

    # Create — всегда есть
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
