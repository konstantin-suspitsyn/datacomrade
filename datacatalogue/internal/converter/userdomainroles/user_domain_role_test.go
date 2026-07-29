package userdomainroles

import (
	"testing"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

var (
	userDomainRoleCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	userDomainRoleUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testUserDomainRoleRow — строка dc.user_domain_roles со значениями, различимыми между полями.
func testUserDomainRoleRow() user_domain_roles.DcUserDomainRole {
	return user_domain_roles.DcUserDomainRole{
		ID:            100,
		UserID:        101,
		DomainRolesID: 102,
		CreatedAt:     userDomainRoleCreatedAt,
		UpdatedAt:     userDomainRoleUpdatedAt,
		IsDeleted:     false,
		DomainID:      106,
	}
}

func TestUserDomainRoleToProto(t *testing.T) {
	row := testUserDomainRoleRow()
	got := UserDomainRoleToProto(row)

	if got == nil {
		t.Fatal("UserDomainRoleToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

	if got.GetDomainRolesId() != row.DomainRolesID {
		t.Errorf("DomainRolesId = %d, want %d", got.GetDomainRolesId(), row.DomainRolesID)
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

	if got.GetDomainId() != row.DomainID {
		t.Errorf("DomainId = %d, want %d", got.GetDomainId(), row.DomainID)
	}

}

func TestUserDomainRoleToProtoDeleted(t *testing.T) {
	row := testUserDomainRoleRow()
	row.IsDeleted = true

	if got := UserDomainRoleToProto(row); !got.GetIsDeleted() {
		t.Error("IsDeleted = false, want true")
	}
}

func TestUserDomainRolesToProto(t *testing.T) {
	first := testUserDomainRoleRow()

	second := testUserDomainRoleRow()
	second.ID = 999

	tests := []struct {
		name    string
		input   []user_domain_roles.DcUserDomainRole
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.DcUserDomainRole{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.DcUserDomainRole{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserDomainRolesToProto(tt.input)

			// Пустой вход даёт пустой, а не nil-слайс.
			if got == nil {
				t.Fatal("UserDomainRolesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestUserDomainRolesToProtoKeepsOrder(t *testing.T) {
	first := testUserDomainRoleRow()
	second := testUserDomainRoleRow()
	second.ID = 999

	got := UserDomainRolesToProto([]user_domain_roles.DcUserDomainRole{first, second})

	if got[0].GetId() != first.ID {
		t.Errorf("[0] = %d, want %d", got[0].GetId(), first.ID)
	}

	if got[1].GetId() != second.ID {
		t.Errorf("[1] = %d, want %d", got[1].GetId(), second.ID)
	}
}

func TestToCreateUserDomainRoleParams(t *testing.T) {
	req := &userdomainrolesv1.CreateUserDomainRoleRequest{
		UserId:        100,
		DomainRolesId: 101,
		DomainId:      102,
	}

	want := user_domain_roles.CreateUserDomainRoleParams{
		UserID:        100,
		DomainRolesID: 101,
		DomainID:      102,
	}

	if got := ToCreateUserDomainRoleParams(req); got != want {
		t.Errorf("ToCreateUserDomainRoleParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateUserDomainRoleParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateUserDomainRoleParams(nil); got != (user_domain_roles.CreateUserDomainRoleParams{}) {
		t.Errorf("ToCreateUserDomainRoleParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateUserDomainRoleByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateUserDomainRoleByIdRequest{
		Id:            100,
		UserId:        101,
		DomainRolesId: 102,
		DomainId:      103,
	}

	want := user_domain_roles.UpdateUserDomainRoleByIdParams{
		ID:            100,
		UserID:        101,
		DomainRolesID: 102,
		DomainID:      103,
	}

	if got := ToUpdateUserDomainRoleByIdParams(req); got != want {
		t.Errorf("ToUpdateUserDomainRoleByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateUserDomainRoleByIdParamsNil(t *testing.T) {
	if got := ToUpdateUserDomainRoleByIdParams(nil); got != (user_domain_roles.UpdateUserDomainRoleByIdParams{}) {
		t.Errorf("ToUpdateUserDomainRoleByIdParams(nil) = %+v, want zero value", got)
	}
}
