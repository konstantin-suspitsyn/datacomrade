// Package hastogroupapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.has_to_group (/v1/has-to-groups...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.HasToGroupServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package hastogroupapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.has_to_group.
type Handler struct {
	client tablesv1.HasToGroupServiceClient
}

func New(client tablesv1.HasToGroupServiceClient) *Handler {
	return &Handler{client: client}
}

// ListHasToGroups отдаёт страницу dc.has_to_group. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListHasToGroups(w http.ResponseWriter, r *http.Request, params openapiv1.ListHasToGroupsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetHasToGroupsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetHasToGroups(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHasToGroupPage(resp.GetData(), resp.GetPagination()))
}

// SearchHasToGroups отдаёт страницу dc.has_to_group, отфильтрованную по подстроке имени.
func (h *Handler) SearchHasToGroups(w http.ResponseWriter, r *http.Request, params openapiv1.SearchHasToGroupsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetHasToGroupsSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetHasToGroupsSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHasToGroupPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedHasToGroups отдаёт страницу мягко удалённых строк dc.has_to_group.
func (h *Handler) ListDeletedHasToGroups(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedHasToGroupsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetHasToGroupsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetHasToGroupsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHasToGroupPage(resp.GetData(), resp.GetPagination()))
}

// GetHasToGroupById отдаёт активную строку dc.has_to_group по id.
func (h *Handler) GetHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	resp, err := h.client.GetHasToGroupById(r.Context(), &tablesv1.GetHasToGroupByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHasToGroup(resp.GetHasToGroup()))
}

// CreateHasToGroup создаёт строку dc.has_to_group. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateHasToGroup(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.HasToGroupWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateHasToGroup(r.Context(), &tablesv1.CreateHasToGroupRequest{
		ColumnIdA:   body.ColumnIdA,
		ColumnIdB:   body.ColumnIdB,
		Description: body.Description,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toHasToGroup(resp.GetHasToGroup()))
}

// UpdateHasToGroupById обновляет строку dc.has_to_group.
func (h *Handler) UpdateHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.HasToGroupWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateHasToGroupById(r.Context(), &tablesv1.UpdateHasToGroupByIdRequest{
		Id:          id,
		ColumnIdA:   body.ColumnIdA,
		ColumnIdB:   body.ColumnIdB,
		Description: body.Description,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteHasToGroupById мягко удаляет строку dc.has_to_group.
func (h *Handler) DeleteHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteHasToGroupById(r.Context(), &tablesv1.DeleteHasToGroupByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteHasToGroupById восстанавливает мягко удалённую строку dc.has_to_group.
func (h *Handler) UndeleteHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteHasToGroupById(r.Context(), &tablesv1.UndeleteHasToGroupByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
