package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// aliasFieldErrors достаёт из ошибки список полей с претензиями.
func aliasFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateAliasRequest() *tablesv1.CreateAliasRequest {
	return &tablesv1.CreateAliasRequest{
		Name:        "name-0",
		Description: "description-1",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateAlias(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateAliasRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateAliasRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.CreateAliasRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.CreateAliasRequest) { r.Name = strings.Repeat("a", aliasNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.CreateAliasRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateAliasRequest) { r.Description = strings.Repeat("a", aliasDescriptionMaxLen+1) }, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateAliasRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateAliasRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateAliasRequest()
			tt.mutate(req)

			err := ValidateCreateAlias(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateAlias() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateAlias() = nil, want error on %q", tt.wantField)
			}

			fields := aliasFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateAliasNil(t *testing.T) {
	if err := ValidateCreateAlias(nil); err == nil {
		t.Error("ValidateCreateAlias(nil) = nil, want error")
	}
}

func validValidateUpdateAliasByIdRequest() *tablesv1.UpdateAliasByIdRequest {
	return &tablesv1.UpdateAliasByIdRequest{
		Name:        "name-0",
		Description: "description-1",
		IsDeleted:   true,
		ExternalId:  "00000000-0000-4000-8000-000000000004",
		Id:          104,
	}
}

func TestValidateUpdateAliasById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateAliasByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateAliasByIdRequest) {}},
		{name: "empty name", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Name = "" }, wantField: "name"},
		{name: "name too long", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Name = strings.Repeat("a", aliasNameMaxLen+1) }, wantField: "name"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateAliasByIdRequest) {
			r.Description = strings.Repeat("a", aliasDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateAliasByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateAliasByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateAliasById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateAliasById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateAliasById() = nil, want error on %q", tt.wantField)
			}

			fields := aliasFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateAliasByIdNil(t *testing.T) {
	if err := ValidateUpdateAliasById(nil); err == nil {
		t.Error("ValidateUpdateAliasById(nil) = nil, want error")
	}
}

func validValidateDeleteAliasByIdRequest() *tablesv1.DeleteAliasByIdRequest {
	return &tablesv1.DeleteAliasByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteAliasById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteAliasByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteAliasByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteAliasByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteAliasByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteAliasByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteAliasByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteAliasById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteAliasById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteAliasById() = nil, want error on %q", tt.wantField)
			}

			fields := aliasFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteAliasByIdNil(t *testing.T) {
	if err := ValidateDeleteAliasById(nil); err == nil {
		t.Error("ValidateDeleteAliasById(nil) = nil, want error")
	}
}

func validValidateUndeleteAliasByIdRequest() *tablesv1.UndeleteAliasByIdRequest {
	return &tablesv1.UndeleteAliasByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteAliasById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteAliasByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteAliasByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteAliasByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteAliasByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteAliasByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteAliasByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteAliasById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteAliasById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteAliasById() = nil, want error on %q", tt.wantField)
			}

			fields := aliasFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteAliasByIdNil(t *testing.T) {
	if err := ValidateUndeleteAliasById(nil); err == nil {
		t.Error("ValidateUndeleteAliasById(nil) = nil, want error")
	}
}

func TestValidateGetAliasesDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetAliasesDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetAliasesDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetAliasesDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetAliasesDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetAliasesDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetAliasesDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetAliasesDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetAliasesDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetAliasesDeletedNil(t *testing.T) {
	if err := ValidateGetAliasesDeleted(nil); err == nil {
		t.Error("ValidateGetAliasesDeleted(nil) = nil, want error")
	}
}

func TestValidateGetAliases(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetAliasesRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetAliasesRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetAliasesRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetAliasesRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetAliasesRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetAliasesRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetAliases(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetAliases() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetAliasesNil(t *testing.T) {
	if err := ValidateGetAliases(nil); err == nil {
		t.Error("ValidateGetAliases(nil) = nil, want error")
	}
}
