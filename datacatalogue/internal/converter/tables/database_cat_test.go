package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	databaseCatCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	databaseCatUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testDatabaseCatRow — строка dc.database_cat со значениями, различимыми между полями.
func testDatabaseCatRow() tables_model.DcDatabaseCat {
	return tables_model.DcDatabaseCat{
		ID:             100,
		Name:           "name-0",
		HostID:         102,
		DatabaseTypeID: 103,
		Description:    "description-0",
		IsDeleted:      false,
		CreatedAt:      databaseCatCreatedAt,
		UpdatedAt:      databaseCatUpdatedAt,
		UserID:         108,
	}
}

func TestDatabaseCatToProto(t *testing.T) {
	row := testDatabaseCatRow()
	got := DatabaseCatToProto(row)

	if got == nil {
		t.Fatal("DatabaseCatToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
	}

	if got.GetHostId() != row.HostID {
		t.Errorf("HostId = %d, want %d", got.GetHostId(), row.HostID)
	}

	if got.GetDatabaseTypeId() != row.DatabaseTypeID {
		t.Errorf("DatabaseTypeId = %d, want %d", got.GetDatabaseTypeId(), row.DatabaseTypeID)
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

func TestDatabaseCatToProtoDeleted(t *testing.T) {
	row := testDatabaseCatRow()
	row.IsDeleted = true

	if got := DatabaseCatToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestDatabaseCatsToProto(t *testing.T) {
	first := testDatabaseCatRow()

	second := testDatabaseCatRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcDatabaseCat
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcDatabaseCat{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcDatabaseCat{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseCatsToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("DatabaseCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestDatabaseCatsToProtoKeepsOrder(t *testing.T) {
	first := testDatabaseCatRow()
	second := testDatabaseCatRow()
	second.Name = "second-value"

	got := DatabaseCatsToProto([]tables_model.DcDatabaseCat{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateDatabaseCatParams(t *testing.T) {
	req := &tablesv1.CreateDatabaseCatRequest{
		Name:           "name-0",
		HostId:         101,
		DatabaseTypeId: 102,
		Description:    "description-0",
		UserId:         104,
	}

	want := tables_model.CreateDatabaseCatParams{
		Name:           "name-0",
		HostID:         101,
		DatabaseTypeID: 102,
		Description:    "description-0",
		UserID:         104,
	}

	if got := ToCreateDatabaseCatParams(req); got != want {
		t.Errorf("ToCreateDatabaseCatParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateDatabaseCatParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateDatabaseCatParams(nil); got != (tables_model.CreateDatabaseCatParams{}) {
		t.Errorf("ToCreateDatabaseCatParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateDatabaseCatByIdParams(t *testing.T) {
	req := &tablesv1.UpdateDatabaseCatByIdRequest{
		Id:             100,
		Name:           "name-0",
		HostId:         102,
		DatabaseTypeId: 103,
		Description:    "description-0",
		UserId:         105,
	}

	want := tables_model.UpdateDatabaseCatByIdParams{
		ID:             100,
		Name:           "name-0",
		HostID:         102,
		DatabaseTypeID: 103,
		Description:    "description-0",
		UserID:         105,
	}

	if got := ToUpdateDatabaseCatByIdParams(req); got != want {
		t.Errorf("ToUpdateDatabaseCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDatabaseCatByIdParamsNil(t *testing.T) {
	if got := ToUpdateDatabaseCatByIdParams(nil); got != (tables_model.UpdateDatabaseCatByIdParams{}) {
		t.Errorf("ToUpdateDatabaseCatByIdParams(nil) = %+v, want zero value", got)
	}
}
