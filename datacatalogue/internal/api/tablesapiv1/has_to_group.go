package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
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

// GetHasToGroupsSearchName отдаёт страницу строк dc.has_to_group.
func (t *TablesApiV1) GetHasToGroupsSearchName(ctx context.Context, req *tablesv1.GetHasToGroupsSearchNameRequest) (*tablesv1.GetHasToGroupsSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetHasToGroupsSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetHasToGroupsSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetHasToGroupsSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHasToGroupsSearchNameResponse{
		Data: tablesconv.HasToGroupsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetHasToGroups отдаёт страницу строк dc.has_to_group.
func (t *TablesApiV1) GetHasToGroups(ctx context.Context, req *tablesv1.GetHasToGroupsRequest) (*tablesv1.GetHasToGroupsResponse, error) {
	if err := tablesvalidation.ValidateGetHasToGroups(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetHasToGroupsParams(req)

	rows, count, err := t.services.TablesService.GetHasToGroups(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHasToGroupsResponse{
		Data: tablesconv.HasToGroupsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetHasToGroupsDeleted отдаёт страницу строк dc.has_to_group.
func (t *TablesApiV1) GetHasToGroupsDeleted(ctx context.Context, req *tablesv1.GetHasToGroupsDeletedRequest) (*tablesv1.GetHasToGroupsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetHasToGroupsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetHasToGroupsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetHasToGroupsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHasToGroupsDeletedResponse{
		Data: tablesconv.HasToGroupsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
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

// UpdateHasToGroupById обновляет строку dc.has_to_group.
func (t *TablesApiV1) UpdateHasToGroupById(ctx context.Context, req *tablesv1.UpdateHasToGroupByIdRequest) (*tablesv1.UpdateHasToGroupByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateHasToGroupById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateHasToGroupById(ctx, tablesconv.ToUpdateHasToGroupByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateHasToGroupByIdResponse{}, nil
}

// DeleteHasToGroupById мягко удаляет строку dc.has_to_group.
func (t *TablesApiV1) DeleteHasToGroupById(ctx context.Context, req *tablesv1.DeleteHasToGroupByIdRequest) (*tablesv1.DeleteHasToGroupByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteHasToGroupById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteHasToGroupById(ctx, tablesconv.ToDeleteHasToGroupByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteHasToGroupByIdResponse{}, nil
}

// UndeleteHasToGroupById восстанавливает мягко удалённую строку dc.has_to_group.
func (t *TablesApiV1) UndeleteHasToGroupById(ctx context.Context, req *tablesv1.UndeleteHasToGroupByIdRequest) (*tablesv1.UndeleteHasToGroupByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteHasToGroupById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteHasToGroupById(ctx, tablesconv.ToUndeleteHasToGroupByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteHasToGroupByIdResponse{}, nil
}
