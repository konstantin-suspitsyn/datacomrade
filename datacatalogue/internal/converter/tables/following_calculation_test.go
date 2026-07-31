package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

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

func TestFollowingCalculationToProtoDeleted(t *testing.T) {
	row := testFollowingCalculationRow()
	row.IsDeleted = true

	if got := FollowingCalculationToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
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

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("FollowingCalculationsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestFollowingCalculationsToProtoKeepsOrder(t *testing.T) {
	first := testFollowingCalculationRow()
	second := testFollowingCalculationRow()
	second.ID = 999

	got := FollowingCalculationsToProto([]tables_model.DcFollowingCalculation{first, second})

	if got[0].GetId() != first.ID {
		t.Errorf("[0] = %d, want %d", got[0].GetId(), first.ID)
	}

	if got[1].GetId() != second.ID {
		t.Errorf("[1] = %d, want %d", got[1].GetId(), second.ID)
	}
}

func TestToCreateFollowingCalculationParams(t *testing.T) {
	req := &tablesv1.CreateFollowingCalculationRequest{
		ColumnCatId:       100,
		CalculationTypeId: 101,
		UserExternalId:    "00000000-0000-4000-8000-000000000003",
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
		Id:                100,
		ColumnCatId:       101,
		CalculationTypeId: 102,
		UserExternalId:    "00000000-0000-4000-8000-000000000004",
	}

	want := tables_model.UpdateFollowingCalculationByIdParams{
		ID:                100,
		ColumnCatID:       101,
		CalculationTypeID: 102,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToUpdateFollowingCalculationByIdParams(req); got != want {
		t.Errorf("ToUpdateFollowingCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateFollowingCalculationByIdParamsNil(t *testing.T) {
	if got := ToUpdateFollowingCalculationByIdParams(nil); got != (tables_model.UpdateFollowingCalculationByIdParams{}) {
		t.Errorf("ToUpdateFollowingCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}
