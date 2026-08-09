package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DatabaseCalculationToProto переводит строку dc.database_calculation в сущность gRPC.
func DatabaseCalculationToProto(row tables_model.DcDatabaseCalculation) *tablesv1.DatabaseCalculation {
	return &tablesv1.DatabaseCalculation{
		Id:                row.ID,
		DatabaseCatId:     row.DatabaseCatID,
		CalculationTypeId: row.CalculationTypeID,
		CreatedAt:         converter.TimeToProto(row.CreatedAt),
		UpdatedAt:         converter.TimeToProto(row.UpdatedAt),
		IsDeleted:         row.IsDeleted,
		UserId:            row.UserID,
	}
}

// DatabaseCalculationsToProto переводит список строк dc.database_calculation в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DatabaseCalculationsToProto(rows []tables_model.DcDatabaseCalculation) []*tablesv1.DatabaseCalculation {
	items := make([]*tablesv1.DatabaseCalculation, 0, len(rows))

	for _, row := range rows {
		items = append(items, DatabaseCalculationToProto(row))
	}

	return items
}

// ToCreateDatabaseCalculationParams собирает параметры вставки dc.database_calculation из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDatabaseCalculationParams(req *tablesv1.CreateDatabaseCalculationRequest) tables_model.CreateDatabaseCalculationParams {
	return tables_model.CreateDatabaseCalculationParams{
		DatabaseCatID:     req.GetDatabaseCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ExternalID:        converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateDatabaseCalculationByIdParams собирает параметры обновления dc.database_calculation из запроса gRPC.
func ToUpdateDatabaseCalculationByIdParams(req *tablesv1.UpdateDatabaseCalculationByIdRequest) tables_model.UpdateDatabaseCalculationByIdParams {
	return tables_model.UpdateDatabaseCalculationByIdParams{
		DatabaseCatID:     req.GetDatabaseCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ExternalID:        converter.ProtoToUUID(req.GetExternalId()),
		ID:                req.GetId(),
	}
}

// ToDeleteDatabaseCalculationByIdParams собирает параметры мягкого удаления dc.database_calculation из запроса gRPC.
func ToDeleteDatabaseCalculationByIdParams(req *tablesv1.DeleteDatabaseCalculationByIdRequest) tables_model.DeleteDatabaseCalculationByIdParams {
	return tables_model.DeleteDatabaseCalculationByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteDatabaseCalculationByIdParams собирает параметры обратного удаления dc.database_calculation из запроса gRPC.
func ToUndeleteDatabaseCalculationByIdParams(req *tablesv1.UndeleteDatabaseCalculationByIdRequest) tables_model.UndeleteDatabaseCalculationByIdParams {
	return tables_model.UndeleteDatabaseCalculationByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetDatabaseCalculationsSearchNameParams собирает параметры страницы dc.database_calculation из запроса gRPC.
func ToGetDatabaseCalculationsSearchNameParams(req *tablesv1.GetDatabaseCalculationsSearchNameRequest) tables_model.GetDatabaseCalculationsSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCalculationsSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}

// ToGetDatabaseCalculationsFilteredParams собирает параметры страницы dc.database_calculation из запроса gRPC.
func ToGetDatabaseCalculationsFilteredParams(req *tablesv1.GetDatabaseCalculationsFilteredRequest) tables_model.GetDatabaseCalculationsFilteredParams {
	if req == nil {
		req = &tablesv1.GetDatabaseCalculationsFilteredRequest{}
	}

	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCalculationsFilteredParams{
		Order:             req.GetOrder(),
		Page:              page,
		PageLimit:         limit,
		DatabaseCatID:     converter.ProtoToNullInt64(req.DatabaseCatId),
		CalculationTypeID: converter.ProtoToNullInt64(req.CalculationTypeId),
	}
}

// ToGetDatabaseCalculationsDeletedParams собирает параметры страницы dc.database_calculation из запроса gRPC.
func ToGetDatabaseCalculationsDeletedParams(req *tablesv1.GetDatabaseCalculationsDeletedRequest) tables_model.GetDatabaseCalculationsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCalculationsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetDatabaseCalculationsParams собирает параметры страницы dc.database_calculation из запроса gRPC.
func ToGetDatabaseCalculationsParams(req *tablesv1.GetDatabaseCalculationsRequest) tables_model.GetDatabaseCalculationsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetDatabaseCalculationsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
