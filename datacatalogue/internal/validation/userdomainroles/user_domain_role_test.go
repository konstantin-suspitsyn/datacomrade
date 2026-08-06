package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userDomainRoleFieldErrors достаёт из ошибки список полей с претензиями.
func userDomainRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateUserDomainRoleRequest() *userdomainrolesv1.CreateUserDomainRoleRequest {
	return &userdomainrolesv1.CreateUserDomainRoleRequest{
		UserId:              100,
		DomainRolesId:       101,
		DomainId:            102,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000004",
	}
}

func TestValidateCreateUserDomainRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateUserDomainRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateUserDomainRoleRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "zero domain_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "empty updated_by_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.UpdatedByExternalId = "" }, wantField: "updated_by_id"},
		{name: "malformed updated_by_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRoleRequest) { r.UpdatedByExternalId = "not-a-uuid" }, wantField: "updated_by_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateUserDomainRoleRequest()
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
		})
	}
}

func TestValidateCreateUserDomainRoleNil(t *testing.T) {
	if err := ValidateCreateUserDomainRole(nil); err == nil {
		t.Error("ValidateCreateUserDomainRole(nil) = nil, want error")
	}
}

func validValidateUpdateUserDomainRoleByIdRequest() *userdomainrolesv1.UpdateUserDomainRoleByIdRequest {
	return &userdomainrolesv1.UpdateUserDomainRoleByIdRequest{
		Id:                  100,
		UserId:              101,
		DomainRolesId:       102,
		DomainId:            103,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000005",
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
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "zero domain_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "empty updated_by_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.UpdatedByExternalId = "" }, wantField: "updated_by_id"},
		{name: "malformed updated_by_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) { r.UpdatedByExternalId = "not-a-uuid" }, wantField: "updated_by_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateUserDomainRoleByIdRequest()
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
		})
	}
}

func TestValidateUpdateUserDomainRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateUserDomainRoleById(nil); err == nil {
		t.Error("ValidateUpdateUserDomainRoleById(nil) = nil, want error")
	}
}
