package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// UserTableRolesToProto переводит строку dc.user_table_roles в сущность gRPC.
func UserTableRolesToProto(row user_domain_roles.DcUserTableRole) *userdomainrolesv1.UserTableRoles {
	return &userdomainrolesv1.UserTableRoles{
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

// UserTableRolesesToProto переводит список строк dc.user_table_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func UserTableRolesesToProto(rows []user_domain_roles.DcUserTableRole) []*userdomainrolesv1.UserTableRoles {
	items := make([]*userdomainrolesv1.UserTableRoles, 0, len(rows))

	for _, row := range rows {
		items = append(items, UserTableRolesToProto(row))
	}

	return items
}

// GetUserTableRolesByIdRowToProto переводит строку выборки в сообщение gRPC.
func GetUserTableRolesByIdRowToProto(row user_domain_roles.GetUserTableRolesByIdRow) *userdomainrolesv1.GetUserTableRolesByIdRow {
	return &userdomainrolesv1.GetUserTableRolesByIdRow{
		Id:                    row.ID,
		UserId:                row.UserID,
		TableRolesId:          row.TableRolesID,
		CreatedAt:             converter.TimeToProto(row.CreatedAt),
		UpdatedAt:             converter.TimeToProto(row.UpdatedAt),
		TableId:               row.TableID,
		DcUserId:              converter.Int64ToProto(row.DcUserID),
		DcUserName:            row.DcUserName,
		DcUserExternalId:      converter.UUIDToProto(row.DcUserExternalID),
		TableRolesName:        row.TableRolesName,
		TableRolesDescription: row.TableRolesDescription,
		TableCatId:            converter.Int64ToProto(row.TableCatID),
		TableCatName:          row.TableCatName,
	}
}

// GetUserTableRolesesRowToProto переводит строку выборки в сообщение gRPC.
func GetUserTableRolesesRowToProto(row user_domain_roles.GetUserTableRolesesRow) *userdomainrolesv1.GetUserTableRolesesRow {
	return &userdomainrolesv1.GetUserTableRolesesRow{
		Id:                    row.ID,
		UserId:                row.UserID,
		CreatedAt:             converter.TimeToProto(row.CreatedAt),
		UpdatedAt:             converter.TimeToProto(row.UpdatedAt),
		TableId:               row.TableID,
		DcUserId:              converter.Int64ToProto(row.DcUserID),
		DcUserName:            row.DcUserName,
		DcUserExternalId:      converter.UUIDToProto(row.DcUserExternalID),
		TableRolesId_2:        row.TableRolesID2,
		TableRolesName:        row.TableRolesName,
		TableRolesDescription: row.TableRolesDescription,
		TableCatId:            converter.Int64ToProto(row.TableCatID),
		TableCatName:          row.TableCatName,
	}
}

// GetUserTableRolesesRowsToProto переводит список строк выборки в список сообщений gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func GetUserTableRolesesRowsToProto(rows []user_domain_roles.GetUserTableRolesesRow) []*userdomainrolesv1.GetUserTableRolesesRow {
	items := make([]*userdomainrolesv1.GetUserTableRolesesRow, 0, len(rows))

	for _, row := range rows {
		items = append(items, GetUserTableRolesesRowToProto(row))
	}

	return items
}

// GetUserTableRolesesFilteredRowToProto переводит строку выборки в сообщение gRPC.
func GetUserTableRolesesFilteredRowToProto(row user_domain_roles.GetUserTableRolesesFilteredRow) *userdomainrolesv1.GetUserTableRolesesFilteredRow {
	return &userdomainrolesv1.GetUserTableRolesesFilteredRow{
		Id:                    row.ID,
		UserId:                row.UserID,
		CreatedAt:             converter.TimeToProto(row.CreatedAt),
		UpdatedAt:             converter.TimeToProto(row.UpdatedAt),
		TableId:               row.TableID,
		DcUserId:              converter.Int64ToProto(row.DcUserID),
		DcUserName:            row.DcUserName,
		DcUserExternalId:      converter.UUIDToProto(row.DcUserExternalID),
		TableRolesName:        row.TableRolesName,
		TableRolesDescription: row.TableRolesDescription,
		TableCatId:            converter.Int64ToProto(row.TableCatID),
		TableCatName:          row.TableCatName,
	}
}

