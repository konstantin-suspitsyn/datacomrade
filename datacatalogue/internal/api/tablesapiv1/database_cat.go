package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetDatabaseCatById отдаёт активную строку dc.database_cat по id.
func (t *TablesApiV1) GetDatabaseCatById(ctx context.Context, req *tablesv1.GetDatabaseCatByIdRequest) (*tablesv1.GetDatabaseCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDatabaseCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatByIdResponse{DatabaseCat: tablesconv.DatabaseCatToProto(row)}, nil
}

// GetDatabaseCats отдаёт страницу строк dc.database_cat.
func (t *TablesApiV1) GetDatabaseCats(ctx context.Context, req *tablesv1.GetDatabaseCatsRequest) (*tablesv1.GetDatabaseCatsResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCats(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCatsParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCats(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatsResponse{
		Data: tablesconv.DatabaseCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseCatsDeleted отдаёт страницу строк dc.database_cat.
func (t *TablesApiV1) GetDatabaseCatsDeleted(ctx context.Context, req *tablesv1.GetDatabaseCatsDeletedRequest) (*tablesv1.GetDatabaseCatsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCatsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCatsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCatsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatsDeletedResponse{
		Data: tablesconv.DatabaseCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseCatsByHostIdDatabaseTypeId отдаёт страницу строк dc.database_cat.
func (t *TablesApiV1) GetDatabaseCatsByHostIdDatabaseTypeId(ctx context.Context, req *tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest) (*tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCatsByHostIdDatabaseTypeId(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCatsByHostIdDatabaseTypeIdParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCatsByHostIdDatabaseTypeId(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdResponse{
		Data: tablesconv.DatabaseCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseCatsSearchName отдаёт страницу строк dc.database_cat.
func (t *TablesApiV1) GetDatabaseCatsSearchName(ctx context.Context, req *tablesv1.GetDatabaseCatsSearchNameRequest) (*tablesv1.GetDatabaseCatsSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseCatsSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseCatsSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseCatsSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatsSearchNameResponse{
		Data: tablesconv.DatabaseCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateDatabaseCat вставляет строку dc.database_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateDatabaseCat(ctx context.Context, req *tablesv1.CreateDatabaseCatRequest) (*tablesv1.CreateDatabaseCatResponse, error) {
	if err := tablesvalidation.ValidateCreateDatabaseCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateDatabaseCat(ctx, tablesconv.ToCreateDatabaseCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateDatabaseCatResponse{DatabaseCat: tablesconv.DatabaseCatToProto(row)}, nil
}

// UpdateDatabaseCatById обновляет строку dc.database_cat.
func (t *TablesApiV1) UpdateDatabaseCatById(ctx context.Context, req *tablesv1.UpdateDatabaseCatByIdRequest) (*tablesv1.UpdateDatabaseCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDatabaseCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateDatabaseCatById(ctx, tablesconv.ToUpdateDatabaseCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDatabaseCatByIdResponse{}, nil
}

// DeleteDatabaseCatById мягко удаляет строку dc.database_cat.
func (t *TablesApiV1) DeleteDatabaseCatById(ctx context.Context, req *tablesv1.DeleteDatabaseCatByIdRequest) (*tablesv1.DeleteDatabaseCatByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteDatabaseCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDatabaseCatById(ctx, tablesconv.ToDeleteDatabaseCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDatabaseCatByIdResponse{}, nil
}

// UndeleteDatabaseCatById восстанавливает мягко удалённую строку dc.database_cat.
func (t *TablesApiV1) UndeleteDatabaseCatById(ctx context.Context, req *tablesv1.UndeleteDatabaseCatByIdRequest) (*tablesv1.UndeleteDatabaseCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteDatabaseCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDatabaseCatById(ctx, tablesconv.ToUndeleteDatabaseCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDatabaseCatByIdResponse{}, nil
}
