package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateTableCat проверяет запрос на вставку строки dc.table_cat.
func ValidateCreateTableCat(req *tablesv1.CreateTableCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), tableCatNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), tableCatDescriptionMaxLen)
	v.Int64ID("schema_id", req.GetSchemaId())
	v.Int64ID("table_type_id", req.GetTableTypeId())
	v.Int64ID("domain_id", req.GetDomainId())
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateTableCatById проверяет запрос на обновление строки dc.table_cat.
func ValidateUpdateTableCatById(req *tablesv1.UpdateTableCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), tableCatNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), tableCatDescriptionMaxLen)
	v.Int64ID("schema_id", req.GetSchemaId())
	v.Int64ID("table_type_id", req.GetTableTypeId())
	v.Int64ID("domain_id", req.GetDomainId())
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteTableCatById проверяет запрос на мягкое удаление строки dc.table_cat.
func ValidateDeleteTableCatById(req *tablesv1.DeleteTableCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteTableCatById проверяет запрос на обратное удаление строки dc.table_cat.
func ValidateUndeleteTableCatById(req *tablesv1.UndeleteTableCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetTableCats проверяет запрос страницы dc.table_cat.
func ValidateGetTableCats(req *tablesv1.GetTableCatsRequest) error {
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

// ValidateGetTableCatsDeleted проверяет запрос страницы dc.table_cat.
func ValidateGetTableCatsDeleted(req *tablesv1.GetTableCatsDeletedRequest) error {
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
