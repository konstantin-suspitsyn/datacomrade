package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateSchemaCatRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateSchemaCatRequest() *tablesv1.CreateSchemaCatRequest {
	return &tablesv1.CreateSchemaCatRequest{
		DatabaseId:     100,
		Name:           "name-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

func validUpdateSchemaCatByIdRequest() *tablesv1.UpdateSchemaCatByIdRequest {
	return &tablesv1.UpdateSchemaCatByIdRequest{
		Id:             42,
		DatabaseId:     100,
		Name:           "name-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

// schemaCatFieldErrors достаёт из ошибки список полей с претензиями.
func schemaCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateSchemaCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateSchemaCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateSchemaCatRequest) {}},
		{name: "zero database_id", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.DatabaseId = 0 }, wantField: "database_id"},
		{name: "negative database_id", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.DatabaseId = -1 }, wantField: "database_id"},
		{name: "empty name", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.Name = strings.Repeat("a", schemaCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateSchemaCatRequest()
			tt.mutate(req)

			err := ValidateCreateSchemaCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateSchemaCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateSchemaCat() = nil, want error on %q", tt.wantField)
			}

			fields := schemaCatFieldErrors(t, err)

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

func TestValidateUpdateSchemaCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateSchemaCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateSchemaCatByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero database_id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.DatabaseId = 0 }, wantField: "database_id"},
		{name: "negative database_id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.DatabaseId = -1 }, wantField: "database_id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Name = strings.Repeat("a", schemaCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateSchemaCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateSchemaCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateSchemaCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateSchemaCatById() = nil, want error on %q", tt.wantField)
			}

			fields := schemaCatFieldErrors(t, err)

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
func TestValidateCreateSchemaCatAtVarcharLimit(t *testing.T) {
	req := validCreateSchemaCatRequest()
	req.Name = strings.Repeat("a", schemaCatNameMaxLen)

	if err := ValidateCreateSchemaCat(req); err != nil {
		t.Errorf("ValidateCreateSchemaCat() = %v, want nil at exactly %d chars", err, schemaCatNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateSchemaCatCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateSchemaCatRequest()
	req.Name = strings.Repeat("я", schemaCatNameMaxLen)

	if err := ValidateCreateSchemaCat(req); err != nil {
		t.Errorf("ValidateCreateSchemaCat() = %v, want nil at exactly %d cyrillic chars", err, schemaCatNameMaxLen)
	}
}

func TestValidateCreateSchemaCatCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateSchemaCat(&tablesv1.CreateSchemaCatRequest{})

	if err == nil {
		t.Fatal("ValidateCreateSchemaCat() = nil, want errors")
	}

	fields := schemaCatFieldErrors(t, err)

	wantFields := []string{"database_id", "name", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateSchemaCatNil(t *testing.T) {
	if err := ValidateCreateSchemaCat(nil); err == nil {
		t.Error("ValidateCreateSchemaCat(nil) = nil, want error")
	}
}

func TestValidateUpdateSchemaCatByIdNil(t *testing.T) {
	if err := ValidateUpdateSchemaCatById(nil); err == nil {
		t.Error("ValidateUpdateSchemaCatById(nil) = nil, want error")
	}
}
