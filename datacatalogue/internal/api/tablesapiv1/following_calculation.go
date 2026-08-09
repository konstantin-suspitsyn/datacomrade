package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetFollowingCalculationById отдаёт активную строку dc.following_calculation по id.
func (t *TablesApiV1) GetFollowingCalculationById(ctx context.Context, req *tablesv1.GetFollowingCalculationByIdRequest) (*tablesv1.GetFollowingCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetFollowingCalculationById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetFollowingCalculationByIdResponse{FollowingCalculation: tablesconv.FollowingCalculationToProto(row)}, nil
}

// GetFollowingCalculations отдаёт страницу строк dc.following_calculation.
func (t *TablesApiV1) GetFollowingCalculations(ctx context.Context, req *tablesv1.GetFollowingCalculationsRequest) (*tablesv1.GetFollowingCalculationsResponse, error) {
	if err := tablesvalidation.ValidateGetFollowingCalculations(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetFollowingCalculationsParams(req)

	rows, count, err := t.services.TablesService.GetFollowingCalculations(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetFollowingCalculationsResponse{
		Data: tablesconv.FollowingCalculationsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetFollowingCalculationsDeleted отдаёт страницу строк dc.following_calculation.
func (t *TablesApiV1) GetFollowingCalculationsDeleted(ctx context.Context, req *tablesv1.GetFollowingCalculationsDeletedRequest) (*tablesv1.GetFollowingCalculationsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetFollowingCalculationsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetFollowingCalculationsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetFollowingCalculationsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetFollowingCalculationsDeletedResponse{
		Data: tablesconv.FollowingCalculationsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetFollowingCalculationsFilters отдаёт страницу строк dc.following_calculation.
func (t *TablesApiV1) GetFollowingCalculationsFilters(ctx context.Context, req *tablesv1.GetFollowingCalculationsFiltersRequest) (*tablesv1.GetFollowingCalculationsFiltersResponse, error) {
	if err := tablesvalidation.ValidateGetFollowingCalculationsFilters(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetFollowingCalculationsFiltersParams(req)

	rows, count, err := t.services.TablesService.GetFollowingCalculationsFilters(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetFollowingCalculationsFiltersResponse{
		Data: tablesconv.FollowingCalculationsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateFollowingCalculation вставляет строку dc.following_calculation и отдаёт её целиком.
func (t *TablesApiV1) CreateFollowingCalculation(ctx context.Context, req *tablesv1.CreateFollowingCalculationRequest) (*tablesv1.CreateFollowingCalculationResponse, error) {
	if err := tablesvalidation.ValidateCreateFollowingCalculation(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateFollowingCalculation(ctx, tablesconv.ToCreateFollowingCalculationParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateFollowingCalculationResponse{FollowingCalculation: tablesconv.FollowingCalculationToProto(row)}, nil
}

// UpdateFollowingCalculationById обновляет строку dc.following_calculation.
func (t *TablesApiV1) UpdateFollowingCalculationById(ctx context.Context, req *tablesv1.UpdateFollowingCalculationByIdRequest) (*tablesv1.UpdateFollowingCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateFollowingCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateFollowingCalculationById(ctx, tablesconv.ToUpdateFollowingCalculationByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateFollowingCalculationByIdResponse{}, nil
}

// DeleteFollowingCalculationById мягко удаляет строку dc.following_calculation.
func (t *TablesApiV1) DeleteFollowingCalculationById(ctx context.Context, req *tablesv1.DeleteFollowingCalculationByIdRequest) (*tablesv1.DeleteFollowingCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteFollowingCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteFollowingCalculationById(ctx, tablesconv.ToDeleteFollowingCalculationByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteFollowingCalculationByIdResponse{}, nil
}

// UndeleteFollowingCalculationById восстанавливает мягко удалённую строку dc.following_calculation.
func (t *TablesApiV1) UndeleteFollowingCalculationById(ctx context.Context, req *tablesv1.UndeleteFollowingCalculationByIdRequest) (*tablesv1.UndeleteFollowingCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteFollowingCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteFollowingCalculationById(ctx, tablesconv.ToUndeleteFollowingCalculationByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteFollowingCalculationByIdResponse{}, nil
}
