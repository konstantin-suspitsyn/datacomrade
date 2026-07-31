package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// followingCalculationWritableFields проверяет поля, общие для вставки и обновления dc.following_calculation.
func followingCalculationWritableFields(
	v *validator.Validator,
	columnCatId int64,
	calculationTypeId int64,
	userExternalId string,
) {
	v.Int64ID("column_cat_id", columnCatId)
	v.Int64ID("calculation_type_id", calculationTypeId)
	v.StringUUID("user_id", userExternalId)
}

// ValidateCreateFollowingCalculation проверяет запрос на вставку строки dc.following_calculation.
func ValidateCreateFollowingCalculation(req *tablesv1.CreateFollowingCalculationRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	followingCalculationWritableFields(
		v,
		req.GetColumnCatId(),
		req.GetCalculationTypeId(),
		req.GetUserExternalId(),
	)

	return v.Err()
}

// ValidateUpdateFollowingCalculationById проверяет запрос на обновление строки dc.following_calculation.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateFollowingCalculationById(req *tablesv1.UpdateFollowingCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	followingCalculationWritableFields(
		v,
		req.GetColumnCatId(),
		req.GetCalculationTypeId(),
		req.GetUserExternalId(),
	)

	return v.Err()
}
