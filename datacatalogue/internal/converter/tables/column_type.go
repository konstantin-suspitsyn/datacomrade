package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// ColumnTypeToProto переводит строку dc.column_type в сущность gRPC.
func ColumnTypeToProto(row tables_model.DcColumnType) *tablesv1.ColumnType {
	return &tablesv1.ColumnType{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		IsDeleted:   row.IsDeleted,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		UserId:      row.UserID,
	}
}

// ColumnTypesToProto переводит список строк dc.column_type в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func ColumnTypesToProto(rows []tables_model.DcColumnType) []*tablesv1.ColumnType {
	items := make([]*tablesv1.ColumnType, 0, len(rows))

	for _, row := range rows {
		items = append(items, ColumnTypeToProto(row))
	}

	return items
}

// ToCreateColumnTypeParams собирает параметры вставки dc.column_type из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateColumnTypeParams(req *tablesv1.CreateColumnTypeRequest) tables_model.CreateColumnTypeParams {
	return tables_model.CreateColumnTypeParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		UserID:      req.GetUserId(),
	}
}

// ToUpdateColumnTypeByIdParams собирает параметры обновления dc.column_type из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateColumnTypeByIdParams(req *tablesv1.UpdateColumnTypeByIdRequest) tables_model.UpdateColumnTypeByIdParams {
	return tables_model.UpdateColumnTypeByIdParams{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		UserID:      req.GetUserId(),
	}
}
