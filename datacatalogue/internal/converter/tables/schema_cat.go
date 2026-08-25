package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// SchemaCatToProto переводит строку dc.schema_cat в сущность gRPC.
func SchemaCatToProto(row tables_model.DcSchemaCat) *tablesv1.SchemaCat {
	return &tablesv1.SchemaCat{
		Id:         row.ID,
		DatabaseId: row.DatabaseID,
		Name:       row.Name,
		IsDeleted:  row.IsDeleted,
		CreatedAt:  converter.TimeToProto(row.CreatedAt),
		UpdatedAt:  converter.TimeToProto(row.UpdatedAt),
		UserId:     row.UserID,
	}
}

// SchemaCatsToProto переводит список строк dc.schema_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func SchemaCatsToProto(rows []tables_model.DcSchemaCat) []*tablesv1.SchemaCat {
	items := make([]*tablesv1.SchemaCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, SchemaCatToProto(row))
	}

	return items
}

// ToCreateSchemaCatParams собирает параметры вставки dc.schema_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateSchemaCatParams(req *tablesv1.CreateSchemaCatRequest) tables_model.CreateSchemaCatParams {
	return tables_model.CreateSchemaCatParams{
		DatabaseID: req.GetDatabaseId(),
		Name:       req.GetName(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateSchemaCatByIdParams собирает параметры обновления dc.schema_cat из запроса gRPC.
func ToUpdateSchemaCatByIdParams(req *tablesv1.UpdateSchemaCatByIdRequest) tables_model.UpdateSchemaCatByIdParams {
	return tables_model.UpdateSchemaCatByIdParams{
		DatabaseID: req.GetDatabaseId(),
		Name:       req.GetName(),
		CreatedAt:  converter.ProtoToTime(req.GetCreatedAt()),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToDeleteSchemaCatByIdParams собирает параметры мягкого удаления dc.schema_cat из запроса gRPC.
func ToDeleteSchemaCatByIdParams(req *tablesv1.DeleteSchemaCatByIdRequest) tables_model.DeleteSchemaCatByIdParams {
	return tables_model.DeleteSchemaCatByIdParams{
		UpdatedAt:  converter.ProtoToTime(req.GetUpdatedAt()),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteSchemaCatByIdParams собирает параметры обратного удаления dc.schema_cat из запроса gRPC.
func ToUndeleteSchemaCatByIdParams(req *tablesv1.UndeleteSchemaCatByIdRequest) tables_model.UndeleteSchemaCatByIdParams {
	return tables_model.UndeleteSchemaCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetSchemaCatsParams собирает параметры страницы dc.schema_cat из запроса gRPC.
func ToGetSchemaCatsParams(req *tablesv1.GetSchemaCatsRequest) tables_model.GetSchemaCatsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetSchemaCatsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetSchemaCatsDeletedParams собирает параметры страницы dc.schema_cat из запроса gRPC.
func ToGetSchemaCatsDeletedParams(req *tablesv1.GetSchemaCatsDeletedRequest) tables_model.GetSchemaCatsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetSchemaCatsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
