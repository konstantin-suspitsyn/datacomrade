package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
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

// GetDatabaseCalculationsSearchName отдаёт страницу строк dc.database_calculation.
func (t *TablesApiV1) GetDatabaseCalculationsSearchName(ctx context.Context, req *tablesv1.GetDatabaseCalculationsSearchNameRequest) (*tablesv1.GetDatabaseCalculationsSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCalculationsSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCalculationsSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCalculationsSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCalculationsSearchNameResponse{
		Data: tablesconv.DatabaseCalculationsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseCalculationsFiltered отдаёт страницу строк dc.database_calculation.
func (t *TablesApiV1) GetDatabaseCalculationsFiltered(ctx context.Context, req *tablesv1.GetDatabaseCalculationsFilteredRequest) (*tablesv1.GetDatabaseCalculationsFilteredResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCalculationsFiltered(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCalculationsFilteredParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCalculationsFiltered(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCalculationsFilteredResponse{
		Data: tablesconv.DatabaseCalculationsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseCalculationsDeleted отдаёт страницу строк dc.database_calculation.
func (t *TablesApiV1) GetDatabaseCalculationsDeleted(ctx context.Context, req *tablesv1.GetDatabaseCalculationsDeletedRequest) (*tablesv1.GetDatabaseCalculationsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCalculationsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCalculationsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCalculationsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCalculationsDeletedResponse{
		Data: tablesconv.DatabaseCalculationsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseCalculations отдаёт страницу строк dc.database_calculation.
func (t *TablesApiV1) GetDatabaseCalculations(ctx context.Context, req *tablesv1.GetDatabaseCalculationsRequest) (*tablesv1.GetDatabaseCalculationsResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCalculations(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCalculationsParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCalculations(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCalculationsResponse{
		Data: tablesconv.DatabaseCalculationsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
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

// UpdateDatabaseCalculationById обновляет строку dc.database_calculation.
func (t *TablesApiV1) UpdateDatabaseCalculationById(ctx context.Context, req *tablesv1.UpdateDatabaseCalculationByIdRequest) (*tablesv1.UpdateDatabaseCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDatabaseCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateDatabaseCalculationById(ctx, tablesconv.ToUpdateDatabaseCalculationByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDatabaseCalculationByIdResponse{}, nil
}

// DeleteDatabaseCalculationById мягко удаляет строку dc.database_calculation.
func (t *TablesApiV1) DeleteDatabaseCalculationById(ctx context.Context, req *tablesv1.DeleteDatabaseCalculationByIdRequest) (*tablesv1.DeleteDatabaseCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteDatabaseCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDatabaseCalculationById(ctx, tablesconv.ToDeleteDatabaseCalculationByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDatabaseCalculationByIdResponse{}, nil
}

// UndeleteDatabaseCalculationById восстанавливает мягко удалённую строку dc.database_calculation.
func (t *TablesApiV1) UndeleteDatabaseCalculationById(ctx context.Context, req *tablesv1.UndeleteDatabaseCalculationByIdRequest) (*tablesv1.UndeleteDatabaseCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteDatabaseCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDatabaseCalculationById(ctx, tablesconv.ToUndeleteDatabaseCalculationByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDatabaseCalculationByIdResponse{}, nil
}
