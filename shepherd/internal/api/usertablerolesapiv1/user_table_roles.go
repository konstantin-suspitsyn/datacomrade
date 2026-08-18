// Package usertablerolesapiv1 реализует часть openapiv1.ServerInterface,
// отвечающую за dc.user_table_roles (/v1/user-table-roles...) — REST-хендлеры
// вручную поверх сгенерированного gRPC-клиента
// userdomainrolesv1.UserTableRolesServiceClient (Metadata Service),
// см. API Gateway.md. Доступ гейтится ролью admin в shepherd/cmd/main.go
// (rbac.RequireRole), а не в этом пакете.
package usertablerolesapiv1

import (
	"encoding/json"
	"net/http"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.user_table_roles.
type Handler struct {
	client userdomainrolesv1.UserTableRolesServiceClient
}

func New(client userdomainrolesv1.UserTableRolesServiceClient) *Handler {
	return &Handler{client: client}
}

// ListUserTableRoles отдаёт страницу dc.user_table_roles. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListUserTableRoles(w http.ResponseWriter, r *http.Request, params openapiv1.ListUserTableRolesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &userdomainrolesv1.GetUserTableRolesesRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetUserTableRoleses(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserTableRolesPage(resp.GetData(), resp.GetPagination()))
}

// FilterUserTableRoles отдаёт страницу dc.user_table_roles с необязательными фильтрами.
// Непереданный фильтр уходит в gRPC как отсутствующий (proto optional) и не сужает выборку.
func (h *Handler) FilterUserTableRoles(w http.ResponseWriter, r *http.Request, params openapiv1.FilterUserTableRolesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &userdomainrolesv1.GetUserTableRolesesFilteredRequest{
		DcUserName:       params.DcUserName,
		DcUserExternalId: params.DcUserExternalId,
		TableRolesName:   params.TableRolesName,
		TableCatId:       params.TableCatId,
		TableCatName:     params.TableCatName,
		Order:            order,
		PageLimit:        pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetUserTableRolesesFiltered(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserTableRolesFilteredPage(resp.GetData(), resp.GetPagination()))
}

// GetUserTableRolesById отдаёт активную строку dc.user_table_roles по id.
func (h *Handler) GetUserTableRolesById(w http.ResponseWriter, r *http.Request, id openapiv1.UserTableRolesIdParam) {
	resp, err := h.client.GetUserTableRolesById(r.Context(), &userdomainrolesv1.GetUserTableRolesByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserTableRolesDetail(resp.GetUserTableRoles()))
}

// CreateUserTableRoles выдаёт пользователю роль на таблицу. external_id (автор)
// берётся из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateUserTableRoles(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.UserTableRolesWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateUserTableRoles(r.Context(), &userdomainrolesv1.CreateUserTableRolesRequest{
		UserId:       body.UserId,
		TableRolesId: body.TableRolesId,
		TableId:      body.TableId,
		ExternalId:   claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserTableRoles(resp.GetUserTableRoles()))
}

// UpdateUserTableRolesById обновляет строку dc.user_table_roles.
func (h *Handler) UpdateUserTableRolesById(w http.ResponseWriter, r *http.Request, id openapiv1.UserTableRolesIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.UserTableRolesWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateUserTableRolesById(r.Context(), &userdomainrolesv1.UpdateUserTableRolesByIdRequest{
		Id:           id,
		UserId:       body.UserId,
		TableRolesId: body.TableRolesId,
		TableId:      body.TableId,
		ExternalId:   claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUserTableRolesById мягко удаляет строку dc.user_table_roles.
func (h *Handler) DeleteUserTableRolesById(w http.ResponseWriter, r *http.Request, id openapiv1.UserTableRolesIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteUserTableRolesById(r.Context(), &userdomainrolesv1.DeleteUserTableRolesByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteUserTableRolesById восстанавливает мягко удалённую строку dc.user_table_roles.
func (h *Handler) UndeleteUserTableRolesById(w http.ResponseWriter, r *http.Request, id openapiv1.UserTableRolesIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteUserTableRolesById(r.Context(), &userdomainrolesv1.UndeleteUserTableRolesByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
