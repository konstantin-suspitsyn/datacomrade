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

// GetAliases отдаёт все активные строки dc.alias.
func (t *TablesApiV1) GetAliases(ctx context.Context, req *tablesv1.GetAliasesRequest) (*tablesv1.GetAliasesResponse, error) {
	rows, err := t.services.TablesService.GetAliases(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetAliasesResponse{Aliases: tablesconv.AliasesToProto(rows)}, nil
}

// GetDeletedAliasById отдаёт мягко удалённую строку dc.alias по id.
func (t *TablesApiV1) GetDeletedAliasById(ctx context.Context, req *tablesv1.GetDeletedAliasByIdRequest) (*tablesv1.GetDeletedAliasByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedAliasById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedAliasByIdResponse{Alias: tablesconv.AliasToProto(row)}, nil
}

// GetDeletedAliases отдаёт все мягко удалённые строки dc.alias.
func (t *TablesApiV1) GetDeletedAliases(ctx context.Context, req *tablesv1.GetDeletedAliasesRequest) (*tablesv1.GetDeletedAliasesResponse, error) {
	rows, err := t.services.TablesService.GetDeletedAliases(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedAliasesResponse{Aliases: tablesconv.AliasesToProto(rows)}, nil
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

// UpdateAliasById обновляет активную строку dc.alias и отдаёт её целиком.
func (t *TablesApiV1) UpdateAliasById(ctx context.Context, req *tablesv1.UpdateAliasByIdRequest) (*tablesv1.UpdateAliasByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateAliasById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateAliasById(ctx, tablesconv.ToUpdateAliasByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateAliasByIdResponse{Alias: tablesconv.AliasToProto(row)}, nil
}

// DeleteAliasById мягко удаляет строку dc.alias.
func (t *TablesApiV1) DeleteAliasById(ctx context.Context, req *tablesv1.DeleteAliasByIdRequest) (*tablesv1.DeleteAliasByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteAliasById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteAliasByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteAliasById восстанавливает мягко удалённую строку dc.alias.
func (t *TablesApiV1) UndeleteAliasById(ctx context.Context, req *tablesv1.UndeleteAliasByIdRequest) (*tablesv1.UndeleteAliasByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteAliasById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteAliasByIdResponse{Empty: &emptypb.Empty{}}, nil
}
