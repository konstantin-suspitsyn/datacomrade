package userdomainrolesapiv1

import (
	"encoding/json"
	"net/http"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toUserDomainRoles(row *userdomainrolesv1.UserDomainRoles) openapiv1.UserDomainRoles {
	return openapiv1.UserDomainRoles{
		Id:            row.GetId(),
		UserId:        row.GetUserId(),
		DomainRolesId: row.GetDomainRolesId(),
		CreatedAt:     row.GetCreatedAt().AsTime(),
		UpdatedAt:     row.GetUpdatedAt().AsTime(),
		IsDeleted:     row.GetIsDeleted(),
		DomainId:      row.GetDomainId(),
		UpdatedById:   row.GetUpdatedById(),
	}
}

func toUserDomainRolesRow(row *userdomainrolesv1.GetUserDomainRolesesRow) openapiv1.UserDomainRolesRow {
	return openapiv1.UserDomainRolesRow{
		Id:                     row.GetId(),
		UserId:                 row.GetUserId(),
		DomainRolesId:          row.GetDomainRolesId(),
		CreatedAt:              row.GetCreatedAt().AsTime(),
		UpdatedAt:              row.GetUpdatedAt().AsTime(),
		IsDeleted:              row.GetIsDeleted(),
		DomainId:               row.GetDomainId(),
		UpdatedById:            row.GetUpdatedById(),
		DomainRolesName:        row.GetDomainRolesName(),
		DomainRolesDescription: row.GetDomainRolesDescription(),
		DcUserName:             row.GetDcUserName(),
	}
}

func toUserDomainRolesSlice(rows []*userdomainrolesv1.GetUserDomainRolesesRow) []openapiv1.UserDomainRolesRow {
	data := make([]openapiv1.UserDomainRolesRow, 0, len(rows))
	for _, row := range rows {
		data = append(data, toUserDomainRolesRow(row))
	}

	return data
}

func toUserDomainRolesPage(rows []*userdomainrolesv1.GetUserDomainRolesesRow, p *userdomainrolesv1.Pagination) openapiv1.UserDomainRolesPage {
	return openapiv1.UserDomainRolesPage{
		Data: toUserDomainRolesSlice(rows),
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
