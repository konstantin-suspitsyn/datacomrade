package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateUserTableRoles проверяет запрос на вставку строки dc.user_table_roles.
func ValidateCreateUserTableRoles(req *userdomainrolesv1.CreateUserTableRolesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("table_roles_id", req.GetTableRolesId())
	v.Int64ID("table_id", req.GetTableId())
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateUserTableRolesById проверяет запрос на обновление строки dc.user_table_roles.
func ValidateUpdateUserTableRolesById(req *userdomainrolesv1.UpdateUserTableRolesByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("table_roles_id", req.GetTableRolesId())
	v.Int64ID("table_id", req.GetTableId())
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteUserTableRolesById проверяет запрос на мягкое удаление строки dc.user_table_roles.
func ValidateDeleteUserTableRolesById(req *userdomainrolesv1.DeleteUserTableRolesByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteUserTableRolesById проверяет запрос на обратное удаление строки dc.user_table_roles.
func ValidateUndeleteUserTableRolesById(req *userdomainrolesv1.UndeleteUserTableRolesByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetUserTableRoleses проверяет запрос страницы dc.user_table_roles.
func ValidateGetUserTableRoleses(req *userdomainrolesv1.GetUserTableRolesesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)
	v.Int32Min("page", req.GetPage(), 0)
	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}

// ValidateGetUserTableRolesesFiltered проверяет запрос страницы dc.user_table_roles.
func ValidateGetUserTableRolesesFiltered(req *userdomainrolesv1.GetUserTableRolesesFilteredRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)
	v.Int32Min("page", req.GetPage(), 0)
	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}
