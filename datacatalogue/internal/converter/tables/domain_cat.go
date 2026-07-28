package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// DomainCatToProto переводит строку dc.domain_cat в сущность gRPC.
func DomainCatToProto(row tables_model.DcDomainCat) *tablesv1.DomainCat {
	return &tablesv1.DomainCat{
		Id:         row.ID,
		DomainName: row.DomainName,
		IsDeleted:  row.IsDeleted,
		CreatedAt:  converter.TimeToProto(row.CreatedAt),
		UpdatedAt:  converter.TimeToProto(row.UpdatedAt),
		UserId:     row.UserID,
	}
}

// DomainCatsToProto переводит список строк dc.domain_cat в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DomainCatsToProto(rows []tables_model.DcDomainCat) []*tablesv1.DomainCat {
	items := make([]*tablesv1.DomainCat, 0, len(rows))

	for _, row := range rows {
		items = append(items, DomainCatToProto(row))
	}

	return items
}

// ToCreateDomainCatParams собирает параметры вставки dc.domain_cat из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDomainCatParams(req *tablesv1.CreateDomainCatRequest) tables_model.CreateDomainCatParams {
	return tables_model.CreateDomainCatParams{
		DomainName: req.GetDomainName(),
		UserID:     req.GetUserId(),
	}
}

// ToUpdateDomainCatByIdParams собирает параметры обновления dc.domain_cat из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateDomainCatByIdParams(req *tablesv1.UpdateDomainCatByIdRequest) tables_model.UpdateDomainCatByIdParams {
	return tables_model.UpdateDomainCatByIdParams{
		ID:         req.GetId(),
		DomainName: req.GetDomainName(),
		UserID:     req.GetUserId(),
	}
}
