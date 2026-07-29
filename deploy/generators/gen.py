# -*- coding: utf-8 -*-
"""Шаг 3 конвейера: генерация слоёв converter / validation / service / api.

Читает build/model.json и переписывает по файлу на таблицу в каждом слое.
ВНИМАНИЕ: перезаписывает файлы без спроса. Список того, что генерируется,
а что писано руками, — в documentation/dev_instructions/crud/generate_gprpc_go.md

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
    return s[0].lower() + s[1:]


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
        # Колонка NOT NULL, а поле в .proto помечено optional.
        return "converter.StringToProto(%s)" % src
    if pair == ("uuid.UUID", "string"):
        # В protobuf нет типа для UUID — передаём канонической строкой.
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
        # Геттер optional-поля отдаёт "" вместо nil — колонка NOT NULL это отсекает
        # на валидации, до конвертера пустое значение не доходит.
        return getter
    if pair == ("uuid.UUID", "string"):
        # Формат строки проверен на валидации, разбор здесь уже не может не удаться.
        return "converter.ProtoToUUID(%s)" % getter
    raise Exception("нет правила для %s -> %s" % pair)


def lookup_arg_expr(lk, var):
    """Выражение для аргумента выборки по уникальной колонке из запроса gRPC."""
    getter = "%s.Get%s()" % (var, lk["proto_field"])
    pair = (lk["arg_type"], lk["proto_type"])
    if pair in (("int64", "int64"), ("string", "string")):
        return getter
    if pair == ("uuid.UUID", "string"):
        return "converter.ProtoToUUID(%s)" % getter
    raise Exception("нет правила для аргумента выборки %s -> %s" % pair)


def needs_converter_pkg(entity):
    for group in ("row_fields", "create_fields", "update_fields"):
        for f in entity[group]:
            if f["proto_type"] in ("*timestamppb.Timestamp", "*string"):
                return True
            if f["go_type"] == "uuid.UUID":
                return True
    for lk in entity["lookups"]:
        if lk["arg_type"] == "uuid.UUID":
            return True
    return False


def needs_uuid_pkg(entity, groups=("row_fields", "create_fields", "update_fields")):
    """Нужен ли импорт github.com/google/uuid в файле."""
    for group in groups:
        for f in entity[group]:
            if f["go_type"] == "uuid.UUID":
                return True
    for lk in entity["lookups"]:
        if lk["arg_type"] == "uuid.UUID":
            return True
    return False


# ---------- converter ----------
def gen_converter(meta, e):
    proto = meta["proto_alias"]
    repo = meta["repo_alias"]
    E, P = e["entity"], e["plural"]

    imports = []
    if any(lk["arg_type"] == "uuid.UUID" for lk in e["lookups"]):
        imports.append('\t"github.com/google/uuid"')
        imports.append("")
    if needs_converter_pkg(e):
        imports.append('\t"%s"' % CONV_IMPORT)
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

    lines.append("// ToUpdate%sByIdParams собирает параметры обновления %s из запроса gRPC." % (E, e["table"]))
    lines.append("// updated_at выставляет SQL, is_deleted через обновление не меняется.")
    lines.append(
        "func ToUpdate%sByIdParams(req *%s.Update%sByIdRequest) %s.%s {"
        % (E, proto, E, repo, e["update_params"])
    )
    lines.append("\treturn %s.%s{" % (repo, e["update_params"]))
    for f in e["update_fields"]:
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

    return "\n".join(lines)


# ---------- значения для тестов ----------
def uuid_literal(idx, variant=0):
    """Валидный UUID канонического вида, различимый по индексу поля."""
    return "00000000-0000-4000-8000-%012d" % (idx + variant * 1000 + 1)


def sample_value(f, idx, variant=0):
    """Различимые значения полей: перепутанные местами поля тест обязан заметить.

    Значение для стороны proto: у колонки uuid это строка.
    """
    t = f["go_type"]
    if t == "string":
        return '"%s-%d"' % (f["col"].replace("_", "-"), variant)
    if t == "int64":
        return str(100 + idx + variant * 1000)
    if t == "int16":
        return str(10 + idx + variant)
    if t == "bool":
        return "true"
    if t == "uuid.UUID":
        return '"%s"' % uuid_literal(idx, variant)
    return None


def lookup_sample_value(lk):
    """Значение аргумента выборки на стороне proto."""
    if lk["arg_type"] == "uuid.UUID":
        return '"%s"' % uuid_literal(0, 7)
    if lk["arg_type"] == "int64":
        return "77"
    return '"%s-lookup"' % lk["col"].replace("_", "-")


def lookup_zero_value(lk):
    """Нулевое значение аргумента выборки — что отдаёт конвертер на nil-запросе."""
    if lk["arg_type"] == "uuid.UUID":
        return "uuid.Nil"
    if lk["arg_type"] == "int64":
        return "0"
    return '""'


def repo_sample_value(f, idx, variant=0):
    """То же значение для стороны sqlc: строку uuid надо разобрать в uuid.UUID."""
    value = sample_value(f, idx, variant)
    if f["go_type"] == "uuid.UUID":
        return "uuid.MustParse(%s)" % value
    return value


def gen_converter_test(meta, e):
    proto = meta["proto_alias"]
    repo = meta["repo_alias"]
    E, P = e["entity"], e["plural"]

    has_time = any(f["go_type"] in ("time.Time", "sql.NullTime") for f in e["row_fields"])
    has_nulltime = any(f["go_type"] == "sql.NullTime" for f in e["row_fields"])

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
    if needs_uuid_pkg(e):
        lines.append('\t"github.com/google/uuid"')
    lines.append('\t"%s"' % meta["repo_import"])
    lines.append("\t" + meta["proto_import"])
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

    # Эталонная строка
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

    # Тест сущности -> proto: каждое поле по отдельности
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

    # is_deleted переносится, а не берётся из воздуха
    lines.append("func Test%sToProtoDeleted(t *testing.T) {" % E)
    lines.append("\trow := test%sRow()" % E)
    lines.append("\trow.IsDeleted = true")
    lines.append("")
    lines.append("\tif got := %sToProto(row); !got.GetIsDeleted() {" % E)
    lines.append('\t\tt.Error("IsDeleted = false, want true")')
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    # Список
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
    lines.append("\t\t\t// Пустой вход даёт пустой, а не nil-слайс.")
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

    # Порядок списка сохраняется
    if mark:
        lines.append("func Test%sToProtoKeepsOrder(t *testing.T) {" % P)
        lines.append("\tfirst := test%sRow()" % E)
        lines.append("\tsecond := test%sRow()" % E)
        if mark["go_type"] == "string":
            lines.append('\tsecond.%s = "second-value"' % mark["go"])
            cmp_first = "first.%s" % mark["go"]
            cmp_second = "second.%s" % mark["go"]
            getter = "Get%s()" % mark["proto"]
            verb = "%q"
        else:
            lines.append("\tsecond.%s = 999" % mark["go"])
            cmp_first = "first.%s" % mark["go"]
            cmp_second = "second.%s" % mark["go"]
            getter = "Get%s()" % mark["proto"]
            verb = "%d"
        lines.append("")
        lines.append(
            "\tgot := %sToProto([]%s.%s{first, second})" % (P, repo, e["row_struct"])
        )
        lines.append("")
        lines.append("\tif got[0].%s != %s {" % (getter, cmp_first))
        lines.append('\t\tt.Errorf("[0] = %s, want %s", got[0].%s, %s)' % (verb, verb, getter, cmp_first))
        lines.append("\t}")
        lines.append("")
        lines.append("\tif got[1].%s != %s {" % (getter, cmp_second))
        lines.append('\t\tt.Errorf("[1] = %s, want %s", got[1].%s, %s)' % (verb, verb, getter, cmp_second))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    # Create params
    lines.append("func TestToCreate%sParams(t *testing.T) {" % E)
    lines.append("\treq := &%s.Create%sRequest{" % (proto, E))
    for i, f in enumerate(e["create_fields"]):
        v = sample_value(f, i)
        if f["proto_type"] == "*string":
            lines.append("\t\t%s: %sPtr(%s)," % (f["proto"], lower_first(E), v))
        else:
            lines.append("\t\t%s: %s," % (f["proto"], v))
    lines.append("\t}")
    lines.append("")
    lines.append("\twant := %s.%s{" % (repo, e["create_params"]))
    for i, f in enumerate(e["create_fields"]):
        v = repo_sample_value(f, i)
        if f["go_type"] == "int16":
            v = str(10 + i)
        lines.append("\t\t%s: %s," % (f["go"], v))
    lines.append("\t}")
    lines.append("")
    lines.append("\tif got := ToCreate%sParams(req); got != want {" % E)
    lines.append('\t\tt.Errorf("ToCreate%sParams() = %%+v, want %%+v", got, want)' % E)
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    lines.append("func TestToCreate%sParamsNil(t *testing.T) {" % E)
    lines.append("\t// Геттеры protobuf безопасны на nil: сервер не должен падать.")
    lines.append(
        "\tif got := ToCreate%sParams(nil); got != (%s.%s{}) {" % (E, repo, e["create_params"])
    )
    lines.append('\t\tt.Errorf("ToCreate%sParams(nil) = %%+v, want zero value", got)' % E)
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    # Update params
    lines.append("func TestToUpdate%sByIdParams(t *testing.T) {" % E)
    lines.append("\treq := &%s.Update%sByIdRequest{" % (proto, E))
    for i, f in enumerate(e["update_fields"]):
        v = sample_value(f, i)
        if f["proto_type"] == "*string":
            lines.append("\t\t%s: %sPtr(%s)," % (f["proto"], lower_first(E), v))
        else:
            lines.append("\t\t%s: %s," % (f["proto"], v))
    lines.append("\t}")
    lines.append("")
    lines.append("\twant := %s.%s{" % (repo, e["update_params"]))
    for i, f in enumerate(e["update_fields"]):
        v = repo_sample_value(f, i)
        if f["go_type"] == "int16":
            v = str(10 + i)
        lines.append("\t\t%s: %s," % (f["go"], v))
    lines.append("\t}")
    lines.append("")
    lines.append("\tif got := ToUpdate%sByIdParams(req); got != want {" % E)
    lines.append('\t\tt.Errorf("ToUpdate%sByIdParams() = %%+v, want %%+v", got, want)' % E)
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    lines.append("func TestToUpdate%sByIdParamsNil(t *testing.T) {" % E)
    lines.append(
        "\tif got := ToUpdate%sByIdParams(nil); got != (%s.%s{}) {" % (E, repo, e["update_params"])
    )
    lines.append('\t\tt.Errorf("ToUpdate%sByIdParams(nil) = %%+v, want zero value", got)' % E)
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    # Выборки по уникальной колонке
    for lk in e["lookups"]:
        proto_lit = lookup_sample_value(lk)
        want = "uuid.MustParse(%s)" % proto_lit if lk["arg_type"] == "uuid.UUID" else proto_lit
        lines.append("func TestTo%sArg(t *testing.T) {" % lk["query"])
        lines.append("\treq := &%s.%sRequest{%s: %s}" % (proto, lk["query"], lk["proto_field"], proto_lit))
        lines.append("")
        lines.append("\tif got := To%sArg(req); got != %s {" % (lk["query"], want))
        lines.append('\t\tt.Errorf("To%sArg() = %%v, want %%v", got, %s)' % (lk["query"], want))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func TestTo%sArgNil(t *testing.T) {" % lk["query"])
        lines.append("\t// Геттеры protobuf безопасны на nil: сервер не должен падать.")
        lines.append("\tif got := To%sArg(nil); got != %s {" % (lk["query"], lookup_zero_value(lk)))
        lines.append(
            '\t\tt.Errorf("To%sArg(nil) = %%v, want zero value", got)' % lk["query"]
        )
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    # Хелпер для optional-полей
    if any(f["proto_type"] == "*string" for f in e["create_fields"] + e["update_fields"]):
        lines.append("// %sPtr — адрес строкового литерала для optional-полей proto." % lower_first(E))
        lines.append("func %sPtr(s string) *string {" % lower_first(E))
        lines.append("\treturn &s")
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


# ---------- validation ----------
def validation_calls(entity, fields, indent="\t"):
    out = []
    for f in fields:
        col, getter = f["col"], "req.Get%s()" % f["proto"]
        if f["go_type"] == "string":
            if f["varchar"] is None:
                raise Exception("нет границы varchar для %s" % col)
            out.append(
                '%sv.StringVarchar("%s", %s, %s)' % (indent, col, getter, limit_const(entity, f["proto"]))
            )
        elif f["go_type"] == "int64":
            out.append('%sv.Int64ID("%s", %s)' % (indent, col, getter))
        elif f["go_type"] == "int16":
            out.append(
                '%sv.Int32Between("%s", %s, math.MinInt16, math.MaxInt16)' % (indent, col, getter)
            )
        elif f["go_type"] == "uuid.UUID":
            out.append('%sv.StringUUID("%s", %s)' % (indent, col, getter))
        elif f["go_type"] == "bool":
            pass  # любое булево значение допустимо
        else:
            raise Exception("нет правила валидации для %s (%s)" % (col, f["go_type"]))
    return out


def lookup_validation_call(entity, lk):
    """Проверка аргумента выборки по уникальной колонке."""
    getter = "req.Get%s()" % lk["proto_field"]
    if lk["arg_type"] == "uuid.UUID":
        return 'v.StringUUID("%s", %s)' % (lk["col"], getter)
    if lk["arg_type"] == "int64":
        return 'v.Int64ID("%s", %s)' % (lk["col"], getter)
    if lk["arg_type"] == "string":
        if lk["varchar"] is None:
            raise Exception("нет границы varchar для %s" % lk["col"])
        return 'v.StringVarchar("%s", %s, %s)' % (
            lk["col"],
            getter,
            limit_const(entity, lk["proto_field"]),
        )
    raise Exception("нет правила валидации для выборки по %s (%s)" % (lk["col"], lk["arg_type"]))


def gen_validation(meta, e):
    proto = meta["proto_alias"]
    E = e["entity"]

    # Изменяемые поля общие для вставки и обновления — всё, кроме id.
    common = [f for f in e["update_fields"] if f["col"] != "id"]
    create_only = [f for f in e["create_fields"] if f["col"] not in {c["col"] for c in common}]

    needs_math = any(f["go_type"] == "int16" for f in e["create_fields"] + e["update_fields"])

    lines = []
    lines.append("package %s" % meta["conv_pkg"])
    lines.append("")
    lines.append("import (")
    if needs_math:
        lines.append('\t"math"')
        lines.append("")
    lines.append("\t" + meta["proto_import"])
    lines.append('\t"%s"' % VALIDATOR_IMPORT)
    lines.append(")")
    lines.append("")

    helper = "%sWritableFields" % lower_first(E)
    args = []
    str_args = [f for f in common if f["go_type"] == "string"]
    other_args = [f for f in common if f["go_type"] != "string"]

    lines.append("// %s проверяет поля, общие для вставки и обновления %s." % (helper, e["table"]))
    lines.append("func %s(" % helper)
    lines.append("\tv *validator.Validator,")
    for f in common:
        gotype = {
            "int16": "int32",
            "int64": "int64",
            "string": "string",
            "bool": "bool",
            "uuid.UUID": "string",
        }[f["go_type"]]
        lines.append("\t%s %s," % (lower_first(f["proto"]), gotype))
    lines.append(") {")
    for f in common:
        col, val = f["col"], lower_first(f["proto"])
        if f["go_type"] == "string":
            lines.append('\tv.StringVarchar("%s", %s, %s)' % (col, val, limit_const(E, f["proto"])))
        elif f["go_type"] == "int64":
            lines.append('\tv.Int64ID("%s", %s)' % (col, val))
        elif f["go_type"] == "int16":
            lines.append('\tv.Int32Between("%s", %s, math.MinInt16, math.MaxInt16)' % (col, val))
        elif f["go_type"] == "uuid.UUID":
            lines.append('\tv.StringUUID("%s", %s)' % (col, val))
    lines.append("}")
    lines.append("")

    def call_helper(reqvar):
        out = ["\t%s(" % helper, "\t\tv,"]
        for f in common:
            out.append("\t\t%s.Get%s()," % (reqvar, f["proto"]))
        out.append("\t)")
        return out

    # Create
    lines.append("// ValidateCreate%s проверяет запрос на вставку строки %s." % (E, e["table"]))
    lines.append("func ValidateCreate%s(req *%s.Create%sRequest) error {" % (E, proto, E))
    lines.append("\tv := validator.New()")
    lines.append("")
    lines.append("\tif req == nil {")
    lines.append('\t\tv.AddError("request", validator.MsgRequired)')
    lines.append("\t\treturn v.Err()")
    lines.append("\t}")
    lines.append("")
    lines.extend(call_helper("req"))
    if create_only:
        lines.append("")
        lines.extend(validation_calls(E, create_only))
    lines.append("")
    lines.append("\treturn v.Err()")
    lines.append("}")
    lines.append("")

    # Update
    lines.append("// ValidateUpdate%sById проверяет запрос на обновление строки %s." % (E, e["table"]))
    lines.append("// К изменяемым полям добавляется id обновляемой записи.")
    lines.append("func ValidateUpdate%sById(req *%s.Update%sByIdRequest) error {" % (E, proto, E))
    lines.append("\tv := validator.New()")
    lines.append("")
    lines.append("\tif req == nil {")
    lines.append('\t\tv.AddError("request", validator.MsgRequired)')
    lines.append("\t\treturn v.Err()")
    lines.append("\t}")
    lines.append("")
    lines.append('\tv.Int64ID("id", req.GetId())')
    lines.append("")
    lines.extend(call_helper("req"))
    lines.append("")
    lines.append("\treturn v.Err()")
    lines.append("}")
    lines.append("")

    # Выборки по уникальной колонке: тип аргумента зависит от колонки,
    # поэтому общий validation.ValidateID здесь не годится.
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
        lines.append("\t" + lookup_validation_call(E, lk))
        lines.append("")
        lines.append("\treturn v.Err()")
        lines.append("}")
        lines.append("")

    return "\n".join(lines), common, create_only


def gen_validation_test(meta, e, common, create_only):
    proto = meta["proto_alias"]
    E = e["entity"]

    def literal(f, i, bad=None):
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

    def build_valid(kind, fields, extra_id):
        out = ["\treturn &%s.%s{" % (proto, kind)]
        if extra_id:
            out.append("\t\tId: 42,")
        for i, f in enumerate(fields):
            v = literal(f, i)
            if f["proto_type"] == "*string":
                out.append("\t\t%s: %sStrPtr(%s)," % (f["proto"], lower_first(E), v))
            else:
                out.append("\t\t%s: %s," % (f["proto"], v))
        out.append("\t}")
        return out

    create_fields = e["create_fields"]
    update_fields = [f for f in e["update_fields"] if f["col"] != "id"]

    # Импорты дописываются в конце: набор зависит от того, какие типы
    # колонок реально встретились (у таблиц-связок нет строковых полей).
    lines = []
    lines.append("// validCreate%sRequest — заведомо корректный запрос." % E)
    lines.append("// Тесты портят по одному полю, чтобы проверять правила по отдельности.")
    lines.append("func validCreate%sRequest() *%s.Create%sRequest {" % (E, proto, E))
    lines.extend(build_valid("Create%sRequest" % E, create_fields, False))
    lines.append("}")
    lines.append("")

    lines.append("func validUpdate%sByIdRequest() *%s.Update%sByIdRequest {" % (E, proto, E))
    lines.extend(build_valid("Update%sByIdRequest" % E, update_fields, True))
    lines.append("}")
    lines.append("")

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

    def gen_cases(kind, fields, reqtype, valid_fn, validate_fn, with_id):
        out = []
        out.append("func Test%s(t *testing.T) {" % validate_fn)
        out.append("\ttests := []struct {")
        out.append("\t\tname      string")
        out.append("\t\tmutate    func(*%s.%s)" % (proto, reqtype))
        out.append("\t\twantField string")
        out.append("\t}{")
        out.append('\t\t{name: "valid", mutate: func(*%s.%s) {}},' % (proto, reqtype))
        if with_id:
            out.append(
                '\t\t{name: "zero id", mutate: func(r *%s.%s) { r.Id = 0 }, wantField: "id"},'
                % (proto, reqtype)
            )
            out.append(
                '\t\t{name: "negative id", mutate: func(r *%s.%s) { r.Id = -5 }, wantField: "id"},'
                % (proto, reqtype)
            )
        for f in fields:
            col = f["col"]
            if f["go_type"] == "string":
                const = limit_const(E, f["proto"])
                if f["proto_type"] == "*string":
                    empty = 'r.%s = %sStrPtr("")' % (f["proto"], lower_first(E))
                    blank = 'r.%s = %sStrPtr("   ")' % (f["proto"], lower_first(E))
                    long = 'r.%s = %sStrPtr(strings.Repeat("a", %s+1))' % (
                        f["proto"],
                        lower_first(E),
                        const,
                    )
                else:
                    empty = 'r.%s = ""' % f["proto"]
                    blank = 'r.%s = "   "' % f["proto"]
                    long = 'r.%s = strings.Repeat("a", %s+1)' % (f["proto"], const)
                out.append(
                    '\t\t{name: "empty %s", mutate: func(r *%s.%s) { %s }, wantField: "%s"},'
                    % (col, proto, reqtype, empty, col)
                )
                out.append(
                    '\t\t{name: "blank %s", mutate: func(r *%s.%s) { %s }, wantField: "%s"},'
                    % (col, proto, reqtype, blank, col)
                )
                out.append(
                    '\t\t{name: "%s too long", mutate: func(r *%s.%s) { %s }, wantField: "%s"},'
                    % (col, proto, reqtype, long, col)
                )
            elif f["go_type"] == "int64":
                out.append(
                    '\t\t{name: "zero %s", mutate: func(r *%s.%s) { r.%s = 0 }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
                out.append(
                    '\t\t{name: "negative %s", mutate: func(r *%s.%s) { r.%s = -1 }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
            elif f["go_type"] == "int16":
                out.append(
                    '\t\t{name: "%s over smallint", mutate: func(r *%s.%s) { r.%s = math.MaxInt16 + 1 }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
                out.append(
                    '\t\t{name: "%s under smallint", mutate: func(r *%s.%s) { r.%s = math.MinInt16 - 1 }, wantField: "%s"},'
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
                # Сокращённая запись без дефисов — не канонический вид.
                out.append(
                    '\t\t{name: "%s without dashes", mutate: func(r *%s.%s) { r.%s = "00000000000040008000000000000001" }, wantField: "%s"},'
                    % (col, proto, reqtype, f["proto"], col)
                )
        out.append("\t}")
        out.append("")
        out.append("\tfor _, tt := range tests {")
        out.append("\t\tt.Run(tt.name, func(t *testing.T) {")
        out.append("\t\t\treq := %s()" % valid_fn)
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
        out.append("")
        out.append("\t\t\tif len(fields[tt.wantField]) == 0 {")
        out.append('\t\t\t\tt.Errorf("no error on %q, got %v", tt.wantField, fields)')
        out.append("\t\t\t}")
        out.append("")
        out.append("\t\t\t// Порча одного поля не должна задевать остальные.")
        out.append("\t\t\tif len(fields) != 1 {")
        out.append(
            '\t\t\t\tt.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)'
        )
        out.append("\t\t\t}")
        out.append("\t\t})")
        out.append("\t}")
        out.append("}")
        out.append("")
        return out

    lines.extend(
        gen_cases(
            "create",
            create_fields,
            "Create%sRequest" % E,
            "validCreate%sRequest" % E,
            "ValidateCreate%s" % E,
            False,
        )
    )
    lines.extend(
        gen_cases(
            "update",
            update_fields,
            "Update%sByIdRequest" % E,
            "validUpdate%sByIdRequest" % E,
            "ValidateUpdate%sById" % E,
            True,
        )
    )

    # Границы varchar
    str_fields = [f for f in create_fields if f["go_type"] == "string"]
    if str_fields:
        f = str_fields[0]
        const = limit_const(E, f["proto"])
        setter = (
            "req.%s = %sStrPtr(strings.Repeat(%%s, %s))" % (f["proto"], lower_first(E), const)
            if f["proto_type"] == "*string"
            else "req.%s = strings.Repeat(%%s, %s)" % (f["proto"], const)
        )
        lines.append("// Ровно граничная длина проходит: varchar(n) допускает n символов.")
        lines.append("func TestValidateCreate%sAtVarcharLimit(t *testing.T) {" % E)
        lines.append("\treq := validCreate%sRequest()" % E)
        lines.append("\t" + (setter % '"a"'))
        lines.append("")
        lines.append("\tif err := ValidateCreate%s(req); err != nil {" % E)
        lines.append(
            '\t\tt.Errorf("ValidateCreate%s() = %%v, want nil at exactly %%d chars", err, %s)' % (E, const)
        )
        lines.append("\t}")
        lines.append("}")
        lines.append("")
        lines.append("// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.")
        lines.append("func TestValidateCreate%sCyrillicAtVarcharLimit(t *testing.T) {" % E)
        lines.append("\treq := validCreate%sRequest()" % E)
        lines.append("\t" + (setter % '"я"'))
        lines.append("")
        lines.append("\tif err := ValidateCreate%s(req); err != nil {" % E)
        lines.append(
            '\t\tt.Errorf("ValidateCreate%s() = %%v, want nil at exactly %%d cyrillic chars", err, %s)'
            % (E, const)
        )
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    # Пустой запрос: копятся все ошибки
    checked = [f for f in create_fields if f["go_type"] in ("string", "int64", "uuid.UUID")]
    if checked:
        lines.append("func TestValidateCreate%sCollectsAllErrors(t *testing.T) {" % E)
        lines.append("\t// Валидатор копит ошибки, а не падает на первой: клиент видит")
        lines.append("\t// все проблемы запроса за один ответ.")
        lines.append("\terr := ValidateCreate%s(&%s.Create%sRequest{})" % (E, proto, E))
        lines.append("")
        lines.append("\tif err == nil {")
        lines.append('\t\tt.Fatal("ValidateCreate%s() = nil, want errors")' % E)
        lines.append("\t}")
        lines.append("")
        lines.append("\tfields := %sFieldErrors(t, err)" % lower_first(E))
        lines.append("")
        lines.append("\twantFields := []string{%s}" % ", ".join('"%s"' % f["col"] for f in checked))
        lines.append("")
        lines.append("\tfor _, field := range wantFields {")
        lines.append("\t\tif len(fields[field]) == 0 {")
        lines.append('\t\t\tt.Errorf("no error on %q", field)')
        lines.append("\t\t}")
        lines.append("\t}")
        lines.append("")
        lines.append("\tif len(fields) != len(wantFields) {")
        lines.append(
            '\t\tt.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)'
        )
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    lines.append("func TestValidateCreate%sNil(t *testing.T) {" % E)
    lines.append("\tif err := ValidateCreate%s(nil); err == nil {" % E)
    lines.append('\t\tt.Error("ValidateCreate%s(nil) = nil, want error")' % E)
    lines.append("\t}")
    lines.append("}")
    lines.append("")
    lines.append("func TestValidateUpdate%sByIdNil(t *testing.T) {" % E)
    lines.append("\tif err := ValidateUpdate%sById(nil); err == nil {" % E)
    lines.append('\t\tt.Error("ValidateUpdate%sById(nil) = nil, want error")' % E)
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    # Выборки по уникальной колонке
    for lk in e["lookups"]:
        good = lookup_sample_value(lk)
        if lk["arg_type"] == "uuid.UUID":
            bad = [
                ("empty", '""'),
                ("malformed", '"not-a-uuid"'),
                ("without dashes", '"00000000000040008000000000000001"'),
            ]
        elif lk["arg_type"] == "int64":
            bad = [("zero", "0"), ("negative", "-1")]
        else:
            bad = [("empty", '""'), ("blank", '"   "')]

        lines.append("func TestValidate%s(t *testing.T) {" % lk["query"])
        lines.append("\ttests := []struct {")
        lines.append("\t\tname    string")
        lines.append("\t\tvalue   %s" % lk["proto_type"])
        lines.append("\t\twantErr bool")
        lines.append("\t}{")
        lines.append('\t\t{name: "valid", value: %s},' % good)
        for name, value in bad:
            lines.append('\t\t{name: "%s", value: %s, wantErr: true},' % (name, value))
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
        lines.append("")
        lines.append("\t\t\tif !tt.wantErr {")
        lines.append("\t\t\t\treturn")
        lines.append("\t\t\t}")
        lines.append("")
        lines.append("\t\t\tfields := %sFieldErrors(t, err)" % lower_first(E))
        lines.append("")
        lines.append('\t\t\tif len(fields["%s"]) == 0 {' % lk["col"])
        lines.append('\t\t\t\tt.Errorf("no error on %s, got %%v", fields)' % lk["col"])
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

    body = "\n".join(lines)

    if "%sStrPtr(" % lower_first(E) in body:
        body += (
            "\n// %sStrPtr — адрес строкового литерала для optional-полей proto.\n"
            "func %sStrPtr(s string) *string {\n\treturn &s\n}\n"
            % (lower_first(E), lower_first(E))
        )

    head = ["package %s" % meta["conv_pkg"], "", "import (", '\t"errors"']
    if "math.MaxInt16" in body or "math.MinInt16" in body:
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
        for f in e["create_fields"] + e["update_fields"]:
            if f["go_type"] == "string" and f["col"] not in seen:
                seen[f["col"]] = f
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

    def one(name, doc, call, notfound_msg):
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

    lines.extend(one("Get%sById" % E, "возвращает активную строку %s по id." % low, None, low))
    lines.append("// Get%s возвращает все активные строки %s." % (P, low))
    lines.append("func (%s *%s) Get%s(ctx context.Context) ([]%s, error) {" % (recv, st, P, row))
    lines.append("\trows, err := %s.%s.Get%s(ctx)" % (recv, field, P))
    lines.append("")
    lines.append("\tif err != nil {")
    lines.append('\t\treturn nil, fmt.Errorf("get %s: %%w", err)' % low)
    lines.append("\t}")
    lines.append("")
    lines.append("\treturn rows, nil")
    lines.append("}")
    lines.append("")

    lines.extend(
        one("GetDeleted%sById" % E, "возвращает мягко удалённую строку %s по id." % low, None, "deleted " + low)
    )
    lines.append("// GetDeleted%s возвращает все мягко удалённые строки %s." % (P, low))
    lines.append("func (%s *%s) GetDeleted%s(ctx context.Context) ([]%s, error) {" % (recv, st, P, row))
    lines.append("\trows, err := %s.%s.GetDeleted%s(ctx)" % (recv, field, P))
    lines.append("")
    lines.append("\tif err != nil {")
    lines.append('\t\treturn nil, fmt.Errorf("get deleted %s: %%w", err)' % low)
    lines.append("\t}")
    lines.append("")
    lines.append("\treturn rows, nil")
    lines.append("}")
    lines.append("")

    for fk in e["fks"]:
        lines.append("// %s возвращает активные строки %s, отобранные по %s." % (fk["query"], low, fk["col"]))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, %s int64) ([]%s, error) {"
            % (recv, st, fk["query"], fk["arg_name"], row)
        )
        lines.append("\trows, err := %s.%s.%s(ctx, %s)" % (recv, field, fk["query"], fk["arg_name"]))
        lines.append("")
        lines.append("\tif err != nil {")
        lines.append(
            '\t\treturn nil, fmt.Errorf("get %s by %s = %%d: %%w", %s, err)' % (low, fk["col"], fk["arg_name"])
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
            '\t\treturn %s{}, fmt.Errorf("get %s %s = %s: %%w", %s, err)'
            % (row, low, lk["col"], verb, lk["arg_name"])
        )
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn row, nil")
        lines.append("}")
        lines.append("")

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

    lines.append("// Update%sById обновляет активную строку %s и возвращает её целиком." % (E, low))
    lines.append("//")
    lines.append("// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую")
    lines.append("// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,")
    lines.append("// чтобы api-слой ответил NotFound, а не Internal.")
    lines.append(
        "func (%s *%s) Update%sById(ctx context.Context, params %s.%s) (%s, error) {"
        % (recv, st, E, repo, e["update_params"], row)
    )
    lines.append("\trow, err := %s.%s.Update%sById(ctx, params)" % (recv, field, E))
    lines.append("")
    lines.append("\tif err != nil {")
    lines.append("\t\tif errors.Is(err, sql.ErrNoRows) {")
    lines.append(
        '\t\t\treturn %s{}, fmt.Errorf("%s id = %%d: %%w", params.ID, customerrors.ErrNotFound)' % (row, low)
    )
    lines.append("\t\t}")
    lines.append("")
    lines.append(
        '\t\treturn %s{}, fmt.Errorf("%%w: %s id = %%d: %%w", customerrors.ErrUpdate, params.ID, err)'
        % (row, low)
    )
    lines.append("\t}")
    lines.append("")
    lines.append("\treturn row, nil")
    lines.append("}")
    lines.append("")

    lines.append("// Delete%sById мягко удаляет строку %s." % (E, low))
    lines.append("//")
    lines.append("// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута")
    lines.append("// строка, поэтому существование активной записи проверяем заранее —")
    lines.append("// иначе удаление несуществующего id молча возвращало бы успех.")
    lines.append("func (%s *%s) Delete%sById(ctx context.Context, id int64) error {" % (recv, st, E))
    lines.append("\tif _, err := %s.Get%sById(ctx, id); err != nil {" % (recv, E))
    lines.append("\t\treturn errors.Join(customerrors.ErrDelete, err)")
    lines.append("\t}")
    lines.append("")
    lines.append("\tif err := %s.%s.Delete%sById(ctx, id); err != nil {" % (recv, field, E))
    lines.append(
        '\t\treturn fmt.Errorf("%%w: %s id = %%d: %%w", customerrors.ErrDelete, id, err)' % low
    )
    lines.append("\t}")
    lines.append("")
    lines.append("\treturn nil")
    lines.append("}")
    lines.append("")

    lines.append("// Undelete%sById восстанавливает мягко удалённую строку %s." % (E, low))
    lines.append("// Существование удалённой записи проверяется заранее по той же причине,")
    lines.append("// что и в Delete%sById." % E)
    lines.append("func (%s *%s) Undelete%sById(ctx context.Context, id int64) error {" % (recv, st, E))
    lines.append("\tif _, err := %s.GetDeleted%sById(ctx, id); err != nil {" % (recv, E))
    lines.append("\t\treturn errors.Join(customerrors.ErrUndelete, err)")
    lines.append("\t}")
    lines.append("")
    lines.append("\tif err := %s.%s.Undelete%sById(ctx, id); err != nil {" % (recv, field, E))
    lines.append(
        '\t\treturn fmt.Errorf("%%w: %s id = %%d: %%w", customerrors.ErrUndelete, id, err)' % low
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
    rf = e["resp_fields"]

    lines = []
    lines.append("package %s" % meta["api_pkg"])
    lines.append("")
    lines.append("import (")
    lines.append('\t"context"')
    lines.append("")
    lines.append('\t"%s"' % APIERR_IMPORT)
    lines.append('\t%s "%s/%s"' % (conv, CONV_IMPORT, meta["conv_pkg"]))
    lines.append('\t"%s"' % VALID_IMPORT)
    lines.append('\t%s "%s/%s"' % (valid, VALID_IMPORT, meta["conv_pkg"]))
    lines.append("\t" + meta["proto_import"])
    lines.append('\t"google.golang.org/protobuf/types/known/emptypb"')
    lines.append(")")
    lines.append("")

    def by_id(rpc, doc, svc_call, resp_field, deleted=False):
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

    def list_all(rpc, doc, svc_call, resp_field):
        out = []
        out.append("// %s %s" % (rpc, doc))
        out.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, rpc, proto, rpc, proto, rpc)
        )
        out.append("\trows, err := %s.services.%s.%s(ctx)" % (recv, svc, svc_call))
        out.append("\tif err != nil {")
        out.append("\t\treturn nil, apierror.Wrap(err)")
        out.append("\t}")
        out.append("")
        out.append(
            "\treturn &%s.%sResponse{%s: %s.%sToProto(rows)}, nil" % (proto, rpc, resp_field, conv, P)
        )
        out.append("}")
        out.append("")
        return out

    lines.extend(
        by_id("Get%sById" % E, "отдаёт активную строку %s по id." % e["table"], "Get%sById" % E, rf["Get%sById" % E])
    )
    lines.extend(list_all("Get%s" % P, "отдаёт все активные строки %s." % e["table"], "Get%s" % P, rf["Get%s" % P]))
    lines.extend(
        by_id(
            "GetDeleted%sById" % E,
            "отдаёт мягко удалённую строку %s по id." % e["table"],
            "GetDeleted%sById" % E,
            rf["GetDeleted%sById" % E],
        )
    )
    lines.extend(
        list_all(
            "GetDeleted%s" % P,
            "отдаёт все мягко удалённые строки %s." % e["table"],
            "GetDeleted%s" % P,
            rf["GetDeleted%s" % P],
        )
    )

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

    # Create
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
        % (proto, E, rf["Create%s" % E], conv, E)
    )
    lines.append("}")
    lines.append("")

    # Update
    lines.append("// Update%sById обновляет активную строку %s и отдаёт её целиком." % (E, e["table"]))
    lines.append(
        "func (%s *%s) Update%sById(ctx context.Context, req *%s.Update%sByIdRequest) (*%s.Update%sByIdResponse, error) {"
        % (recv, at, E, proto, E, proto, E)
    )
    lines.append("\tif err := %s.ValidateUpdate%sById(req); err != nil {" % (valid, E))
    lines.append("\t\treturn nil, apierror.Wrap(err)")
    lines.append("\t}")
    lines.append("")
    lines.append(
        "\trow, err := %s.services.%s.Update%sById(ctx, %s.ToUpdate%sByIdParams(req))"
        % (recv, svc, E, conv, E)
    )
    lines.append("\tif err != nil {")
    lines.append("\t\treturn nil, apierror.Wrap(err)")
    lines.append("\t}")
    lines.append("")
    lines.append(
        "\treturn &%s.Update%sByIdResponse{%s: %s.%sToProto(row)}, nil"
        % (proto, E, rf["Update%sById" % E], conv, E)
    )
    lines.append("}")
    lines.append("")

    # Delete / Undelete
    for op, doc in (("Delete", "мягко удаляет строку"), ("Undelete", "восстанавливает мягко удалённую строку")):
        rpc = "%s%sById" % (op, E)
        lines.append("// %s %s %s." % (rpc, doc, e["table"]))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, rpc, proto, rpc, proto, rpc)
        )
        lines.append("\tif err := validation.ValidateID(req.GetId()); err != nil {")
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append("\tif err := %s.services.%s.%s(ctx, req.GetId()); err != nil {" % (recv, svc, rpc))
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\treturn &%s.%sResponse{%s: &emptypb.Empty{}}, nil" % (proto, rpc, rf[rpc])
        )
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

            vtext, common, create_only = gen_validation(meta, e)
            write("%s/%s.go" % (meta["valid_dir"], e["file"]), vtext)
            write(
                "%s/%s_test.go" % (meta["valid_dir"], e["file"]),
                gen_validation_test(meta, e, common, create_only),
            )

            write("%s/%s.go" % (meta["service_dir"], e["file"]), gen_service(meta, e))
            write("%s/%s.go" % (meta["api_dir"], e["file"]), gen_api(meta, e))

        write("%s/limits.go" % meta["valid_dir"], gen_limits(meta, entities))

    print("файлов записано:", len(written))
    print("каталог:", INTERNAL)
    print("\nдальше: gofmt -w <internal>, затем go build ./... и go test ./internal/...")


main()
