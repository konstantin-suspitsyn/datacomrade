package userdomainroles

import (
	"errors"
	"strings"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateDomainRoleRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateDomainRoleRequest() *userdomainrolesv1.CreateDomainRoleRequest {
	return &userdomainrolesv1.CreateDomainRoleRequest{
		Name:        "name-0",
		Description: "description-1",
	}
}

func validUpdateDomainRoleByIdRequest() *userdomainrolesv1.UpdateDomainRoleByIdRequest {
	return &userdomainrolesv1.UpdateDomainRoleByIdRequest{
		Id:          42,
		Name:        "name-0",
		Description: "description-1",
	}
}

// domainRoleFieldErrors достаёт из ошибки список полей с претензиями.
func domainRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateDomainRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateDomainRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateDomainRoleRequest) {}},
		{name: "empty name", mutate: func(r *userdomainrolesv1.CreateDomainRoleRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *userdomainrolesv1.CreateDomainRoleRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *userdomainrolesv1.CreateDomainRoleRequest) {
			r.Name = strings.Repeat("a", domainRoleNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty description", mutate: func(r *userdomainrolesv1.CreateDomainRoleRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *userdomainrolesv1.CreateDomainRoleRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *userdomainrolesv1.CreateDomainRoleRequest) {
			r.Description = strings.Repeat("a", domainRoleDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateDomainRoleRequest()
			tt.mutate(req)

			err := ValidateCreateDomainRole(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDomainRole() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDomainRole() = nil, want error on %q", tt.wantField)
			}

			fields := domainRoleFieldErrors(t, err)

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

func TestValidateUpdateDomainRoleById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateDomainRoleByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateDomainRoleByIdRequest) {}},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) {
			r.Name = strings.Repeat("a", domainRoleNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty description", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *userdomainrolesv1.UpdateDomainRoleByIdRequest) {
			r.Description = strings.Repeat("a", domainRoleDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateDomainRoleByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDomainRoleById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDomainRoleById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDomainRoleById() = nil, want error on %q", tt.wantField)
			}

			fields := domainRoleFieldErrors(t, err)

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
func TestValidateCreateDomainRoleAtVarcharLimit(t *testing.T) {
	req := validCreateDomainRoleRequest()
	req.Name = strings.Repeat("a", domainRoleNameMaxLen)

	if err := ValidateCreateDomainRole(req); err != nil {
		t.Errorf("ValidateCreateDomainRole() = %v, want nil at exactly %d chars", err, domainRoleNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateDomainRoleCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateDomainRoleRequest()
	req.Name = strings.Repeat("я", domainRoleNameMaxLen)

	if err := ValidateCreateDomainRole(req); err != nil {
		t.Errorf("ValidateCreateDomainRole() = %v, want nil at exactly %d cyrillic chars", err, domainRoleNameMaxLen)
	}
}

func TestValidateCreateDomainRoleCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateDomainRole(&userdomainrolesv1.CreateDomainRoleRequest{})

	if err == nil {
		t.Fatal("ValidateCreateDomainRole() = nil, want errors")
	}

	fields := domainRoleFieldErrors(t, err)

	wantFields := []string{"name", "description"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateDomainRoleNil(t *testing.T) {
	if err := ValidateCreateDomainRole(nil); err == nil {
		t.Error("ValidateCreateDomainRole(nil) = nil, want error")
	}
}

func TestValidateUpdateDomainRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateDomainRoleById(nil); err == nil {
		t.Error("ValidateUpdateDomainRoleById(nil) = nil, want error")
	}
}
