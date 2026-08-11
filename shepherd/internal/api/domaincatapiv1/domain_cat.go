// Package domaincatapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.domain_cat (/v1/domain-cats...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.DomainCatServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package domaincatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.domain_cat.
type Handler struct {
	client tablesv1.DomainCatServiceClient
}

func New(client tablesv1.DomainCatServiceClient) *Handler {
	return &Handler{client: client}
}

// ListDomainCats отдаёт страницу dc.domain_cat. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListDomainCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDomainCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDomainCatsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDomainCats(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDomainCatPage(resp.GetData(), resp.GetPagination()))
}

// SearchDomainCats отдаёт страницу dc.domain_cat, отфильтрованную по подстроке имени.
func (h *Handler) SearchDomainCats(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDomainCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDomainCatsSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDomainCatsSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDomainCatPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedDomainCats отдаёт страницу мягко удалённых строк dc.domain_cat.
func (h *Handler) ListDeletedDomainCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDomainCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDomainCatsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDomainCatsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDomainCatPage(resp.GetData(), resp.GetPagination()))
}

// GetDomainCatById отдаёт активную строку dc.domain_cat по id.
func (h *Handler) GetDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	resp, err := h.client.GetDomainCatById(r.Context(), &tablesv1.GetDomainCatByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDomainCat(resp.GetDomainCat()))
}

// CreateDomainCat создаёт строку dc.domain_cat. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateDomainCat(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DomainCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateDomainCat(r.Context(), &tablesv1.CreateDomainCatRequest{
		DomainName: body.DomainName,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toDomainCat(resp.GetDomainCat()))
}

// UpdateDomainCatById обновляет строку dc.domain_cat.
func (h *Handler) UpdateDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DomainCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateDomainCatById(r.Context(), &tablesv1.UpdateDomainCatByIdRequest{
		Id:         id,
		DomainName: body.DomainName,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteDomainCatById мягко удаляет строку dc.domain_cat.
func (h *Handler) DeleteDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteDomainCatById(r.Context(), &tablesv1.DeleteDomainCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteDomainCatById восстанавливает мягко удалённую строку dc.domain_cat.
func (h *Handler) UndeleteDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteDomainCatById(r.Context(), &tablesv1.UndeleteDomainCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
