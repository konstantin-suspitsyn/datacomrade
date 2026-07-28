package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateCalculationTypeRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateCalculationTypeRequest() *tablesv1.CreateCalculationTypeRequest {
	return &tablesv1.CreateCalculationTypeRequest{
		Name:        "name-0",
		Description: "description-1",
	}
}

func validUpdateCalculationTypeByIdRequest() *tablesv1.UpdateCalculationTypeByIdRequest {
	return &tablesv1.UpdateCalculationTypeByIdRequest{
		Id:          42,
		Name:        "name-0",
		Description: "description-1",
	}
}

// calculationTypeFieldErrors достаёт из ошибки список полей с претензиями.
func calculationTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateCalculationType(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateCalculationTypeRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateCalculationTypeRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateCalculationTypeRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateCalculationTypeRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateCalculationTypeRequest) {
			r.Name = strings.Repeat("a", calculationTypeNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateCalculationTypeRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateCalculationTypeRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateCalculationTypeRequest) {
			r.Description = strings.Repeat("a", calculationTypeDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateCalculationTypeRequest()
			tt.mutate(req)

			err := ValidateCreateCalculationType(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateCalculationType() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateCalculationType() = nil, want error on %q", tt.wantField)
			}

			fields := calculationTypeFieldErrors(t, err)

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

func TestValidateUpdateCalculationTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateCalculationTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateCalculationTypeByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) {
			r.Name = strings.Repeat("a", calculationTypeNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateCalculationTypeByIdRequest) {
			r.Description = strings.Repeat("a", calculationTypeDescriptionMaxLen+1)
		}, wantField: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateCalculationTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateCalculationTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateCalculationTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateCalculationTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := calculationTypeFieldErrors(t, err)

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
func TestValidateCreateCalculationTypeAtVarcharLimit(t *testing.T) {
	req := validCreateCalculationTypeRequest()
	req.Name = strings.Repeat("a", calculationTypeNameMaxLen)

	if err := ValidateCreateCalculationType(req); err != nil {
		t.Errorf("ValidateCreateCalculationType() = %v, want nil at exactly %d chars", err, calculationTypeNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateCalculationTypeCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateCalculationTypeRequest()
	req.Name = strings.Repeat("я", calculationTypeNameMaxLen)

	if err := ValidateCreateCalculationType(req); err != nil {
		t.Errorf("ValidateCreateCalculationType() = %v, want nil at exactly %d cyrillic chars", err, calculationTypeNameMaxLen)
	}
}

func TestValidateCreateCalculationTypeCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateCalculationType(&tablesv1.CreateCalculationTypeRequest{})

	if err == nil {
		t.Fatal("ValidateCreateCalculationType() = nil, want errors")
	}

	fields := calculationTypeFieldErrors(t, err)

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

func TestValidateCreateCalculationTypeNil(t *testing.T) {
	if err := ValidateCreateCalculationType(nil); err == nil {
		t.Error("ValidateCreateCalculationType(nil) = nil, want error")
	}
}

func TestValidateUpdateCalculationTypeByIdNil(t *testing.T) {
	if err := ValidateUpdateCalculationTypeById(nil); err == nil {
		t.Error("ValidateUpdateCalculationTypeById(nil) = nil, want error")
	}
}
