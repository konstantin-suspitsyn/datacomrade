// Package olapexchangeapiv1 реализует часть openapiv1.ServerInterface,
// отвечающую за OLAP-обмен (/v1/olap-aliases...) — REST-хендлеры вручную
// поверх сгенерированного gRPC-клиента olapexchangev1.AliasServiceClient
// (Metadata Service), см. API Gateway.md. Доступ каждому пользователю
// ограничен уже на уровне SQL (dc.get_table_ids_by_external_user_and_role,
// роль can_read), поэтому здесь нет rbac.RequireRole — достаточно быть
// аутентифицированным.
package olapexchangeapiv1

import (
	"encoding/json"
	"net/http"

	olapexchangev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/olap_exchange/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/pagination"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к OLAP-обмену.
type Handler struct {
	client olapexchangev1.AliasServiceClient
}

func New(client olapexchangev1.AliasServiceClient) *Handler {
	return &Handler{client: client}
}

func toOlapAlias(row *olapexchangev1.GetOlapAliasesRow) openapiv1.OlapAlias {
	return openapiv1.OlapAlias{
		Id:                    row.GetId(),
		Name:                  row.GetName(),
		Description:           row.GetDescription(),
		IsDeleted:             row.GetIsDeleted(),
		ColumnTypeId:          row.GetColumnTypeId(),
		ColumnTypeName:        row.GetColumnTypeName(),
		ColumnTypeDescription: row.GetColumnTypeDescription(),
	}
}

func toOlapAliasList(rows []*olapexchangev1.GetOlapAliasesRow) openapiv1.OlapAliasList {
	data := make([]openapiv1.OlapAlias, 0, len(rows))
	for _, row := range rows {
		data = append(data, toOlapAlias(row))
	}

	return openapiv1.OlapAliasList{Data: data}
}

// ListOlapAliases отдаёт алиасы, доступные текущему пользователю по его
// ролям (can_read). external_id берётся из claims токена, а не из запроса.
func (h *Handler) ListOlapAliases(w http.ResponseWriter, r *http.Request, params openapiv1.ListOlapAliasesParams) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	order, ok := pagination.ResolveOrder(params.Order)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "order должен быть ASC или DESC")
		return
	}

	resp, err := h.client.GetOlapAliases(r.Context(), &olapexchangev1.GetOlapAliasesRequest{
		ExternalId: claims.Subject,
		Order:      order,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toOlapAliasList(resp.GetRows()))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
