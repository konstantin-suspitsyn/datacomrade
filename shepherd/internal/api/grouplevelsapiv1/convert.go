package grouplevelsapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toGroupLevels(row *tablesv1.GroupLevels) openapiv1.GroupLevels {
	return openapiv1.GroupLevels{
		Id:             row.GetId(),
		ColumnId:       row.GetColumnId(),
		ParentColumnId: row.GetParentColumnId(),
		Level:          row.GetLevel(),
		Description:    row.GetDescription(),
		CreatedAt:      row.GetCreatedAt().AsTime(),
		UpdatedAt:      row.GetUpdatedAt().AsTime(),
		IsDeleted:      row.GetIsDeleted(),
		UserId:         row.GetUserId(),
	}
}

func toGroupLevelsSlice(rows []*tablesv1.GroupLevels) []openapiv1.GroupLevels {
	data := make([]openapiv1.GroupLevels, 0, len(rows))
	for _, row := range rows {
		data = append(data, toGroupLevels(row))
	}

	return data
}

func toGroupLevelsPage(rows []*tablesv1.GroupLevels, p *tablesv1.Pagination) openapiv1.GroupLevelsPage {
	return openapiv1.GroupLevelsPage{
		Data: toGroupLevelsSlice(rows),
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
