# -*- coding: utf-8 -*-
"""Шаг 3 конвейера для доменов без табличного CRUD: converter / validation /
service / api по build/model_custom.json.

Пара к gen.py, но пишет по одному файлу на домен (не на таблицу — здесь нет
сущности, только независимые запросы). ВНИМАНИЕ: перезаписывает файлы без
спроса, как и gen.py.

Запуск:  python deploy/generators/gen_custom.py [--config ...] [--work-dir ...]
"""
import io
import os

from genconfig import parse_args, read_json, repo_path

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


def limit_const(col):
    return "%sMaxLen" % lower_first("".join(p.capitalize() for p in col.split("_")))


# ---------- преобразование значений (те же правила, что в gen.py) ----------
def req_to_param_expr(f, var):
    getter = "%s.Get%s()" % (var, f["proto"])
    pair = (f["go_type"], f["proto_type"])
    if pair in (("bool", "bool"), ("int64", "int64"), ("string", "string")):
        return getter
    if pair == ("int16", "int32"):
        return "int16(%s)" % getter
    if pair == ("uuid.UUID", "string"):
        return "converter.ProtoToUUID(%s)" % getter
    raise Exception("нет правила для аргумента %s -> %s" % pair)


def elem_to_proto_expr(elem_conv, var):
    if elem_conv is None:
        return var
    return elem_conv % var


def needs_converter_pkg(call):
    if call["elem_conv"]:
        return True
    return any(f["go_type"] == "uuid.UUID" for f in call["req_fields"])


# ---------- converter ----------
def gen_converter(meta, calls):
    proto = meta["proto_alias"]
    repo = meta["repo_alias"]

    lines = ["package %s" % meta["conv_pkg"], "", "import ("]
    if needs_converter_pkg_any(calls):
        lines.append('\t"%s"' % CONV_IMPORT)
    lines.append('\t"%s"' % meta["repo_import"])
    lines.append("\t" + meta["proto_import"])
    lines.append(")")
    lines.append("")

    for c in calls:
        name = c["query"]
        lines.append("// To%sParams собирает параметры %s.%s из запроса gRPC." % (name, meta["sqlc"], name))
        lines.append(
            "func To%sParams(req *%s.%sRequest) %s.%s {" % (name, proto, name, repo, c["params_type"])
        )
        lines.append("\treturn %s.%s{" % (repo, c["params_type"]))
        for f in c["req_fields"]:
            lines.append("\t\t%s: %s," % (f["go"], req_to_param_expr(f, "req")))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        resp_go_type = ("[]" + c["go_elem"]) if c["is_slice"] else c["go_elem"]
        arg = lower_first(c["resp_field"])
        lines.append(
            "// %sToProto переводит результат %s.%s в ответ gRPC." % (name, meta["sqlc"], name)
        )
        lines.append(
            "func %sToProto(%s %s) *%s.%sResponse {" % (name, arg, resp_go_type, proto, name)
        )
        if c["elem_conv"] is None:
            lines.append("\treturn &%s.%sResponse{%s: %s}" % (proto, name, c["resp_field"], arg))
        else:
            proto_elem = c["proto_elem"]
            lines.append("\titems := make([]%s, 0, len(%s))" % (proto_elem, arg))
            lines.append("")
            lines.append("\tfor _, row := range %s {" % arg)
            lines.append("\t\titems = append(items, %s)" % elem_to_proto_expr(c["elem_conv"], "row"))
            lines.append("\t}")
            lines.append("")
            lines.append("\treturn &%s.%sResponse{%s: items}" % (proto, name, c["resp_field"]))
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


def needs_converter_pkg_any(calls):
    return any(needs_converter_pkg(c) for c in calls)


# ---------- тесты converter ----------
def uuid_literal(idx):
    return "00000000-0000-4000-8000-%012d" % (idx + 1)


def sample_proto(f, idx):
    if f["go_type"] == "string":
        return '"%s-%d"' % (f["col"].replace("_", "-"), idx)
    if f["go_type"] == "int64":
        return str(100 + idx)
    if f["go_type"] == "bool":
        return "true"
    if f["go_type"] == "uuid.UUID":
        return '"%s"' % uuid_literal(idx)
    raise Exception("нет тестового значения для %s" % f["go_type"])


def sample_repo(f, idx):
    v = sample_proto(f, idx)
    if f["go_type"] == "uuid.UUID":
        return "uuid.MustParse(%s)" % v
    return v


def elem_samples(go_elem):
    if go_elem == "int64":
        return "[]int64{1, 2, 3}", "3"
    if go_elem == "string":
        return '[]string{"a", "b"}', "2"
    if go_elem == "bool":
        return "[]bool{true, false}", "2"
    raise Exception("нет тестовых значений для элемента %s" % go_elem)


