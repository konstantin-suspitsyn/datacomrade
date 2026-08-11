package schemacatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toSchemaCat(row *tablesv1.SchemaCat) openapiv1.SchemaCat {
	return openapiv1.SchemaCat{
		Id:         row.GetId(),
		DatabaseId: row.GetDatabaseId(),
		Name:       row.GetName(),
		IsDeleted:  row.GetIsDeleted(),
		CreatedAt:  row.GetCreatedAt().AsTime(),
		UpdatedAt:  row.GetUpdatedAt().AsTime(),
		UserId:     row.GetUserId(),
	}
}

func toSchemaCatSlice(rows []*tablesv1.SchemaCat) []openapiv1.SchemaCat {
	data := make([]openapiv1.SchemaCat, 0, len(rows))
	for _, row := range rows {
		data = append(data, toSchemaCat(row))
	}

	return data
}

func toSchemaCatPage(rows []*tablesv1.SchemaCat, p *tablesv1.Pagination) openapiv1.SchemaCatPage {
	return openapiv1.SchemaCatPage{
		Data: toSchemaCatSlice(rows),
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
