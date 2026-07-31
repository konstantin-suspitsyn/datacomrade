package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DatabaseCatToProto переводит строку dc.database_cat в сущность gRPC.
func DatabaseCatToProto(row tables_model.DcDatabaseCat) *tablesv1.DatabaseCat {
	return &tablesv1.DatabaseCat{
		Id:             row.ID,
		Name:           row.Name,
		HostId:         row.HostID,
		DatabaseTypeId: row.DatabaseTypeID,
		Description:    row.Description,
		IsDeleted:      row.IsDeleted,
		CreatedAt:      converter.TimeToProto(row.CreatedAt),
		UpdatedAt:      converter.TimeToProto(row.UpdatedAt),
		UserId:         row.UserID,
	}
}

// DatabaseCatsToProto переводит список строк dc.database_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DatabaseCatsToProto(rows []tables_model.DcDatabaseCat) []*tablesv1.DatabaseCat {
	items := make([]*tablesv1.DatabaseCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, DatabaseCatToProto(row))
	}

	return items
}

// ToCreateDatabaseCatParams собирает параметры вставки dc.database_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDatabaseCatParams(req *tablesv1.CreateDatabaseCatRequest) tables_model.CreateDatabaseCatParams {
	return tables_model.CreateDatabaseCatParams{
		Name:           req.GetName(),
		HostID:         req.GetHostId(),
		DatabaseTypeID: req.GetDatabaseTypeId(),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetUserExternalId()),
	}
}

// ToUpdateDatabaseCatByIdParams собирает параметры обновления dc.database_cat из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateDatabaseCatByIdParams(req *tablesv1.UpdateDatabaseCatByIdRequest) tables_model.UpdateDatabaseCatByIdParams {
	return tables_model.UpdateDatabaseCatByIdParams{
		ID:             req.GetId(),
		Name:           req.GetName(),
		HostID:         req.GetHostId(),
		DatabaseTypeID: req.GetDatabaseTypeId(),
		Description:    req.GetDescription(),
		ExternalID:     converter.ProtoToUUID(req.GetUserExternalId()),
	}
}
