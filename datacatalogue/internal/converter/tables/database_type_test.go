package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

var (
	databaseTypeCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	databaseTypeUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testDatabaseTypeRow — строка dc.database_type со значениями, различимыми между полями.
func testDatabaseTypeRow() tables_model.DcDatabaseType {
	return tables_model.DcDatabaseType{
		ID:        100,
		Name:      "name-0",
		DbVersion: "db-version-0",
		IsDeleted: false,
		CreatedAt: databaseTypeCreatedAt,
		UpdatedAt: databaseTypeUpdatedAt,
		UserID:    106,
	}
}

func TestDatabaseTypeToProto(t *testing.T) {
	row := testDatabaseTypeRow()
	got := DatabaseTypeToProto(row)

	if got == nil {
		t.Fatal("DatabaseTypeToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
	}

	if got.GetDbVersion() != row.DbVersion {
		t.Errorf("DbVersion = %q, want %q", got.GetDbVersion(), row.DbVersion)
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

func TestDatabaseTypeToProtoDeleted(t *testing.T) {
	row := testDatabaseTypeRow()
	row.IsDeleted = true

	if got := DatabaseTypeToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestDatabaseTypesToProto(t *testing.T) {
	first := testDatabaseTypeRow()

	second := testDatabaseTypeRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcDatabaseType
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcDatabaseType{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcDatabaseType{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DatabaseTypesToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("DatabaseTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestDatabaseTypesToProtoKeepsOrder(t *testing.T) {
	first := testDatabaseTypeRow()
	second := testDatabaseTypeRow()
	second.Name = "second-value"

	got := DatabaseTypesToProto([]tables_model.DcDatabaseType{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateDatabaseTypeParams(t *testing.T) {
	req := &tablesv1.CreateDatabaseTypeRequest{
		Name:           "name-0",
		DbVersion:      "db-version-0",
		UserExternalId: "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateDatabaseTypeParams{
		Name:       "name-0",
		DbVersion:  "db-version-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateDatabaseTypeParams(req); got != want {
		t.Errorf("ToCreateDatabaseTypeParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateDatabaseTypeParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateDatabaseTypeParams(nil); got != (tables_model.CreateDatabaseTypeParams{}) {
		t.Errorf("ToCreateDatabaseTypeParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateDatabaseTypeByIdParams(t *testing.T) {
	req := &tablesv1.UpdateDatabaseTypeByIdRequest{
		Id:             100,
		Name:           "name-0",
		DbVersion:      "db-version-0",
		UserExternalId: "00000000-0000-4000-8000-000000000004",
	}

	want := tables_model.UpdateDatabaseTypeByIdParams{
		ID:         100,
		Name:       "name-0",
		DbVersion:  "db-version-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToUpdateDatabaseTypeByIdParams(req); got != want {
		t.Errorf("ToUpdateDatabaseTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDatabaseTypeByIdParamsNil(t *testing.T) {
	if got := ToUpdateDatabaseTypeByIdParams(nil); got != (tables_model.UpdateDatabaseTypeByIdParams{}) {
		t.Errorf("ToUpdateDatabaseTypeByIdParams(nil) = %+v, want zero value", got)
	}
}