def gen_converter_test(meta, calls):
    proto = meta["proto_alias"]
    repo = meta["repo_alias"]

    needs_uuid = any(f["go_type"] == "uuid.UUID" for c in calls for f in c["req_fields"])

    lines = ["package %s" % meta["conv_pkg"], "", "import (", '\t"testing"', ""]
    if needs_uuid:
        lines.append('\t"github.com/google/uuid"')
        lines.append("")
    lines.append('\t"%s"' % meta["repo_import"])
    lines.append("\t" + meta["proto_import"])
    lines.append(")")
    lines.append("")

    for c in calls:
        name = c["query"]

        lines.append("func TestTo%sParams(t *testing.T) {" % name)
        lines.append("\treq := &%s.%sRequest{" % (proto, name))
        for i, f in enumerate(c["req_fields"]):
            lines.append("\t\t%s: %s," % (f["proto"], sample_proto(f, i)))
        lines.append("\t}")
        lines.append("")
        lines.append("\twant := %s.%s{" % (repo, c["params_type"]))
        for i, f in enumerate(c["req_fields"]):
            lines.append("\t\t%s: %s," % (f["go"], sample_repo(f, i)))
        lines.append("\t}")
        lines.append("")
        lines.append("\tif got := To%sParams(req); got != want {" % name)
        lines.append('\t\tt.Errorf("To%sParams() = %%+v, want %%+v", got, want)' % name)
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func TestTo%sParamsNil(t *testing.T) {" % name)
        lines.append(
            "\tif got := To%sParams(nil); got != (%s.%s{}) {" % (name, repo, c["params_type"])
        )
        lines.append('\t\tt.Errorf("To%sParams(nil) = %%+v, want zero value", got)' % name)
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        elem_lit, count = elem_samples(c["go_elem"])
        arg = lower_first(c["resp_field"])
        lines.append("func Test%sToProto(t *testing.T) {" % name)
        lines.append("\tgot := %sToProto(%s)" % (name, elem_lit))
        lines.append("")
        lines.append("\tif got == nil {")
        lines.append('\t\tt.Fatal("%sToProto() = nil, want value")' % name)
        lines.append("\t}")
        lines.append("")
        lines.append("\tif len(got.Get%s()) != %s {" % (c["resp_field"], count))
        lines.append('\t\tt.Errorf("len(%s) = %%d, want %s", len(got.Get%s()))' % (c["resp_field"], count, c["resp_field"]))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


# ---------- validation ----------
def validation_call(f):
    col, getter = f["col"], "req.Get%s()" % f["proto"]
    if f["go_type"] == "string":
        if f["varchar"] is None:
            raise Exception("нет границы varchar для %s" % col)
        return 'v.StringVarchar("%s", %s, %s)' % (col, getter, limit_const(col))
    if f["go_type"] == "int64":
        return 'v.Int64ID("%s", %s)' % (col, getter)
    if f["go_type"] == "uuid.UUID":
        return 'v.StringUUID("%s", %s)' % (col, getter)
    if f["go_type"] == "bool":
        return None
    raise Exception("нет правила валидации для %s (%s)" % (col, f["go_type"]))


def gen_validation(meta, calls):
    proto = meta["proto_alias"]

    lines = [
        "package %s" % meta["conv_pkg"],
        "",
        "import (",
        "\t" + meta["proto_import"],
        '\t"%s"' % VALIDATOR_IMPORT,
        ")",
        "",
    ]

    for c in calls:
        name = c["query"]
        lines.append(
            "// Validate%s проверяет запрос на выборку %s.%s." % (name, meta["sqlc"], name)
        )
        lines.append("func Validate%s(req *%s.%sRequest) error {" % (name, proto, name))
        lines.append("\tv := validator.New()")
        lines.append("")
        lines.append("\tif req == nil {")
        lines.append('\t\tv.AddError("request", validator.MsgRequired)')
        lines.append("\t\treturn v.Err()")
        lines.append("\t}")
        lines.append("")
        for f in c["req_fields"]:
            call = validation_call(f)
            if call:
                lines.append("\t" + call)
        lines.append("")
        lines.append("\treturn v.Err()")
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


def gen_limits(meta, calls):
    seen = {}
    for c in calls:
        for f in c["req_fields"]:
            if f["go_type"] != "string":
                continue
            if f["col"] in seen and seen[f["col"]] != f["varchar"]:
                raise Exception(
                    "колонка %s встречается с разной длиной (%s и %s) в разных запросах %s — "
                    "нужен отдельный column_hints"
                    % (f["col"], seen[f["col"]], f["varchar"], meta["sqlc"])
                )
            seen[f["col"]] = f["varchar"]

    lines = [
        "package %s" % meta["conv_pkg"],
        "",
        "// Ограничения длины строковых аргументов. Значения взяты из",
        "// datacatalogue/db/sqlc/%s/schema.sql и должны меняться вместе" % meta["sqlc"],
        "// с ней: валидация обязана отсекать слишком длинную строку раньше,",
        "// чем её отвергнет Postgres.",
        "const (",
    ]
    for col, varchar in seen.items():
        lines.append("\t%s = %d // %s character varying(%d)" % (limit_const(col), varchar, col, varchar))
    lines.append(")")
    lines.append("")
    return "\n".join(lines)


def gen_validation_test(meta, calls):
    proto = meta["proto_alias"]

    lines = [
        "package %s" % meta["conv_pkg"],
        "",
        "import (",
        '\t"errors"',
        '\t"testing"',
        "",
        "\t" + meta["proto_import"],
        '\t"%s"' % VALIDATOR_IMPORT,
        ")",
        "",
        "// authLogicFieldErrors достаёт из ошибки список полей с претензиями.",
        "func authLogicFieldErrors(t *testing.T, err error) map[string][]string {",
        "\tt.Helper()",
        "",
        "\tvar validationErr *validator.ValidationError",
        "\tif !errors.As(err, &validationErr) {",
        '\t\tt.Fatalf("error = %v, want *validator.ValidationError", err)',
        "\t}",
        "",
        "\treturn validationErr.Errors",
        "}",
        "",
    ]

    for c in calls:
        name = c["query"]

        lines.append("func valid%sRequest() *%s.%sRequest {" % (name, proto, name))
        lines.append("\treturn &%s.%sRequest{" % (proto, name))
        for i, f in enumerate(c["req_fields"]):
            lines.append("\t\t%s: %s," % (f["proto"], sample_proto(f, i)))
        lines.append("\t}")
        lines.append("}")
        lines.append("")

        lines.append("func TestValidate%s(t *testing.T) {" % name)
        lines.append("\ttests := []struct {")
        lines.append("\t\tname      string")
        lines.append("\t\tmutate    func(*%s.%sRequest)" % (proto, name))
        lines.append("\t\twantField string")
        lines.append("\t}{")
        lines.append('\t\t{name: "valid", mutate: func(*%s.%sRequest) {}},' % (proto, name))
        for f in c["req_fields"]:
            col = f["col"]
            if f["go_type"] == "string":
                lines.append(
                    '\t\t{name: "empty %s", mutate: func(r *%s.%sRequest) { r.%s = "" }, wantField: "%s"},'
                    % (col, proto, name, f["proto"], col)
                )
                lines.append(
                    '\t\t{name: "%s too long", mutate: func(r *%s.%sRequest) { r.%s = strings.Repeat("a", %s+1) }, wantField: "%s"},'
                    % (col, proto, name, f["proto"], limit_const(col), col)
                )
            elif f["go_type"] == "int64":
                lines.append(
                    '\t\t{name: "zero %s", mutate: func(r *%s.%sRequest) { r.%s = 0 }, wantField: "%s"},'
                    % (col, proto, name, f["proto"], col)
                )
            elif f["go_type"] == "uuid.UUID":
                lines.append(
                    '\t\t{name: "empty %s", mutate: func(r *%s.%sRequest) { r.%s = "" }, wantField: "%s"},'
                    % (col, proto, name, f["proto"], col)
                )
                lines.append(
                    '\t\t{name: "malformed %s", mutate: func(r *%s.%sRequest) { r.%s = "not-a-uuid" }, wantField: "%s"},'
                    % (col, proto, name, f["proto"], col)
                )
        lines.append("\t}")
        lines.append("")
        lines.append("\tfor _, tt := range tests {")
        lines.append("\t\tt.Run(tt.name, func(t *testing.T) {")
        lines.append("\t\t\treq := valid%sRequest()" % name)
        lines.append("\t\t\ttt.mutate(req)")
        lines.append("")
        lines.append("\t\t\terr := Validate%s(req)" % name)
        lines.append("")
        lines.append('\t\t\tif tt.wantField == "" {')
        lines.append("\t\t\t\tif err != nil {")
        lines.append('\t\t\t\t\tt.Fatalf("Validate%s() = %%v, want nil", err)' % name)
        lines.append("\t\t\t\t}")
        lines.append("\t\t\t\treturn")
        lines.append("\t\t\t}")
        lines.append("")
        lines.append("\t\t\tif err == nil {")
        lines.append('\t\t\t\tt.Fatalf("Validate%s() = nil, want error on %%q", tt.wantField)' % name)
        lines.append("\t\t\t}")
        lines.append("")
        lines.append("\t\t\tfields := authLogicFieldErrors(t, err)")
        lines.append("\t\t\tif len(fields[tt.wantField]) == 0 {")
        lines.append('\t\t\t\tt.Errorf("no error on %q, got %v", tt.wantField, fields)')
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
    if "strings.Repeat(" in body:
        body = body.replace('\t"errors"\n\t"testing"', '\t"errors"\n\t"strings"\n\t"testing"')
    return body


# ---------- service ----------
def gen_service(meta, calls):
    field = meta["repo_field"]
    st = meta["service_type"]
    recv = "s"

    lines = [
        "package %s" % meta["service_pkg"],
        "",
        "import (",
        '\t"context"',
        '\t"fmt"',
        "",
        '\t"%s"' % meta["repo_import"],
        ")",
        "",
    ]

    for c in calls:
        name = c["query"]
        repo = meta["repo_alias"]
        ret = ("[]" + c["go_elem"]) if c["is_slice"] else c["go_elem"]
        lines.append("// %s оборачивает sqlc-запрос %s.%s." % (name, meta["sqlc"], name))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, params %s.%s) (%s, error) {"
            % (recv, st, name, repo, c["params_type"], ret)
        )
        lines.append("\trows, err := %s.%s.%s(ctx, params)" % (recv, field, name))
        lines.append("\tif err != nil {")
        lines.append('\t\treturn nil, fmt.Errorf("%s: %%w", err)' % name)
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn rows, nil")
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


