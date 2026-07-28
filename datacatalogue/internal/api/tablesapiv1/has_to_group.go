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

// GetHasToGroupById отдаёт активную строку dc.has_to_group по id.
func (t *TablesApiV1) GetHasToGroupById(ctx context.Context, req *tablesv1.GetHasToGroupByIdRequest) (*tablesv1.GetHasToGroupByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetHasToGroupById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHasToGroupByIdResponse{HasToGroup: tablesconv.HasToGroupToProto(row)}, nil
}

// GetHasToGroups отдаёт все активные строки dc.has_to_group.
func (t *TablesApiV1) GetHasToGroups(ctx context.Context, req *tablesv1.GetHasToGroupsRequest) (*tablesv1.GetHasToGroupsResponse, error) {
	rows, err := t.services.TablesService.GetHasToGroups(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHasToGroupsResponse{HasToGroups: tablesconv.HasToGroupsToProto(rows)}, nil
}

// GetDeletedHasToGroupById отдаёт мягко удалённую строку dc.has_to_group по id.
func (t *TablesApiV1) GetDeletedHasToGroupById(ctx context.Context, req *tablesv1.GetDeletedHasToGroupByIdRequest) (*tablesv1.GetDeletedHasToGroupByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedHasToGroupById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedHasToGroupByIdResponse{HasToGroup: tablesconv.HasToGroupToProto(row)}, nil
}

// GetDeletedHasToGroups отдаёт все мягко удалённые строки dc.has_to_group.
func (t *TablesApiV1) GetDeletedHasToGroups(ctx context.Context, req *tablesv1.GetDeletedHasToGroupsRequest) (*tablesv1.GetDeletedHasToGroupsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedHasToGroups(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedHasToGroupsResponse{HasToGroups: tablesconv.HasToGroupsToProto(rows)}, nil
}

// CreateHasToGroup вставляет строку dc.has_to_group и отдаёт её целиком.
func (t *TablesApiV1) CreateHasToGroup(ctx context.Context, req *tablesv1.CreateHasToGroupRequest) (*tablesv1.CreateHasToGroupResponse, error) {
	if err := tablesvalidation.ValidateCreateHasToGroup(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateHasToGroup(ctx, tablesconv.ToCreateHasToGroupParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateHasToGroupResponse{HasToGroup: tablesconv.HasToGroupToProto(row)}, nil
}

// UpdateHasToGroupById обновляет активную строку dc.has_to_group и отдаёт её целиком.
func (t *TablesApiV1) UpdateHasToGroupById(ctx context.Context, req *tablesv1.UpdateHasToGroupByIdRequest) (*tablesv1.UpdateHasToGroupByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateHasToGroupById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateHasToGroupById(ctx, tablesconv.ToUpdateHasToGroupByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateHasToGroupByIdResponse{HasToGroup: tablesconv.HasToGroupToProto(row)}, nil
}

// DeleteHasToGroupById мягко удаляет строку dc.has_to_group.
func (t *TablesApiV1) DeleteHasToGroupById(ctx context.Context, req *tablesv1.DeleteHasToGroupByIdRequest) (*tablesv1.DeleteHasToGroupByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteHasToGroupById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteHasToGroupByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteHasToGroupById восстанавливает мягко удалённую строку dc.has_to_group.
func (t *TablesApiV1) UndeleteHasToGroupById(ctx context.Context, req *tablesv1.UndeleteHasToGroupByIdRequest) (*tablesv1.UndeleteHasToGroupByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteHasToGroupById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteHasToGroupByIdResponse{Empty: &emptypb.Empty{}}, nil
}
