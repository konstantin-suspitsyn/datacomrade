package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// columnCatWritableFields проверяет поля, общие для вставки и обновления dc.column_cat.
func columnCatWritableFields(
	v *validator.Validator,
	tableId int64,
	name string,
	aliasId int64,
	columnTypeId int64,
	description string,
	calculationTypeId int64,
	showInUi bool,
	userId int64,
) {
	v.Int64ID("table_id", tableId)
	v.StringVarchar("name", name, columnCatNameMaxLen)
	v.Int64ID("alias_id", aliasId)
	v.Int64ID("column_type_id", columnTypeId)
	v.StringVarchar("description", description, columnCatDescriptionMaxLen)
	v.Int64ID("calculation_type_id", calculationTypeId)
	v.Int64ID("user_id", userId)
}

// ValidateCreateColumnCat проверяет запрос на вставку строки dc.column_cat.
func ValidateCreateColumnCat(req *tablesv1.CreateColumnCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	columnCatWritableFields(
		v,
		req.GetTableId(),
		req.GetName(),
		req.GetAliasId(),
		req.GetColumnTypeId(),
		req.GetDescription(),
		req.GetCalculationTypeId(),
		req.GetShowInUi(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateColumnCatById проверяет запрос на обновление строки dc.column_cat.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateColumnCatById(req *tablesv1.UpdateColumnCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	columnCatWritableFields(
		v,
		req.GetTableId(),
		req.GetName(),
		req.GetAliasId(),
		req.GetColumnTypeId(),
		req.GetDescription(),
		req.GetCalculationTypeId(),
		req.GetShowInUi(),
		req.GetUserId(),
	)

	return v.Err()
}
