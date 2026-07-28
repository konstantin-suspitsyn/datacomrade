package validator_test

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

func TestStringChecks(t *testing.T) {
	tests := []struct {
		name  string
		check func(v *validator.Validator) bool
		want  bool
	}{
		{"required ок", func(v *validator.Validator) bool { return v.StringRequired("name", "sales") }, true},
		{"required пусто", func(v *validator.Validator) bool { return v.StringRequired("name", "") }, false},
		{"required пробелы", func(v *validator.Validator) bool { return v.StringRequired("name", "   ") }, false},
		{"max len ок", func(v *validator.Validator) bool { return v.StringMaxLen("name", "abc", 3) }, true},
		{"max len превышен", func(v *validator.Validator) bool { return v.StringMaxLen("name", "abcd", 3) }, false},
		{"max len кириллица считается в символах", func(v *validator.Validator) bool { return v.StringMaxLen("name", "продажи", 7) }, true},
		{"min len ок", func(v *validator.Validator) bool { return v.StringMinLen("name", "abc", 3) }, true},
		{"min len не хватает", func(v *validator.Validator) bool { return v.StringMinLen("name", "ab", 3) }, false},
		{"exact len ок", func(v *validator.Validator) bool { return v.StringExactLen("code", "RU", 2) }, true},
		{"exact len не совпал", func(v *validator.Validator) bool { return v.StringExactLen("code", "RUS", 2) }, false},
		{"varchar ок", func(v *validator.Validator) bool { return v.StringVarchar("name", "sales", 128) }, true},
		{"varchar пусто", func(v *validator.Validator) bool { return v.StringVarchar("name", "", 128) }, false},
		{"varchar длинный", func(v *validator.Validator) bool { return v.StringVarchar("name", strings.Repeat("a", 129), 128) }, false},
		{"optional varchar пусто", func(v *validator.Validator) bool { return v.StringOptionalVarchar("descr", "", 10) }, true},
		{"optional varchar длинный", func(v *validator.Validator) bool { return v.StringOptionalVarchar("descr", "abcdefghijk", 10) }, false},
		{"in ок", func(v *validator.Validator) bool { return v.StringIn("type", "view", "table", "view") }, true},
		{"in не найдено", func(v *validator.Validator) bool { return v.StringIn("type", "index", "table", "view") }, false},
		{"matches ок", func(v *validator.Validator) bool {
			return v.StringMatches("name", "ab1", validator.IdentifierRX, "идентификатор")
		}, true},
		{"matches nil regexp", func(v *validator.Validator) bool {
			return v.StringMatches("name", "ab1", nil, "идентификатор")
		}, false},
		{"email ок", func(v *validator.Validator) bool { return v.StringEmail("email", "user@example.com") }, true},
		{"email плохой", func(v *validator.Validator) bool { return v.StringEmail("email", "user@@example") }, false},
		{"identifier ок", func(v *validator.Validator) bool { return v.StringIdentifier("table", "dc_table_cat") }, true},
		{"identifier с цифры", func(v *validator.Validator) bool { return v.StringIdentifier("table", "1table") }, false},
		{"identifier кириллица", func(v *validator.Validator) bool { return v.StringIdentifier("table", "таблица") }, false},
		{"no spaces ок", func(v *validator.Validator) bool { return v.StringNoSpaces("login", "user") }, true},
		{"no spaces с табом", func(v *validator.Validator) bool { return v.StringNoSpaces("login", "us\ter") }, false},
		{"trimmed ок", func(v *validator.Validator) bool { return v.StringTrimmed("name", "sales") }, true},
		{"trimmed с пробелом", func(v *validator.Validator) bool { return v.StringTrimmed("name", " sales") }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			got := tt.check(v)
			if got != tt.want {
				t.Errorf("проверка вернула %v, ожидалось %v", got, tt.want)
			}

			if v.Valid() != tt.want {
				t.Errorf("Valid() = %v, ожидалось %v (ошибки: %v)", v.Valid(), tt.want, v.Errors())
			}
		})
	}
}

