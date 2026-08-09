package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// domainCatFieldErrors достаёт из ошибки список полей с претензиями.
func domainCatFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateDomainCatRequest() *tablesv1.CreateDomainCatRequest {
	return &tablesv1.CreateDomainCatRequest{
		DomainName: "domain-name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
	}
}

func TestValidateCreateDomainCat(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateDomainCatRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateDomainCatRequest) {}},
		{name: "empty domain_name", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.DomainName = "" }, wantField: "domain_name"},
		{name: "domain_name too long", mutate: func(r *tablesv1.CreateDomainCatRequest) {
			r.DomainName = strings.Repeat("a", domainCatDomainNameMaxLen+1)
		}, wantField: "domain_name"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateDomainCatRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateDomainCatRequest()
			tt.mutate(req)

			err := ValidateCreateDomainCat(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateDomainCat() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateDomainCat() = nil, want error on %q", tt.wantField)
			}

			fields := domainCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateDomainCatNil(t *testing.T) {
	if err := ValidateCreateDomainCat(nil); err == nil {
		t.Error("ValidateCreateDomainCat(nil) = nil, want error")
	}
}

func validValidateUpdateDomainCatByIdRequest() *tablesv1.UpdateDomainCatByIdRequest {
	return &tablesv1.UpdateDomainCatByIdRequest{
		DomainName: "domain-name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
		Id:         102,
	}
}

func TestValidateUpdateDomainCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateDomainCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateDomainCatByIdRequest) {}},
		{name: "empty domain_name", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.DomainName = "" }, wantField: "domain_name"},
		{name: "domain_name too long", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) {
			r.DomainName = strings.Repeat("a", domainCatDomainNameMaxLen+1)
		}, wantField: "domain_name"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateDomainCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateDomainCatByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateDomainCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateDomainCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateDomainCatById() = nil, want error on %q", tt.wantField)
			}

			fields := domainCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateDomainCatByIdNil(t *testing.T) {
	if err := ValidateUpdateDomainCatById(nil); err == nil {
		t.Error("ValidateUpdateDomainCatById(nil) = nil, want error")
	}
}

func validValidateDeleteDomainCatByIdRequest() *tablesv1.DeleteDomainCatByIdRequest {
	return &tablesv1.DeleteDomainCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteDomainCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteDomainCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteDomainCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteDomainCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteDomainCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteDomainCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteDomainCatByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteDomainCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteDomainCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteDomainCatById() = nil, want error on %q", tt.wantField)
			}

			fields := domainCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteDomainCatByIdNil(t *testing.T) {
	if err := ValidateDeleteDomainCatById(nil); err == nil {
		t.Error("ValidateDeleteDomainCatById(nil) = nil, want error")
	}
}

func validValidateUndeleteDomainCatByIdRequest() *tablesv1.UndeleteDomainCatByIdRequest {
	return &tablesv1.UndeleteDomainCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteDomainCatById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteDomainCatByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteDomainCatByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteDomainCatByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteDomainCatByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteDomainCatByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteDomainCatByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteDomainCatById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteDomainCatById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteDomainCatById() = nil, want error on %q", tt.wantField)
			}

			fields := domainCatFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteDomainCatByIdNil(t *testing.T) {
	if err := ValidateUndeleteDomainCatById(nil); err == nil {
		t.Error("ValidateUndeleteDomainCatById(nil) = nil, want error")
	}
}

func TestValidateGetDomainCatsSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDomainCatsSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDomainCatsSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDomainCatsSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDomainCatsSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDomainCatsSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDomainCatsSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDomainCatsSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDomainCatsSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDomainCatsSearchNameNil(t *testing.T) {
	if err := ValidateGetDomainCatsSearchName(nil); err == nil {
		t.Error("ValidateGetDomainCatsSearchName(nil) = nil, want error")
	}
}

func TestValidateGetDomainCats(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDomainCatsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDomainCatsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDomainCatsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDomainCatsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDomainCatsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDomainCatsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDomainCats(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDomainCats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDomainCatsNil(t *testing.T) {
	if err := ValidateGetDomainCats(nil); err == nil {
		t.Error("ValidateGetDomainCats(nil) = nil, want error")
	}
}

func TestValidateGetDomainCatsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetDomainCatsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetDomainCatsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetDomainCatsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetDomainCatsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetDomainCatsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetDomainCatsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetDomainCatsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetDomainCatsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetDomainCatsDeletedNil(t *testing.T) {
	if err := ValidateGetDomainCatsDeleted(nil); err == nil {
		t.Error("ValidateGetDomainCatsDeleted(nil) = nil, want error")
	}
}
