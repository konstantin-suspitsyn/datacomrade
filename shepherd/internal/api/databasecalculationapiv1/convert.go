package databasecalculationapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toDatabaseCalculation(row *tablesv1.DatabaseCalculation) openapiv1.DatabaseCalculation {
	return openapiv1.DatabaseCalculation{
		Id:                row.GetId(),
		DatabaseCatId:     row.GetDatabaseCatId(),
		CalculationTypeId: row.GetCalculationTypeId(),
		CreatedAt:         row.GetCreatedAt().AsTime(),
		UpdatedAt:         row.GetUpdatedAt().AsTime(),
		IsDeleted:         row.GetIsDeleted(),
		UserId:            row.GetUserId(),
	}
}

func toDatabaseCalculationSlice(rows []*tablesv1.DatabaseCalculation) []openapiv1.DatabaseCalculation {
	data := make([]openapiv1.DatabaseCalculation, 0, len(rows))
	for _, row := range rows {
		data = append(data, toDatabaseCalculation(row))
	}

	return data
}

func toDatabaseCalculationPage(rows []*tablesv1.DatabaseCalculation, p *tablesv1.Pagination) openapiv1.DatabaseCalculationPage {
	return openapiv1.DatabaseCalculationPage{
		Data: toDatabaseCalculationSlice(rows),
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