func TestNumberChecks(t *testing.T) {
	tests := []struct {
		name  string
		check func(v *validator.Validator) bool
		want  bool
	}{
		{"int64 required ок", func(v *validator.Validator) bool { return v.Int64Required("user_id", 5) }, true},
		{"int64 required ноль", func(v *validator.Validator) bool { return v.Int64Required("user_id", 0) }, false},
		{"int64 id ок", func(v *validator.Validator) bool { return v.Int64ID("schema_id", 12) }, true},
		{"int64 id ноль", func(v *validator.Validator) bool { return v.Int64ID("schema_id", 0) }, false},
		{"int64 id отрицательный", func(v *validator.Validator) bool { return v.Int64ID("schema_id", -1) }, false},
		{"int64 min ок", func(v *validator.Validator) bool { return v.Int64Min("port", 5432, 1) }, true},
		{"int64 max превышен", func(v *validator.Validator) bool { return v.Int64Max("port", 70000, 65535) }, false},
		{"int64 between ок", func(v *validator.Validator) bool { return v.Int64Between("port", 5432, 1, 65535) }, true},
		{"int64 between ниже", func(v *validator.Validator) bool { return v.Int64Between("port", 0, 1, 65535) }, false},
		{"int64 in ок", func(v *validator.Validator) bool { return v.Int64In("type_id", 2, 1, 2, 3) }, true},
		{"int64 in не найдено", func(v *validator.Validator) bool { return v.Int64In("type_id", 9, 1, 2, 3) }, false},
		{"int64 в int32 ок", func(v *validator.Validator) bool { return v.Int64FitsInt32("port", 65535) }, true},
		{"int64 в int32 не влезает", func(v *validator.Validator) bool { return v.Int64FitsInt32("port", math.MaxInt32+1) }, false},
		{"int32 positive ок", func(v *validator.Validator) bool { return v.Int32Positive("port", 5432) }, true},
		{"int32 positive ноль", func(v *validator.Validator) bool { return v.Int32Positive("port", 0) }, false},
		{"float64 positive ок", func(v *validator.Validator) bool { return v.Float64Positive("ratio", 0.5) }, true},
		{"float64 positive ноль", func(v *validator.Validator) bool { return v.Float64Positive("ratio", 0) }, false},
		{"float64 between ок", func(v *validator.Validator) bool { return v.Float64Between("ratio", 0.5, 0, 1) }, true},
		{"float64 NaN", func(v *validator.Validator) bool { return v.Float64Between("ratio", math.NaN(), 0, 1) }, false},
		{"float64 бесконечность", func(v *validator.Validator) bool { return v.Float64Finite("ratio", math.Inf(1)) }, false},
		{"generic number min", func(v *validator.Validator) bool { return validator.NumberMin(v, "n", 3, 5) }, false},
		{"generic number in", func(v *validator.Validator) bool { return validator.NumberIn(v, "n", uint8(2), 1, 2) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			got := tt.check(v)
			if got != tt.want {
				t.Errorf("проверка вернула %v, ожидалось %v", got, tt.want)
			}

			if v.Valid() != tt.want {
				t.Errorf("Valid() = %v, ожидалось %v (ошибки: %v)", v.Valid(), tt.want, v.Errors())
			}
		})
	}
}

func TestTimeAndBoolChecks(t *testing.T) {
	now := time.Date(2026, 7, 28, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name  string
		check func(v *validator.Validator) bool
		want  bool
	}{
		{"time required ок", func(v *validator.Validator) bool { return v.TimeRequired("created_at", now) }, true},
		{"time required нулевая", func(v *validator.Validator) bool { return v.TimeRequired("created_at", time.Time{}) }, false},
		{"time after ок", func(v *validator.Validator) bool { return v.TimeAfter("updated_at", now, past) }, true},
		{"time after не позже", func(v *validator.Validator) bool { return v.TimeAfter("updated_at", past, now) }, false},
		{"time before ок", func(v *validator.Validator) bool { return v.TimeBefore("created_at", past, now) }, true},
		{"time between ок", func(v *validator.Validator) bool { return v.TimeBetween("created_at", now, past, future) }, true},
		{"time between вне диапазона", func(v *validator.Validator) bool { return v.TimeBetween("created_at", future, past, now) }, false},
		{"не в будущем ок", func(v *validator.Validator) bool {
			return v.TimeNotInFuture("created_at", time.Now().Add(-time.Minute))
		}, true},
		{"не в будущем нарушено", func(v *validator.Validator) bool { return v.TimeNotInFuture("created_at", time.Now().Add(time.Hour)) }, false},
		{"bool true ок", func(v *validator.Validator) bool { return v.BoolTrue("is_active", true) }, true},
		{"bool true нарушено", func(v *validator.Validator) bool { return v.BoolTrue("is_active", false) }, false},
		{"bool false ок", func(v *validator.Validator) bool { return v.BoolFalse("is_deleted", false) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			got := tt.check(v)
			if got != tt.want {
				t.Errorf("проверка вернула %v, ожидалось %v", got, tt.want)
			}

			if v.Valid() != tt.want {
				t.Errorf("Valid() = %v, ожидалось %v (ошибки: %v)", v.Valid(), tt.want, v.Errors())
			}
		})
	}
}

