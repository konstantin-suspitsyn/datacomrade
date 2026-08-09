package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateDatabaseCalculation проверяет запрос на вставку строки dc.database_calculation.
func ValidateCreateDatabaseCalculation(req *tablesv1.CreateDatabaseCalculationRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("database_cat_id", req.GetDatabaseCatId())
	v.Int64ID("calculation_type_id", req.GetCalculationTypeId())
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateDatabaseCalculationById проверяет запрос на обновление строки dc.database_calculation.
func ValidateUpdateDatabaseCalculationById(req *tablesv1.UpdateDatabaseCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("database_cat_id", req.GetDatabaseCatId())
	v.Int64ID("calculation_type_id", req.GetCalculationTypeId())
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteDatabaseCalculationById проверяет запрос на мягкое удаление строки dc.database_calculation.
func ValidateDeleteDatabaseCalculationById(req *tablesv1.DeleteDatabaseCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteDatabaseCalculationById проверяет запрос на обратное удаление строки dc.database_calculation.
func ValidateUndeleteDatabaseCalculationById(req *tablesv1.UndeleteDatabaseCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetDatabaseCalculationsSearchName проверяет запрос страницы dc.database_calculation.
func ValidateGetDatabaseCalculationsSearchName(req *tablesv1.GetDatabaseCalculationsSearchNameRequest) error {
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

// ValidateGetDatabaseCalculationsFiltered проверяет запрос страницы dc.database_calculation.
func ValidateGetDatabaseCalculationsFiltered(req *tablesv1.GetDatabaseCalculationsFilteredRequest) error {
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

// ValidateGetDatabaseCalculationsDeleted проверяет запрос страницы dc.database_calculation.
func ValidateGetDatabaseCalculationsDeleted(req *tablesv1.GetDatabaseCalculationsDeletedRequest) error {
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

// ValidateGetDatabaseCalculations проверяет запрос страницы dc.database_calculation.
func ValidateGetDatabaseCalculations(req *tablesv1.GetDatabaseCalculationsRequest) error {
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
