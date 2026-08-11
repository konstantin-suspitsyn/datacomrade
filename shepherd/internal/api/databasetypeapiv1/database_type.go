// Package databasetypeapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.database_type (/v1/database-types...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.DatabaseTypeServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package databasetypeapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.database_type.
type Handler struct {
	client tablesv1.DatabaseTypeServiceClient
}

func New(client tablesv1.DatabaseTypeServiceClient) *Handler {
	return &Handler{client: client}
}

// ListDatabaseTypes отдаёт страницу dc.database_type. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListDatabaseTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseTypesRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseTypes(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseTypePage(resp.GetData(), resp.GetPagination()))
}

// SearchDatabaseTypes отдаёт страницу dc.database_type, отфильтрованную по подстроке имени.
func (h *Handler) SearchDatabaseTypes(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDatabaseTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseTypesSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseTypesSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseTypePage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedDatabaseTypes отдаёт страницу мягко удалённых строк dc.database_type.
func (h *Handler) ListDeletedDatabaseTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDatabaseTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseTypesDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseTypesDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseTypePage(resp.GetData(), resp.GetPagination()))
}

// GetDatabaseTypeById отдаёт активную строку dc.database_type по id.
func (h *Handler) GetDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	resp, err := h.client.GetDatabaseTypeById(r.Context(), &tablesv1.GetDatabaseTypeByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseType(resp.GetDatabaseType()))
}

// CreateDatabaseType создаёт строку dc.database_type. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateDatabaseType(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DatabaseTypeWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateDatabaseType(r.Context(), &tablesv1.CreateDatabaseTypeRequest{
		Name:       body.Name,
		DbVersion:  body.DbVersion,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toDatabaseType(resp.GetDatabaseType()))
}

// UpdateDatabaseTypeById обновляет строку dc.database_type.
func (h *Handler) UpdateDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DatabaseTypeWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateDatabaseTypeById(r.Context(), &tablesv1.UpdateDatabaseTypeByIdRequest{
		Id:         id,
		Name:       body.Name,
		DbVersion:  body.DbVersion,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteDatabaseTypeById мягко удаляет строку dc.database_type.
func (h *Handler) DeleteDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteDatabaseTypeById(r.Context(), &tablesv1.DeleteDatabaseTypeByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteDatabaseTypeById восстанавливает мягко удалённую строку dc.database_type.
func (h *Handler) UndeleteDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteDatabaseTypeById(r.Context(), &tablesv1.UndeleteDatabaseTypeByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
