package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// CalculationTypeToProto переводит строку dc.calculation_type в сущность gRPC.
func CalculationTypeToProto(row tables_model.DcCalculationType) *tablesv1.CalculationType {
	return &tablesv1.CalculationType{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		IsDeleted:   row.IsDeleted,
	}
}

// CalculationTypesToProto переводит список строк dc.calculation_type в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func CalculationTypesToProto(rows []tables_model.DcCalculationType) []*tablesv1.CalculationType {
	items := make([]*tablesv1.CalculationType, 0, len(rows))

	for _, row := range rows {
		items = append(items, CalculationTypeToProto(row))
	}

	return items
}

// ToCreateCalculationTypeParams собирает параметры вставки dc.calculation_type из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateCalculationTypeParams(req *tablesv1.CreateCalculationTypeRequest) tables_model.CreateCalculationTypeParams {
	return tables_model.CreateCalculationTypeParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
}

// ToUpdateCalculationTypeByIdParams собирает параметры обновления dc.calculation_type из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateCalculationTypeByIdParams(req *tablesv1.UpdateCalculationTypeByIdRequest) tables_model.UpdateCalculationTypeByIdParams {
	return tables_model.UpdateCalculationTypeByIdParams{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
}
