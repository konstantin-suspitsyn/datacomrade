package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tableCatFieldErrors достаёт из ошибки список полей с претензиями.
func tableCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateTableCatRequest() *tablesv1.CreateTableCatRequest {
	return &tablesv1.CreateTableCatRequest{
		Name:        "name-0",
		Description: "description-1",
		SchemaId:    102,
		TableTypeId: 103,
		DomainId:    104,
		ExternalId:  "00000000-0000-4000-8000-000000000006",
	}
}

func TestValidateCreateTableCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateTableCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateTableCatRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Name = strings.Repeat("a", tableCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateTableCatRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateTableCatRequest) {
			r.Description = strings.Repeat("a", tableCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero schema_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.SchemaId = 0 }, wantField: "schema_id"},
		{name: "zero table_type_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.TableTypeId = 0 }, wantField: "table_type_id"},
		{name: "zero domain_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateTableCatRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateTableCatRequest()
			tt.mutate(req)

			err := ValidateCreateTableCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateTableCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateTableCat() = nil, want error on %q", tt.wantField)
			}

			fields := tableCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateTableCatNil(t *testing.T) {
	if err := ValidateCreateTableCat(nil); err == nil {
		t.Error("ValidateCreateTableCat(nil) = nil, want error")
	}
}

func validValidateUpdateTableCatByIdRequest() *tablesv1.UpdateTableCatByIdRequest {
	return &tablesv1.UpdateTableCatByIdRequest{
		Name:        "name-0",
		Description: "description-1",
		SchemaId:    102,
		TableTypeId: 103,
		DomainId:    104,
		IsGetDict:   true,
		ExternalId:  "00000000-0000-4000-8000-000000000007",
		Id:          107,
	}
}

func TestValidateUpdateTableCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateTableCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateTableCatByIdRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Name = strings.Repeat("a", tableCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) {
			r.Description = strings.Repeat("a", tableCatDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "zero schema_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.SchemaId = 0 }, wantField: "schema_id"},
		{name: "zero table_type_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.TableTypeId = 0 }, wantField: "table_type_id"},
		{name: "zero domain_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.DomainId = 0 }, wantField: "domain_id"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateTableCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateTableCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateTableCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateTableCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateTableCatById() = nil, want error on %q", tt.wantField)
			}

			fields := tableCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateTableCatByIdNil(t *testing.T) {
	if err := ValidateUpdateTableCatById(nil); err == nil {
		t.Error("ValidateUpdateTableCatById(nil) = nil, want error")
	}
}

func validValidateDeleteTableCatByIdRequest() *tablesv1.DeleteTableCatByIdRequest {
	return &tablesv1.DeleteTableCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteTableCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteTableCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteTableCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteTableCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteTableCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteTableCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteTableCatByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteTableCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteTableCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteTableCatById() = nil, want error on %q", tt.wantField)
			}

			fields := tableCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteTableCatByIdNil(t *testing.T) {
	if err := ValidateDeleteTableCatById(nil); err == nil {
		t.Error("ValidateDeleteTableCatById(nil) = nil, want error")
	}
}

func validValidateUndeleteTableCatByIdRequest() *tablesv1.UndeleteTableCatByIdRequest {
	return &tablesv1.UndeleteTableCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteTableCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteTableCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteTableCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteTableCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteTableCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteTableCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteTableCatByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteTableCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteTableCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteTableCatById() = nil, want error on %q", tt.wantField)
			}

			fields := tableCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteTableCatByIdNil(t *testing.T) {
	if err := ValidateUndeleteTableCatById(nil); err == nil {
		t.Error("ValidateUndeleteTableCatById(nil) = nil, want error")
	}
}

func TestValidateGetTableCats(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetTableCatsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetTableCatsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetTableCatsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetTableCatsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetTableCatsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetTableCatsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetTableCats(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetTableCats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetTableCatsNil(t *testing.T) {
	if err := ValidateGetTableCats(nil); err == nil {
		t.Error("ValidateGetTableCats(nil) = nil, want error")
	}
}

func TestValidateGetTableCatsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetTableCatsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetTableCatsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetTableCatsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetTableCatsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetTableCatsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetTableCatsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetTableCatsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetTableCatsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetTableCatsDeletedNil(t *testing.T) {
	if err := ValidateGetTableCatsDeleted(nil); err == nil {
		t.Error("ValidateGetTableCatsDeleted(nil) = nil, want error")
	}
}
