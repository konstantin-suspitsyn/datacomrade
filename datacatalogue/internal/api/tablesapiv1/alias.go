package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetAliasById отдаёт активную строку dc.alias по id.
func (t *TablesApiV1) GetAliasById(ctx context.Context, req *tablesv1.GetAliasByIdRequest) (*tablesv1.GetAliasByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetAliasById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetAliasByIdResponse{Alias: tablesconv.AliasToProto(row)}, nil
}

// GetAliasesDeleted отдаёт страницу строк dc.alias.
func (t *TablesApiV1) GetAliasesDeleted(ctx context.Context, req *tablesv1.GetAliasesDeletedRequest) (*tablesv1.GetAliasesDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetAliasesDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetAliasesDeletedParams(req)

	rows, count, err := t.services.TablesService.GetAliasesDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetAliasesDeletedResponse{
		Data: tablesconv.AliasesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetAliases отдаёт страницу строк dc.alias.
func (t *TablesApiV1) GetAliases(ctx context.Context, req *tablesv1.GetAliasesRequest) (*tablesv1.GetAliasesResponse, error) {
	if err := tablesvalidation.ValidateGetAliases(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetAliasesParams(req)

	rows, count, err := t.services.TablesService.GetAliases(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetAliasesResponse{
		Data: tablesconv.AliasesToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateAlias вставляет строку dc.alias и отдаёт её целиком.
func (t *TablesApiV1) CreateAlias(ctx context.Context, req *tablesv1.CreateAliasRequest) (*tablesv1.CreateAliasResponse, error) {
	if err := tablesvalidation.ValidateCreateAlias(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateAlias(ctx, tablesconv.ToCreateAliasParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateAliasResponse{Alias: tablesconv.AliasToProto(row)}, nil
}

// UpdateAliasById обновляет строку dc.alias.
func (t *TablesApiV1) UpdateAliasById(ctx context.Context, req *tablesv1.UpdateAliasByIdRequest) (*tablesv1.UpdateAliasByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateAliasById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateAliasById(ctx, tablesconv.ToUpdateAliasByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateAliasByIdResponse{}, nil
}

// DeleteAliasById мягко удаляет строку dc.alias.
func (t *TablesApiV1) DeleteAliasById(ctx context.Context, req *tablesv1.DeleteAliasByIdRequest) (*tablesv1.DeleteAliasByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteAliasById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteAliasById(ctx, tablesconv.ToDeleteAliasByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteAliasByIdResponse{}, nil
}

// UndeleteAliasById восстанавливает мягко удалённую строку dc.alias.
func (t *TablesApiV1) UndeleteAliasById(ctx context.Context, req *tablesv1.UndeleteAliasByIdRequest) (*tablesv1.UndeleteAliasByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteAliasById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteAliasById(ctx, tablesconv.ToUndeleteAliasByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteAliasByIdResponse{}, nil
}
