package constants

import (
	"fmt"
	"time"
)

// Config собирает конфигурацию Shepherd из переменных окружения.
// См. documentation/readme_obsidian/datacomrade/01_architecture/API Gateway.md
type Config struct {
	HTTPPort int

	LoggerLevel  string
	LoggerAsJSON bool

	// Allowlist origin'ов фронтенда для CORS.
	CORSOrigins []string

	// Issuer realm'а Keycloak, например http://localhost:8081/realms/datacomrade.
	// JWKS-адрес вычисляется из него по стандартному пути Keycloak.
	KeycloakIssuerURL string

	DataCatalogueGRPCAddr string

	RedisAddr     string
	RedisPassword string

	// TTL кеша резолва пользователя (external_id -> user_id) в Redis —
	// см. Кеширование Redis.md, ключ user:known:{external_id}.
	EnsureUserCacheTTL time.Duration

	// Включает /openapi.yaml и /docs (Swagger UI). Держать выключенным в проде.
	EnableDocs bool
}

func InitConfig() Config {
	return Config{
		HTTPPort: getIntEnv("SHEPHERD_HTTP_PORT"),

		LoggerLevel:  getStringEnv("SHEPHERD_LOGGER_LEVEL"),
		LoggerAsJSON: getBoolEnv("SHEPHERD_LOGGER_AS_JSON"),

		CORSOrigins: getStringSliceEnv("SHEPHERD_CORS_ORIGINS"),

		KeycloakIssuerURL: getStringEnv("SHEPHERD_KEYCLOAK_ISSUER_URL"),

		DataCatalogueGRPCAddr: fmt.Sprintf("%s:%d", getStringEnv("DATACATALOGUE_GRPC_HOST"), getIntEnv("DATACATALOGUE_GRPC_PORT")),

		RedisAddr:     fmt.Sprintf("%s:%d", getStringEnv("SHEPHERD_REDIS_HOST"), getIntEnv("REDIS_PORT")),
		RedisPassword: getStringEnv("REDIS_PASSWORD"),

		EnsureUserCacheTTL: 30 * time.Minute,

		EnableDocs: getBoolEnv("SHEPHERD_ENABLE_DOCS"),
	}
}
