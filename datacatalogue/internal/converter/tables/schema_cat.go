package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// SchemaCatToProto переводит строку dc.schema_cat в сущность gRPC.
func SchemaCatToProto(row tables_model.DcSchemaCat) *tablesv1.SchemaCat {
	return &tablesv1.SchemaCat{
		Id:         row.ID,
		DatabaseId: row.DatabaseID,
		Name:       row.Name,
		IsDeleted:  row.IsDeleted,
		CreatedAt:  converter.TimeToProto(row.CreatedAt),
		UpdatedAt:  converter.TimeToProto(row.UpdatedAt),
		UserId:     row.UserID,
	}
}

// SchemaCatsToProto переводит список строк dc.schema_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func SchemaCatsToProto(rows []tables_model.DcSchemaCat) []*tablesv1.SchemaCat {
	items := make([]*tablesv1.SchemaCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, SchemaCatToProto(row))
	}

	return items
}

// ToCreateSchemaCatParams собирает параметры вставки dc.schema_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateSchemaCatParams(req *tablesv1.CreateSchemaCatRequest) tables_model.CreateSchemaCatParams {
	return tables_model.CreateSchemaCatParams{
		DatabaseID: req.GetDatabaseId(),
		Name:       req.GetName(),
		ExternalID: converter.ProtoToUUID(req.GetUserExternalId()),
	}
}

// ToUpdateSchemaCatByIdParams собирает параметры обновления dc.schema_cat из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateSchemaCatByIdParams(req *tablesv1.UpdateSchemaCatByIdRequest) tables_model.UpdateSchemaCatByIdParams {
	return tables_model.UpdateSchemaCatByIdParams{
		ID:         req.GetId(),
		DatabaseID: req.GetDatabaseId(),
		Name:       req.GetName(),
		ExternalID: converter.ProtoToUUID(req.GetUserExternalId()),
	}
}
