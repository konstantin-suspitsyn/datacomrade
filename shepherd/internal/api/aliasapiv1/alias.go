// Package aliasapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.alias (/v1/aliases...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.AliasServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package aliasapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.alias.
type Handler struct {
	client tablesv1.AliasServiceClient
}

func New(client tablesv1.AliasServiceClient) *Handler {
	return &Handler{client: client}
}

// ListAliases отдаёт страницу dc.alias. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListAliases(w http.ResponseWriter, r *http.Request, params openapiv1.ListAliasesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetAliasesRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetAliases(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAliasPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedAliases отдаёт страницу мягко удалённых строк dc.alias.
func (h *Handler) ListDeletedAliases(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedAliasesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetAliasesDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetAliasesDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAliasPage(resp.GetData(), resp.GetPagination()))
}

// GetAliasById отдаёт активную строку dc.alias по id.
func (h *Handler) GetAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	resp, err := h.client.GetAliasById(r.Context(), &tablesv1.GetAliasByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAlias(resp.GetAlias()))
}

// CreateAlias создаёт строку dc.alias. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateAlias(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.AliasWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateAlias(r.Context(), &tablesv1.CreateAliasRequest{
		Name:        body.Name,
		Description: body.Description,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toAlias(resp.GetAlias()))
}

// UpdateAliasById обновляет строку dc.alias.
func (h *Handler) UpdateAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.AliasUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateAliasById(r.Context(), &tablesv1.UpdateAliasByIdRequest{
		Id:          id,
		Name:        body.Name,
		Description: body.Description,
		IsDeleted:   body.IsDeleted,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAliasById мягко удаляет строку dc.alias.
func (h *Handler) DeleteAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteAliasById(r.Context(), &tablesv1.DeleteAliasByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteAliasById восстанавливает мягко удалённую строку dc.alias.
func (h *Handler) UndeleteAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteAliasById(r.Context(), &tablesv1.UndeleteAliasByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
