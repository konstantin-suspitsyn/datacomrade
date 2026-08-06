package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateDomainRole проверяет запрос на вставку строки dc.domain_roles.
func ValidateCreateDomainRole(req *userdomainrolesv1.CreateDomainRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), domainRoleNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), domainRoleDescriptionMaxLen)

	return v.Err()
}

// ValidateUpdateDomainRoleById проверяет запрос на обновление строки dc.domain_roles.
func ValidateUpdateDomainRoleById(req *userdomainrolesv1.UpdateDomainRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())
	v.StringVarchar("name", req.GetName(), domainRoleNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), domainRoleDescriptionMaxLen)

	return v.Err()
}
