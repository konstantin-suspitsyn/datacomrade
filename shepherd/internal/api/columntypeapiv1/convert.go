package columntypeapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toColumnType(row *tablesv1.ColumnType) openapiv1.ColumnType {
	return openapiv1.ColumnType{
		Id:          row.GetId(),
		Name:        row.GetName(),
		Description: row.GetDescription(),
		IsDeleted:   row.GetIsDeleted(),
		CreatedAt:   row.GetCreatedAt().AsTime(),
		UpdatedAt:   row.GetUpdatedAt().AsTime(),
		UserId:      row.GetUserId(),
	}
}

func toColumnTypeSlice(rows []*tablesv1.ColumnType) []openapiv1.ColumnType {
	data := make([]openapiv1.ColumnType, 0, len(rows))
	for _, row := range rows {
		data = append(data, toColumnType(row))
	}

	return data
}

func toColumnTypePage(rows []*tablesv1.ColumnType, p *tablesv1.Pagination) openapiv1.ColumnTypePage {
	return openapiv1.ColumnTypePage{
		Data: toColumnTypeSlice(rows),
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
