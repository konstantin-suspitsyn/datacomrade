package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DomainCatToProto переводит строку dc.domain_cat в сущность gRPC.
func DomainCatToProto(row tables_model.DcDomainCat) *tablesv1.DomainCat {
	return &tablesv1.DomainCat{
		Id:         row.ID,
		DomainName: row.DomainName,
		IsDeleted:  row.IsDeleted,
		CreatedAt:  converter.TimeToProto(row.CreatedAt),
		UpdatedAt:  converter.TimeToProto(row.UpdatedAt),
		UserId:     row.UserID,
	}
}

// DomainCatsToProto переводит список строк dc.domain_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DomainCatsToProto(rows []tables_model.DcDomainCat) []*tablesv1.DomainCat {
	items := make([]*tablesv1.DomainCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, DomainCatToProto(row))
	}

	return items
}

// ToCreateDomainCatParams собирает параметры вставки dc.domain_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDomainCatParams(req *tablesv1.CreateDomainCatRequest) tables_model.CreateDomainCatParams {
	return tables_model.CreateDomainCatParams{
		DomainName: req.GetDomainName(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateDomainCatByIdParams собирает параметры обновления dc.domain_cat из запроса gRPC.
func ToUpdateDomainCatByIdParams(req *tablesv1.UpdateDomainCatByIdRequest) tables_model.UpdateDomainCatByIdParams {
	return tables_model.UpdateDomainCatByIdParams{
		DomainName: req.GetDomainName(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToDeleteDomainCatByIdParams собирает параметры мягкого удаления dc.domain_cat из запроса gRPC.
func ToDeleteDomainCatByIdParams(req *tablesv1.DeleteDomainCatByIdRequest) tables_model.DeleteDomainCatByIdParams {
	return tables_model.DeleteDomainCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteDomainCatByIdParams собирает параметры обратного удаления dc.domain_cat из запроса gRPC.
func ToUndeleteDomainCatByIdParams(req *tablesv1.UndeleteDomainCatByIdRequest) tables_model.UndeleteDomainCatByIdParams {
	return tables_model.UndeleteDomainCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetDomainCatsSearchNameParams собирает параметры страницы dc.domain_cat из запроса gRPC.
func ToGetDomainCatsSearchNameParams(req *tablesv1.GetDomainCatsSearchNameRequest) tables_model.GetDomainCatsSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDomainCatsSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}

// ToGetDomainCatsParams собирает параметры страницы dc.domain_cat из запроса gRPC.
func ToGetDomainCatsParams(req *tablesv1.GetDomainCatsRequest) tables_model.GetDomainCatsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDomainCatsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetDomainCatsDeletedParams собирает параметры страницы dc.domain_cat из запроса gRPC.
func ToGetDomainCatsDeletedParams(req *tablesv1.GetDomainCatsDeletedRequest) tables_model.GetDomainCatsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDomainCatsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
