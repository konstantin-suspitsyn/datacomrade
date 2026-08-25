package olapexchange

import (
	"testing"

	"github.com/google/uuid"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/olapexchange"
	olapexchangev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/olap_exchange/v1"
)

func TestToGetOlapAliasesParams(t *testing.T) {
	req := &olapexchangev1.GetOlapAliasesRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Order:      "ASC",
	}

	want := olapexchange.GetOlapAliasesParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		Order:      "ASC",
	}

	if got := ToGetOlapAliasesParams(req); got != want {
		t.Errorf("ToGetOlapAliasesParams() = %+v, want %+v", got, want)
	}
}

func TestToGetOlapAliasesParamsNil(t *testing.T) {
	if got := ToGetOlapAliasesParams(nil); got != (olapexchange.GetOlapAliasesParams{}) {
		t.Errorf("ToGetOlapAliasesParams(nil) = %+v, want zero value", got)
	}
}

// testGetOlapAliasesRowRow — строка запроса GetOlapAliasesRow со значениями, различимыми между полями.
func testGetOlapAliasesRowRow() olapexchange.GetOlapAliasesRow {
	return olapexchange.GetOlapAliasesRow{
		ID:                    100,
		Name:                  "name-1",
		Description:           "description-2",
		IsDeleted:             true,
		ColumnTypeID:          104,
		ColumnTypeName:        "columntypename-5",
		ColumnTypeDescription: "columntypedescription-6",
	}
}

func TestGetOlapAliasesRowToProto(t *testing.T) {
	row := testGetOlapAliasesRowRow()
	got := GetOlapAliasesRowToProto(row)

	if got == nil {
		t.Fatal("GetOlapAliasesRowToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %v, want %v", got.GetId(), row.ID)
	}
	if got.GetName() != row.Name {
		t.Errorf("Name = %v, want %v", got.GetName(), row.Name)
	}
	if got.GetDescription() != row.Description {
		t.Errorf("Description = %v, want %v", got.GetDescription(), row.Description)
	}
	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}
	if got.GetColumnTypeId() != row.ColumnTypeID {
		t.Errorf("ColumnTypeId = %v, want %v", got.GetColumnTypeId(), row.ColumnTypeID)
	}
	if got.GetColumnTypeName() != row.ColumnTypeName {
		t.Errorf("ColumnTypeName = %v, want %v", got.GetColumnTypeName(), row.ColumnTypeName)
	}
	if got.GetColumnTypeDescription() != row.ColumnTypeDescription {
		t.Errorf("ColumnTypeDescription = %v, want %v", got.GetColumnTypeDescription(), row.ColumnTypeDescription)
	}
}

func TestGetOlapAliasesToProto(t *testing.T) {
	got := GetOlapAliasesToProto([]olapexchange.GetOlapAliasesRow{testGetOlapAliasesRowRow()})

	if got == nil {
		t.Fatal("GetOlapAliasesToProto() = nil, want value")
	}

	if len(got.GetRows()) != 1 {
		t.Errorf("len(Rows) = %d, want 1", len(got.GetRows()))
	}
}
