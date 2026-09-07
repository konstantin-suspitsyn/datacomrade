package columncatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toColumnCat(row *tablesv1.ColumnCat) openapiv1.ColumnCat {
	// alias_id и calculation_type_id читаются полем, а не геттером: колонка
	// может быть без алиаса и без типа расчёта, а GetXxx() отдал бы 0 и на
	// «значения нет», и на значение с id 0. Прямой доступ к полю не переживает
	// nil-row, который геттеры гасят сами.
	if row == nil {
		row = &tablesv1.ColumnCat{}
	}

	return openapiv1.ColumnCat{
		Id:                row.GetId(),
		TableId:           row.GetTableId(),
		Name:              row.GetName(),
		AliasId:           row.AliasId,
		ColumnTypeId:      row.GetColumnTypeId(),
		Description:       row.GetDescription(),
		CalculationTypeId: row.CalculationTypeId,
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
