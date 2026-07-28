package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateColumnCatRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateColumnCatRequest() *tablesv1.CreateColumnCatRequest {
	return &tablesv1.CreateColumnCatRequest{
		TableId:           100,
		Name:              "name-1",
		AliasId:           102,
		ColumnTypeId:      103,
		Description:       "description-4",
		CalculationTypeId: 105,
		ShowInUi:          true,
		UserId:            107,
	}
}

func validUpdateColumnCatByIdRequest() *tablesv1.UpdateColumnCatByIdRequest {
	return &tablesv1.UpdateColumnCatByIdRequest{
		Id:                42,
		TableId:           100,
		Name:              "name-1",
		AliasId:           102,
		ColumnTypeId:      103,
		Description:       "description-4",
		CalculationTypeId: 105,
		ShowInUi:          true,
		UserId:            107,
	}
}

// columnCatFieldErrors достаёт из ошибки список полей с претензиями.
func columnCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateColumnCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateColumnCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateColumnCatRequest) {}},
		{name: "zero table_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "negative table_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.TableId = -1 }, wantField: "table_id"},
		{name: "empty name", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Name = strings.Repeat("a", columnCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero alias_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.AliasId = 0 }, wantField: "alias_id"},
		{name: "negative alias_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.AliasId = -1 }, wantField: "alias_id"},
		{name: "zero column_type_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.ColumnTypeId = 0 }, wantField: "column_type_id"},
		{name: "negative column_type_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.ColumnTypeId = -1 }, wantField: "column_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateColumnCatRequest) {
			r.Description = strings.Repeat("a", columnCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "negative calculation_type_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.CalculationTypeId = -1 }, wantField: "calculation_type_id"},
		{name: "zero user_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateColumnCatRequest()
			tt.mutate(req)

			err := ValidateCreateColumnCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateColumnCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateColumnCat() = nil, want error on %q", tt.wantField)
			}

			fields := columnCatFieldErrors(t, err)

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

func TestValidateUpdateColumnCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateColumnCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateColumnCatByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero table_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "negative table_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.TableId = -1 }, wantField: "table_id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "blank name", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Name = "   " }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Name = strings.Repeat("a", columnCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero alias_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.AliasId = 0 }, wantField: "alias_id"},
		{name: "negative alias_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.AliasId = -1 }, wantField: "alias_id"},
		{name: "zero column_type_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.ColumnTypeId = 0 }, wantField: "column_type_id"},
		{name: "negative column_type_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.ColumnTypeId = -1 }, wantField: "column_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "blank description", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Description = "   " }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) {
			r.Description = strings.Repeat("a", columnCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "negative calculation_type_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.CalculationTypeId = -1 }, wantField: "calculation_type_id"},
		{name: "zero user_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.UserId = 0 }, wantField: "user_id"},
		{name: "negative user_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.UserId = -1 }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateColumnCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateColumnCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateColumnCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateColumnCatById() = nil, want error on %q", tt.wantField)
			}

			fields := columnCatFieldErrors(t, err)

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
func TestValidateCreateColumnCatAtVarcharLimit(t *testing.T) {
	req := validCreateColumnCatRequest()
	req.Name = strings.Repeat("a", columnCatNameMaxLen)

	if err := ValidateCreateColumnCat(req); err != nil {
		t.Errorf("ValidateCreateColumnCat() = %v, want nil at exactly %d chars", err, columnCatNameMaxLen)
	}
}

// Длина считается в символах, а не в байтах: кириллица занимает по 2 байта.
func TestValidateCreateColumnCatCyrillicAtVarcharLimit(t *testing.T) {
	req := validCreateColumnCatRequest()
	req.Name = strings.Repeat("я", columnCatNameMaxLen)

	if err := ValidateCreateColumnCat(req); err != nil {
		t.Errorf("ValidateCreateColumnCat() = %v, want nil at exactly %d cyrillic chars", err, columnCatNameMaxLen)
	}
}

func TestValidateCreateColumnCatCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateColumnCat(&tablesv1.CreateColumnCatRequest{})

	if err == nil {
		t.Fatal("ValidateCreateColumnCat() = nil, want errors")
	}

	fields := columnCatFieldErrors(t, err)

	wantFields := []string{"table_id", "name", "alias_id", "column_type_id", "description", "calculation_type_id", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateColumnCatNil(t *testing.T) {
	if err := ValidateCreateColumnCat(nil); err == nil {
		t.Error("ValidateCreateColumnCat(nil) = nil, want error")
	}
}

func TestValidateUpdateColumnCatByIdNil(t *testing.T) {
	if err := ValidateUpdateColumnCatById(nil); err == nil {
		t.Error("ValidateUpdateColumnCatById(nil) = nil, want error")
	}
}
