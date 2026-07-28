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

// GetFollowingCalculations отдаёт все активные строки dc.following_calculation.
func (t *TablesApiV1) GetFollowingCalculations(ctx context.Context, req *tablesv1.GetFollowingCalculationsRequest) (*tablesv1.GetFollowingCalculationsResponse, error) {
	rows, err := t.services.TablesService.GetFollowingCalculations(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetFollowingCalculationsResponse{FollowingCalculations: tablesconv.FollowingCalculationsToProto(rows)}, nil
}

// GetDeletedFollowingCalculationById отдаёт мягко удалённую строку dc.following_calculation по id.
func (t *TablesApiV1) GetDeletedFollowingCalculationById(ctx context.Context, req *tablesv1.GetDeletedFollowingCalculationByIdRequest) (*tablesv1.GetDeletedFollowingCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedFollowingCalculationById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedFollowingCalculationByIdResponse{FollowingCalculation: tablesconv.FollowingCalculationToProto(row)}, nil
}

// GetDeletedFollowingCalculations отдаёт все мягко удалённые строки dc.following_calculation.
func (t *TablesApiV1) GetDeletedFollowingCalculations(ctx context.Context, req *tablesv1.GetDeletedFollowingCalculationsRequest) (*tablesv1.GetDeletedFollowingCalculationsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedFollowingCalculations(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedFollowingCalculationsResponse{FollowingCalculations: tablesconv.FollowingCalculationsToProto(rows)}, nil
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

// UpdateFollowingCalculationById обновляет активную строку dc.following_calculation и отдаёт её целиком.
func (t *TablesApiV1) UpdateFollowingCalculationById(ctx context.Context, req *tablesv1.UpdateFollowingCalculationByIdRequest) (*tablesv1.UpdateFollowingCalculationByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateFollowingCalculationById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateFollowingCalculationById(ctx, tablesconv.ToUpdateFollowingCalculationByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateFollowingCalculationByIdResponse{FollowingCalculation: tablesconv.FollowingCalculationToProto(row)}, nil
}

// DeleteFollowingCalculationById мягко удаляет строку dc.following_calculation.
func (t *TablesApiV1) DeleteFollowingCalculationById(ctx context.Context, req *tablesv1.DeleteFollowingCalculationByIdRequest) (*tablesv1.DeleteFollowingCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteFollowingCalculationById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteFollowingCalculationByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteFollowingCalculationById восстанавливает мягко удалённую строку dc.following_calculation.
func (t *TablesApiV1) UndeleteFollowingCalculationById(ctx context.Context, req *tablesv1.UndeleteFollowingCalculationByIdRequest) (*tablesv1.UndeleteFollowingCalculationByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteFollowingCalculationById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteFollowingCalculationByIdResponse{Empty: &emptypb.Empty{}}, nil
}