# ---------- api ----------
def gen_api(meta, calls):
    proto = meta["proto_alias"]
    recv = meta["api_recv"]
    at = meta["api_type"]
    svc = meta["service_field"]
    conv = meta["conv_pkg"] + "conv"
    valid = meta["conv_pkg"] + "validation"

    lines = [
        "package %s" % meta["api_pkg"],
        "",
        "import (",
        '\t"context"',
        "",
        '\t"%s"' % APIERR_IMPORT,
        '\t%s "%s/%s"' % (conv, CONV_IMPORT, meta["conv_pkg"]),
        '\t%s "%s/%s"' % (valid, VALID_IMPORT, meta["conv_pkg"]),
        "\t" + meta["proto_import"],
        ")",
        "",
    ]

    for c in calls:
        name = c["query"]
        arg = lower_first(c["resp_field"])
        lines.append("// %s оборачивает sqlc-запрос %s.%s." % (name, meta["sqlc"], name))
        lines.append(
            "func (%s *%s) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {"
            % (recv, at, name, proto, name, proto, name)
        )
        lines.append("\tif err := %s.Validate%s(req); err != nil {" % (valid, name))
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append(
            "\t%s, err := %s.services.%s.%s(ctx, %s.To%sParams(req))" % (arg, recv, svc, name, conv, name)
        )
        lines.append("\tif err != nil {")
        lines.append("\t\treturn nil, apierror.Wrap(err)")
        lines.append("\t}")
        lines.append("")
        lines.append("\treturn %s.%sToProto(%s), nil" % (conv, name, arg))
        lines.append("}")
        lines.append("")

    return "\n".join(lines)


