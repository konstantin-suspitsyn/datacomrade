package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetSchemaCatById отдаёт активную строку dc.schema_cat по id.
func (t *TablesApiV1) GetSchemaCatById(ctx context.Context, req *tablesv1.GetSchemaCatByIdRequest) (*tablesv1.GetSchemaCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetSchemaCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetSchemaCatByIdResponse{SchemaCat: tablesconv.SchemaCatToProto(row)}, nil
}

// GetSchemaCats отдаёт страницу строк dc.schema_cat.
func (t *TablesApiV1) GetSchemaCats(ctx context.Context, req *tablesv1.GetSchemaCatsRequest) (*tablesv1.GetSchemaCatsResponse, error) {
	if err := tablesvalidation.ValidateGetSchemaCats(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetSchemaCatsParams(req)

	rows, count, err := t.services.TablesService.GetSchemaCats(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetSchemaCatsResponse{
		Data: tablesconv.SchemaCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetSchemaCatsDeleted отдаёт страницу строк dc.schema_cat.
func (t *TablesApiV1) GetSchemaCatsDeleted(ctx context.Context, req *tablesv1.GetSchemaCatsDeletedRequest) (*tablesv1.GetSchemaCatsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetSchemaCatsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetSchemaCatsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetSchemaCatsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetSchemaCatsDeletedResponse{
		Data: tablesconv.SchemaCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateSchemaCat вставляет строку dc.schema_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateSchemaCat(ctx context.Context, req *tablesv1.CreateSchemaCatRequest) (*tablesv1.CreateSchemaCatResponse, error) {
	if err := tablesvalidation.ValidateCreateSchemaCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateSchemaCat(ctx, tablesconv.ToCreateSchemaCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateSchemaCatResponse{SchemaCat: tablesconv.SchemaCatToProto(row)}, nil
}

// UpdateSchemaCatById обновляет строку dc.schema_cat.
func (t *TablesApiV1) UpdateSchemaCatById(ctx context.Context, req *tablesv1.UpdateSchemaCatByIdRequest) (*tablesv1.UpdateSchemaCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateSchemaCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateSchemaCatById(ctx, tablesconv.ToUpdateSchemaCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateSchemaCatByIdResponse{}, nil
}

// DeleteSchemaCatById мягко удаляет строку dc.schema_cat.
func (t *TablesApiV1) DeleteSchemaCatById(ctx context.Context, req *tablesv1.DeleteSchemaCatByIdRequest) (*tablesv1.DeleteSchemaCatByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteSchemaCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteSchemaCatById(ctx, tablesconv.ToDeleteSchemaCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteSchemaCatByIdResponse{}, nil
}

// UndeleteSchemaCatById восстанавливает мягко удалённую строку dc.schema_cat.
func (t *TablesApiV1) UndeleteSchemaCatById(ctx context.Context, req *tablesv1.UndeleteSchemaCatByIdRequest) (*tablesv1.UndeleteSchemaCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteSchemaCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteSchemaCatById(ctx, tablesconv.ToUndeleteSchemaCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteSchemaCatByIdResponse{}, nil
}
