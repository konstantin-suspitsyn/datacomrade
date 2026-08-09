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
	hasToGroupCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	hasToGroupUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testHasToGroupRow — строка dc.has_to_group со значениями, различимыми между полями.
func testHasToGroupRow() tables_model.DcHasToGroup {
	return tables_model.DcHasToGroup{
		ID:          100,
		ColumnIDA:   101,
		ColumnIDB:   102,
		Description: "description-0",
		IsDeleted:   false,
		CreatedAt:   hasToGroupCreatedAt,
		UpdatedAt:   hasToGroupUpdatedAt,
		UserID:      107,
	}
}

func TestHasToGroupToProto(t *testing.T) {
	row := testHasToGroupRow()
	got := HasToGroupToProto(row)

	if got == nil {
		t.Fatal("HasToGroupToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetColumnIdA() != row.ColumnIDA {
		t.Errorf("ColumnIdA = %d, want %d", got.GetColumnIdA(), row.ColumnIDA)
	}

	if got.GetColumnIdB() != row.ColumnIDB {
		t.Errorf("ColumnIdB = %d, want %d", got.GetColumnIdB(), row.ColumnIDB)
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

func TestHasToGroupsToProto(t *testing.T) {
	first := testHasToGroupRow()

	second := testHasToGroupRow()
	second.ID = 999
	second.Description = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcHasToGroup
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcHasToGroup{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcHasToGroup{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasToGroupsToProto(tt.input)

			if got == nil {
				t.Fatal("HasToGroupsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateHasToGroupParams(t *testing.T) {
	req := &tablesv1.CreateHasToGroupRequest{
		ColumnIdA:   100,
		ColumnIdB:   101,
		Description: "description-0",
		ExternalId:  "00000000-0000-4000-8000-000000000004",
	}

	want := tables_model.CreateHasToGroupParams{
		ColumnIDA:   100,
		ColumnIDB:   101,
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToCreateHasToGroupParams(req); got != want {
		t.Errorf("ToCreateHasToGroupParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateHasToGroupParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateHasToGroupParams(nil); got != (tables_model.CreateHasToGroupParams{}) {
		t.Errorf("ToCreateHasToGroupParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateHasToGroupByIdParams(t *testing.T) {
	req := &tablesv1.UpdateHasToGroupByIdRequest{
		ColumnIdA:   100,
		ColumnIdB:   101,
		Description: "description-0",
		ExternalId:  "00000000-0000-4000-8000-000000000004",
		Id:          104,
	}

	want := tables_model.UpdateHasToGroupByIdParams{
		ColumnIDA:   100,
		ColumnIDB:   101,
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000004"),
		ID:          104,
	}

	if got := ToUpdateHasToGroupByIdParams(req); got != want {
		t.Errorf("ToUpdateHasToGroupByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateHasToGroupByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateHasToGroupByIdParams(nil); got != (tables_model.UpdateHasToGroupByIdParams{}) {
		t.Errorf("ToUpdateHasToGroupByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteHasToGroupByIdParams(t *testing.T) {
	req := &tablesv1.DeleteHasToGroupByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteHasToGroupByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteHasToGroupByIdParams(req); got != want {
		t.Errorf("ToDeleteHasToGroupByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteHasToGroupByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteHasToGroupByIdParams(nil); got != (tables_model.DeleteHasToGroupByIdParams{}) {
		t.Errorf("ToDeleteHasToGroupByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteHasToGroupByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteHasToGroupByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteHasToGroupByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteHasToGroupByIdParams(req); got != want {
		t.Errorf("ToUndeleteHasToGroupByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteHasToGroupByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteHasToGroupByIdParams(nil); got != (tables_model.UndeleteHasToGroupByIdParams{}) {
		t.Errorf("ToUndeleteHasToGroupByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetHasToGroupsSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetHasToGroupsSearchNameParams(&tablesv1.GetHasToGroupsSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetHasToGroupsSearchNameDefaultsPage(t *testing.T) {
	got := ToGetHasToGroupsSearchNameParams(&tablesv1.GetHasToGroupsSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetHasToGroupsSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetHasToGroupsSearchNameParams(&tablesv1.GetHasToGroupsSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetHasToGroupsSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetHasToGroupsSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetHasToGroupsSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}

func TestGetHasToGroupsDefaultsPageLimit(t *testing.T) {
	got := ToGetHasToGroupsParams(&tablesv1.GetHasToGroupsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetHasToGroupsDefaultsPage(t *testing.T) {
	got := ToGetHasToGroupsParams(&tablesv1.GetHasToGroupsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetHasToGroupsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetHasToGroupsParams(&tablesv1.GetHasToGroupsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetHasToGroupsDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetHasToGroupsDeletedParams(&tablesv1.GetHasToGroupsDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetHasToGroupsDeletedDefaultsPage(t *testing.T) {
	got := ToGetHasToGroupsDeletedParams(&tablesv1.GetHasToGroupsDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetHasToGroupsDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetHasToGroupsDeletedParams(&tablesv1.GetHasToGroupsDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
