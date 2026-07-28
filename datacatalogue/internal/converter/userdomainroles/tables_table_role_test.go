package userdomainroles

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

var (
	tablesTableRoleCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	tablesTableRoleUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testTablesTableRoleRow — строка dc.tables_table_roles со значениями, различимыми между полями.
func testTablesTableRoleRow() user_domain_roles.DcTablesTableRole {
	return user_domain_roles.DcTablesTableRole{
		ID:           100,
		TableCatID:   101,
		TableRolesID: 102,
		CreatedAt:    tablesTableRoleCreatedAt,
		UpdatedAt:    tablesTableRoleUpdatedAt,
		IsDeleted:    false,
	}
}

func TestTablesTableRoleToProto(t *testing.T) {
	row := testTablesTableRoleRow()
	got := TablesTableRoleToProto(row)

	if got == nil {
		t.Fatal("TablesTableRoleToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetTableCatId() != row.TableCatID {
		t.Errorf("TableCatId = %d, want %d", got.GetTableCatId(), row.TableCatID)
	}

	if got.GetTableRolesId() != row.TableRolesID {
		t.Errorf("TableRolesId = %d, want %d", got.GetTableRolesId(), row.TableRolesID)
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

func TestTablesTableRoleToProtoDeleted(t *testing.T) {
	row := testTablesTableRoleRow()
	row.IsDeleted = true

	if got := TablesTableRoleToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestTablesTableRolesToProto(t *testing.T) {
	first := testTablesTableRoleRow()

	second := testTablesTableRoleRow()
	second.ID = 999

	tests := []struct {
		name    string
		input   []user_domain_roles.DcTablesTableRole
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.DcTablesTableRole{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.DcTablesTableRole{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TablesTableRolesToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("TablesTableRolesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestTablesTableRolesToProtoKeepsOrder(t *testing.T) {
	first := testTablesTableRoleRow()
	second := testTablesTableRoleRow()
	second.ID = 999

	got := TablesTableRolesToProto([]user_domain_roles.DcTablesTableRole{first, second})

	if got[0].GetId() != first.ID {
		t.Errorf("[0] = %d, want %d", got[0].GetId(), first.ID)
	}

	if got[1].GetId() != second.ID {
		t.Errorf("[1] = %d, want %d", got[1].GetId(), second.ID)
	}
}

func TestToCreateTablesTableRoleParams(t *testing.T) {
	req := &userdomainrolesv1.CreateTablesTableRoleRequest{
		TableCatId:   100,
		TableRolesId: 101,
	}

	want := user_domain_roles.CreateTablesTableRoleParams{
		TableCatID:   100,
		TableRolesID: 101,
	}

	if got := ToCreateTablesTableRoleParams(req); got != want {
		t.Errorf("ToCreateTablesTableRoleParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateTablesTableRoleParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateTablesTableRoleParams(nil); got != (user_domain_roles.CreateTablesTableRoleParams{}) {
		t.Errorf("ToCreateTablesTableRoleParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateTablesTableRoleByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateTablesTableRoleByIdRequest{
		Id:           100,
		TableCatId:   101,
		TableRolesId: 102,
	}

	want := user_domain_roles.UpdateTablesTableRoleByIdParams{
		ID:           100,
		TableCatID:   101,
		TableRolesID: 102,
	}

	if got := ToUpdateTablesTableRoleByIdParams(req); got != want {
		t.Errorf("ToUpdateTablesTableRoleByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateTablesTableRoleByIdParamsNil(t *testing.T) {
	if got := ToUpdateTablesTableRoleByIdParams(nil); got != (user_domain_roles.UpdateTablesTableRoleByIdParams{}) {
		t.Errorf("ToUpdateTablesTableRoleByIdParams(nil) = %+v, want zero value", got)
	}
}
