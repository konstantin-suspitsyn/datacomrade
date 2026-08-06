package userdomainroles

import (
	"errors"
	"strings"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tableRoleFieldErrors достаёт из ошибки список полей с претензиями.
func tableRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateTableRoleRequest() *userdomainrolesv1.CreateTableRoleRequest {
	return &userdomainrolesv1.CreateTableRoleRequest{
		Name:        "name-0",
		Description: "description-1",
	}
}

func TestValidateCreateTableRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateTableRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateTableRoleRequest) {}},
		{name: "empty name", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Name = strings.Repeat("a", tableRoleNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) {
			r.Description = strings.Repeat("a", tableRoleDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateTableRoleRequest()
			tt.mutate(req)

			err := ValidateCreateTableRole(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateTableRole() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateTableRole() = nil, want error on %q", tt.wantField)
			}

			fields := tableRoleFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateTableRoleNil(t *testing.T) {
	if err := ValidateCreateTableRole(nil); err == nil {
		t.Error("ValidateCreateTableRole(nil) = nil, want error")
	}
}

func validValidateUpdateTableRoleByIdRequest() *userdomainrolesv1.UpdateTableRoleByIdRequest {
	return &userdomainrolesv1.UpdateTableRoleByIdRequest{
		Id:          100,
		Name:        "name-1",
		Description: "description-2",
	}
}

func TestValidateUpdateTableRoleById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateTableRoleByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateTableRoleByIdRequest) {}},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "empty name", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) {
			r.Name = strings.Repeat("a", tableRoleNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty description", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) {
			r.Description = strings.Repeat("a", tableRoleDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateTableRoleByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateTableRoleById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateTableRoleById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateTableRoleById() = nil, want error on %q", tt.wantField)
			}

			fields := tableRoleFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateTableRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateTableRoleById(nil); err == nil {
		t.Error("ValidateUpdateTableRoleById(nil) = nil, want error")
	}
}
