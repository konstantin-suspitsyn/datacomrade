package userdomainroles

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

var (
	userTableRoleCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	userTableRoleUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testUserTableRoleRow — строка dc.user_table_roles со значениями, различимыми между полями.
func testUserTableRoleRow() user_domain_roles.DcUserTableRole {
	return user_domain_roles.DcUserTableRole{
		ID:           100,
		UserID:       101,
		TableRolesID: 102,
		CreatedAt:    userTableRoleCreatedAt,
		UpdatedAt:    userTableRoleUpdatedAt,
		IsDeleted:    false,
		TableID:      106,
		UpdatedByID:  107,
	}
}

func TestUserTableRoleToProto(t *testing.T) {
	row := testUserTableRoleRow()
	got := UserTableRoleToProto(row)

	if got == nil {
		t.Fatal("UserTableRoleToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
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

	if got.GetTableId() != row.TableID {
		t.Errorf("TableId = %d, want %d", got.GetTableId(), row.TableID)
	}

	if got.GetUpdatedById() != row.UpdatedByID {
		t.Errorf("UpdatedById = %d, want %d", got.GetUpdatedById(), row.UpdatedByID)
	}

}

func TestUserTableRolesToProto(t *testing.T) {
	first := testUserTableRoleRow()

	second := testUserTableRoleRow()
	second.ID = 999

	tests := []struct {
		name    string
		input   []user_domain_roles.DcUserTableRole
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.DcUserTableRole{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.DcUserTableRole{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserTableRolesToProto(tt.input)

			if got == nil {
				t.Fatal("UserTableRolesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateUserTableRoleParams(t *testing.T) {
	req := &userdomainrolesv1.CreateUserTableRoleRequest{
		UserId:              100,
		TableRolesId:        101,
		TableId:             102,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000004",
	}

	want := user_domain_roles.CreateUserTableRoleParams{
		UserID:       100,
		TableRolesID: 101,
		TableID:      102,
		ExternalID:   uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToCreateUserTableRoleParams(req); got != want {
		t.Errorf("ToCreateUserTableRoleParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateUserTableRoleParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateUserTableRoleParams(nil); got != (user_domain_roles.CreateUserTableRoleParams{}) {
		t.Errorf("ToCreateUserTableRoleParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateUserTableRoleByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateUserTableRoleByIdRequest{
		Id:                  100,
		UserId:              101,
		TableRolesId:        102,
		TableId:             103,
		UpdatedByExternalId: "00000000-0000-4000-8000-000000000005",
	}

	want := user_domain_roles.UpdateUserTableRoleByIdParams{
		ID:           100,
		UserID:       101,
		TableRolesID: 102,
		TableID:      103,
		ExternalID:   uuid.MustParse("00000000-0000-4000-8000-000000000005"),
	}

	if got := ToUpdateUserTableRoleByIdParams(req); got != want {
		t.Errorf("ToUpdateUserTableRoleByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateUserTableRoleByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateUserTableRoleByIdParams(nil); got != (user_domain_roles.UpdateUserTableRoleByIdParams{}) {
		t.Errorf("ToUpdateUserTableRoleByIdParams(nil) = %+v, want zero value", got)
	}
}
