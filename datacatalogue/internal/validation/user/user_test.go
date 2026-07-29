package user

import (
	"errors"
	"strings"
	"testing"

	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateUserRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateUserRequest() *userv1.CreateUserRequest {
	return &userv1.CreateUserRequest{
		Name:       "name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
	}
}

func validUpdateUserByIdRequest() *userv1.UpdateUserByIdRequest {
	return &userv1.UpdateUserByIdRequest{
		Id:         42,
		Name:       "name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
	}
}

// userFieldErrors достаёт из ошибки список полей с претензиями.
func userFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userv1.CreateUserRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userv1.CreateUserRequest) {}},
		{name: "empty name", mutate: func(r *userv1.CreateUserRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *userv1.CreateUserRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *userv1.CreateUserRequest) { r.Name = strings.Repeat("a", userNameMaxLen+1) }, wantField: "name"},
		{name: "empty external_id", mutate: func(r *userv1.CreateUserRequest) { r.ExternalId = "" }, wantField: "external_id"},
		{name: "malformed external_id", mutate: func(r *userv1.CreateUserRequest) { r.ExternalId = "not-a-uuid" }, wantField: "external_id"},
		{name: "external_id without dashes", mutate: func(r *userv1.CreateUserRequest) { r.ExternalId = "00000000000040008000000000000001" }, wantField: "external_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateUserRequest()
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

			// Порча одного поля не должна задевать остальные.
			if len(fields) != 1 {
				t.Errorf("errors on %d fields, want only %q: %v", len(fields), tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateUserById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*userv1.UpdateUserByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*userv1.UpdateUserByIdRequest) {}},
		{name: "zero id", mutate: func(r *userv1.UpdateUserByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *userv1.UpdateUserByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *userv1.UpdateUserByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *userv1.UpdateUserByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *userv1.UpdateUserByIdRequest) { r.Name = strings.Repeat("a", userNameMaxLen+1) }, wantField: "name"},
		{name: "empty external_id", mutate: func(r *userv1.UpdateUserByIdRequest) { r.ExternalId = "" }, wantField: "external_id"},
		{name: "malformed external_id", mutate: func(r *userv1.UpdateUserByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "external_id"},
		{name: "external_id without dashes", mutate: func(r *userv1.UpdateUserByIdRequest) { r.ExternalId = "00000000000040008000000000000001" }, wantField: "external_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateUserByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateUserById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateUserById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateUserById() = nil, want error on %q", tt.wantField)
			}

			fields := userFieldErrors(t, err)

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
func TestValidateCreateUserAtVarcharLimit(t *testing.T) {
	req := validCreateUserRequest()
	req.Name = strings.Repeat("a", userNameMaxLen)

	if err := ValidateCreateUser(req); err != nil {
		t.Errorf("ValidateCreateUser() = %v, want nil at exactly %d chars", err, userNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateUserCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateUserRequest()
	req.Name = strings.Repeat("я", userNameMaxLen)

	if err := ValidateCreateUser(req); err != nil {
		t.Errorf("ValidateCreateUser() = %v, want nil at exactly %d cyrillic chars", err, userNameMaxLen)
	}
}

func TestValidateCreateUserCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateUser(&userv1.CreateUserRequest{})

	if err == nil {
		t.Fatal("ValidateCreateUser() = nil, want errors")
	}

	fields := userFieldErrors(t, err)

	wantFields := []string{"name", "external_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateUserNil(t *testing.T) {
	if err := ValidateCreateUser(nil); err == nil {
		t.Error("ValidateCreateUser(nil) = nil, want error")
	}
}

func TestValidateUpdateUserByIdNil(t *testing.T) {
	if err := ValidateUpdateUserById(nil); err == nil {
		t.Error("ValidateUpdateUserById(nil) = nil, want error")
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
		{name: "without dashes", value: "00000000000040008000000000000001", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetUserByExternalId(&userv1.GetUserByExternalIdRequest{ExternalId: tt.value})

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetUserByExternalId() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				return
			}

			fields := userFieldErrors(t, err)

			if len(fields["external_id"]) == 0 {
				t.Errorf("no error on external_id, got %v", fields)
			}
		})
	}
}

func TestValidateGetUserByExternalIdNil(t *testing.T) {
	if err := ValidateGetUserByExternalId(nil); err == nil {
		t.Error("ValidateGetUserByExternalId(nil) = nil, want error")
	}
}
