// Package tablecatapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.table_cat (/v1/table-cats...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.TableCatServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package tablecatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.table_cat.
type Handler struct {
	client tablesv1.TableCatServiceClient
}

func New(client tablesv1.TableCatServiceClient) *Handler {
	return &Handler{client: client}
}

// ListTableCats отдаёт страницу dc.table_cat. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListTableCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListTableCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetTableCatsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetTableCats(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTableCatPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedTableCats отдаёт страницу мягко удалённых строк dc.table_cat.
func (h *Handler) ListDeletedTableCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedTableCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetTableCatsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetTableCatsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTableCatPage(resp.GetData(), resp.GetPagination()))
}

// GetTableCatById отдаёт активную строку dc.table_cat по id.
func (h *Handler) GetTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	resp, err := h.client.GetTableCatById(r.Context(), &tablesv1.GetTableCatByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTableCat(resp.GetTableCat()))
}

// CreateTableCat создаёт строку dc.table_cat. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateTableCat(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.TableCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateTableCat(r.Context(), &tablesv1.CreateTableCatRequest{
		Name:        body.Name,
		Description: body.Description,
		SchemaId:    body.SchemaId,
		TableTypeId: body.TableTypeId,
		DomainId:    body.DomainId,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toTableCat(resp.GetTableCat()))
}

// UpdateTableCatById обновляет строку dc.table_cat.
func (h *Handler) UpdateTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.TableCatUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateTableCatById(r.Context(), &tablesv1.UpdateTableCatByIdRequest{
		Id:          id,
		Name:        body.Name,
		Description: body.Description,
		SchemaId:    body.SchemaId,
		TableTypeId: body.TableTypeId,
		DomainId:    body.DomainId,
		IsGetDict:   body.IsGetDict,
		ExternalId:  claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTableCatById мягко удаляет строку dc.table_cat.
func (h *Handler) DeleteTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteTableCatById(r.Context(), &tablesv1.DeleteTableCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteTableCatById восстанавливает мягко удалённую строку dc.table_cat.
func (h *Handler) UndeleteTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteTableCatById(r.Context(), &tablesv1.UndeleteTableCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
