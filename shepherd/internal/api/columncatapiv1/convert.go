package columncatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toColumnCat(row *tablesv1.ColumnCat) openapiv1.ColumnCat {
	return openapiv1.ColumnCat{
		Id:                row.GetId(),
		TableId:           row.GetTableId(),
		Name:              row.GetName(),
		AliasId:           row.GetAliasId(),
		ColumnTypeId:      row.GetColumnTypeId(),
		Description:       row.GetDescription(),
		CalculationTypeId: row.GetCalculationTypeId(),
		IsDeleted:         row.GetIsDeleted(),
		ShowInUi:          row.GetShowInUi(),
		CreatedAt:         row.GetCreatedAt().AsTime(),
		UpdatedAt:         row.GetUpdatedAt().AsTime(),
		UserId:            row.GetUserId(),
	}
}

func toColumnCatSlice(rows []*tablesv1.ColumnCat) []openapiv1.ColumnCat {
	data := make([]openapiv1.ColumnCat, 0, len(rows))
	for _, row := range rows {
		data = append(data, toColumnCat(row))
	}

	return data
}

func toColumnCatPage(rows []*tablesv1.ColumnCat, p *tablesv1.Pagination) openapiv1.ColumnCatPage {
	return openapiv1.ColumnCatPage{
		Data: toColumnCatSlice(rows),
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
