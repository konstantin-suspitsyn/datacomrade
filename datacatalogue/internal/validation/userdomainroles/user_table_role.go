package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userTableRoleWritableFields проверяет поля, общие для вставки и обновления dc.user_table_roles.
func userTableRoleWritableFields(
	v *validator.Validator,
	userId int64,
	tableRolesId int64,
	tableId int64,
) {
	v.Int64ID("user_id", userId)
	v.Int64ID("table_roles_id", tableRolesId)
	v.Int64ID("table_id", tableId)
}

// ValidateCreateUserTableRole проверяет запрос на вставку строки dc.user_table_roles.
func ValidateCreateUserTableRole(req *userdomainrolesv1.CreateUserTableRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	userTableRoleWritableFields(
		v,
		req.GetUserId(),
		req.GetTableRolesId(),
		req.GetTableId(),
	)

	return v.Err()
}

// ValidateUpdateUserTableRoleById проверяет запрос на обновление строки dc.user_table_roles.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateUserTableRoleById(req *userdomainrolesv1.UpdateUserTableRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	userTableRoleWritableFields(
		v,
		req.GetUserId(),
		req.GetTableRolesId(),
		req.GetTableId(),
	)

	return v.Err()
}
