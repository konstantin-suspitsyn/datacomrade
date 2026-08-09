package tables

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// columnCatPtr — адрес значения для optional-полей proto.
func columnCatPtr[T any](v T) *T {
	return &v
}

var (
	columnCatCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	columnCatUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testColumnCatRow — строка dc.column_cat со значениями, различимыми между полями.
func testColumnCatRow() tables_model.DcColumnCat {
	return tables_model.DcColumnCat{
		ID:                100,
		TableID:           101,
		Name:              "name-0",
		AliasID:           103,
		ColumnTypeID:      104,
		Description:       "description-0",
		CalculationTypeID: 106,
		IsDeleted:         false,
		ShowInUi:          true,
		CreatedAt:         columnCatCreatedAt,
		UpdatedAt:         columnCatUpdatedAt,
		UserID:            111,
	}
}

func TestColumnCatToProto(t *testing.T) {
	row := testColumnCatRow()
	got := ColumnCatToProto(row)

	if got == nil {
		t.Fatal("ColumnCatToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetTableId() != row.TableID {
		t.Errorf("TableId = %d, want %d", got.GetTableId(), row.TableID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
	}

	if got.GetAliasId() != row.AliasID {
		t.Errorf("AliasId = %d, want %d", got.GetAliasId(), row.AliasID)
	}

	if got.GetColumnTypeId() != row.ColumnTypeID {
		t.Errorf("ColumnTypeId = %d, want %d", got.GetColumnTypeId(), row.ColumnTypeID)
	}

	if got.GetDescription() != row.Description {
		t.Errorf("Description = %q, want %q", got.GetDescription(), row.Description)
	}

	if got.GetCalculationTypeId() != row.CalculationTypeID {
		t.Errorf("CalculationTypeId = %d, want %d", got.GetCalculationTypeId(), row.CalculationTypeID)
	}

	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}

	if got.GetShowInUi() != row.ShowInUi {
		t.Errorf("ShowInUi = %v, want %v", got.GetShowInUi(), row.ShowInUi)
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

func TestColumnCatsToProto(t *testing.T) {
	first := testColumnCatRow()

	second := testColumnCatRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcColumnCat
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcColumnCat{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcColumnCat{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColumnCatsToProto(tt.input)

			if got == nil {
				t.Fatal("ColumnCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateColumnCatParams(t *testing.T) {
	req := &tablesv1.CreateColumnCatRequest{
		TableId:           100,
		Name:              "name-0",
		AliasId:           102,
		ColumnTypeId:      103,
		Description:       "description-0",
		CalculationTypeId: 105,
		ShowInUi:          true,
		ExternalId:        "00000000-0000-4000-8000-000000000008",
	}

	want := tables_model.CreateColumnCatParams{
		TableID:           100,
		Name:              "name-0",
		AliasID:           102,
		ColumnTypeID:      103,
		Description:       "description-0",
		CalculationTypeID: 105,
		ShowInUi:          true,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000008"),
	}

	if got := ToCreateColumnCatParams(req); got != want {
		t.Errorf("ToCreateColumnCatParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateColumnCatParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateColumnCatParams(nil); got != (tables_model.CreateColumnCatParams{}) {
		t.Errorf("ToCreateColumnCatParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateColumnCatByIdParams(t *testing.T) {
	req := &tablesv1.UpdateColumnCatByIdRequest{
		TableId:           100,
		Name:              "name-0",
		AliasId:           102,
		ColumnTypeId:      103,
		Description:       "description-0",
		CalculationTypeId: 105,
		ShowInUi:          true,
		ExternalId:        "00000000-0000-4000-8000-000000000008",
		Id:                108,
	}

	want := tables_model.UpdateColumnCatByIdParams{
		TableID:           100,
		Name:              "name-0",
		AliasID:           102,
		ColumnTypeID:      103,
		Description:       "description-0",
		CalculationTypeID: 105,
		ShowInUi:          true,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000008"),
		ID:                108,
	}

	if got := ToUpdateColumnCatByIdParams(req); got != want {
		t.Errorf("ToUpdateColumnCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateColumnCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateColumnCatByIdParams(nil); got != (tables_model.UpdateColumnCatByIdParams{}) {
		t.Errorf("ToUpdateColumnCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteColumnCatByIdParams(t *testing.T) {
	req := &tablesv1.DeleteColumnCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteColumnCatByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteColumnCatByIdParams(req); got != want {
		t.Errorf("ToDeleteColumnCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteColumnCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteColumnCatByIdParams(nil); got != (tables_model.DeleteColumnCatByIdParams{}) {
		t.Errorf("ToDeleteColumnCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteColumnCatByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteColumnCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteColumnCatByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteColumnCatByIdParams(req); got != want {
		t.Errorf("ToUndeleteColumnCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteColumnCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteColumnCatByIdParams(nil); got != (tables_model.UndeleteColumnCatByIdParams{}) {
		t.Errorf("ToUndeleteColumnCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetColumnCatsDefaultsPageLimit(t *testing.T) {
	got := ToGetColumnCatsParams(&tablesv1.GetColumnCatsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetColumnCatsDefaultsPage(t *testing.T) {
	got := ToGetColumnCatsParams(&tablesv1.GetColumnCatsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetColumnCatsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetColumnCatsParams(&tablesv1.GetColumnCatsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetColumnCatsDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetColumnCatsDeletedParams(&tablesv1.GetColumnCatsDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetColumnCatsDeletedDefaultsPage(t *testing.T) {
	got := ToGetColumnCatsDeletedParams(&tablesv1.GetColumnCatsDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetColumnCatsDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetColumnCatsDeletedParams(&tablesv1.GetColumnCatsDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetColumnCatsOptionalFiltersDefaultsPageLimit(t *testing.T) {
	got := ToGetColumnCatsOptionalFiltersParams(&tablesv1.GetColumnCatsOptionalFiltersRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetColumnCatsOptionalFiltersDefaultsPage(t *testing.T) {
	got := ToGetColumnCatsOptionalFiltersParams(&tablesv1.GetColumnCatsOptionalFiltersRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetColumnCatsOptionalFiltersKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetColumnCatsOptionalFiltersParams(&tablesv1.GetColumnCatsOptionalFiltersRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetColumnCatsOptionalFiltersPassesFilterFields(t *testing.T) {
	req := &tablesv1.GetColumnCatsOptionalFiltersRequest{PageLimit: 10, Page: 1,
		TableId:           columnCatPtr(int64(100)),
		AliasId:           columnCatPtr(int64(101)),
		ColumnTypeId:      columnCatPtr(int64(102)),
		CalculationTypeId: columnCatPtr(int64(103)),
		ShowInUi:          columnCatPtr(true),
	}

	got := ToGetColumnCatsOptionalFiltersParams(req)

	if got.TableID != (sql.NullInt64{Int64: 100, Valid: true}) {
		t.Errorf("TableID = %v, want %v", got.TableID, (sql.NullInt64{Int64: 100, Valid: true}))
	}
	if got.AliasID != (sql.NullInt64{Int64: 101, Valid: true}) {
		t.Errorf("AliasID = %v, want %v", got.AliasID, (sql.NullInt64{Int64: 101, Valid: true}))
	}
	if got.ColumnTypeID != (sql.NullInt64{Int64: 102, Valid: true}) {
		t.Errorf("ColumnTypeID = %v, want %v", got.ColumnTypeID, (sql.NullInt64{Int64: 102, Valid: true}))
	}
	if got.CalculationTypeID != (sql.NullInt64{Int64: 103, Valid: true}) {
		t.Errorf("CalculationTypeID = %v, want %v", got.CalculationTypeID, (sql.NullInt64{Int64: 103, Valid: true}))
	}
	if got.ShowInUi != (sql.NullBool{Bool: true, Valid: true}) {
		t.Errorf("ShowInUi = %v, want %v", got.ShowInUi, (sql.NullBool{Bool: true, Valid: true}))
	}
}

// Непереданный необязательный фильтр обязан остаться SQL NULL, а не
// превратиться в 0/false: иначе выборка молча сузилась бы.
func TestGetColumnCatsOptionalFiltersOmitsAbsentFilters(t *testing.T) {
	got := ToGetColumnCatsOptionalFiltersParams(&tablesv1.GetColumnCatsOptionalFiltersRequest{PageLimit: 10, Page: 1})

	if got.TableID.Valid {
		t.Errorf("TableID = %+v, want invalid", got.TableID)
	}
	if got.AliasID.Valid {
		t.Errorf("AliasID = %+v, want invalid", got.AliasID)
	}
	if got.ColumnTypeID.Valid {
		t.Errorf("ColumnTypeID = %+v, want invalid", got.ColumnTypeID)
	}
	if got.CalculationTypeID.Valid {
		t.Errorf("CalculationTypeID = %+v, want invalid", got.CalculationTypeID)
	}
	if got.ShowInUi.Valid {
		t.Errorf("ShowInUi = %+v, want invalid", got.ShowInUi)
	}
}

func TestGetColumnCatsSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetColumnCatsSearchNameParams(&tablesv1.GetColumnCatsSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetColumnCatsSearchNameDefaultsPage(t *testing.T) {
	got := ToGetColumnCatsSearchNameParams(&tablesv1.GetColumnCatsSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetColumnCatsSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetColumnCatsSearchNameParams(&tablesv1.GetColumnCatsSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetColumnCatsSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetColumnCatsSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetColumnCatsSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}
