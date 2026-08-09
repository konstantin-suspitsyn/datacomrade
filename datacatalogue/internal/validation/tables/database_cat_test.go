package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// databaseCatFieldErrors достаёт из ошибки список полей с претензиями.
func databaseCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateDatabaseCatRequest() *tablesv1.CreateDatabaseCatRequest {
	return &tablesv1.CreateDatabaseCatRequest{
		Name:           "name-0",
		HostId:         101,
		DatabaseTypeId: 102,
		Description:    "description-3",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
	}
}

func TestValidateCreateDatabaseCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDatabaseCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDatabaseCatRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Name = strings.Repeat("a", databaseCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero host_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.HostId = 0 }, wantField: "host_id"},
		{name: "zero database_type_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.DatabaseTypeId = 0 }, wantField: "database_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateDatabaseCatRequest) {
			r.Description = strings.Repeat("a", databaseCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateDatabaseCatRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateDatabaseCatRequest()
			tt.mutate(req)

			err := ValidateCreateDatabaseCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDatabaseCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDatabaseCat() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateDatabaseCatNil(t *testing.T) {
	if err := ValidateCreateDatabaseCat(nil); err == nil {
		t.Error("ValidateCreateDatabaseCat(nil) = nil, want error")
	}
}

func validValidateUpdateDatabaseCatByIdRequest() *tablesv1.UpdateDatabaseCatByIdRequest {
	return &tablesv1.UpdateDatabaseCatByIdRequest{
		Name:           "name-0",
		HostId:         101,
		DatabaseTypeId: 102,
		Description:    "description-3",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
		Id:             105,
	}
}

func TestValidateUpdateDatabaseCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDatabaseCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDatabaseCatByIdRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Name = strings.Repeat("a", databaseCatNameMaxLen+1) }, wantField: "name"},
		{name: "zero host_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.HostId = 0 }, wantField: "host_id"},
		{name: "zero database_type_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.DatabaseTypeId = 0 }, wantField: "database_type_id"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) {
			r.Description = strings.Repeat("a", databaseCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDatabaseCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateDatabaseCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDatabaseCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDatabaseCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDatabaseCatById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateDatabaseCatByIdNil(t *testing.T) {
	if err := ValidateUpdateDatabaseCatById(nil); err == nil {
		t.Error("ValidateUpdateDatabaseCatById(nil) = nil, want error")
	}
}

func validValidateDeleteDatabaseCatByIdRequest() *tablesv1.DeleteDatabaseCatByIdRequest {
	return &tablesv1.DeleteDatabaseCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteDatabaseCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteDatabaseCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteDatabaseCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteDatabaseCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteDatabaseCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteDatabaseCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteDatabaseCatByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteDatabaseCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteDatabaseCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteDatabaseCatById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteDatabaseCatByIdNil(t *testing.T) {
	if err := ValidateDeleteDatabaseCatById(nil); err == nil {
		t.Error("ValidateDeleteDatabaseCatById(nil) = nil, want error")
	}
}

func validValidateUndeleteDatabaseCatByIdRequest() *tablesv1.UndeleteDatabaseCatByIdRequest {
	return &tablesv1.UndeleteDatabaseCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteDatabaseCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteDatabaseCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteDatabaseCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteDatabaseCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteDatabaseCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteDatabaseCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteDatabaseCatByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteDatabaseCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteDatabaseCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteDatabaseCatById() = nil, want error on %q", tt.wantField)
			}

			fields := databaseCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteDatabaseCatByIdNil(t *testing.T) {
	if err := ValidateUndeleteDatabaseCatById(nil); err == nil {
		t.Error("ValidateUndeleteDatabaseCatById(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCats(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCatsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCatsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCatsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCatsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCatsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCatsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCats(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCatsNil(t *testing.T) {
	if err := ValidateGetDatabaseCats(nil); err == nil {
		t.Error("ValidateGetDatabaseCats(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCatsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCatsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCatsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCatsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCatsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCatsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCatsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCatsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCatsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCatsDeletedNil(t *testing.T) {
	if err := ValidateGetDatabaseCatsDeleted(nil); err == nil {
		t.Error("ValidateGetDatabaseCatsDeleted(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCatsByHostIdDatabaseTypeId(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{Order: "sideways"}, wantErr: true},
		{name: "valid order_by", req: &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{OrderBy: "name"}},
		{name: "unknown order_by", req: &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{OrderBy: "nonexistent"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCatsByHostIdDatabaseTypeId(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCatsByHostIdDatabaseTypeId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCatsByHostIdDatabaseTypeIdNil(t *testing.T) {
	if err := ValidateGetDatabaseCatsByHostIdDatabaseTypeId(nil); err == nil {
		t.Error("ValidateGetDatabaseCatsByHostIdDatabaseTypeId(nil) = nil, want error")
	}
}

func TestValidateGetDatabaseCatsSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDatabaseCatsSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDatabaseCatsSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDatabaseCatsSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDatabaseCatsSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDatabaseCatsSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDatabaseCatsSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDatabaseCatsSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDatabaseCatsSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDatabaseCatsSearchNameNil(t *testing.T) {
	if err := ValidateGetDatabaseCatsSearchName(nil); err == nil {
		t.Error("ValidateGetDatabaseCatsSearchName(nil) = nil, want error")
	}
}
