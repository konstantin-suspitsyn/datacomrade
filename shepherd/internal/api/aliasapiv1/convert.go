package aliasapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toAlias(row *tablesv1.Alias) openapiv1.Alias {
	return openapiv1.Alias{
		Id:          row.GetId(),
		Name:        row.GetName(),
		Description: row.GetDescription(),
		CreatedAt:   row.GetCreatedAt().AsTime(),
		UpdatedAt:   row.GetUpdatedAt().AsTime(),
		IsDeleted:   row.GetIsDeleted(),
		UserId:      row.GetUserId(),
	}
}

func toAliasSlice(rows []*tablesv1.Alias) []openapiv1.Alias {
	data := make([]openapiv1.Alias, 0, len(rows))
	for _, row := range rows {
		data = append(data, toAlias(row))
	}

	return data
}

func toAliasPage(rows []*tablesv1.Alias, p *tablesv1.Pagination) openapiv1.AliasPage {
	return openapiv1.AliasPage{
		Data: toAliasSlice(rows),
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
