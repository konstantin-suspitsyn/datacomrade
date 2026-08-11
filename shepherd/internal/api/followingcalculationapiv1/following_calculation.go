// Package followingcalculationapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.following_calculation (/v1/following-calculations...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.FollowingCalculationServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package followingcalculationapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.following_calculation.
type Handler struct {
	client tablesv1.FollowingCalculationServiceClient
}

func New(client tablesv1.FollowingCalculationServiceClient) *Handler {
	return &Handler{client: client}
}

// ListFollowingCalculations отдаёт страницу dc.following_calculation. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListFollowingCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListFollowingCalculationsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetFollowingCalculationsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetFollowingCalculations(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toFollowingCalculationPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedFollowingCalculations отдаёт страницу мягко удалённых строк dc.following_calculation.
func (h *Handler) ListDeletedFollowingCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedFollowingCalculationsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetFollowingCalculationsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetFollowingCalculationsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toFollowingCalculationPage(resp.GetData(), resp.GetPagination()))
}

// FilterFollowingCalculations отдаёт страницу dc.following_calculation с необязательными фильтрами. Непереданный
// фильтр уходит в gRPC как отсутствующий (proto optional) и не сужает выборку.
func (h *Handler) FilterFollowingCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.FilterFollowingCalculationsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetFollowingCalculationsFiltersRequest{
		ColumnCatId:       params.ColumnCatId,
		CalculationTypeId: params.CalculationTypeId,
		Order:             order,
		PageLimit:         pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetFollowingCalculationsFilters(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toFollowingCalculationPage(resp.GetData(), resp.GetPagination()))
}

// GetFollowingCalculationById отдаёт активную строку dc.following_calculation по id.
func (h *Handler) GetFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	resp, err := h.client.GetFollowingCalculationById(r.Context(), &tablesv1.GetFollowingCalculationByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toFollowingCalculation(resp.GetFollowingCalculation()))
}

// CreateFollowingCalculation создаёт строку dc.following_calculation. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateFollowingCalculation(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.FollowingCalculationWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateFollowingCalculation(r.Context(), &tablesv1.CreateFollowingCalculationRequest{
		ColumnCatId:       body.ColumnCatId,
		CalculationTypeId: body.CalculationTypeId,
		ExternalId:        claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toFollowingCalculation(resp.GetFollowingCalculation()))
}

// UpdateFollowingCalculationById обновляет строку dc.following_calculation.
func (h *Handler) UpdateFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.FollowingCalculationWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateFollowingCalculationById(r.Context(), &tablesv1.UpdateFollowingCalculationByIdRequest{
		Id:                id,
		ColumnCatId:       body.ColumnCatId,
		CalculationTypeId: body.CalculationTypeId,
		ExternalId:        claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteFollowingCalculationById мягко удаляет строку dc.following_calculation.
func (h *Handler) DeleteFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteFollowingCalculationById(r.Context(), &tablesv1.DeleteFollowingCalculationByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteFollowingCalculationById восстанавливает мягко удалённую строку dc.following_calculation.
func (h *Handler) UndeleteFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteFollowingCalculationById(r.Context(), &tablesv1.UndeleteFollowingCalculationByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
