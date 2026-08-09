package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetCalculationTypeById отдаёт активную строку dc.calculation_type по id.
func (t *TablesApiV1) GetCalculationTypeById(ctx context.Context, req *tablesv1.GetCalculationTypeByIdRequest) (*tablesv1.GetCalculationTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetCalculationTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetCalculationTypeByIdResponse{CalculationType: tablesconv.CalculationTypeToProto(row)}, nil
}

// GetCalculationTypes отдаёт строки dc.calculation_type.
func (t *TablesApiV1) GetCalculationTypes(ctx context.Context, req *tablesv1.GetCalculationTypesRequest) (*tablesv1.GetCalculationTypesResponse, error) {
	if err := tablesvalidation.ValidateGetCalculationTypes(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetCalculationTypes(ctx, req.GetOrder())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetCalculationTypesResponse{Rows: tablesconv.CalculationTypesToProto(rows)}, nil
}
