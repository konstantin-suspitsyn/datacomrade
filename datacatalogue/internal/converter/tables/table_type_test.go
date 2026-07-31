package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	tableTypeCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	tableTypeUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testTableTypeRow — строка dc.table_type со значениями, различимыми между полями.
func testTableTypeRow() tables_model.DcTableType {
	return tables_model.DcTableType{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		IsDeleted:   false,
		CreatedAt:   tableTypeCreatedAt,
		UpdatedAt:   tableTypeUpdatedAt,
		UserID:      106,
	}
}

func TestTableTypeToProto(t *testing.T) {
	row := testTableTypeRow()
	got := TableTypeToProto(row)

	if got == nil {
		t.Fatal("TableTypeToProto() = nil, want value")
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

func TestTableTypeToProtoDeleted(t *testing.T) {
	row := testTableTypeRow()
	row.IsDeleted = true

	if got := TableTypeToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestTableTypesToProto(t *testing.T) {
	first := testTableTypeRow()

	second := testTableTypeRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcTableType
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcTableType{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcTableType{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TableTypesToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("TableTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestTableTypesToProtoKeepsOrder(t *testing.T) {
	first := testTableTypeRow()
	second := testTableTypeRow()
	second.Name = "second-value"

	got := TableTypesToProto([]tables_model.DcTableType{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateTableTypeParams(t *testing.T) {
	req := &tablesv1.CreateTableTypeRequest{
		Name:           "name-0",
		Description:    "description-0",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateTableTypeParams{
		Name:        "name-0",
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateTableTypeParams(req); got != want {
		t.Errorf("ToCreateTableTypeParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateTableTypeParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateTableTypeParams(nil); got != (tables_model.CreateTableTypeParams{}) {
		t.Errorf("ToCreateTableTypeParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateTableTypeByIdParams(t *testing.T) {
	req := &tablesv1.UpdateTableTypeByIdRequest{
		Id:             100,
		Name:           "name-0",
		Description:    "description-0",
		UserExternalId: "00000000-0000-4000-8000-000000000004",
	}

	want := tables_model.UpdateTableTypeByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToUpdateTableTypeByIdParams(req); got != want {
		t.Errorf("ToUpdateTableTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateTableTypeByIdParamsNil(t *testing.T) {
	if got := ToUpdateTableTypeByIdParams(nil); got != (tables_model.UpdateTableTypeByIdParams{}) {
		t.Errorf("ToUpdateTableTypeByIdParams(nil) = %+v, want zero value", got)
	}
}
