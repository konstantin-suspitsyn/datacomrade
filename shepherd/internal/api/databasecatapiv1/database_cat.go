// Package databasecatapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.database_cat (/v1/database-cats...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.DatabaseCatServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package databasecatapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.database_cat.
type Handler struct {
	client tablesv1.DatabaseCatServiceClient
}

func New(client tablesv1.DatabaseCatServiceClient) *Handler {
	return &Handler{client: client}
}

// ListDatabaseCats отдаёт страницу dc.database_cat. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListDatabaseCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseCatsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCats(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCatPage(resp.GetData(), resp.GetPagination()))
}

// SearchDatabaseCats отдаёт страницу dc.database_cat, отфильтрованную по подстроке имени.
func (h *Handler) SearchDatabaseCats(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDatabaseCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseCatsSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCatsSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCatPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedDatabaseCats отдаёт страницу мягко удалённых строк dc.database_cat.
func (h *Handler) ListDeletedDatabaseCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDatabaseCatsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseCatsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCatsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCatPage(resp.GetData(), resp.GetPagination()))
}

// ListDatabaseCatsByHostAndDatabaseType отдаёт страницу dc.database_cat с
// фильтром по host_id и database_type_id. Оба фильтра опциональны и
// независимы — непереданный уходит в gRPC как отсутствующий (proto optional)
// и не сужает выборку.
func (h *Handler) ListDatabaseCatsByHostAndDatabaseType(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseCatsByHostAndDatabaseTypeParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	// Пустая строка — это «сортировка по умолчанию» (id), её datacatalogue
	// принимает наравне с name; проверяем только явно переданное значение.
	var orderBy string
	if params.OrderBy != nil && *params.OrderBy != "" {
		if !params.OrderBy.Valid() {
			apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order_by допускает только name")
			return
		}

		orderBy = string(*params.OrderBy)
	}

	req := &tablesv1.GetDatabaseCatsByHostIdDatabaseTypeIdRequest{
		HostId:         params.HostId,
		DatabaseTypeId: params.DatabaseTypeId,
		Order:          order,
		OrderBy:        orderBy,
		PageLimit:      pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCatsByHostIdDatabaseTypeId(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCatPage(resp.GetData(), resp.GetPagination()))
}

// GetDatabaseCatById отдаёт активную строку dc.database_cat по id.
func (h *Handler) GetDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	resp, err := h.client.GetDatabaseCatById(r.Context(), &tablesv1.GetDatabaseCatByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCat(resp.GetDatabaseCat()))
}

// CreateDatabaseCat создаёт строку dc.database_cat. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateDatabaseCat(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DatabaseCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateDatabaseCat(r.Context(), &tablesv1.CreateDatabaseCatRequest{
		Name:           body.Name,
		HostId:         body.HostId,
		DatabaseTypeId: body.DatabaseTypeId,
		Description:    body.Description,
		ExternalId:     claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toDatabaseCat(resp.GetDatabaseCat()))
}

// UpdateDatabaseCatById обновляет строку dc.database_cat.
func (h *Handler) UpdateDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DatabaseCatWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateDatabaseCatById(r.Context(), &tablesv1.UpdateDatabaseCatByIdRequest{
		Id:             id,
		Name:           body.Name,
		HostId:         body.HostId,
		DatabaseTypeId: body.DatabaseTypeId,
		Description:    body.Description,
		ExternalId:     claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteDatabaseCatById мягко удаляет строку dc.database_cat.
func (h *Handler) DeleteDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteDatabaseCatById(r.Context(), &tablesv1.DeleteDatabaseCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteDatabaseCatById восстанавливает мягко удалённую строку dc.database_cat.
func (h *Handler) UndeleteDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteDatabaseCatById(r.Context(), &tablesv1.UndeleteDatabaseCatByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
