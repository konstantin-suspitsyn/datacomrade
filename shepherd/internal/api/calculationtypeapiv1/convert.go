package calculationtypeapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toCalculationType(row *tablesv1.CalculationType) openapiv1.CalculationType {
	return openapiv1.CalculationType{
		Id:          row.GetId(),
		Name:        row.GetName(),
		Description: row.GetDescription(),
		CreatedAt:   row.GetCreatedAt().AsTime(),
		UpdatedAt:   row.GetUpdatedAt().AsTime(),
		IsDeleted:   row.GetIsDeleted(),
	}
}

func toCalculationTypeSlice(rows []*tablesv1.CalculationType) []openapiv1.CalculationType {
	data := make([]openapiv1.CalculationType, 0, len(rows))
	for _, row := range rows {
		data = append(data, toCalculationType(row))
	}

	return data
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