def main():
    global INTERNAL, CONV_IMPORT, VALID_IMPORT, ERRORS_IMPORT, VALIDATOR_IMPORT, APIERR_IMPORT

    args, config = parse_args("Генерация Go-слоёв по model_custom.json")

    INTERNAL = repo_path(config["internal"])
    shared = config["shared_imports"]
    CONV_IMPORT = shared["converter"]
    VALID_IMPORT = shared["validation"]
    ERRORS_IMPORT = shared["custom_errors"]
    VALIDATOR_IMPORT = shared["validator"]
    APIERR_IMPORT = shared["apierror"]

    model = read_json(args.work_dir, "model_custom.json")
    if not model:
        print("model_custom.json пуст — писать нечего")
        return

    for sqlc, d in model.items():
        meta = dict(d["meta"])
        meta["sqlc"] = sqlc
        calls = d["calls"]

        write("%s/%s.go" % (meta["conv_dir"], sqlc), gen_converter(meta, calls))
        write("%s/%s_test.go" % (meta["conv_dir"], sqlc), gen_converter_test(meta, calls))

        write("%s/%s.go" % (meta["valid_dir"], sqlc), gen_validation(meta, calls))
        write("%s/%s_test.go" % (meta["valid_dir"], sqlc), gen_validation_test(meta, calls))
        write("%s/limits.go" % meta["valid_dir"], gen_limits(meta, calls))

        write("%s/%s.go" % (meta["service_dir"], sqlc), gen_service(meta, calls))
        write("%s/%s.go" % (meta["api_dir"], sqlc), gen_api(meta, calls))

    print("файлов записано:", len(written))
    print("каталог:", INTERNAL)
    print("\nдальше: gofmt -w <internal>, затем go build ./... и go test ./internal/...")


if __name__ == "__main__":
    main()
