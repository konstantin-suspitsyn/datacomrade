package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetGroupLevelsById отдаёт активную строку dc.group_levels по id.
func (t *TablesApiV1) GetGroupLevelsById(ctx context.Context, req *tablesv1.GetGroupLevelsByIdRequest) (*tablesv1.GetGroupLevelsByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetGroupLevelsById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetGroupLevelsByIdResponse{GroupLevels: tablesconv.GroupLevelsToProto(row)}, nil
}

// GetGroupLevelsSearchName отдаёт страницу строк dc.group_levels.
func (t *TablesApiV1) GetGroupLevelsSearchName(ctx context.Context, req *tablesv1.GetGroupLevelsSearchNameRequest) (*tablesv1.GetGroupLevelsSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetGroupLevelsSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetGroupLevelsSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetGroupLevelsSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetGroupLevelsSearchNameResponse{
		Data: tablesconv.GroupLevelsesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetGroupLevels отдаёт страницу строк dc.group_levels.
func (t *TablesApiV1) GetGroupLevels(ctx context.Context, req *tablesv1.GetGroupLevelsRequest) (*tablesv1.GetGroupLevelsResponse, error) {
	if err := tablesvalidation.ValidateGetGroupLevels(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetGroupLevelsParams(req)

	rows, count, err := t.services.TablesService.GetGroupLevels(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetGroupLevelsResponse{
		Data: tablesconv.GroupLevelsesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetGroupLevelsFiltered отдаёт страницу строк dc.group_levels.
func (t *TablesApiV1) GetGroupLevelsFiltered(ctx context.Context, req *tablesv1.GetGroupLevelsFilteredRequest) (*tablesv1.GetGroupLevelsFilteredResponse, error) {
	if err := tablesvalidation.ValidateGetGroupLevelsFiltered(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetGroupLevelsFilteredParams(req)

	rows, count, err := t.services.TablesService.GetGroupLevelsFiltered(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetGroupLevelsFilteredResponse{
		Data: tablesconv.GroupLevelsesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetGroupLevelsDeleted отдаёт строки dc.group_levels.
func (t *TablesApiV1) GetGroupLevelsDeleted(ctx context.Context, req *tablesv1.GetGroupLevelsDeletedRequest) (*tablesv1.GetGroupLevelsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetGroupLevelsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetGroupLevelsDeleted(ctx, req.GetOrder())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetGroupLevelsDeletedResponse{Rows: tablesconv.GroupLevelsesToProto(rows)}, nil
}

// CreateGroupLevels вставляет строку dc.group_levels и отдаёт её целиком.
func (t *TablesApiV1) CreateGroupLevels(ctx context.Context, req *tablesv1.CreateGroupLevelsRequest) (*tablesv1.CreateGroupLevelsResponse, error) {
	if err := tablesvalidation.ValidateCreateGroupLevels(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateGroupLevels(ctx, tablesconv.ToCreateGroupLevelsParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateGroupLevelsResponse{GroupLevels: tablesconv.GroupLevelsToProto(row)}, nil
}

// UpdateGroupLevelsById обновляет строку dc.group_levels.
func (t *TablesApiV1) UpdateGroupLevelsById(ctx context.Context, req *tablesv1.UpdateGroupLevelsByIdRequest) (*tablesv1.UpdateGroupLevelsByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateGroupLevelsById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateGroupLevelsById(ctx, tablesconv.ToUpdateGroupLevelsByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateGroupLevelsByIdResponse{}, nil
}

// DeleteGroupLevelsById мягко удаляет строку dc.group_levels.
func (t *TablesApiV1) DeleteGroupLevelsById(ctx context.Context, req *tablesv1.DeleteGroupLevelsByIdRequest) (*tablesv1.DeleteGroupLevelsByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteGroupLevelsById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteGroupLevelsById(ctx, tablesconv.ToDeleteGroupLevelsByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteGroupLevelsByIdResponse{}, nil
}

// UndeleteGroupLevelsById восстанавливает мягко удалённую строку dc.group_levels.
func (t *TablesApiV1) UndeleteGroupLevelsById(ctx context.Context, req *tablesv1.UndeleteGroupLevelsByIdRequest) (*tablesv1.UndeleteGroupLevelsByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteGroupLevelsById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteGroupLevelsById(ctx, tablesconv.ToUndeleteGroupLevelsByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteGroupLevelsByIdResponse{}, nil
}
