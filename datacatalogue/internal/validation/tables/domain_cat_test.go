package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateDomainCatRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateDomainCatRequest() *tablesv1.CreateDomainCatRequest {
	return &tablesv1.CreateDomainCatRequest{
		DomainName:     "domain-name-0",
		UserExternalId: "00000000-0000-4000-8000-000000000002",
	}
}

func validUpdateDomainCatByIdRequest() *tablesv1.UpdateDomainCatByIdRequest {
	return &tablesv1.UpdateDomainCatByIdRequest{
		Id:             42,
		DomainName:     "domain-name-0",
		UserExternalId: "00000000-0000-4000-8000-000000000002",
	}
}

// domainCatFieldErrors достаёт из ошибки список полей с претензиями.
func domainCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateDomainCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDomainCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDomainCatRequest) {}},
		{name: "empty domain_name", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.DomainName = "" }, wantField: "domain_name"},
		{name: "blank domain_name", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.DomainName = "   " }, wantField: "domain_name"},
		{name: "domain_name too long", mutate: func(r *tablesv1.CreateDomainCatRequest) {
			r.DomainName = strings.Repeat("a", domainCatDomainNameMaxLen+1)
		}, wantField: "domain_name"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateDomainCatRequest()
			tt.mutate(req)

			err := ValidateCreateDomainCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDomainCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDomainCat() = nil, want error on %q", tt.wantField)
			}

			fields := domainCatFieldErrors(t, err)

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

func TestValidateUpdateDomainCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDomainCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDomainCatByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty domain_name", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.DomainName = "" }, wantField: "domain_name"},
		{name: "blank domain_name", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.DomainName = "   " }, wantField: "domain_name"},
		{name: "domain_name too long", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) {
			r.DomainName = strings.Repeat("a", domainCatDomainNameMaxLen+1)
		}, wantField: "domain_name"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.UserExternalId = "00000000000040008000000000000001" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateDomainCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDomainCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDomainCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDomainCatById() = nil, want error on %q", tt.wantField)
			}

			fields := domainCatFieldErrors(t, err)

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
func TestValidateCreateDomainCatAtVarcharLimit(t *testing.T) {
	req := validCreateDomainCatRequest()
	req.DomainName = strings.Repeat("a", domainCatDomainNameMaxLen)

	if err := ValidateCreateDomainCat(req); err != nil {
		t.Errorf("ValidateCreateDomainCat() = %v, want nil at exactly %d chars", err, domainCatDomainNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateDomainCatCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateDomainCatRequest()
	req.DomainName = strings.Repeat("я", domainCatDomainNameMaxLen)

	if err := ValidateCreateDomainCat(req); err != nil {
		t.Errorf("ValidateCreateDomainCat() = %v, want nil at exactly %d cyrillic chars", err, domainCatDomainNameMaxLen)
	}
}

func TestValidateCreateDomainCatCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateDomainCat(&tablesv1.CreateDomainCatRequest{})

	if err == nil {
		t.Fatal("ValidateCreateDomainCat() = nil, want errors")
	}

	fields := domainCatFieldErrors(t, err)

	wantFields := []string{"domain_name", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateDomainCatNil(t *testing.T) {
	if err := ValidateCreateDomainCat(nil); err == nil {
		t.Error("ValidateCreateDomainCat(nil) = nil, want error")
	}
}

func TestValidateUpdateDomainCatByIdNil(t *testing.T) {
	if err := ValidateUpdateDomainCatById(nil); err == nil {
		t.Error("ValidateUpdateDomainCatById(nil) = nil, want error")
	}
}
