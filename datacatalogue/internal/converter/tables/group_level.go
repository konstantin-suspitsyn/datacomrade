package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GroupLevelToProto переводит строку dc.group_levels в сущность gRPC.
func GroupLevelToProto(row tables_model.DcGroupLevel) *tablesv1.GroupLevel {
	return &tablesv1.GroupLevel{
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

// GroupLevelsToProto переводит список строк dc.group_levels в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func GroupLevelsToProto(rows []tables_model.DcGroupLevel) []*tablesv1.GroupLevel {
	items := make([]*tablesv1.GroupLevel, 0, len(rows))

	for _, row := range rows {
		items = append(items, GroupLevelToProto(row))
	}

	return items
}

// ToCreateGroupLevelParams собирает параметры вставки dc.group_levels из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateGroupLevelParams(req *tablesv1.CreateGroupLevelRequest) tables_model.CreateGroupLevelParams {
	return tables_model.CreateGroupLevelParams{
		ColumnID:       req.GetColumnId(),
		ParentColumnID: req.GetParentColumnId(),
		Level:          int16(req.GetLevel()),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetUserExternalId()),
	}
}

// ToUpdateGroupLevelByIdParams собирает параметры обновления dc.group_levels из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateGroupLevelByIdParams(req *tablesv1.UpdateGroupLevelByIdRequest) tables_model.UpdateGroupLevelByIdParams {
	return tables_model.UpdateGroupLevelByIdParams{
		ID:             req.GetId(),
		ColumnID:       req.GetColumnId(),
		ParentColumnID: req.GetParentColumnId(),
		Level:          int16(req.GetLevel()),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetUserExternalId()),
	}
}
