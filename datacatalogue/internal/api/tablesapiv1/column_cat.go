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

// GetColumnCatById отдаёт активную строку dc.column_cat по id.
func (t *TablesApiV1) GetColumnCatById(ctx context.Context, req *tablesv1.GetColumnCatByIdRequest) (*tablesv1.GetColumnCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetColumnCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatByIdResponse{ColumnCat: tablesconv.ColumnCatToProto(row)}, nil
}

// GetColumnCats отдаёт все активные строки dc.column_cat.
func (t *TablesApiV1) GetColumnCats(ctx context.Context, req *tablesv1.GetColumnCatsRequest) (*tablesv1.GetColumnCatsResponse, error) {
	rows, err := t.services.TablesService.GetColumnCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatsResponse{ColumnCats: tablesconv.ColumnCatsToProto(rows)}, nil
}

// GetDeletedColumnCatById отдаёт мягко удалённую строку dc.column_cat по id.
func (t *TablesApiV1) GetDeletedColumnCatById(ctx context.Context, req *tablesv1.GetDeletedColumnCatByIdRequest) (*tablesv1.GetDeletedColumnCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedColumnCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedColumnCatByIdResponse{ColumnCat: tablesconv.ColumnCatToProto(row)}, nil
}

// GetDeletedColumnCats отдаёт все мягко удалённые строки dc.column_cat.
func (t *TablesApiV1) GetDeletedColumnCats(ctx context.Context, req *tablesv1.GetDeletedColumnCatsRequest) (*tablesv1.GetDeletedColumnCatsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedColumnCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedColumnCatsResponse{ColumnCats: tablesconv.ColumnCatsToProto(rows)}, nil
}

// GetColumnCatsByTableId отдаёт активные строки dc.column_cat, отобранные по table_id.
func (t *TablesApiV1) GetColumnCatsByTableId(ctx context.Context, req *tablesv1.GetColumnCatsByTableIdRequest) (*tablesv1.GetColumnCatsByTableIdResponse, error) {
	if err := validation.ValidateID(req.GetTableId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetColumnCatsByTableId(ctx, req.GetTableId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatsByTableIdResponse{ColumnCats: tablesconv.ColumnCatsToProto(rows)}, nil
}

// GetColumnCatsByAliasId отдаёт активные строки dc.column_cat, отобранные по alias_id.
func (t *TablesApiV1) GetColumnCatsByAliasId(ctx context.Context, req *tablesv1.GetColumnCatsByAliasIdRequest) (*tablesv1.GetColumnCatsByAliasIdResponse, error) {
	if err := validation.ValidateID(req.GetAliasId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetColumnCatsByAliasId(ctx, req.GetAliasId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetColumnCatsByAliasIdResponse{ColumnCats: tablesconv.ColumnCatsToProto(rows)}, nil
}

// CreateColumnCat вставляет строку dc.column_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateColumnCat(ctx context.Context, req *tablesv1.CreateColumnCatRequest) (*tablesv1.CreateColumnCatResponse, error) {
	if err := tablesvalidation.ValidateCreateColumnCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateColumnCat(ctx, tablesconv.ToCreateColumnCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateColumnCatResponse{ColumnCat: tablesconv.ColumnCatToProto(row)}, nil
}

// UpdateColumnCatById обновляет активную строку dc.column_cat и отдаёт её целиком.
func (t *TablesApiV1) UpdateColumnCatById(ctx context.Context, req *tablesv1.UpdateColumnCatByIdRequest) (*tablesv1.UpdateColumnCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateColumnCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateColumnCatById(ctx, tablesconv.ToUpdateColumnCatByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateColumnCatByIdResponse{ColumnCat: tablesconv.ColumnCatToProto(row)}, nil
}

// DeleteColumnCatById мягко удаляет строку dc.column_cat.
func (t *TablesApiV1) DeleteColumnCatById(ctx context.Context, req *tablesv1.DeleteColumnCatByIdRequest) (*tablesv1.DeleteColumnCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteColumnCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteColumnCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteColumnCatById восстанавливает мягко удалённую строку dc.column_cat.
func (t *TablesApiV1) UndeleteColumnCatById(ctx context.Context, req *tablesv1.UndeleteColumnCatByIdRequest) (*tablesv1.UndeleteColumnCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteColumnCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteColumnCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}
