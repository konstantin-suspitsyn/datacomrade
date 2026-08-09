package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// FollowingCalculationToProto переводит строку dc.following_calculation в сущность gRPC.
func FollowingCalculationToProto(row tables_model.DcFollowingCalculation) *tablesv1.FollowingCalculation {
	return &tablesv1.FollowingCalculation{
		Id:                row.ID,
		ColumnCatId:       row.ColumnCatID,
		CalculationTypeId: row.CalculationTypeID,
		CreatedAt:         converter.TimeToProto(row.CreatedAt),
		UpdatedAt:         converter.TimeToProto(row.UpdatedAt),
		IsDeleted:         row.IsDeleted,
		UserId:            row.UserID,
	}
}

// FollowingCalculationsToProto переводит список строк dc.following_calculation в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func FollowingCalculationsToProto(rows []tables_model.DcFollowingCalculation) []*tablesv1.FollowingCalculation {
	items := make([]*tablesv1.FollowingCalculation, 0, len(rows))

	for _, row := range rows {
		items = append(items, FollowingCalculationToProto(row))
	}

	return items
}

// ToCreateFollowingCalculationParams собирает параметры вставки dc.following_calculation из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateFollowingCalculationParams(req *tablesv1.CreateFollowingCalculationRequest) tables_model.CreateFollowingCalculationParams {
	return tables_model.CreateFollowingCalculationParams{
		ColumnCatID:       req.GetColumnCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ExternalID:        converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateFollowingCalculationByIdParams собирает параметры обновления dc.following_calculation из запроса gRPC.
func ToUpdateFollowingCalculationByIdParams(req *tablesv1.UpdateFollowingCalculationByIdRequest) tables_model.UpdateFollowingCalculationByIdParams {
	return tables_model.UpdateFollowingCalculationByIdParams{
		ColumnCatID:       req.GetColumnCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ExternalID:        converter.ProtoToUUID(req.GetExternalId()),
		ID:                req.GetId(),
	}
}

// ToDeleteFollowingCalculationByIdParams собирает параметры мягкого удаления dc.following_calculation из запроса gRPC.
func ToDeleteFollowingCalculationByIdParams(req *tablesv1.DeleteFollowingCalculationByIdRequest) tables_model.DeleteFollowingCalculationByIdParams {
	return tables_model.DeleteFollowingCalculationByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteFollowingCalculationByIdParams собирает параметры обратного удаления dc.following_calculation из запроса gRPC.
func ToUndeleteFollowingCalculationByIdParams(req *tablesv1.UndeleteFollowingCalculationByIdRequest) tables_model.UndeleteFollowingCalculationByIdParams {
	return tables_model.UndeleteFollowingCalculationByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetFollowingCalculationsParams собирает параметры страницы dc.following_calculation из запроса gRPC.
func ToGetFollowingCalculationsParams(req *tablesv1.GetFollowingCalculationsRequest) tables_model.GetFollowingCalculationsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetFollowingCalculationsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetFollowingCalculationsDeletedParams собирает параметры страницы dc.following_calculation из запроса gRPC.
func ToGetFollowingCalculationsDeletedParams(req *tablesv1.GetFollowingCalculationsDeletedRequest) tables_model.GetFollowingCalculationsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetFollowingCalculationsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetFollowingCalculationsFiltersParams собирает параметры страницы dc.following_calculation из запроса gRPC.
func ToGetFollowingCalculationsFiltersParams(req *tablesv1.GetFollowingCalculationsFiltersRequest) tables_model.GetFollowingCalculationsFiltersParams {
	if req == nil {
		req = &tablesv1.GetFollowingCalculationsFiltersRequest{}
	}

	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetFollowingCalculationsFiltersParams{
		Order:             req.GetOrder(),
		Page:              page,
		PageLimit:         limit,
		ColumnCatID:       converter.ProtoToNullInt64(req.ColumnCatId),
		CalculationTypeID: converter.ProtoToNullInt64(req.CalculationTypeId),
	}
}
