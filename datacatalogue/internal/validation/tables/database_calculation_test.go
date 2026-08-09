package tables

import (
	"errors"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// databaseCalculationFieldErrors достаёт из ошибки список полей с претензиями.
func databaseCalculationFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateDatabaseCalculationRequest() *tablesv1.CreateDatabaseCalculationRequest {
	return &tablesv1.CreateDatabaseCalculationRequest{
		DatabaseCatId:     100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateDatabaseCalculation(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDatabaseCalculationRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDatabaseCalculationRequest) {}},
		{name: "zero database_cat_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.DatabaseCatId = 0 }, wantField: "database_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateDatabaseCalculationRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateDatabaseCalculationRequest()
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
		})
	}
}

func TestValidateCreateDatabaseCalculationNil(t *testing.T) {
	if err := ValidateCreateDatabaseCalculation(nil); err == nil {
		t.Error("ValidateCreateDatabaseCalculation(nil) = nil, want error")
	}
}

func validValidateUpdateDatabaseCalculationByIdRequest() *tablesv1.UpdateDatabaseCalculationByIdRequest {
	return &tablesv1.UpdateDatabaseCalculationByIdRequest{
		DatabaseCatId:     100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
		Id:                103,
	}
}

func TestValidateUpdateDatabaseCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDatabaseCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDatabaseCalculationByIdRequest) {}},
		{name: "zero database_cat_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.DatabaseCatId = 0 }, wantField: "database_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDatabaseCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateDatabaseCalculationByIdRequest()
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
		})
	}
}

func TestValidateUpdateDatabaseCalculationByIdNil(t *testing.T) {
	if err := ValidateUpdateDatabaseCalculationById(nil); err == nil {
		t.Error("ValidateUpdateDatabaseCalculationById(nil) = nil, want error")
	}
}

func validValidateDeleteDatabaseCalculationByIdRequest() *tablesv1.DeleteDatabaseCalculationByIdRequest {
	return &tablesv1.DeleteDatabaseCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteDatabaseCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteDatabaseCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteDatabaseCalculationByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteDatabaseCalculationByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteDatabaseCalculationByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteDatabaseCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteDatabaseCalculationByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteDatabaseCalculationById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteDatabaseCalculationById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteDatabaseCalculationById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCalculationFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteDatabaseCalculationByIdNil(t *testing.T) {
	if err := ValidateDeleteDatabaseCalculationById(nil); err == nil {
		t.Error("ValidateDeleteDatabaseCalculationById(nil) = nil, want error")
	}
}

func validValidateUndeleteDatabaseCalculationByIdRequest() *tablesv1.UndeleteDatabaseCalculationByIdRequest {
	return &tablesv1.UndeleteDatabaseCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteDatabaseCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteDatabaseCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteDatabaseCalculationByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteDatabaseCalculationByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteDatabaseCalculationByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteDatabaseCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteDatabaseCalculationByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteDatabaseCalculationById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteDatabaseCalculationById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteDatabaseCalculationById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCalculationFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteDatabaseCalculationByIdNil(t *testing.T) {
	if err := ValidateUndeleteDatabaseCalculationById(nil); err == nil {
		t.Error("ValidateUndeleteDatabaseCalculationById(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCalculationsSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCalculationsSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCalculationsSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCalculationsSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCalculationsSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCalculationsSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCalculationsSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCalculationsSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCalculationsSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCalculationsSearchNameNil(t *testing.T) {
	if err := ValidateGetDatabaseCalculationsSearchName(nil); err == nil {
		t.Error("ValidateGetDatabaseCalculationsSearchName(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCalculationsFiltered(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCalculationsFilteredRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCalculationsFilteredRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCalculationsFilteredRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCalculationsFilteredRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCalculationsFilteredRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCalculationsFilteredRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCalculationsFiltered(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCalculationsFiltered() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCalculationsFilteredNil(t *testing.T) {
	if err := ValidateGetDatabaseCalculationsFiltered(nil); err == nil {
		t.Error("ValidateGetDatabaseCalculationsFiltered(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCalculationsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCalculationsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCalculationsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCalculationsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCalculationsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCalculationsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCalculationsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCalculationsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCalculationsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCalculationsDeletedNil(t *testing.T) {
	if err := ValidateGetDatabaseCalculationsDeleted(nil); err == nil {
		t.Error("ValidateGetDatabaseCalculationsDeleted(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCalculations(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCalculationsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCalculationsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCalculationsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCalculationsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCalculationsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCalculationsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCalculations(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCalculations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCalculationsNil(t *testing.T) {
	if err := ValidateGetDatabaseCalculations(nil); err == nil {
		t.Error("ValidateGetDatabaseCalculations(nil) = nil, want error")
	}
}
