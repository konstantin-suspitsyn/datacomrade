package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetTableCatById отдаёт активную строку dc.table_cat по id.
func (t *TablesApiV1) GetTableCatById(ctx context.Context, req *tablesv1.GetTableCatByIdRequest) (*tablesv1.GetTableCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetTableCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatByIdResponse{TableCat: tablesconv.TableCatToProto(row)}, nil
}

// GetTableCats отдаёт страницу строк dc.table_cat.
func (t *TablesApiV1) GetTableCats(ctx context.Context, req *tablesv1.GetTableCatsRequest) (*tablesv1.GetTableCatsResponse, error) {
	if err := tablesvalidation.ValidateGetTableCats(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetTableCatsParams(req)

	rows, count, err := t.services.TablesService.GetTableCats(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatsResponse{
		Data: tablesconv.TableCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetTableCatsDeleted отдаёт страницу строк dc.table_cat.
func (t *TablesApiV1) GetTableCatsDeleted(ctx context.Context, req *tablesv1.GetTableCatsDeletedRequest) (*tablesv1.GetTableCatsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetTableCatsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetTableCatsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetTableCatsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatsDeletedResponse{
		Data: tablesconv.TableCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateTableCat вставляет строку dc.table_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateTableCat(ctx context.Context, req *tablesv1.CreateTableCatRequest) (*tablesv1.CreateTableCatResponse, error) {
	if err := tablesvalidation.ValidateCreateTableCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateTableCat(ctx, tablesconv.ToCreateTableCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateTableCatResponse{TableCat: tablesconv.TableCatToProto(row)}, nil
}

// UpdateTableCatById обновляет строку dc.table_cat.
func (t *TablesApiV1) UpdateTableCatById(ctx context.Context, req *tablesv1.UpdateTableCatByIdRequest) (*tablesv1.UpdateTableCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateTableCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateTableCatById(ctx, tablesconv.ToUpdateTableCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateTableCatByIdResponse{}, nil
}

// DeleteTableCatById мягко удаляет строку dc.table_cat.
func (t *TablesApiV1) DeleteTableCatById(ctx context.Context, req *tablesv1.DeleteTableCatByIdRequest) (*tablesv1.DeleteTableCatByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteTableCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteTableCatById(ctx, tablesconv.ToDeleteTableCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteTableCatByIdResponse{}, nil
}

// UndeleteTableCatById восстанавливает мягко удалённую строку dc.table_cat.
func (t *TablesApiV1) UndeleteTableCatById(ctx context.Context, req *tablesv1.UndeleteTableCatByIdRequest) (*tablesv1.UndeleteTableCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteTableCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteTableCatById(ctx, tablesconv.ToUndeleteTableCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteTableCatByIdResponse{}, nil
}
