package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateUserDomainRoles проверяет запрос на вставку строки dc.user_domain_roles.
func ValidateCreateUserDomainRoles(req *userdomainrolesv1.CreateUserDomainRolesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("domain_roles_id", req.GetDomainRolesId())
	v.Int64ID("domain_id", req.GetDomainId())
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateUserDomainRolesById проверяет запрос на обновление строки dc.user_domain_roles.
func ValidateUpdateUserDomainRolesById(req *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("domain_roles_id", req.GetDomainRolesId())
	v.Int64ID("domain_id", req.GetDomainId())
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteUserDomainRolesById проверяет запрос на мягкое удаление строки dc.user_domain_roles.
func ValidateDeleteUserDomainRolesById(req *userdomainrolesv1.DeleteUserDomainRolesByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteUserDomainRolesById проверяет запрос на обратное удаление строки dc.user_domain_roles.
func ValidateUndeleteUserDomainRolesById(req *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetUserDomainRoleses проверяет запрос страницы dc.user_domain_roles.
func ValidateGetUserDomainRoleses(req *userdomainrolesv1.GetUserDomainRolesesRequest) error {
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
