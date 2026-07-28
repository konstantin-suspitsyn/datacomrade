package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// domainsDomainRoleWritableFields проверяет поля, общие для вставки и обновления dc.domains_domain_roles.
func domainsDomainRoleWritableFields(
	v *validator.Validator,
	domainCatId int64,
	domainRolesId int64,
) {
	v.Int64ID("domain_cat_id", domainCatId)
	v.Int64ID("domain_roles_id", domainRolesId)
}

// ValidateCreateDomainsDomainRole проверяет запрос на вставку строки dc.domains_domain_roles.
func ValidateCreateDomainsDomainRole(req *userdomainrolesv1.CreateDomainsDomainRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	domainsDomainRoleWritableFields(
		v,
		req.GetDomainCatId(),
		req.GetDomainRolesId(),
	)

	return v.Err()
}

// ValidateUpdateDomainsDomainRoleById проверяет запрос на обновление строки dc.domains_domain_roles.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateDomainsDomainRoleById(req *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	domainsDomainRoleWritableFields(
		v,
		req.GetDomainCatId(),
		req.GetDomainRolesId(),
	)

	return v.Err()
}
