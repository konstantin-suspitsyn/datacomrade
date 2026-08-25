package olapexchange

import (
	olapexchangev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/olap_exchange/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateGetOlapAliases проверяет запрос на выборку olap_exchange.GetOlapAliases.
func ValidateGetOlapAliases(req *olapexchangev1.GetOlapAliasesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("external_id", req.GetExternalId())
	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}
