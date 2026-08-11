// Package calculationtypeapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за dc.calculation_type (/v1/calculation-types...) — REST-хендлеры вручную поверх
// сгенерированного gRPC-клиента tablesv1.CalculationTypeServiceClient
// (Metadata Service), см. API Gateway.md. Доступ гейтится ролью
// admin/maintainer в shepherd/cmd/main.go (rbac.RequireRole), а не в этом пакете.
package calculationtypeapiv1

import (
	"net/http"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к dc.calculation_type.
type Handler struct {
	client tablesv1.CalculationTypeServiceClient
}

func New(client tablesv1.CalculationTypeServiceClient) *Handler {
	return &Handler{client: client}
}

// ListCalculationTypes отдаёт все строки dc.calculation_type. order по умолчанию — см. пакет pagination.
func (h *Handler) ListCalculationTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListCalculationTypesParams) {
	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	req := &tablesv1.GetCalculationTypesRequest{
		Order: order,
	}

	resp, err := h.client.GetCalculationTypes(r.Context(), req)
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCalculationTypeSlice(resp.GetRows()))
}

// GetCalculationTypeById отдаёт активную строку dc.calculation_type по id.
func (h *Handler) GetCalculationTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.CalculationTypeIdParam) {
	resp, err := h.client.GetCalculationTypeById(r.Context(), &tablesv1.GetCalculationTypeByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCalculationType(resp.GetCalculationType()))
}
