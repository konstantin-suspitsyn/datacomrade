package tables

import (
	"errors"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// validCreateFollowingCalculationRequest — заведомо корректный запрос.
// Тесты портят по одному полю, чтобы проверять правила по отдельности.
func validCreateFollowingCalculationRequest() *tablesv1.CreateFollowingCalculationRequest {
	return &tablesv1.CreateFollowingCalculationRequest{
		ColumnCatId:       100,
		CalculationTypeId: 101,
		UserExternalId:    "00000000-0000-4000-8000-000000000003",
	}
}

func validUpdateFollowingCalculationByIdRequest() *tablesv1.UpdateFollowingCalculationByIdRequest {
	return &tablesv1.UpdateFollowingCalculationByIdRequest{
		Id:                42,
		ColumnCatId:       100,
		CalculationTypeId: 101,
		UserExternalId:    "00000000-0000-4000-8000-000000000003",
	}
}

// followingCalculationFieldErrors достаёт из ошибки список полей с претензиями.
func followingCalculationFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func TestValidateCreateFollowingCalculation(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateFollowingCalculationRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateFollowingCalculationRequest) {}},
		{name: "zero column_cat_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.ColumnCatId = 0 }, wantField: "column_cat_id"},
		{name: "negative column_cat_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.ColumnCatId = -1 }, wantField: "column_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "negative calculation_type_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.CalculationTypeId = -1 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) {
			r.UserExternalId = "00000000000040008000000000000001"
		}, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateFollowingCalculationRequest()
			tt.mutate(req)

			err := ValidateCreateFollowingCalculation(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateFollowingCalculation() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateFollowingCalculation() = nil, want error on %q", tt.wantField)
			}

			fields := followingCalculationFieldErrors(t, err)

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

func TestValidateUpdateFollowingCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateFollowingCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateFollowingCalculationByIdRequest) {}},
		{name: "zero id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
		{name: "negative id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.Id = -5 }, wantField: "id"},
		{name: "zero column_cat_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.ColumnCatId = 0 }, wantField: "column_cat_id"},
		{name: "negative column_cat_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.ColumnCatId = -1 }, wantField: "column_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "negative calculation_type_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.CalculationTypeId = -1 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.UserExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.UserExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "user_id without dashes", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) {
			r.UserExternalId = "00000000000040008000000000000001"
		}, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validUpdateFollowingCalculationByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateFollowingCalculationById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateFollowingCalculationById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateFollowingCalculationById() = nil, want error on %q", tt.wantField)
			}

			fields := followingCalculationFieldErrors(t, err)

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

func TestValidateCreateFollowingCalculationCollectsAllErrors(t *testing.T) {
	// Валидатор копит ошибки, а не падает на первой: клиент видит
	// все проблемы запроса за один ответ.
	err := ValidateCreateFollowingCalculation(&tablesv1.CreateFollowingCalculationRequest{})

	if err == nil {
		t.Fatal("ValidateCreateFollowingCalculation() = nil, want errors")
	}

	fields := followingCalculationFieldErrors(t, err)

	wantFields := []string{"column_cat_id", "calculation_type_id", "user_id"}

	for _, field := range wantFields {
		if len(fields[field]) == 0 {
			t.Errorf("no error on %q", field)
		}
	}

	if len(fields) != len(wantFields) {
		t.Errorf("errors on %d fields, want %d: %v", len(fields), len(wantFields), fields)
	}
}

func TestValidateCreateFollowingCalculationNil(t *testing.T) {
	if err := ValidateCreateFollowingCalculation(nil); err == nil {
		t.Error("ValidateCreateFollowingCalculation(nil) = nil, want error")
	}
}

func TestValidateUpdateFollowingCalculationByIdNil(t *testing.T) {
	if err := ValidateUpdateFollowingCalculationById(nil); err == nil {
		t.Error("ValidateUpdateFollowingCalculationById(nil) = nil, want error")
	}
}
