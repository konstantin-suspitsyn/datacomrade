package databasecatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toDatabaseCat(row *tablesv1.DatabaseCat) openapiv1.DatabaseCat {
	return openapiv1.DatabaseCat{
		Id:             row.GetId(),
		Name:           row.GetName(),
		HostId:         row.GetHostId(),
		DatabaseTypeId: row.GetDatabaseTypeId(),
		Description:    row.GetDescription(),
		IsDeleted:      row.GetIsDeleted(),
		CreatedAt:      row.GetCreatedAt().AsTime(),
		UpdatedAt:      row.GetUpdatedAt().AsTime(),
		UserId:         row.GetUserId(),
	}
}

func toDatabaseCatPage(rows []*tablesv1.DatabaseCat, p *tablesv1.Pagination) openapiv1.DatabaseCatPage {
	data := make([]openapiv1.DatabaseCat, 0, len(rows))
	for _, row := range rows {
		data = append(data, toDatabaseCat(row))
	}

	return openapiv1.DatabaseCatPage{
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
