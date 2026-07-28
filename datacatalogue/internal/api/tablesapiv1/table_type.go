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

// GetTableTypes отдаёт все активные строки dc.table_type.
func (t *TablesApiV1) GetTableTypes(ctx context.Context, req *tablesv1.GetTableTypesRequest) (*tablesv1.GetTableTypesResponse, error) {
	rows, err := t.services.TablesService.GetTableTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableTypesResponse{TableTypes: tablesconv.TableTypesToProto(rows)}, nil
}

// GetDeletedTableTypeById отдаёт мягко удалённую строку dc.table_type по id.
func (t *TablesApiV1) GetDeletedTableTypeById(ctx context.Context, req *tablesv1.GetDeletedTableTypeByIdRequest) (*tablesv1.GetDeletedTableTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedTableTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedTableTypeByIdResponse{TableType: tablesconv.TableTypeToProto(row)}, nil
}

// GetDeletedTableTypes отдаёт все мягко удалённые строки dc.table_type.
func (t *TablesApiV1) GetDeletedTableTypes(ctx context.Context, req *tablesv1.GetDeletedTableTypesRequest) (*tablesv1.GetDeletedTableTypesResponse, error) {
	rows, err := t.services.TablesService.GetDeletedTableTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedTableTypesResponse{TableTypes: tablesconv.TableTypesToProto(rows)}, nil
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

// UpdateTableTypeById обновляет активную строку dc.table_type и отдаёт её целиком.
func (t *TablesApiV1) UpdateTableTypeById(ctx context.Context, req *tablesv1.UpdateTableTypeByIdRequest) (*tablesv1.UpdateTableTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateTableTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateTableTypeById(ctx, tablesconv.ToUpdateTableTypeByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateTableTypeByIdResponse{TableType: tablesconv.TableTypeToProto(row)}, nil
}

// DeleteTableTypeById мягко удаляет строку dc.table_type.
func (t *TablesApiV1) DeleteTableTypeById(ctx context.Context, req *tablesv1.DeleteTableTypeByIdRequest) (*tablesv1.DeleteTableTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteTableTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteTableTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteTableTypeById восстанавливает мягко удалённую строку dc.table_type.
func (t *TablesApiV1) UndeleteTableTypeById(ctx context.Context, req *tablesv1.UndeleteTableTypeByIdRequest) (*tablesv1.UndeleteTableTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteTableTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteTableTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}
