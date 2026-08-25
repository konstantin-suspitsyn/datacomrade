package olapexchangeapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	olapexchangeconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/olapexchange"
	olapexchangevalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/olapexchange"
	olapexchangev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/olap_exchange/v1"
)

// GetOlapAliases оборачивает sqlc-запрос olap_exchange.GetOlapAliases.
func (o *OlapExchangeApiV1) GetOlapAliases(ctx context.Context, req *olapexchangev1.GetOlapAliasesRequest) (*olapexchangev1.GetOlapAliasesResponse, error) {
	if err := olapexchangevalidation.ValidateGetOlapAliases(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := o.services.OlapExchangeService.GetOlapAliases(ctx, olapexchangeconv.ToGetOlapAliasesParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return olapexchangeconv.GetOlapAliasesToProto(rows), nil
}
