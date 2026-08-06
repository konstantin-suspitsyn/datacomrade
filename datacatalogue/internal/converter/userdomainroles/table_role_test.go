package userdomainroles

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

var (
	tableRoleCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	tableRoleUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testTableRoleRow — строка dc.table_roles со значениями, различимыми между полями.
func testTableRoleRow() user_domain_roles.DcTableRole {
	return user_domain_roles.DcTableRole{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
		CreatedAt:   tableRoleCreatedAt,
		UpdatedAt:   tableRoleUpdatedAt,
		IsDeleted:   false,
	}
}

func TestTableRoleToProto(t *testing.T) {
	row := testTableRoleRow()
	got := TableRoleToProto(row)

	if got == nil {
		t.Fatal("TableRoleToProto() = nil, want value")
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

func TestTableRolesToProto(t *testing.T) {
	first := testTableRoleRow()

	second := testTableRoleRow()
	second.ID = 999
	second.Name = "second-value"

	tests := []struct {
		name    string
		input   []user_domain_roles.DcTableRole
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.DcTableRole{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.DcTableRole{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TableRolesToProto(tt.input)

			if got == nil {
				t.Fatal("TableRolesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateTableRoleParams(t *testing.T) {
	req := &userdomainrolesv1.CreateTableRoleRequest{
		Name:        "name-0",
		Description: "description-0",
	}

	want := user_domain_roles.CreateTableRoleParams{
		Name:        "name-0",
		Description: "description-0",
	}

	if got := ToCreateTableRoleParams(req); got != want {
		t.Errorf("ToCreateTableRoleParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateTableRoleParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateTableRoleParams(nil); got != (user_domain_roles.CreateTableRoleParams{}) {
		t.Errorf("ToCreateTableRoleParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateTableRoleByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateTableRoleByIdRequest{
		Id:          100,
		Name:        "name-0",
		Description: "description-0",
	}

	want := user_domain_roles.UpdateTableRoleByIdParams{
		ID:          100,
		Name:        "name-0",
		Description: "description-0",
	}

	if got := ToUpdateTableRoleByIdParams(req); got != want {
		t.Errorf("ToUpdateTableRoleByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateTableRoleByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateTableRoleByIdParams(nil); got != (user_domain_roles.UpdateTableRoleByIdParams{}) {
		t.Errorf("ToUpdateTableRoleByIdParams(nil) = %+v, want zero value", got)
	}
}
