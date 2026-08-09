package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// columnCatFieldErrors достаёт из ошибки список полей с претензиями.
func columnCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateColumnCatRequest() *tablesv1.CreateColumnCatRequest {
	return &tablesv1.CreateColumnCatRequest{
		TableId:           100,
		Name:              "name-1",
		AliasId:           102,
		ColumnTypeId:      103,
		Description:       "description-4",
		CalculationTypeId: 105,
		ShowInUi:          true,
		ExternalId:        "00000000-0000-4000-8000-000000000008",
	}
}

func TestValidateCreateColumnCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateColumnCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateColumnCatRequest) {}},
		{name: "zero table_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "empty name", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Name = strings.Repeat("a", columnCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero alias_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.AliasId = 0 }, wantField: "alias_id"},
		{name: "zero column_type_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.ColumnTypeId = 0 }, wantField: "column_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateColumnCatRequest) {
			r.Description = strings.Repeat("a", columnCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateColumnCatRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateColumnCatRequest()
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
		})
	}
}

func TestValidateCreateColumnCatNil(t *testing.T) {
	if err := ValidateCreateColumnCat(nil); err == nil {
		t.Error("ValidateCreateColumnCat(nil) = nil, want error")
	}
}

func validValidateUpdateColumnCatByIdRequest() *tablesv1.UpdateColumnCatByIdRequest {
	return &tablesv1.UpdateColumnCatByIdRequest{
		TableId:           100,
		Name:              "name-1",
		AliasId:           102,
		ColumnTypeId:      103,
		Description:       "description-4",
		CalculationTypeId: 105,
		ShowInUi:          true,
		ExternalId:        "00000000-0000-4000-8000-000000000008",
		Id:                108,
	}
}

func TestValidateUpdateColumnCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateColumnCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateColumnCatByIdRequest) {}},
		{name: "zero table_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.TableId = 0 }, wantField: "table_id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Name = strings.Repeat("a", columnCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero alias_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.AliasId = 0 }, wantField: "alias_id"},
		{name: "zero column_type_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.ColumnTypeId = 0 }, wantField: "column_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) {
			r.Description = strings.Repeat("a", columnCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero calculation_type_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.CalculationTypeId = 0 }, wantField: "calculation_type_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateColumnCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateColumnCatByIdRequest()
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
		})
	}
}

func TestValidateUpdateColumnCatByIdNil(t *testing.T) {
	if err := ValidateUpdateColumnCatById(nil); err == nil {
		t.Error("ValidateUpdateColumnCatById(nil) = nil, want error")
	}
}

func validValidateDeleteColumnCatByIdRequest() *tablesv1.DeleteColumnCatByIdRequest {
	return &tablesv1.DeleteColumnCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteColumnCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteColumnCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteColumnCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteColumnCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteColumnCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteColumnCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteColumnCatByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteColumnCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteColumnCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteColumnCatById() = nil, want error on %q", tt.wantField)
			}

			fields := columnCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteColumnCatByIdNil(t *testing.T) {
	if err := ValidateDeleteColumnCatById(nil); err == nil {
		t.Error("ValidateDeleteColumnCatById(nil) = nil, want error")
	}
}

func validValidateUndeleteColumnCatByIdRequest() *tablesv1.UndeleteColumnCatByIdRequest {
	return &tablesv1.UndeleteColumnCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteColumnCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteColumnCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteColumnCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteColumnCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteColumnCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteColumnCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteColumnCatByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteColumnCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteColumnCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteColumnCatById() = nil, want error on %q", tt.wantField)
			}

			fields := columnCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteColumnCatByIdNil(t *testing.T) {
	if err := ValidateUndeleteColumnCatById(nil); err == nil {
		t.Error("ValidateUndeleteColumnCatById(nil) = nil, want error")
	}
}

func TestValidateGetColumnCats(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetColumnCatsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetColumnCatsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetColumnCatsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetColumnCatsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetColumnCatsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetColumnCatsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetColumnCats(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetColumnCats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetColumnCatsNil(t *testing.T) {
	if err := ValidateGetColumnCats(nil); err == nil {
		t.Error("ValidateGetColumnCats(nil) = nil, want error")
	}
}

func TestValidateGetColumnCatsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetColumnCatsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetColumnCatsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetColumnCatsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetColumnCatsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetColumnCatsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetColumnCatsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetColumnCatsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetColumnCatsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetColumnCatsDeletedNil(t *testing.T) {
	if err := ValidateGetColumnCatsDeleted(nil); err == nil {
		t.Error("ValidateGetColumnCatsDeleted(nil) = nil, want error")
	}
}

func TestValidateGetColumnCatsOptionalFilters(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetColumnCatsOptionalFiltersRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetColumnCatsOptionalFiltersRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetColumnCatsOptionalFiltersRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetColumnCatsOptionalFiltersRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetColumnCatsOptionalFiltersRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetColumnCatsOptionalFiltersRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetColumnCatsOptionalFilters(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetColumnCatsOptionalFilters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetColumnCatsOptionalFiltersNil(t *testing.T) {
	if err := ValidateGetColumnCatsOptionalFilters(nil); err == nil {
		t.Error("ValidateGetColumnCatsOptionalFilters(nil) = nil, want error")
	}
}

func TestValidateGetColumnCatsSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetColumnCatsSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetColumnCatsSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetColumnCatsSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetColumnCatsSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetColumnCatsSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetColumnCatsSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetColumnCatsSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetColumnCatsSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetColumnCatsSearchNameNil(t *testing.T) {
	if err := ValidateGetColumnCatsSearchName(nil); err == nil {
		t.Error("ValidateGetColumnCatsSearchName(nil) = nil, want error")
	}
}
