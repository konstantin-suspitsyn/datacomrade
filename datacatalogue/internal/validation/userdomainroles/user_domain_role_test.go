package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateUserDomainRoleRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateUserDomainRoleRequest() *userdomainrolesv1.CreateUserDomainRoleRequest {
	return &userdomainrolesv1.CreateUserDomainRoleRequest{
		UserId:              100,
		DomainRolesId:       101,
		DomainId:            102,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000004",
	}
}

func validUpdateUserDomainRoleByIdRequest() *userdomainrolesv1.UpdateUserDomainRoleByIdRequest {
	return &userdomainrolesv1.UpdateUserDomainRoleByIdRequest{
		Id:                  42,
		UserId:              100,
		DomainRolesId:       101,
		DomainId:            102,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000004",
	}
}

// userDomainRoleFieldErrors достаёт из ошибки список полей с претензиями.
func userDomainRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateUserDomainRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateUserDomainRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateUserDomainRoleRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.UserId = -1 }, wantField: "user_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "negative domain_roles_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.DomainRolesId = -1 }, wantField: "domain_roles_id"},
		{name: "zero domain_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "negative domain_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.DomainId = -1 }, wantField: "domain_id"},
		{name: "empty updated_by_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.UpdatedByExternalId = "" }, wantField: "updated_by_id"},
		{name: "malformed updated_by_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.UpdatedByExternalId = "not-a-uuid" }, wantField: "updated_by_id"},
		{name: "updated_by_id without dashes", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) {
			r.UpdatedByExternalId = "00000000000040008000000000000001"
		}, wantField: "updated_by_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateUserDomainRoleRequest()
			tt.mutate(req)

			err := ValidateCreateUserDomainRole(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateUserDomainRole() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateUserDomainRole() = nil, want error on %q", tt.wantField)
			}

			fields := userDomainRoleFieldErrors(t, err)

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

func TestValidateUpdateUserDomainRoleById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateUserDomainRoleByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateUserDomainRoleByIdRequest) {}},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.UserId = -1 }, wantField: "user_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "negative domain_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.DomainRolesId = -1 }, wantField: "domain_roles_id"},
		{name: "zero domain_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "negative domain_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.DomainId = -1 }, wantField: "domain_id"},
		{name: "empty updated_by_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.UpdatedByExternalId = "" }, wantField: "updated_by_id"},
		{name: "malformed updated_by_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.UpdatedByExternalId = "not-a-uuid" }, wantField: "updated_by_id"},
		{name: "updated_by_id without dashes", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) {
			r.UpdatedByExternalId = "00000000000040008000000000000001"
		}, wantField: "updated_by_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateUserDomainRoleByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateUserDomainRoleById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateUserDomainRoleById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateUserDomainRoleById() = nil, want error on %q", tt.wantField)
			}

			fields := userDomainRoleFieldErrors(t, err)

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

func TestValidateCreateUserDomainRoleCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateUserDomainRole(&userdomainrolesv1.CreateUserDomainRoleRequest{})

	if err == nil {
		t.Fatal("ValidateCreateUserDomainRole() = nil, want errors")
	}

	fields := userDomainRoleFieldErrors(t, err)

	wantFields := []string{"user_id", "domain_roles_id", "domain_id", "updated_by_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateUserDomainRoleNil(t *testing.T) {
	if err := ValidateCreateUserDomainRole(nil); err == nil {
		t.Error("ValidateCreateUserDomainRole(nil) = nil, want error")
	}
}

func TestValidateUpdateUserDomainRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateUserDomainRoleById(nil); err == nil {
		t.Error("ValidateUpdateUserDomainRoleById(nil) = nil, want error")
	}
}
