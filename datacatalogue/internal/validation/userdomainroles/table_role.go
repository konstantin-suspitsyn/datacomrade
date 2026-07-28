package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tableRoleWritableFields проверяет поля, общие для вставки и обновления dc.table_roles.
func tableRoleWritableFields(
	v *validator.Validator,
	name string,
	description string,
) {
	v.StringVarchar("name", name, tableRoleNameMaxLen)
	v.StringVarchar("description", description, tableRoleDescriptionMaxLen)
}

// ValidateCreateTableRole проверяет запрос на вставку строки dc.table_roles.
func ValidateCreateTableRole(req *userdomainrolesv1.CreateTableRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	tableRoleWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
	)

	return v.Err()
}

// ValidateUpdateTableRoleById проверяет запрос на обновление строки dc.table_roles.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateTableRoleById(req *userdomainrolesv1.UpdateTableRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	tableRoleWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
	)

	return v.Err()
}
