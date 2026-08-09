package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateFollowingCalculation проверяет запрос на вставку строки dc.following_calculation.
func ValidateCreateFollowingCalculation(req *tablesv1.CreateFollowingCalculationRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("column_cat_id", req.GetColumnCatId())
	v.Int64ID("calculation_type_id", req.GetCalculationTypeId())
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateFollowingCalculationById проверяет запрос на обновление строки dc.following_calculation.
func ValidateUpdateFollowingCalculationById(req *tablesv1.UpdateFollowingCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("column_cat_id", req.GetColumnCatId())
	v.Int64ID("calculation_type_id", req.GetCalculationTypeId())
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteFollowingCalculationById проверяет запрос на мягкое удаление строки dc.following_calculation.
func ValidateDeleteFollowingCalculationById(req *tablesv1.DeleteFollowingCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteFollowingCalculationById проверяет запрос на обратное удаление строки dc.following_calculation.
func ValidateUndeleteFollowingCalculationById(req *tablesv1.UndeleteFollowingCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetFollowingCalculations проверяет запрос страницы dc.following_calculation.
func ValidateGetFollowingCalculations(req *tablesv1.GetFollowingCalculationsRequest) error {
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

// ValidateGetFollowingCalculationsDeleted проверяет запрос страницы dc.following_calculation.
func ValidateGetFollowingCalculationsDeleted(req *tablesv1.GetFollowingCalculationsDeletedRequest) error {
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

// ValidateGetFollowingCalculationsFilters проверяет запрос страницы dc.following_calculation.
func ValidateGetFollowingCalculationsFilters(req *tablesv1.GetFollowingCalculationsFiltersRequest) error {
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
