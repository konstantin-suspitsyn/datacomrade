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

// databaseCatPtr — адрес значения для optional-полей proto.
func databaseCatPtr[T any](v T) *T {
	return &v
}

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

			if got == nil {
				t.Fatal("DatabaseCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateDatabaseCatParams(t *testing.T) {
	req := &tablesv1.CreateDatabaseCatRequest{
		Name:           "name-0",
		HostId:         101,
		DatabaseTypeId: 102,
		Description:    "description-0",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
	}

	want := tables_model.CreateDatabaseCatParams{
		Name:           "name-0",
		HostID:         101,
		DatabaseTypeID: 102,
		Description:    "description-0",
		ExternalID:     uuid.MustParse("00000000-0000-4000-8000-000000000005"),
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
		Name:           "name-0",
		HostId:         101,
		DatabaseTypeId: 102,
		Description:    "description-0",
		ExternalId:     "00000000-0000-4000-8000-000000000005",
		Id:             105,
	}

	want := tables_model.UpdateDatabaseCatByIdParams{
		Name:           "name-0",
		HostID:         101,
		DatabaseTypeID: 102,
		Description:    "description-0",
		ExternalID:     uuid.MustParse("00000000-0000-4000-8000-000000000005"),
		ID:             105,
	}

	if got := ToUpdateDatabaseCatByIdParams(req); got != want {
		t.Errorf("ToUpdateDatabaseCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDatabaseCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateDatabaseCatByIdParams(nil); got != (tables_model.UpdateDatabaseCatByIdParams{}) {
		t.Errorf("ToUpdateDatabaseCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteDatabaseCatByIdParams(t *testing.T) {
	req := &tablesv1.DeleteDatabaseCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteDatabaseCatByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteDatabaseCatByIdParams(req); got != want {
		t.Errorf("ToDeleteDatabaseCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteDatabaseCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteDatabaseCatByIdParams(nil); got != (tables_model.DeleteDatabaseCatByIdParams{}) {
		t.Errorf("ToDeleteDatabaseCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteDatabaseCatByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteDatabaseCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteDatabaseCatByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteDatabaseCatByIdParams(req); got != want {
		t.Errorf("ToUndeleteDatabaseCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteDatabaseCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteDatabaseCatByIdParams(nil); got != (tables_model.UndeleteDatabaseCatByIdParams{}) {
		t.Errorf("ToUndeleteDatabaseCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetDatabaseCatsDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsParams(&tablesv1.GetDatabaseCatsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCatsDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCatsParams(&tablesv1.GetDatabaseCatsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCatsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsParams(&tablesv1.GetDatabaseCatsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseCatsDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsDeletedParams(&tablesv1.GetDatabaseCatsDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCatsDeletedDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCatsDeletedParams(&tablesv1.GetDatabaseCatsDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCatsDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsDeletedParams(&tablesv1.GetDatabaseCatsDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseCatsByHostIdDatabaseTypeIdDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsByHostIdDatabaseTypeIdParams(&tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCatsByHostIdDatabaseTypeIdDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCatsByHostIdDatabaseTypeIdParams(&tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCatsByHostIdDatabaseTypeIdKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsByHostIdDatabaseTypeIdParams(&tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseCatsByHostIdDatabaseTypeIdPassesFilterFields(t *testing.T) {
	req := &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{PageLimit: 10, Page: 1,
		HostId:         databaseCatPtr(int64(100)),
		DatabaseTypeId: databaseCatPtr(int64(101)),
	}

	got := ToGetDatabaseCatsByHostIdDatabaseTypeIdParams(req)

	if got.HostID != (sql.NullInt64{Int64: 100, Valid: true}) {
		t.Errorf("HostID = %v, want %v", got.HostID, (sql.NullInt64{Int64: 100, Valid: true}))
	}
	if got.DatabaseTypeID != (sql.NullInt64{Int64: 101, Valid: true}) {
		t.Errorf("DatabaseTypeID = %v, want %v", got.DatabaseTypeID, (sql.NullInt64{Int64: 101, Valid: true}))
	}
}

// Непереданный необязательный фильтр обязан остаться SQL NULL, а не
// превратиться в 0/false: иначе выборка молча сузилась бы.
func TestGetDatabaseCatsByHostIdDatabaseTypeIdOmitsAbsentFilters(t *testing.T) {
	got := ToGetDatabaseCatsByHostIdDatabaseTypeIdParams(&tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{PageLimit: 10, Page: 1})

	if got.HostID.Valid {
		t.Errorf("HostID = %+v, want invalid", got.HostID)
	}
	if got.DatabaseTypeID.Valid {
		t.Errorf("DatabaseTypeID = %+v, want invalid", got.DatabaseTypeID)
	}
}

func TestGetDatabaseCatsSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsSearchNameParams(&tablesv1.GetDatabaseCatsSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCatsSearchNameDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCatsSearchNameParams(&tablesv1.GetDatabaseCatsSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCatsSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCatsSearchNameParams(&tablesv1.GetDatabaseCatsSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseCatsSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetDatabaseCatsSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetDatabaseCatsSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}
