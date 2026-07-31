// Package auth проверяет access-token Keycloak на входе в Shepherd:
// локальная валидация JWT по JWKS, без token introspection на каждый запрос —
// см. documentation/readme_obsidian/datacomrade/01_architecture/API Gateway.md
package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type contextKey string

const claimsContextKey contextKey = "shepherd_auth_claims"

// Claims — то, что Shepherd вычитывает из access-token и прокидывает дальше
// по запросу (в контекст и в gRPC metadata бэкендов).
type Claims struct {
	// Subject — sub из токена, он же external_id в dc.user (Metadata Service).
	Subject string
	Name    string
	Email   string
	// Roles — realm-роли Keycloak (admin/maintainer/viewer, см. константы
	// RoleAdmin/RoleMaintainer/RoleViewer в middleware/rbac). Не путать с
	// dc.domain_roles/dc.table_roles — это гейт на уровне маршрута Gateway,
	// а не прав на домены/таблицы. См. Разграничение доступа.md.
	Roles []string
}

// HasRole сообщает, есть ли у пользователя хотя бы одна из перечисленных ролей.
func (c Claims) HasRole(roles ...string) bool {
	for _, want := range roles {
		for _, have := range c.Roles {
			if have == want {
				return true
			}
		}
	}
	return false
}

// ContextWithClaims кладёт claims в контекст запроса.
func ContextWithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, c)
}

// ClaimsFromContext достаёт claims, положенные туда Middleware.Authenticate.
func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsContextKey).(Claims)
	return c, ok
}

// Middleware проверяет access-token Keycloak: подпись по JWKS, issuer, срок действия.
type Middleware struct {
	jwks   *JWKSFetcher
	issuer string
}

// NewMiddleware строит middleware проверки токена для заданного issuer'а Keycloak.
func NewMiddleware(jwks *JWKSFetcher, issuerURL string) *Middleware {
	return &Middleware{jwks: jwks, issuer: issuerURL}
}

// Authenticate — chi-совместимый middleware. Требует валидный Bearer-токен,
// иначе отвечает 401 и не пропускает запрос дальше.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := bearerToken(r)
		if raw == "" {
			apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "отсутствует Bearer-токен")
			return
		}

		set, err := m.jwks.Get(r.Context())
		if err != nil {
			apierror.Write(w, http.StatusServiceUnavailable, "jwks_unavailable", "не удалось получить ключи Keycloak")
			return
		}

		tok, err := jwt.Parse([]byte(raw), jwt.WithKeySet(set), jwt.WithIssuer(m.issuer))
		if err != nil {
			apierror.Write(w, http.StatusUnauthorized, "invalid_token", "токен недействителен")
			return
		}

		ctx := ContextWithClaims(r.Context(), claimsFromToken(tok))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func claimsFromToken(tok jwt.Token) Claims {
	sub, _ := tok.Subject()

	var name string
	_ = tok.Get("preferred_username", &name)

	var email string
	_ = tok.Get("email", &email)

	var realmAccess struct {
		Roles []string `json:"roles"`
	}
	_ = tok.Get("realm_access", &realmAccess)

	return Claims{
		Subject: sub,
		Name:    name,
		Email:   email,
		Roles:   realmAccess.Roles,
	}
}
