package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"google.golang.org/protobuf/types/known/emptypb"
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

// GetHosts отдаёт все активные строки dc.host.
func (t *TablesApiV1) GetHosts(ctx context.Context, req *tablesv1.GetHostsRequest) (*tablesv1.GetHostsResponse, error) {
	rows, err := t.services.TablesService.GetHosts(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetHostsResponse{Hosts: tablesconv.HostsToProto(rows)}, nil
}

// GetDeletedHostById отдаёт мягко удалённую строку dc.host по id.
func (t *TablesApiV1) GetDeletedHostById(ctx context.Context, req *tablesv1.GetDeletedHostByIdRequest) (*tablesv1.GetDeletedHostByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedHostById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedHostByIdResponse{Host: tablesconv.HostToProto(row)}, nil
}

// GetDeletedHosts отдаёт все мягко удалённые строки dc.host.
func (t *TablesApiV1) GetDeletedHosts(ctx context.Context, req *tablesv1.GetDeletedHostsRequest) (*tablesv1.GetDeletedHostsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedHosts(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedHostsResponse{Hosts: tablesconv.HostsToProto(rows)}, nil
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

// UpdateHostById обновляет активную строку dc.host и отдаёт её целиком.
func (t *TablesApiV1) UpdateHostById(ctx context.Context, req *tablesv1.UpdateHostByIdRequest) (*tablesv1.UpdateHostByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateHostById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateHostById(ctx, tablesconv.ToUpdateHostByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateHostByIdResponse{Host: tablesconv.HostToProto(row)}, nil
}

// DeleteHostById мягко удаляет строку dc.host.
func (t *TablesApiV1) DeleteHostById(ctx context.Context, req *tablesv1.DeleteHostByIdRequest) (*tablesv1.DeleteHostByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteHostById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteHostByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteHostById восстанавливает мягко удалённую строку dc.host.
func (t *TablesApiV1) UndeleteHostById(ctx context.Context, req *tablesv1.UndeleteHostByIdRequest) (*tablesv1.UndeleteHostByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteHostById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteHostByIdResponse{Empty: &emptypb.Empty{}}, nil
}
