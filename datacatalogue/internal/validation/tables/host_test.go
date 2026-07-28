package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateHostRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateHostRequest() *tablesv1.CreateHostRequest {
	return &tablesv1.CreateHostRequest{
		Name:        "name-0",
		Description: "description-1",
		HostEnv:     "host-env-2",
		PortEnv:     "port-env-3",
		UsernameEnv: "username-env-4",
		PasswordEnv: "password-env-5",
		UserId:      106,
	}
}

func validUpdateHostByIdRequest() *tablesv1.UpdateHostByIdRequest {
	return &tablesv1.UpdateHostByIdRequest{
		Id:          42,
		Name:        "name-0",
		Description: "description-1",
		HostEnv:     "host-env-2",
		PortEnv:     "port-env-3",
		UsernameEnv: "username-env-4",
		PasswordEnv: "password-env-5",
		UserId:      106,
	}
}

// hostFieldErrors достаёт из ошибки список полей с претензиями.
func hostFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateHost(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateHostRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateHostRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateHostRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateHostRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateHostRequest) { r.Name = strings.Repeat("a", hostNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateHostRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateHostRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateHostRequest) { r.Description = strings.Repeat("a", hostDescriptionMaxLen+1) }, wantField: "description"},
		{name: "empty host_env", mutate: func(r *tablesv1.CreateHostRequest) { r.HostEnv = "" }, wantField: "host_env"},
		{name: "blank host_env", mutate: func(r *tablesv1.CreateHostRequest) { r.HostEnv = "   " }, wantField: "host_env"},
		{name: "host_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.HostEnv = strings.Repeat("a", hostHostEnvMaxLen+1) }, wantField: "host_env"},
		{name: "empty port_env", mutate: func(r *tablesv1.CreateHostRequest) { r.PortEnv = "" }, wantField: "port_env"},
		{name: "blank port_env", mutate: func(r *tablesv1.CreateHostRequest) { r.PortEnv = "   " }, wantField: "port_env"},
		{name: "port_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.PortEnv = strings.Repeat("a", hostPortEnvMaxLen+1) }, wantField: "port_env"},
		{name: "empty username_env", mutate: func(r *tablesv1.CreateHostRequest) { r.UsernameEnv = "" }, wantField: "username_env"},
		{name: "blank username_env", mutate: func(r *tablesv1.CreateHostRequest) { r.UsernameEnv = "   " }, wantField: "username_env"},
		{name: "username_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.UsernameEnv = strings.Repeat("a", hostUsernameEnvMaxLen+1) }, wantField: "username_env"},
		{name: "empty password_env", mutate: func(r *tablesv1.CreateHostRequest) { r.PasswordEnv = "" }, wantField: "password_env"},
		{name: "blank password_env", mutate: func(r *tablesv1.CreateHostRequest) { r.PasswordEnv = "   " }, wantField: "password_env"},
		{name: "password_env too long", mutate: func(r *tablesv1.CreateHostRequest) { r.PasswordEnv = strings.Repeat("a", hostPasswordEnvMaxLen+1) }, wantField: "password_env"},
		{name: "zero user_id", mutate: func(r *tablesv1.CreateHostRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.CreateHostRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateHostRequest()
			tt.mutate(req)

			err := ValidateCreateHost(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateHost() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateHost() = nil, want error on %q", tt.wantField)
			}

			fields := hostFieldErrors(t, err)

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

func TestValidateUpdateHostById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateHostByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateHostByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Name = strings.Repeat("a", hostNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.Description = strings.Repeat("a", hostDescriptionMaxLen+1) }, wantField: "description"},
		{name: "empty host_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.HostEnv = "" }, wantField: "host_env"},
		{name: "blank host_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.HostEnv = "   " }, wantField: "host_env"},
		{name: "host_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.HostEnv = strings.Repeat("a", hostHostEnvMaxLen+1) }, wantField: "host_env"},
		{name: "empty port_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PortEnv = "" }, wantField: "port_env"},
		{name: "blank port_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PortEnv = "   " }, wantField: "port_env"},
		{name: "port_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PortEnv = strings.Repeat("a", hostPortEnvMaxLen+1) }, wantField: "port_env"},
		{name: "empty username_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.UsernameEnv = "" }, wantField: "username_env"},
		{name: "blank username_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.UsernameEnv = "   " }, wantField: "username_env"},
		{name: "username_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.UsernameEnv = strings.Repeat("a", hostUsernameEnvMaxLen+1) }, wantField: "username_env"},
		{name: "empty password_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PasswordEnv = "" }, wantField: "password_env"},
		{name: "blank password_env", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PasswordEnv = "   " }, wantField: "password_env"},
		{name: "password_env too long", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.PasswordEnv = strings.Repeat("a", hostPasswordEnvMaxLen+1) }, wantField: "password_env"},
		{name: "zero user_id", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.UpdateHostByIdRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateHostByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateHostById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateHostById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateHostById() = nil, want error on %q", tt.wantField)
			}

			fields := hostFieldErrors(t, err)

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
func TestValidateCreateHostAtVarcharLimit(t *testing.T) {
	req := validCreateHostRequest()
	req.Name = strings.Repeat("a", hostNameMaxLen)

	if err := ValidateCreateHost(req); err != nil {
		t.Errorf("ValidateCreateHost() = %v, want nil at exactly %d chars", err, hostNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateHostCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateHostRequest()
	req.Name = strings.Repeat("я", hostNameMaxLen)

	if err := ValidateCreateHost(req); err != nil {
		t.Errorf("ValidateCreateHost() = %v, want nil at exactly %d cyrillic chars", err, hostNameMaxLen)
	}
}

func TestValidateCreateHostCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateHost(&tablesv1.CreateHostRequest{})

	if err == nil {
		t.Fatal("ValidateCreateHost() = nil, want errors")
	}

	fields := hostFieldErrors(t, err)

	wantFields := []string{"name", "description", "host_env", "port_env", "username_env", "password_env", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateHostNil(t *testing.T) {
	if err := ValidateCreateHost(nil); err == nil {
		t.Error("ValidateCreateHost(nil) = nil, want error")
	}
}

func TestValidateUpdateHostByIdNil(t *testing.T) {
	if err := ValidateUpdateHostById(nil); err == nil {
		t.Error("ValidateUpdateHostById(nil) = nil, want error")
	}
}
