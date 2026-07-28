package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateDatabaseCatRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateDatabaseCatRequest() *tablesv1.CreateDatabaseCatRequest {
	return &tablesv1.CreateDatabaseCatRequest{
		Name:           "name-0",
		HostId:         101,
		DatabaseTypeId: 102,
		Description:    "description-3",
		UserId:         104,
	}
}

func validUpdateDatabaseCatByIdRequest() *tablesv1.UpdateDatabaseCatByIdRequest {
	return &tablesv1.UpdateDatabaseCatByIdRequest{
		Id:             42,
		Name:           "name-0",
		HostId:         101,
		DatabaseTypeId: 102,
		Description:    "description-3",
		UserId:         104,
	}
}

// databaseCatFieldErrors достаёт из ошибки список полей с претензиями.
func databaseCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateDatabaseCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDatabaseCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDatabaseCatRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Name = strings.Repeat("a", databaseCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero host_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.HostId = 0 }, wantField: "host_id"},
		{name: "negative host_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.HostId = -1 }, wantField: "host_id"},
		{name: "zero database_type_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.DatabaseTypeId = 0 }, wantField: "database_type_id"},
		{name: "negative database_type_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.DatabaseTypeId = -1 }, wantField: "database_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateDatabaseCatRequest) {
			r.Description = strings.Repeat("a", databaseCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero user_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateDatabaseCatRequest()
			tt.mutate(req)

			err := ValidateCreateDatabaseCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDatabaseCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDatabaseCat() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCatFieldErrors(t, err)

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

func TestValidateUpdateDatabaseCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDatabaseCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDatabaseCatByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Name = strings.Repeat("a", databaseCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero host_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.HostId = 0 }, wantField: "host_id"},
		{name: "negative host_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.HostId = -1 }, wantField: "host_id"},
		{name: "zero database_type_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.DatabaseTypeId = 0 }, wantField: "database_type_id"},
		{name: "negative database_type_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.DatabaseTypeId = -1 }, wantField: "database_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) {
			r.Description = strings.Repeat("a", databaseCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero user_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateDatabaseCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDatabaseCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDatabaseCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDatabaseCatById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCatFieldErrors(t, err)

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
func TestValidateCreateDatabaseCatAtVarcharLimit(t *testing.T) {
	req := validCreateDatabaseCatRequest()
	req.Name = strings.Repeat("a", databaseCatNameMaxLen)

	if err := ValidateCreateDatabaseCat(req); err != nil {
		t.Errorf("ValidateCreateDatabaseCat() = %v, want nil at exactly %d chars", err, databaseCatNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateDatabaseCatCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateDatabaseCatRequest()
	req.Name = strings.Repeat("я", databaseCatNameMaxLen)

	if err := ValidateCreateDatabaseCat(req); err != nil {
		t.Errorf("ValidateCreateDatabaseCat() = %v, want nil at exactly %d cyrillic chars", err, databaseCatNameMaxLen)
	}
}

func TestValidateCreateDatabaseCatCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateDatabaseCat(&tablesv1.CreateDatabaseCatRequest{})

	if err == nil {
		t.Fatal("ValidateCreateDatabaseCat() = nil, want errors")
	}

	fields := databaseCatFieldErrors(t, err)

	wantFields := []string{"name", "host_id", "database_type_id", "description", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateDatabaseCatNil(t *testing.T) {
	if err := ValidateCreateDatabaseCat(nil); err == nil {
		t.Error("ValidateCreateDatabaseCat(nil) = nil, want error")
	}
}

func TestValidateUpdateDatabaseCatByIdNil(t *testing.T) {
	if err := ValidateUpdateDatabaseCatById(nil); err == nil {
		t.Error("ValidateUpdateDatabaseCatById(nil) = nil, want error")
	}
}
