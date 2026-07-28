package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// calculationTypeWritableFields проверяет поля, общие для вставки и обновления dc.calculation_type.
func calculationTypeWritableFields(
	v *validator.Validator,
	name string,
	description string,
) {
	v.StringVarchar("name", name, calculationTypeNameMaxLen)
	v.StringVarchar("description", description, calculationTypeDescriptionMaxLen)
}

// ValidateCreateCalculationType проверяет запрос на вставку строки dc.calculation_type.
func ValidateCreateCalculationType(req *tablesv1.CreateCalculationTypeRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	calculationTypeWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
	)

	return v.Err()
}

// ValidateUpdateCalculationTypeById проверяет запрос на обновление строки dc.calculation_type.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateCalculationTypeById(req *tablesv1.UpdateCalculationTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	calculationTypeWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
	)

	return v.Err()
}
