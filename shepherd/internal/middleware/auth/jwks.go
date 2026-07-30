package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/redis/go-redis/v9"
)

// defaultTTL используется, если Keycloak не прислал Cache-Control с max-age.
const defaultTTL = 5 * time.Minute

// localCacheTTL — сколько разобранный jwk.Set держится в памяти процесса
// между проверками Redis. Отдельно от TTL самого ключа в Redis: Redis —
// источник истины с TTL по Cache-Control, это — просто защита от того,
// чтобы каждый запрос ходил в Redis.
const localCacheTTL = 1 * time.Minute

// JWKSFetcher отдаёт актуальный jwk.Set для проверки подписи токенов Keycloak.
// Кеш двухуровневый: разобранный набор ключей в памяти процесса (localCacheTTL)
// поверх сырого JSON в Redis (ключ jwks:{realm}, TTL по Cache-Control) — см.
// documentation/readme_obsidian/datacomrade/01_architecture/Кеширование Redis.md
type JWKSFetcher struct {
	redisClient *redis.Client
	httpClient  *http.Client
	jwksURL     string
	cacheKey    string

	mu        sync.RWMutex
	set       jwk.Set
	expiresAt time.Time
}

// NewJWKSFetcher строит fetcher для JWKS-эндпоинта Keycloak по issuer URL,
// например http://localhost:8081/realms/datacomrade. Realm для ключа кеша
// берётся из последнего сегмента issuer.
func NewJWKSFetcher(redisClient *redis.Client, issuerURL string) *JWKSFetcher {
	trimmed := strings.TrimRight(issuerURL, "/")
	segments := strings.Split(trimmed, "/")
	realm := segments[len(segments)-1]

	return &JWKSFetcher{
		redisClient: redisClient,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
		jwksURL:     trimmed + "/protocol/openid-connect/certs",
		cacheKey:    "jwks:" + realm,
	}
}

// Get возвращает текущий набор ключей, обновляя его при необходимости.
func (f *JWKSFetcher) Get(ctx context.Context) (jwk.Set, error) {
	f.mu.RLock()
	if f.set != nil && time.Now().Before(f.expiresAt) {
		set := f.set
		f.mu.RUnlock()
		return set, nil
	}
	f.mu.RUnlock()

	return f.refresh(ctx)
}

func (f *JWKSFetcher) refresh(ctx context.Context) (jwk.Set, error) {
	raw, err := f.fromRedis(ctx)
	if err != nil {
		raw, err = f.fromKeycloak(ctx)
		if err != nil {
			return nil, err
		}
	}

	set, err := jwk.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse jwks: %w", err)
	}

	f.mu.Lock()
	f.set = set
	f.expiresAt = time.Now().Add(localCacheTTL)
	f.mu.Unlock()

	return set, nil
}

func (f *JWKSFetcher) fromRedis(ctx context.Context) ([]byte, error) {
	if f.redisClient == nil {
		return nil, fmt.Errorf("redis client not configured")
	}

	raw, err := f.redisClient.Get(ctx, f.cacheKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("jwks cache miss: %w", err)
	}

	return raw, nil
}

// fromKeycloak запрашивает JWKS напрямую и кеширует ответ в Redis на срок из
// Cache-Control: max-age (запасной вариант — defaultTTL).
func (f *JWKSFetcher) fromKeycloak(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.jwksURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build jwks request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch jwks from %s: %w", f.jwksURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch jwks from %s: status %d", f.jwksURL, resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read jwks response: %w", err)
	}

	if f.redisClient != nil {
		ttl := maxAge(resp.Header.Get("Cache-Control"))
		if err := f.redisClient.Set(ctx, f.cacheKey, raw, ttl).Err(); err != nil {
			// Redis — оптимизация, а не источник правды: при сбое просто
			// не кешируем, следующий запрос сходит в Keycloak напрямую.
			_ = err
		}
	}

	return raw, nil
}

// maxAge разбирает "max-age=N" из Cache-Control. Возвращает defaultTTL,
// если заголовка нет или он не парсится.
func maxAge(cacheControl string) time.Duration {
	for _, part := range strings.Split(cacheControl, ",") {
		part = strings.TrimSpace(part)
		if rest, ok := strings.CutPrefix(part, "max-age="); ok {
			if seconds, err := strconv.Atoi(rest); err == nil && seconds > 0 {
				return time.Duration(seconds) * time.Second
			}
		}
	}

	return defaultTTL
}
