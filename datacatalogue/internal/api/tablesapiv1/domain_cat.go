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

// GetDomainCatById отдаёт активную строку dc.domain_cat по id.
func (t *TablesApiV1) GetDomainCatById(ctx context.Context, req *tablesv1.GetDomainCatByIdRequest) (*tablesv1.GetDomainCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDomainCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDomainCatByIdResponse{DomainCat: tablesconv.DomainCatToProto(row)}, nil
}

// GetDomainCats отдаёт все активные строки dc.domain_cat.
func (t *TablesApiV1) GetDomainCats(ctx context.Context, req *tablesv1.GetDomainCatsRequest) (*tablesv1.GetDomainCatsResponse, error) {
	rows, err := t.services.TablesService.GetDomainCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDomainCatsResponse{DomainCats: tablesconv.DomainCatsToProto(rows)}, nil
}

// GetDeletedDomainCatById отдаёт мягко удалённую строку dc.domain_cat по id.
func (t *TablesApiV1) GetDeletedDomainCatById(ctx context.Context, req *tablesv1.GetDeletedDomainCatByIdRequest) (*tablesv1.GetDeletedDomainCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedDomainCatById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDomainCatByIdResponse{DomainCat: tablesconv.DomainCatToProto(row)}, nil
}

// GetDeletedDomainCats отдаёт все мягко удалённые строки dc.domain_cat.
func (t *TablesApiV1) GetDeletedDomainCats(ctx context.Context, req *tablesv1.GetDeletedDomainCatsRequest) (*tablesv1.GetDeletedDomainCatsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedDomainCats(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedDomainCatsResponse{DomainCats: tablesconv.DomainCatsToProto(rows)}, nil
}

// CreateDomainCat вставляет строку dc.domain_cat и отдаёт её целиком.
func (t *TablesApiV1) CreateDomainCat(ctx context.Context, req *tablesv1.CreateDomainCatRequest) (*tablesv1.CreateDomainCatResponse, error) {
	if err := tablesvalidation.ValidateCreateDomainCat(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateDomainCat(ctx, tablesconv.ToCreateDomainCatParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateDomainCatResponse{DomainCat: tablesconv.DomainCatToProto(row)}, nil
}

// UpdateDomainCatById обновляет активную строку dc.domain_cat и отдаёт её целиком.
func (t *TablesApiV1) UpdateDomainCatById(ctx context.Context, req *tablesv1.UpdateDomainCatByIdRequest) (*tablesv1.UpdateDomainCatByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateDomainCatById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateDomainCatById(ctx, tablesconv.ToUpdateDomainCatByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateDomainCatByIdResponse{DomainCat: tablesconv.DomainCatToProto(row)}, nil
}

// DeleteDomainCatById мягко удаляет строку dc.domain_cat.
func (t *TablesApiV1) DeleteDomainCatById(ctx context.Context, req *tablesv1.DeleteDomainCatByIdRequest) (*tablesv1.DeleteDomainCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteDomainCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteDomainCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteDomainCatById восстанавливает мягко удалённую строку dc.domain_cat.
func (t *TablesApiV1) UndeleteDomainCatById(ctx context.Context, req *tablesv1.UndeleteDomainCatByIdRequest) (*tablesv1.UndeleteDomainCatByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteDomainCatById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteDomainCatByIdResponse{Empty: &emptypb.Empty{}}, nil
}
