package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetDatabaseTypeById отдаёт активную строку dc.database_type по id.
func (t *TablesApiV1) GetDatabaseTypeById(ctx context.Context, req *tablesv1.GetDatabaseTypeByIdRequest) (*tablesv1.GetDatabaseTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDatabaseTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseTypeByIdResponse{DatabaseType: tablesconv.DatabaseTypeToProto(row)}, nil
}

// GetDatabaseTypesSearchName отдаёт страницу строк dc.database_type.
func (t *TablesApiV1) GetDatabaseTypesSearchName(ctx context.Context, req *tablesv1.GetDatabaseTypesSearchNameRequest) (*tablesv1.GetDatabaseTypesSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseTypesSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseTypesSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseTypesSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseTypesSearchNameResponse{
		Data: tablesconv.DatabaseTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseTypes отдаёт страницу строк dc.database_type.
func (t *TablesApiV1) GetDatabaseTypes(ctx context.Context, req *tablesv1.GetDatabaseTypesRequest) (*tablesv1.GetDatabaseTypesResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseTypes(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseTypesParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseTypes(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseTypesResponse{
		Data: tablesconv.DatabaseTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDatabaseTypesDeleted отдаёт страницу строк dc.database_type.
func (t *TablesApiV1) GetDatabaseTypesDeleted(ctx context.Context, req *tablesv1.GetDatabaseTypesDeletedRequest) (*tablesv1.GetDatabaseTypesDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetDatabaseTypesDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDatabaseTypesDeletedParams(req)

	rows, count, err := t.services.TablesService.GetDatabaseTypesDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseTypesDeletedResponse{
		Data: tablesconv.DatabaseTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateDatabaseType вставляет строку dc.database_type и отдаёт её целиком.
func (t *TablesApiV1) CreateDatabaseType(ctx context.Context, req *tablesv1.CreateDatabaseTypeRequest) (*tablesv1.CreateDatabaseTypeResponse, error) {
	if err := tablesvalidation.ValidateCreateDatabaseType(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateDatabaseType(ctx, tablesconv.ToCreateDatabaseTypeParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateDatabaseTypeResponse{DatabaseType: tablesconv.DatabaseTypeToProto(row)}, nil
}

// UpdateDatabaseTypeById обновляет строку dc.database_type.
func (t *TablesApiV1) UpdateDatabaseTypeById(ctx context.Context, req *tablesv1.UpdateDatabaseTypeByIdRequest) (*tablesv1.UpdateDatabaseTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDatabaseTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateDatabaseTypeById(ctx, tablesconv.ToUpdateDatabaseTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDatabaseTypeByIdResponse{}, nil
}

// DeleteDatabaseTypeById мягко удаляет строку dc.database_type.
func (t *TablesApiV1) DeleteDatabaseTypeById(ctx context.Context, req *tablesv1.DeleteDatabaseTypeByIdRequest) (*tablesv1.DeleteDatabaseTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteDatabaseTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDatabaseTypeById(ctx, tablesconv.ToDeleteDatabaseTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDatabaseTypeByIdResponse{}, nil
}

// UndeleteDatabaseTypeById восстанавливает мягко удалённую строку dc.database_type.
func (t *TablesApiV1) UndeleteDatabaseTypeById(ctx context.Context, req *tablesv1.UndeleteDatabaseTypeByIdRequest) (*tablesv1.UndeleteDatabaseTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteDatabaseTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDatabaseTypeById(ctx, tablesconv.ToUndeleteDatabaseTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDatabaseTypeByIdResponse{}, nil
}
