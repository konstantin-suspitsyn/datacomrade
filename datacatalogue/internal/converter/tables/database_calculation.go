package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DatabaseCalculationToProto переводит строку dc.database_calculation в сущность gRPC.
func DatabaseCalculationToProto(row tables_model.DcDatabaseCalculation) *tablesv1.DatabaseCalculation {
	return &tablesv1.DatabaseCalculation{
		Id:                row.ID,
		DatabaseCatId:     row.DatabaseCatID,
		CalculationTypeId: row.CalculationTypeID,
		CreatedAt:         converter.TimeToProto(row.CreatedAt),
		UpdatedAt:         converter.TimeToProto(row.UpdatedAt),
		IsDeleted:         row.IsDeleted,
		UserId:            row.UserID,
	}
}

// DatabaseCalculationsToProto переводит список строк dc.database_calculation в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DatabaseCalculationsToProto(rows []tables_model.DcDatabaseCalculation) []*tablesv1.DatabaseCalculation {
	items := make([]*tablesv1.DatabaseCalculation, 0, len(rows))

	for _, row := range rows {
		items = append(items, DatabaseCalculationToProto(row))
	}

	return items
}

// ToCreateDatabaseCalculationParams собирает параметры вставки dc.database_calculation из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDatabaseCalculationParams(req *tablesv1.CreateDatabaseCalculationRequest) tables_model.CreateDatabaseCalculationParams {
	return tables_model.CreateDatabaseCalculationParams{
		DatabaseCatID:     req.GetDatabaseCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		UserID:            req.GetUserId(),
	}
}

// ToUpdateDatabaseCalculationByIdParams собирает параметры обновления dc.database_calculation из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateDatabaseCalculationByIdParams(req *tablesv1.UpdateDatabaseCalculationByIdRequest) tables_model.UpdateDatabaseCalculationByIdParams {
	return tables_model.UpdateDatabaseCalculationByIdParams{
		ID:                req.GetId(),
		DatabaseCatID:     req.GetDatabaseCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		UserID:            req.GetUserId(),
	}
}
