package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"google.golang.org/protobuf/types/known/emptypb"
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

// GetCalculationTypes отдаёт все активные строки dc.calculation_type.
func (t *TablesApiV1) GetCalculationTypes(ctx context.Context, req *tablesv1.GetCalculationTypesRequest) (*tablesv1.GetCalculationTypesResponse, error) {
	rows, err := t.services.TablesService.GetCalculationTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetCalculationTypesResponse{CalculationTypes: tablesconv.CalculationTypesToProto(rows)}, nil
}

// GetDeletedCalculationTypeById отдаёт мягко удалённую строку dc.calculation_type по id.
func (t *TablesApiV1) GetDeletedCalculationTypeById(ctx context.Context, req *tablesv1.GetDeletedCalculationTypeByIdRequest) (*tablesv1.GetDeletedCalculationTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedCalculationTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedCalculationTypeByIdResponse{CalculationType: tablesconv.CalculationTypeToProto(row)}, nil
}

// GetDeletedCalculationTypes отдаёт все мягко удалённые строки dc.calculation_type.
func (t *TablesApiV1) GetDeletedCalculationTypes(ctx context.Context, req *tablesv1.GetDeletedCalculationTypesRequest) (*tablesv1.GetDeletedCalculationTypesResponse, error) {
	rows, err := t.services.TablesService.GetDeletedCalculationTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedCalculationTypesResponse{CalculationTypes: tablesconv.CalculationTypesToProto(rows)}, nil
}

// CreateCalculationType вставляет строку dc.calculation_type и отдаёт её целиком.
func (t *TablesApiV1) CreateCalculationType(ctx context.Context, req *tablesv1.CreateCalculationTypeRequest) (*tablesv1.CreateCalculationTypeResponse, error) {
	if err := tablesvalidation.ValidateCreateCalculationType(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateCalculationType(ctx, tablesconv.ToCreateCalculationTypeParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateCalculationTypeResponse{CalculationType: tablesconv.CalculationTypeToProto(row)}, nil
}

// UpdateCalculationTypeById обновляет активную строку dc.calculation_type и отдаёт её целиком.
func (t *TablesApiV1) UpdateCalculationTypeById(ctx context.Context, req *tablesv1.UpdateCalculationTypeByIdRequest) (*tablesv1.UpdateCalculationTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateCalculationTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateCalculationTypeById(ctx, tablesconv.ToUpdateCalculationTypeByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateCalculationTypeByIdResponse{CalculationType: tablesconv.CalculationTypeToProto(row)}, nil
}

// DeleteCalculationTypeById мягко удаляет строку dc.calculation_type.
func (t *TablesApiV1) DeleteCalculationTypeById(ctx context.Context, req *tablesv1.DeleteCalculationTypeByIdRequest) (*tablesv1.DeleteCalculationTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteCalculationTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteCalculationTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteCalculationTypeById восстанавливает мягко удалённую строку dc.calculation_type.
func (t *TablesApiV1) UndeleteCalculationTypeById(ctx context.Context, req *tablesv1.UndeleteCalculationTypeByIdRequest) (*tablesv1.UndeleteCalculationTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteCalculationTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteCalculationTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}
