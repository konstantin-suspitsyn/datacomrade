package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// TableTypeToProto переводит строку dc.table_type в сущность gRPC.
func TableTypeToProto(row tables_model.DcTableType) *tablesv1.TableType {
	return &tablesv1.TableType{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		IsDeleted:   row.IsDeleted,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		UserId:      row.UserID,
	}
}

// TableTypesToProto переводит список строк dc.table_type в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func TableTypesToProto(rows []tables_model.DcTableType) []*tablesv1.TableType {
	items := make([]*tablesv1.TableType, 0, len(rows))

	for _, row := range rows {
		items = append(items, TableTypeToProto(row))
	}

	return items
}

// ToCreateTableTypeParams собирает параметры вставки dc.table_type из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateTableTypeParams(req *tablesv1.CreateTableTypeRequest) tables_model.CreateTableTypeParams {
	return tables_model.CreateTableTypeParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		UserID:      req.GetUserId(),
	}
}

// ToUpdateTableTypeByIdParams собирает параметры обновления dc.table_type из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateTableTypeByIdParams(req *tablesv1.UpdateTableTypeByIdRequest) tables_model.UpdateTableTypeByIdParams {
	return tables_model.UpdateTableTypeByIdParams{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		UserID:      req.GetUserId(),
	}
}
