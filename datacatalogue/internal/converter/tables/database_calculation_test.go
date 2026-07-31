package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

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

func TestDatabaseCalculationToProtoDeleted(t *testing.T) {
	row := testDatabaseCalculationRow()
	row.IsDeleted = true

	if got := DatabaseCalculationToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
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

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("DatabaseCalculationsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestDatabaseCalculationsToProtoKeepsOrder(t *testing.T) {
	first := testDatabaseCalculationRow()
	second := testDatabaseCalculationRow()
	second.ID = 999

	got := DatabaseCalculationsToProto([]tables_model.DcDatabaseCalculation{first, second})

	if got[0].GetId() != first.ID {
		t.Errorf("[0] = %d, want %d", got[0].GetId(), first.ID)
	}

	if got[1].GetId() != second.ID {
		t.Errorf("[1] = %d, want %d", got[1].GetId(), second.ID)
	}
}

func TestToCreateDatabaseCalculationParams(t *testing.T) {
	req := &tablesv1.CreateDatabaseCalculationRequest{
		DatabaseCatId:     100,
		CalculationTypeId: 101,
		UserExternalId:    "00000000-0000-4000-8000-000000000003",
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
		Id:                100,
		DatabaseCatId:     101,
		CalculationTypeId: 102,
		UserExternalId:    "00000000-0000-4000-8000-000000000004",
	}

	want := tables_model.UpdateDatabaseCalculationByIdParams{
		ID:                100,
		DatabaseCatID:     101,
		CalculationTypeID: 102,
		ExternalID:        uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToUpdateDatabaseCalculationByIdParams(req); got != want {
		t.Errorf("ToUpdateDatabaseCalculationByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDatabaseCalculationByIdParamsNil(t *testing.T) {
	if got := ToUpdateDatabaseCalculationByIdParams(nil); got != (tables_model.UpdateDatabaseCalculationByIdParams{}) {
		t.Errorf("ToUpdateDatabaseCalculationByIdParams(nil) = %+v, want zero value", got)
	}
}
