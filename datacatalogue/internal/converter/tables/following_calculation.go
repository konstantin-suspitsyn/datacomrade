package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// FollowingCalculationToProto переводит строку dc.following_calculation в сущность gRPC.
func FollowingCalculationToProto(row tables_model.DcFollowingCalculation) *tablesv1.FollowingCalculation {
	return &tablesv1.FollowingCalculation{
		Id:                row.ID,
		ColumnCatId:       row.ColumnCatID,
		CalculationTypeId: row.CalculationTypeID,
		CreatedAt:         converter.TimeToProto(row.CreatedAt),
		UpdatedAt:         converter.TimeToProto(row.UpdatedAt),
		IsDeleted:         row.IsDeleted,
		UserId:            row.UserID,
	}
}

// FollowingCalculationsToProto переводит список строк dc.following_calculation в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func FollowingCalculationsToProto(rows []tables_model.DcFollowingCalculation) []*tablesv1.FollowingCalculation {
	items := make([]*tablesv1.FollowingCalculation, 0, len(rows))

	for _, row := range rows {
		items = append(items, FollowingCalculationToProto(row))
	}

	return items
}

// ToCreateFollowingCalculationParams собирает параметры вставки dc.following_calculation из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateFollowingCalculationParams(req *tablesv1.CreateFollowingCalculationRequest) tables_model.CreateFollowingCalculationParams {
	return tables_model.CreateFollowingCalculationParams{
		ColumnCatID:       req.GetColumnCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		UserID:            req.GetUserId(),
	}
}

// ToUpdateFollowingCalculationByIdParams собирает параметры обновления dc.following_calculation из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateFollowingCalculationByIdParams(req *tablesv1.UpdateFollowingCalculationByIdRequest) tables_model.UpdateFollowingCalculationByIdParams {
	return tables_model.UpdateFollowingCalculationByIdParams{
		ID:                req.GetId(),
		ColumnCatID:       req.GetColumnCatId(),
		CalculationTypeID: req.GetCalculationTypeId(),
		UserID:            req.GetUserId(),
	}
}
