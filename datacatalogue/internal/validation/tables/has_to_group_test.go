package tables

import (
	"errors"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// hasToGroupFieldErrors достаёт из ошибки список полей с претензиями.
func hasToGroupFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateHasToGroupRequest() *tablesv1.CreateHasToGroupRequest {
	return &tablesv1.CreateHasToGroupRequest{
		ColumnIdA:   100,
		ColumnIdB:   101,
		Description: "description-2",
		ExternalId:  "00000000-0000-4000-8000-000000000004",
	}
}

func TestValidateCreateHasToGroup(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateHasToGroupRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateHasToGroupRequest) {}},
		{name: "zero column_id_a", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ColumnIdA = 0 }, wantField: "column_id_a"},
		{name: "zero column_id_b", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ColumnIdB = 0 }, wantField: "column_id_b"},
		{name: "empty description", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateHasToGroupRequest) {
			r.Description = strings.Repeat("a", hasToGroupDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateHasToGroupRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateHasToGroupRequest()
			tt.mutate(req)

			err := ValidateCreateHasToGroup(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateHasToGroup() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateHasToGroup() = nil, want error on %q", tt.wantField)
			}

			fields := hasToGroupFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateHasToGroupNil(t *testing.T) {
	if err := ValidateCreateHasToGroup(nil); err == nil {
		t.Error("ValidateCreateHasToGroup(nil) = nil, want error")
	}
}

func validValidateUpdateHasToGroupByIdRequest() *tablesv1.UpdateHasToGroupByIdRequest {
	return &tablesv1.UpdateHasToGroupByIdRequest{
		ColumnIdA:   100,
		ColumnIdB:   101,
		Description: "description-2",
		ExternalId:  "00000000-0000-4000-8000-000000000004",
		Id:          104,
	}
}

func TestValidateUpdateHasToGroupById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateHasToGroupByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateHasToGroupByIdRequest) {}},
		{name: "zero column_id_a", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ColumnIdA = 0 }, wantField: "column_id_a"},
		{name: "zero column_id_b", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ColumnIdB = 0 }, wantField: "column_id_b"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) {
			r.Description = strings.Repeat("a", hasToGroupDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateHasToGroupByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateHasToGroupByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateHasToGroupById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateHasToGroupById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateHasToGroupById() = nil, want error on %q", tt.wantField)
			}

			fields := hasToGroupFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateHasToGroupByIdNil(t *testing.T) {
	if err := ValidateUpdateHasToGroupById(nil); err == nil {
		t.Error("ValidateUpdateHasToGroupById(nil) = nil, want error")
	}
}

func validValidateDeleteHasToGroupByIdRequest() *tablesv1.DeleteHasToGroupByIdRequest {
	return &tablesv1.DeleteHasToGroupByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteHasToGroupById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteHasToGroupByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteHasToGroupByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteHasToGroupByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteHasToGroupByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteHasToGroupByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteHasToGroupByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteHasToGroupById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteHasToGroupById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteHasToGroupById() = nil, want error on %q", tt.wantField)
			}

			fields := hasToGroupFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteHasToGroupByIdNil(t *testing.T) {
	if err := ValidateDeleteHasToGroupById(nil); err == nil {
		t.Error("ValidateDeleteHasToGroupById(nil) = nil, want error")
	}
}

func validValidateUndeleteHasToGroupByIdRequest() *tablesv1.UndeleteHasToGroupByIdRequest {
	return &tablesv1.UndeleteHasToGroupByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteHasToGroupById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteHasToGroupByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteHasToGroupByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteHasToGroupByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteHasToGroupByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteHasToGroupByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteHasToGroupByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteHasToGroupById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteHasToGroupById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteHasToGroupById() = nil, want error on %q", tt.wantField)
			}

			fields := hasToGroupFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteHasToGroupByIdNil(t *testing.T) {
	if err := ValidateUndeleteHasToGroupById(nil); err == nil {
		t.Error("ValidateUndeleteHasToGroupById(nil) = nil, want error")
	}
}

func TestValidateGetHasToGroupsSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetHasToGroupsSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetHasToGroupsSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetHasToGroupsSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetHasToGroupsSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetHasToGroupsSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetHasToGroupsSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetHasToGroupsSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetHasToGroupsSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetHasToGroupsSearchNameNil(t *testing.T) {
	if err := ValidateGetHasToGroupsSearchName(nil); err == nil {
		t.Error("ValidateGetHasToGroupsSearchName(nil) = nil, want error")
	}
}

func TestValidateGetHasToGroups(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetHasToGroupsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetHasToGroupsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetHasToGroupsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetHasToGroupsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetHasToGroupsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetHasToGroupsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetHasToGroups(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetHasToGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetHasToGroupsNil(t *testing.T) {
	if err := ValidateGetHasToGroups(nil); err == nil {
		t.Error("ValidateGetHasToGroups(nil) = nil, want error")
	}
}

func TestValidateGetHasToGroupsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetHasToGroupsDeletedRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetHasToGroupsDeletedRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetHasToGroupsDeletedRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetHasToGroupsDeletedRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetHasToGroupsDeletedRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetHasToGroupsDeletedRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetHasToGroupsDeleted(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetHasToGroupsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetHasToGroupsDeletedNil(t *testing.T) {
	if err := ValidateGetHasToGroupsDeleted(nil); err == nil {
		t.Error("ValidateGetHasToGroupsDeleted(nil) = nil, want error")
	}
}
