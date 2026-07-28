package tables

import (
	"database/sql"
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
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

func TestAliasToProtoDeleted(t *testing.T) {
	row := testAliasRow()
	row.IsDeleted = true

	if got := AliasToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
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

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("AliasesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestAliasesToProtoKeepsOrder(t *testing.T) {
	first := testAliasRow()
	second := testAliasRow()
	second.Name = "second-value"

	got := AliasesToProto([]tables_model.DcAlias{first, second})

	if got[0].GetName() != first.Name {
		t.Errorf("[0] = %q, want %q", got[0].GetName(), first.Name)
	}

	if got[1].GetName() != second.Name {
		t.Errorf("[1] = %q, want %q", got[1].GetName(), second.Name)
	}
}

func TestToCreateAliasParams(t *testing.T) {
	req := &tablesv1.CreateAliasRequest{
		Name:        "name-0",
		Description: "description-0",
		UserId:      102,
	}

	want := tables_model.CreateAliasParams{
		Name:        "name-0",
		Description: "description-0",
		UserID:      102,
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
		Id:          100,
		Name:        "name-0",
		Description: "description-0",
		UserId:      103,
	}

	want := tables_model.UpdateAliasByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		UserID:      103,
	}

	if got := ToUpdateAliasByIdParams(req); got != want {
		t.Errorf("ToUpdateAliasByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateAliasByIdParamsNil(t *testing.T) {
	if got := ToUpdateAliasByIdParams(nil); got != (tables_model.UpdateAliasByIdParams{}) {
		t.Errorf("ToUpdateAliasByIdParams(nil) = %+v, want zero value", got)
	}
}
