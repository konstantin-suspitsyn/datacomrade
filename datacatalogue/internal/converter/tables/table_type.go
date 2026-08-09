package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
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
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateTableTypeByIdParams собирает параметры обновления dc.table_type из запроса gRPC.
func ToUpdateTableTypeByIdParams(req *tablesv1.UpdateTableTypeByIdRequest) tables_model.UpdateTableTypeByIdParams {
	return tables_model.UpdateTableTypeByIdParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
		ID:          req.GetId(),
	}
}

// ToDeleteTableTypeByIdParams собирает параметры мягкого удаления dc.table_type из запроса gRPC.
func ToDeleteTableTypeByIdParams(req *tablesv1.DeleteTableTypeByIdRequest) tables_model.DeleteTableTypeByIdParams {
	return tables_model.DeleteTableTypeByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteTableTypeByIdParams собирает параметры обратного удаления dc.table_type из запроса gRPC.
func ToUndeleteTableTypeByIdParams(req *tablesv1.UndeleteTableTypeByIdRequest) tables_model.UndeleteTableTypeByIdParams {
	return tables_model.UndeleteTableTypeByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetTableTypesParams собирает параметры страницы dc.table_type из запроса gRPC.
func ToGetTableTypesParams(req *tablesv1.GetTableTypesRequest) tables_model.GetTableTypesParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetTableTypesParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetTableTypesDeletedParams собирает параметры страницы dc.table_type из запроса gRPC.
func ToGetTableTypesDeletedParams(req *tablesv1.GetTableTypesDeletedRequest) tables_model.GetTableTypesDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetTableTypesDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
