// Package grouplevelsapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.group_levels (/v1/group-levels...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.GroupLevelsServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package grouplevelsapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.group_levels.
type Handler struct {
	client tablesv1.GroupLevelsServiceClient
}

func New(client tablesv1.GroupLevelsServiceClient) *Handler {
	return &Handler{client: client}
}

// ListGroupLevels отдаёт страницу dc.group_levels. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.ListGroupLevelsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetGroupLevelsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetGroupLevels(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toGroupLevelsPage(resp.GetData(), resp.GetPagination()))
}

// SearchGroupLevels отдаёт страницу dc.group_levels, отфильтрованную по подстроке имени.
func (h *Handler) SearchGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.SearchGroupLevelsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetGroupLevelsSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetGroupLevelsSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toGroupLevelsPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedGroupLevels отдаёт список мягко удалённых строк dc.group_levels.
func (h *Handler) ListDeletedGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedGroupLevelsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetGroupLevelsDeletedRequest{
		Order: order,
	}

	resp, err := h.client.GetGroupLevelsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toGroupLevelsSlice(resp.GetRows()))
}

// FilterGroupLevels отдаёт страницу dc.group_levels с необязательными фильтрами. Непереданный
// фильтр уходит в gRPC как отсутствующий (proto optional) и не сужает выборку.
func (h *Handler) FilterGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.FilterGroupLevelsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetGroupLevelsFilteredRequest{
		ColumnId:       params.ColumnId,
		ParentColumnId: params.ParentColumnId,
		Order:          order,
		PageLimit:      pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetGroupLevelsFiltered(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toGroupLevelsPage(resp.GetData(), resp.GetPagination()))
}

// GetGroupLevelsById отдаёт активную строку dc.group_levels по id.
func (h *Handler) GetGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	resp, err := h.client.GetGroupLevelsById(r.Context(), &tablesv1.GetGroupLevelsByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toGroupLevels(resp.GetGroupLevels()))
}

// CreateGroupLevels создаёт строку dc.group_levels. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateGroupLevels(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.GroupLevelsWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateGroupLevels(r.Context(), &tablesv1.CreateGroupLevelsRequest{
		ColumnId:       body.ColumnId,
		ParentColumnId: body.ParentColumnId,
		Level:          body.Level,
		Description:    body.Description,
		ExternalId:     claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toGroupLevels(resp.GetGroupLevels()))
}

// UpdateGroupLevelsById обновляет строку dc.group_levels.
func (h *Handler) UpdateGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.GroupLevelsWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateGroupLevelsById(r.Context(), &tablesv1.UpdateGroupLevelsByIdRequest{
		Id:             id,
		ColumnId:       body.ColumnId,
		ParentColumnId: body.ParentColumnId,
		Level:          body.Level,
		Description:    body.Description,
		ExternalId:     claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteGroupLevelsById мягко удаляет строку dc.group_levels.
func (h *Handler) DeleteGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteGroupLevelsById(r.Context(), &tablesv1.DeleteGroupLevelsByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteGroupLevelsById восстанавливает мягко удалённую строку dc.group_levels.
func (h *Handler) UndeleteGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteGroupLevelsById(r.Context(), &tablesv1.UndeleteGroupLevelsByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
