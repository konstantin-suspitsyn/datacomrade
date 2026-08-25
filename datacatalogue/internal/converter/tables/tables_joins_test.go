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

// tablesJoinsPtr — адрес значения для optional-полей proto.
func tablesJoinsPtr[T any](v T) *T {
	return &v
}

var (
	tablesJoinsCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	tablesJoinsUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testTablesJoinsRow — строка dc.tables_joins со значениями, различимыми между полями.
func testTablesJoinsRow() tables_model.DcTablesJoin {
	return tables_model.DcTablesJoin{
		ID:              100,
		TableCatIDLeft:  101,
		TableCatIDRight: 102,
		IsDeleted:       false,
		CreatedAt:       tablesJoinsCreatedAt,
		UpdatedAt:       tablesJoinsUpdatedAt,
		UserID:          106,
	}
}

func TestTablesJoinsToProto(t *testing.T) {
	row := testTablesJoinsRow()
	got := TablesJoinsToProto(row)

	if got == nil {
		t.Fatal("TablesJoinsToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetTableCatIdLeft() != row.TableCatIDLeft {
		t.Errorf("TableCatIdLeft = %d, want %d", got.GetTableCatIdLeft(), row.TableCatIDLeft)
	}

	if got.GetTableCatIdRight() != row.TableCatIDRight {
		t.Errorf("TableCatIdRight = %d, want %d", got.GetTableCatIdRight(), row.TableCatIDRight)
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

func TestTablesJoinsesToProto(t *testing.T) {
	first := testTablesJoinsRow()

	second := testTablesJoinsRow()
	second.ID = 999

	tests := []struct {
		name    string
		input   []tables_model.DcTablesJoin
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcTablesJoin{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcTablesJoin{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TablesJoinsesToProto(tt.input)

			if got == nil {
				t.Fatal("TablesJoinsesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// testGetTablesJoinsByIdRow — строка выборки GetTablesJoinsByIdRow со значениями, различимыми между полями.
func testGetTablesJoinsByIdRow() tables_model.GetTablesJoinsByIdRow {
	return tables_model.GetTablesJoinsByIdRow{
		ID:                100,
		TableCatIDLeft:    101,
		TableCatIDRight:   102,
		IsDeleted:         false,
		CreatedAt:         tablesJoinsCreatedAt,
		UpdatedAt:         tablesJoinsUpdatedAt,
		UserID:            106,
		TableCatLeftName:  "table-cat-left-name-0",
		TableCatRightName: "table-cat-right-name-0",
	}
}

func TestGetTablesJoinsByIdRowToProto(t *testing.T) {
	row := testGetTablesJoinsByIdRow()
	got := GetTablesJoinsByIdRowToProto(row)

	if got == nil {
		t.Fatal("GetTablesJoinsByIdRowToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetTableCatIdLeft() != row.TableCatIDLeft {
		t.Errorf("TableCatIdLeft = %d, want %d", got.GetTableCatIdLeft(), row.TableCatIDLeft)
	}

	if got.GetTableCatIdRight() != row.TableCatIDRight {
		t.Errorf("TableCatIdRight = %d, want %d", got.GetTableCatIdRight(), row.TableCatIDRight)
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

	if got.GetTableCatLeftName() != row.TableCatLeftName {
		t.Errorf("TableCatLeftName = %q, want %q", got.GetTableCatLeftName(), row.TableCatLeftName)
	}

	if got.GetTableCatRightName() != row.TableCatRightName {
		t.Errorf("TableCatRightName = %q, want %q", got.GetTableCatRightName(), row.TableCatRightName)
	}

}

// testGetTablesJoinsRow — строка выборки GetTablesJoinsRow со значениями, различимыми между полями.
func testGetTablesJoinsRow() tables_model.GetTablesJoinsRow {
	return tables_model.GetTablesJoinsRow{
		ID:                100,
		TableCatIDLeft:    101,
		TableCatIDRight:   102,
		IsDeleted:         false,
		CreatedAt:         tablesJoinsCreatedAt,
		UpdatedAt:         tablesJoinsUpdatedAt,
		UserID:            106,
		TableCatLeftName:  "table-cat-left-name-0",
		TableCatRightName: "table-cat-right-name-0",
	}
}

func TestGetTablesJoinsRowToProto(t *testing.T) {
	row := testGetTablesJoinsRow()
	got := GetTablesJoinsRowToProto(row)

	if got == nil {
		t.Fatal("GetTablesJoinsRowToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetTableCatIdLeft() != row.TableCatIDLeft {
		t.Errorf("TableCatIdLeft = %d, want %d", got.GetTableCatIdLeft(), row.TableCatIDLeft)
	}

	if got.GetTableCatIdRight() != row.TableCatIDRight {
		t.Errorf("TableCatIdRight = %d, want %d", got.GetTableCatIdRight(), row.TableCatIDRight)
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

	if got.GetTableCatLeftName() != row.TableCatLeftName {
		t.Errorf("TableCatLeftName = %q, want %q", got.GetTableCatLeftName(), row.TableCatLeftName)
	}

	if got.GetTableCatRightName() != row.TableCatRightName {
		t.Errorf("TableCatRightName = %q, want %q", got.GetTableCatRightName(), row.TableCatRightName)
	}

}

func TestGetTablesJoinsRowsToProto(t *testing.T) {
	first := testGetTablesJoinsRow()

	second := testGetTablesJoinsRow()
	second.ID = 999
	second.TableCatLeftName = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.GetTablesJoinsRow
		wantLen int
	}{
		{name: "two rows", input: []tables_model.GetTablesJoinsRow{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.GetTablesJoinsRow{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTablesJoinsRowsToProto(tt.input)

			if got == nil {
				t.Fatal("GetTablesJoinsRowsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateTablesJoinsParams(t *testing.T) {
	req := &tablesv1.CreateTablesJoinsRequest{
		TableCatIdLeft:  100,
		TableCatIdRight: 101,
		ExternalId:      "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateTablesJoinsParams{
		TableCatIDLeft:  100,
		TableCatIDRight: 101,
		ExternalID:      uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateTablesJoinsParams(req); got != want {
		t.Errorf("ToCreateTablesJoinsParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateTablesJoinsParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateTablesJoinsParams(nil); got != (tables_model.CreateTablesJoinsParams{}) {
		t.Errorf("ToCreateTablesJoinsParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateTablesJoinsByIdParams(t *testing.T) {
	req := &tablesv1.UpdateTablesJoinsByIdRequest{
		TableCatIdLeft:  100,
		TableCatIdRight: 101,
		IsDeleted:       true,
		ExternalId:      "00000000-0000-4000-8000-000000000004",
		Id:              104,
	}

	want := tables_model.UpdateTablesJoinsByIdParams{
		TableCatIDLeft:  100,
		TableCatIDRight: 101,
		IsDeleted:       true,
		ExternalID:      uuid.MustParse("00000000-0000-4000-8000-000000000004"),
		ID:              104,
	}

	if got := ToUpdateTablesJoinsByIdParams(req); got != want {
		t.Errorf("ToUpdateTablesJoinsByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateTablesJoinsByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateTablesJoinsByIdParams(nil); got != (tables_model.UpdateTablesJoinsByIdParams{}) {
		t.Errorf("ToUpdateTablesJoinsByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteTablesJoinsByIdParams(t *testing.T) {
	req := &tablesv1.DeleteTablesJoinsByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteTablesJoinsByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteTablesJoinsByIdParams(req); got != want {
		t.Errorf("ToDeleteTablesJoinsByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteTablesJoinsByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteTablesJoinsByIdParams(nil); got != (tables_model.DeleteTablesJoinsByIdParams{}) {
		t.Errorf("ToDeleteTablesJoinsByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetTablesJoinsDefaultsPageLimit(t *testing.T) {
	got := ToGetTablesJoinsParams(&tablesv1.GetTablesJoinsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetTablesJoinsDefaultsPage(t *testing.T) {
	got := ToGetTablesJoinsParams(&tablesv1.GetTablesJoinsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetTablesJoinsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetTablesJoinsParams(&tablesv1.GetTablesJoinsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetTablesJoinsPassesFilterFields(t *testing.T) {
	req := &tablesv1.GetTablesJoinsRequest{PageLimit: 10, Page: 1,
		TableCatLeftName:  tablesJoinsPtr("tablecatleftname-0"),
		TableCatRightName: tablesJoinsPtr("tablecatrightname-1"),
	}

	got := ToGetTablesJoinsParams(req)

	if got.TableCatLeftName != (sql.NullString{String: "tablecatleftname-0", Valid: true}) {
		t.Errorf("TableCatLeftName = %v, want %v", got.TableCatLeftName, (sql.NullString{String: "tablecatleftname-0", Valid: true}))
	}
	if got.TableCatRightName != (sql.NullString{String: "tablecatrightname-1", Valid: true}) {
		t.Errorf("TableCatRightName = %v, want %v", got.TableCatRightName, (sql.NullString{String: "tablecatrightname-1", Valid: true}))
	}
}

// Непереданный необязательный фильтр обязан остаться SQL NULL, а не
// превратиться в 0/false: иначе выборка молча сузилась бы.
func TestGetTablesJoinsOmitsAbsentFilters(t *testing.T) {
	got := ToGetTablesJoinsParams(&tablesv1.GetTablesJoinsRequest{PageLimit: 10, Page: 1})

	if got.TableCatLeftName.Valid {
		t.Errorf("TableCatLeftName = %+v, want invalid", got.TableCatLeftName)
	}
	if got.TableCatRightName.Valid {
		t.Errorf("TableCatRightName = %+v, want invalid", got.TableCatRightName)
	}
}
