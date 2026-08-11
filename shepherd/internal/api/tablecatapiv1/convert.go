package tablecatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toTableCat(row *tablesv1.TableCat) openapiv1.TableCat {
	return openapiv1.TableCat{
		Id:          row.GetId(),
		Name:        row.GetName(),
		Description: row.GetDescription(),
		SchemaId:    row.GetSchemaId(),
		TableTypeId: row.GetTableTypeId(),
		DomainId:    row.GetDomainId(),
		IsDeleted:   row.GetIsDeleted(),
		CreatedAt:   row.GetCreatedAt().AsTime(),
		UpdatedAt:   row.GetUpdatedAt().AsTime(),
		IsGetDict:   row.GetIsGetDict(),
		UserId:      row.GetUserId(),
	}
}

func toTableCatSlice(rows []*tablesv1.TableCat) []openapiv1.TableCat {
	data := make([]openapiv1.TableCat, 0, len(rows))
	for _, row := range rows {
		data = append(data, toTableCat(row))
	}

	return data
}

func toTableCatPage(rows []*tablesv1.TableCat, p *tablesv1.Pagination) openapiv1.TableCatPage {
	return openapiv1.TableCatPage{
		Data: toTableCatSlice(rows),
		Pagination: openapiv1.Pagination{
			Page:       p.GetPage(),
			PerPage:    p.GetPerPage(),
			TotalItems: p.GetTotalItems(),
			TotalPages: p.GetTotalPages(),
		},
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
