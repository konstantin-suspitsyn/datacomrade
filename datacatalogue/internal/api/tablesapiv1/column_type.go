package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetColumnTypeById отдаёт активную строку dc.column_type по id.
func (t *TablesApiV1) GetColumnTypeById(ctx context.Context, req *tablesv1.GetColumnTypeByIdRequest) (*tablesv1.GetColumnTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetColumnTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnTypeByIdResponse{ColumnType: tablesconv.ColumnTypeToProto(row)}, nil
}

// GetColumnTypes отдаёт страницу строк dc.column_type.
func (t *TablesApiV1) GetColumnTypes(ctx context.Context, req *tablesv1.GetColumnTypesRequest) (*tablesv1.GetColumnTypesResponse, error) {
	if err := tablesvalidation.ValidateGetColumnTypes(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetColumnTypesParams(req)

	rows, count, err := t.services.TablesService.GetColumnTypes(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnTypesResponse{
		Data: tablesconv.ColumnTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetColumnTypesSearchName отдаёт страницу строк dc.column_type.
func (t *TablesApiV1) GetColumnTypesSearchName(ctx context.Context, req *tablesv1.GetColumnTypesSearchNameRequest) (*tablesv1.GetColumnTypesSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetColumnTypesSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetColumnTypesSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetColumnTypesSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnTypesSearchNameResponse{
		Data: tablesconv.ColumnTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetColumnTypesDeleted отдаёт страницу строк dc.column_type.
func (t *TablesApiV1) GetColumnTypesDeleted(ctx context.Context, req *tablesv1.GetColumnTypesDeletedRequest) (*tablesv1.GetColumnTypesDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetColumnTypesDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetColumnTypesDeletedParams(req)

	rows, count, err := t.services.TablesService.GetColumnTypesDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnTypesDeletedResponse{
		Data: tablesconv.ColumnTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateColumnType вставляет строку dc.column_type и отдаёт её целиком.
func (t *TablesApiV1) CreateColumnType(ctx context.Context, req *tablesv1.CreateColumnTypeRequest) (*tablesv1.CreateColumnTypeResponse, error) {
	if err := tablesvalidation.ValidateCreateColumnType(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateColumnType(ctx, tablesconv.ToCreateColumnTypeParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateColumnTypeResponse{ColumnType: tablesconv.ColumnTypeToProto(row)}, nil
}

// UpdateColumnTypeById обновляет строку dc.column_type.
func (t *TablesApiV1) UpdateColumnTypeById(ctx context.Context, req *tablesv1.UpdateColumnTypeByIdRequest) (*tablesv1.UpdateColumnTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateColumnTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateColumnTypeById(ctx, tablesconv.ToUpdateColumnTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateColumnTypeByIdResponse{}, nil
}

// DeleteColumnTypeById мягко удаляет строку dc.column_type.
func (t *TablesApiV1) DeleteColumnTypeById(ctx context.Context, req *tablesv1.DeleteColumnTypeByIdRequest) (*tablesv1.DeleteColumnTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteColumnTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteColumnTypeById(ctx, tablesconv.ToDeleteColumnTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteColumnTypeByIdResponse{}, nil
}

// UndeleteColumnTypeById восстанавливает мягко удалённую строку dc.column_type.
func (t *TablesApiV1) UndeleteColumnTypeById(ctx context.Context, req *tablesv1.UndeleteColumnTypeByIdRequest) (*tablesv1.UndeleteColumnTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteColumnTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteColumnTypeById(ctx, tablesconv.ToUndeleteColumnTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteColumnTypeByIdResponse{}, nil
}
