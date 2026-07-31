package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateTableCatRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateTableCatRequest() *tablesv1.CreateTableCatRequest {
	return &tablesv1.CreateTableCatRequest{
		Name:           "name-0",
		Description:    "description-1",
		SchemaId:       102,
		TableTypeId:    103,
		DomainId:       104,
		IsGetDict:      true,
		UserExternalId: "00000000-0000-4000-8000-000000000007",
	}
}

func validUpdateTableCatByIdRequest() *tablesv1.UpdateTableCatByIdRequest {
	return &tablesv1.UpdateTableCatByIdRequest{
		Id:             42,
		Name:           "name-0",
		Description:    "description-1",
		SchemaId:       102,
		TableTypeId:    103,
		DomainId:       104,
		IsGetDict:      true,
		UserExternalId: "00000000-0000-4000-8000-000000000007",
	}
}

// tableCatFieldErrors достаёт из ошибки список полей с претензиями.
func tableCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateTableCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateTableCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateTableCatRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Name = strings.Repeat("a", tableCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateTableCatRequest) {
			r.Description = strings.Repeat("a", tableCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero schema_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.SchemaId = 0 }, wantField: "schema_id"},
		{name: "negative schema_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.SchemaId = -1 }, wantField: "schema_id"},
		{name: "zero table_type_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.TableTypeId = 0 }, wantField: "table_type_id"},
		{name: "negative table_type_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.TableTypeId = -1 }, wantField: "table_type_id"},
		{name: "zero domain_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "negative domain_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.DomainId = -1 }, wantField: "domain_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateTableCatRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateTableCatRequest()
			tt.mutate(req)

			err := ValidateCreateTableCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateTableCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateTableCat() = nil, want error on %q", tt.wantField)
			}

			fields := tableCatFieldErrors(t, err)

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

func TestValidateUpdateTableCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateTableCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateTableCatByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Name = strings.Repeat("a", tableCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) {
			r.Description = strings.Repeat("a", tableCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero schema_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.SchemaId = 0 }, wantField: "schema_id"},
		{name: "negative schema_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.SchemaId = -1 }, wantField: "schema_id"},
		{name: "zero table_type_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.TableTypeId = 0 }, wantField: "table_type_id"},
		{name: "negative table_type_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.TableTypeId = -1 }, wantField: "table_type_id"},
		{name: "zero domain_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "negative domain_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.DomainId = -1 }, wantField: "domain_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateTableCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateTableCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateTableCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateTableCatById() = nil, want error on %q", tt.wantField)
			}

			fields := tableCatFieldErrors(t, err)

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
func TestValidateCreateTableCatAtVarcharLimit(t *testing.T) {
	req := validCreateTableCatRequest()
	req.Name = strings.Repeat("a", tableCatNameMaxLen)

	if err := ValidateCreateTableCat(req); err != nil {
		t.Errorf("ValidateCreateTableCat() = %v, want nil at exactly %d chars", err, tableCatNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateTableCatCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateTableCatRequest()
	req.Name = strings.Repeat("я", tableCatNameMaxLen)

	if err := ValidateCreateTableCat(req); err != nil {
		t.Errorf("ValidateCreateTableCat() = %v, want nil at exactly %d cyrillic chars", err, tableCatNameMaxLen)
	}
}

func TestValidateCreateTableCatCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateTableCat(&tablesv1.CreateTableCatRequest{})

	if err == nil {
		t.Fatal("ValidateCreateTableCat() = nil, want errors")
	}

	fields := tableCatFieldErrors(t, err)

	wantFields := []string{"name", "description", "schema_id", "table_type_id", "domain_id", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateTableCatNil(t *testing.T) {
	if err := ValidateCreateTableCat(nil); err == nil {
		t.Error("ValidateCreateTableCat(nil) = nil, want error")
	}
}

func TestValidateUpdateTableCatByIdNil(t *testing.T) {
	if err := ValidateUpdateTableCatById(nil); err == nil {
		t.Error("ValidateUpdateTableCatById(nil) = nil, want error")
	}
}
