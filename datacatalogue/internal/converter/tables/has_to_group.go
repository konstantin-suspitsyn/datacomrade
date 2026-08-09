package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// HasToGroupToProto переводит строку dc.has_to_group в сущность gRPC.
func HasToGroupToProto(row tables_model.DcHasToGroup) *tablesv1.HasToGroup {
	return &tablesv1.HasToGroup{
		Id:          row.ID,
		ColumnIdA:   row.ColumnIDA,
		ColumnIdB:   row.ColumnIDB,
		Description: row.Description,
		IsDeleted:   row.IsDeleted,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		UserId:      row.UserID,
	}
}

// HasToGroupsToProto переводит список строк dc.has_to_group в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func HasToGroupsToProto(rows []tables_model.DcHasToGroup) []*tablesv1.HasToGroup {
	items := make([]*tablesv1.HasToGroup, 0, len(rows))

	for _, row := range rows {
		items = append(items, HasToGroupToProto(row))
	}

	return items
}

// ToCreateHasToGroupParams собирает параметры вставки dc.has_to_group из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateHasToGroupParams(req *tablesv1.CreateHasToGroupRequest) tables_model.CreateHasToGroupParams {
	return tables_model.CreateHasToGroupParams{
		ColumnIDA:   req.GetColumnIdA(),
		ColumnIDB:   req.GetColumnIdB(),
		Description: req.GetDescription(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateHasToGroupByIdParams собирает параметры обновления dc.has_to_group из запроса gRPC.
func ToUpdateHasToGroupByIdParams(req *tablesv1.UpdateHasToGroupByIdRequest) tables_model.UpdateHasToGroupByIdParams {
	return tables_model.UpdateHasToGroupByIdParams{
		ColumnIDA:   req.GetColumnIdA(),
		ColumnIDB:   req.GetColumnIdB(),
		Description: req.GetDescription(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
		ID:          req.GetId(),
	}
}

// ToDeleteHasToGroupByIdParams собирает параметры мягкого удаления dc.has_to_group из запроса gRPC.
func ToDeleteHasToGroupByIdParams(req *tablesv1.DeleteHasToGroupByIdRequest) tables_model.DeleteHasToGroupByIdParams {
	return tables_model.DeleteHasToGroupByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteHasToGroupByIdParams собирает параметры обратного удаления dc.has_to_group из запроса gRPC.
func ToUndeleteHasToGroupByIdParams(req *tablesv1.UndeleteHasToGroupByIdRequest) tables_model.UndeleteHasToGroupByIdParams {
	return tables_model.UndeleteHasToGroupByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetHasToGroupsSearchNameParams собирает параметры страницы dc.has_to_group из запроса gRPC.
func ToGetHasToGroupsSearchNameParams(req *tablesv1.GetHasToGroupsSearchNameRequest) tables_model.GetHasToGroupsSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetHasToGroupsSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}

// ToGetHasToGroupsParams собирает параметры страницы dc.has_to_group из запроса gRPC.
func ToGetHasToGroupsParams(req *tablesv1.GetHasToGroupsRequest) tables_model.GetHasToGroupsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetHasToGroupsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetHasToGroupsDeletedParams собирает параметры страницы dc.has_to_group из запроса gRPC.
func ToGetHasToGroupsDeletedParams(req *tablesv1.GetHasToGroupsDeletedRequest) tables_model.GetHasToGroupsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetHasToGroupsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
