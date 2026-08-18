package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// UserDomainRolesToProto переводит строку dc.user_domain_roles в сущность gRPC.
func UserDomainRolesToProto(row user_domain_roles.DcUserDomainRole) *userdomainrolesv1.UserDomainRoles {
	return &userdomainrolesv1.UserDomainRoles{
		Id:            row.ID,
		UserId:        row.UserID,
		DomainRolesId: row.DomainRolesID,
		CreatedAt:     converter.TimeToProto(row.CreatedAt),
		UpdatedAt:     converter.TimeToProto(row.UpdatedAt),
		IsDeleted:     row.IsDeleted,
		DomainId:      row.DomainID,
		UpdatedById:   row.UpdatedByID,
	}
}

// UserDomainRolesesToProto переводит список строк dc.user_domain_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func UserDomainRolesesToProto(rows []user_domain_roles.DcUserDomainRole) []*userdomainrolesv1.UserDomainRoles {
	items := make([]*userdomainrolesv1.UserDomainRoles, 0, len(rows))

	for _, row := range rows {
		items = append(items, UserDomainRolesToProto(row))
	}

	return items
}

// GetUserDomainRolesesRowToProto переводит строку выборки в сообщение gRPC.
func GetUserDomainRolesesRowToProto(row user_domain_roles.GetUserDomainRolesesRow) *userdomainrolesv1.GetUserDomainRolesesRow {
	return &userdomainrolesv1.GetUserDomainRolesesRow{
		Id:                     row.ID,
		UserId:                 row.UserID,
		DomainRolesId:          row.DomainRolesID,
		CreatedAt:              converter.TimeToProto(row.CreatedAt),
		UpdatedAt:              converter.TimeToProto(row.UpdatedAt),
		IsDeleted:              row.IsDeleted,
		DomainId:               row.DomainID,
		UpdatedById:            row.UpdatedByID,
		DomainRolesName:        row.DomainRolesName,
		DomainRolesDescription: row.DomainRolesDescription,
		DcUserName:             row.DcUserName,
	}
}

// GetUserDomainRolesesRowsToProto переводит список строк выборки в список сообщений gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func GetUserDomainRolesesRowsToProto(rows []user_domain_roles.GetUserDomainRolesesRow) []*userdomainrolesv1.GetUserDomainRolesesRow {
	items := make([]*userdomainrolesv1.GetUserDomainRolesesRow, 0, len(rows))

	for _, row := range rows {
		items = append(items, GetUserDomainRolesesRowToProto(row))
	}

	return items
}

// ToCreateUserDomainRolesParams собирает параметры вставки dc.user_domain_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateUserDomainRolesParams(req *userdomainrolesv1.CreateUserDomainRolesRequest) user_domain_roles.CreateUserDomainRolesParams {
	return user_domain_roles.CreateUserDomainRolesParams{
		UserID:        req.GetUserId(),
		DomainRolesID: req.GetDomainRolesId(),
		DomainID:      req.GetDomainId(),
		ExternalID:    converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateUserDomainRolesByIdParams собирает параметры обновления dc.user_domain_roles из запроса gRPC.
func ToUpdateUserDomainRolesByIdParams(req *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) user_domain_roles.UpdateUserDomainRolesByIdParams {
	return user_domain_roles.UpdateUserDomainRolesByIdParams{
		UserID:        req.GetUserId(),
		DomainRolesID: req.GetDomainRolesId(),
		DomainID:      req.GetDomainId(),
		ExternalID:    converter.ProtoToUUID(req.GetExternalId()),
		ID:            req.GetId(),
	}
}

// ToDeleteUserDomainRolesByIdParams собирает параметры мягкого удаления dc.user_domain_roles из запроса gRPC.
func ToDeleteUserDomainRolesByIdParams(req *userdomainrolesv1.DeleteUserDomainRolesByIdRequest) user_domain_roles.DeleteUserDomainRolesByIdParams {
	return user_domain_roles.DeleteUserDomainRolesByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteUserDomainRolesByIdParams собирает параметры обратного удаления dc.user_domain_roles из запроса gRPC.
func ToUndeleteUserDomainRolesByIdParams(req *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest) user_domain_roles.UndeleteUserDomainRolesByIdParams {
	return user_domain_roles.UndeleteUserDomainRolesByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetUserDomainRolesesParams собирает параметры страницы dc.user_domain_roles из запроса gRPC.
func ToGetUserDomainRolesesParams(req *userdomainrolesv1.GetUserDomainRolesesRequest) user_domain_roles.GetUserDomainRolesesParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return user_domain_roles.GetUserDomainRolesesParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
