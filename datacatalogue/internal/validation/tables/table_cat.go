package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// tableCatWritableFields проверяет поля, общие для вставки и обновления dc.table_cat.
func tableCatWritableFields(
	v *validator.Validator,
	name string,
	description string,
	schemaId int64,
	tableTypeId int64,
	domainId int64,
	isGetDict bool,
	userId int64,
) {
	v.StringVarchar("name", name, tableCatNameMaxLen)
	v.StringVarchar("description", description, tableCatDescriptionMaxLen)
	v.Int64ID("schema_id", schemaId)
	v.Int64ID("table_type_id", tableTypeId)
	v.Int64ID("domain_id", domainId)
	v.Int64ID("user_id", userId)
}

// ValidateCreateTableCat проверяет запрос на вставку строки dc.table_cat.
func ValidateCreateTableCat(req *tablesv1.CreateTableCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	tableCatWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetSchemaId(),
		req.GetTableTypeId(),
		req.GetDomainId(),
		req.GetIsGetDict(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateTableCatById проверяет запрос на обновление строки dc.table_cat.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateTableCatById(req *tablesv1.UpdateTableCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	tableCatWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetSchemaId(),
		req.GetTableTypeId(),
		req.GetDomainId(),
		req.GetIsGetDict(),
		req.GetUserId(),
	)

	return v.Err()
}
