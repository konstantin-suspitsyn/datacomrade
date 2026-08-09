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

// groupLevelsPtr — адрес значения для optional-полей proto.
func groupLevelsPtr[T any](v T) *T {
	return &v
}

var (
	groupLevelsCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	groupLevelsUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testGroupLevelsRow — строка dc.group_levels со значениями, различимыми между полями.
func testGroupLevelsRow() tables_model.DcGroupLevel {
	return tables_model.DcGroupLevel{
		ID:             100,
		ColumnID:       101,
		ParentColumnID: 102,
		Level:          13,
		Description:    "description-0",
		CreatedAt:      groupLevelsCreatedAt,
		UpdatedAt:      groupLevelsUpdatedAt,
		IsDeleted:      false,
		UserID:         108,
	}
}

func TestGroupLevelsToProto(t *testing.T) {
	row := testGroupLevelsRow()
	got := GroupLevelsToProto(row)

	if got == nil {
		t.Fatal("GroupLevelsToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetColumnId() != row.ColumnID {
		t.Errorf("ColumnId = %d, want %d", got.GetColumnId(), row.ColumnID)
	}

	if got.GetParentColumnId() != row.ParentColumnID {
		t.Errorf("ParentColumnId = %d, want %d", got.GetParentColumnId(), row.ParentColumnID)
	}

	if got.GetLevel() != int32(row.Level) {
		t.Errorf("Level = %d, want %d", got.GetLevel(), row.Level)
	}

	if got.GetDescription() != row.Description {
		t.Errorf("Description = %q, want %q", got.GetDescription(), row.Description)
	}

	if !got.GetCreatedAt().AsTime().Equal(row.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.GetCreatedAt().AsTime(), row.CreatedAt)
	}

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt)
	}

	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

}

func TestGroupLevelsesToProto(t *testing.T) {
	first := testGroupLevelsRow()

	second := testGroupLevelsRow()
	second.ID = 999
	second.Description = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcGroupLevel
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcGroupLevel{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcGroupLevel{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupLevelsesToProto(tt.input)

			if got == nil {
				t.Fatal("GroupLevelsesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateGroupLevelsParams(t *testing.T) {
	req := &tablesv1.CreateGroupLevelsRequest{
		ColumnId:       100,
		ParentColumnId: 101,
		Level:          12,
		Description:    "description-0",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
	}

	want := tables_model.CreateGroupLevelsParams{
		ColumnID:       100,
		ParentColumnID: 101,
		Level:          12,
		Description:    "description-0",
		ExternalID:     uuid.MustParse("00000000-0000-4000-8000-000000000005"),
	}

	if got := ToCreateGroupLevelsParams(req); got != want {
		t.Errorf("ToCreateGroupLevelsParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateGroupLevelsParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateGroupLevelsParams(nil); got != (tables_model.CreateGroupLevelsParams{}) {
		t.Errorf("ToCreateGroupLevelsParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateGroupLevelsByIdParams(t *testing.T) {
	req := &tablesv1.UpdateGroupLevelsByIdRequest{
		ColumnId:       100,
		ParentColumnId: 101,
		Level:          12,
		Description:    "description-0",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
		Id:             105,
	}

	want := tables_model.UpdateGroupLevelsByIdParams{
		ColumnID:       100,
		ParentColumnID: 101,
		Level:          12,
		Description:    "description-0",
		ExternalID:     uuid.MustParse("00000000-0000-4000-8000-000000000005"),
		ID:             105,
	}

	if got := ToUpdateGroupLevelsByIdParams(req); got != want {
		t.Errorf("ToUpdateGroupLevelsByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateGroupLevelsByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateGroupLevelsByIdParams(nil); got != (tables_model.UpdateGroupLevelsByIdParams{}) {
		t.Errorf("ToUpdateGroupLevelsByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteGroupLevelsByIdParams(t *testing.T) {
	req := &tablesv1.DeleteGroupLevelsByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteGroupLevelsByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteGroupLevelsByIdParams(req); got != want {
		t.Errorf("ToDeleteGroupLevelsByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteGroupLevelsByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteGroupLevelsByIdParams(nil); got != (tables_model.DeleteGroupLevelsByIdParams{}) {
		t.Errorf("ToDeleteGroupLevelsByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteGroupLevelsByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteGroupLevelsByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteGroupLevelsByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteGroupLevelsByIdParams(req); got != want {
		t.Errorf("ToUndeleteGroupLevelsByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteGroupLevelsByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteGroupLevelsByIdParams(nil); got != (tables_model.UndeleteGroupLevelsByIdParams{}) {
		t.Errorf("ToUndeleteGroupLevelsByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetGroupLevelsSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetGroupLevelsSearchNameParams(&tablesv1.GetGroupLevelsSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetGroupLevelsSearchNameDefaultsPage(t *testing.T) {
	got := ToGetGroupLevelsSearchNameParams(&tablesv1.GetGroupLevelsSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetGroupLevelsSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetGroupLevelsSearchNameParams(&tablesv1.GetGroupLevelsSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetGroupLevelsSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetGroupLevelsSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetGroupLevelsSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}

func TestGetGroupLevelsDefaultsPageLimit(t *testing.T) {
	got := ToGetGroupLevelsParams(&tablesv1.GetGroupLevelsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetGroupLevelsDefaultsPage(t *testing.T) {
	got := ToGetGroupLevelsParams(&tablesv1.GetGroupLevelsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetGroupLevelsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetGroupLevelsParams(&tablesv1.GetGroupLevelsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetGroupLevelsFilteredDefaultsPageLimit(t *testing.T) {
	got := ToGetGroupLevelsFilteredParams(&tablesv1.GetGroupLevelsFilteredRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetGroupLevelsFilteredDefaultsPage(t *testing.T) {
	got := ToGetGroupLevelsFilteredParams(&tablesv1.GetGroupLevelsFilteredRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetGroupLevelsFilteredKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetGroupLevelsFilteredParams(&tablesv1.GetGroupLevelsFilteredRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetGroupLevelsFilteredPassesFilterFields(t *testing.T) {
	req := &tablesv1.GetGroupLevelsFilteredRequest{PageLimit: 10, Page: 1,
		ColumnId:       groupLevelsPtr(int64(100)),
		ParentColumnId: groupLevelsPtr(int64(101)),
	}

	got := ToGetGroupLevelsFilteredParams(req)

	if got.ColumnID != (sql.NullInt64{Int64: 100, Valid: true}) {
		t.Errorf("ColumnID = %v, want %v", got.ColumnID, (sql.NullInt64{Int64: 100, Valid: true}))
	}
	if got.ParentColumnID != (sql.NullInt64{Int64: 101, Valid: true}) {
		t.Errorf("ParentColumnID = %v, want %v", got.ParentColumnID, (sql.NullInt64{Int64: 101, Valid: true}))
	}
}

// Непереданный необязательный фильтр обязан остаться SQL NULL, а не
// превратиться в 0/false: иначе выборка молча сузилась бы.
func TestGetGroupLevelsFilteredOmitsAbsentFilters(t *testing.T) {
	got := ToGetGroupLevelsFilteredParams(&tablesv1.GetGroupLevelsFilteredRequest{PageLimit: 10, Page: 1})

	if got.ColumnID.Valid {
		t.Errorf("ColumnID = %+v, want invalid", got.ColumnID)
	}
	if got.ParentColumnID.Valid {
		t.Errorf("ParentColumnID = %+v, want invalid", got.ParentColumnID)
	}
}
