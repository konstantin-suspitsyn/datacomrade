package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userTableRolesFieldErrors достаёт из ошибки список полей с претензиями.
func userTableRolesFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateUserTableRolesRequest() *userdomainrolesv1.CreateUserTableRolesRequest {
	return &userdomainrolesv1.CreateUserTableRolesRequest{
		UserId:       100,
		TableRolesId: 101,
		TableId:      102,
		ExternalId:   "00000000-0000-4000-8000-000000000004",
	}
}

func TestValidateCreateUserTableRoles(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateUserTableRolesRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateUserTableRolesRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.CreateUserTableRolesRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.CreateUserTableRolesRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "zero table_id", mutate: func(r *userdomainrolesv1.CreateUserTableRolesRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.CreateUserTableRolesRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.CreateUserTableRolesRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateUserTableRolesRequest()
			tt.mutate(req)

			err := ValidateCreateUserTableRoles(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateUserTableRoles() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateUserTableRoles() = nil, want error on %q", tt.wantField)
			}

			fields := userTableRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateUserTableRolesNil(t *testing.T) {
	if err := ValidateCreateUserTableRoles(nil); err == nil {
		t.Error("ValidateCreateUserTableRoles(nil) = nil, want error")
	}
}

func validValidateUpdateUserTableRolesByIdRequest() *userdomainrolesv1.UpdateUserTableRolesByIdRequest {
	return &userdomainrolesv1.UpdateUserTableRolesByIdRequest{
		UserId:       100,
		TableRolesId: 101,
		TableId:      102,
		ExternalId:   "00000000-0000-4000-8000-000000000004",
		Id:           104,
	}
}

func TestValidateUpdateUserTableRolesById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateUserTableRolesByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateUserTableRolesByIdRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRolesByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRolesByIdRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "zero table_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRolesByIdRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRolesByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRolesByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateUserTableRolesByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateUserTableRolesByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateUserTableRolesById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateUserTableRolesById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateUserTableRolesById() = nil, want error on %q", tt.wantField)
			}

			fields := userTableRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateUserTableRolesByIdNil(t *testing.T) {
	if err := ValidateUpdateUserTableRolesById(nil); err == nil {
		t.Error("ValidateUpdateUserTableRolesById(nil) = nil, want error")
	}
}

func validValidateDeleteUserTableRolesByIdRequest() *userdomainrolesv1.DeleteUserTableRolesByIdRequest {
	return &userdomainrolesv1.DeleteUserTableRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteUserTableRolesById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.DeleteUserTableRolesByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.DeleteUserTableRolesByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.DeleteUserTableRolesByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.DeleteUserTableRolesByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *userdomainrolesv1.DeleteUserTableRolesByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteUserTableRolesByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteUserTableRolesById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteUserTableRolesById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteUserTableRolesById() = nil, want error on %q", tt.wantField)
			}

			fields := userTableRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteUserTableRolesByIdNil(t *testing.T) {
	if err := ValidateDeleteUserTableRolesById(nil); err == nil {
		t.Error("ValidateDeleteUserTableRolesById(nil) = nil, want error")
	}
}

func validValidateUndeleteUserTableRolesByIdRequest() *userdomainrolesv1.UndeleteUserTableRolesByIdRequest {
	return &userdomainrolesv1.UndeleteUserTableRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteUserTableRolesById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UndeleteUserTableRolesByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UndeleteUserTableRolesByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *userdomainrolesv1.UndeleteUserTableRolesByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *userdomainrolesv1.UndeleteUserTableRolesByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UndeleteUserTableRolesByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteUserTableRolesByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteUserTableRolesById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteUserTableRolesById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteUserTableRolesById() = nil, want error on %q", tt.wantField)
			}

			fields := userTableRolesFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteUserTableRolesByIdNil(t *testing.T) {
	if err := ValidateUndeleteUserTableRolesById(nil); err == nil {
		t.Error("ValidateUndeleteUserTableRolesById(nil) = nil, want error")
	}
}

func TestValidateGetUserTableRoleses(t *testing.T) {
	tests := []struct {
		name    string
		req     *userdomainrolesv1.GetUserTableRolesesRequest
		wantErr bool
	}{
		{name: "valid", req: &userdomainrolesv1.GetUserTableRolesesRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &userdomainrolesv1.GetUserTableRolesesRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &userdomainrolesv1.GetUserTableRolesesRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &userdomainrolesv1.GetUserTableRolesesRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &userdomainrolesv1.GetUserTableRolesesRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUserTableRoleses(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUserTableRoleses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetUserTableRolesesNil(t *testing.T) {
	if err := ValidateGetUserTableRoleses(nil); err == nil {
		t.Error("ValidateGetUserTableRoleses(nil) = nil, want error")
	}
}

func TestValidateGetUserTableRolesesFiltered(t *testing.T) {
	tests := []struct {
		name    string
		req     *userdomainrolesv1.GetUserTableRolesesFilteredRequest
		wantErr bool
	}{
		{name: "valid", req: &userdomainrolesv1.GetUserTableRolesesFilteredRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &userdomainrolesv1.GetUserTableRolesesFilteredRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &userdomainrolesv1.GetUserTableRolesesFilteredRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &userdomainrolesv1.GetUserTableRolesesFilteredRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &userdomainrolesv1.GetUserTableRolesesFilteredRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUserTableRolesesFiltered(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUserTableRolesesFiltered() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetUserTableRolesesFilteredNil(t *testing.T) {
	if err := ValidateGetUserTableRolesesFiltered(nil); err == nil {
		t.Error("ValidateGetUserTableRolesesFiltered(nil) = nil, want error")
	}
}
