package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

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

func TestColumnCatToProtoDeleted(t *testing.T) {
	row := testColumnCatRow()
	row.IsDeleted = true

	if got := ColumnCatToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
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

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("ColumnCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestColumnCatsToProtoKeepsOrder(t *testing.T) {
	first := testColumnCatRow()
	second := testColumnCatRow()
	second.Name = "second-value"

	got := ColumnCatsToProto([]tables_model.DcColumnCat{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
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
		UserId:            107,
	}

	want := tables_model.CreateColumnCatParams{
		TableID:           100,
		Name:              "name-0",
		AliasID:           102,
		ColumnTypeID:      103,
		Description:       "description-0",
		CalculationTypeID: 105,
		ShowInUi:          true,
		UserID:            107,
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
		Id:                100,
		TableId:           101,
		Name:              "name-0",
		AliasId:           103,
		ColumnTypeId:      104,
		Description:       "description-0",
		CalculationTypeId: 106,
		ShowInUi:          true,
		UserId:            108,
	}

	want := tables_model.UpdateColumnCatByIdParams{
		ID:                100,
		TableID:           101,
		Name:              "name-0",
		AliasID:           103,
		ColumnTypeID:      104,
		Description:       "description-0",
		CalculationTypeID: 106,
		ShowInUi:          true,
		UserID:            108,
	}

	if got := ToUpdateColumnCatByIdParams(req); got != want {
		t.Errorf("ToUpdateColumnCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateColumnCatByIdParamsNil(t *testing.T) {
	if got := ToUpdateColumnCatByIdParams(nil); got != (tables_model.UpdateColumnCatByIdParams{}) {
		t.Errorf("ToUpdateColumnCatByIdParams(nil) = %+v, want zero value", got)
	}
}
