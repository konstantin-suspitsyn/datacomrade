package userdomainroles

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// userTableRolesPtr — адрес значения для optional-полей proto.
func userTableRolesPtr[T any](v T) *T {
	return &v
}

var (
	userTableRolesCreatedAt = time.Date(2026, time.July, 1, 9, 0, 0, 0, time.UTC)
	userTableRolesUpdatedAt = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)
)

// testUserTableRolesRow — строка dc.user_table_roles со значениями, различимыми между полями.
func testUserTableRolesRow() user_domain_roles.DcUserTableRole {
	return user_domain_roles.DcUserTableRole{
		ID:           100,
		UserID:       101,
		TableRolesID: 102,
		CreatedAt:    userTableRolesCreatedAt,
		UpdatedAt:    userTableRolesUpdatedAt,
		IsDeleted:    false,
		TableID:      106,
		UpdatedByID:  107,
	}
}

func TestUserTableRolesToProto(t *testing.T) {
	row := testUserTableRolesRow()
	got := UserTableRolesToProto(row)

	if got == nil {
		t.Fatal("UserTableRolesToProto() = nil, want value")
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

func TestUserTableRolesesToProto(t *testing.T) {
	first := testUserTableRolesRow()

	second := testUserTableRolesRow()
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
			got := UserTableRolesesToProto(tt.input)

			if got == nil {
				t.Fatal("UserTableRolesesToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// testGetUserTableRolesByIdRow — строка выборки GetUserTableRolesByIdRow со значениями, различимыми между полями.
func testGetUserTableRolesByIdRow() user_domain_roles.GetUserTableRolesByIdRow {
	return user_domain_roles.GetUserTableRolesByIdRow{
		ID:                    100,
		UserID:                101,
		TableRolesID:          102,
		CreatedAt:             userTableRolesCreatedAt,
		UpdatedAt:             userTableRolesUpdatedAt,
		TableID:               105,
		DcUserID:              106,
		DcUserName:            "dc-user-name-0",
		DcUserExternalID:      uuid.MustParse("00000000-0000-4000-8000-000000000009"),
		TableRolesName:        "table-roles-name-0",
		TableRolesDescription: "table-roles-description-0",
		TableCatID:            111,
		TableCatName:          "table-cat-name-0",
	}
}

func TestGetUserTableRolesByIdRowToProto(t *testing.T) {
	row := testGetUserTableRolesByIdRow()
	got := GetUserTableRolesByIdRowToProto(row)

	if got == nil {
		t.Fatal("GetUserTableRolesByIdRowToProto() = nil, want value")
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

	if got.GetTableId() != row.TableID {
		t.Errorf("TableId = %d, want %d", got.GetTableId(), row.TableID)
	}

	if got.GetDcUserId() != row.DcUserID {
		t.Errorf("DcUserId = %d, want %d", got.GetDcUserId(), row.DcUserID)
	}

	if got.GetDcUserName() != row.DcUserName {
		t.Errorf("DcUserName = %q, want %q", got.GetDcUserName(), row.DcUserName)
	}

	if got.GetDcUserExternalId() != row.DcUserExternalID.String() {
		t.Errorf("DcUserExternalId = %q, want %q", got.GetDcUserExternalId(), row.DcUserExternalID.String())
	}

	if got.GetTableRolesName() != row.TableRolesName {
		t.Errorf("TableRolesName = %q, want %q", got.GetTableRolesName(), row.TableRolesName)
	}

	if got.GetTableRolesDescription() != row.TableRolesDescription {
		t.Errorf("TableRolesDescription = %q, want %q", got.GetTableRolesDescription(), row.TableRolesDescription)
	}

	if got.GetTableCatId() != row.TableCatID {
		t.Errorf("TableCatId = %d, want %d", got.GetTableCatId(), row.TableCatID)
	}

	if got.GetTableCatName() != row.TableCatName {
		t.Errorf("TableCatName = %q, want %q", got.GetTableCatName(), row.TableCatName)
	}

}

// testGetUserTableRolesesRow — строка выборки GetUserTableRolesesRow со значениями, различимыми между полями.
func testGetUserTableRolesesRow() user_domain_roles.GetUserTableRolesesRow {
	return user_domain_roles.GetUserTableRolesesRow{
		ID:                    100,
		UserID:                101,
		CreatedAt:             userTableRolesCreatedAt,
		UpdatedAt:             userTableRolesUpdatedAt,
		TableID:               104,
		DcUserID:              105,
		DcUserName:            "dc-user-name-0",
		DcUserExternalID:      uuid.MustParse("00000000-0000-4000-8000-000000000008"),
		TableRolesID2:         108,
		TableRolesName:        "table-roles-name-0",
		TableRolesDescription: "table-roles-description-0",
		TableCatID:            111,
		TableCatName:          "table-cat-name-0",
	}
}

func TestGetUserTableRolesesRowToProto(t *testing.T) {
	row := testGetUserTableRolesesRow()
	got := GetUserTableRolesesRowToProto(row)

	if got == nil {
		t.Fatal("GetUserTableRolesesRowToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

	if !got.GetCreatedAt().AsTime().Equal(row.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.GetCreatedAt().AsTime(), row.CreatedAt)
	}

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt)
	}

	if got.GetTableId() != row.TableID {
		t.Errorf("TableId = %d, want %d", got.GetTableId(), row.TableID)
	}

	if got.GetDcUserId() != row.DcUserID {
		t.Errorf("DcUserId = %d, want %d", got.GetDcUserId(), row.DcUserID)
	}

	if got.GetDcUserName() != row.DcUserName {
		t.Errorf("DcUserName = %q, want %q", got.GetDcUserName(), row.DcUserName)
	}

	if got.GetDcUserExternalId() != row.DcUserExternalID.String() {
		t.Errorf("DcUserExternalId = %q, want %q", got.GetDcUserExternalId(), row.DcUserExternalID.String())
	}

	if got.GetTableRolesId_2() != row.TableRolesID2 {
		t.Errorf("TableRolesId_2 = %d, want %d", got.GetTableRolesId_2(), row.TableRolesID2)
	}

	if got.GetTableRolesName() != row.TableRolesName {
		t.Errorf("TableRolesName = %q, want %q", got.GetTableRolesName(), row.TableRolesName)
	}

	if got.GetTableRolesDescription() != row.TableRolesDescription {
		t.Errorf("TableRolesDescription = %q, want %q", got.GetTableRolesDescription(), row.TableRolesDescription)
	}

	if got.GetTableCatId() != row.TableCatID {
		t.Errorf("TableCatId = %d, want %d", got.GetTableCatId(), row.TableCatID)
	}

	if got.GetTableCatName() != row.TableCatName {
		t.Errorf("TableCatName = %q, want %q", got.GetTableCatName(), row.TableCatName)
	}

}

func TestGetUserTableRolesesRowsToProto(t *testing.T) {
	first := testGetUserTableRolesesRow()

	second := testGetUserTableRolesesRow()
	second.ID = 999
	second.DcUserName = "second-value"

	tests := []struct {
		name    string
		input   []user_domain_roles.GetUserTableRolesesRow
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.GetUserTableRolesesRow{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.GetUserTableRolesesRow{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUserTableRolesesRowsToProto(tt.input)

			if got == nil {
				t.Fatal("GetUserTableRolesesRowsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// testGetUserTableRolesesFilteredRow — строка выборки GetUserTableRolesesFilteredRow со значениями, различимыми между полями.
func testGetUserTableRolesesFilteredRow() user_domain_roles.GetUserTableRolesesFilteredRow {
	return user_domain_roles.GetUserTableRolesesFilteredRow{
		ID:                    100,
		UserID:                101,
		CreatedAt:             userTableRolesCreatedAt,
		UpdatedAt:             userTableRolesUpdatedAt,
		TableID:               104,
		DcUserID:              105,
		DcUserName:            "dc-user-name-0",
		DcUserExternalID:      uuid.MustParse("00000000-0000-4000-8000-000000000008"),
		TableRolesName:        "table-roles-name-0",
		TableRolesDescription: "table-roles-description-0",
		TableCatID:            110,
		TableCatName:          "table-cat-name-0",
	}
}

func TestGetUserTableRolesesFilteredRowToProto(t *testing.T) {
	row := testGetUserTableRolesesFilteredRow()
	got := GetUserTableRolesesFilteredRowToProto(row)

	if got == nil {
		t.Fatal("GetUserTableRolesesFilteredRowToProto() = nil, want value")
	}

	if got.GetId() != row.ID {
		t.Errorf("Id = %d, want %d", got.GetId(), row.ID)
	}

	if got.GetUserId() != row.UserID {
		t.Errorf("UserId = %d, want %d", got.GetUserId(), row.UserID)
	}

	if !got.GetCreatedAt().AsTime().Equal(row.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", got.GetCreatedAt().AsTime(), row.CreatedAt)
	}

	if !got.GetUpdatedAt().AsTime().Equal(row.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", got.GetUpdatedAt().AsTime(), row.UpdatedAt)
	}

	if got.GetTableId() != row.TableID {
		t.Errorf("TableId = %d, want %d", got.GetTableId(), row.TableID)
	}

	if got.GetDcUserId() != row.DcUserID {
		t.Errorf("DcUserId = %d, want %d", got.GetDcUserId(), row.DcUserID)
	}

	if got.GetDcUserName() != row.DcUserName {
		t.Errorf("DcUserName = %q, want %q", got.GetDcUserName(), row.DcUserName)
	}

	if got.GetDcUserExternalId() != row.DcUserExternalID.String() {
		t.Errorf("DcUserExternalId = %q, want %q", got.GetDcUserExternalId(), row.DcUserExternalID.String())
	}

	if got.GetTableRolesName() != row.TableRolesName {
		t.Errorf("TableRolesName = %q, want %q", got.GetTableRolesName(), row.TableRolesName)
	}

	if got.GetTableRolesDescription() != row.TableRolesDescription {
		t.Errorf("TableRolesDescription = %q, want %q", got.GetTableRolesDescription(), row.TableRolesDescription)
	}

	if got.GetTableCatId() != row.TableCatID {
		t.Errorf("TableCatId = %d, want %d", got.GetTableCatId(), row.TableCatID)
	}

	if got.GetTableCatName() != row.TableCatName {
		t.Errorf("TableCatName = %q, want %q", got.GetTableCatName(), row.TableCatName)
	}

}

func TestGetUserTableRolesesFilteredRowsToProto(t *testing.T) {
	first := testGetUserTableRolesesFilteredRow()

	second := testGetUserTableRolesesFilteredRow()
	second.ID = 999
	second.DcUserName = "second-value"

	tests := []struct {
		name    string
		input   []user_domain_roles.GetUserTableRolesesFilteredRow
		wantLen int
	}{
		{name: "two rows", input: []user_domain_roles.GetUserTableRolesesFilteredRow{first, second}, wantLen: 2},
		{name: "empty slice", input: []user_domain_roles.GetUserTableRolesesFilteredRow{}, wantLen: 0},
		{name: "nil slice", input: nil, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUserTableRolesesFilteredRowsToProto(tt.input)

			if got == nil {
				t.Fatal("GetUserTableRolesesFilteredRowsToProto() = nil, want empty slice")
			}

			if len(got) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestToCreateUserTableRolesParams(t *testing.T) {
	req := &userdomainrolesv1.CreateUserTableRolesRequest{
		UserId:       100,
		TableRolesId: 101,
		TableId:      102,
		ExternalId:   "00000000-0000-4000-8000-000000000004",
	}

	want := user_domain_roles.CreateUserTableRolesParams{
		UserID:       100,
		TableRolesID: 101,
		TableID:      102,
		ExternalID:   uuid.MustParse("00000000-0000-4000-8000-000000000004"),
	}

	if got := ToCreateUserTableRolesParams(req); got != want {
		t.Errorf("ToCreateUserTableRolesParams() = %+v, want %+v", got, want)
	}
}

func TestToCreateUserTableRolesParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToCreateUserTableRolesParams(nil); got != (user_domain_roles.CreateUserTableRolesParams{}) {
		t.Errorf("ToCreateUserTableRolesParams(nil) = %+v, want zero value", got)
	}
}

func TestToUpdateUserTableRolesByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UpdateUserTableRolesByIdRequest{
		UserId:       100,
		TableRolesId: 101,
		TableId:      102,
		ExternalId:   "00000000-0000-4000-8000-000000000004",
		Id:           104,
	}

	want := user_domain_roles.UpdateUserTableRolesByIdParams{
		UserID:       100,
		TableRolesID: 101,
		TableID:      102,
		ExternalID:   uuid.MustParse("00000000-0000-4000-8000-000000000004"),
		ID:           104,
	}

	if got := ToUpdateUserTableRolesByIdParams(req); got != want {
		t.Errorf("ToUpdateUserTableRolesByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUpdateUserTableRolesByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUpdateUserTableRolesByIdParams(nil); got != (user_domain_roles.UpdateUserTableRolesByIdParams{}) {
		t.Errorf("ToUpdateUserTableRolesByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToDeleteUserTableRolesByIdParams(t *testing.T) {
	req := &userdomainrolesv1.DeleteUserTableRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := user_domain_roles.DeleteUserTableRolesByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToDeleteUserTableRolesByIdParams(req); got != want {
		t.Errorf("ToDeleteUserTableRolesByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToDeleteUserTableRolesByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToDeleteUserTableRolesByIdParams(nil); got != (user_domain_roles.DeleteUserTableRolesByIdParams{}) {
		t.Errorf("ToDeleteUserTableRolesByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestToUndeleteUserTableRolesByIdParams(t *testing.T) {
	req := &userdomainrolesv1.UndeleteUserTableRolesByIdRequest{
		ExternalId: "00000000-0000-4000-8000-000000000001",
		Id:         101,
	}

	want := user_domain_roles.UndeleteUserTableRolesByIdParams{
		ExternalID: uuid.MustParse("00000000-0000-4000-8000-000000000001"),
		ID:         101,
	}

	if got := ToUndeleteUserTableRolesByIdParams(req); got != want {
		t.Errorf("ToUndeleteUserTableRolesByIdParams() = %+v, want %+v", got, want)
	}
}

func TestToUndeleteUserTableRolesByIdParamsNil(t *testing.T) {
	// Геттеры protobuf безопасны на nil: сервер не должен падать.
	if got := ToUndeleteUserTableRolesByIdParams(nil); got != (user_domain_roles.UndeleteUserTableRolesByIdParams{}) {
		t.Errorf("ToUndeleteUserTableRolesByIdParams(nil) = %+v, want zero value", got)
	}
}

func TestGetUserTableRolesesDefaultsPageLimit(t *testing.T) {
	got := ToGetUserTableRolesesParams(&userdomainrolesv1.GetUserTableRolesesRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetUserTableRolesesDefaultsPage(t *testing.T) {
	got := ToGetUserTableRolesesParams(&userdomainrolesv1.GetUserTableRolesesRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetUserTableRolesesKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetUserTableRolesesParams(&userdomainrolesv1.GetUserTableRolesesRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetUserTableRolesesFilteredDefaultsPageLimit(t *testing.T) {
	got := ToGetUserTableRolesesFilteredParams(&userdomainrolesv1.GetUserTableRolesesFilteredRequest{Page: 3})

	if got.PageLimit != validation.DefaultPageSize {
		t.Errorf("PageLimit = %d, want %d", got.PageLimit, validation.DefaultPageSize)
	}

	if got.Page != 3 {
		t.Errorf("Page = %d, want 3", got.Page)
	}
}

func TestGetUserTableRolesesFilteredDefaultsPage(t *testing.T) {
	got := ToGetUserTableRolesesFilteredParams(&userdomainrolesv1.GetUserTableRolesesFilteredRequest{PageLimit: 10})

	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetUserTableRolesesFilteredKeepsExplicitPageLimit(t *testing.T) {
	got := ToGetUserTableRolesesFilteredParams(&userdomainrolesv1.GetUserTableRolesesFilteredRequest{PageLimit: 10, Page: 5})

	if got.PageLimit != 10 {
		t.Errorf("PageLimit = %d, want 10", got.PageLimit)
	}
}

func TestGetUserTableRolesesFilteredPassesFilterFields(t *testing.T) {
	req := &userdomainrolesv1.GetUserTableRolesesFilteredRequest{PageLimit: 10, Page: 1,
		DcUserName:       userTableRolesPtr("dcusername-0"),
		DcUserExternalId: userTableRolesPtr("00000000-0000-4000-8000-000000000002"),
		TableRolesName:   userTableRolesPtr("tablerolesname-2"),
		TableCatId:       userTableRolesPtr(int64(103)),
		TableCatName:     userTableRolesPtr("tablecatname-4"),
	}

	got := ToGetUserTableRolesesFilteredParams(req)

	if got.DcUserName != (sql.NullString{String: "dcusername-0", Valid: true}) {
		t.Errorf("DcUserName = %v, want %v", got.DcUserName, (sql.NullString{String: "dcusername-0", Valid: true}))
	}
	if got.DcUserExternalID != (uuid.NullUUID{UUID: uuid.MustParse("00000000-0000-4000-8000-000000000002"), Valid: true}) {
		t.Errorf("DcUserExternalID = %v, want %v", got.DcUserExternalID, (uuid.NullUUID{UUID: uuid.MustParse("00000000-0000-4000-8000-000000000002"), Valid: true}))
	}
	if got.TableRolesName != (sql.NullString{String: "tablerolesname-2", Valid: true}) {
		t.Errorf("TableRolesName = %v, want %v", got.TableRolesName, (sql.NullString{String: "tablerolesname-2", Valid: true}))
	}
	if got.TableCatID != (sql.NullInt64{Int64: 103, Valid: true}) {
		t.Errorf("TableCatID = %v, want %v", got.TableCatID, (sql.NullInt64{Int64: 103, Valid: true}))
	}
	if got.TableCatName != (sql.NullString{String: "tablecatname-4", Valid: true}) {
		t.Errorf("TableCatName = %v, want %v", got.TableCatName, (sql.NullString{String: "tablecatname-4", Valid: true}))
	}
}

// Непереданный необязательный фильтр обязан остаться SQL NULL, а не
// превратиться в 0/false: иначе выборка молча сузилась бы.
func TestGetUserTableRolesesFilteredOmitsAbsentFilters(t *testing.T) {
	got := ToGetUserTableRolesesFilteredParams(&userdomainrolesv1.GetUserTableRolesesFilteredRequest{PageLimit: 10, Page: 1})

	if got.DcUserName.Valid {
		t.Errorf("DcUserName = %+v, want invalid", got.DcUserName)
	}
	if got.DcUserExternalID.Valid {
		t.Errorf("DcUserExternalID = %+v, want invalid", got.DcUserExternalID)
	}
	if got.TableRolesName.Valid {
		t.Errorf("TableRolesName = %+v, want invalid", got.TableRolesName)
	}
	if got.TableCatID.Valid {
		t.Errorf("TableCatID = %+v, want invalid", got.TableCatID)
	}
	if got.TableCatName.Valid {
		t.Errorf("TableCatName = %+v, want invalid", got.TableCatName)
	}
}
