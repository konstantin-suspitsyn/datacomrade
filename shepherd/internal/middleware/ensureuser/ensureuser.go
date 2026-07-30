// Package ensureuser гарантирует, что для аутентифицированного пользователя
// существует строка dc.user в Metadata Service, и резолвит её численный id —
// он нужен бэкендам вроде UserDomainRolesService, которые принимают user_id,
// а не external_id. Подробности флоу и оговорка про гонку при создании — в
// documentation/readme_obsidian/datacomrade/01_architecture/API Gateway.md
// ("Автосоздание пользователя") и Замечания к реализации.md п.9.
package ensureuser

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	pkglogger "github.com/konstantin-suspitsyn/datacomrade/platform/pkg/logger"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const userIDContextKey contextKey = "shepherd_dc_user_id"

// ContextWithUserID кладёт в контекст численный id из dc.user.
func ContextWithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, id)
}

// UserIDFromContext достаёт id, положенный туда Middleware.Handler.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDContextKey).(int64)
	return id, ok
}

// Middleware реализует get-or-create для dc.user по external_id (=sub из JWT).
type Middleware struct {
	users userv1.UserServiceClient
	redis *redis.Client
	ttl   time.Duration
}

// New строит middleware. ttl — срок жизни кеша резолва в Redis
// (ключ user:known:{external_id}); промах кеша не критичен для безопасности,
// это чистая оптимизация, см. Кеширование Redis.md.
func New(users userv1.UserServiceClient, redisClient *redis.Client, ttl time.Duration) *Middleware {
	return &Middleware{users: users, redis: redisClient, ttl: ttl}
}

// Handler — chi-совместимый middleware. Должен вешаться после auth.Middleware.Authenticate.
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := auth.ClaimsFromContext(r.Context())
		if !ok {
			apierror.Write(w, http.StatusUnauthorized, "unauthenticated", "требуется аутентификация")
			return
		}

		userID, err := m.resolve(r.Context(), claims)
		if err != nil {
			apierror.Write(w, http.StatusInternalServerError, "user_provisioning_failed", err.Error())
			return
		}

		ctx := ContextWithUserID(r.Context(), userID)
		ctx = pkglogger.ContextWithUserID(ctx, strconv.FormatInt(userID, 10))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func cacheKey(externalID string) string {
	return "user:known:" + externalID
}

func (m *Middleware) resolve(ctx context.Context, claims auth.Claims) (int64, error) {
	if cached, err := m.redis.Get(ctx, cacheKey(claims.Subject)).Result(); err == nil {
		if id, parseErr := strconv.ParseInt(cached, 10, 64); parseErr == nil {
			return id, nil
		}
	}

	id, err := m.getOrCreate(ctx, claims)
	if err != nil {
		return 0, err
	}

	// Best-effort: Redis — оптимизация, а не источник правды для этого ключа.
	m.redis.Set(ctx, cacheKey(claims.Subject), id, m.ttl)

	return id, nil
}

func (m *Middleware) getOrCreate(ctx context.Context, claims auth.Claims) (int64, error) {
	resp, err := m.users.GetUserByExternalId(ctx, &userv1.GetUserByExternalIdRequest{ExternalId: claims.Subject})
	if err == nil {
		return resp.GetUser().GetId(), nil
	}
	if status.Code(err) != codes.NotFound {
		return 0, fmt.Errorf("get user by external id: %w", err)
	}

	name := claims.Name
	if name == "" {
		name = claims.Email
	}
	if name == "" {
		name = claims.Subject
	}

	created, err := m.users.CreateUser(ctx, &userv1.CreateUserRequest{Name: name, ExternalId: claims.Subject})
	if err == nil {
		return created.GetUser().GetId(), nil
	}

	// Гонка: два параллельных первых запроса одного нового пользователя.
	// CreateUser в Metadata Service — обычный INSERT без ON CONFLICT (см.
	// Замечания к реализации.md п.9), поэтому проигравший ретраит чтение
	// один раз вместо того, чтобы падать всем запросом.
	retryResp, retryErr := m.users.GetUserByExternalId(ctx, &userv1.GetUserByExternalIdRequest{ExternalId: claims.Subject})
	if retryErr != nil {
		return 0, fmt.Errorf("create user failed (%w) and retry lookup failed: %w", err, retryErr)
	}

	return retryResp.GetUser().GetId(), nil
}
