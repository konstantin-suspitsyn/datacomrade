package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DatabaseTypeToProto переводит строку dc.database_type в сущность gRPC.
func DatabaseTypeToProto(row tables_model.DcDatabaseType) *tablesv1.DatabaseType {
	return &tablesv1.DatabaseType{
		Id:        row.ID,
		Name:      row.Name,
		DbVersion: row.DbVersion,
		IsDeleted: row.IsDeleted,
		CreatedAt: converter.TimeToProto(row.CreatedAt),
		UpdatedAt: converter.TimeToProto(row.UpdatedAt),
		UserId:    row.UserID,
	}
}

// DatabaseTypesToProto переводит список строк dc.database_type в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DatabaseTypesToProto(rows []tables_model.DcDatabaseType) []*tablesv1.DatabaseType {
	items := make([]*tablesv1.DatabaseType, 0, len(rows))

	for _, row := range rows {
		items = append(items, DatabaseTypeToProto(row))
	}

	return items
}

// ToCreateDatabaseTypeParams собирает параметры вставки dc.database_type из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDatabaseTypeParams(req *tablesv1.CreateDatabaseTypeRequest) tables_model.CreateDatabaseTypeParams {
	return tables_model.CreateDatabaseTypeParams{
		Name:       req.GetName(),
		DbVersion:  req.GetDbVersion(),
		ExternalID: converter.ProtoToUUID(req.GetUserExternalId()),
	}
}

// ToUpdateDatabaseTypeByIdParams собирает параметры обновления dc.database_type из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateDatabaseTypeByIdParams(req *tablesv1.UpdateDatabaseTypeByIdRequest) tables_model.UpdateDatabaseTypeByIdParams {
	return tables_model.UpdateDatabaseTypeByIdParams{
		ID:         req.GetId(),
		Name:       req.GetName(),
		DbVersion:  req.GetDbVersion(),
		ExternalID: converter.ProtoToUUID(req.GetUserExternalId()),
	}
}
