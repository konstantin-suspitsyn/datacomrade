package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateColumnCat проверяет запрос на вставку строки dc.column_cat.
func ValidateCreateColumnCat(req *tablesv1.CreateColumnCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("table_id", req.GetTableId())
	v.StringVarchar("name", req.GetName(), columnCatNameMaxLen)
	v.OptionalInt64ID("alias_id", req.AliasId)
	v.Int64ID("column_type_id", req.GetColumnTypeId())
	v.StringVarchar("description", req.GetDescription(), columnCatDescriptionMaxLen)
	v.Int64ID("calculation_type_id", req.GetCalculationTypeId())
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateColumnCatById проверяет запрос на обновление строки dc.column_cat.
func ValidateUpdateColumnCatById(req *tablesv1.UpdateColumnCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("table_id", req.GetTableId())
	v.StringVarchar("name", req.GetName(), columnCatNameMaxLen)
	v.OptionalInt64ID("alias_id", req.AliasId)
	v.Int64ID("column_type_id", req.GetColumnTypeId())
	v.StringVarchar("description", req.GetDescription(), columnCatDescriptionMaxLen)
	v.Int64ID("calculation_type_id", req.GetCalculationTypeId())
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteColumnCatById проверяет запрос на мягкое удаление строки dc.column_cat.
func ValidateDeleteColumnCatById(req *tablesv1.DeleteColumnCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteColumnCatById проверяет запрос на обратное удаление строки dc.column_cat.
func ValidateUndeleteColumnCatById(req *tablesv1.UndeleteColumnCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetColumnCats проверяет запрос страницы dc.column_cat.
func ValidateGetColumnCats(req *tablesv1.GetColumnCatsRequest) error {
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

// ValidateGetColumnCatsDeleted проверяет запрос страницы dc.column_cat.
func ValidateGetColumnCatsDeleted(req *tablesv1.GetColumnCatsDeletedRequest) error {
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

// ValidateGetColumnCatsOptionalFilters проверяет запрос страницы dc.column_cat.
func ValidateGetColumnCatsOptionalFilters(req *tablesv1.GetColumnCatsOptionalFiltersRequest) error {
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

// ValidateGetColumnCatsSearchName проверяет запрос страницы dc.column_cat.
func ValidateGetColumnCatsSearchName(req *tablesv1.GetColumnCatsSearchNameRequest) error {
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
