package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// ColumnCatToProto переводит строку dc.column_cat в сущность gRPC.
func ColumnCatToProto(row tables_model.DcColumnCat) *tablesv1.ColumnCat {
	return &tablesv1.ColumnCat{
		Id:                row.ID,
		TableId:           row.TableID,
		Name:              row.Name,
		AliasId:           row.AliasID,
		ColumnTypeId:      row.ColumnTypeID,
		Description:       row.Description,
		CalculationTypeId: row.CalculationTypeID,
		IsDeleted:         row.IsDeleted,
		ShowInUi:          row.ShowInUi,
		CreatedAt:         converter.TimeToProto(row.CreatedAt),
		UpdatedAt:         converter.TimeToProto(row.UpdatedAt),
		UserId:            row.UserID,
	}
}

// ColumnCatsToProto переводит список строк dc.column_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func ColumnCatsToProto(rows []tables_model.DcColumnCat) []*tablesv1.ColumnCat {
	items := make([]*tablesv1.ColumnCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, ColumnCatToProto(row))
	}

	return items
}

// ToCreateColumnCatParams собирает параметры вставки dc.column_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateColumnCatParams(req *tablesv1.CreateColumnCatRequest) tables_model.CreateColumnCatParams {
	return tables_model.CreateColumnCatParams{
		TableID:           req.GetTableId(),
		Name:              req.GetName(),
		AliasID:           req.GetAliasId(),
		ColumnTypeID:      req.GetColumnTypeId(),
		Description:       req.GetDescription(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ShowInUi:          req.GetShowInUi(),
		ExternalID:        converter.ProtoToUUID(req.GetUserExternalId()),
	}
}

// ToUpdateColumnCatByIdParams собирает параметры обновления dc.column_cat из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateColumnCatByIdParams(req *tablesv1.UpdateColumnCatByIdRequest) tables_model.UpdateColumnCatByIdParams {
	return tables_model.UpdateColumnCatByIdParams{
		ID:                req.GetId(),
		TableID:           req.GetTableId(),
		Name:              req.GetName(),
		AliasID:           req.GetAliasId(),
		ColumnTypeID:      req.GetColumnTypeId(),
		Description:       req.GetDescription(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ShowInUi:          req.GetShowInUi(),
		ExternalID:        converter.ProtoToUUID(req.GetUserExternalId()),
	}
}
