package domaincatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toDomainCat(row *tablesv1.DomainCat) openapiv1.DomainCat {
	return openapiv1.DomainCat{
		Id:         row.GetId(),
		DomainName: row.GetDomainName(),
		IsDeleted:  row.GetIsDeleted(),
		CreatedAt:  row.GetCreatedAt().AsTime(),
		UpdatedAt:  row.GetUpdatedAt().AsTime(),
		UserId:     row.GetUserId(),
	}
}

func toDomainCatSlice(rows []*tablesv1.DomainCat) []openapiv1.DomainCat {
	data := make([]openapiv1.DomainCat, 0, len(rows))
	for _, row := range rows {
		data = append(data, toDomainCat(row))
	}

	return data
}

func toDomainCatPage(rows []*tablesv1.DomainCat, p *tablesv1.Pagination) openapiv1.DomainCatPage {
	return openapiv1.DomainCatPage{
		Data: toDomainCatSlice(rows),
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
