package tables

import (
	"errors"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateDatabaseCalculationRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateDatabaseCalculationRequest() *tablesv1.CreateDatabaseCalculationRequest {
	return &tablesv1.CreateDatabaseCalculationRequest{
		DatabaseCatId:     100,
		CalculationTypeId: 101,
		UserExternalId:    "00000000-0000-4000-8000-000000000003",
	}
}

func validUpdateDatabaseCalculationByIdRequest() *tablesv1.UpdateDatabaseCalculationByIdRequest {
	return &tablesv1.UpdateDatabaseCalculationByIdRequest{
		Id:                42,
		DatabaseCatId:     100,
		CalculationTypeId: 101,
		UserExternalId:    "00000000-0000-4000-8000-000000000003",
	}
}

// databaseCalculationFieldErrors достаёт из ошибки список полей с претензиями.
func databaseCalculationFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateDatabaseCalculation(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDatabaseCalculationRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDatabaseCalculationRequest) {}},
		{name: "zero database_cat_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.DatabaseCatId = 0 }, wantField: "database_cat_id"},
		{name: "negative database_cat_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.DatabaseCatId = -1 }, wantField: "database_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "negative calculation_type_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.CalculationTypeId = -1 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) {
			r.UserExternalId = "00000000000040008000000000000001"
		}, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateDatabaseCalculationRequest()
			tt.mutate(req)

			err := ValidateCreateDatabaseCalculation(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDatabaseCalculation() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDatabaseCalculation() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCalculationFieldErrors(t, err)

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

func TestValidateUpdateDatabaseCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDatabaseCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDatabaseCalculationByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero database_cat_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.DatabaseCatId = 0 }, wantField: "database_cat_id"},
		{name: "negative database_cat_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.DatabaseCatId = -1 }, wantField: "database_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "negative calculation_type_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.CalculationTypeId = -1 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) {
			r.UserExternalId = "00000000000040008000000000000001"
		}, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateDatabaseCalculationByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDatabaseCalculationById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDatabaseCalculationById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDatabaseCalculationById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCalculationFieldErrors(t, err)

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

func TestValidateCreateDatabaseCalculationCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateDatabaseCalculation(&tablesv1.CreateDatabaseCalculationRequest{})

	if err == nil {
		t.Fatal("ValidateCreateDatabaseCalculation() = nil, want errors")
	}

	fields := databaseCalculationFieldErrors(t, err)

	wantFields := []string{"database_cat_id", "calculation_type_id", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateDatabaseCalculationNil(t *testing.T) {
	if err := ValidateCreateDatabaseCalculation(nil); err == nil {
		t.Error("ValidateCreateDatabaseCalculation(nil) = nil, want error")
	}
}

func TestValidateUpdateDatabaseCalculationByIdNil(t *testing.T) {
	if err := ValidateUpdateDatabaseCalculationById(nil); err == nil {
		t.Error("ValidateUpdateDatabaseCalculationById(nil) = nil, want error")
	}
}
