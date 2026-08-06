package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userFieldErrors достаёт из ошибки список полей с претензиями.
func userFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateUserRequest() *tablesv1.CreateUserRequest {
	return &tablesv1.CreateUserRequest{
		Name:       "name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
	}
}

func TestValidateCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateUserRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateUserRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateUserRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateUserRequest) { r.Name = strings.Repeat("a", userNameMaxLen+1) }, wantField: "name"},
		{name: "empty external_id", mutate: func(r *tablesv1.CreateUserRequest) { r.ExternalId = "" }, wantField: "external_id"},
		{name: "malformed external_id", mutate: func(r *tablesv1.CreateUserRequest) { r.ExternalId = "not-a-uuid" }, wantField: "external_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateUserRequest()
			tt.mutate(req)

			err := ValidateCreateUser(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateUser() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateUser() = nil, want error on %q", tt.wantField)
			}

			fields := userFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateUserNil(t *testing.T) {
	if err := ValidateCreateUser(nil); err == nil {
		t.Error("ValidateCreateUser(nil) = nil, want error")
	}
}

func TestValidateGetUserByExternalId(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "00000000-0000-4000-8000-000000007001"},
		{name: "empty", value: "", wantErr: true},
		{name: "malformed", value: "not-a-uuid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUserByExternalId(&tablesv1.GetUserByExternalIdRequest{ExternalId: tt.value})

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUserByExternalId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetUserByExternalIdNil(t *testing.T) {
	if err := ValidateGetUserByExternalId(nil); err == nil {
		t.Error("ValidateGetUserByExternalId(nil) = nil, want error")
	}
}

func TestValidateGetUsers(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetUsersRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetUsersRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetUsersRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetUsersRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetUsersRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetUsersRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUsers(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetUsersNil(t *testing.T) {
	if err := ValidateGetUsers(nil); err == nil {
		t.Error("ValidateGetUsers(nil) = nil, want error")
	}
}
