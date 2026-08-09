// Package hostapiv1 реализует часть openapiv1.ServerInterface, отвечающую за
// dc.host (/v1/hosts...) — REST-хендлеры вручную поверх сгенерированного
// gRPC-клиента tablesv1.HostServiceClient (Metadata Service), см. API Gateway.md.
// Доступ гейтится ролью admin/maintainer в shepherd/cmd/main.go
// (rbac.RequireRole), а не в этом пакете.
package hostapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.host.
type Handler struct {
	client tablesv1.HostServiceClient
}

func New(client tablesv1.HostServiceClient) *Handler {
	return &Handler{client: client}
}

// ListHosts отдаёт страницу dc.host. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListHosts(w http.ResponseWriter, r *http.Request, params openapiv1.ListHostsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetHostsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetHosts(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHostPage(resp.GetData(), resp.GetPagination()))
}

// SearchHosts отдаёт страницу dc.host, отфильтрованную по подстроке имени.
func (h *Handler) SearchHosts(w http.ResponseWriter, r *http.Request, params openapiv1.SearchHostsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetHostsSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetHostsSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHostPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedHosts отдаёт страницу мягко удалённых строк dc.host.
func (h *Handler) ListDeletedHosts(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedHostsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetHostDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetHostDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHostPage(resp.GetData(), resp.GetPagination()))
}

// GetHostById отдаёт активную строку dc.host по id.
func (h *Handler) GetHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	resp, err := h.client.GetHostById(r.Context(), &tablesv1.GetHostByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toHost(resp.GetHost()))
}

// CreateHost создаёт строку dc.host. external_id (автор) берётся из claims
// токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateHost(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.HostWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateHost(r.Context(), &tablesv1.CreateHostRequest{
		Name:        body.Name,
		Description: body.Description,
		HostEnv:     body.HostEnv,
		PortEnv:     body.PortEnv,
		UsernameEnv: body.UsernameEnv,
		PasswordEnv: body.PasswordEnv,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toHost(resp.GetHost()))
}

// UpdateHostById обновляет строку dc.host.
func (h *Handler) UpdateHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.HostWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateHostById(r.Context(), &tablesv1.UpdateHostByIdRequest{
		Id:          id,
		Name:        body.Name,
		Description: body.Description,
		HostEnv:     body.HostEnv,
		PortEnv:     body.PortEnv,
		UsernameEnv: body.UsernameEnv,
		PasswordEnv: body.PasswordEnv,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteHostById мягко удаляет строку dc.host.
func (h *Handler) DeleteHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteHostById(r.Context(), &tablesv1.DeleteHostByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteHostById восстанавливает мягко удалённую строку dc.host.
func (h *Handler) UndeleteHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteHostById(r.Context(), &tablesv1.UndeleteHostByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
