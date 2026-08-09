package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DatabaseCatToProto переводит строку dc.database_cat в сущность gRPC.
func DatabaseCatToProto(row tables_model.DcDatabaseCat) *tablesv1.DatabaseCat {
	return &tablesv1.DatabaseCat{
		Id:             row.ID,
		Name:           row.Name,
		HostId:         row.HostID,
		DatabaseTypeId: row.DatabaseTypeID,
		Description:    row.Description,
		IsDeleted:      row.IsDeleted,
		CreatedAt:      converter.TimeToProto(row.CreatedAt),
		UpdatedAt:      converter.TimeToProto(row.UpdatedAt),
		UserId:         row.UserID,
	}
}

// DatabaseCatsToProto переводит список строк dc.database_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DatabaseCatsToProto(rows []tables_model.DcDatabaseCat) []*tablesv1.DatabaseCat {
	items := make([]*tablesv1.DatabaseCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, DatabaseCatToProto(row))
	}

	return items
}

// ToCreateDatabaseCatParams собирает параметры вставки dc.database_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDatabaseCatParams(req *tablesv1.CreateDatabaseCatRequest) tables_model.CreateDatabaseCatParams {
	return tables_model.CreateDatabaseCatParams{
		Name:           req.GetName(),
		HostID:         req.GetHostId(),
		DatabaseTypeID: req.GetDatabaseTypeId(),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateDatabaseCatByIdParams собирает параметры обновления dc.database_cat из запроса gRPC.
func ToUpdateDatabaseCatByIdParams(req *tablesv1.UpdateDatabaseCatByIdRequest) tables_model.UpdateDatabaseCatByIdParams {
	return tables_model.UpdateDatabaseCatByIdParams{
		Name:           req.GetName(),
		HostID:         req.GetHostId(),
		DatabaseTypeID: req.GetDatabaseTypeId(),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetExternalId()),
		ID:             req.GetId(),
	}
}

// ToDeleteDatabaseCatByIdParams собирает параметры мягкого удаления dc.database_cat из запроса gRPC.
func ToDeleteDatabaseCatByIdParams(req *tablesv1.DeleteDatabaseCatByIdRequest) tables_model.DeleteDatabaseCatByIdParams {
	return tables_model.DeleteDatabaseCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteDatabaseCatByIdParams собирает параметры обратного удаления dc.database_cat из запроса gRPC.
func ToUndeleteDatabaseCatByIdParams(req *tablesv1.UndeleteDatabaseCatByIdRequest) tables_model.UndeleteDatabaseCatByIdParams {
	return tables_model.UndeleteDatabaseCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetDatabaseCatsParams собирает параметры страницы dc.database_cat из запроса gRPC.
func ToGetDatabaseCatsParams(req *tablesv1.GetDatabaseCatsRequest) tables_model.GetDatabaseCatsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCatsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetDatabaseCatsDeletedParams собирает параметры страницы dc.database_cat из запроса gRPC.
func ToGetDatabaseCatsDeletedParams(req *tablesv1.GetDatabaseCatsDeletedRequest) tables_model.GetDatabaseCatsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCatsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetDatabaseCatsByHostIdDatabaseTypeIdParams собирает параметры страницы dc.database_cat из запроса gRPC.
func ToGetDatabaseCatsByHostIdDatabaseTypeIdParams(req *tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest) tables_model.GetDatabaseCatsByHostIdDatabaseTypeIdParams {
	if req == nil {
		req = &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{}
	}

	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCatsByHostIdDatabaseTypeIdParams{
		Order:          req.GetOrder(),
		OrderBy:        req.GetOrderBy(),
		Page:           page,
		PageLimit:      limit,
		HostID:         converter.ProtoToNullInt64(req.HostId),
		DatabaseTypeID: converter.ProtoToNullInt64(req.DatabaseTypeId),
	}
}

// ToGetDatabaseCatsSearchNameParams собирает параметры страницы dc.database_cat из запроса gRPC.
func ToGetDatabaseCatsSearchNameParams(req *tablesv1.GetDatabaseCatsSearchNameRequest) tables_model.GetDatabaseCatsSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCatsSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}
