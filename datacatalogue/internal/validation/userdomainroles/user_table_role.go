package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateUserTableRole проверяет запрос на вставку строки dc.user_table_roles.
func ValidateCreateUserTableRole(req *userdomainrolesv1.CreateUserTableRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("table_roles_id", req.GetTableRolesId())
	v.Int64ID("table_id", req.GetTableId())
	v.StringUUID("updated_by_id", req.GetUpdatedByExternalId())

	return v.Err()
}

// ValidateUpdateUserTableRoleById проверяет запрос на обновление строки dc.user_table_roles.
func ValidateUpdateUserTableRoleById(req *userdomainrolesv1.UpdateUserTableRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())
	v.Int64ID("user_id", req.GetUserId())
	v.Int64ID("table_roles_id", req.GetTableRolesId())
	v.Int64ID("table_id", req.GetTableId())
	v.StringUUID("updated_by_id", req.GetUpdatedByExternalId())

	return v.Err()
}
