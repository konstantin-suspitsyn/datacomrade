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

// followingCalculationPtr — адрес значения для optional-полей proto.
func followingCalculationPtr[T any](v T) *T {
	return &v
}

var (
	followingCalculationCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	followingCalculationUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testFollowingCalculationRow — строка dc.following_calculation со значениями, различимыми между полями.
func testFollowingCalculationRow() tables_model.DcFollowingCalculation {
	return tables_model.DcFollowingCalculation{
		ID:                100,
		ColumnCatID:       101,
		CalculationTypeID: 102,
		CreatedAt:         followingCalculationCreatedAt,
		UpdatedAt:         followingCalculationUpdatedAt,
		IsDeleted:         false,
		UserID:            106,
	}
}

func TestFollowingCalculationToProto(t *testing.T) {
	row := testFollowingCalculationRow()
	got := FollowingCalculationToProto(row)

	if got == nil {
		t.Fatal("FollowingCalculationToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetColumnCatId() != row.ColumnCatID {
		t.Errorf("ColumnCatId = %d, want %d", got.GetColumnCatId(), row.ColumnCatID)
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

func TestFollowingCalculationsToProto(t *testing.T) {
	first := testFollowingCalculationRow()

	second := testFollowingCalculationRow()
	second.ID = 999

	tests := []struct {
		name    string
		input   []tables_model.DcFollowingCalculation
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcFollowingCalculation{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcFollowingCalculation{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FollowingCalculationsToProto(tt.input)

			if got == nil {
				t.Fatal("FollowingCalculationsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateFollowingCalculationParams(t *testing.T) {
	req := &tablesv1.CreateFollowingCalculationRequest{
		ColumnCatId:       100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateFollowingCalculationParams{
		ColumnCatID:       100,
		CalculationTypeID: 101,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateFollowingCalculationParams(req); got != want {
		t.Errorf("ToCreateFollowingCalculationParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateFollowingCalculationParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateFollowingCalculationParams(nil); got != (tables_model.CreateFollowingCalculationParams{}) {
		t.Errorf("ToCreateFollowingCalculationParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateFollowingCalculationByIdParams(t *testing.T) {
	req := &tablesv1.UpdateFollowingCalculationByIdRequest{
		ColumnCatId:       100,
		CalculationTypeId: 101,
		ExternalId:        "00000000-0000-4000-8000-000000000003",
		Id:                103,
	}

	want := tables_model.UpdateFollowingCalculationByIdParams{
		ColumnCatID:       100,
		CalculationTypeID: 101,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000003"),
		ID:                103,
	}

	if got := ToUpdateFollowingCalculationByIdParams(req); got != want {
		t.Errorf("ToUpdateFollowingCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateFollowingCalculationByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateFollowingCalculationByIdParams(nil); got != (tables_model.UpdateFollowingCalculationByIdParams{}) {
		t.Errorf("ToUpdateFollowingCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteFollowingCalculationByIdParams(t *testing.T) {
	req := &tablesv1.DeleteFollowingCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteFollowingCalculationByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteFollowingCalculationByIdParams(req); got != want {
		t.Errorf("ToDeleteFollowingCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteFollowingCalculationByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteFollowingCalculationByIdParams(nil); got != (tables_model.DeleteFollowingCalculationByIdParams{}) {
		t.Errorf("ToDeleteFollowingCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteFollowingCalculationByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteFollowingCalculationByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteFollowingCalculationByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteFollowingCalculationByIdParams(req); got != want {
		t.Errorf("ToUndeleteFollowingCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteFollowingCalculationByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteFollowingCalculationByIdParams(nil); got != (tables_model.UndeleteFollowingCalculationByIdParams{}) {
		t.Errorf("ToUndeleteFollowingCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetFollowingCalculationsDefaultsPageLimit(t *testing.T) {
	got := ToGetFollowingCalculationsParams(&tablesv1.GetFollowingCalculationsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetFollowingCalculationsDefaultsPage(t *testing.T) {
	got := ToGetFollowingCalculationsParams(&tablesv1.GetFollowingCalculationsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetFollowingCalculationsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetFollowingCalculationsParams(&tablesv1.GetFollowingCalculationsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetFollowingCalculationsDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetFollowingCalculationsDeletedParams(&tablesv1.GetFollowingCalculationsDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetFollowingCalculationsDeletedDefaultsPage(t *testing.T) {
	got := ToGetFollowingCalculationsDeletedParams(&tablesv1.GetFollowingCalculationsDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetFollowingCalculationsDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetFollowingCalculationsDeletedParams(&tablesv1.GetFollowingCalculationsDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetFollowingCalculationsFiltersDefaultsPageLimit(t *testing.T) {
	got := ToGetFollowingCalculationsFiltersParams(&tablesv1.GetFollowingCalculationsFiltersRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetFollowingCalculationsFiltersDefaultsPage(t *testing.T) {
	got := ToGetFollowingCalculationsFiltersParams(&tablesv1.GetFollowingCalculationsFiltersRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetFollowingCalculationsFiltersKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetFollowingCalculationsFiltersParams(&tablesv1.GetFollowingCalculationsFiltersRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetFollowingCalculationsFiltersPassesFilterFields(t *testing.T) {
	req := &tablesv1.GetFollowingCalculationsFiltersRequest{PageLimit: 10, Page: 1,
		ColumnCatId:       followingCalculationPtr(int64(100)),
		CalculationTypeId: followingCalculationPtr(int64(101)),
	}

	got := ToGetFollowingCalculationsFiltersParams(req)

	if got.ColumnCatID != (sql.NullInt64{Int64: 100, Valid: true}) {
		t.Errorf("ColumnCatID = %v, want %v", got.ColumnCatID, (sql.NullInt64{Int64: 100, Valid: true}))
	}
	if got.CalculationTypeID != (sql.NullInt64{Int64: 101, Valid: true}) {
		t.Errorf("CalculationTypeID = %v, want %v", got.CalculationTypeID, (sql.NullInt64{Int64: 101, Valid: true}))
	}
}

// Непереданный необязательный фильтр обязан остаться SQL NULL, а не
// превратиться в 0/false: иначе выборка молча сузилась бы.
func TestGetFollowingCalculationsFiltersOmitsAbsentFilters(t *testing.T) {
	got := ToGetFollowingCalculationsFiltersParams(&tablesv1.GetFollowingCalculationsFiltersRequest{PageLimit: 10, Page: 1})

	if got.ColumnCatID.Valid {
		t.Errorf("ColumnCatID = %+v, want invalid", got.ColumnCatID)
	}
	if got.CalculationTypeID.Valid {
		t.Errorf("CalculationTypeID = %+v, want invalid", got.CalculationTypeID)
	}
}
