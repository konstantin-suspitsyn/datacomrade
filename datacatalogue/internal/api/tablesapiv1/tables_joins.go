package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetTablesJoinsById отдаёт активную строку dc.tables_joins по id.
func (t *TablesApiV1) GetTablesJoinsById(ctx context.Context, req *tablesv1.GetTablesJoinsByIdRequest) (*tablesv1.GetTablesJoinsByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetTablesJoinsById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTablesJoinsByIdResponse{TablesJoins: tablesconv.GetTablesJoinsByIdRowToProto(row)}, nil
}

// GetTablesJoins отдаёт страницу строк dc.tables_joins.
func (t *TablesApiV1) GetTablesJoins(ctx context.Context, req *tablesv1.GetTablesJoinsRequest) (*tablesv1.GetTablesJoinsResponse, error) {
	if err := tablesvalidation.ValidateGetTablesJoins(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetTablesJoinsParams(req)

	rows, count, err := t.services.TablesService.GetTablesJoins(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTablesJoinsResponse{
		Data: tablesconv.GetTablesJoinsRowsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateTablesJoins вставляет строку dc.tables_joins и отдаёт её целиком.
func (t *TablesApiV1) CreateTablesJoins(ctx context.Context, req *tablesv1.CreateTablesJoinsRequest) (*tablesv1.CreateTablesJoinsResponse, error) {
	if err := tablesvalidation.ValidateCreateTablesJoins(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateTablesJoins(ctx, tablesconv.ToCreateTablesJoinsParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateTablesJoinsResponse{TablesJoins: tablesconv.TablesJoinsToProto(row)}, nil
}

// UpdateTablesJoinsById обновляет строку dc.tables_joins.
func (t *TablesApiV1) UpdateTablesJoinsById(ctx context.Context, req *tablesv1.UpdateTablesJoinsByIdRequest) (*tablesv1.UpdateTablesJoinsByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateTablesJoinsById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateTablesJoinsById(ctx, tablesconv.ToUpdateTablesJoinsByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateTablesJoinsByIdResponse{}, nil
}

// DeleteTablesJoinsById мягко удаляет строку dc.tables_joins.
func (t *TablesApiV1) DeleteTablesJoinsById(ctx context.Context, req *tablesv1.DeleteTablesJoinsByIdRequest) (*tablesv1.DeleteTablesJoinsByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteTablesJoinsById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteTablesJoinsById(ctx, tablesconv.ToDeleteTablesJoinsByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteTablesJoinsByIdResponse{}, nil
}

// UndeleteTablesJoinsById восстанавливает мягко удалённую строку dc.tables_joins.
func (t *TablesApiV1) UndeleteTablesJoinsById(ctx context.Context, req *tablesv1.UndeleteTablesJoinsByIdRequest) (*tablesv1.UndeleteTablesJoinsByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteTablesJoinsById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteTablesJoinsByIdResponse{}, nil
}
