package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// TablesJoinsToProto переводит строку dc.tables_joins в сущность gRPC.
func TablesJoinsToProto(row tables_model.DcTablesJoin) *tablesv1.TablesJoins {
	return &tablesv1.TablesJoins{
		Id:              row.ID,
		TableCatIdLeft:  row.TableCatIDLeft,
		TableCatIdRight: row.TableCatIDRight,
		IsDeleted:       row.IsDeleted,
		CreatedAt:       converter.TimeToProto(row.CreatedAt),
		UpdatedAt:       converter.TimeToProto(row.UpdatedAt),
		UserId:          row.UserID,
	}
}

// TablesJoinsesToProto переводит список строк dc.tables_joins в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func TablesJoinsesToProto(rows []tables_model.DcTablesJoin) []*tablesv1.TablesJoins {
	items := make([]*tablesv1.TablesJoins, 0, len(rows))

	for _, row := range rows {
		items = append(items, TablesJoinsToProto(row))
	}

	return items
}

// GetTablesJoinsByIdRowToProto переводит строку выборки в сообщение gRPC.
func GetTablesJoinsByIdRowToProto(row tables_model.GetTablesJoinsByIdRow) *tablesv1.GetTablesJoinsByIdRow {
	return &tablesv1.GetTablesJoinsByIdRow{
		Id:                row.ID,
		TableCatIdLeft:    row.TableCatIDLeft,
		TableCatIdRight:   row.TableCatIDRight,
		IsDeleted:         row.IsDeleted,
		CreatedAt:         converter.TimeToProto(row.CreatedAt),
		UpdatedAt:         converter.TimeToProto(row.UpdatedAt),
		UserId:            row.UserID,
		TableCatLeftName:  row.TableCatLeftName,
		TableCatRightName: row.TableCatRightName,
	}
}

// GetTablesJoinsRowToProto переводит строку выборки в сообщение gRPC.
func GetTablesJoinsRowToProto(row tables_model.GetTablesJoinsRow) *tablesv1.GetTablesJoinsRow {
	return &tablesv1.GetTablesJoinsRow{
		Id:                row.ID,
		TableCatIdLeft:    row.TableCatIDLeft,
		TableCatIdRight:   row.TableCatIDRight,
		IsDeleted:         row.IsDeleted,
		CreatedAt:         converter.TimeToProto(row.CreatedAt),
		UpdatedAt:         converter.TimeToProto(row.UpdatedAt),
		UserId:            row.UserID,
		TableCatLeftName:  row.TableCatLeftName,
		TableCatRightName: row.TableCatRightName,
	}
}

// GetTablesJoinsRowsToProto переводит список строк выборки в список сообщений gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func GetTablesJoinsRowsToProto(rows []tables_model.GetTablesJoinsRow) []*tablesv1.GetTablesJoinsRow {
	items := make([]*tablesv1.GetTablesJoinsRow, 0, len(rows))

	for _, row := range rows {
		items = append(items, GetTablesJoinsRowToProto(row))
	}

	return items
}

// ToCreateTablesJoinsParams собирает параметры вставки dc.tables_joins из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateTablesJoinsParams(req *tablesv1.CreateTablesJoinsRequest) tables_model.CreateTablesJoinsParams {
	return tables_model.CreateTablesJoinsParams{
		TableCatIDLeft:  req.GetTableCatIdLeft(),
		TableCatIDRight: req.GetTableCatIdRight(),
		ExternalID:      converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateTablesJoinsByIdParams собирает параметры обновления dc.tables_joins из запроса gRPC.
func ToUpdateTablesJoinsByIdParams(req *tablesv1.UpdateTablesJoinsByIdRequest) tables_model.UpdateTablesJoinsByIdParams {
	return tables_model.UpdateTablesJoinsByIdParams{
		TableCatIDLeft:  req.GetTableCatIdLeft(),
		TableCatIDRight: req.GetTableCatIdRight(),
		IsDeleted:       req.GetIsDeleted(),
		ExternalID:      converter.ProtoToUUID(req.GetExternalId()),
		ID:              req.GetId(),
	}
}

// ToDeleteTablesJoinsByIdParams собирает параметры мягкого удаления dc.tables_joins из запроса gRPC.
func ToDeleteTablesJoinsByIdParams(req *tablesv1.DeleteTablesJoinsByIdRequest) tables_model.DeleteTablesJoinsByIdParams {
	return tables_model.DeleteTablesJoinsByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetTablesJoinsParams собирает параметры страницы dc.tables_joins из запроса gRPC.
func ToGetTablesJoinsParams(req *tablesv1.GetTablesJoinsRequest) tables_model.GetTablesJoinsParams {
	if req == nil {
		req = &tablesv1.GetTablesJoinsRequest{}
	}

	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetTablesJoinsParams{
		Order:             req.GetOrder(),
		Page:              page,
		PageLimit:         limit,
		TableCatLeftName:  converter.ProtoToNullString(req.TableCatLeftName),
		TableCatRightName: converter.ProtoToNullString(req.TableCatRightName),
	}
}
