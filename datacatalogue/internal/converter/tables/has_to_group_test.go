package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
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

func TestHasToGroupToProtoDeleted(t *testing.T) {
	row := testHasToGroupRow()
	row.IsDeleted = true

	if got := HasToGroupToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
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

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("HasToGroupsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestHasToGroupsToProtoKeepsOrder(t *testing.T) {
	first := testHasToGroupRow()
	second := testHasToGroupRow()
	second.Description = "second-value"

	got := HasToGroupsToProto([]tables_model.DcHasToGroup{first, second})

	if got[0].GetDescription() != first.Description {
		t.Errorf("[0] = %q, want %q", got[0].GetDescription(), first.Description)
	}

	if got[1].GetDescription() != second.Description {
		t.Errorf("[1] = %q, want %q", got[1].GetDescription(), second.Description)
	}
}

func TestToCreateHasToGroupParams(t *testing.T) {
	req := &tablesv1.CreateHasToGroupRequest{
		ColumnIdA:   100,
		ColumnIdB:   101,
		Description: "description-0",
		UserId:      103,
	}

	want := tables_model.CreateHasToGroupParams{
		ColumnIDA:   100,
		ColumnIDB:   101,
		Description: "description-0",
		UserID:      103,
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
		Id:          100,
		ColumnIdA:   101,
		ColumnIdB:   102,
		Description: "description-0",
		UserId:      104,
	}

	want := tables_model.UpdateHasToGroupByIdParams{
		ID:          100,
		ColumnIDA:   101,
		ColumnIDB:   102,
		Description: "description-0",
		UserID:      104,
	}

	if got := ToUpdateHasToGroupByIdParams(req); got != want {
		t.Errorf("ToUpdateHasToGroupByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateHasToGroupByIdParamsNil(t *testing.T) {
	if got := ToUpdateHasToGroupByIdParams(nil); got != (tables_model.UpdateHasToGroupByIdParams{}) {
		t.Errorf("ToUpdateHasToGroupByIdParams(nil) = %+v, want zero value", got)
	}
}
