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
	domainCatCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	domainCatUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testDomainCatRow — строка dc.domain_cat со значениями, различимыми между полями.
func testDomainCatRow() tables_model.DcDomainCat {
	return tables_model.DcDomainCat{
		ID:         100,
		DomainName: "domain-name-0",
		IsDeleted:  false,
		CreatedAt:  domainCatCreatedAt,
		UpdatedAt:  domainCatUpdatedAt,
		UserID:     105,
	}
}

func TestDomainCatToProto(t *testing.T) {
	row := testDomainCatRow()
	got := DomainCatToProto(row)

	if got == nil {
		t.Fatal("DomainCatToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetDomainName() != row.DomainName {
		t.Errorf("DomainName = %q, want %q", got.GetDomainName(), row.DomainName)
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

func TestDomainCatsToProto(t *testing.T) {
	first := testDomainCatRow()

	second := testDomainCatRow()
	second.ID = 999
	second.DomainName = "second-value"

	tests := []struct {
		name    string
		input   []tables_model.DcDomainCat
		wantLen int
	}{
		{name: "two rows", input: []tables_model.DcDomainCat{first, second}, wantLen: 2},
		{name: "empty slice", input: []tables_model.DcDomainCat{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DomainCatsToProto(tt.input)

			if got == nil {
				t.Fatal("DomainCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateDomainCatParams(t *testing.T) {
	req := &tablesv1.CreateDomainCatRequest{
		DomainName: "domain-name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
	}

	want := tables_model.CreateDomainCatParams{
		DomainName: "domain-name-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000002"),
	}

	if got := ToCreateDomainCatParams(req); got != want {
		t.Errorf("ToCreateDomainCatParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateDomainCatParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateDomainCatParams(nil); got != (tables_model.CreateDomainCatParams{}) {
		t.Errorf("ToCreateDomainCatParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateDomainCatByIdParams(t *testing.T) {
	req := &tablesv1.UpdateDomainCatByIdRequest{
		DomainName: "domain-name-0",
		ExternalId: "00000000-0000-4000-8000-000000000002",
		Id:         102,
	}

	want := tables_model.UpdateDomainCatByIdParams{
		DomainName: "domain-name-0",
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000002"),
		ID:         102,
	}

	if got := ToUpdateDomainCatByIdParams(req); got != want {
		t.Errorf("ToUpdateDomainCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDomainCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateDomainCatByIdParams(nil); got != (tables_model.UpdateDomainCatByIdParams{}) {
		t.Errorf("ToUpdateDomainCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteDomainCatByIdParams(t *testing.T) {
	req := &tablesv1.DeleteDomainCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.DeleteDomainCatByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteDomainCatByIdParams(req); got != want {
		t.Errorf("ToDeleteDomainCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteDomainCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteDomainCatByIdParams(nil); got != (tables_model.DeleteDomainCatByIdParams{}) {
		t.Errorf("ToDeleteDomainCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteDomainCatByIdParams(t *testing.T) {
	req := &tablesv1.UndeleteDomainCatByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := tables_model.UndeleteDomainCatByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteDomainCatByIdParams(req); got != want {
		t.Errorf("ToUndeleteDomainCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteDomainCatByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteDomainCatByIdParams(nil); got != (tables_model.UndeleteDomainCatByIdParams{}) {
		t.Errorf("ToUndeleteDomainCatByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetDomainCatsSearchNameDefaultsPageLimit(t *testing.T) {
	got := ToGetDomainCatsSearchNameParams(&tablesv1.GetDomainCatsSearchNameRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDomainCatsSearchNameDefaultsPage(t *testing.T) {
	got := ToGetDomainCatsSearchNameParams(&tablesv1.GetDomainCatsSearchNameRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDomainCatsSearchNameKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDomainCatsSearchNameParams(&tablesv1.GetDomainCatsSearchNameRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDomainCatsSearchNamePassesFilterFields(t *testing.T) {
	req := &tablesv1.GetDomainCatsSearchNameRequest{PageLimit: 10, Page: 1,
		SearchName: "SearchName-0",
	}

	got := ToGetDomainCatsSearchNameParams(req)

	if got.SearchName != "SearchName-0" {
		t.Errorf("SearchName = %q, want %q", got.SearchName, "SearchName-0")
	}
}

func TestGetDomainCatsDefaultsPageLimit(t *testing.T) {
	got := ToGetDomainCatsParams(&tablesv1.GetDomainCatsRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDomainCatsDefaultsPage(t *testing.T) {
	got := ToGetDomainCatsParams(&tablesv1.GetDomainCatsRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDomainCatsKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDomainCatsParams(&tablesv1.GetDomainCatsRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetDomainCatsDeletedDefaultsPageLimit(t *testing.T) {
	got := ToGetDomainCatsDeletedParams(&tablesv1.GetDomainCatsDeletedRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetDomainCatsDeletedDefaultsPage(t *testing.T) {
	got := ToGetDomainCatsDeletedParams(&tablesv1.GetDomainCatsDeletedRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetDomainCatsDeletedKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetDomainCatsDeletedParams(&tablesv1.GetDomainCatsDeletedRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
