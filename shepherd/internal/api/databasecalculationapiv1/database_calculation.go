// Package databasecalculationapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.database_calculation (/v1/database-calculations...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.DatabaseCalculationServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package databasecalculationapiv1

import (
	"encoding/json"
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.database_calculation.
type Handler struct {
	client tablesv1.DatabaseCalculationServiceClient
}

func New(client tablesv1.DatabaseCalculationServiceClient) *Handler {
	return &Handler{client: client}
}

// ListDatabaseCalculations отдаёт страницу dc.database_calculation. order/page_limit по умолчанию — см. пакет pagination.
func (h *Handler) ListDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseCalculationsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseCalculationsRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCalculations(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCalculationPage(resp.GetData(), resp.GetPagination()))
}

// SearchDatabaseCalculations отдаёт страницу dc.database_calculation, отфильтрованную по подстроке имени.
func (h *Handler) SearchDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDatabaseCalculationsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseCalculationsSearchNameRequest{
		SearchName: params.SearchName,
		Order:      order,
		PageLimit:  pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCalculationsSearchName(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCalculationPage(resp.GetData(), resp.GetPagination()))
}

// ListDeletedDatabaseCalculations отдаёт страницу мягко удалённых строк dc.database_calculation.
func (h *Handler) ListDeletedDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDatabaseCalculationsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseCalculationsDeletedRequest{
		Order:     order,
		PageLimit: pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCalculationsDeleted(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCalculationPage(resp.GetData(), resp.GetPagination()))
}

// FilterDatabaseCalculations отдаёт страницу dc.database_calculation с необязательными фильтрами. Непереданный
// фильтр уходит в gRPC как отсутствующий (proto optional) и не сужает выборку.
func (h *Handler) FilterDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.FilterDatabaseCalculationsParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetDatabaseCalculationsFilteredRequest{
		DatabaseCatId:     params.DatabaseCatId,
		CalculationTypeId: params.CalculationTypeId,
		Order:             order,
		PageLimit:         pagination.ResolvePageLimit(params.PageLimit),
	}
	if params.Page != nil {
		req.Page = *params.Page
	}

	resp, err := h.client.GetDatabaseCalculationsFiltered(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCalculationPage(resp.GetData(), resp.GetPagination()))
}

// GetDatabaseCalculationById отдаёт активную строку dc.database_calculation по id.
func (h *Handler) GetDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	resp, err := h.client.GetDatabaseCalculationById(r.Context(), &tablesv1.GetDatabaseCalculationByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDatabaseCalculation(resp.GetDatabaseCalculation()))
}

// CreateDatabaseCalculation создаёт строку dc.database_calculation. external_id (автор) берётся
// из claims токена, а не из тела запроса — фронтенд его не передаёт.
func (h *Handler) CreateDatabaseCalculation(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DatabaseCalculationWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	resp, err := h.client.CreateDatabaseCalculation(r.Context(), &tablesv1.CreateDatabaseCalculationRequest{
		DatabaseCatId:     body.DatabaseCatId,
		CalculationTypeId: body.CalculationTypeId,
		ExternalId:        claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toDatabaseCalculation(resp.GetDatabaseCalculation()))
}

// UpdateDatabaseCalculationById обновляет строку dc.database_calculation.
func (h *Handler) UpdateDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.DatabaseCalculationWrite
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	_, err := h.client.UpdateDatabaseCalculationById(r.Context(), &tablesv1.UpdateDatabaseCalculationByIdRequest{
		Id:                id,
		DatabaseCatId:     body.DatabaseCatId,
		CalculationTypeId: body.CalculationTypeId,
		ExternalId:        claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteDatabaseCalculationById мягко удаляет строку dc.database_calculation.
func (h *Handler) DeleteDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.DeleteDatabaseCalculationById(r.Context(), &tablesv1.DeleteDatabaseCalculationByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UndeleteDatabaseCalculationById восстанавливает мягко удалённую строку dc.database_calculation.
func (h *Handler) UndeleteDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	_, err := h.client.UndeleteDatabaseCalculationById(r.Context(), &tablesv1.UndeleteDatabaseCalculationByIdRequest{
		Id:         id,
		ExternalId: claims.Subject,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
