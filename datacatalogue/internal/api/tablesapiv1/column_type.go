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

// GetColumnTypeById отдаёт активную строку dc.column_type по id.
func (t *TablesApiV1) GetColumnTypeById(ctx context.Context, req *tablesv1.GetColumnTypeByIdRequest) (*tablesv1.GetColumnTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetColumnTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnTypeByIdResponse{ColumnType: tablesconv.ColumnTypeToProto(row)}, nil
}

// GetColumnTypes отдаёт все активные строки dc.column_type.
func (t *TablesApiV1) GetColumnTypes(ctx context.Context, req *tablesv1.GetColumnTypesRequest) (*tablesv1.GetColumnTypesResponse, error) {
	rows, err := t.services.TablesService.GetColumnTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnTypesResponse{ColumnTypes: tablesconv.ColumnTypesToProto(rows)}, nil
}

// GetDeletedColumnTypeById отдаёт мягко удалённую строку dc.column_type по id.
func (t *TablesApiV1) GetDeletedColumnTypeById(ctx context.Context, req *tablesv1.GetDeletedColumnTypeByIdRequest) (*tablesv1.GetDeletedColumnTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedColumnTypeById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedColumnTypeByIdResponse{ColumnType: tablesconv.ColumnTypeToProto(row)}, nil
}

// GetDeletedColumnTypes отдаёт все мягко удалённые строки dc.column_type.
func (t *TablesApiV1) GetDeletedColumnTypes(ctx context.Context, req *tablesv1.GetDeletedColumnTypesRequest) (*tablesv1.GetDeletedColumnTypesResponse, error) {
	rows, err := t.services.TablesService.GetDeletedColumnTypes(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedColumnTypesResponse{ColumnTypes: tablesconv.ColumnTypesToProto(rows)}, nil
}

// CreateColumnType вставляет строку dc.column_type и отдаёт её целиком.
func (t *TablesApiV1) CreateColumnType(ctx context.Context, req *tablesv1.CreateColumnTypeRequest) (*tablesv1.CreateColumnTypeResponse, error) {
	if err := tablesvalidation.ValidateCreateColumnType(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateColumnType(ctx, tablesconv.ToCreateColumnTypeParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateColumnTypeResponse{ColumnType: tablesconv.ColumnTypeToProto(row)}, nil
}

// UpdateColumnTypeById обновляет активную строку dc.column_type и отдаёт её целиком.
func (t *TablesApiV1) UpdateColumnTypeById(ctx context.Context, req *tablesv1.UpdateColumnTypeByIdRequest) (*tablesv1.UpdateColumnTypeByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateColumnTypeById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateColumnTypeById(ctx, tablesconv.ToUpdateColumnTypeByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateColumnTypeByIdResponse{ColumnType: tablesconv.ColumnTypeToProto(row)}, nil
}

// DeleteColumnTypeById мягко удаляет строку dc.column_type.
func (t *TablesApiV1) DeleteColumnTypeById(ctx context.Context, req *tablesv1.DeleteColumnTypeByIdRequest) (*tablesv1.DeleteColumnTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteColumnTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteColumnTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteColumnTypeById восстанавливает мягко удалённую строку dc.column_type.
func (t *TablesApiV1) UndeleteColumnTypeById(ctx context.Context, req *tablesv1.UndeleteColumnTypeByIdRequest) (*tablesv1.UndeleteColumnTypeByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteColumnTypeById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteColumnTypeByIdResponse{Empty: &emptypb.Empty{}}, nil
}
