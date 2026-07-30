package authlogic

import (
	"errors"
	"strings"
	"testing"

	authlogicv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/auth_logic/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// authLogicFieldErrors достаёт из ошибки список полей с претензиями.
func authLogicFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validGetTableIdsByExternalUserIdAndRolesRequest() *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest {
	return &authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Name:       "name-1",
	}
}

func TestValidateGetTableIdsByExternalUserIdAndRoles(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) {}},
		{name: "empty external_id", mutate: func(r *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) { r.ExternalId = "" }, wantField: "external_id"},
		{name: "malformed external_id", mutate: func(r *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) { r.ExternalId = "not-a-uuid" }, wantField: "external_id"},
		{name: "empty name", mutate: func(r *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) {
			r.Name = strings.Repeat("a", nameMaxLen+1)
		}, wantField: "name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validGetTableIdsByExternalUserIdAndRolesRequest()
			tt.mutate(req)

			err := ValidateGetTableIdsByExternalUserIdAndRoles(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateGetTableIdsByExternalUserIdAndRoles() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateGetTableIdsByExternalUserIdAndRoles() = nil, want error on %q", tt.wantField)
			}

			fields := authLogicFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateGetTableIdsByExternalUserIdAndRolesNil(t *testing.T) {
	if err := ValidateGetTableIdsByExternalUserIdAndRoles(nil); err == nil {
		t.Error("ValidateGetTableIdsByExternalUserIdAndRoles(nil) = nil, want error")
	}
}

func validGetTableIdsByUserIdAndRolesRequest() *authlogicv1.GetTableIdsByUserIdAndRolesRequest {
	return &authlogicv1.GetTableIdsByUserIdAndRolesRequest{
		UserId: 100,
		Name:   "name-1",
	}
}

func TestValidateGetTableIdsByUserIdAndRoles(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*authlogicv1.GetTableIdsByUserIdAndRolesRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*authlogicv1.GetTableIdsByUserIdAndRolesRequest) {}},
		{name: "zero user_id", mutate: func(r *authlogicv1.GetTableIdsByUserIdAndRolesRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "empty name", mutate: func(r *authlogicv1.GetTableIdsByUserIdAndRolesRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *authlogicv1.GetTableIdsByUserIdAndRolesRequest) { r.Name = strings.Repeat("a", nameMaxLen+1) }, wantField: "name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validGetTableIdsByUserIdAndRolesRequest()
			tt.mutate(req)

			err := ValidateGetTableIdsByUserIdAndRoles(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateGetTableIdsByUserIdAndRoles() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateGetTableIdsByUserIdAndRoles() = nil, want error on %q", tt.wantField)
			}

			fields := authLogicFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateGetTableIdsByUserIdAndRolesNil(t *testing.T) {
	if err := ValidateGetTableIdsByUserIdAndRoles(nil); err == nil {
		t.Error("ValidateGetTableIdsByUserIdAndRoles(nil) = nil, want error")
	}
}
