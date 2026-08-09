package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetColumnCatById отдаёт активную строку dc.column_cat по id.
func (t *TablesApiV1) GetColumnCatById(ctx context.Context, req *tablesv1.GetColumnCatByIdRequest) (*tablesv1.GetColumnCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetColumnCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatByIdResponse{ColumnCat: tablesconv.ColumnCatToProto(row)}, nil
}

// GetColumnCats отдаёт страницу строк dc.column_cat.
func (t *TablesApiV1) GetColumnCats(ctx context.Context, req *tablesv1.GetColumnCatsRequest) (*tablesv1.GetColumnCatsResponse, error) {
	if err := tablesvalidation.ValidateGetColumnCats(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetColumnCatsParams(req)

	rows, count, err := t.services.TablesService.GetColumnCats(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatsResponse{
		Data: tablesconv.ColumnCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetColumnCatsDeleted отдаёт страницу строк dc.column_cat.
func (t *TablesApiV1) GetColumnCatsDeleted(ctx context.Context, req *tablesv1.GetColumnCatsDeletedRequest) (*tablesv1.GetColumnCatsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetColumnCatsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetColumnCatsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetColumnCatsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatsDeletedResponse{
		Data: tablesconv.ColumnCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetColumnCatsOptionalFilters отдаёт страницу строк dc.column_cat.
func (t *TablesApiV1) GetColumnCatsOptionalFilters(ctx context.Context, req *tablesv1.GetColumnCatsOptionalFiltersRequest) (*tablesv1.GetColumnCatsOptionalFiltersResponse, error) {
	if err := tablesvalidation.ValidateGetColumnCatsOptionalFilters(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetColumnCatsOptionalFiltersParams(req)

	rows, count, err := t.services.TablesService.GetColumnCatsOptionalFilters(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatsOptionalFiltersResponse{
		Data: tablesconv.ColumnCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetColumnCatsSearchName отдаёт страницу строк dc.column_cat.
func (t *TablesApiV1) GetColumnCatsSearchName(ctx context.Context, req *tablesv1.GetColumnCatsSearchNameRequest) (*tablesv1.GetColumnCatsSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetColumnCatsSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetColumnCatsSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetColumnCatsSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatsSearchNameResponse{
		Data: tablesconv.ColumnCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateColumnCat вставляет строку dc.column_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateColumnCat(ctx context.Context, req *tablesv1.CreateColumnCatRequest) (*tablesv1.CreateColumnCatResponse, error) {
	if err := tablesvalidation.ValidateCreateColumnCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateColumnCat(ctx, tablesconv.ToCreateColumnCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateColumnCatResponse{ColumnCat: tablesconv.ColumnCatToProto(row)}, nil
}

// UpdateColumnCatById обновляет строку dc.column_cat.
func (t *TablesApiV1) UpdateColumnCatById(ctx context.Context, req *tablesv1.UpdateColumnCatByIdRequest) (*tablesv1.UpdateColumnCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateColumnCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateColumnCatById(ctx, tablesconv.ToUpdateColumnCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateColumnCatByIdResponse{}, nil
}

// DeleteColumnCatById мягко удаляет строку dc.column_cat.
func (t *TablesApiV1) DeleteColumnCatById(ctx context.Context, req *tablesv1.DeleteColumnCatByIdRequest) (*tablesv1.DeleteColumnCatByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteColumnCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteColumnCatById(ctx, tablesconv.ToDeleteColumnCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteColumnCatByIdResponse{}, nil
}

// UndeleteColumnCatById восстанавливает мягко удалённую строку dc.column_cat.
func (t *TablesApiV1) UndeleteColumnCatById(ctx context.Context, req *tablesv1.UndeleteColumnCatByIdRequest) (*tablesv1.UndeleteColumnCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteColumnCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteColumnCatById(ctx, tablesconv.ToUndeleteColumnCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteColumnCatByIdResponse{}, nil
}
