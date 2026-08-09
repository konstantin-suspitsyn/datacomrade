package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
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

			if got == nil {
				t.Fatal("CalculationTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}
