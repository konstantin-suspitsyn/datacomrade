package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateColumnTypeRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateColumnTypeRequest() *tablesv1.CreateColumnTypeRequest {
	return &tablesv1.CreateColumnTypeRequest{
		Name:        "name-0",
		Description: "description-1",
		UserId:      102,
	}
}

func validUpdateColumnTypeByIdRequest() *tablesv1.UpdateColumnTypeByIdRequest {
	return &tablesv1.UpdateColumnTypeByIdRequest{
		Id:          42,
		Name:        "name-0",
		Description: "description-1",
		UserId:      102,
	}
}

// columnTypeFieldErrors достаёт из ошибки список полей с претензиями.
func columnTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateColumnType(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateColumnTypeRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateColumnTypeRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Name = strings.Repeat("a", columnTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateColumnTypeRequest) {
			r.Description = strings.Repeat("a", columnTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero user_id", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateColumnTypeRequest()
			tt.mutate(req)

			err := ValidateCreateColumnType(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateColumnType() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateColumnType() = nil, want error on %q", tt.wantField)
			}

			fields := columnTypeFieldErrors(t, err)

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

func TestValidateUpdateColumnTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateColumnTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateColumnTypeByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Name = strings.Repeat("a", columnTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) {
			r.Description = strings.Repeat("a", columnTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero user_id", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateColumnTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateColumnTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateColumnTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateColumnTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := columnTypeFieldErrors(t, err)

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
func TestValidateCreateColumnTypeAtVarcharLimit(t *testing.T) {
	req := validCreateColumnTypeRequest()
	req.Name = strings.Repeat("a", columnTypeNameMaxLen)

	if err := ValidateCreateColumnType(req); err != nil {
		t.Errorf("ValidateCreateColumnType() = %v, want nil at exactly %d chars", err, columnTypeNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateColumnTypeCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateColumnTypeRequest()
	req.Name = strings.Repeat("я", columnTypeNameMaxLen)

	if err := ValidateCreateColumnType(req); err != nil {
		t.Errorf("ValidateCreateColumnType() = %v, want nil at exactly %d cyrillic chars", err, columnTypeNameMaxLen)
	}
}

func TestValidateCreateColumnTypeCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateColumnType(&tablesv1.CreateColumnTypeRequest{})

	if err == nil {
		t.Fatal("ValidateCreateColumnType() = nil, want errors")
	}

	fields := columnTypeFieldErrors(t, err)

	wantFields := []string{"name", "description", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateColumnTypeNil(t *testing.T) {
	if err := ValidateCreateColumnType(nil); err == nil {
		t.Error("ValidateCreateColumnType(nil) = nil, want error")
	}
}

func TestValidateUpdateColumnTypeByIdNil(t *testing.T) {
	if err := ValidateUpdateColumnTypeById(nil); err == nil {
		t.Error("ValidateUpdateColumnTypeById(nil) = nil, want error")
	}
}
