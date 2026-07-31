// Package ensureuser резолвит численный user_id из dc.user (Metadata Service)
// по external_id (=sub из JWT Keycloak). Вызывается явно из двух хендлеров с
// разной семантикой — не как middleware на каждый запрос:
//   - loginapiv1.Login — GetOrCreate: создаёт запись dc.user, если её нет.
//   - meapiv1.GetMe — Resolve: только читает, не создаёт.
//
// См. documentation/readme_obsidian/datacomrade/01_architecture/API Gateway.md
// ("Автосоздание пользователя").
package ensureuser

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrNotFound — в dc.user ещё нет записи для этого external_id. Вызывающий
// должен предложить пользователю сначала сходить в POST /v1/login.
var ErrNotFound = errors.New("user not provisioned")

// Resolver резолвит и, если явно попросили, создаёт dc.user по external_id.
type Resolver struct {
	users userv1.UserServiceClient
	redis *redis.Client
	ttl   time.Duration
}

// New строит Resolver. ttl — срок жизни кеша резолва в Redis
// (ключ user:known:{external_id}); промах кеша не критичен для безопасности,
// это чистая оптимизация, см. Кеширование Redis.md.
func New(users userv1.UserServiceClient, redisClient *redis.Client, ttl time.Duration) *Resolver {
	return &Resolver{users: users, redis: redisClient, ttl: ttl}
}

func cacheKey(externalID string) string {
	return "user:known:" + externalID
}

func (r *Resolver) cached(ctx context.Context, externalID string) (int64, bool) {
	cached, err := r.redis.Get(ctx, cacheKey(externalID)).Result()
	if err != nil {
		return 0, false
	}
	id, err := strconv.ParseInt(cached, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func (r *Resolver) cache(ctx context.Context, externalID string, id int64) {
	// Best-effort: Redis — оптимизация, а не источник правды для этого ключа.
	r.redis.Set(ctx, cacheKey(externalID), id, r.ttl)
}

// Resolve находит user_id по external_id. Не создаёт запись — если её нет,
// возвращает ErrNotFound.
func (r *Resolver) Resolve(ctx context.Context, externalID string) (int64, error) {
	if id, ok := r.cached(ctx, externalID); ok {
		return id, nil
	}

	resp, err := r.users.GetUserByExternalId(ctx, &userv1.GetUserByExternalIdRequest{ExternalId: externalID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("get user by external id: %w", err)
	}

	id := resp.GetUser().GetId()
	r.cache(ctx, externalID, id)

	return id, nil
}

// GetOrCreate находит user_id по external_id, создавая dc.user, если записи
// ещё нет.
func (r *Resolver) GetOrCreate(ctx context.Context, claims auth.Claims) (int64, error) {
	if id, ok := r.cached(ctx, claims.Subject); ok {
		return id, nil
	}

	id, err := r.getOrCreate(ctx, claims)
	if err != nil {
		return 0, err
	}

	r.cache(ctx, claims.Subject, id)

	return id, nil
}

func (r *Resolver) getOrCreate(ctx context.Context, claims auth.Claims) (int64, error) {
	resp, err := r.users.GetUserByExternalId(ctx, &userv1.GetUserByExternalIdRequest{ExternalId: claims.Subject})
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

	created, err := r.users.CreateUser(ctx, &userv1.CreateUserRequest{Name: name, ExternalId: claims.Subject})
	if err == nil {
		return created.GetUser().GetId(), nil
	}

	// Гонка: два параллельных первых запроса одного нового пользователя.
	// CreateUser в Metadata Service — обычный INSERT без ON CONFLICT (см.
	// Замечания к реализации.md п.9), поэтому проигравший ретраит чтение
	// один раз вместо того, чтобы падать всем запросом.
	retryResp, retryErr := r.users.GetUserByExternalId(ctx, &userv1.GetUserByExternalIdRequest{ExternalId: claims.Subject})
	if retryErr != nil {
		return 0, fmt.Errorf("create user failed (%w) and retry lookup failed: %w", err, retryErr)
	}

	return retryResp.GetUser().GetId(), nil
}
