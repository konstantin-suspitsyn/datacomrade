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

// GetSchemaCats отдаёт все активные строки dc.schema_cat.
func (t *TablesApiV1) GetSchemaCats(ctx context.Context, req *tablesv1.GetSchemaCatsRequest) (*tablesv1.GetSchemaCatsResponse, error) {
	rows, err := t.services.TablesService.GetSchemaCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetSchemaCatsResponse{SchemaCats: tablesconv.SchemaCatsToProto(rows)}, nil
}

// GetDeletedSchemaCatById отдаёт мягко удалённую строку dc.schema_cat по id.
func (t *TablesApiV1) GetDeletedSchemaCatById(ctx context.Context, req *tablesv1.GetDeletedSchemaCatByIdRequest) (*tablesv1.GetDeletedSchemaCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedSchemaCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedSchemaCatByIdResponse{SchemaCat: tablesconv.SchemaCatToProto(row)}, nil
}

// GetDeletedSchemaCats отдаёт все мягко удалённые строки dc.schema_cat.
func (t *TablesApiV1) GetDeletedSchemaCats(ctx context.Context, req *tablesv1.GetDeletedSchemaCatsRequest) (*tablesv1.GetDeletedSchemaCatsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedSchemaCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedSchemaCatsResponse{SchemaCats: tablesconv.SchemaCatsToProto(rows)}, nil
}

// GetSchemaCatsByDatabaseId отдаёт активные строки dc.schema_cat, отобранные по database_id.
func (t *TablesApiV1) GetSchemaCatsByDatabaseId(ctx context.Context, req *tablesv1.GetSchemaCatsByDatabaseIdRequest) (*tablesv1.GetSchemaCatsByDatabaseIdResponse, error) {
	if err := validation.ValidateID(req.GetDatabaseId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetSchemaCatsByDatabaseId(ctx, req.GetDatabaseId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetSchemaCatsByDatabaseIdResponse{SchemaCats: tablesconv.SchemaCatsToProto(rows)}, nil
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

// UpdateSchemaCatById обновляет активную строку dc.schema_cat и отдаёт её целиком.
func (t *TablesApiV1) UpdateSchemaCatById(ctx context.Context, req *tablesv1.UpdateSchemaCatByIdRequest) (*tablesv1.UpdateSchemaCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateSchemaCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateSchemaCatById(ctx, tablesconv.ToUpdateSchemaCatByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateSchemaCatByIdResponse{SchemaCat: tablesconv.SchemaCatToProto(row)}, nil
}

// DeleteSchemaCatById мягко удаляет строку dc.schema_cat.
func (t *TablesApiV1) DeleteSchemaCatById(ctx context.Context, req *tablesv1.DeleteSchemaCatByIdRequest) (*tablesv1.DeleteSchemaCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteSchemaCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteSchemaCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteSchemaCatById восстанавливает мягко удалённую строку dc.schema_cat.
func (t *TablesApiV1) UndeleteSchemaCatById(ctx context.Context, req *tablesv1.UndeleteSchemaCatByIdRequest) (*tablesv1.UndeleteSchemaCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteSchemaCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteSchemaCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}
