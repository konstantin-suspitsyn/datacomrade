package tables

import (
	"errors"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// followingCalculationFieldErrors достаёт из ошибки список полей с претензиями.
func followingCalculationFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateFollowingCalculationRequest() *tablesv1.CreateFollowingCalculationRequest {
	return &tablesv1.CreateFollowingCalculationRequest{
		ColumnCatId:       100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateFollowingCalculation(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateFollowingCalculationRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateFollowingCalculationRequest) {}},
		{name: "zero column_cat_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.ColumnCatId = 0 }, wantField: "column_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateFollowingCalculationRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateFollowingCalculationRequest()
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
		})
	}
}

func TestValidateCreateFollowingCalculationNil(t *testing.T) {
	if err := ValidateCreateFollowingCalculation(nil); err == nil {
		t.Error("ValidateCreateFollowingCalculation(nil) = nil, want error")
	}
}

func validValidateUpdateFollowingCalculationByIdRequest() *tablesv1.UpdateFollowingCalculationByIdRequest {
	return &tablesv1.UpdateFollowingCalculationByIdRequest{
		ColumnCatId:       100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
		Id:                103,
	}
}

func TestValidateUpdateFollowingCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateFollowingCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateFollowingCalculationByIdRequest) {}},
		{name: "zero column_cat_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.ColumnCatId = 0 }, wantField: "column_cat_id"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateFollowingCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateFollowingCalculationByIdRequest()
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
		})
	}
}

func TestValidateUpdateFollowingCalculationByIdNil(t *testing.T) {
	if err := ValidateUpdateFollowingCalculationById(nil); err == nil {
		t.Error("ValidateUpdateFollowingCalculationById(nil) = nil, want error")
	}
}

func validValidateDeleteFollowingCalculationByIdRequest() *tablesv1.DeleteFollowingCalculationByIdRequest {
	return &tablesv1.DeleteFollowingCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteFollowingCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteFollowingCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteFollowingCalculationByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteFollowingCalculationByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteFollowingCalculationByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteFollowingCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteFollowingCalculationByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteFollowingCalculationById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteFollowingCalculationById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteFollowingCalculationById() = nil, want error on %q", tt.wantField)
			}

			fields := followingCalculationFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteFollowingCalculationByIdNil(t *testing.T) {
	if err := ValidateDeleteFollowingCalculationById(nil); err == nil {
		t.Error("ValidateDeleteFollowingCalculationById(nil) = nil, want error")
	}
}

func validValidateUndeleteFollowingCalculationByIdRequest() *tablesv1.UndeleteFollowingCalculationByIdRequest {
	return &tablesv1.UndeleteFollowingCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteFollowingCalculationById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteFollowingCalculationByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteFollowingCalculationByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteFollowingCalculationByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteFollowingCalculationByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteFollowingCalculationByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteFollowingCalculationByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteFollowingCalculationById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteFollowingCalculationById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteFollowingCalculationById() = nil, want error on %q", tt.wantField)
			}

			fields := followingCalculationFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteFollowingCalculationByIdNil(t *testing.T) {
	if err := ValidateUndeleteFollowingCalculationById(nil); err == nil {
		t.Error("ValidateUndeleteFollowingCalculationById(nil) = nil, want error")
	}
}

func TestValidateGetFollowingCalculations(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetFollowingCalculationsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetFollowingCalculationsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetFollowingCalculationsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetFollowingCalculationsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetFollowingCalculationsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetFollowingCalculationsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetFollowingCalculations(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetFollowingCalculations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetFollowingCalculationsNil(t *testing.T) {
	if err := ValidateGetFollowingCalculations(nil); err == nil {
		t.Error("ValidateGetFollowingCalculations(nil) = nil, want error")
	}
}

func TestValidateGetFollowingCalculationsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetFollowingCalculationsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetFollowingCalculationsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetFollowingCalculationsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetFollowingCalculationsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetFollowingCalculationsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetFollowingCalculationsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetFollowingCalculationsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetFollowingCalculationsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetFollowingCalculationsDeletedNil(t *testing.T) {
	if err := ValidateGetFollowingCalculationsDeleted(nil); err == nil {
		t.Error("ValidateGetFollowingCalculationsDeleted(nil) = nil, want error")
	}
}

func TestValidateGetFollowingCalculationsFilters(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetFollowingCalculationsFiltersRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetFollowingCalculationsFiltersRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetFollowingCalculationsFiltersRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetFollowingCalculationsFiltersRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetFollowingCalculationsFiltersRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetFollowingCalculationsFiltersRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetFollowingCalculationsFilters(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetFollowingCalculationsFilters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetFollowingCalculationsFiltersNil(t *testing.T) {
	if err := ValidateGetFollowingCalculationsFilters(nil); err == nil {
		t.Error("ValidateGetFollowingCalculationsFilters(nil) = nil, want error")
	}
}
