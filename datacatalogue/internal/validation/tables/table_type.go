package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateTableType проверяет запрос на вставку строки dc.table_type.
func ValidateCreateTableType(req *tablesv1.CreateTableTypeRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), tableTypeNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), tableTypeDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateTableTypeById проверяет запрос на обновление строки dc.table_type.
func ValidateUpdateTableTypeById(req *tablesv1.UpdateTableTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), tableTypeNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), tableTypeDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteTableTypeById проверяет запрос на мягкое удаление строки dc.table_type.
func ValidateDeleteTableTypeById(req *tablesv1.DeleteTableTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteTableTypeById проверяет запрос на обратное удаление строки dc.table_type.
func ValidateUndeleteTableTypeById(req *tablesv1.UndeleteTableTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetTableTypes проверяет запрос страницы dc.table_type.
func ValidateGetTableTypes(req *tablesv1.GetTableTypesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)
	v.Int32Min("page", req.GetPage(), 0)
	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}

// ValidateGetTableTypesDeleted проверяет запрос страницы dc.table_type.
func ValidateGetTableTypesDeleted(req *tablesv1.GetTableTypesDeletedRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)
	v.Int32Min("page", req.GetPage(), 0)
	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}
