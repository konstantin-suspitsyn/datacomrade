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

// GetDatabaseCalculationById отдаёт активную строку dc.database_calculation по id.
func (t *TablesApiV1) GetDatabaseCalculationById(ctx context.Context, req *tablesv1.GetDatabaseCalculationByIdRequest) (*tablesv1.GetDatabaseCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDatabaseCalculationById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCalculationByIdResponse{DatabaseCalculation: tablesconv.DatabaseCalculationToProto(row)}, nil
}

// GetDatabaseCalculations отдаёт все активные строки dc.database_calculation.
func (t *TablesApiV1) GetDatabaseCalculations(ctx context.Context, req *tablesv1.GetDatabaseCalculationsRequest) (*tablesv1.GetDatabaseCalculationsResponse, error) {
	rows, err := t.services.TablesService.GetDatabaseCalculations(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCalculationsResponse{DatabaseCalculations: tablesconv.DatabaseCalculationsToProto(rows)}, nil
}

// GetDeletedDatabaseCalculationById отдаёт мягко удалённую строку dc.database_calculation по id.
func (t *TablesApiV1) GetDeletedDatabaseCalculationById(ctx context.Context, req *tablesv1.GetDeletedDatabaseCalculationByIdRequest) (*tablesv1.GetDeletedDatabaseCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedDatabaseCalculationById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDatabaseCalculationByIdResponse{DatabaseCalculation: tablesconv.DatabaseCalculationToProto(row)}, nil
}

// GetDeletedDatabaseCalculations отдаёт все мягко удалённые строки dc.database_calculation.
func (t *TablesApiV1) GetDeletedDatabaseCalculations(ctx context.Context, req *tablesv1.GetDeletedDatabaseCalculationsRequest) (*tablesv1.GetDeletedDatabaseCalculationsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedDatabaseCalculations(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDatabaseCalculationsResponse{DatabaseCalculations: tablesconv.DatabaseCalculationsToProto(rows)}, nil
}

// CreateDatabaseCalculation вставляет строку dc.database_calculation и отдаёт её целиком.
func (t *TablesApiV1) CreateDatabaseCalculation(ctx context.Context, req *tablesv1.CreateDatabaseCalculationRequest) (*tablesv1.CreateDatabaseCalculationResponse, error) {
	if err := tablesvalidation.ValidateCreateDatabaseCalculation(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateDatabaseCalculation(ctx, tablesconv.ToCreateDatabaseCalculationParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateDatabaseCalculationResponse{DatabaseCalculation: tablesconv.DatabaseCalculationToProto(row)}, nil
}

// UpdateDatabaseCalculationById обновляет активную строку dc.database_calculation и отдаёт её целиком.
func (t *TablesApiV1) UpdateDatabaseCalculationById(ctx context.Context, req *tablesv1.UpdateDatabaseCalculationByIdRequest) (*tablesv1.UpdateDatabaseCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDatabaseCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateDatabaseCalculationById(ctx, tablesconv.ToUpdateDatabaseCalculationByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDatabaseCalculationByIdResponse{DatabaseCalculation: tablesconv.DatabaseCalculationToProto(row)}, nil
}

// DeleteDatabaseCalculationById мягко удаляет строку dc.database_calculation.
func (t *TablesApiV1) DeleteDatabaseCalculationById(ctx context.Context, req *tablesv1.DeleteDatabaseCalculationByIdRequest) (*tablesv1.DeleteDatabaseCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDatabaseCalculationById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDatabaseCalculationByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteDatabaseCalculationById восстанавливает мягко удалённую строку dc.database_calculation.
func (t *TablesApiV1) UndeleteDatabaseCalculationById(ctx context.Context, req *tablesv1.UndeleteDatabaseCalculationByIdRequest) (*tablesv1.UndeleteDatabaseCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDatabaseCalculationById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDatabaseCalculationByIdResponse{Empty: &emptypb.Empty{}}, nil
}
