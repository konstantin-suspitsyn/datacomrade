package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateDatabaseCat проверяет запрос на вставку строки dc.database_cat.
func ValidateCreateDatabaseCat(req *tablesv1.CreateDatabaseCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), databaseCatNameMaxLen)
	v.Int64ID("host_id", req.GetHostId())
	v.Int64ID("database_type_id", req.GetDatabaseTypeId())
	v.StringVarchar("description", req.GetDescription(), databaseCatDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateDatabaseCatById проверяет запрос на обновление строки dc.database_cat.
func ValidateUpdateDatabaseCatById(req *tablesv1.UpdateDatabaseCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), databaseCatNameMaxLen)
	v.Int64ID("host_id", req.GetHostId())
	v.Int64ID("database_type_id", req.GetDatabaseTypeId())
	v.StringVarchar("description", req.GetDescription(), databaseCatDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteDatabaseCatById проверяет запрос на мягкое удаление строки dc.database_cat.
func ValidateDeleteDatabaseCatById(req *tablesv1.DeleteDatabaseCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteDatabaseCatById проверяет запрос на обратное удаление строки dc.database_cat.
func ValidateUndeleteDatabaseCatById(req *tablesv1.UndeleteDatabaseCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetDatabaseCats проверяет запрос страницы dc.database_cat.
func ValidateGetDatabaseCats(req *tablesv1.GetDatabaseCatsRequest) error {
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

// ValidateGetDatabaseCatsDeleted проверяет запрос страницы dc.database_cat.
func ValidateGetDatabaseCatsDeleted(req *tablesv1.GetDatabaseCatsDeletedRequest) error {
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

// ValidateGetDatabaseCatsByHostIdDatabaseTypeId проверяет запрос страницы dc.database_cat.
func ValidateGetDatabaseCatsByHostIdDatabaseTypeId(req *tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)
	v.Int32Min("page", req.GetPage(), 0)
	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")
	v.StringIn("order_by", req.GetOrderBy(), "", "name")

	return v.Err()
}

// ValidateGetDatabaseCatsSearchName проверяет запрос страницы dc.database_cat.
func ValidateGetDatabaseCatsSearchName(req *tablesv1.GetDatabaseCatsSearchNameRequest) error {
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
