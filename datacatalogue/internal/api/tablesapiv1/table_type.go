package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetTableTypeById отдаёт активную строку dc.table_type по id.
func (t *TablesApiV1) GetTableTypeById(ctx context.Context, req *tablesv1.GetTableTypeByIdRequest) (*tablesv1.GetTableTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetTableTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableTypeByIdResponse{TableType: tablesconv.TableTypeToProto(row)}, nil
}

// GetTableTypes отдаёт страницу строк dc.table_type.
func (t *TablesApiV1) GetTableTypes(ctx context.Context, req *tablesv1.GetTableTypesRequest) (*tablesv1.GetTableTypesResponse, error) {
	if err := tablesvalidation.ValidateGetTableTypes(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetTableTypesParams(req)

	rows, count, err := t.services.TablesService.GetTableTypes(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableTypesResponse{
		Data: tablesconv.TableTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetTableTypesDeleted отдаёт страницу строк dc.table_type.
func (t *TablesApiV1) GetTableTypesDeleted(ctx context.Context, req *tablesv1.GetTableTypesDeletedRequest) (*tablesv1.GetTableTypesDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetTableTypesDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetTableTypesDeletedParams(req)

	rows, count, err := t.services.TablesService.GetTableTypesDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableTypesDeletedResponse{
		Data: tablesconv.TableTypesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateTableType вставляет строку dc.table_type и отдаёт её целиком.
func (t *TablesApiV1) CreateTableType(ctx context.Context, req *tablesv1.CreateTableTypeRequest) (*tablesv1.CreateTableTypeResponse, error) {
	if err := tablesvalidation.ValidateCreateTableType(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateTableType(ctx, tablesconv.ToCreateTableTypeParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateTableTypeResponse{TableType: tablesconv.TableTypeToProto(row)}, nil
}

// UpdateTableTypeById обновляет строку dc.table_type.
func (t *TablesApiV1) UpdateTableTypeById(ctx context.Context, req *tablesv1.UpdateTableTypeByIdRequest) (*tablesv1.UpdateTableTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateTableTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateTableTypeById(ctx, tablesconv.ToUpdateTableTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateTableTypeByIdResponse{}, nil
}

// DeleteTableTypeById мягко удаляет строку dc.table_type.
func (t *TablesApiV1) DeleteTableTypeById(ctx context.Context, req *tablesv1.DeleteTableTypeByIdRequest) (*tablesv1.DeleteTableTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteTableTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteTableTypeById(ctx, tablesconv.ToDeleteTableTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteTableTypeByIdResponse{}, nil
}

// UndeleteTableTypeById восстанавливает мягко удалённую строку dc.table_type.
func (t *TablesApiV1) UndeleteTableTypeById(ctx context.Context, req *tablesv1.UndeleteTableTypeByIdRequest) (*tablesv1.UndeleteTableTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteTableTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteTableTypeById(ctx, tablesconv.ToUndeleteTableTypeByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteTableTypeByIdResponse{}, nil
}
