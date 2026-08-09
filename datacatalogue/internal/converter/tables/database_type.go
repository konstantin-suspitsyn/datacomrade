package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DatabaseTypeToProto переводит строку dc.database_type в сущность gRPC.
func DatabaseTypeToProto(row tables_model.DcDatabaseType) *tablesv1.DatabaseType {
	return &tablesv1.DatabaseType{
		Id:        row.ID,
		Name:      row.Name,
		DbVersion: row.DbVersion,
		IsDeleted: row.IsDeleted,
		CreatedAt: converter.TimeToProto(row.CreatedAt),
		UpdatedAt: converter.TimeToProto(row.UpdatedAt),
		UserId:    row.UserID,
	}
}

// DatabaseTypesToProto переводит список строк dc.database_type в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DatabaseTypesToProto(rows []tables_model.DcDatabaseType) []*tablesv1.DatabaseType {
	items := make([]*tablesv1.DatabaseType, 0, len(rows))

	for _, row := range rows {
		items = append(items, DatabaseTypeToProto(row))
	}

	return items
}

// ToCreateDatabaseTypeParams собирает параметры вставки dc.database_type из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDatabaseTypeParams(req *tablesv1.CreateDatabaseTypeRequest) tables_model.CreateDatabaseTypeParams {
	return tables_model.CreateDatabaseTypeParams{
		Name:       req.GetName(),
		DbVersion:  req.GetDbVersion(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateDatabaseTypeByIdParams собирает параметры обновления dc.database_type из запроса gRPC.
func ToUpdateDatabaseTypeByIdParams(req *tablesv1.UpdateDatabaseTypeByIdRequest) tables_model.UpdateDatabaseTypeByIdParams {
	return tables_model.UpdateDatabaseTypeByIdParams{
		Name:       req.GetName(),
		DbVersion:  req.GetDbVersion(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToDeleteDatabaseTypeByIdParams собирает параметры мягкого удаления dc.database_type из запроса gRPC.
func ToDeleteDatabaseTypeByIdParams(req *tablesv1.DeleteDatabaseTypeByIdRequest) tables_model.DeleteDatabaseTypeByIdParams {
	return tables_model.DeleteDatabaseTypeByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteDatabaseTypeByIdParams собирает параметры обратного удаления dc.database_type из запроса gRPC.
func ToUndeleteDatabaseTypeByIdParams(req *tablesv1.UndeleteDatabaseTypeByIdRequest) tables_model.UndeleteDatabaseTypeByIdParams {
	return tables_model.UndeleteDatabaseTypeByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetDatabaseTypesSearchNameParams собирает параметры страницы dc.database_type из запроса gRPC.
func ToGetDatabaseTypesSearchNameParams(req *tablesv1.GetDatabaseTypesSearchNameRequest) tables_model.GetDatabaseTypesSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseTypesSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}

// ToGetDatabaseTypesParams собирает параметры страницы dc.database_type из запроса gRPC.
func ToGetDatabaseTypesParams(req *tablesv1.GetDatabaseTypesRequest) tables_model.GetDatabaseTypesParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseTypesParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetDatabaseTypesDeletedParams собирает параметры страницы dc.database_type из запроса gRPC.
func ToGetDatabaseTypesDeletedParams(req *tablesv1.GetDatabaseTypesDeletedRequest) tables_model.GetDatabaseTypesDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseTypesDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
