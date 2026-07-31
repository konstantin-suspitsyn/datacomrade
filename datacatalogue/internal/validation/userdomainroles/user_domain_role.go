package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userDomainRoleWritableFields проверяет поля, общие для вставки и обновления dc.user_domain_roles.
func userDomainRoleWritableFields(
	v *validator.Validator,
	userId int64,
	domainRolesId int64,
	domainId int64,
	updatedByExternalId string,
) {
	v.Int64ID("user_id", userId)
	v.Int64ID("domain_roles_id", domainRolesId)
	v.Int64ID("domain_id", domainId)
	v.StringUUID("updated_by_id", updatedByExternalId)
}

// ValidateCreateUserDomainRole проверяет запрос на вставку строки dc.user_domain_roles.
func ValidateCreateUserDomainRole(req *userdomainrolesv1.CreateUserDomainRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	userDomainRoleWritableFields(
		v,
		req.GetUserId(),
		req.GetDomainRolesId(),
		req.GetDomainId(),
		req.GetUpdatedByExternalId(),
	)

	return v.Err()
}

// ValidateUpdateUserDomainRoleById проверяет запрос на обновление строки dc.user_domain_roles.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateUserDomainRoleById(req *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	userDomainRoleWritableFields(
		v,
		req.GetUserId(),
		req.GetDomainRolesId(),
		req.GetDomainId(),
		req.GetUpdatedByExternalId(),
	)

	return v.Err()
}
