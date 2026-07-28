package tables

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
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

func TestDomainCatToProtoDeleted(t *testing.T) {
	row := testDomainCatRow()
	row.IsDeleted = true

	if got := DomainCatToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
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

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("DomainCatsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestDomainCatsToProtoKeepsOrder(t *testing.T) {
	first := testDomainCatRow()
	second := testDomainCatRow()
	second.DomainName = "second-value"

	got := DomainCatsToProto([]tables_model.DcDomainCat{first, second})

	if got[0].GetDomainName() != first.DomainName {
		t.Errorf("[0] = %q, want %q", got[0].GetDomainName(), first.DomainName)
	}

	if got[1].GetDomainName() != second.DomainName {
		t.Errorf("[1] = %q, want %q", got[1].GetDomainName(), second.DomainName)
	}
}

func TestToCreateDomainCatParams(t *testing.T) {
	req := &tablesv1.CreateDomainCatRequest{
		DomainName: "domain-name-0",
		UserId:     101,
	}

	want := tables_model.CreateDomainCatParams{
		DomainName: "domain-name-0",
		UserID:     101,
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
		Id:         100,
		DomainName: "domain-name-0",
		UserId:     102,
	}

	want := tables_model.UpdateDomainCatByIdParams{
		ID:         100,
		DomainName: "domain-name-0",
		UserID:     102,
	}

	if got := ToUpdateDomainCatByIdParams(req); got != want {
		t.Errorf("ToUpdateDomainCatByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateDomainCatByIdParamsNil(t *testing.T) {
	if got := ToUpdateDomainCatByIdParams(nil); got != (tables_model.UpdateDomainCatByIdParams{}) {
		t.Errorf("ToUpdateDomainCatByIdParams(nil) = %+v, want zero value", got)
	}
}
