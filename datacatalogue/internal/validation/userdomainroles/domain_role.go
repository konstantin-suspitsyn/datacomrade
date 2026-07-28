package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// domainRoleWritableFields проверяет поля, общие для вставки и обновления dc.domain_roles.
func domainRoleWritableFields(
	v *validator.Validator,
	name string,
	description string,
) {
	v.StringVarchar("name", name, domainRoleNameMaxLen)
	v.StringVarchar("description", description, domainRoleDescriptionMaxLen)
}

// ValidateCreateDomainRole проверяет запрос на вставку строки dc.domain_roles.
func ValidateCreateDomainRole(req *userdomainrolesv1.CreateDomainRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	domainRoleWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
	)

	return v.Err()
}

// ValidateUpdateDomainRoleById проверяет запрос на обновление строки dc.domain_roles.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateDomainRoleById(req *userdomainrolesv1.UpdateDomainRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	domainRoleWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
	)

	return v.Err()
}
