package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateDatabaseType проверяет запрос на вставку строки dc.database_type.
func ValidateCreateDatabaseType(req *tablesv1.CreateDatabaseTypeRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), databaseTypeNameMaxLen)
	v.StringVarchar("db_version", req.GetDbVersion(), databaseTypeDbVersionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateDatabaseTypeById проверяет запрос на обновление строки dc.database_type.
func ValidateUpdateDatabaseTypeById(req *tablesv1.UpdateDatabaseTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), databaseTypeNameMaxLen)
	v.StringVarchar("db_version", req.GetDbVersion(), databaseTypeDbVersionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteDatabaseTypeById проверяет запрос на мягкое удаление строки dc.database_type.
func ValidateDeleteDatabaseTypeById(req *tablesv1.DeleteDatabaseTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteDatabaseTypeById проверяет запрос на обратное удаление строки dc.database_type.
func ValidateUndeleteDatabaseTypeById(req *tablesv1.UndeleteDatabaseTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetDatabaseTypesSearchName проверяет запрос страницы dc.database_type.
func ValidateGetDatabaseTypesSearchName(req *tablesv1.GetDatabaseTypesSearchNameRequest) error {
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

// ValidateGetDatabaseTypes проверяет запрос страницы dc.database_type.
func ValidateGetDatabaseTypes(req *tablesv1.GetDatabaseTypesRequest) error {
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

// ValidateGetDatabaseTypesDeleted проверяет запрос страницы dc.database_type.
func ValidateGetDatabaseTypesDeleted(req *tablesv1.GetDatabaseTypesDeletedRequest) error {
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
