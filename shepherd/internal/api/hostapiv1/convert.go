package hostapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toHost(row *tablesv1.Host) openapiv1.Host {
	return openapiv1.Host{
		Id:          row.GetId(),
		Name:        row.GetName(),
		Description: row.GetDescription(),
		HostEnv:     row.GetHostEnv(),
		PortEnv:     row.GetPortEnv(),
		UsernameEnv: row.GetUsernameEnv(),
		PasswordEnv: row.GetPasswordEnv(),
		IsDeleted:   row.GetIsDeleted(),
		CreatedAt:   row.GetCreatedAt().AsTime(),
		UpdatedAt:   row.GetUpdatedAt().AsTime(),
		UserId:      row.GetUserId(),
	}
}

func toHostPage(rows []*tablesv1.Host, p *tablesv1.Pagination) openapiv1.HostPage {
	data := make([]openapiv1.Host, 0, len(rows))
	for _, row := range rows {
		data = append(data, toHost(row))
	}

	return openapiv1.HostPage{
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
