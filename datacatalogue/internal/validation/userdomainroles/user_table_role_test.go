package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userTableRoleFieldErrors достаёт из ошибки список полей с претензиями.
func userTableRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateUserTableRoleRequest() *userdomainrolesv1.CreateUserTableRoleRequest {
	return &userdomainrolesv1.CreateUserTableRoleRequest{
		UserId:              100,
		TableRolesId:        101,
		TableId:             102,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000004",
	}
}

func TestValidateCreateUserTableRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateUserTableRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateUserTableRoleRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "zero table_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "empty updated_by_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.UpdatedByExternalId = "" }, wantField: "updated_by_id"},
		{name: "malformed updated_by_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.UpdatedByExternalId = "not-a-uuid" }, wantField: "updated_by_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateUserTableRoleRequest()
			tt.mutate(req)

			err := ValidateCreateUserTableRole(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateUserTableRole() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateUserTableRole() = nil, want error on %q", tt.wantField)
			}

			fields := userTableRoleFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateUserTableRoleNil(t *testing.T) {
	if err := ValidateCreateUserTableRole(nil); err == nil {
		t.Error("ValidateCreateUserTableRole(nil) = nil, want error")
	}
}

func validValidateUpdateUserTableRoleByIdRequest() *userdomainrolesv1.UpdateUserTableRoleByIdRequest {
	return &userdomainrolesv1.UpdateUserTableRoleByIdRequest{
		Id:                  100,
		UserId:              101,
		TableRolesId:        102,
		TableId:             103,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000005",
	}
}

func TestValidateUpdateUserTableRoleById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateUserTableRoleByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateUserTableRoleByIdRequest) {}},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "zero table_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "empty updated_by_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.UpdatedByExternalId = "" }, wantField: "updated_by_id"},
		{name: "malformed updated_by_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.UpdatedByExternalId = "not-a-uuid" }, wantField: "updated_by_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateUserTableRoleByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateUserTableRoleById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateUserTableRoleById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateUserTableRoleById() = nil, want error on %q", tt.wantField)
			}

			fields := userTableRoleFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateUserTableRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateUserTableRoleById(nil); err == nil {
		t.Error("ValidateUpdateUserTableRoleById(nil) = nil, want error")
	}
}
