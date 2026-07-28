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

// GetTableCatById отдаёт активную строку dc.table_cat по id.
func (t *TablesApiV1) GetTableCatById(ctx context.Context, req *tablesv1.GetTableCatByIdRequest) (*tablesv1.GetTableCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetTableCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatByIdResponse{TableCat: tablesconv.TableCatToProto(row)}, nil
}

// GetTableCats отдаёт все активные строки dc.table_cat.
func (t *TablesApiV1) GetTableCats(ctx context.Context, req *tablesv1.GetTableCatsRequest) (*tablesv1.GetTableCatsResponse, error) {
	rows, err := t.services.TablesService.GetTableCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatsResponse{TableCats: tablesconv.TableCatsToProto(rows)}, nil
}

// GetDeletedTableCatById отдаёт мягко удалённую строку dc.table_cat по id.
func (t *TablesApiV1) GetDeletedTableCatById(ctx context.Context, req *tablesv1.GetDeletedTableCatByIdRequest) (*tablesv1.GetDeletedTableCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedTableCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedTableCatByIdResponse{TableCat: tablesconv.TableCatToProto(row)}, nil
}

// GetDeletedTableCats отдаёт все мягко удалённые строки dc.table_cat.
func (t *TablesApiV1) GetDeletedTableCats(ctx context.Context, req *tablesv1.GetDeletedTableCatsRequest) (*tablesv1.GetDeletedTableCatsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedTableCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedTableCatsResponse{TableCats: tablesconv.TableCatsToProto(rows)}, nil
}

// GetTableCatsBySchemaId отдаёт активные строки dc.table_cat, отобранные по schema_id.
func (t *TablesApiV1) GetTableCatsBySchemaId(ctx context.Context, req *tablesv1.GetTableCatsBySchemaIdRequest) (*tablesv1.GetTableCatsBySchemaIdResponse, error) {
	if err := validation.ValidateID(req.GetSchemaId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetTableCatsBySchemaId(ctx, req.GetSchemaId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatsBySchemaIdResponse{TableCats: tablesconv.TableCatsToProto(rows)}, nil
}

// GetTableCatsByTableTypeId отдаёт активные строки dc.table_cat, отобранные по table_type_id.
func (t *TablesApiV1) GetTableCatsByTableTypeId(ctx context.Context, req *tablesv1.GetTableCatsByTableTypeIdRequest) (*tablesv1.GetTableCatsByTableTypeIdResponse, error) {
	if err := validation.ValidateID(req.GetTableTypeId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetTableCatsByTableTypeId(ctx, req.GetTableTypeId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatsByTableTypeIdResponse{TableCats: tablesconv.TableCatsToProto(rows)}, nil
}

// GetTableCatsByDomainId отдаёт активные строки dc.table_cat, отобранные по domain_id.
func (t *TablesApiV1) GetTableCatsByDomainId(ctx context.Context, req *tablesv1.GetTableCatsByDomainIdRequest) (*tablesv1.GetTableCatsByDomainIdResponse, error) {
	if err := validation.ValidateID(req.GetDomainId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetTableCatsByDomainId(ctx, req.GetDomainId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetTableCatsByDomainIdResponse{TableCats: tablesconv.TableCatsToProto(rows)}, nil
}

// CreateTableCat вставляет строку dc.table_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateTableCat(ctx context.Context, req *tablesv1.CreateTableCatRequest) (*tablesv1.CreateTableCatResponse, error) {
	if err := tablesvalidation.ValidateCreateTableCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateTableCat(ctx, tablesconv.ToCreateTableCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateTableCatResponse{TableCat: tablesconv.TableCatToProto(row)}, nil
}

// UpdateTableCatById обновляет активную строку dc.table_cat и отдаёт её целиком.
func (t *TablesApiV1) UpdateTableCatById(ctx context.Context, req *tablesv1.UpdateTableCatByIdRequest) (*tablesv1.UpdateTableCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateTableCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateTableCatById(ctx, tablesconv.ToUpdateTableCatByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateTableCatByIdResponse{TableCat: tablesconv.TableCatToProto(row)}, nil
}

// DeleteTableCatById мягко удаляет строку dc.table_cat.
func (t *TablesApiV1) DeleteTableCatById(ctx context.Context, req *tablesv1.DeleteTableCatByIdRequest) (*tablesv1.DeleteTableCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteTableCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteTableCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteTableCatById восстанавливает мягко удалённую строку dc.table_cat.
func (t *TablesApiV1) UndeleteTableCatById(ctx context.Context, req *tablesv1.UndeleteTableCatByIdRequest) (*tablesv1.UndeleteTableCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteTableCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteTableCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}
