package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateDatabaseTypeRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateDatabaseTypeRequest() *tablesv1.CreateDatabaseTypeRequest {
	return &tablesv1.CreateDatabaseTypeRequest{
		Name:           "name-0",
		DbVersion:      "db-version-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

func validUpdateDatabaseTypeByIdRequest() *tablesv1.UpdateDatabaseTypeByIdRequest {
	return &tablesv1.UpdateDatabaseTypeByIdRequest{
		Id:             42,
		Name:           "name-0",
		DbVersion:      "db-version-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

// databaseTypeFieldErrors достаёт из ошибки список полей с претензиями.
func databaseTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateDatabaseType(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDatabaseTypeRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDatabaseTypeRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.Name = strings.Repeat("a", databaseTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty db_version", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.DbVersion = "" }, wantField: "db_version"},
		{name: "blank db_version", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.DbVersion = "   " }, wantField: "db_version"},
		{name: "db_version too long", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) {
			r.DbVersion = strings.Repeat("a", databaseTypeDbVersionMaxLen+1)
		}, wantField: "db_version"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateDatabaseTypeRequest()
			tt.mutate(req)

			err := ValidateCreateDatabaseType(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDatabaseType() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDatabaseType() = nil, want error on %q", tt.wantField)
			}

			fields := databaseTypeFieldErrors(t, err)

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

func TestValidateUpdateDatabaseTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDatabaseTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDatabaseTypeByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) {
			r.Name = strings.Repeat("a", databaseTypeNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty db_version", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.DbVersion = "" }, wantField: "db_version"},
		{name: "blank db_version", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.DbVersion = "   " }, wantField: "db_version"},
		{name: "db_version too long", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) {
			r.DbVersion = strings.Repeat("a", databaseTypeDbVersionMaxLen+1)
		}, wantField: "db_version"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateDatabaseTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDatabaseTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDatabaseTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDatabaseTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseTypeFieldErrors(t, err)

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
func TestValidateCreateDatabaseTypeAtVarcharLimit(t *testing.T) {
	req := validCreateDatabaseTypeRequest()
	req.Name = strings.Repeat("a", databaseTypeNameMaxLen)

	if err := ValidateCreateDatabaseType(req); err != nil {
		t.Errorf("ValidateCreateDatabaseType() = %v, want nil at exactly %d chars", err, databaseTypeNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateDatabaseTypeCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateDatabaseTypeRequest()
	req.Name = strings.Repeat("я", databaseTypeNameMaxLen)

	if err := ValidateCreateDatabaseType(req); err != nil {
		t.Errorf("ValidateCreateDatabaseType() = %v, want nil at exactly %d cyrillic chars", err, databaseTypeNameMaxLen)
	}
}

func TestValidateCreateDatabaseTypeCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateDatabaseType(&tablesv1.CreateDatabaseTypeRequest{})

	if err == nil {
		t.Fatal("ValidateCreateDatabaseType() = nil, want errors")
	}

	fields := databaseTypeFieldErrors(t, err)

	wantFields := []string{"name", "db_version", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateDatabaseTypeNil(t *testing.T) {
	if err := ValidateCreateDatabaseType(nil); err == nil {
		t.Error("ValidateCreateDatabaseType(nil) = nil, want error")
	}
}

func TestValidateUpdateDatabaseTypeByIdNil(t *testing.T) {
	if err := ValidateUpdateDatabaseTypeById(nil); err == nil {
		t.Error("ValidateUpdateDatabaseTypeById(nil) = nil, want error")
	}
}
