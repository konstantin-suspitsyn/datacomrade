package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
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
		ExternalID:  converter.ProtoToUUID(req.GetUserExternalId()),
	}
}

// ToUpdateHasToGroupByIdParams собирает параметры обновления dc.has_to_group из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateHasToGroupByIdParams(req *tablesv1.UpdateHasToGroupByIdRequest) tables_model.UpdateHasToGroupByIdParams {
	return tables_model.UpdateHasToGroupByIdParams{
		ID:          req.GetId(),
		ColumnIDA:   req.GetColumnIdA(),
		ColumnIDB:   req.GetColumnIdB(),
		Description: req.GetDescription(),
		ExternalID:  converter.ProtoToUUID(req.GetUserExternalId()),
	}
}
