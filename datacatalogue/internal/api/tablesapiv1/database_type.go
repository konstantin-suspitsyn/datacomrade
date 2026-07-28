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

// GetDatabaseTypes отдаёт все активные строки dc.database_type.
func (t *TablesApiV1) GetDatabaseTypes(ctx context.Context, req *tablesv1.GetDatabaseTypesRequest) (*tablesv1.GetDatabaseTypesResponse, error) {
	rows, err := t.services.TablesService.GetDatabaseTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseTypesResponse{DatabaseTypes: tablesconv.DatabaseTypesToProto(rows)}, nil
}

// GetDeletedDatabaseTypeById отдаёт мягко удалённую строку dc.database_type по id.
func (t *TablesApiV1) GetDeletedDatabaseTypeById(ctx context.Context, req *tablesv1.GetDeletedDatabaseTypeByIdRequest) (*tablesv1.GetDeletedDatabaseTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedDatabaseTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDatabaseTypeByIdResponse{DatabaseType: tablesconv.DatabaseTypeToProto(row)}, nil
}

// GetDeletedDatabaseTypes отдаёт все мягко удалённые строки dc.database_type.
func (t *TablesApiV1) GetDeletedDatabaseTypes(ctx context.Context, req *tablesv1.GetDeletedDatabaseTypesRequest) (*tablesv1.GetDeletedDatabaseTypesResponse, error) {
	rows, err := t.services.TablesService.GetDeletedDatabaseTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDatabaseTypesResponse{DatabaseTypes: tablesconv.DatabaseTypesToProto(rows)}, nil
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

// UpdateDatabaseTypeById обновляет активную строку dc.database_type и отдаёт её целиком.
func (t *TablesApiV1) UpdateDatabaseTypeById(ctx context.Context, req *tablesv1.UpdateDatabaseTypeByIdRequest) (*tablesv1.UpdateDatabaseTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDatabaseTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateDatabaseTypeById(ctx, tablesconv.ToUpdateDatabaseTypeByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDatabaseTypeByIdResponse{DatabaseType: tablesconv.DatabaseTypeToProto(row)}, nil
}

// DeleteDatabaseTypeById мягко удаляет строку dc.database_type.
func (t *TablesApiV1) DeleteDatabaseTypeById(ctx context.Context, req *tablesv1.DeleteDatabaseTypeByIdRequest) (*tablesv1.DeleteDatabaseTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDatabaseTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDatabaseTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteDatabaseTypeById восстанавливает мягко удалённую строку dc.database_type.
func (t *TablesApiV1) UndeleteDatabaseTypeById(ctx context.Context, req *tablesv1.UndeleteDatabaseTypeByIdRequest) (*tablesv1.UndeleteDatabaseTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDatabaseTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDatabaseTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}
