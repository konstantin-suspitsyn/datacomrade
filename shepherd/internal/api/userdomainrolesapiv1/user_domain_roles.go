// Package userdomainrolesapiv1 реализует часть openapiv1.ServerInterface,
// отвечающую за dc.user_domain_roles (/v1/user-domain-roles...) — REST-хендлеры
// вручную поверх сгенерированного gRPC-клиента
// userdomainrolesv1.UserDomainRolesServiceClient (Metadata Service),
// см. API Gateway.md. Доступ гейтится ролью admin в shepherd/cmd/main.go
// (rbac.RequireRole), а не в этом пакете.
package userdomainrolesapiv1

import (
	"encoding/json"
	"net/http"

	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.user_domain_roles.
type Handler struct {
	client userdomainrolesv1.UserDomainRolesServiceClient
}

func New(client userdomainrolesv1.UserDomainRolesServiceClient) *Handler {
	return &Handler{client: client}
}

// ListUserDomainRoles отдаёт страницу dc.user_domain_roles. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListUserDomainRoles(w http.ResponseWriter, r *http.Request, params openapiv1.ListUserDomainRolesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &userdomainrolesv1.GetUserDomainRolesesRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetUserDomainRoleses(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserDomainRolesPage(resp.GetData(), resp.GetPagination()))
}

// CreateUserDomainRoles выдаёт пользователю роль на домен. external_id (автор)
// берётся из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateUserDomainRoles(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.UserDomainRolesWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateUserDomainRoles(r.Context(), &userdomainrolesv1.CreateUserDomainRolesRequest{
		UserId:        body.UserId,
		DomainRolesId: body.DomainRolesId,
		DomainId:      body.DomainId,
		ExternalId:    claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserDomainRoles(resp.GetUserDomainRoles()))
}

// UpdateUserDomainRolesById обновляет строку dc.user_domain_roles.
func (h *Handler) UpdateUserDomainRolesById(w http.ResponseWriter, r *http.Request, id openapiv1.UserDomainRolesIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.UserDomainRolesWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateUserDomainRolesById(r.Context(), &userdomainrolesv1.UpdateUserDomainRolesByIdRequest{
		Id:            id,
		UserId:        body.UserId,
		DomainRolesId: body.DomainRolesId,
		DomainId:      body.DomainId,
		ExternalId:    claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUserDomainRolesById мягко удаляет строку dc.user_domain_roles.
func (h *Handler) DeleteUserDomainRolesById(w http.ResponseWriter, r *http.Request, id openapiv1.UserDomainRolesIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteUserDomainRolesById(r.Context(), &userdomainrolesv1.DeleteUserDomainRolesByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteUserDomainRolesById восстанавливает мягко удалённую строку dc.user_domain_roles.
func (h *Handler) UndeleteUserDomainRolesById(w http.ResponseWriter, r *http.Request, id openapiv1.UserDomainRolesIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteUserDomainRolesById(r.Context(), &userdomainrolesv1.UndeleteUserDomainRolesByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
