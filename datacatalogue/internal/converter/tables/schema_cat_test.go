package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	schemaCatCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	schemaCatUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testSchemaCatRow — строка dc.schema_cat со значениями, различимыми между полями.
func testSchemaCatRow() tables_model.DcSchemaCat {
	return tables_model.DcSchemaCat{
		ID:         100,
		DatabaseID: 101,
		Name:       "name-0",
		IsDeleted:  false,
		CreatedAt:  schemaCatCreatedAt,
		UpdatedAt:  schemaCatUpdatedAt,
		UserID:     106,
	}
}

func TestSchemaCatToProto(t *testing.T) {
	row := testSchemaCatRow()
	got := SchemaCatToProto(row)

	if got == nil {
		t.Fatal("SchemaCatToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetDatabaseId() != row.DatabaseID {
		t.Errorf("DatabaseId = %d, want %d", got.GetDatabaseId(), row.DatabaseID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
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

func TestSchemaCatToProtoDeleted(t *testing.T) {
	row := testSchemaCatRow()
	row.IsDeleted = true

	if got := SchemaCatToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestSchemaCatsToProto(t *testing.T) {
	first := testSchemaCatRow()

	second := testSchemaCatRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcSchemaCat
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcSchemaCat{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcSchemaCat{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SchemaCatsToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("SchemaCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestSchemaCatsToProtoKeepsOrder(t *testing.T) {
	first := testSchemaCatRow()
	second := testSchemaCatRow()
	second.Name = "second-value"

	got := SchemaCatsToProto([]tables_model.DcSchemaCat{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateSchemaCatParams(t *testing.T) {
	req := &tablesv1.CreateSchemaCatRequest{
		DatabaseId: 100,
		Name:       "name-0",
		UserId:     102,
	}

	want := tables_model.CreateSchemaCatParams{
		DatabaseID: 100,
		Name:       "name-0",
		UserID:     102,
	}

	if got := ToCreateSchemaCatParams(req); got != want {
		t.Errorf("ToCreateSchemaCatParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateSchemaCatParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateSchemaCatParams(nil); got != (tables_model.CreateSchemaCatParams{}) {
		t.Errorf("ToCreateSchemaCatParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateSchemaCatByIdParams(t *testing.T) {
	req := &tablesv1.UpdateSchemaCatByIdRequest{
		Id:         100,
		DatabaseId: 101,
		Name:       "name-0",
		UserId:     103,
	}

	want := tables_model.UpdateSchemaCatByIdParams{
		ID:         100,
		DatabaseID: 101,
		Name:       "name-0",
		UserID:     103,
	}

	if got := ToUpdateSchemaCatByIdParams(req); got != want {
		t.Errorf("ToUpdateSchemaCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateSchemaCatByIdParamsNil(t *testing.T) {
	if got := ToUpdateSchemaCatByIdParams(nil); got != (tables_model.UpdateSchemaCatByIdParams{}) {
		t.Errorf("ToUpdateSchemaCatByIdParams(nil) = %+v, want zero value", got)
	}
}
