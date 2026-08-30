// Package queryengineapiv1 реализует часть openapiv1.ServerInterface,
// отвечающую за OLAP-запросы (/v1/olap-queries...) — REST-хендлеры вручную
// поверх сгенерированного gRPC-клиента queryenginev1.QueryEngineServiceClient
// (query_engine). Доступ — только аутентификация, как у /olap-aliases:
// отдельной роли Gateway не требуется. created_by берётся из claims токена,
// а не из тела запроса.
package queryengineapiv1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	queryenginev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/query_engine/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/ensureuser"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к
// /v1/olap-queries.
type Handler struct {
	client   queryenginev1.QueryEngineServiceClient
	resolver *ensureuser.Resolver
}

func New(client queryenginev1.QueryEngineServiceClient, resolver *ensureuser.Resolver) *Handler {
	return &Handler{client: client, resolver: resolver}
}

// CreateOlapQuery принимает конфигурацию сводной от конструктора, добавляет
// целевую базу и автора и отдаёт её в query_engine на сохранение. created_by
// резолвится из claims токена (read-only, как в GET /me): если dc.user ещё
// нет — 404, вызывающий должен сначала сходить в POST /login.
func (h *Handler) CreateOlapQuery(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	var body openapiv1.OlapQueryCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "некорректное тело запроса")
		return
	}

	if body.DatabaseId <= 0 {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "database_id обязателен и должен быть положительным")
		return
	}
	if bad := unknownFilterOperators(body); len(bad) > 0 {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "неизвестный оператор фильтра: "+strings.Join(bad, ", "))
		return
	}
	if invalidMeasures(body) {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "показатель kind=field требует alias_id и aggregation, kind=formula — непустую formula")
		return
	}

	userID, err := h.resolver.Resolve(r.Context(), claims.Subject)
	if err != nil {
		if errors.Is(err, ensureuser.ErrNotFound) {
			apierror.Write(w, http.StatusNotFound, "user_not_provisioned", "пользователь ещё не создан — сначала вызовите POST /login")
			return
		}
		apierror.Write(w, http.StatusInternalServerError, "internal", "не удалось получить пользователя")
		return
	}

	resp, err := h.client.CreateQuery(r.Context(), &queryenginev1.CreateQueryRequest{
		Query:      toOlapQuery(body),
		DatabaseId: body.DatabaseId,
		CreatedBy:  userID,
	})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	// Алиасы формы встречаются в нескольких доменах, domain_id не задан —
	// query_engine ничего не сохранил, вернул список доменов на выбор.
	if choice := resp.GetDomainChoice(); choice != nil {
		writeJSON(w, http.StatusConflict, toDomainChoice(choice))
		return
	}

	writeJSON(w, http.StatusCreated, toOlapQueryRecord(resp.GetSaved()))
}

// GetOlapQueryById отдаёт ранее сохранённую строку queries.queries по id.
func (h *Handler) GetOlapQueryById(w http.ResponseWriter, r *http.Request, id int64) {
	if id <= 0 {
		apierror.Write(w, http.StatusBadRequest, "invalid_argument", "id должен быть положительным")
		return
	}

	resp, err := h.client.GetQueryById(r.Context(), &queryenginev1.GetQueryByIdRequest{Id: id})
	if err != nil {
		apierror.WriteGRPC(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toOlapQueryRecord(resp.GetQuery()))
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
