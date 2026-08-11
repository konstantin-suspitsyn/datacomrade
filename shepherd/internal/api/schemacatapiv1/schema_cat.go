// Package schemacatapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.schema_cat (/v1/schema-cats...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.SchemaCatServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package schemacatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.schema_cat.
type Handler struct {
	client tablesv1.SchemaCatServiceClient
}

func New(client tablesv1.SchemaCatServiceClient) *Handler {
	return &Handler{client: client}
}

// ListSchemaCats отдаёт страницу dc.schema_cat. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListSchemaCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListSchemaCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetSchemaCatsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetSchemaCats(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toSchemaCatPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedSchemaCats отдаёт страницу мягко удалённых строк dc.schema_cat.
func (h *Handler) ListDeletedSchemaCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedSchemaCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetSchemaCatsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetSchemaCatsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toSchemaCatPage(resp.GetData(), resp.GetPagination()))
}

// GetSchemaCatById отдаёт активную строку dc.schema_cat по id.
func (h *Handler) GetSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	resp, err := h.client.GetSchemaCatById(r.Context(), &tablesv1.GetSchemaCatByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toSchemaCat(resp.GetSchemaCat()))
}

// CreateSchemaCat создаёт строку dc.schema_cat. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateSchemaCat(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.SchemaCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateSchemaCat(r.Context(), &tablesv1.CreateSchemaCatRequest{
		DatabaseId: body.DatabaseId,
		Name:       body.Name,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toSchemaCat(resp.GetSchemaCat()))
}

// UpdateSchemaCatById обновляет строку dc.schema_cat.
func (h *Handler) UpdateSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.SchemaCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateSchemaCatById(r.Context(), &tablesv1.UpdateSchemaCatByIdRequest{
		Id:         id,
		DatabaseId: body.DatabaseId,
		Name:       body.Name,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteSchemaCatById мягко удаляет строку dc.schema_cat.
func (h *Handler) DeleteSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteSchemaCatById(r.Context(), &tablesv1.DeleteSchemaCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteSchemaCatById восстанавливает мягко удалённую строку dc.schema_cat.
func (h *Handler) UndeleteSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteSchemaCatById(r.Context(), &tablesv1.UndeleteSchemaCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
