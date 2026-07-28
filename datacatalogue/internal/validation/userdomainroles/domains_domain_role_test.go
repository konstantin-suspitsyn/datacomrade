package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateDomainsDomainRoleRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateDomainsDomainRoleRequest() *userdomainrolesv1.CreateDomainsDomainRoleRequest {
	return &userdomainrolesv1.CreateDomainsDomainRoleRequest{
		DomainCatId:   100,
		DomainRolesId: 101,
	}
}

func validUpdateDomainsDomainRoleByIdRequest() *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest {
	return &userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest{
		Id:            42,
		DomainCatId:   100,
		DomainRolesId: 101,
	}
}

// domainsDomainRoleFieldErrors достаёт из ошибки список полей с претензиями.
func domainsDomainRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateDomainsDomainRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateDomainsDomainRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateDomainsDomainRoleRequest) {}},
		{name: "zero domain_cat_id", mutate: func(r *userdomainrolesv1.CreateDomainsDomainRoleRequest) { r.DomainCatId = 0 }, wantField: "domain_cat_id"},
		{name: "negative domain_cat_id", mutate: func(r *userdomainrolesv1.CreateDomainsDomainRoleRequest) { r.DomainCatId = -1 }, wantField: "domain_cat_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.CreateDomainsDomainRoleRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "negative domain_roles_id", mutate: func(r *userdomainrolesv1.CreateDomainsDomainRoleRequest) { r.DomainRolesId = -1 }, wantField: "domain_roles_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateDomainsDomainRoleRequest()
			tt.mutate(req)

			err := ValidateCreateDomainsDomainRole(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDomainsDomainRole() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDomainsDomainRole() = nil, want error on %q", tt.wantField)
			}

			fields := domainsDomainRoleFieldErrors(t, err)

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

func TestValidateUpdateDomainsDomainRoleById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) {}},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero domain_cat_id", mutate: func(r *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) { r.DomainCatId = 0 }, wantField: "domain_cat_id"},
		{name: "negative domain_cat_id", mutate: func(r *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) { r.DomainCatId = -1 }, wantField: "domain_cat_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "negative domain_roles_id", mutate: func(r *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) { r.DomainRolesId = -1 }, wantField: "domain_roles_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateDomainsDomainRoleByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDomainsDomainRoleById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDomainsDomainRoleById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDomainsDomainRoleById() = nil, want error on %q", tt.wantField)
			}

			fields := domainsDomainRoleFieldErrors(t, err)

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

func TestValidateCreateDomainsDomainRoleCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateDomainsDomainRole(&userdomainrolesv1.CreateDomainsDomainRoleRequest{})

	if err == nil {
		t.Fatal("ValidateCreateDomainsDomainRole() = nil, want errors")
	}

	fields := domainsDomainRoleFieldErrors(t, err)

	wantFields := []string{"domain_cat_id", "domain_roles_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateDomainsDomainRoleNil(t *testing.T) {
	if err := ValidateCreateDomainsDomainRole(nil); err == nil {
		t.Error("ValidateCreateDomainsDomainRole(nil) = nil, want error")
	}
}

func TestValidateUpdateDomainsDomainRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateDomainsDomainRoleById(nil); err == nil {
		t.Error("ValidateUpdateDomainsDomainRoleById(nil) = nil, want error")
	}
}
