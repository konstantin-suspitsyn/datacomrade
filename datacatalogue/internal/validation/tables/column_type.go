package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateColumnType проверяет запрос на вставку строки dc.column_type.
func ValidateCreateColumnType(req *tablesv1.CreateColumnTypeRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), columnTypeNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), columnTypeDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateColumnTypeById проверяет запрос на обновление строки dc.column_type.
func ValidateUpdateColumnTypeById(req *tablesv1.UpdateColumnTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), columnTypeNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), columnTypeDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteColumnTypeById проверяет запрос на мягкое удаление строки dc.column_type.
func ValidateDeleteColumnTypeById(req *tablesv1.DeleteColumnTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteColumnTypeById проверяет запрос на обратное удаление строки dc.column_type.
func ValidateUndeleteColumnTypeById(req *tablesv1.UndeleteColumnTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetColumnTypes проверяет запрос страницы dc.column_type.
func ValidateGetColumnTypes(req *tablesv1.GetColumnTypesRequest) error {
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

// ValidateGetColumnTypesSearchName проверяет запрос страницы dc.column_type.
func ValidateGetColumnTypesSearchName(req *tablesv1.GetColumnTypesSearchNameRequest) error {
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

// ValidateGetColumnTypesDeleted проверяет запрос страницы dc.column_type.
func ValidateGetColumnTypesDeleted(req *tablesv1.GetColumnTypesDeletedRequest) error {
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
