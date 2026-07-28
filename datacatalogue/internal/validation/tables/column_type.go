package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// columnTypeWritableFields проверяет поля, общие для вставки и обновления dc.column_type.
func columnTypeWritableFields(
	v *validator.Validator,
	name string,
	description string,
	userId int64,
) {
	v.StringVarchar("name", name, columnTypeNameMaxLen)
	v.StringVarchar("description", description, columnTypeDescriptionMaxLen)
	v.Int64ID("user_id", userId)
}

// ValidateCreateColumnType проверяет запрос на вставку строки dc.column_type.
func ValidateCreateColumnType(req *tablesv1.CreateColumnTypeRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	columnTypeWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateColumnTypeById проверяет запрос на обновление строки dc.column_type.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateColumnTypeById(req *tablesv1.UpdateColumnTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	columnTypeWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}
