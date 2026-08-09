package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateSchemaCat проверяет запрос на вставку строки dc.schema_cat.
func ValidateCreateSchemaCat(req *tablesv1.CreateSchemaCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("database_id", req.GetDatabaseId())
	v.StringVarchar("name", req.GetName(), schemaCatNameMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateSchemaCatById проверяет запрос на обновление строки dc.schema_cat.
func ValidateUpdateSchemaCatById(req *tablesv1.UpdateSchemaCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("database_id", req.GetDatabaseId())
	v.StringVarchar("name", req.GetName(), schemaCatNameMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteSchemaCatById проверяет запрос на мягкое удаление строки dc.schema_cat.
func ValidateDeleteSchemaCatById(req *tablesv1.DeleteSchemaCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteSchemaCatById проверяет запрос на обратное удаление строки dc.schema_cat.
func ValidateUndeleteSchemaCatById(req *tablesv1.UndeleteSchemaCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetSchemaCats проверяет запрос страницы dc.schema_cat.
func ValidateGetSchemaCats(req *tablesv1.GetSchemaCatsRequest) error {
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

// ValidateGetSchemaCatsDeleted проверяет запрос страницы dc.schema_cat.
func ValidateGetSchemaCatsDeleted(req *tablesv1.GetSchemaCatsDeletedRequest) error {
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
