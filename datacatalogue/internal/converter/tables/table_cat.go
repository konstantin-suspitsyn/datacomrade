package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// TableCatToProto переводит строку dc.table_cat в сущность gRPC.
func TableCatToProto(row tables_model.DcTableCat) *tablesv1.TableCat {
	return &tablesv1.TableCat{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		SchemaId:    row.SchemaID,
		TableTypeId: row.TableTypeID,
		DomainId:    row.DomainID,
		IsDeleted:   row.IsDeleted,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		IsGetDict:   row.IsGetDict,
		UserId:      row.UserID,
	}
}

// TableCatsToProto переводит список строк dc.table_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func TableCatsToProto(rows []tables_model.DcTableCat) []*tablesv1.TableCat {
	items := make([]*tablesv1.TableCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, TableCatToProto(row))
	}

	return items
}

// ToCreateTableCatParams собирает параметры вставки dc.table_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateTableCatParams(req *tablesv1.CreateTableCatRequest) tables_model.CreateTableCatParams {
	return tables_model.CreateTableCatParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		SchemaID:    req.GetSchemaId(),
		TableTypeID: req.GetTableTypeId(),
		DomainID:    req.GetDomainId(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateTableCatByIdParams собирает параметры обновления dc.table_cat из запроса gRPC.
func ToUpdateTableCatByIdParams(req *tablesv1.UpdateTableCatByIdRequest) tables_model.UpdateTableCatByIdParams {
	return tables_model.UpdateTableCatByIdParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		SchemaID:    req.GetSchemaId(),
		TableTypeID: req.GetTableTypeId(),
		DomainID:    req.GetDomainId(),
		IsGetDict:   req.GetIsGetDict(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
		ID:          req.GetId(),
	}
}

// ToDeleteTableCatByIdParams собирает параметры мягкого удаления dc.table_cat из запроса gRPC.
func ToDeleteTableCatByIdParams(req *tablesv1.DeleteTableCatByIdRequest) tables_model.DeleteTableCatByIdParams {
	return tables_model.DeleteTableCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteTableCatByIdParams собирает параметры обратного удаления dc.table_cat из запроса gRPC.
func ToUndeleteTableCatByIdParams(req *tablesv1.UndeleteTableCatByIdRequest) tables_model.UndeleteTableCatByIdParams {
	return tables_model.UndeleteTableCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetTableCatsParams собирает параметры страницы dc.table_cat из запроса gRPC.
func ToGetTableCatsParams(req *tablesv1.GetTableCatsRequest) tables_model.GetTableCatsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetTableCatsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetTableCatsDeletedParams собирает параметры страницы dc.table_cat из запроса gRPC.
func ToGetTableCatsDeletedParams(req *tablesv1.GetTableCatsDeletedRequest) tables_model.GetTableCatsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetTableCatsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
