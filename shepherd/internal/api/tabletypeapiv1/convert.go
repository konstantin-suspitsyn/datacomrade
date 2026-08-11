package tabletypeapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toTableType(row *tablesv1.TableType) openapiv1.TableType {
	return openapiv1.TableType{
		Id:          row.GetId(),
		Name:        row.GetName(),
		Description: row.GetDescription(),
		IsDeleted:   row.GetIsDeleted(),
		CreatedAt:   row.GetCreatedAt().AsTime(),
		UpdatedAt:   row.GetUpdatedAt().AsTime(),
		UserId:      row.GetUserId(),
	}
}

func toTableTypeSlice(rows []*tablesv1.TableType) []openapiv1.TableType {
	data := make([]openapiv1.TableType, 0, len(rows))
	for _, row := range rows {
		data = append(data, toTableType(row))
	}

	return data
}

func toTableTypePage(rows []*tablesv1.TableType, p *tablesv1.Pagination) openapiv1.TableTypePage {
	return openapiv1.TableTypePage{
		Data: toTableTypeSlice(rows),
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
