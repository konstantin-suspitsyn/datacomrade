// Package rbac проверяет realm-роли Keycloak (admin/maintainer/viewer, см.
// roles.go) на уровне маршрута — это отдельный, более грубый гейт, чем
// dc.domain_roles/dc.table_roles в Metadata Service. См. API Gateway.md и
// Разграничение доступа.md.
//
// Ожидаемое сопоставление роль → группа методов Data Catalog:
//   - RoleAdmin: datacatalogue/internal/api/userdomainrolesapiv1 (просмотр
//     пользователей, назначение прав).
//   - RoleMaintainer: datacatalogue/internal/api/tablesapiv1 (каталог: хосты,
//     базы, домены, таблицы, колонки).
//   - Создание пользователя (UserService.CreateUser через EnsureUser)
//     нарочно НЕ гейтится ролью — доступно любому аутентифицированному
//     пользователю независимо от роли.
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
