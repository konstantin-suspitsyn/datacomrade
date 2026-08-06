package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateTableRole проверяет запрос на вставку строки dc.table_roles.
func ValidateCreateTableRole(req *userdomainrolesv1.CreateTableRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), tableRoleNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), tableRoleDescriptionMaxLen)

	return v.Err()
}

// ValidateUpdateTableRoleById проверяет запрос на обновление строки dc.table_roles.
func ValidateUpdateTableRoleById(req *userdomainrolesv1.UpdateTableRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())
	v.StringVarchar("name", req.GetName(), tableRoleNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), tableRoleDescriptionMaxLen)

	return v.Err()
}
