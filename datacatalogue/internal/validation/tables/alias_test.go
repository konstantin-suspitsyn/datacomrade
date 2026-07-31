package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateAliasRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateAliasRequest() *tablesv1.CreateAliasRequest {
	return &tablesv1.CreateAliasRequest{
		Name:           "name-0",
		Description:    "description-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

func validUpdateAliasByIdRequest() *tablesv1.UpdateAliasByIdRequest {
	return &tablesv1.UpdateAliasByIdRequest{
		Id:             42,
		Name:           "name-0",
		Description:    "description-1",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

// aliasFieldErrors достаёт из ошибки список полей с претензиями.
func aliasFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateAlias(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateAliasRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateAliasRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateAliasRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateAliasRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateAliasRequest) { r.Name = strings.Repeat("a", aliasNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateAliasRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateAliasRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateAliasRequest) { r.Description = strings.Repeat("a", aliasDescriptionMaxLen+1) }, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateAliasRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateAliasRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateAliasRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateAliasRequest()
			tt.mutate(req)

			err := ValidateCreateAlias(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateAlias() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateAlias() = nil, want error on %q", tt.wantField)
			}

			fields := aliasFieldErrors(t, err)

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

func TestValidateUpdateAliasById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateAliasByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateAliasByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Name = strings.Repeat("a", aliasNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateAliasByIdRequest) {
			r.Description = strings.Repeat("a", aliasDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateAliasByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateAliasById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateAliasById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateAliasById() = nil, want error on %q", tt.wantField)
			}

			fields := aliasFieldErrors(t, err)

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
func TestValidateCreateAliasAtVarcharLimit(t *testing.T) {
	req := validCreateAliasRequest()
	req.Name = strings.Repeat("a", aliasNameMaxLen)

	if err := ValidateCreateAlias(req); err != nil {
		t.Errorf("ValidateCreateAlias() = %v, want nil at exactly %d chars", err, aliasNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateAliasCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateAliasRequest()
	req.Name = strings.Repeat("я", aliasNameMaxLen)

	if err := ValidateCreateAlias(req); err != nil {
		t.Errorf("ValidateCreateAlias() = %v, want nil at exactly %d cyrillic chars", err, aliasNameMaxLen)
	}
}

func TestValidateCreateAliasCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateAlias(&tablesv1.CreateAliasRequest{})

	if err == nil {
		t.Fatal("ValidateCreateAlias() = nil, want errors")
	}

	fields := aliasFieldErrors(t, err)

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

func TestValidateCreateAliasNil(t *testing.T) {
	if err := ValidateCreateAlias(nil); err == nil {
		t.Error("ValidateCreateAlias(nil) = nil, want error")
	}
}

func TestValidateUpdateAliasByIdNil(t *testing.T) {
	if err := ValidateUpdateAliasById(nil); err == nil {
		t.Error("ValidateUpdateAliasById(nil) = nil, want error")
	}
}
