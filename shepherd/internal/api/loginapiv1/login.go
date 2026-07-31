// Package loginapiv1 реализует часть openapiv1.ServerInterface, отвечающую
// за POST /v1/login — единственное место, где Gateway создаёт запись dc.user
// при первом входе пользователя. См. API Gateway.md ("Автосоздание
// пользователя").
package loginapiv1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/ensureuser"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
)

// Handler реализует методы openapiv1.ServerInterface, относящиеся к /login.
type Handler struct {
	resolver *ensureuser.Resolver
}

func New(resolver *ensureuser.Resolver) *Handler {
	return &Handler{resolver: resolver}
}

// Login резолвит dc.user по external_id из токена, создавая запись, если её
// ещё нет.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
		return
	}

	userID, err := h.resolver.GetOrCreate(r.Context(), claims)
	if err != nil {
		apierror.Write(w, http.StatusInternalServerError, "user_provisioning_failed", err.Error())
		return
	}

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
