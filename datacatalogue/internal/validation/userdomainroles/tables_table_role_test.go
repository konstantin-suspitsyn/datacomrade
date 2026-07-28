package userdomainroles

import (
	"errors"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateTablesTableRoleRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateTablesTableRoleRequest() *userdomainrolesv1.CreateTablesTableRoleRequest {
	return &userdomainrolesv1.CreateTablesTableRoleRequest{
		TableCatId:   100,
		TableRolesId: 101,
	}
}

func validUpdateTablesTableRoleByIdRequest() *userdomainrolesv1.UpdateTablesTableRoleByIdRequest {
	return &userdomainrolesv1.UpdateTablesTableRoleByIdRequest{
		Id:           42,
		TableCatId:   100,
		TableRolesId: 101,
	}
}

// tablesTableRoleFieldErrors достаёт из ошибки список полей с претензиями.
func tablesTableRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateTablesTableRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateTablesTableRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateTablesTableRoleRequest) {}},
		{name: "zero table_cat_id", mutate: func(r *userdomainrolesv1.CreateTablesTableRoleRequest) { r.TableCatId = 0 }, wantField: "table_cat_id"},
		{name: "negative table_cat_id", mutate: func(r *userdomainrolesv1.CreateTablesTableRoleRequest) { r.TableCatId = -1 }, wantField: "table_cat_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.CreateTablesTableRoleRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "negative table_roles_id", mutate: func(r *userdomainrolesv1.CreateTablesTableRoleRequest) { r.TableRolesId = -1 }, wantField: "table_roles_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateTablesTableRoleRequest()
			tt.mutate(req)

			err := ValidateCreateTablesTableRole(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateTablesTableRole() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateTablesTableRole() = nil, want error on %q", tt.wantField)
			}

			fields := tablesTableRoleFieldErrors(t, err)

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

func TestValidateUpdateTablesTableRoleById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.UpdateTablesTableRoleByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.UpdateTablesTableRoleByIdRequest) {}},
		{name: "zero id", mutate: func(r *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero table_cat_id", mutate: func(r *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) { r.TableCatId = 0 }, wantField: "table_cat_id"},
		{name: "negative table_cat_id", mutate: func(r *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) { r.TableCatId = -1 }, wantField: "table_cat_id"},
		{name: "zero table_roles_id", mutate: func(r *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) { r.TableRolesId = 0 }, wantField: "table_roles_id"},
		{name: "negative table_roles_id", mutate: func(r *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) { r.TableRolesId = -1 }, wantField: "table_roles_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateTablesTableRoleByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateTablesTableRoleById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateTablesTableRoleById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateTablesTableRoleById() = nil, want error on %q", tt.wantField)
			}

			fields := tablesTableRoleFieldErrors(t, err)

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

func TestValidateCreateTablesTableRoleCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateTablesTableRole(&userdomainrolesv1.CreateTablesTableRoleRequest{})

	if err == nil {
		t.Fatal("ValidateCreateTablesTableRole() = nil, want errors")
	}

	fields := tablesTableRoleFieldErrors(t, err)

	wantFields := []string{"table_cat_id", "table_roles_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateTablesTableRoleNil(t *testing.T) {
	if err := ValidateCreateTablesTableRole(nil); err == nil {
		t.Error("ValidateCreateTablesTableRole(nil) = nil, want error")
	}
}

func TestValidateUpdateTablesTableRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateTablesTableRoleById(nil); err == nil {
		t.Error("ValidateUpdateTablesTableRoleById(nil) = nil, want error")
	}
}
