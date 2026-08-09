package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	tableTypeCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	tableTypeUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testTableTypeRow — строка dc.table_type со значениями, различимыми между полями.
func testTableTypeRow() tables_model.DcTableType {
	return tables_model.DcTableType{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		IsDeleted:   false,
		CreatedAt:   tableTypeCreatedAt,
		UpdatedAt:   tableTypeUpdatedAt,
		UserID:      106,
	}
}

func TestTableTypeToProto(t *testing.T) {
	row := testTableTypeRow()
	got := TableTypeToProto(row)

	if got == nil {
		t.Fatal("TableTypeToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
	}

	if got.GetDescription() != row.Description {
		t.Errorf("Description = %q, want %q", got.GetDescription(), row.Description)
	}

	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}

	if !got.GetCreatedAt().AsTime().Equal(row.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.GetCreatedAt().AsTime(), row.CreatedAt)
	}

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

}

func TestTableTypesToProto(t *testing.T) {
	first := testTableTypeRow()

	second := testTableTypeRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcTableType
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcTableType{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcTableType{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TableTypesToProto(tt.input)

			if got == nil {
				t.Fatal("TableTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateTableTypeParams(t *testing.T) {
	req := &tablesv1.CreateTableTypeRequest{
		Name:        "name-0",
		Description: "description-0",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateTableTypeParams{
		Name:        "name-0",
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateTableTypeParams(req); got != want {
		t.Errorf("ToCreateTableTypeParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateTableTypeParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateTableTypeParams(nil); got != (tables_model.CreateTableTypeParams{}) {
		t.Errorf("ToCreateTableTypeParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateTableTypeByIdParams(t *testing.T) {
	req := &tablesv1.UpdateTableTypeByIdRequest{
		Name:        "name-0",
		Description: "description-0",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
		Id:          103,
	}

	want := tables_model.UpdateTableTypeByIdParams{
		Name:        "name-0",
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000003"),
		ID:          103,
	}

	if got := ToUpdateTableTypeByIdParams(req); got != want {
		t.Errorf("ToUpdateTableTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateTableTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateTableTypeByIdParams(nil); got != (tables_model.UpdateTableTypeByIdParams{}) {
		t.Errorf("ToUpdateTableTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteTableTypeByIdParams(t *testing.T) {
	req := &tablesv1.DeleteTableTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteTableTypeByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteTableTypeByIdParams(req); got != want {
		t.Errorf("ToDeleteTableTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteTableTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteTableTypeByIdParams(nil); got != (tables_model.DeleteTableTypeByIdParams{}) {
		t.Errorf("ToDeleteTableTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteTableTypeByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteTableTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteTableTypeByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteTableTypeByIdParams(req); got != want {
		t.Errorf("ToUndeleteTableTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteTableTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteTableTypeByIdParams(nil); got != (tables_model.UndeleteTableTypeByIdParams{}) {
		t.Errorf("ToUndeleteTableTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetTableTypesDefaultsPageLimit(t *testing.T) {
	got := ToGetTableTypesParams(&tablesv1.GetTableTypesRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetTableTypesDefaultsPage(t *testing.T) {
	got := ToGetTableTypesParams(&tablesv1.GetTableTypesRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetTableTypesKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetTableTypesParams(&tablesv1.GetTableTypesRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetTableTypesDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetTableTypesDeletedParams(&tablesv1.GetTableTypesDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetTableTypesDeletedDefaultsPage(t *testing.T) {
	got := ToGetTableTypesDeletedParams(&tablesv1.GetTableTypesDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetTableTypesDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetTableTypesDeletedParams(&tablesv1.GetTableTypesDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
