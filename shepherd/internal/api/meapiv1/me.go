// Package meapiv1 реализует часть openapiv1.ServerInterface, отвечающую за
// identity текущего пользователя (GET /v1/me) — первый сквозной эндпоинт
// Shepherd: auth -> ensureuser -> ответ, см. API Gateway.md.
package meapiv1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/ensureuser"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к /me.
type Handler struct{}

func New() *Handler {
	return &Handler{}
}

// GetMe отдаёт identity, собранную middleware auth и ensureuser.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	userID, _ := ensureuser.UserIDFromContext(r.Context())

	externalID, err := uuid.Parse(claims.Subject)
	if err != nil {
		apierror.Write(w, http.StatusInternalServerError, "internal", "sub токена не является UUID")
		return
	}

	resp := openapiv1.Me{
		ExternalId: externalID,
		UserId:     userID,
		Name:       claims.Name,
		Roles:      claims.Roles,
	}
	if claims.Email != "" {
		resp.Email = &claims.Email
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(resp)
}
