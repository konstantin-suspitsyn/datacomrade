package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	groupLevelCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	groupLevelUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testGroupLevelRow — строка dc.group_levels со значениями, различимыми между полями.
func testGroupLevelRow() tables_model.DcGroupLevel {
	return tables_model.DcGroupLevel{
		ID:             100,
		ColumnID:       101,
		ParentColumnID: 102,
		Level:          13,
		Description:    "description-0",
		CreatedAt:      groupLevelCreatedAt,
		UpdatedAt:      groupLevelUpdatedAt,
		IsDeleted:      false,
		UserID:         108,
	}
}

func TestGroupLevelToProto(t *testing.T) {
	row := testGroupLevelRow()
	got := GroupLevelToProto(row)

	if got == nil {
		t.Fatal("GroupLevelToProto() = nil, want value")
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

func TestGroupLevelToProtoDeleted(t *testing.T) {
	row := testGroupLevelRow()
	row.IsDeleted = true

	if got := GroupLevelToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestGroupLevelsToProto(t *testing.T) {
	first := testGroupLevelRow()

	second := testGroupLevelRow()
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
			got := GroupLevelsToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("GroupLevelsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestGroupLevelsToProtoKeepsOrder(t *testing.T) {
	first := testGroupLevelRow()
	second := testGroupLevelRow()
	second.Description = "second-value"

	got := GroupLevelsToProto([]tables_model.DcGroupLevel{first, second})

	if got[0].GetDescription() != first.Description {
		t.Errorf("[0] = %q, want %q", got[0].GetDescription(), first.Description)
	}

	if got[1].GetDescription() != second.Description {
		t.Errorf("[1] = %q, want %q", got[1].GetDescription(), second.Description)
	}
}

func TestToCreateGroupLevelParams(t *testing.T) {
	req := &tablesv1.CreateGroupLevelRequest{
		ColumnId:       100,
		ParentColumnId: 101,
		Level:          12,
		Description:    "description-0",
		UserExternalId: "00000000-0000-4000-8000-000000000005",
	}

	want := tables_model.CreateGroupLevelParams{
		ColumnID:       100,
		ParentColumnID: 101,
		Level:          12,
		Description:    "description-0",
		ExternalID:     uuid.MustParse("00000000-0000-4000-8000-000000000005"),
	}

	if got := ToCreateGroupLevelParams(req); got != want {
		t.Errorf("ToCreateGroupLevelParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateGroupLevelParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateGroupLevelParams(nil); got != (tables_model.CreateGroupLevelParams{}) {
		t.Errorf("ToCreateGroupLevelParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateGroupLevelByIdParams(t *testing.T) {
	req := &tablesv1.UpdateGroupLevelByIdRequest{
		Id:             100,
		ColumnId:       101,
		ParentColumnId: 102,
		Level:          13,
		Description:    "description-0",
		UserExternalId: "00000000-0000-4000-8000-000000000006",
	}

	want := tables_model.UpdateGroupLevelByIdParams{
		ID:             100,
		ColumnID:       101,
		ParentColumnID: 102,
		Level:          13,
		Description:    "description-0",
		ExternalID:     uuid.MustParse("00000000-0000-4000-8000-000000000006"),
	}

	if got := ToUpdateGroupLevelByIdParams(req); got != want {
		t.Errorf("ToUpdateGroupLevelByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateGroupLevelByIdParamsNil(t *testing.T) {
	if got := ToUpdateGroupLevelByIdParams(nil); got != (tables_model.UpdateGroupLevelByIdParams{}) {
		t.Errorf("ToUpdateGroupLevelByIdParams(nil) = %+v, want zero value", got)
	}
}
