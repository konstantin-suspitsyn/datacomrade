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
	columnTypeCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	columnTypeUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testColumnTypeRow — строка dc.column_type со значениями, различимыми между полями.
func testColumnTypeRow() tables_model.DcColumnType {
	return tables_model.DcColumnType{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		IsDeleted:   false,
		CreatedAt:   columnTypeCreatedAt,
		UpdatedAt:   columnTypeUpdatedAt,
		UserID:      106,
	}
}

func TestColumnTypeToProto(t *testing.T) {
	row := testColumnTypeRow()
	got := ColumnTypeToProto(row)

	if got == nil {
		t.Fatal("ColumnTypeToProto() = nil, want value")
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

func TestColumnTypesToProto(t *testing.T) {
	first := testColumnTypeRow()

	second := testColumnTypeRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcColumnType
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcColumnType{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcColumnType{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColumnTypesToProto(tt.input)

			if got == nil {
				t.Fatal("ColumnTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateColumnTypeParams(t *testing.T) {
	req := &tablesv1.CreateColumnTypeRequest{
		Name:        "name-0",
		Description: "description-0",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateColumnTypeParams{
		Name:        "name-0",
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateColumnTypeParams(req); got != want {
		t.Errorf("ToCreateColumnTypeParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateColumnTypeParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateColumnTypeParams(nil); got != (tables_model.CreateColumnTypeParams{}) {
		t.Errorf("ToCreateColumnTypeParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateColumnTypeByIdParams(t *testing.T) {
	req := &tablesv1.UpdateColumnTypeByIdRequest{
		Name:        "name-0",
		Description: "description-0",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
		Id:          103,
	}

	want := tables_model.UpdateColumnTypeByIdParams{
		Name:        "name-0",
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000003"),
		ID:          103,
	}

	if got := ToUpdateColumnTypeByIdParams(req); got != want {
		t.Errorf("ToUpdateColumnTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateColumnTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateColumnTypeByIdParams(nil); got != (tables_model.UpdateColumnTypeByIdParams{}) {
		t.Errorf("ToUpdateColumnTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteColumnTypeByIdParams(t *testing.T) {
	req := &tablesv1.DeleteColumnTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteColumnTypeByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteColumnTypeByIdParams(req); got != want {
		t.Errorf("ToDeleteColumnTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteColumnTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteColumnTypeByIdParams(nil); got != (tables_model.DeleteColumnTypeByIdParams{}) {
		t.Errorf("ToDeleteColumnTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteColumnTypeByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteColumnTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteColumnTypeByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteColumnTypeByIdParams(req); got != want {
		t.Errorf("ToUndeleteColumnTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteColumnTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteColumnTypeByIdParams(nil); got != (tables_model.UndeleteColumnTypeByIdParams{}) {
		t.Errorf("ToUndeleteColumnTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetColumnTypesDefaultsPageLimit(t *testing.T) {
	got := ToGetColumnTypesParams(&tablesv1.GetColumnTypesRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetColumnTypesDefaultsPage(t *testing.T) {
	got := ToGetColumnTypesParams(&tablesv1.GetColumnTypesRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetColumnTypesKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetColumnTypesParams(&tablesv1.GetColumnTypesRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetColumnTypesSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetColumnTypesSearchNameParams(&tablesv1.GetColumnTypesSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetColumnTypesSearchNameDefaultsPage(t *testing.T) {
	got := ToGetColumnTypesSearchNameParams(&tablesv1.GetColumnTypesSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetColumnTypesSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetColumnTypesSearchNameParams(&tablesv1.GetColumnTypesSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetColumnTypesSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetColumnTypesSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetColumnTypesSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}

func TestGetColumnTypesDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetColumnTypesDeletedParams(&tablesv1.GetColumnTypesDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetColumnTypesDeletedDefaultsPage(t *testing.T) {
	got := ToGetColumnTypesDeletedParams(&tablesv1.GetColumnTypesDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetColumnTypesDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetColumnTypesDeletedParams(&tablesv1.GetColumnTypesDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
