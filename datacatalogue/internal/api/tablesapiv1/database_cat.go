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

// GetDatabaseCatById отдаёт активную строку dc.database_cat по id.
func (t *TablesApiV1) GetDatabaseCatById(ctx context.Context, req *tablesv1.GetDatabaseCatByIdRequest) (*tablesv1.GetDatabaseCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDatabaseCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatByIdResponse{DatabaseCat: tablesconv.DatabaseCatToProto(row)}, nil
}

// GetDatabaseCats отдаёт все активные строки dc.database_cat.
func (t *TablesApiV1) GetDatabaseCats(ctx context.Context, req *tablesv1.GetDatabaseCatsRequest) (*tablesv1.GetDatabaseCatsResponse, error) {
	rows, err := t.services.TablesService.GetDatabaseCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatsResponse{DatabaseCats: tablesconv.DatabaseCatsToProto(rows)}, nil
}

// GetDeletedDatabaseCatById отдаёт мягко удалённую строку dc.database_cat по id.
func (t *TablesApiV1) GetDeletedDatabaseCatById(ctx context.Context, req *tablesv1.GetDeletedDatabaseCatByIdRequest) (*tablesv1.GetDeletedDatabaseCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedDatabaseCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDatabaseCatByIdResponse{DatabaseCat: tablesconv.DatabaseCatToProto(row)}, nil
}

// GetDeletedDatabaseCats отдаёт все мягко удалённые строки dc.database_cat.
func (t *TablesApiV1) GetDeletedDatabaseCats(ctx context.Context, req *tablesv1.GetDeletedDatabaseCatsRequest) (*tablesv1.GetDeletedDatabaseCatsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedDatabaseCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDatabaseCatsResponse{DatabaseCats: tablesconv.DatabaseCatsToProto(rows)}, nil
}

// GetDatabaseCatsByHostId отдаёт активные строки dc.database_cat, отобранные по host_id.
func (t *TablesApiV1) GetDatabaseCatsByHostId(ctx context.Context, req *tablesv1.GetDatabaseCatsByHostIdRequest) (*tablesv1.GetDatabaseCatsByHostIdResponse, error) {
	if err := validation.ValidateID(req.GetHostId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetDatabaseCatsByHostId(ctx, req.GetHostId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatsByHostIdResponse{DatabaseCats: tablesconv.DatabaseCatsToProto(rows)}, nil
}

// GetDatabaseCatsByDatabaseTypeId отдаёт активные строки dc.database_cat, отобранные по database_type_id.
func (t *TablesApiV1) GetDatabaseCatsByDatabaseTypeId(ctx context.Context, req *tablesv1.GetDatabaseCatsByDatabaseTypeIdRequest) (*tablesv1.GetDatabaseCatsByDatabaseTypeIdResponse, error) {
	if err := validation.ValidateID(req.GetDatabaseTypeId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	rows, err := t.services.TablesService.GetDatabaseCatsByDatabaseTypeId(ctx, req.GetDatabaseTypeId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDatabaseCatsByDatabaseTypeIdResponse{DatabaseCats: tablesconv.DatabaseCatsToProto(rows)}, nil
}

// CreateDatabaseCat вставляет строку dc.database_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateDatabaseCat(ctx context.Context, req *tablesv1.CreateDatabaseCatRequest) (*tablesv1.CreateDatabaseCatResponse, error) {
	if err := tablesvalidation.ValidateCreateDatabaseCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateDatabaseCat(ctx, tablesconv.ToCreateDatabaseCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateDatabaseCatResponse{DatabaseCat: tablesconv.DatabaseCatToProto(row)}, nil
}

// UpdateDatabaseCatById обновляет активную строку dc.database_cat и отдаёт её целиком.
func (t *TablesApiV1) UpdateDatabaseCatById(ctx context.Context, req *tablesv1.UpdateDatabaseCatByIdRequest) (*tablesv1.UpdateDatabaseCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDatabaseCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateDatabaseCatById(ctx, tablesconv.ToUpdateDatabaseCatByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDatabaseCatByIdResponse{DatabaseCat: tablesconv.DatabaseCatToProto(row)}, nil
}

// DeleteDatabaseCatById мягко удаляет строку dc.database_cat.
func (t *TablesApiV1) DeleteDatabaseCatById(ctx context.Context, req *tablesv1.DeleteDatabaseCatByIdRequest) (*tablesv1.DeleteDatabaseCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDatabaseCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDatabaseCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteDatabaseCatById восстанавливает мягко удалённую строку dc.database_cat.
func (t *TablesApiV1) UndeleteDatabaseCatById(ctx context.Context, req *tablesv1.UndeleteDatabaseCatByIdRequest) (*tablesv1.UndeleteDatabaseCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDatabaseCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDatabaseCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}
