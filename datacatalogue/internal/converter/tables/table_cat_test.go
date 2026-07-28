package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	tableCatCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	tableCatUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testTableCatRow — строка dc.table_cat со значениями, различимыми между полями.
func testTableCatRow() tables_model.DcTableCat {
	return tables_model.DcTableCat{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		SchemaID:    103,
		TableTypeID: 104,
		DomainID:    105,
		IsDeleted:   false,
		CreatedAt:   tableCatCreatedAt,
		UpdatedAt:   tableCatUpdatedAt,
		IsGetDict:   true,
		UserID:      110,
	}
}

func TestTableCatToProto(t *testing.T) {
	row := testTableCatRow()
	got := TableCatToProto(row)

	if got == nil {
		t.Fatal("TableCatToProto() = nil, want value")
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

	if got.GetSchemaId() != row.SchemaID {
		t.Errorf("SchemaId = %d, want %d", got.GetSchemaId(), row.SchemaID)
	}

	if got.GetTableTypeId() != row.TableTypeID {
		t.Errorf("TableTypeId = %d, want %d", got.GetTableTypeId(), row.TableTypeID)
	}

	if got.GetDomainId() != row.DomainID {
		t.Errorf("DomainId = %d, want %d", got.GetDomainId(), row.DomainID)
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

	if got.GetIsGetDict() != row.IsGetDict {
		t.Errorf("IsGetDict = %v, want %v", got.GetIsGetDict(), row.IsGetDict)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

}

func TestTableCatToProtoDeleted(t *testing.T) {
	row := testTableCatRow()
	row.IsDeleted = true

	if got := TableCatToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestTableCatsToProto(t *testing.T) {
	first := testTableCatRow()

	second := testTableCatRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcTableCat
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcTableCat{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcTableCat{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TableCatsToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("TableCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestTableCatsToProtoKeepsOrder(t *testing.T) {
	first := testTableCatRow()
	second := testTableCatRow()
	second.Name = "second-value"

	got := TableCatsToProto([]tables_model.DcTableCat{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateTableCatParams(t *testing.T) {
	req := &tablesv1.CreateTableCatRequest{
		Name:        "name-0",
		Description: "description-0",
		SchemaId:    102,
		TableTypeId: 103,
		DomainId:    104,
		IsGetDict:   true,
		UserId:      106,
	}

	want := tables_model.CreateTableCatParams{
		Name:        "name-0",
		Description: "description-0",
		SchemaID:    102,
		TableTypeID: 103,
		DomainID:    104,
		IsGetDict:   true,
		UserID:      106,
	}

	if got := ToCreateTableCatParams(req); got != want {
		t.Errorf("ToCreateTableCatParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateTableCatParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateTableCatParams(nil); got != (tables_model.CreateTableCatParams{}) {
		t.Errorf("ToCreateTableCatParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateTableCatByIdParams(t *testing.T) {
	req := &tablesv1.UpdateTableCatByIdRequest{
		Id:          100,
		Name:        "name-0",
		Description: "description-0",
		SchemaId:    103,
		TableTypeId: 104,
		DomainId:    105,
		IsGetDict:   true,
		UserId:      107,
	}

	want := tables_model.UpdateTableCatByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		SchemaID:    103,
		TableTypeID: 104,
		DomainID:    105,
		IsGetDict:   true,
		UserID:      107,
	}

	if got := ToUpdateTableCatByIdParams(req); got != want {
		t.Errorf("ToUpdateTableCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateTableCatByIdParamsNil(t *testing.T) {
	if got := ToUpdateTableCatByIdParams(nil); got != (tables_model.UpdateTableCatByIdParams{}) {
		t.Errorf("ToUpdateTableCatByIdParams(nil) = %+v, want zero value", got)
	}
}