func TestSliceAndOptionalChecks(t *testing.T) {
	name := "sales"
	long := strings.Repeat("a", 20)
	id := int64(7)
	zeroID := int64(0)

	tests := []struct {
		name  string
		check func(v *validator.Validator) bool
		want  bool
	}{
		{"slice required ок", func(v *validator.Validator) bool { return validator.SliceRequired(v, "columns", []string{"a"}) }, true},
		{"slice required пусто", func(v *validator.Validator) bool { return validator.SliceRequired(v, "columns", []string{}) }, false},
		{"slice min len", func(v *validator.Validator) bool { return validator.SliceMinLen(v, "columns", []int{1}, 2) }, false},
		{"slice max len", func(v *validator.Validator) bool { return validator.SliceMaxLen(v, "columns", []int{1, 2}, 2) }, true},
		{"slice unique ок", func(v *validator.Validator) bool { return validator.SliceUnique(v, "ids", []int64{1, 2, 3}) }, true},
		{"slice unique дубликат", func(v *validator.Validator) bool { return validator.SliceUnique(v, "ids", []int64{1, 2, 1}) }, false},
		{"optional string nil", func(v *validator.Validator) bool { return v.OptionalStringVarchar("descr", nil, 10) }, true},
		{"optional string ок", func(v *validator.Validator) bool { return v.OptionalStringVarchar("descr", &name, 10) }, true},
		{"optional string длинный", func(v *validator.Validator) bool { return v.OptionalStringVarchar("descr", &long, 10) }, false},
		{"optional id nil", func(v *validator.Validator) bool { return v.OptionalInt64ID("alias_id", nil) }, true},
		{"optional id ок", func(v *validator.Validator) bool { return v.OptionalInt64ID("alias_id", &id) }, true},
		{"optional id ноль", func(v *validator.Validator) bool { return v.OptionalInt64ID("alias_id", &zeroID) }, false},
		{
			"generic optional",
			func(v *validator.Validator) bool {
				return validator.Optional(v, "descr", &long, func(v *validator.Validator, field, value string) bool {
					return v.StringMaxLen(field, value, 10)
				})
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			got := tt.check(v)
			if got != tt.want {
				t.Errorf("проверка вернула %v, ожидалось %v", got, tt.want)
			}

			if v.Valid() != tt.want {
				t.Errorf("Valid() = %v, ожидалось %v (ошибки: %v)", v.Valid(), tt.want, v.Errors())
			}
		})
	}
}

func TestSliceEach(t *testing.T) {
	v := validator.New()

	ok := validator.SliceEach(v, "columns", []string{"good", ""}, func(v *validator.Validator, field, value string) bool {
		return v.StringRequired(field, value)
	})

	if ok {
		t.Fatal("SliceEach должен был вернуть false")
	}

	if len(v.FieldErrors("columns.1")) != 1 {
		t.Errorf("ожидалась одна ошибка у columns.1, получено: %v", v.Errors())
	}

	if len(v.FieldErrors("columns.0")) != 0 {
		t.Errorf("у columns.0 не должно быть ошибок, получено: %v", v.FieldErrors("columns.0"))
	}
}

