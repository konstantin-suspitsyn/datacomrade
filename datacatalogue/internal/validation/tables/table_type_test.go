package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tableTypeFieldErrors достаёт из ошибки список полей с претензиями.
func tableTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateTableTypeRequest() *tablesv1.CreateTableTypeRequest {
	return &tablesv1.CreateTableTypeRequest{
		Name:        "name-0",
		Description: "description-1",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateTableType(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateTableTypeRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateTableTypeRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Name = strings.Repeat("a", tableTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateTableTypeRequest) {
			r.Description = strings.Repeat("a", tableTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateTableTypeRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateTableTypeRequest()
			tt.mutate(req)

			err := ValidateCreateTableType(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateTableType() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateTableType() = nil, want error on %q", tt.wantField)
			}

			fields := tableTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateTableTypeNil(t *testing.T) {
	if err := ValidateCreateTableType(nil); err == nil {
		t.Error("ValidateCreateTableType(nil) = nil, want error")
	}
}

func validValidateUpdateTableTypeByIdRequest() *tablesv1.UpdateTableTypeByIdRequest {
	return &tablesv1.UpdateTableTypeByIdRequest{
		Name:        "name-0",
		Description: "description-1",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
		Id:          103,
	}
}

func TestValidateUpdateTableTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateTableTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateTableTypeByIdRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Name = strings.Repeat("a", tableTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) {
			r.Description = strings.Repeat("a", tableTypeDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateTableTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateTableTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateTableTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateTableTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateTableTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := tableTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateTableTypeByIdNil(t *testing.T) {
	if err := ValidateUpdateTableTypeById(nil); err == nil {
		t.Error("ValidateUpdateTableTypeById(nil) = nil, want error")
	}
}

func validValidateDeleteTableTypeByIdRequest() *tablesv1.DeleteTableTypeByIdRequest {
	return &tablesv1.DeleteTableTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteTableTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteTableTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteTableTypeByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteTableTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteTableTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteTableTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteTableTypeByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteTableTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteTableTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteTableTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := tableTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteTableTypeByIdNil(t *testing.T) {
	if err := ValidateDeleteTableTypeById(nil); err == nil {
		t.Error("ValidateDeleteTableTypeById(nil) = nil, want error")
	}
}

func validValidateUndeleteTableTypeByIdRequest() *tablesv1.UndeleteTableTypeByIdRequest {
	return &tablesv1.UndeleteTableTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteTableTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteTableTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteTableTypeByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteTableTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteTableTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteTableTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteTableTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteTableTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteTableTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteTableTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := tableTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteTableTypeByIdNil(t *testing.T) {
	if err := ValidateUndeleteTableTypeById(nil); err == nil {
		t.Error("ValidateUndeleteTableTypeById(nil) = nil, want error")
	}
}

func TestValidateGetTableTypes(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetTableTypesRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetTableTypesRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetTableTypesRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetTableTypesRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetTableTypesRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetTableTypesRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetTableTypes(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetTableTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetTableTypesNil(t *testing.T) {
	if err := ValidateGetTableTypes(nil); err == nil {
		t.Error("ValidateGetTableTypes(nil) = nil, want error")
	}
}

func TestValidateGetTableTypesDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetTableTypesDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetTableTypesDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetTableTypesDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetTableTypesDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetTableTypesDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetTableTypesDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetTableTypesDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetTableTypesDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetTableTypesDeletedNil(t *testing.T) {
	if err := ValidateGetTableTypesDeleted(nil); err == nil {
		t.Error("ValidateGetTableTypesDeleted(nil) = nil, want error")
	}
}
