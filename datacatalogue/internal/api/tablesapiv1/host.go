package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetHostById отдаёт активную строку dc.host по id.
func (t *TablesApiV1) GetHostById(ctx context.Context, req *tablesv1.GetHostByIdRequest) (*tablesv1.GetHostByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetHostById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHostByIdResponse{Host: tablesconv.HostToProto(row)}, nil
}

// GetHosts отдаёт страницу строк dc.host.
func (t *TablesApiV1) GetHosts(ctx context.Context, req *tablesv1.GetHostsRequest) (*tablesv1.GetHostsResponse, error) {
	if err := tablesvalidation.ValidateGetHosts(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetHostsParams(req)

	rows, count, err := t.services.TablesService.GetHosts(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHostsResponse{
		Data: tablesconv.HostsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetHostsSearchName отдаёт страницу строк dc.host.
func (t *TablesApiV1) GetHostsSearchName(ctx context.Context, req *tablesv1.GetHostsSearchNameRequest) (*tablesv1.GetHostsSearchNameResponse, error) {
	if err := tablesvalidation.ValidateGetHostsSearchName(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetHostsSearchNameParams(req)

	rows, count, err := t.services.TablesService.GetHostsSearchName(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHostsSearchNameResponse{
		Data: tablesconv.HostsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetHostDeleted отдаёт страницу строк dc.host.
func (t *TablesApiV1) GetHostDeleted(ctx context.Context, req *tablesv1.GetHostDeletedRequest) (*tablesv1.GetHostDeletedResponse, error) {
	if err := tablesvalidation.ValidateGetHostDeleted(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetHostDeletedParams(req)

	rows, count, err := t.services.TablesService.GetHostDeleted(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHostDeletedResponse{
		Data: tablesconv.HostsToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateHost вставляет строку dc.host и отдаёт её целиком.
func (t *TablesApiV1) CreateHost(ctx context.Context, req *tablesv1.CreateHostRequest) (*tablesv1.CreateHostResponse, error) {
	if err := tablesvalidation.ValidateCreateHost(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateHost(ctx, tablesconv.ToCreateHostParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateHostResponse{Host: tablesconv.HostToProto(row)}, nil
}

// UpdateHostById обновляет строку dc.host.
func (t *TablesApiV1) UpdateHostById(ctx context.Context, req *tablesv1.UpdateHostByIdRequest) (*tablesv1.UpdateHostByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateHostById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UpdateHostById(ctx, tablesconv.ToUpdateHostByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateHostByIdResponse{}, nil
}

// DeleteHostById мягко удаляет строку dc.host.
func (t *TablesApiV1) DeleteHostById(ctx context.Context, req *tablesv1.DeleteHostByIdRequest) (*tablesv1.DeleteHostByIdResponse, error) {
	if err := tablesvalidation.ValidateDeleteHostById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteHostById(ctx, tablesconv.ToDeleteHostByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteHostByIdResponse{}, nil
}

// UndeleteHostById восстанавливает мягко удалённую строку dc.host.
func (t *TablesApiV1) UndeleteHostById(ctx context.Context, req *tablesv1.UndeleteHostByIdRequest) (*tablesv1.UndeleteHostByIdResponse, error) {
	if err := tablesvalidation.ValidateUndeleteHostById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteHostById(ctx, tablesconv.ToUndeleteHostByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteHostByIdResponse{}, nil
}
