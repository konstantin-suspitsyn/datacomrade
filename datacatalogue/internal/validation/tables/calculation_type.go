package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateGetCalculationTypes проверяет запрос списка dc.calculation_type.
func ValidateGetCalculationTypes(req *tablesv1.GetCalculationTypesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}
