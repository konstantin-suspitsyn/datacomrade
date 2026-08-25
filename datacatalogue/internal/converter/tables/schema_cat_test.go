package tables

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	schemaCatCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	schemaCatUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testSchemaCatRow — строка dc.schema_cat со значениями, различимыми между полями.
func testSchemaCatRow() tables_model.DcSchemaCat {
	return tables_model.DcSchemaCat{
		ID:         100,
		DatabaseID: 101,
		Name:       "name-0",
		IsDeleted:  false,
		CreatedAt:  schemaCatCreatedAt,
		UpdatedAt:  schemaCatUpdatedAt,
		UserID:     106,
	}
}

func TestSchemaCatToProto(t *testing.T) {
	row := testSchemaCatRow()
	got := SchemaCatToProto(row)

	if got == nil {
		t.Fatal("SchemaCatToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetDatabaseId() != row.DatabaseID {
		t.Errorf("DatabaseId = %d, want %d", got.GetDatabaseId(), row.DatabaseID)
	}

	if got.GetName() != row.Name {
		t.Errorf("Name = %q, want %q", got.GetName(), row.Name)
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

func TestSchemaCatsToProto(t *testing.T) {
	first := testSchemaCatRow()

	second := testSchemaCatRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcSchemaCat
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcSchemaCat{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcSchemaCat{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SchemaCatsToProto(tt.input)

			if got == nil {
				t.Fatal("SchemaCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateSchemaCatParams(t *testing.T) {
	req := &tablesv1.CreateSchemaCatRequest{
		DatabaseId: 100,
		Name:       "name-0",
		ExternalId: "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateSchemaCatParams{
		DatabaseID: 100,
		Name:       "name-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateSchemaCatParams(req); got != want {
		t.Errorf("ToCreateSchemaCatParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateSchemaCatParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateSchemaCatParams(nil); got != (tables_model.CreateSchemaCatParams{}) {
		t.Errorf("ToCreateSchemaCatParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateSchemaCatByIdParams(t *testing.T) {
	req := &tablesv1.UpdateSchemaCatByIdRequest{
		DatabaseId: 100,
		Name:       "name-0",
		CreatedAt:  timestamppb.New(schemaCatCreatedAt),
		ExternalId: "00000000-0000-4000-8000-000000000004",
		Id:         104,
	}

	want := tables_model.UpdateSchemaCatByIdParams{
		DatabaseID: 100,
		Name:       "name-0",
		CreatedAt:  schemaCatCreatedAt,
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000004"),
		ID:         104,
	}

	got := ToUpdateSchemaCatByIdParams(req)
	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("ToUpdateSchemaCatByIdParams.CreatedAt = %v, want %v", got.CreatedAt, want.CreatedAt)
	}
	got.CreatedAt, want.CreatedAt = time.Time{}, time.Time{}

	if got != want {
		t.Errorf("ToUpdateSchemaCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateSchemaCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateSchemaCatByIdParams(nil); got != (tables_model.UpdateSchemaCatByIdParams{}) {
		t.Errorf("ToUpdateSchemaCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteSchemaCatByIdParams(t *testing.T) {
	req := &tablesv1.DeleteSchemaCatByIdRequest{
		UpdatedAt:  timestamppb.New(schemaCatUpdatedAt),
		ExternalId: "00000000-0000-4000-8000-000000000002",
		Id:         102,
	}

	want := tables_model.DeleteSchemaCatByIdParams{
		UpdatedAt:  schemaCatUpdatedAt,
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000002"),
		ID:         102,
	}

	got := ToDeleteSchemaCatByIdParams(req)
	if !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("ToDeleteSchemaCatByIdParams.UpdatedAt = %v, want %v", got.UpdatedAt, want.UpdatedAt)
	}
	got.UpdatedAt, want.UpdatedAt = time.Time{}, time.Time{}

	if got != want {
		t.Errorf("ToDeleteSchemaCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteSchemaCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteSchemaCatByIdParams(nil); got != (tables_model.DeleteSchemaCatByIdParams{}) {
		t.Errorf("ToDeleteSchemaCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteSchemaCatByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteSchemaCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteSchemaCatByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteSchemaCatByIdParams(req); got != want {
		t.Errorf("ToUndeleteSchemaCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteSchemaCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteSchemaCatByIdParams(nil); got != (tables_model.UndeleteSchemaCatByIdParams{}) {
		t.Errorf("ToUndeleteSchemaCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetSchemaCatsDefaultsPageLimit(t *testing.T) {
	got := ToGetSchemaCatsParams(&tablesv1.GetSchemaCatsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetSchemaCatsDefaultsPage(t *testing.T) {
	got := ToGetSchemaCatsParams(&tablesv1.GetSchemaCatsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetSchemaCatsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetSchemaCatsParams(&tablesv1.GetSchemaCatsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetSchemaCatsDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetSchemaCatsDeletedParams(&tablesv1.GetSchemaCatsDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetSchemaCatsDeletedDefaultsPage(t *testing.T) {
	got := ToGetSchemaCatsDeletedParams(&tablesv1.GetSchemaCatsDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetSchemaCatsDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetSchemaCatsDeletedParams(&tablesv1.GetSchemaCatsDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
