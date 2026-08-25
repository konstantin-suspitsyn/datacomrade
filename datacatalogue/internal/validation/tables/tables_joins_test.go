package tables

import (
	"errors"
	"testing"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tablesJoinsFieldErrors достаёт из ошибки список полей с претензиями.
func tablesJoinsFieldErrors(t *testing.T, err error) map[string][]string {
	t.Helper()

	var validationErr *validator.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *validator.ValidationError", err)
	}

	return validationErr.Errors
}

func validValidateCreateTablesJoinsRequest() *tablesv1.CreateTablesJoinsRequest {
	return &tablesv1.CreateTablesJoinsRequest{
		TableCatIdLeft:  100,
		TableCatIdRight: 101,
		ExternalId:      "00000000-0000-4000-8000-000000000003",
	}
}

func TestValidateCreateTablesJoins(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.CreateTablesJoinsRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.CreateTablesJoinsRequest) {}},
		{name: "zero table_cat_id_left", mutate: func(r *tablesv1.CreateTablesJoinsRequest) { r.TableCatIdLeft = 0 }, wantField: "table_cat_id_left"},
		{name: "zero table_cat_id_right", mutate: func(r *tablesv1.CreateTablesJoinsRequest) { r.TableCatIdRight = 0 }, wantField: "table_cat_id_right"},
		{name: "empty user_id", mutate: func(r *tablesv1.CreateTablesJoinsRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.CreateTablesJoinsRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateCreateTablesJoinsRequest()
			tt.mutate(req)

			err := ValidateCreateTablesJoins(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateCreateTablesJoins() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateCreateTablesJoins() = nil, want error on %q", tt.wantField)
			}

			fields := tablesJoinsFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateCreateTablesJoinsNil(t *testing.T) {
	if err := ValidateCreateTablesJoins(nil); err == nil {
		t.Error("ValidateCreateTablesJoins(nil) = nil, want error")
	}
}

func validValidateUpdateTablesJoinsByIdRequest() *tablesv1.UpdateTablesJoinsByIdRequest {
	return &tablesv1.UpdateTablesJoinsByIdRequest{
		TableCatIdLeft:  100,
		TableCatIdRight: 101,
		IsDeleted:       true,
		ExternalId:      "00000000-0000-4000-8000-000000000004",
		Id:              104,
	}
}

func TestValidateUpdateTablesJoinsById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.UpdateTablesJoinsByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.UpdateTablesJoinsByIdRequest) {}},
		{name: "zero table_cat_id_left", mutate: func(r *tablesv1.UpdateTablesJoinsByIdRequest) { r.TableCatIdLeft = 0 }, wantField: "table_cat_id_left"},
		{name: "zero table_cat_id_right", mutate: func(r *tablesv1.UpdateTablesJoinsByIdRequest) { r.TableCatIdRight = 0 }, wantField: "table_cat_id_right"},
		{name: "empty user_id", mutate: func(r *tablesv1.UpdateTablesJoinsByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.UpdateTablesJoinsByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.UpdateTablesJoinsByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateUpdateTablesJoinsByIdRequest()
			tt.mutate(req)

			err := ValidateUpdateTablesJoinsById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateUpdateTablesJoinsById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateUpdateTablesJoinsById() = nil, want error on %q", tt.wantField)
			}

			fields := tablesJoinsFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateUpdateTablesJoinsByIdNil(t *testing.T) {
	if err := ValidateUpdateTablesJoinsById(nil); err == nil {
		t.Error("ValidateUpdateTablesJoinsById(nil) = nil, want error")
	}
}

func validValidateDeleteTablesJoinsByIdRequest() *tablesv1.DeleteTablesJoinsByIdRequest {
	return &tablesv1.DeleteTablesJoinsByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}
}

func TestValidateDeleteTablesJoinsById(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*tablesv1.DeleteTablesJoinsByIdRequest)
		wantField string
	}{
		{name: "valid", mutate: func(*tablesv1.DeleteTablesJoinsByIdRequest) {}},
		{name: "empty user_id", mutate: func(r *tablesv1.DeleteTablesJoinsByIdRequest) { r.ExternalId = "" }, wantField: "user_id"},
		{name: "malformed user_id", mutate: func(r *tablesv1.DeleteTablesJoinsByIdRequest) { r.ExternalId = "not-a-uuid" }, wantField: "user_id"},
		{name: "zero id", mutate: func(r *tablesv1.DeleteTablesJoinsByIdRequest) { r.Id = 0 }, wantField: "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validValidateDeleteTablesJoinsByIdRequest()
			tt.mutate(req)

			err := ValidateDeleteTablesJoinsById(req)

			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("ValidateDeleteTablesJoinsById() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateDeleteTablesJoinsById() = nil, want error on %q", tt.wantField)
			}

			fields := tablesJoinsFieldErrors(t, err)
			if len(fields[tt.wantField]) == 0 {
				t.Errorf("no error on %q, got %v", tt.wantField, fields)
			}
		})
	}
}

func TestValidateDeleteTablesJoinsByIdNil(t *testing.T) {
	if err := ValidateDeleteTablesJoinsById(nil); err == nil {
		t.Error("ValidateDeleteTablesJoinsById(nil) = nil, want error")
	}
}

func TestValidateGetTablesJoins(t *testing.T) {
	tests := []struct {
		name    string
		req     *tablesv1.GetTablesJoinsRequest
		wantErr bool
	}{
		{name: "valid", req: &tablesv1.GetTablesJoinsRequest{PageLimit: 50, Page: 1}},
		{name: "zero page_limit and page ok", req: &tablesv1.GetTablesJoinsRequest{PageLimit: 0, Page: 0}},
		{name: "negative page_limit", req: &tablesv1.GetTablesJoinsRequest{PageLimit: -1}, wantErr: true},
		{name: "negative page", req: &tablesv1.GetTablesJoinsRequest{Page: -1}, wantErr: true},
		{name: "invalid order", req: &tablesv1.GetTablesJoinsRequest{Order: "sideways"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGetTablesJoins(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGetTablesJoins() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGetTablesJoinsNil(t *testing.T) {
	if err := ValidateGetTablesJoins(nil); err == nil {
		t.Error("ValidateGetTablesJoins(nil) = nil, want error")
	}
}
