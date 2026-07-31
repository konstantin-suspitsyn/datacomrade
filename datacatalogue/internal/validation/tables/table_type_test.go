package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateTableTypeRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateTableTypeRequest() *tablesv1.CreateTableTypeRequest {
	return &tablesv1.CreateTableTypeRequest{
		Name:           "name-0",
		Description:    "description-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

func validUpdateTableTypeByIdRequest() *tablesv1.UpdateTableTypeByIdRequest {
	return &tablesv1.UpdateTableTypeByIdRequest{
		Id:             42,
		Name:           "name-0",
		Description:    "description-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

// tableTypeFieldErrors достаёт из ошибки список полей с претензиями.
func tableTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateTableType(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateTableTypeRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateTableTypeRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Name = strings.Repeat("a", tableTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateTableTypeRequest) {
			r.Description = strings.Repeat("a", tableTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateTableTypeRequest()
			tt.mutate(req)

			err := ValidateCreateTableType(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateTableType() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateTableType() = nil, want error on %q", tt.wantField)
			}

			fields := tableTypeFieldErrors(t, err)

			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateTableTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateTableTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateTableTypeByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Name = strings.Repeat("a", tableTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) {
			r.Description = strings.Repeat("a", tableTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateTableTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateTableTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateTableTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateTableTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := tableTypeFieldErrors(t, err)

			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
	}
}

// Ровно граничная длина проходит: varchar(n) допускает n символов.
func TestValidateCreateTableTypeAtVarcharLimit(t *testing.T) {
	req := validCreateTableTypeRequest()
	req.Name = strings.Repeat("a", tableTypeNameMaxLen)

	if err := ValidateCreateTableType(req); err != nil {
		t.Errorf("ValidateCreateTableType() = %v, want nil at exactly %d chars", err, tableTypeNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateTableTypeCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateTableTypeRequest()
	req.Name = strings.Repeat("я", tableTypeNameMaxLen)

	if err := ValidateCreateTableType(req); err != nil {
		t.Errorf("ValidateCreateTableType() = %v, want nil at exactly %d cyrillic chars", err, tableTypeNameMaxLen)
	}
}

func TestValidateCreateTableTypeCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateTableType(&tablesv1.CreateTableTypeRequest{})

	if err == nil {
		t.Fatal("ValidateCreateTableType() = nil, want errors")
	}

	fields := tableTypeFieldErrors(t, err)

	wantFields := []string{"name", "description", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateTableTypeNil(t *testing.T) {
	if err := ValidateCreateTableType(nil); err == nil {
		t.Error("ValidateCreateTableType(nil) = nil, want error")
	}
}

func TestValidateUpdateTableTypeByIdNil(t *testing.T) {
	if err := ValidateUpdateTableTypeById(nil); err == nil {
		t.Error("ValidateUpdateTableTypeById(nil) = nil, want error")
	}
}
