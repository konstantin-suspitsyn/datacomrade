package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// UserTableRoleToProto переводит строку dc.user_table_roles в сущность gRPC.
func UserTableRoleToProto(row user_domain_roles.DcUserTableRole) *userdomainrolesv1.UserTableRole {
	return &userdomainrolesv1.UserTableRole{
		Id:           row.ID,
		UserId:       row.UserID,
		TableRolesId: row.TableRolesID,
		CreatedAt:    converter.TimeToProto(row.CreatedAt),
		UpdatedAt:    converter.TimeToProto(row.UpdatedAt),
		IsDeleted:    row.IsDeleted,
		TableId:      row.TableID,
		UpdatedById:  row.UpdatedByID,
	}
}

// UserTableRolesToProto переводит список строк dc.user_table_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func UserTableRolesToProto(rows []user_domain_roles.DcUserTableRole) []*userdomainrolesv1.UserTableRole {
	items := make([]*userdomainrolesv1.UserTableRole, 0, len(rows))

	for _, row := range rows {
		items = append(items, UserTableRoleToProto(row))
	}

	return items
}

// ToCreateUserTableRoleParams собирает параметры вставки dc.user_table_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateUserTableRoleParams(req *userdomainrolesv1.CreateUserTableRoleRequest) user_domain_roles.CreateUserTableRoleParams {
	return user_domain_roles.CreateUserTableRoleParams{
		UserID:       req.GetUserId(),
		TableRolesID: req.GetTableRolesId(),
		TableID:      req.GetTableId(),
		ExternalID:   converter.ProtoToUUID(req.GetUpdatedByExternalId()),
	}
}

// ToUpdateUserTableRoleByIdParams собирает параметры обновления dc.user_table_roles из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateUserTableRoleByIdParams(req *userdomainrolesv1.UpdateUserTableRoleByIdRequest) user_domain_roles.UpdateUserTableRoleByIdParams {
	return user_domain_roles.UpdateUserTableRoleByIdParams{
		ID:           req.GetId(),
		UserID:       req.GetUserId(),
		TableRolesID: req.GetTableRolesId(),
		TableID:      req.GetTableId(),
		ExternalID:   converter.ProtoToUUID(req.GetUpdatedByExternalId()),
	}
}
