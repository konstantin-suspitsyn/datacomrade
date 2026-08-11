package databasetypeapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toDatabaseType(row *tablesv1.DatabaseType) openapiv1.DatabaseType {
	return openapiv1.DatabaseType{
		Id:        row.GetId(),
		Name:      row.GetName(),
		DbVersion: row.GetDbVersion(),
		IsDeleted: row.GetIsDeleted(),
		CreatedAt: row.GetCreatedAt().AsTime(),
		UpdatedAt: row.GetUpdatedAt().AsTime(),
		UserId:    row.GetUserId(),
	}
}

func toDatabaseTypePage(rows []*tablesv1.DatabaseType, p *tablesv1.Pagination) openapiv1.DatabaseTypePage {
	data := make([]openapiv1.DatabaseType, 0, len(rows))
	for _, row := range rows {
		data = append(data, toDatabaseType(row))
	}

	return openapiv1.DatabaseTypePage{
		Data: data,
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
