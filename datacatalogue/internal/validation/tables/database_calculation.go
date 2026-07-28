package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// databaseCalculationWritableFields проверяет поля, общие для вставки и обновления dc.database_calculation.
func databaseCalculationWritableFields(
	v *validator.Validator,
	databaseCatId int64,
	calculationTypeId int64,
	userId int64,
) {
	v.Int64ID("database_cat_id", databaseCatId)
	v.Int64ID("calculation_type_id", calculationTypeId)
	v.Int64ID("user_id", userId)
}

// ValidateCreateDatabaseCalculation проверяет запрос на вставку строки dc.database_calculation.
func ValidateCreateDatabaseCalculation(req *tablesv1.CreateDatabaseCalculationRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	databaseCalculationWritableFields(
		v,
		req.GetDatabaseCatId(),
		req.GetCalculationTypeId(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateDatabaseCalculationById проверяет запрос на обновление строки dc.database_calculation.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateDatabaseCalculationById(req *tablesv1.UpdateDatabaseCalculationByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	databaseCalculationWritableFields(
		v,
		req.GetDatabaseCatId(),
		req.GetCalculationTypeId(),
		req.GetUserId(),
	)

	return v.Err()
}
