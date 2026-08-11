// Package tabletypeapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.table_type (/v1/table-types...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.TableTypeServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package tabletypeapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.table_type.
type Handler struct {
	client tablesv1.TableTypeServiceClient
}

func New(client tablesv1.TableTypeServiceClient) *Handler {
	return &Handler{client: client}
}

// ListTableTypes отдаёт страницу dc.table_type. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListTableTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListTableTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetTableTypesRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetTableTypes(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTableTypePage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedTableTypes отдаёт страницу мягко удалённых строк dc.table_type.
func (h *Handler) ListDeletedTableTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedTableTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetTableTypesDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetTableTypesDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTableTypePage(resp.GetData(), resp.GetPagination()))
}

// GetTableTypeById отдаёт активную строку dc.table_type по id.
func (h *Handler) GetTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	resp, err := h.client.GetTableTypeById(r.Context(), &tablesv1.GetTableTypeByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTableType(resp.GetTableType()))
}

// CreateTableType создаёт строку dc.table_type. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateTableType(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.TableTypeWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateTableType(r.Context(), &tablesv1.CreateTableTypeRequest{
		Name:        body.Name,
		Description: body.Description,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toTableType(resp.GetTableType()))
}

// UpdateTableTypeById обновляет строку dc.table_type.
func (h *Handler) UpdateTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.TableTypeWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateTableTypeById(r.Context(), &tablesv1.UpdateTableTypeByIdRequest{
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

// DeleteTableTypeById мягко удаляет строку dc.table_type.
func (h *Handler) DeleteTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteTableTypeById(r.Context(), &tablesv1.DeleteTableTypeByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteTableTypeById восстанавливает мягко удалённую строку dc.table_type.
func (h *Handler) UndeleteTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteTableTypeById(r.Context(), &tablesv1.UndeleteTableTypeByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
