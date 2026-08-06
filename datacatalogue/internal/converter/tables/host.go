package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// HostToProto переводит строку dc.host в сущность gRPC.
func HostToProto(row tables_model.DcHost) *tablesv1.Host {
	return &tablesv1.Host{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		HostEnv:     row.HostEnv,
		PortEnv:     row.PortEnv,
		UsernameEnv: row.UsernameEnv,
		PasswordEnv: row.PasswordEnv,
		IsDeleted:   row.IsDeleted,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		UserId:      row.UserID,
	}
}

// HostsToProto переводит список строк dc.host в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func HostsToProto(rows []tables_model.DcHost) []*tablesv1.Host {
	items := make([]*tablesv1.Host, 0, len(rows))

	for _, row := range rows {
		items = append(items, HostToProto(row))
	}

	return items
}

// ToCreateHostParams собирает параметры вставки dc.host из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateHostParams(req *tablesv1.CreateHostRequest) tables_model.CreateHostParams {
	return tables_model.CreateHostParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		HostEnv:     req.GetHostEnv(),
		PortEnv:     req.GetPortEnv(),
		UsernameEnv: req.GetUsernameEnv(),
		PasswordEnv: req.GetPasswordEnv(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateHostByIdParams собирает параметры обновления dc.host из запроса gRPC.
func ToUpdateHostByIdParams(req *tablesv1.UpdateHostByIdRequest) tables_model.UpdateHostByIdParams {
	return tables_model.UpdateHostByIdParams{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		HostEnv:     req.GetHostEnv(),
		PortEnv:     req.GetPortEnv(),
		UsernameEnv: req.GetUsernameEnv(),
		PasswordEnv: req.GetPasswordEnv(),
		ExternalID:  converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToDeleteHostByIdParams собирает параметры мягкого удаления dc.host из запроса gRPC.
func ToDeleteHostByIdParams(req *tablesv1.DeleteHostByIdRequest) tables_model.DeleteHostByIdParams {
	return tables_model.DeleteHostByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToUndeleteHostByIdParams собирает параметры обратного удаления dc.host из запроса gRPC.
func ToUndeleteHostByIdParams(req *tablesv1.UndeleteHostByIdRequest) tables_model.UndeleteHostByIdParams {
	return tables_model.UndeleteHostByIdParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		ID:         req.GetId(),
	}
}

// ToGetHostsParams собирает параметры страницы dc.host из запроса gRPC.
func ToGetHostsParams(req *tablesv1.GetHostsRequest) tables_model.GetHostsParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetHostsParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}

// ToGetHostsSearchNameParams собирает параметры страницы dc.host из запроса gRPC.
func ToGetHostsSearchNameParams(req *tablesv1.GetHostsSearchNameRequest) tables_model.GetHostsSearchNameParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetHostsSearchNameParams{
		Order:      req.GetOrder(),
		Page:       page,
		PageLimit:  limit,
		SearchName: req.GetSearchName(),
	}
}

// ToGetHostDeletedParams собирает параметры страницы dc.host из запроса gRPC.
func ToGetHostDeletedParams(req *tablesv1.GetHostDeletedRequest) tables_model.GetHostDeletedParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetHostDeletedParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
