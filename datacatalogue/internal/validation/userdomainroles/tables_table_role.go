package userdomainroles

import (
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tablesTableRoleWritableFields проверяет поля, общие для вставки и обновления dc.tables_table_roles.
func tablesTableRoleWritableFields(
	v *validator.Validator,
	tableCatId int64,
	tableRolesId int64,
) {
	v.Int64ID("table_cat_id", tableCatId)
	v.Int64ID("table_roles_id", tableRolesId)
}

// ValidateCreateTablesTableRole проверяет запрос на вставку строки dc.tables_table_roles.
func ValidateCreateTablesTableRole(req *userdomainrolesv1.CreateTablesTableRoleRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	tablesTableRoleWritableFields(
		v,
		req.GetTableCatId(),
		req.GetTableRolesId(),
	)

	return v.Err()
}

// ValidateUpdateTablesTableRoleById проверяет запрос на обновление строки dc.tables_table_roles.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateTablesTableRoleById(req *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	tablesTableRoleWritableFields(
		v,
		req.GetTableCatId(),
		req.GetTableRolesId(),
	)

	return v.Err()
}
