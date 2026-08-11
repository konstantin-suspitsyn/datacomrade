// Package columncatapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.column_cat (/v1/column-cats...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.ColumnCatServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package columncatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.column_cat.
type Handler struct {
	client tablesv1.ColumnCatServiceClient
}

func New(client tablesv1.ColumnCatServiceClient) *Handler {
	return &Handler{client: client}
}

// ListColumnCats отдаёт страницу dc.column_cat. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListColumnCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetColumnCatsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetColumnCats(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnCatPage(resp.GetData(), resp.GetPagination()))
}

// SearchColumnCats отдаёт страницу dc.column_cat, отфильтрованную по подстроке имени.
func (h *Handler) SearchColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.SearchColumnCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetColumnCatsSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetColumnCatsSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnCatPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedColumnCats отдаёт страницу мягко удалённых строк dc.column_cat.
func (h *Handler) ListDeletedColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedColumnCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetColumnCatsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetColumnCatsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnCatPage(resp.GetData(), resp.GetPagination()))
}

// FilterColumnCats отдаёт страницу dc.column_cat с необязательными фильтрами. Непереданный
// фильтр уходит в gRPC как отсутствующий (proto optional) и не сужает выборку.
func (h *Handler) FilterColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.FilterColumnCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetColumnCatsOptionalFiltersRequest{
		TableId:           params.TableId,
		AliasId:           params.AliasId,
		ColumnTypeId:      params.ColumnTypeId,
		CalculationTypeId: params.CalculationTypeId,
		ShowInUi:          params.ShowInUi,
		Order:             order,
		PageLimit:         pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetColumnCatsOptionalFilters(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnCatPage(resp.GetData(), resp.GetPagination()))
}

// GetColumnCatById отдаёт активную строку dc.column_cat по id.
func (h *Handler) GetColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	resp, err := h.client.GetColumnCatById(r.Context(), &tablesv1.GetColumnCatByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toColumnCat(resp.GetColumnCat()))
}

// CreateColumnCat создаёт строку dc.column_cat. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateColumnCat(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.ColumnCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateColumnCat(r.Context(), &tablesv1.CreateColumnCatRequest{
		TableId:           body.TableId,
		Name:              body.Name,
		AliasId:           body.AliasId,
		ColumnTypeId:      body.ColumnTypeId,
		Description:       body.Description,
		CalculationTypeId: body.CalculationTypeId,
		ShowInUi:          body.ShowInUi,
		ExternalId:        claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toColumnCat(resp.GetColumnCat()))
}

// UpdateColumnCatById обновляет строку dc.column_cat.
func (h *Handler) UpdateColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.ColumnCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateColumnCatById(r.Context(), &tablesv1.UpdateColumnCatByIdRequest{
		Id:                id,
		TableId:           body.TableId,
		Name:              body.Name,
		AliasId:           body.AliasId,
		ColumnTypeId:      body.ColumnTypeId,
		Description:       body.Description,
		CalculationTypeId: body.CalculationTypeId,
		ShowInUi:          body.ShowInUi,
		ExternalId:        claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteColumnCatById мягко удаляет строку dc.column_cat.
func (h *Handler) DeleteColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteColumnCatById(r.Context(), &tablesv1.DeleteColumnCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteColumnCatById восстанавливает мягко удалённую строку dc.column_cat.
func (h *Handler) UndeleteColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteColumnCatById(r.Context(), &tablesv1.UndeleteColumnCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
