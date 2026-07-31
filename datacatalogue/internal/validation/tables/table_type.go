package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tableTypeWritableFields проверяет поля, общие для вставки и обновления dc.table_type.
func tableTypeWritableFields(
	v *validator.Validator,
	name string,
	description string,
	userExternalId string,
) {
	v.StringVarchar("name", name, tableTypeNameMaxLen)
	v.StringVarchar("description", description, tableTypeDescriptionMaxLen)
	v.StringUUID("user_id", userExternalId)
}

// ValidateCreateTableType проверяет запрос на вставку строки dc.table_type.
func ValidateCreateTableType(req *tablesv1.CreateTableTypeRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	tableTypeWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetUserExternalId(),
	)

	return v.Err()
}

// ValidateUpdateTableTypeById проверяет запрос на обновление строки dc.table_type.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateTableTypeById(req *tablesv1.UpdateTableTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	tableTypeWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetUserExternalId(),
	)

	return v.Err()
}
