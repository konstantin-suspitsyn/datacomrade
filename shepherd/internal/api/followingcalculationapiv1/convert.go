package followingcalculationapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toFollowingCalculation(row *tablesv1.FollowingCalculation) openapiv1.FollowingCalculation {
	return openapiv1.FollowingCalculation{
		Id:                row.GetId(),
		ColumnCatId:       row.GetColumnCatId(),
		CalculationTypeId: row.GetCalculationTypeId(),
		CreatedAt:         row.GetCreatedAt().AsTime(),
		UpdatedAt:         row.GetUpdatedAt().AsTime(),
		IsDeleted:         row.GetIsDeleted(),
		UserId:            row.GetUserId(),
	}
}

func toFollowingCalculationSlice(rows []*tablesv1.FollowingCalculation) []openapiv1.FollowingCalculation {
	data := make([]openapiv1.FollowingCalculation, 0, len(rows))
	for _, row := range rows {
		data = append(data, toFollowingCalculation(row))
	}

	return data
}

func toFollowingCalculationPage(rows []*tablesv1.FollowingCalculation, p *tablesv1.Pagination) openapiv1.FollowingCalculationPage {
	return openapiv1.FollowingCalculationPage{
		Data: toFollowingCalculationSlice(rows),
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
