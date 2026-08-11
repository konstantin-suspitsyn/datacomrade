// Package columntypeapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.column_type (/v1/column-types...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.ColumnTypeServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package columntypeapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.column_type.
type Handler struct {
	client tablesv1.ColumnTypeServiceClient
}

func New(client tablesv1.ColumnTypeServiceClient) *Handler {
	return &Handler{client: client}
}

// ListColumnTypes отдаёт страницу dc.column_type. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListColumnTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListColumnTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetColumnTypesRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetColumnTypes(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnTypePage(resp.GetData(), resp.GetPagination()))
}

// SearchColumnTypes отдаёт страницу dc.column_type, отфильтрованную по подстроке имени.
func (h *Handler) SearchColumnTypes(w http.ResponseWriter, r *http.Request, params openapiv1.SearchColumnTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetColumnTypesSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetColumnTypesSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnTypePage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedColumnTypes отдаёт страницу мягко удалённых строк dc.column_type.
func (h *Handler) ListDeletedColumnTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedColumnTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetColumnTypesDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetColumnTypesDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnTypePage(resp.GetData(), resp.GetPagination()))
}

// GetColumnTypeById отдаёт активную строку dc.column_type по id.
func (h *Handler) GetColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	resp, err := h.client.GetColumnTypeById(r.Context(), &tablesv1.GetColumnTypeByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnType(resp.GetColumnType()))
}

// CreateColumnType создаёт строку dc.column_type. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateColumnType(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.ColumnTypeWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateColumnType(r.Context(), &tablesv1.CreateColumnTypeRequest{
		Name:        body.Name,
		Description: body.Description,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toColumnType(resp.GetColumnType()))
}

// UpdateColumnTypeById обновляет строку dc.column_type.
func (h *Handler) UpdateColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.ColumnTypeWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateColumnTypeById(r.Context(), &tablesv1.UpdateColumnTypeByIdRequest{
		Id:          id,
		Name:        body.Name,
		Description: body.Description,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteColumnTypeById мягко удаляет строку dc.column_type.
func (h *Handler) DeleteColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteColumnTypeById(r.Context(), &tablesv1.DeleteColumnTypeByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteColumnTypeById восстанавливает мягко удалённую строку dc.column_type.
func (h *Handler) UndeleteColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteColumnTypeById(r.Context(), &tablesv1.UndeleteColumnTypeByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
