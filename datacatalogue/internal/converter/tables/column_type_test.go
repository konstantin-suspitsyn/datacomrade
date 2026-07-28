package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
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

func TestColumnTypeToProtoDeleted(t *testing.T) {
	row := testColumnTypeRow()
	row.IsDeleted = true

	if got := ColumnTypeToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
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

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("ColumnTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestColumnTypesToProtoKeepsOrder(t *testing.T) {
	first := testColumnTypeRow()
	second := testColumnTypeRow()
	second.Name = "second-value"

	got := ColumnTypesToProto([]tables_model.DcColumnType{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateColumnTypeParams(t *testing.T) {
	req := &tablesv1.CreateColumnTypeRequest{
		Name:        "name-0",
		Description: "description-0",
		UserId:      102,
	}

	want := tables_model.CreateColumnTypeParams{
		Name:        "name-0",
		Description: "description-0",
		UserID:      102,
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
		Id:          100,
		Name:        "name-0",
		Description: "description-0",
		UserId:      103,
	}

	want := tables_model.UpdateColumnTypeByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		UserID:      103,
	}

	if got := ToUpdateColumnTypeByIdParams(req); got != want {
		t.Errorf("ToUpdateColumnTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateColumnTypeByIdParamsNil(t *testing.T) {
	if got := ToUpdateColumnTypeByIdParams(nil); got != (tables_model.UpdateColumnTypeByIdParams{}) {
		t.Errorf("ToUpdateColumnTypeByIdParams(nil) = %+v, want zero value", got)
	}
}
