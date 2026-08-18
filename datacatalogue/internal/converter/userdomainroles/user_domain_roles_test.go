package userdomainroles

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

var (
	userDomainRolesCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	userDomainRolesUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testUserDomainRolesRow — строка dc.user_domain_roles со значениями, различимыми между полями.
func testUserDomainRolesRow() user_domain_roles.DcUserDomainRole {
	return user_domain_roles.DcUserDomainRole{
		ID:            100,
		UserID:        101,
		DomainRolesID: 102,
		CreatedAt:     userDomainRolesCreatedAt,
		UpdatedAt:     userDomainRolesUpdatedAt,
		IsDeleted:     false,
		DomainID:      106,
		UpdatedByID:   107,
	}
}

func TestUserDomainRolesToProto(t *testing.T) {
	row := testUserDomainRolesRow()
	got := UserDomainRolesToProto(row)

	if got == nil {
		t.Fatal("UserDomainRolesToProto() = nil, want value")
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

	if got.GetUpdatedById() != row.UpdatedByID {
		t.Errorf("UpdatedById = %d, want %d", got.GetUpdatedById(), row.UpdatedByID)
	}

}

func TestUserDomainRolesesToProto(t *testing.T) {
	first := testUserDomainRolesRow()

	second := testUserDomainRolesRow()
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
			got := UserDomainRolesesToProto(tt.input)

			if got == nil {
				t.Fatal("UserDomainRolesesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// testGetUserDomainRolesesRow — строка выборки GetUserDomainRolesesRow со значениями, различимыми между полями.
func testGetUserDomainRolesesRow() user_domain_roles.GetUserDomainRolesesRow {
	return user_domain_roles.GetUserDomainRolesesRow{
		ID:                     100,
		UserID:                 101,
		DomainRolesID:          102,
		CreatedAt:              userDomainRolesCreatedAt,
		UpdatedAt:              userDomainRolesUpdatedAt,
		IsDeleted:              false,
		DomainID:               106,
		UpdatedByID:            107,
		DomainRolesName:        "domain-roles-name-0",
		DomainRolesDescription: "domain-roles-description-0",
		DcUserName:             "dc-user-name-0",
	}
}

func TestGetUserDomainRolesesRowToProto(t *testing.T) {
	row := testGetUserDomainRolesesRow()
	got := GetUserDomainRolesesRowToProto(row)

	if got == nil {
		t.Fatal("GetUserDomainRolesesRowToProto() = nil, want value")
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

	if got.GetUpdatedById() != row.UpdatedByID {
		t.Errorf("UpdatedById = %d, want %d", got.GetUpdatedById(), row.UpdatedByID)
	}

	if got.GetDomainRolesName() != row.DomainRolesName {
		t.Errorf("DomainRolesName = %q, want %q", got.GetDomainRolesName(), row.DomainRolesName)
	}

	if got.GetDomainRolesDescription() != row.DomainRolesDescription {
		t.Errorf("DomainRolesDescription = %q, want %q", got.GetDomainRolesDescription(), row.DomainRolesDescription)
	}

	if got.GetDcUserName() != row.DcUserName {
		t.Errorf("DcUserName = %q, want %q", got.GetDcUserName(), row.DcUserName)
	}

}

func TestGetUserDomainRolesesRowsToProto(t *testing.T) {
	first := testGetUserDomainRolesesRow()

	second := testGetUserDomainRolesesRow()
	second.ID = 999
	second.DomainRolesName = "second-value"

	tests := []struct {
		name    string
		input   []user_domain_roles.GetUserDomainRolesesRow
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.GetUserDomainRolesesRow{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.GetUserDomainRolesesRow{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUserDomainRolesesRowsToProto(tt.input)

			if got == nil {
				t.Fatal("GetUserDomainRolesesRowsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateUserDomainRolesParams(t *testing.T) {
	req := &userdomainrolesv1.CreateUserDomainRolesRequest{
		UserId:        100,
		DomainRolesId: 101,
		DomainId:      102,
		ExternalId:    "00000000-0000-4000-8000-000000000004",
	}

	want := user_domain_roles.CreateUserDomainRolesParams{
		UserID:        100,
		DomainRolesID: 101,
		DomainID:      102,
		ExternalID:    uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToCreateUserDomainRolesParams(req); got != want {
		t.Errorf("ToCreateUserDomainRolesParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateUserDomainRolesParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateUserDomainRolesParams(nil); got != (user_domain_roles.CreateUserDomainRolesParams{}) {
		t.Errorf("ToCreateUserDomainRolesParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateUserDomainRolesByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateUserDomainRolesByIdRequest{
		UserId:        100,
		DomainRolesId: 101,
		DomainId:      102,
		ExternalId:    "00000000-0000-4000-8000-000000000004",
		Id:            104,
	}

	want := user_domain_roles.UpdateUserDomainRolesByIdParams{
		UserID:        100,
		DomainRolesID: 101,
		DomainID:      102,
		ExternalID:    uuid.MustParse("00000000-0000-4000-8000-000000000004"),
		ID:            104,
	}

	if got := ToUpdateUserDomainRolesByIdParams(req); got != want {
		t.Errorf("ToUpdateUserDomainRolesByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateUserDomainRolesByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateUserDomainRolesByIdParams(nil); got != (user_domain_roles.UpdateUserDomainRolesByIdParams{}) {
		t.Errorf("ToUpdateUserDomainRolesByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteUserDomainRolesByIdParams(t *testing.T) {
	req := &userdomainrolesv1.DeleteUserDomainRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := user_domain_roles.DeleteUserDomainRolesByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteUserDomainRolesByIdParams(req); got != want {
		t.Errorf("ToDeleteUserDomainRolesByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteUserDomainRolesByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteUserDomainRolesByIdParams(nil); got != (user_domain_roles.DeleteUserDomainRolesByIdParams{}) {
		t.Errorf("ToDeleteUserDomainRolesByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteUserDomainRolesByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UndeleteUserDomainRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := user_domain_roles.UndeleteUserDomainRolesByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteUserDomainRolesByIdParams(req); got != want {
		t.Errorf("ToUndeleteUserDomainRolesByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteUserDomainRolesByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteUserDomainRolesByIdParams(nil); got != (user_domain_roles.UndeleteUserDomainRolesByIdParams{}) {
		t.Errorf("ToUndeleteUserDomainRolesByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetUserDomainRolesesDefaultsPageLimit(t *testing.T) {
	got := ToGetUserDomainRolesesParams(&userdomainrolesv1.GetUserDomainRolesesRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetUserDomainRolesesDefaultsPage(t *testing.T) {
	got := ToGetUserDomainRolesesParams(&userdomainrolesv1.GetUserDomainRolesesRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetUserDomainRolesesKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetUserDomainRolesesParams(&userdomainrolesv1.GetUserDomainRolesesRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}
