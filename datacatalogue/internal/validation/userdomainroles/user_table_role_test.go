package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateUserTableRoleRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateUserTableRoleRequest() *userdomainrolesv1.CreateUserTableRoleRequest {
	return &userdomainrolesv1.CreateUserTableRoleRequest{
		UserId:       100,
		TableRolesId: 101,
		TableId:      102,
	}
}

func validUpdateUserTableRoleByIdRequest() *userdomainrolesv1.UpdateUserTableRoleByIdRequest {
	return &userdomainrolesv1.UpdateUserTableRoleByIdRequest{
		Id:           42,
		UserId:       100,
		TableRolesId: 101,
		TableId:      102,
	}
}

// userTableRoleFieldErrors достаёт из ошибки список полей с претензиями.
func userTableRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateUserTableRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateUserTableRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateUserTableRoleRequest) {}},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.UserId = -1 }, wantField: "user_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "negative table_roles_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.TableRolesId = -1 }, wantField: "table_roles_id"},
		{name: "zero table_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "negative table_id", mutate: func(r *userdomainrolesv1.CreateUserTableRoleRequest) { r.TableId = -1 }, wantField: "table_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateUserTableRoleRequest()
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

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
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
		{name: "negative id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero user_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.UserId = -1 }, wantField: "user_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "negative table_roles_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.TableRolesId = -1 }, wantField: "table_roles_id"},
		{name: "zero table_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "negative table_id", mutate: func(r *userdomainrolesv1.UpdateUserTableRoleByIdRequest) { r.TableId = -1 }, wantField: "table_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateUserTableRoleByIdRequest()
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

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateUserTableRoleCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateUserTableRole(&userdomainrolesv1.CreateUserTableRoleRequest{})

	if err == nil {
		t.Fatal("ValidateCreateUserTableRole() = nil, want errors")
	}

	fields := userTableRoleFieldErrors(t, err)

	wantFields := []string{"user_id", "table_roles_id", "table_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateUserTableRoleNil(t *testing.T) {
	if err := ValidateCreateUserTableRole(nil); err == nil {
		t.Error("ValidateCreateUserTableRole(nil) = nil, want error")
	}
}

func TestValidateUpdateUserTableRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateUserTableRoleById(nil); err == nil {
		t.Error("ValidateUpdateUserTableRoleById(nil) = nil, want error")
	}
}
