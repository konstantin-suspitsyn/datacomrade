// Package rbac проверяет realm-роли Keycloak (Admin/Reader/Writer) на уровне
// маршрута — это отдельный, более грубый гейт, чем dc.domain_roles/dc.table_roles
// в Metadata Service. См. API Gateway.md и Разграничение доступа.md.
package rbac

import (
	"net/http"

	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
)

// RequireRole пропускает запрос, только если у пользователя есть хотя бы
// одна из перечисленных realm-ролей. Должен вешаться после auth.Middleware.Authenticate.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.ClaimsFromContext(r.Context())
			if !ok {
				apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
				return
			}

			if !claims.HasRole(roles...) {
				apierror.Write(w, http.StatusForbidden, "forbidden", "недостаточно прав")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
