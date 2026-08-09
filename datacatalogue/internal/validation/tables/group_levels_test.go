package tables

import (
	"errors"
	"math"
	"strings"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// groupLevelsFieldErrors достаёт из ошибки список полей с претензиями.
func groupLevelsFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateGroupLevelsRequest() *tablesv1.CreateGroupLevelsRequest {
	return &tablesv1.CreateGroupLevelsRequest{
		ColumnId:       100,
		ParentColumnId: 101,
		Level:          12,
		Description:    "description-3",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
	}
}

func TestValidateCreateGroupLevels(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateGroupLevelsRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateGroupLevelsRequest) {}},
		{name: "zero column_id", mutate: func(r *tablesv1.CreateGroupLevelsRequest) { r.ColumnId = 0 }, wantField: "column_id"},
		{name: "zero parent_column_id", mutate: func(r *tablesv1.CreateGroupLevelsRequest) { r.ParentColumnId = 0 }, wantField: "parent_column_id"},
		{name: "level over smallint", mutate: func(r *tablesv1.CreateGroupLevelsRequest) { r.Level = math.MaxInt16 + 1 }, wantField: "level"},
		{name: "empty description", mutate: func(r *tablesv1.CreateGroupLevelsRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.CreateGroupLevelsRequest) {
			r.Description = strings.Repeat("a", groupLevelsDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateGroupLevelsRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateGroupLevelsRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateGroupLevelsRequest()
			tt.mutate(req)

			err := ValidateCreateGroupLevels(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateGroupLevels() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateGroupLevels() = nil, want error on %q", tt.wantField)
			}

			fields := groupLevelsFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateGroupLevelsNil(t *testing.T) {
	if err := ValidateCreateGroupLevels(nil); err == nil {
		t.Error("ValidateCreateGroupLevels(nil) = nil, want error")
	}
}

func validValidateUpdateGroupLevelsByIdRequest() *tablesv1.UpdateGroupLevelsByIdRequest {
	return &tablesv1.UpdateGroupLevelsByIdRequest{
		ColumnId:       100,
		ParentColumnId: 101,
		Level:          12,
		Description:    "description-3",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
		Id:             105,
	}
}

func TestValidateUpdateGroupLevelsById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateGroupLevelsByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateGroupLevelsByIdRequest) {}},
		{name: "zero column_id", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) { r.ColumnId = 0 }, wantField: "column_id"},
		{name: "zero parent_column_id", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) { r.ParentColumnId = 0 }, wantField: "parent_column_id"},
		{name: "level over smallint", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) { r.Level = math.MaxInt16 + 1 }, wantField: "level"},
		{name: "empty description", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) { r.Description = "" }, wantField: "description"},
		{name: "description too long", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) {
			r.Description = strings.Repeat("a", groupLevelsDescriptionMaxLen+1)
		}, wantField: "description"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateGroupLevelsByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateGroupLevelsByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateGroupLevelsById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateGroupLevelsById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateGroupLevelsById() = nil, want error on %q", tt.wantField)
			}

			fields := groupLevelsFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateGroupLevelsByIdNil(t *testing.T) {
	if err := ValidateUpdateGroupLevelsById(nil); err == nil {
		t.Error("ValidateUpdateGroupLevelsById(nil) = nil, want error")
	}
}

func validValidateDeleteGroupLevelsByIdRequest() *tablesv1.DeleteGroupLevelsByIdRequest {
	return &tablesv1.DeleteGroupLevelsByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteGroupLevelsById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteGroupLevelsByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteGroupLevelsByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteGroupLevelsByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteGroupLevelsByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteGroupLevelsByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteGroupLevelsByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteGroupLevelsById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteGroupLevelsById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteGroupLevelsById() = nil, want error on %q", tt.wantField)
			}

			fields := groupLevelsFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteGroupLevelsByIdNil(t *testing.T) {
	if err := ValidateDeleteGroupLevelsById(nil); err == nil {
		t.Error("ValidateDeleteGroupLevelsById(nil) = nil, want error")
	}
}

func validValidateUndeleteGroupLevelsByIdRequest() *tablesv1.UndeleteGroupLevelsByIdRequest {
	return &tablesv1.UndeleteGroupLevelsByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateUndeleteGroupLevelsById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UndeleteGroupLevelsByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UndeleteGroupLevelsByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.UndeleteGroupLevelsByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UndeleteGroupLevelsByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UndeleteGroupLevelsByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUndeleteGroupLevelsByIdRequest()
			tt.mutate(req)

			err := ValidateUndeleteGroupLevelsById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUndeleteGroupLevelsById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUndeleteGroupLevelsById() = nil, want error on %q", tt.wantField)
			}

			fields := groupLevelsFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUndeleteGroupLevelsByIdNil(t *testing.T) {
	if err := ValidateUndeleteGroupLevelsById(nil); err == nil {
		t.Error("ValidateUndeleteGroupLevelsById(nil) = nil, want error")
	}
}

func TestValidateGetGroupLevelsDeleted(t *testing.T) {
	tests := []struct {
		name    string
		order   string
		wantErr bool
	}{
		{name: "empty order", order: ""},
		{name: "ASC", order: "ASC"},
		{name: "DESC", order: "DESC"},
		{name: "invalid order", order: "sideways", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetGroupLevelsDeleted(&tablesv1.GetGroupLevelsDeletedRequest{Order: tt.order})
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetGroupLevelsDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetGroupLevelsDeletedNil(t *testing.T) {
	if err := ValidateGetGroupLevelsDeleted(nil); err == nil {
		t.Error("ValidateGetGroupLevelsDeleted(nil) = nil, want error")
	}
}

func TestValidateGetGroupLevelsSearchName(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetGroupLevelsSearchNameRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetGroupLevelsSearchNameRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetGroupLevelsSearchNameRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetGroupLevelsSearchNameRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetGroupLevelsSearchNameRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetGroupLevelsSearchNameRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetGroupLevelsSearchName(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetGroupLevelsSearchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetGroupLevelsSearchNameNil(t *testing.T) {
	if err := ValidateGetGroupLevelsSearchName(nil); err == nil {
		t.Error("ValidateGetGroupLevelsSearchName(nil) = nil, want error")
	}
}

func TestValidateGetGroupLevels(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetGroupLevelsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetGroupLevelsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetGroupLevelsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetGroupLevelsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetGroupLevelsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetGroupLevelsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetGroupLevels(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetGroupLevels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetGroupLevelsNil(t *testing.T) {
	if err := ValidateGetGroupLevels(nil); err == nil {
		t.Error("ValidateGetGroupLevels(nil) = nil, want error")
	}
}

func TestValidateGetGroupLevelsFiltered(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetGroupLevelsFilteredRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetGroupLevelsFilteredRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetGroupLevelsFilteredRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetGroupLevelsFilteredRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetGroupLevelsFilteredRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetGroupLevelsFilteredRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetGroupLevelsFiltered(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetGroupLevelsFiltered() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetGroupLevelsFilteredNil(t *testing.T) {
	if err := ValidateGetGroupLevelsFiltered(nil); err == nil {
		t.Error("ValidateGetGroupLevelsFiltered(nil) = nil, want error")
	}
}
