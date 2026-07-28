package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// AliasToProto переводит строку dc.alias в сущность gRPC.
func AliasToProto(row tables_model.DcAlias) *tablesv1.Alias {
	return &tablesv1.Alias{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.NullTimeToProto(row.UpdatedAt),
		IsDeleted:   row.IsDeleted,
		UserId:      row.UserID,
	}
}

// AliasesToProto переводит список строк dc.alias в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func AliasesToProto(rows []tables_model.DcAlias) []*tablesv1.Alias {
	items := make([]*tablesv1.Alias, 0, len(rows))

	for _, row := range rows {
		items = append(items, AliasToProto(row))
	}

	return items
}

// ToCreateAliasParams собирает параметры вставки dc.alias из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateAliasParams(req *tablesv1.CreateAliasRequest) tables_model.CreateAliasParams {
	return tables_model.CreateAliasParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		UserID:      req.GetUserId(),
	}
}

// ToUpdateAliasByIdParams собирает параметры обновления dc.alias из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateAliasByIdParams(req *tablesv1.UpdateAliasByIdRequest) tables_model.UpdateAliasByIdParams {
	return tables_model.UpdateAliasByIdParams{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		UserID:      req.GetUserId(),
	}
}
