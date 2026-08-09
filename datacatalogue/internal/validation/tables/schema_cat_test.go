package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// schemaCatFieldErrors достаёт из ошибки список полей с претензиями.
func schemaCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateSchemaCatRequest() *tablesv1.CreateSchemaCatRequest {
	return &tablesv1.CreateSchemaCatRequest{
		DatabaseId: 100,
		Name:       "name-1",
		ExternalId: "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateSchemaCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateSchemaCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateSchemaCatRequest) {}},
		{name: "zero database_id", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.DatabaseId = 0 }, wantField: "database_id"},
		{name: "empty name", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.Name = strings.Repeat("a", schemaCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateSchemaCatRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateSchemaCatRequest()
			tt.mutate(req)

			err := ValidateCreateSchemaCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateSchemaCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateSchemaCat() = nil, want error on %q", tt.wantField)
			}

			fields := schemaCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateSchemaCatNil(t *testing.T) {
	if err := ValidateCreateSchemaCat(nil); err == nil {
		t.Error("ValidateCreateSchemaCat(nil) = nil, want error")
	}
}

func validValidateUpdateSchemaCatByIdRequest() *tablesv1.UpdateSchemaCatByIdRequest {
	return &tablesv1.UpdateSchemaCatByIdRequest{
		DatabaseId: 100,
		Name:       "name-1",
		ExternalId: "00000000-0000-4000-8000-000000000003",
		Id:         103,
	}
}

func TestValidateUpdateSchemaCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateSchemaCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateSchemaCatByIdRequest) {}},
		{name: "zero database_id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.DatabaseId = 0 }, wantField: "database_id"},
		{name: "empty name", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Name = strings.Repeat("a", schemaCatNameMaxLen+1) }, wantField: "name"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateSchemaCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateSchemaCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateSchemaCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateSchemaCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateSchemaCatById() = nil, want error on %q", tt.wantField)
			}

			fields := schemaCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateSchemaCatByIdNil(t *testing.T) {
	if err := ValidateUpdateSchemaCatById(nil); err == nil {
		t.Error("ValidateUpdateSchemaCatById(nil) = nil, want error")
	}
}

func validValidateDeleteSchemaCatByIdRequest() *tablesv1.DeleteSchemaCatByIdRequest {
	return &tablesv1.DeleteSchemaCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteSchemaCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteSchemaCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteSchemaCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteSchemaCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteSchemaCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteSchemaCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteSchemaCatByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteSchemaCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteSchemaCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteSchemaCatById() = nil, want error on %q", tt.wantField)
			}

			fields := schemaCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteSchemaCatByIdNil(t *testing.T) {
	if err := ValidateDeleteSchemaCatById(nil); err == nil {
		t.Error("ValidateDeleteSchemaCatById(nil) = nil, want error")
	}
}

func validValidateUndeleteSchemaCatByIdRequest() *tablesv1.UndeleteSchemaCatByIdRequest {
	return &tablesv1.UndeleteSchemaCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteSchemaCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteSchemaCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteSchemaCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteSchemaCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteSchemaCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteSchemaCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteSchemaCatByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteSchemaCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteSchemaCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteSchemaCatById() = nil, want error on %q", tt.wantField)
			}

			fields := schemaCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteSchemaCatByIdNil(t *testing.T) {
	if err := ValidateUndeleteSchemaCatById(nil); err == nil {
		t.Error("ValidateUndeleteSchemaCatById(nil) = nil, want error")
	}
}

func TestValidateGetSchemaCats(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetSchemaCatsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetSchemaCatsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetSchemaCatsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetSchemaCatsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetSchemaCatsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetSchemaCatsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetSchemaCats(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetSchemaCats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetSchemaCatsNil(t *testing.T) {
	if err := ValidateGetSchemaCats(nil); err == nil {
		t.Error("ValidateGetSchemaCats(nil) = nil, want error")
	}
}

func TestValidateGetSchemaCatsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetSchemaCatsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetSchemaCatsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetSchemaCatsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetSchemaCatsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetSchemaCatsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetSchemaCatsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetSchemaCatsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetSchemaCatsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetSchemaCatsDeletedNil(t *testing.T) {
	if err := ValidateGetSchemaCatsDeleted(nil); err == nil {
		t.Error("ValidateGetSchemaCatsDeleted(nil) = nil, want error")
	}
}
