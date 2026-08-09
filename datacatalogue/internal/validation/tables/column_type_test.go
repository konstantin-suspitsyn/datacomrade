package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// columnTypeFieldErrors достаёт из ошибки список полей с претензиями.
func columnTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateColumnTypeRequest() *tablesv1.CreateColumnTypeRequest {
	return &tablesv1.CreateColumnTypeRequest{
		Name:        "name-0",
		Description: "description-1",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateColumnType(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateColumnTypeRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateColumnTypeRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Name = strings.Repeat("a", columnTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateColumnTypeRequest) {
			r.Description = strings.Repeat("a", columnTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateColumnTypeRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateColumnTypeRequest()
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
		})
	}
}

func TestValidateCreateColumnTypeNil(t *testing.T) {
	if err := ValidateCreateColumnType(nil); err == nil {
		t.Error("ValidateCreateColumnType(nil) = nil, want error")
	}
}

func validValidateUpdateColumnTypeByIdRequest() *tablesv1.UpdateColumnTypeByIdRequest {
	return &tablesv1.UpdateColumnTypeByIdRequest{
		Name:        "name-0",
		Description: "description-1",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
		Id:          103,
	}
}

func TestValidateUpdateColumnTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateColumnTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateColumnTypeByIdRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Name = strings.Repeat("a", columnTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) {
			r.Description = strings.Repeat("a", columnTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateColumnTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateColumnTypeByIdRequest()
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
		})
	}
}

func TestValidateUpdateColumnTypeByIdNil(t *testing.T) {
	if err := ValidateUpdateColumnTypeById(nil); err == nil {
		t.Error("ValidateUpdateColumnTypeById(nil) = nil, want error")
	}
}

func validValidateDeleteColumnTypeByIdRequest() *tablesv1.DeleteColumnTypeByIdRequest {
	return &tablesv1.DeleteColumnTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteColumnTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteColumnTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteColumnTypeByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteColumnTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteColumnTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteColumnTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteColumnTypeByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteColumnTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteColumnTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteColumnTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := columnTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteColumnTypeByIdNil(t *testing.T) {
	if err := ValidateDeleteColumnTypeById(nil); err == nil {
		t.Error("ValidateDeleteColumnTypeById(nil) = nil, want error")
	}
}

func validValidateUndeleteColumnTypeByIdRequest() *tablesv1.UndeleteColumnTypeByIdRequest {
	return &tablesv1.UndeleteColumnTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteColumnTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteColumnTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteColumnTypeByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteColumnTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteColumnTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteColumnTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteColumnTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteColumnTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteColumnTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteColumnTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := columnTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteColumnTypeByIdNil(t *testing.T) {
	if err := ValidateUndeleteColumnTypeById(nil); err == nil {
		t.Error("ValidateUndeleteColumnTypeById(nil) = nil, want error")
	}
}

func TestValidateGetColumnTypes(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetColumnTypesRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetColumnTypesRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetColumnTypesRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetColumnTypesRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetColumnTypesRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetColumnTypesRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetColumnTypes(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetColumnTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetColumnTypesNil(t *testing.T) {
	if err := ValidateGetColumnTypes(nil); err == nil {
		t.Error("ValidateGetColumnTypes(nil) = nil, want error")
	}
}

func TestValidateGetColumnTypesSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetColumnTypesSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetColumnTypesSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetColumnTypesSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetColumnTypesSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetColumnTypesSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetColumnTypesSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetColumnTypesSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetColumnTypesSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetColumnTypesSearchNameNil(t *testing.T) {
	if err := ValidateGetColumnTypesSearchName(nil); err == nil {
		t.Error("ValidateGetColumnTypesSearchName(nil) = nil, want error")
	}
}

func TestValidateGetColumnTypesDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetColumnTypesDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetColumnTypesDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetColumnTypesDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetColumnTypesDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetColumnTypesDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetColumnTypesDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetColumnTypesDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetColumnTypesDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetColumnTypesDeletedNil(t *testing.T) {
	if err := ValidateGetColumnTypesDeleted(nil); err == nil {
		t.Error("ValidateGetColumnTypesDeleted(nil) = nil, want error")
	}
}
