package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GroupLevelsToProto переводит строку dc.group_levels в сущность gRPC.
func GroupLevelsToProto(row tables_model.DcGroupLevel) *tablesv1.GroupLevels {
	return &tablesv1.GroupLevels{
		Id:             row.ID,
		ColumnId:       row.ColumnID,
		ParentColumnId: row.ParentColumnID,
		Level:          int32(row.Level),
		Description:    row.Description,
		CreatedAt:      converter.TimeToProto(row.CreatedAt),
		UpdatedAt:      converter.TimeToProto(row.UpdatedAt),
		IsDeleted:      row.IsDeleted,
		UserId:         row.UserID,
	}
}

// GroupLevelsesToProto переводит список строк dc.group_levels в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func GroupLevelsesToProto(rows []tables_model.DcGroupLevel) []*tablesv1.GroupLevels {
	items := make([]*tablesv1.GroupLevels, 0, len(rows))

	for _, row := range rows {
		items = append(items, GroupLevelsToProto(row))
	}

	return items
}

// ToCreateGroupLevelsParams собирает параметры вставки dc.group_levels из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateGroupLevelsParams(req *tablesv1.CreateGroupLevelsRequest) tables_model.CreateGroupLevelsParams {
	return tables_model.CreateGroupLevelsParams{
		ColumnID:       req.GetColumnId(),
		ParentColumnID: req.GetParentColumnId(),
		Level:          int16(req.GetLevel()),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateGroupLevelsByIdParams собирает параметры обновления dc.group_levels из запроса gRPC.
func ToUpdateGroupLevelsByIdParams(req *tablesv1.UpdateGroupLevelsByIdRequest) tables_model.UpdateGroupLevelsByIdParams {
	return tables_model.UpdateGroupLevelsByIdParams{
		ColumnID:       req.GetColumnId(),
		ParentColumnID: req.GetParentColumnId(),
		Level:          int16(req.GetLevel()),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetExternalId()),
		ID:             req.GetId(),
	}
}

// ToDeleteGroupLevelsByIdParams собирает параметры мягкого удаления dc.group_levels из запроса gRPC.
func ToDeleteGroupLevelsByIdParams(req *tablesv1.DeleteGroupLevelsByIdRequest) tables_model.DeleteGroupLevelsByIdParams {
	return tables_model.DeleteGroupLevelsByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteGroupLevelsByIdParams собирает параметры обратного удаления dc.group_levels из запроса gRPC.
func ToUndeleteGroupLevelsByIdParams(req *tablesv1.UndeleteGroupLevelsByIdRequest) tables_model.UndeleteGroupLevelsByIdParams {
	return tables_model.UndeleteGroupLevelsByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetGroupLevelsSearchNameParams собирает параметры страницы dc.group_levels из запроса gRPC.
func ToGetGroupLevelsSearchNameParams(req *tablesv1.GetGroupLevelsSearchNameRequest) tables_model.GetGroupLevelsSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetGroupLevelsSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}

// ToGetGroupLevelsParams собирает параметры страницы dc.group_levels из запроса gRPC.
func ToGetGroupLevelsParams(req *tablesv1.GetGroupLevelsRequest) tables_model.GetGroupLevelsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetGroupLevelsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetGroupLevelsFilteredParams собирает параметры страницы dc.group_levels из запроса gRPC.
func ToGetGroupLevelsFilteredParams(req *tablesv1.GetGroupLevelsFilteredRequest) tables_model.GetGroupLevelsFilteredParams {
	if req == nil {
		req = &tablesv1.GetGroupLevelsFilteredRequest{}
	}

	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetGroupLevelsFilteredParams{
		Order:          req.GetOrder(),
		Page:           page,
		PageLimit:      limit,
		ColumnID:       converter.ProtoToNullInt64(req.ColumnId),
		ParentColumnID: converter.ProtoToNullInt64(req.ParentColumnId),
	}
}
