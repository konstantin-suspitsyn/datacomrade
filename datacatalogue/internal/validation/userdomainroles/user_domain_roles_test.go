package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userDomainRolesFieldErrors достаёт из ошибки список полей с претензиями.
func userDomainRolesFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateUserDomainRolesRequest() *userdomainrolesv1.CreateUserDomainRolesRequest {
	return &userdomainrolesv1.CreateUserDomainRolesRequest{
		UserId:        100,
		DomainRolesId: 101,
		DomainId:      102,
		ExternalId:    "00000000-0000-4000-8000-000000000004",
	}
}

func TestValidateCreateUserDomainRoles(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateUserDomainRolesRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateUserDomainRolesRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRolesRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRolesRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "zero domain_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRolesRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRolesRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.CreateUserDomainRolesRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateUserDomainRolesRequest()
			tt.mutate(req)

			err := ValidateCreateUserDomainRoles(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateUserDomainRoles() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateUserDomainRoles() = nil, want error on %q", tt.wantField)
			}

			fields := userDomainRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateUserDomainRolesNil(t *testing.T) {
	if err := ValidateCreateUserDomainRoles(nil); err == nil {
		t.Error("ValidateCreateUserDomainRoles(nil) = nil, want error")
	}
}

func validValidateUpdateUserDomainRolesByIdRequest() *userdomainrolesv1.UpdateUserDomainRolesByIdRequest {
	return &userdomainrolesv1.UpdateUserDomainRolesByIdRequest{
		UserId:        100,
		DomainRolesId: 101,
		DomainId:      102,
		ExternalId:    "00000000-0000-4000-8000-000000000004",
		Id:            104,
	}
}

func TestValidateUpdateUserDomainRolesById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateUserDomainRolesByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateUserDomainRolesByIdRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero domain_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) { r.DomainRolesId = 0 }, wantField: "domain_roles_id"},
		{name: "zero domain_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateUserDomainRolesByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateUserDomainRolesById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateUserDomainRolesById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateUserDomainRolesById() = nil, want error on %q", tt.wantField)
			}

			fields := userDomainRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateUserDomainRolesByIdNil(t *testing.T) {
	if err := ValidateUpdateUserDomainRolesById(nil); err == nil {
		t.Error("ValidateUpdateUserDomainRolesById(nil) = nil, want error")
	}
}

func validValidateDeleteUserDomainRolesByIdRequest() *userdomainrolesv1.DeleteUserDomainRolesByIdRequest {
	return &userdomainrolesv1.DeleteUserDomainRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteUserDomainRolesById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.DeleteUserDomainRolesByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.DeleteUserDomainRolesByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.DeleteUserDomainRolesByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.DeleteUserDomainRolesByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *userdomainrolesv1.DeleteUserDomainRolesByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteUserDomainRolesByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteUserDomainRolesById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteUserDomainRolesById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteUserDomainRolesById() = nil, want error on %q", tt.wantField)
			}

			fields := userDomainRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteUserDomainRolesByIdNil(t *testing.T) {
	if err := ValidateDeleteUserDomainRolesById(nil); err == nil {
		t.Error("ValidateDeleteUserDomainRolesById(nil) = nil, want error")
	}
}

func validValidateUndeleteUserDomainRolesByIdRequest() *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest {
	return &userdomainrolesv1.UndeleteUserDomainRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteUserDomainRolesById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UndeleteUserDomainRolesByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UndeleteUserDomainRolesByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteUserDomainRolesByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteUserDomainRolesById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteUserDomainRolesById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteUserDomainRolesById() = nil, want error on %q", tt.wantField)
			}

			fields := userDomainRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteUserDomainRolesByIdNil(t *testing.T) {
	if err := ValidateUndeleteUserDomainRolesById(nil); err == nil {
		t.Error("ValidateUndeleteUserDomainRolesById(nil) = nil, want error")
	}
}

func TestValidateGetUserDomainRoleses(t *testing.T) {
	tests := []struct {
		name    string
		req     *userdomainrolesv1.GetUserDomainRolesesRequest
		wantErr bool
	}{
		{name: "valid", req: &userdomainrolesv1.GetUserDomainRolesesRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &userdomainrolesv1.GetUserDomainRolesesRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &userdomainrolesv1.GetUserDomainRolesesRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &userdomainrolesv1.GetUserDomainRolesesRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &userdomainrolesv1.GetUserDomainRolesesRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUserDomainRoleses(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUserDomainRoleses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetUserDomainRolesesNil(t *testing.T) {
	if err := ValidateGetUserDomainRoleses(nil); err == nil {
		t.Error("ValidateGetUserDomainRoleses(nil) = nil, want error")
	}
}
