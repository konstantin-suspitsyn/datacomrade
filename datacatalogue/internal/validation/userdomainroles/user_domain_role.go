package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateUserDomainRole проверяет запрос на вставку строки dc.user_domain_roles.
func ValidateCreateUserDomainRole(req *userdomainrolesv1.CreateUserDomainRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("domain_roles_id", req.GetDomainRolesId())
	v.Int64ID("domain_id", req.GetDomainId())
	v.StringUUID("updated_by_id", req.GetUpdatedByExternalId())

	return v.Err()
}

// ValidateUpdateUserDomainRoleById проверяет запрос на обновление строки dc.user_domain_roles.
func ValidateUpdateUserDomainRoleById(req *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())
	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("domain_roles_id", req.GetDomainRolesId())
	v.Int64ID("domain_id", req.GetDomainId())
	v.StringUUID("updated_by_id", req.GetUpdatedByExternalId())

	return v.Err()
}
