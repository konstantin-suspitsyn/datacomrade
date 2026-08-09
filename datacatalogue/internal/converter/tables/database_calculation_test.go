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

// databaseCalculationPtr — адрес значения для optional-полей proto.
func databaseCalculationPtr[T any](v T) *T {
	return &v
}

var (
	databaseCalculationCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	databaseCalculationUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testDatabaseCalculationRow — строка dc.database_calculation со значениями, различимыми между полями.
func testDatabaseCalculationRow() tables_model.DcDatabaseCalculation {
	return tables_model.DcDatabaseCalculation{
		ID:                100,
		DatabaseCatID:     101,
		CalculationTypeID: 102,
		CreatedAt:         databaseCalculationCreatedAt,
		UpdatedAt:         databaseCalculationUpdatedAt,
		IsDeleted:         false,
		UserID:            106,
	}
}

func TestDatabaseCalculationToProto(t *testing.T) {
	row := testDatabaseCalculationRow()
	got := DatabaseCalculationToProto(row)

	if got == nil {
		t.Fatal("DatabaseCalculationToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetDatabaseCatId() != row.DatabaseCatID {
		t.Errorf("DatabaseCatId = %d, want %d", got.GetDatabaseCatId(), row.DatabaseCatID)
	}

	if got.GetCalculationTypeId() != row.CalculationTypeID {
		t.Errorf("CalculationTypeId = %d, want %d", got.GetCalculationTypeId(), row.CalculationTypeID)
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

func TestDatabaseCalculationsToProto(t *testing.T) {
	first := testDatabaseCalculationRow()

	second := testDatabaseCalculationRow()
	second.ID = 999

	tests := []struct {
		name    string
		input   []tables_model.DcDatabaseCalculation
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcDatabaseCalculation{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcDatabaseCalculation{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseCalculationsToProto(tt.input)

			if got == nil {
				t.Fatal("DatabaseCalculationsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateDatabaseCalculationParams(t *testing.T) {
	req := &tablesv1.CreateDatabaseCalculationRequest{
		DatabaseCatId:     100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateDatabaseCalculationParams{
		DatabaseCatID:     100,
		CalculationTypeID: 101,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateDatabaseCalculationParams(req); got != want {
		t.Errorf("ToCreateDatabaseCalculationParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateDatabaseCalculationParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateDatabaseCalculationParams(nil); got != (tables_model.CreateDatabaseCalculationParams{}) {
		t.Errorf("ToCreateDatabaseCalculationParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateDatabaseCalculationByIdParams(t *testing.T) {
	req := &tablesv1.UpdateDatabaseCalculationByIdRequest{
		DatabaseCatId:     100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
		Id:                103,
	}

	want := tables_model.UpdateDatabaseCalculationByIdParams{
		DatabaseCatID:     100,
		CalculationTypeID: 101,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000003"),
		ID:                103,
	}

	if got := ToUpdateDatabaseCalculationByIdParams(req); got != want {
		t.Errorf("ToUpdateDatabaseCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDatabaseCalculationByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateDatabaseCalculationByIdParams(nil); got != (tables_model.UpdateDatabaseCalculationByIdParams{}) {
		t.Errorf("ToUpdateDatabaseCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteDatabaseCalculationByIdParams(t *testing.T) {
	req := &tablesv1.DeleteDatabaseCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteDatabaseCalculationByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteDatabaseCalculationByIdParams(req); got != want {
		t.Errorf("ToDeleteDatabaseCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteDatabaseCalculationByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteDatabaseCalculationByIdParams(nil); got != (tables_model.DeleteDatabaseCalculationByIdParams{}) {
		t.Errorf("ToDeleteDatabaseCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteDatabaseCalculationByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteDatabaseCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteDatabaseCalculationByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteDatabaseCalculationByIdParams(req); got != want {
		t.Errorf("ToUndeleteDatabaseCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteDatabaseCalculationByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteDatabaseCalculationByIdParams(nil); got != (tables_model.UndeleteDatabaseCalculationByIdParams{}) {
		t.Errorf("ToUndeleteDatabaseCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetDatabaseCalculationsSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsSearchNameParams(&tablesv1.GetDatabaseCalculationsSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCalculationsSearchNameDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCalculationsSearchNameParams(&tablesv1.GetDatabaseCalculationsSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCalculationsSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsSearchNameParams(&tablesv1.GetDatabaseCalculationsSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseCalculationsSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetDatabaseCalculationsSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetDatabaseCalculationsSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}

func TestGetDatabaseCalculationsFilteredDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsFilteredParams(&tablesv1.GetDatabaseCalculationsFilteredRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCalculationsFilteredDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCalculationsFilteredParams(&tablesv1.GetDatabaseCalculationsFilteredRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCalculationsFilteredKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsFilteredParams(&tablesv1.GetDatabaseCalculationsFilteredRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseCalculationsFilteredPassesFilterFields(t *testing.T) {
	req := &tablesv1.GetDatabaseCalculationsFilteredRequest{PageLimit: 10, Page: 1,
		DatabaseCatId:     databaseCalculationPtr(int64(100)),
		CalculationTypeId: databaseCalculationPtr(int64(101)),
	}

	got := ToGetDatabaseCalculationsFilteredParams(req)

	if got.DatabaseCatID != (sql.NullInt64{Int64: 100, Valid: true}) {
		t.Errorf("DatabaseCatID = %v, want %v", got.DatabaseCatID, (sql.NullInt64{Int64: 100, Valid: true}))
	}
	if got.CalculationTypeID != (sql.NullInt64{Int64: 101, Valid: true}) {
		t.Errorf("CalculationTypeID = %v, want %v", got.CalculationTypeID, (sql.NullInt64{Int64: 101, Valid: true}))
	}
}

// Непереданный необязательный фильтр обязан остаться SQL NULL, а не
// превратиться в 0/false: иначе выборка молча сузилась бы.
func TestGetDatabaseCalculationsFilteredOmitsAbsentFilters(t *testing.T) {
	got := ToGetDatabaseCalculationsFilteredParams(&tablesv1.GetDatabaseCalculationsFilteredRequest{PageLimit: 10, Page: 1})

	if got.DatabaseCatID.Valid {
		t.Errorf("DatabaseCatID = %+v, want invalid", got.DatabaseCatID)
	}
	if got.CalculationTypeID.Valid {
		t.Errorf("CalculationTypeID = %+v, want invalid", got.CalculationTypeID)
	}
}

func TestGetDatabaseCalculationsDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsDeletedParams(&tablesv1.GetDatabaseCalculationsDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCalculationsDeletedDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCalculationsDeletedParams(&tablesv1.GetDatabaseCalculationsDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCalculationsDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsDeletedParams(&tablesv1.GetDatabaseCalculationsDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseCalculationsDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsParams(&tablesv1.GetDatabaseCalculationsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseCalculationsDefaultsPage(t *testing.T) {
	got := ToGetDatabaseCalculationsParams(&tablesv1.GetDatabaseCalculationsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseCalculationsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseCalculationsParams(&tablesv1.GetDatabaseCalculationsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
