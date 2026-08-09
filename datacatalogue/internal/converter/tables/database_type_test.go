package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
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

			if got == nil {
				t.Fatal("DatabaseTypesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateDatabaseTypeParams(t *testing.T) {
	req := &tablesv1.CreateDatabaseTypeRequest{
		Name:       "name-0",
		DbVersion:  "db-version-0",
		ExternalId: "00000000-0000-4000-8000-000000000003",
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
		Name:       "name-0",
		DbVersion:  "db-version-0",
		ExternalId: "00000000-0000-4000-8000-000000000003",
		Id:         103,
	}

	want := tables_model.UpdateDatabaseTypeByIdParams{
		Name:       "name-0",
		DbVersion:  "db-version-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000003"),
		ID:         103,
	}

	if got := ToUpdateDatabaseTypeByIdParams(req); got != want {
		t.Errorf("ToUpdateDatabaseTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDatabaseTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateDatabaseTypeByIdParams(nil); got != (tables_model.UpdateDatabaseTypeByIdParams{}) {
		t.Errorf("ToUpdateDatabaseTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteDatabaseTypeByIdParams(t *testing.T) {
	req := &tablesv1.DeleteDatabaseTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteDatabaseTypeByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteDatabaseTypeByIdParams(req); got != want {
		t.Errorf("ToDeleteDatabaseTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteDatabaseTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteDatabaseTypeByIdParams(nil); got != (tables_model.DeleteDatabaseTypeByIdParams{}) {
		t.Errorf("ToDeleteDatabaseTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteDatabaseTypeByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteDatabaseTypeByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteDatabaseTypeByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteDatabaseTypeByIdParams(req); got != want {
		t.Errorf("ToUndeleteDatabaseTypeByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteDatabaseTypeByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteDatabaseTypeByIdParams(nil); got != (tables_model.UndeleteDatabaseTypeByIdParams{}) {
		t.Errorf("ToUndeleteDatabaseTypeByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetDatabaseTypesSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseTypesSearchNameParams(&tablesv1.GetDatabaseTypesSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseTypesSearchNameDefaultsPage(t *testing.T) {
	got := ToGetDatabaseTypesSearchNameParams(&tablesv1.GetDatabaseTypesSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseTypesSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseTypesSearchNameParams(&tablesv1.GetDatabaseTypesSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseTypesSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetDatabaseTypesSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetDatabaseTypesSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}

func TestGetDatabaseTypesDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseTypesParams(&tablesv1.GetDatabaseTypesRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseTypesDefaultsPage(t *testing.T) {
	got := ToGetDatabaseTypesParams(&tablesv1.GetDatabaseTypesRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseTypesKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseTypesParams(&tablesv1.GetDatabaseTypesRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDatabaseTypesDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetDatabaseTypesDeletedParams(&tablesv1.GetDatabaseTypesDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDatabaseTypesDeletedDefaultsPage(t *testing.T) {
	got := ToGetDatabaseTypesDeletedParams(&tablesv1.GetDatabaseTypesDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDatabaseTypesDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDatabaseTypesDeletedParams(&tablesv1.GetDatabaseTypesDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
