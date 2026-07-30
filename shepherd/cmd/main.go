package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	pkglogger "github.com/konstantin-suspitsyn/datacomrade/platform/pkg/logger"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/meapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/client/datacatalogue"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/config/constants"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/ensureuser"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/requestlog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const openapiSpecPath = "api/openapi/v1/openapi.yaml"

func main() {
	var env string
	flag.StringVar(&env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		panic("godotenv file was not found")
	}

	cfg := constants.InitConfig()

	if err := pkglogger.Init(cfg.LoggerLevel, cfg.LoggerAsJSON); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func() { _ = pkglogger.Sync() }()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	defer func() { _ = redisClient.Close() }()

	dcClient, err := datacatalogue.Dial(cfg.DataCatalogueGRPCAddr)
	if err != nil {
		log.Fatalf("failed to dial datacatalogue: %v", err)
	}
	defer func() { _ = dcClient.Close() }()

	jwks := auth.NewJWKSFetcher(redisClient, cfg.KeycloakIssuerURL)
	authMiddleware := auth.NewMiddleware(jwks, cfg.KeycloakIssuerURL)
	ensureUserMiddleware := ensureuser.New(dcClient.User, redisClient, cfg.EnsureUserCacheTTL)

	router := chi.NewRouter()
	router.Use(requestlog.Middleware)
	router.Use(chimiddleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	router.Get("/healthz", handleHealthz)
	router.Get("/readyz", handleReadyz(redisClient))

	if cfg.EnableDocs {
		mountDocs(router)
	}

	router.Route("/v1", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(ensureUserMiddleware.Handler)

		openapiv1.HandlerFromMux(meapiv1.New(), r)
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		pkglogger.Info(context.Background(), "shepherd listening", zap.Int("port", cfg.HTTPPort), zap.String("env", env))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	pkglogger.Info(context.Background(), "shutting down shepherd")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		pkglogger.Error(context.Background(), "graceful shutdown failed", zap.Error(err))
	}

	pkglogger.Info(context.Background(), "shepherd stopped")
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// handleReadyz проверяет только Redis: он единственная зависимость, без
// которой процесс не может обслуживать запросы (JWKS и EnsureUser его
// требуют). datacatalogue дозвонится лениво через grpc.NewClient — его
// недоступность проявится как ошибка конкретного запроса, а не всего сервиса.
func handleReadyz(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			apierror.Write(w, http.StatusServiceUnavailable, "not_ready", "redis недоступен")
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

// mountDocs включает /openapi.yaml и /docs (Swagger UI). Держать выключенным
// в production (SHEPHERD_ENABLE_DOCS=false) — см. API Gateway.md.
func mountDocs(r chi.Router) {
	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, openapiSpecPath)
	})

	r.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerUIPage))
	})
}

const swaggerUIPage = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Shepherd API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({ url: "/openapi.yaml", dom_id: "#swagger-ui" });
  </script>
</body>
</html>`
