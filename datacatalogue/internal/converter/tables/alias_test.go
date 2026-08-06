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

var (
	aliasCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	aliasUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testAliasRow — строка dc.alias со значениями, различимыми между полями.
func testAliasRow() tables_model.DcAlias {
	return tables_model.DcAlias{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		CreatedAt:   aliasCreatedAt,
		UpdatedAt:   sql.NullTime{Time: aliasUpdatedAt, Valid: true},
		IsDeleted:   false,
		UserID:      106,
	}
}

func TestAliasToProto(t *testing.T) {
	row := testAliasRow()
	got := AliasToProto(row)

	if got == nil {
		t.Fatal("AliasToProto() = nil, want value")
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

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt.Time) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt.Time)
	}

	if got.GetIsDeleted() != row.IsDeleted {
		t.Errorf("IsDeleted = %v, want %v", got.GetIsDeleted(), row.IsDeleted)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

}

func TestAliasesToProto(t *testing.T) {
	first := testAliasRow()

	second := testAliasRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcAlias
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcAlias{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcAlias{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AliasesToProto(tt.input)

			if got == nil {
				t.Fatal("AliasesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateAliasParams(t *testing.T) {
	req := &tablesv1.CreateAliasRequest{
		Name:        "name-0",
		Description: "description-0",
		ExternalId:  "00000000-0000-4000-8000-000000000003",
	}

	want := tables_model.CreateAliasParams{
		Name:        "name-0",
		Description: "description-0",
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000003"),
	}

	if got := ToCreateAliasParams(req); got != want {
		t.Errorf("ToCreateAliasParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateAliasParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateAliasParams(nil); got != (tables_model.CreateAliasParams{}) {
		t.Errorf("ToCreateAliasParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateAliasByIdParams(t *testing.T) {
	req := &tablesv1.UpdateAliasByIdRequest{
		Name:        "name-0",
		Description: "description-0",
		IsDeleted:   true,
		ExternalId:  "00000000-0000-4000-8000-000000000004",
		Id:          104,
	}

	want := tables_model.UpdateAliasByIdParams{
		Name:        "name-0",
		Description: "description-0",
		IsDeleted:   true,
		ExternalID:  uuid.MustParse("00000000-0000-4000-8000-000000000004"),
		ID:          104,
	}

	if got := ToUpdateAliasByIdParams(req); got != want {
		t.Errorf("ToUpdateAliasByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateAliasByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateAliasByIdParams(nil); got != (tables_model.UpdateAliasByIdParams{}) {
		t.Errorf("ToUpdateAliasByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteAliasByIdParams(t *testing.T) {
	req := &tablesv1.DeleteAliasByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteAliasByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteAliasByIdParams(req); got != want {
		t.Errorf("ToDeleteAliasByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteAliasByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteAliasByIdParams(nil); got != (tables_model.DeleteAliasByIdParams{}) {
		t.Errorf("ToDeleteAliasByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteAliasByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteAliasByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteAliasByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteAliasByIdParams(req); got != want {
		t.Errorf("ToUndeleteAliasByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteAliasByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteAliasByIdParams(nil); got != (tables_model.UndeleteAliasByIdParams{}) {
		t.Errorf("ToUndeleteAliasByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetAliasesDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetAliasesDeletedParams(&tablesv1.GetAliasesDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetAliasesDeletedDefaultsPage(t *testing.T) {
	got := ToGetAliasesDeletedParams(&tablesv1.GetAliasesDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetAliasesDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetAliasesDeletedParams(&tablesv1.GetAliasesDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetAliasesDefaultsPageLimit(t *testing.T) {
	got := ToGetAliasesParams(&tablesv1.GetAliasesRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetAliasesDefaultsPage(t *testing.T) {
	got := ToGetAliasesParams(&tablesv1.GetAliasesRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetAliasesKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetAliasesParams(&tablesv1.GetAliasesRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
