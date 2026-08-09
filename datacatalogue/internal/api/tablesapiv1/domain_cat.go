package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetDomainCatById отдаёт активную строку dc.domain_cat по id.
func (t *TablesApiV1) GetDomainCatById(ctx context.Context, req *tablesv1.GetDomainCatByIdRequest) (*tablesv1.GetDomainCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDomainCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDomainCatByIdResponse{DomainCat: tablesconv.DomainCatToProto(row)}, nil
}

// GetDomainCatsSearchName отдаёт страницу строк dc.domain_cat.
func (t *TablesApiV1) GetDomainCatsSearchName(ctx context.Context, req *tablesv1.GetDomainCatsSearchNameRequest) (*tablesv1.GetDomainCatsSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetDomainCatsSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDomainCatsSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetDomainCatsSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDomainCatsSearchNameResponse{
		Data: tablesconv.DomainCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDomainCats отдаёт страницу строк dc.domain_cat.
func (t *TablesApiV1) GetDomainCats(ctx context.Context, req *tablesv1.GetDomainCatsRequest) (*tablesv1.GetDomainCatsResponse, error) {
	if err := tablesvalidation.ValidateGetDomainCats(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDomainCatsParams(req)

	rows, count, err := t.services.TablesService.GetDomainCats(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDomainCatsResponse{
		Data: tablesconv.DomainCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetDomainCatsDeleted отдаёт страницу строк dc.domain_cat.
func (t *TablesApiV1) GetDomainCatsDeleted(ctx context.Context, req *tablesv1.GetDomainCatsDeletedRequest) (*tablesv1.GetDomainCatsDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetDomainCatsDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetDomainCatsDeletedParams(req)

	rows, count, err := t.services.TablesService.GetDomainCatsDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDomainCatsDeletedResponse{
		Data: tablesconv.DomainCatsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateDomainCat вставляет строку dc.domain_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateDomainCat(ctx context.Context, req *tablesv1.CreateDomainCatRequest) (*tablesv1.CreateDomainCatResponse, error) {
	if err := tablesvalidation.ValidateCreateDomainCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateDomainCat(ctx, tablesconv.ToCreateDomainCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateDomainCatResponse{DomainCat: tablesconv.DomainCatToProto(row)}, nil
}

// UpdateDomainCatById обновляет строку dc.domain_cat.
func (t *TablesApiV1) UpdateDomainCatById(ctx context.Context, req *tablesv1.UpdateDomainCatByIdRequest) (*tablesv1.UpdateDomainCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDomainCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateDomainCatById(ctx, tablesconv.ToUpdateDomainCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDomainCatByIdResponse{}, nil
}

// DeleteDomainCatById мягко удаляет строку dc.domain_cat.
func (t *TablesApiV1) DeleteDomainCatById(ctx context.Context, req *tablesv1.DeleteDomainCatByIdRequest) (*tablesv1.DeleteDomainCatByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteDomainCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDomainCatById(ctx, tablesconv.ToDeleteDomainCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDomainCatByIdResponse{}, nil
}

// UndeleteDomainCatById восстанавливает мягко удалённую строку dc.domain_cat.
func (t *TablesApiV1) UndeleteDomainCatById(ctx context.Context, req *tablesv1.UndeleteDomainCatByIdRequest) (*tablesv1.UndeleteDomainCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteDomainCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDomainCatById(ctx, tablesconv.ToUndeleteDomainCatByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDomainCatByIdResponse{}, nil
}
