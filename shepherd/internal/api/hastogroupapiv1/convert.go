package hastogroupapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toHasToGroup(row *tablesv1.HasToGroup) openapiv1.HasToGroup {
	return openapiv1.HasToGroup{
		Id:          row.GetId(),
		ColumnIdA:   row.GetColumnIdA(),
		ColumnIdB:   row.GetColumnIdB(),
		Description: row.GetDescription(),
		IsDeleted:   row.GetIsDeleted(),
		CreatedAt:   row.GetCreatedAt().AsTime(),
		UpdatedAt:   row.GetUpdatedAt().AsTime(),
		UserId:      row.GetUserId(),
	}
}

func toHasToGroupSlice(rows []*tablesv1.HasToGroup) []openapiv1.HasToGroup {
	data := make([]openapiv1.HasToGroup, 0, len(rows))
	for _, row := range rows {
		data = append(data, toHasToGroup(row))
	}

	return data
}

func toHasToGroupPage(rows []*tablesv1.HasToGroup, p *tablesv1.Pagination) openapiv1.HasToGroupPage {
	return openapiv1.HasToGroupPage{
		Data: toHasToGroupSlice(rows),
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
