package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	calculationTypeCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	calculationTypeUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testCalculationTypeRow — строка dc.calculation_type со значениями, различимыми между полями.
func testCalculationTypeRow() tables_model.DcCalculationType {
	return tables_model.DcCalculationType{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		CreatedAt:   calculationTypeCreatedAt,
		UpdatedAt:   calculationTypeUpdatedAt,
		IsDeleted:   false,
	}
}

func TestCalculationTypeToProto(t *testing.T) {
	row := testCalculationTypeRow()
	got := CalculationTypeToProto(row)

	if got == nil {
		t.Fatal("CalculationTypeToProto() = nil, want value")
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

	if !got.GetCreatedAt().AsTime().Equal(row.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.GetCreatedAt().AsTime(), row.CreatedAt)
	}

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt)
	}

	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}

}

func TestCalculationTypeToProtoDeleted(t *testing.T) {
	row := testCalculationTypeRow()
	row.IsDeleted = true

	if got := CalculationTypeToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestCalculationTypesToProto(t *testing.T) {
	first := testCalculationTypeRow()

	second := testCalculationTypeRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcCalculationType
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcCalculationType{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcCalculationType{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculationTypesToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("CalculationTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestCalculationTypesToProtoKeepsOrder(t *testing.T) {
	first := testCalculationTypeRow()
	second := testCalculationTypeRow()
	second.Name = "second-value"

	got := CalculationTypesToProto([]tables_model.DcCalculationType{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateCalculationTypeParams(t *testing.T) {
	req := &tablesv1.CreateCalculationTypeRequest{
		Name:        "name-0",
		Description: "description-0",
	}

	want := tables_model.CreateCalculationTypeParams{
		Name:        "name-0",
		Description: "description-0",
	}

	if got := ToCreateCalculationTypeParams(req); got != want {
		t.Errorf("ToCreateCalculationTypeParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateCalculationTypeParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateCalculationTypeParams(nil); got != (tables_model.CreateCalculationTypeParams{}) {
		t.Errorf("ToCreateCalculationTypeParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateCalculationTypeByIdParams(t *testing.T) {
	req := &tablesv1.UpdateCalculationTypeByIdRequest{
		Id:          100,
		Name:        "name-0",
		Description: "description-0",
	}

	want := tables_model.UpdateCalculationTypeByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
	}

	if got := ToUpdateCalculationTypeByIdParams(req); got != want {
		t.Errorf("ToUpdateCalculationTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateCalculationTypeByIdParamsNil(t *testing.T) {
	if got := ToUpdateCalculationTypeByIdParams(nil); got != (tables_model.UpdateCalculationTypeByIdParams{}) {
		t.Errorf("ToUpdateCalculationTypeByIdParams(nil) = %+v, want zero value", got)
	}
}
