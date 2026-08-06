package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// AliasToProto переводит строку dc.alias в сущность gRPC.
func AliasToProto(row tables_model.DcAlias) *tablesv1.Alias {
	return &tablesv1.Alias{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.NullTimeToProto(row.UpdatedAt),
		IsDeleted:   row.IsDeleted,
		UserId:      row.UserID,
	}
}

// AliasesToProto переводит список строк dc.alias в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func AliasesToProto(rows []tables_model.DcAlias) []*tablesv1.Alias {
	items := make([]*tablesv1.Alias, 0, len(rows))

	for _, row := range rows {
		items = append(items, AliasToProto(row))
	}

	return items
}

// ToCreateAliasParams собирает параметры вставки dc.alias из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateAliasParams(req *tablesv1.CreateAliasRequest) tables_model.CreateAliasParams {
	return tables_model.CreateAliasParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateAliasByIdParams собирает параметры обновления dc.alias из запроса gRPC.
func ToUpdateAliasByIdParams(req *tablesv1.UpdateAliasByIdRequest) tables_model.UpdateAliasByIdParams {
	return tables_model.UpdateAliasByIdParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		IsDeleted:   req.GetIsDeleted(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
		ID:          req.GetId(),
	}
}

// ToDeleteAliasByIdParams собирает параметры мягкого удаления dc.alias из запроса gRPC.
func ToDeleteAliasByIdParams(req *tablesv1.DeleteAliasByIdRequest) tables_model.DeleteAliasByIdParams {
	return tables_model.DeleteAliasByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteAliasByIdParams собирает параметры обратного удаления dc.alias из запроса gRPC.
func ToUndeleteAliasByIdParams(req *tablesv1.UndeleteAliasByIdRequest) tables_model.UndeleteAliasByIdParams {
	return tables_model.UndeleteAliasByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetAliasesDeletedParams собирает параметры страницы dc.alias из запроса gRPC.
func ToGetAliasesDeletedParams(req *tablesv1.GetAliasesDeletedRequest) tables_model.GetAliasesDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetAliasesDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetAliasesParams собирает параметры страницы dc.alias из запроса gRPC.
func ToGetAliasesParams(req *tablesv1.GetAliasesRequest) tables_model.GetAliasesParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetAliasesParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
