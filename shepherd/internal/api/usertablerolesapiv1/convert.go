package usertablerolesapiv1

import (
	"encoding/json"
	"net/http"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

func toUserTableRoles(row *userdomainrolesv1.UserTableRoles) openapiv1.UserTableRoles {
	return openapiv1.UserTableRoles{
		Id:           row.GetId(),
		UserId:       row.GetUserId(),
		TableRolesId: row.GetTableRolesId(),
		CreatedAt:    row.GetCreatedAt().AsTime(),
		UpdatedAt:    row.GetUpdatedAt().AsTime(),
		IsDeleted:    row.GetIsDeleted(),
		TableId:      row.GetTableId(),
		UpdatedById:  row.GetUpdatedById(),
	}
}

// toUserTableRolesRow: table_roles_id приходит из gRPC как table_roles_id_2 —
// имя от sqlc, которое джойн даёт второй одноимённой колонке; наружу отдаём
// его под обычным именем.
func toUserTableRolesRow(row *userdomainrolesv1.GetUserTableRolesesRow) openapiv1.UserTableRolesRow {
	return openapiv1.UserTableRolesRow{
		Id:                    row.GetId(),
		UserId:                row.GetUserId(),
		TableRolesId:          row.GetTableRolesId_2(),
		CreatedAt:             row.GetCreatedAt().AsTime(),
		UpdatedAt:             row.GetUpdatedAt().AsTime(),
		TableId:               row.GetTableId(),
		DcUserId:              row.DcUserId,
		DcUserName:            row.GetDcUserName(),
		DcUserExternalId:      row.GetDcUserExternalId(),
		TableRolesName:        row.GetTableRolesName(),
		TableRolesDescription: row.GetTableRolesDescription(),
		TableCatId:            row.TableCatId,
		TableCatName:          row.GetTableCatName(),
	}
}

func toUserTableRolesPage(rows []*userdomainrolesv1.GetUserTableRolesesRow, p *userdomainrolesv1.Pagination) openapiv1.UserTableRolesPage {
	data := make([]openapiv1.UserTableRolesRow, 0, len(rows))
	for _, row := range rows {
		data = append(data, toUserTableRolesRow(row))
	}

	return openapiv1.UserTableRolesPage{
		Data:       data,
		Pagination: toPagination(p),
	}
}

func toUserTableRolesFilteredRow(row *userdomainrolesv1.GetUserTableRolesesFilteredRow) openapiv1.UserTableRolesFilteredRow {
	return openapiv1.UserTableRolesFilteredRow{
		Id:                    row.GetId(),
		UserId:                row.GetUserId(),
		CreatedAt:             row.GetCreatedAt().AsTime(),
		UpdatedAt:             row.GetUpdatedAt().AsTime(),
		TableId:               row.GetTableId(),
		DcUserId:              row.DcUserId,
		DcUserName:            row.GetDcUserName(),
		DcUserExternalId:      row.GetDcUserExternalId(),
		TableRolesName:        row.GetTableRolesName(),
		TableRolesDescription: row.GetTableRolesDescription(),
		TableCatId:            row.TableCatId,
		TableCatName:          row.GetTableCatName(),
	}
}

func toUserTableRolesFilteredPage(rows []*userdomainrolesv1.GetUserTableRolesesFilteredRow, p *userdomainrolesv1.Pagination) openapiv1.UserTableRolesFilteredPage {
	data := make([]openapiv1.UserTableRolesFilteredRow, 0, len(rows))
	for _, row := range rows {
		data = append(data, toUserTableRolesFilteredRow(row))
	}

	return openapiv1.UserTableRolesFilteredPage{
		Data:       data,
		Pagination: toPagination(p),
	}
}

func toUserTableRolesDetail(row *userdomainrolesv1.GetUserTableRolesByIdRow) openapiv1.UserTableRolesDetail {
	return openapiv1.UserTableRolesDetail{
		Id:                    row.GetId(),
		UserId:                row.GetUserId(),
		TableRolesId:          row.GetTableRolesId(),
		CreatedAt:             row.GetCreatedAt().AsTime(),
		UpdatedAt:             row.GetUpdatedAt().AsTime(),
		TableId:               row.GetTableId(),
		DcUserId:              row.DcUserId,
		DcUserName:            row.GetDcUserName(),
		DcUserExternalId:      row.GetDcUserExternalId(),
		TableRolesName:        row.GetTableRolesName(),
		TableRolesDescription: row.GetTableRolesDescription(),
		TableCatId:            row.TableCatId,
		TableCatName:          row.GetTableCatName(),
	}
}

func toPagination(p *userdomainrolesv1.Pagination) openapiv1.Pagination {
	return openapiv1.Pagination{
		Page:       p.GetPage(),
		PerPage:    p.GetPerPage(),
		TotalItems: p.GetTotalItems(),
		TotalPages: p.GetTotalPages(),
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
