package userdomainroles

import (
	"errors"
	"strings"
	"testing"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateTableRoleRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateTableRoleRequest() *userdomainrolesv1.CreateTableRoleRequest {
	return &userdomainrolesv1.CreateTableRoleRequest{
		Name:        "name-0",
		Description: "description-1",
	}
}

func validUpdateTableRoleByIdRequest() *userdomainrolesv1.UpdateTableRoleByIdRequest {
	return &userdomainrolesv1.UpdateTableRoleByIdRequest{
		Id:          42,
		Name:        "name-0",
		Description: "description-1",
	}
}

// tableRoleFieldErrors достаёт из ошибки список полей с претензиями.
func tableRoleFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateTableRole(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userdomainrolesv1.CreateTableRoleRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userdomainrolesv1.CreateTableRoleRequest) {}},
		{name: "empty name", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Name = strings.Repeat("a", tableRoleNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *userdomainrolesv1.CreateTableRoleRequest) {
			r.Description = strings.Repeat("a", tableRoleDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateTableRoleRequest()
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

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
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
		{name: "negative id", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) {
			r.Name = strings.Repeat("a", tableRoleNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty description", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *userdomainrolesv1.UpdateTableRoleByIdRequest) {
			r.Description = strings.Repeat("a", tableRoleDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateTableRoleByIdRequest()
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

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
	}
}

// Ровно граничная длина проходит: varchar(n) допускает n символов.
func TestValidateCreateTableRoleAtVarcharLimit(t *testing.T) {
	req := validCreateTableRoleRequest()
	req.Name = strings.Repeat("a", tableRoleNameMaxLen)

	if err := ValidateCreateTableRole(req); err != nil {
		t.Errorf("ValidateCreateTableRole() = %v, want nil at exactly %d chars", err, tableRoleNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateTableRoleCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateTableRoleRequest()
	req.Name = strings.Repeat("я", tableRoleNameMaxLen)

	if err := ValidateCreateTableRole(req); err != nil {
		t.Errorf("ValidateCreateTableRole() = %v, want nil at exactly %d cyrillic chars", err, tableRoleNameMaxLen)
	}
}

func TestValidateCreateTableRoleCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateTableRole(&userdomainrolesv1.CreateTableRoleRequest{})

	if err == nil {
		t.Fatal("ValidateCreateTableRole() = nil, want errors")
	}

	fields := tableRoleFieldErrors(t, err)

	wantFields := []string{"name", "description"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateTableRoleNil(t *testing.T) {
	if err := ValidateCreateTableRole(nil); err == nil {
		t.Error("ValidateCreateTableRole(nil) = nil, want error")
	}
}

func TestValidateUpdateTableRoleByIdNil(t *testing.T) {
	if err := ValidateUpdateTableRoleById(nil); err == nil {
		t.Error("ValidateUpdateTableRoleById(nil) = nil, want error")
	}
}