// GetUserTableRolesesFilteredRowsToProto переводит список строк выборки в список сообщений gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func GetUserTableRolesesFilteredRowsToProto(rows []user_domain_roles.GetUserTableRolesesFilteredRow) []*userdomainrolesv1.GetUserTableRolesesFilteredRow {
	items := make([]*userdomainrolesv1.GetUserTableRolesesFilteredRow, 0, len(rows))

	for _, row := range rows {
		items = append(items, GetUserTableRolesesFilteredRowToProto(row))
	}

	return items
}

// ToCreateUserTableRolesParams собирает параметры вставки dc.user_table_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateUserTableRolesParams(req *userdomainrolesv1.CreateUserTableRolesRequest) user_domain_roles.CreateUserTableRolesParams {
	return user_domain_roles.CreateUserTableRolesParams{
		UserID:       req.GetUserId(),
		TableRolesID: req.GetTableRolesId(),
		TableID:      req.GetTableId(),
		ExternalID:   converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateUserTableRolesByIdParams собирает параметры обновления dc.user_table_roles из запроса gRPC.
func ToUpdateUserTableRolesByIdParams(req *userdomainrolesv1.UpdateUserTableRolesByIdRequest) user_domain_roles.UpdateUserTableRolesByIdParams {
	return user_domain_roles.UpdateUserTableRolesByIdParams{
		UserID:       req.GetUserId(),
		TableRolesID: req.GetTableRolesId(),
		TableID:      req.GetTableId(),
		ExternalID:   converter.ProtoToUUID(req.GetExternalId()),
		ID:           req.GetId(),
	}
}

// ToDeleteUserTableRolesByIdParams собирает параметры мягкого удаления dc.user_table_roles из запроса gRPC.
func ToDeleteUserTableRolesByIdParams(req *userdomainrolesv1.DeleteUserTableRolesByIdRequest) user_domain_roles.DeleteUserTableRolesByIdParams {
	return user_domain_roles.DeleteUserTableRolesByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteUserTableRolesByIdParams собирает параметры обратного удаления dc.user_table_roles из запроса gRPC.
func ToUndeleteUserTableRolesByIdParams(req *userdomainrolesv1.UndeleteUserTableRolesByIdRequest) user_domain_roles.UndeleteUserTableRolesByIdParams {
	return user_domain_roles.UndeleteUserTableRolesByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetUserTableRolesesParams собирает параметры страницы dc.user_table_roles из запроса gRPC.
func ToGetUserTableRolesesParams(req *userdomainrolesv1.GetUserTableRolesesRequest) user_domain_roles.GetUserTableRolesesParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return user_domain_roles.GetUserTableRolesesParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetUserTableRolesesFilteredParams собирает параметры страницы dc.user_table_roles из запроса gRPC.
func ToGetUserTableRolesesFilteredParams(req *userdomainrolesv1.GetUserTableRolesesFilteredRequest) user_domain_roles.GetUserTableRolesesFilteredParams {
	if req == nil {
		req = &userdomainrolesv1.GetUserTableRolesesFilteredRequest{}
	}

	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return user_domain_roles.GetUserTableRolesesFilteredParams{
		Order:            req.GetOrder(),
		Page:             page,
		PageLimit:        limit,
		DcUserName:       converter.ProtoToNullString(req.DcUserName),
		DcUserExternalID: converter.ProtoToNullUUID(req.DcUserExternalId),
		TableRolesName:   converter.ProtoToNullString(req.TableRolesName),
		TableCatID:       converter.ProtoToNullInt64(req.TableCatId),
		TableCatName:     converter.ProtoToNullString(req.TableCatName),
	}
}
