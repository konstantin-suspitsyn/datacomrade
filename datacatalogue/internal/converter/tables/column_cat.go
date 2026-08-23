package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// ColumnCatToProto переводит строку dc.column_cat в сущность gRPC.
func ColumnCatToProto(row tables_model.DcColumnCat) *tablesv1.ColumnCat {
	return &tablesv1.ColumnCat{
		Id:                row.ID,
		TableId:           row.TableID,
		Name:              row.Name,
		AliasId:           converter.NullInt64ToProto(row.AliasID),
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
	if req == nil {
		req = &tablesv1.CreateColumnCatRequest{}
	}

	return tables_model.CreateColumnCatParams{
		TableID:           req.GetTableId(),
		Name:              req.GetName(),
		AliasID:           converter.ProtoToNullInt64(req.AliasId),
		ColumnTypeID:      req.GetColumnTypeId(),
		Description:       req.GetDescription(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ShowInUi:          req.GetShowInUi(),
		ExternalID:        converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateColumnCatByIdParams собирает параметры обновления dc.column_cat из запроса gRPC.
func ToUpdateColumnCatByIdParams(req *tablesv1.UpdateColumnCatByIdRequest) tables_model.UpdateColumnCatByIdParams {
	if req == nil {
		req = &tablesv1.UpdateColumnCatByIdRequest{}
	}

	return tables_model.UpdateColumnCatByIdParams{
		TableID:           req.GetTableId(),
		Name:              req.GetName(),
		AliasID:           converter.ProtoToNullInt64(req.AliasId),
		ColumnTypeID:      req.GetColumnTypeId(),
		Description:       req.GetDescription(),
		CalculationTypeID: req.GetCalculationTypeId(),
		ShowInUi:          req.GetShowInUi(),
		ExternalID:        converter.ProtoToUUID(req.GetExternalId()),
		ID:                req.GetId(),
	}
}

// ToDeleteColumnCatByIdParams собирает параметры мягкого удаления dc.column_cat из запроса gRPC.
func ToDeleteColumnCatByIdParams(req *tablesv1.DeleteColumnCatByIdRequest) tables_model.DeleteColumnCatByIdParams {
	return tables_model.DeleteColumnCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteColumnCatByIdParams собирает параметры обратного удаления dc.column_cat из запроса gRPC.
func ToUndeleteColumnCatByIdParams(req *tablesv1.UndeleteColumnCatByIdRequest) tables_model.UndeleteColumnCatByIdParams {
	return tables_model.UndeleteColumnCatByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetColumnCatsParams собирает параметры страницы dc.column_cat из запроса gRPC.
func ToGetColumnCatsParams(req *tablesv1.GetColumnCatsRequest) tables_model.GetColumnCatsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetColumnCatsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetColumnCatsDeletedParams собирает параметры страницы dc.column_cat из запроса gRPC.
func ToGetColumnCatsDeletedParams(req *tablesv1.GetColumnCatsDeletedRequest) tables_model.GetColumnCatsDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetColumnCatsDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetColumnCatsOptionalFiltersParams собирает параметры страницы dc.column_cat из запроса gRPC.
func ToGetColumnCatsOptionalFiltersParams(req *tablesv1.GetColumnCatsOptionalFiltersRequest) tables_model.GetColumnCatsOptionalFiltersParams {
	if req == nil {
		req = &tablesv1.GetColumnCatsOptionalFiltersRequest{}
	}

	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetColumnCatsOptionalFiltersParams{
		Order:             req.GetOrder(),
		Page:              page,
		PageLimit:         limit,
		TableID:           converter.ProtoToNullInt64(req.TableId),
		AliasID:           converter.ProtoToNullInt64(req.AliasId),
		ColumnTypeID:      converter.ProtoToNullInt64(req.ColumnTypeId),
		CalculationTypeID: converter.ProtoToNullInt64(req.CalculationTypeId),
		ShowInUi:          converter.ProtoToNullBool(req.ShowInUi),
	}
}

// ToGetColumnCatsSearchNameParams собирает параметры страницы dc.column_cat из запроса gRPC.
func ToGetColumnCatsSearchNameParams(req *tablesv1.GetColumnCatsSearchNameRequest) tables_model.GetColumnCatsSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetColumnCatsSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}
