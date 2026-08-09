package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// databaseTypeFieldErrors достаёт из ошибки список полей с претензиями.
func databaseTypeFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateDatabaseTypeRequest() *tablesv1.CreateDatabaseTypeRequest {
	return &tablesv1.CreateDatabaseTypeRequest{
		Name:       "name-0",
		DbVersion:  "db-version-1",
		ExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateDatabaseType(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDatabaseTypeRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDatabaseTypeRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.Name = strings.Repeat("a", databaseTypeNameMaxLen+1) }, wantField: "name"},
		{name: "empty db_version", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.DbVersion = "" }, wantField: "db_version"},
		{name: "db_version too long", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) {
			r.DbVersion = strings.Repeat("a", databaseTypeDbVersionMaxLen+1)
		}, wantField: "db_version"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateDatabaseTypeRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateDatabaseTypeRequest()
			tt.mutate(req)

			err := ValidateCreateDatabaseType(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDatabaseType() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDatabaseType() = nil, want error on %q", tt.wantField)
			}

			fields := databaseTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateDatabaseTypeNil(t *testing.T) {
	if err := ValidateCreateDatabaseType(nil); err == nil {
		t.Error("ValidateCreateDatabaseType(nil) = nil, want error")
	}
}

func validValidateUpdateDatabaseTypeByIdRequest() *tablesv1.UpdateDatabaseTypeByIdRequest {
	return &tablesv1.UpdateDatabaseTypeByIdRequest{
		Name:       "name-0",
		DbVersion:  "db-version-1",
		ExternalId: "00000000-0000-4000-8000-000000000003",
		Id:         103,
	}
}

func TestValidateUpdateDatabaseTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDatabaseTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDatabaseTypeByIdRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) {
			r.Name = strings.Repeat("a", databaseTypeNameMaxLen+1)
		}, wantField: "name"},
		{name: "empty db_version", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.DbVersion = "" }, wantField: "db_version"},
		{name: "db_version too long", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) {
			r.DbVersion = strings.Repeat("a", databaseTypeDbVersionMaxLen+1)
		}, wantField: "db_version"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDatabaseTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateDatabaseTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDatabaseTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDatabaseTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDatabaseTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateDatabaseTypeByIdNil(t *testing.T) {
	if err := ValidateUpdateDatabaseTypeById(nil); err == nil {
		t.Error("ValidateUpdateDatabaseTypeById(nil) = nil, want error")
	}
}

func validValidateDeleteDatabaseTypeByIdRequest() *tablesv1.DeleteDatabaseTypeByIdRequest {
	return &tablesv1.DeleteDatabaseTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteDatabaseTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteDatabaseTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteDatabaseTypeByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteDatabaseTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteDatabaseTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteDatabaseTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteDatabaseTypeByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteDatabaseTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteDatabaseTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteDatabaseTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteDatabaseTypeByIdNil(t *testing.T) {
	if err := ValidateDeleteDatabaseTypeById(nil); err == nil {
		t.Error("ValidateDeleteDatabaseTypeById(nil) = nil, want error")
	}
}

func validValidateUndeleteDatabaseTypeByIdRequest() *tablesv1.UndeleteDatabaseTypeByIdRequest {
	return &tablesv1.UndeleteDatabaseTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteDatabaseTypeById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteDatabaseTypeByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteDatabaseTypeByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteDatabaseTypeByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteDatabaseTypeByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteDatabaseTypeByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteDatabaseTypeByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteDatabaseTypeById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteDatabaseTypeById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteDatabaseTypeById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseTypeFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteDatabaseTypeByIdNil(t *testing.T) {
	if err := ValidateUndeleteDatabaseTypeById(nil); err == nil {
		t.Error("ValidateUndeleteDatabaseTypeById(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseTypesSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseTypesSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseTypesSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseTypesSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseTypesSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseTypesSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseTypesSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseTypesSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseTypesSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseTypesSearchNameNil(t *testing.T) {
	if err := ValidateGetDatabaseTypesSearchName(nil); err == nil {
		t.Error("ValidateGetDatabaseTypesSearchName(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseTypes(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseTypesRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseTypesRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseTypesRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseTypesRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseTypesRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseTypesRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseTypes(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseTypesNil(t *testing.T) {
	if err := ValidateGetDatabaseTypes(nil); err == nil {
		t.Error("ValidateGetDatabaseTypes(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseTypesDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseTypesDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseTypesDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseTypesDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseTypesDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseTypesDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseTypesDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseTypesDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseTypesDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseTypesDeletedNil(t *testing.T) {
	if err := ValidateGetDatabaseTypesDeleted(nil); err == nil {
		t.Error("ValidateGetDatabaseTypesDeleted(nil) = nil, want error")
	}
}