func TestErrorsCollection(t *testing.T) {
	v := validator.New()

	v.StringVarchar("name", "", 128)
	v.StringVarchar("description", strings.Repeat("d", 2001), 2000)
	v.Int64ID("schema_id", 0)
	v.Int64ID("domain_id", 42)

	if v.Valid() {
		t.Fatal("Valid() = true, ожидалось false")
	}

	if !v.HasErrors() {
		t.Error("HasErrors() = false, ожидалось true")
	}

	errs := v.Errors()

	if len(errs) != 3 {
		t.Errorf("ожидалось 3 поля с ошибками, получено %d: %v", len(errs), errs)
	}

	if _, ok := errs["domain_id"]; ok {
		t.Error("у корректного поля domain_id не должно быть ошибок")
	}

	if !v.FieldValid("domain_id") {
		t.Error("FieldValid(domain_id) = false, ожидалось true")
	}

	if v.Count() != 3 {
		t.Errorf("Count() = %d, ожидалось 3", v.Count())
	}

	want := []string{"description", "name", "schema_id"}
	got := v.Fields()

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Fields() = %v, ожидалось %v", got, want)
			break
		}
	}

	// Мапа должна быть копией: правки снаружи не влияют на валидатор.
	errs["name"] = append(errs["name"], "подделка")

	if len(v.FieldErrors("name")) != 1 {
		t.Error("Errors() вернул ссылку на внутренние данные вместо копии")
	}
}

func TestAddErrorDeduplicates(t *testing.T) {
	v := validator.New()

	v.AddError("name", "поле обязательно")
	v.AddError("name", "поле обязательно")
	v.AddError("name", "слишком длинное")

	if got := len(v.FieldErrors("name")); got != 2 {
		t.Errorf("ожидалось 2 уникальные ошибки, получено %d: %v", got, v.FieldErrors("name"))
	}
}

func TestZeroValueUsable(t *testing.T) {
	var v validator.Validator

	if !v.Valid() {
		t.Error("пустой валидатор должен быть валидным")
	}

	v.StringRequired("name", "")

	if v.Valid() {
		t.Error("после неуспешной проверки Valid() должен быть false")
	}
}

func TestCheck(t *testing.T) {
	v := validator.New()

	if v.Check(false, "name", "своя проверка") {
		t.Error("Check(false) должен вернуть false")
	}

	if !v.Check(true, "other", "своя проверка") {
		t.Error("Check(true) должен вернуть true")
	}

	if len(v.FieldErrors("name")) != 1 || len(v.FieldErrors("other")) != 0 {
		t.Errorf("неожиданные ошибки: %v", v.Errors())
	}
}

func TestMerge(t *testing.T) {
	column := validator.New()
	column.StringRequired("name", "")

	table := validator.New()
	table.Merge("columns.0", column)

	if len(table.FieldErrors("columns.0.name")) != 1 {
		t.Errorf("ожидалась ошибка у columns.0.name, получено: %v", table.Errors())
	}

	empty := validator.New()
	empty.Merge("", column)

	if len(empty.FieldErrors("name")) != 1 {
		t.Errorf("без префикса имя поля должно остаться прежним: %v", empty.Errors())
	}

	empty.Merge("prefix", nil)
}

func TestReset(t *testing.T) {
	v := validator.New()
	v.StringRequired("name", "")

	v.Reset()

	if !v.Valid() {
		t.Errorf("после Reset() валидатор должен быть чистым: %v", v.Errors())
	}
}

func TestErr(t *testing.T) {
	v := validator.New()

	if err := v.Err(); err != nil {
		t.Errorf("для валидного объекта Err() должен быть nil, получено: %v", err)
	}

	v.StringVarchar("name", "", 128)
	v.Int64ID("schema_id", 0)

	err := v.Err()
	if err == nil {
		t.Fatal("Err() = nil, ожидалась ошибка")
	}

	validationErr, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("ожидался *validator.ValidationError, получен %T", err)
	}

	if len(validationErr.Errors) != 2 {
		t.Errorf("ожидалось 2 поля в ошибке, получено %d", len(validationErr.Errors))
	}

	message := err.Error()

	if !strings.Contains(message, "name:") || !strings.Contains(message, "schema_id:") {
		t.Errorf("в тексте ошибки должны быть имена полей, получено: %q", message)
	}

	// Порядок полей детерминированный.
	if strings.Index(message, "name:") > strings.Index(message, "schema_id:") {
		t.Errorf("поля должны идти в алфавитном порядке: %q", message)
	}
}
