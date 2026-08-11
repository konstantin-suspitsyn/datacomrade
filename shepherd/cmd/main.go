package main

import (
	"context"
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
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/aliasapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/calculationtypeapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/columncatapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/columntypeapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/databasecalculationapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/databasecatapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/databasetypeapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/domaincatapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/followingcalculationapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/grouplevelsapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/hastogroupapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/hostapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/loginapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/meapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/schemacatapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/tablecatapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/tabletypeapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/apierror"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/client/datacatalogue"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/config/constants"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/ensureuser"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/auth"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/rbac"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/middleware/requestlog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const openapiSpecPath = "api/openapi/v1/openapi.yaml"

func main() {
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
		pkglogger.Fatal(context.Background(), "failed to dial datacatalogue", zap.Error(err))
	}
	defer func() { _ = dcClient.Close() }()

	jwks := auth.NewJWKSFetcher(redisClient, cfg.KeycloakIssuerURL)
	authMiddleware := auth.NewMiddleware(jwks, cfg.KeycloakIssuerURL)
	userResolver := ensureuser.New(dcClient.User, redisClient, cfg.EnsureUserCacheTTL)

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

		openapiv1.HandlerFromMux(newAPI(userResolver, dcClient), r)
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		pkglogger.Info(context.Background(), "shepherd listening", zap.Int("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serveErr:
		pkglogger.Error(context.Background(), "shepherd server failed, shutting down", zap.Error(err))
	case <-quit:
		pkglogger.Info(context.Background(), "shutting down shepherd")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		pkglogger.Error(context.Background(), "graceful shutdown failed", zap.Error(err))
	}

	pkglogger.Info(context.Background(), "shepherd stopped")
}

// api собирает независимые под-хендлеры (/me, /login и по пакету на таблицу
// каталога) в единственную реализацию openapiv1.ServerInterface. Пакеты
// экспортируют тип с одним и тем же именем Handler, поэтому анонимное
// встраивание невозможно (коллизия имени поля) — здесь явные именованные поля
// и тривиальные методы-делегаты.
type api struct {
	me    *meapiv1.Handler
	login *loginapiv1.Handler

	host                 *hostapiv1.Handler
	databaseType         *databasetypeapiv1.Handler
	databaseCat          *databasecatapiv1.Handler
	alias                *aliasapiv1.Handler
	schemaCat            *schemacatapiv1.Handler
	tableCat             *tablecatapiv1.Handler
	tableType            *tabletypeapiv1.Handler
	columnCat            *columncatapiv1.Handler
	columnType           *columntypeapiv1.Handler
	domainCat            *domaincatapiv1.Handler
	calculationType      *calculationtypeapiv1.Handler
	databaseCalculation  *databasecalculationapiv1.Handler
	followingCalculation *followingcalculationapiv1.Handler
	groupLevels          *grouplevelsapiv1.Handler
	hasToGroup           *hastogroupapiv1.Handler
}

// newAPI принимает клиент datacatalogue целиком, а не по gRPC-клиенту на
// аргумент: таблиц каталога пятнадцать, и позиционный список такой длины
// читается и правится хуже, чем одно поле на таблицу внутри.
func newAPI(resolver *ensureuser.Resolver, dc *datacatalogue.Client) api {
	return api{
		me:    meapiv1.New(resolver),
		login: loginapiv1.New(resolver),

		host:                 hostapiv1.New(dc.Host),
		databaseType:         databasetypeapiv1.New(dc.DatabaseType),
		databaseCat:          databasecatapiv1.New(dc.DatabaseCat),
		alias:                aliasapiv1.New(dc.Alias),
		schemaCat:            schemacatapiv1.New(dc.SchemaCat),
		tableCat:             tablecatapiv1.New(dc.TableCat),
		tableType:            tabletypeapiv1.New(dc.TableType),
		columnCat:            columncatapiv1.New(dc.ColumnCat),
		columnType:           columntypeapiv1.New(dc.ColumnType),
		domainCat:            domaincatapiv1.New(dc.DomainCat),
		calculationType:      calculationtypeapiv1.New(dc.CalculationType),
		databaseCalculation:  databasecalculationapiv1.New(dc.DatabaseCalculation),
		followingCalculation: followingcalculationapiv1.New(dc.FollowingCalculation),
		groupLevels:          grouplevelsapiv1.New(dc.GroupLevels),
		hasToGroup:           hastogroupapiv1.New(dc.HasToGroup),
	}
}

func (a api) GetMe(w http.ResponseWriter, r *http.Request) { a.me.GetMe(w, r) }
func (a api) Login(w http.ResponseWriter, r *http.Request) { a.login.Login(w, r) }

// requireCatalogRole гейтит методы каталога (dc.host и далее) ролью
// admin/maintainer — см. mapping роль → группа методов в rbac.go.
func requireCatalogRole(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rbac.RequireRole(rbac.RoleAdmin, rbac.RoleMaintainer)(next).ServeHTTP(w, r)
	}
}

func (a api) ListHosts(w http.ResponseWriter, r *http.Request, params openapiv1.ListHostsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.host.ListHosts(w, r, params) })(w, r)
}

func (a api) CreateHost(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.host.CreateHost)(w, r)
}

func (a api) SearchHosts(w http.ResponseWriter, r *http.Request, params openapiv1.SearchHostsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.host.SearchHosts(w, r, params) })(w, r)
}

func (a api) ListDeletedHosts(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedHostsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.host.ListDeletedHosts(w, r, params) })(w, r)
}

func (a api) GetHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.host.GetHostById(w, r, id) })(w, r)
}

func (a api) UpdateHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.host.UpdateHostById(w, r, id) })(w, r)
}

func (a api) DeleteHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.host.DeleteHostById(w, r, id) })(w, r)
}

func (a api) UndeleteHostById(w http.ResponseWriter, r *http.Request, id openapiv1.HostIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.host.UndeleteHostById(w, r, id) })(w, r)
}

func (a api) ListDatabaseTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseType.ListDatabaseTypes(w, r, params) })(w, r)
}

func (a api) CreateDatabaseType(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.databaseType.CreateDatabaseType)(w, r)
}

func (a api) SearchDatabaseTypes(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDatabaseTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseType.SearchDatabaseTypes(w, r, params) })(w, r)
}

func (a api) ListDeletedDatabaseTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDatabaseTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseType.ListDeletedDatabaseTypes(w, r, params) })(w, r)
}

func (a api) GetDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseType.GetDatabaseTypeById(w, r, id) })(w, r)
}

func (a api) UpdateDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseType.UpdateDatabaseTypeById(w, r, id) })(w, r)
}

func (a api) DeleteDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseType.DeleteDatabaseTypeById(w, r, id) })(w, r)
}

func (a api) UndeleteDatabaseTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseType.UndeleteDatabaseTypeById(w, r, id) })(w, r)
}

func (a api) ListDatabaseCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseCat.ListDatabaseCats(w, r, params) })(w, r)
}

func (a api) CreateDatabaseCat(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.databaseCat.CreateDatabaseCat)(w, r)
}

func (a api) SearchDatabaseCats(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDatabaseCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseCat.SearchDatabaseCats(w, r, params) })(w, r)
}

func (a api) ListDeletedDatabaseCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDatabaseCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseCat.ListDeletedDatabaseCats(w, r, params) })(w, r)
}

func (a api) ListDatabaseCatsByHostAndDatabaseType(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseCatsByHostAndDatabaseTypeParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCat.ListDatabaseCatsByHostAndDatabaseType(w, r, params)
	})(w, r)
}

func (a api) GetDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseCat.GetDatabaseCatById(w, r, id) })(w, r)
}

func (a api) UpdateDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseCat.UpdateDatabaseCatById(w, r, id) })(w, r)
}

func (a api) DeleteDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseCat.DeleteDatabaseCatById(w, r, id) })(w, r)
}

func (a api) UndeleteDatabaseCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.databaseCat.UndeleteDatabaseCatById(w, r, id) })(w, r)
}

func (a api) ListAliases(w http.ResponseWriter, r *http.Request, params openapiv1.ListAliasesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.alias.ListAliases(w, r, params) })(w, r)
}

func (a api) ListDeletedAliases(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedAliasesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.alias.ListDeletedAliases(w, r, params) })(w, r)
}

func (a api) CreateAlias(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.alias.CreateAlias)(w, r)
}

func (a api) GetAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.alias.GetAliasById(w, r, id) })(w, r)
}

func (a api) UpdateAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.alias.UpdateAliasById(w, r, id) })(w, r)
}

func (a api) DeleteAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.alias.DeleteAliasById(w, r, id) })(w, r)
}

func (a api) UndeleteAliasById(w http.ResponseWriter, r *http.Request, id openapiv1.AliasIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.alias.UndeleteAliasById(w, r, id) })(w, r)
}

func (a api) ListColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListColumnCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.ListColumnCats(w, r, params) })(w, r)
}

func (a api) SearchColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.SearchColumnCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.SearchColumnCats(w, r, params) })(w, r)
}

func (a api) ListDeletedColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedColumnCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.ListDeletedColumnCats(w, r, params) })(w, r)
}

func (a api) FilterColumnCats(w http.ResponseWriter, r *http.Request, params openapiv1.FilterColumnCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.FilterColumnCats(w, r, params) })(w, r)
}

func (a api) CreateColumnCat(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.columnCat.CreateColumnCat)(w, r)
}

func (a api) GetColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.GetColumnCatById(w, r, id) })(w, r)
}

func (a api) UpdateColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.UpdateColumnCatById(w, r, id) })(w, r)
}

func (a api) DeleteColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.DeleteColumnCatById(w, r, id) })(w, r)
}

func (a api) UndeleteColumnCatById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnCat.UndeleteColumnCatById(w, r, id) })(w, r)
}

func (a api) ListColumnTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListColumnTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnType.ListColumnTypes(w, r, params) })(w, r)
}

func (a api) SearchColumnTypes(w http.ResponseWriter, r *http.Request, params openapiv1.SearchColumnTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnType.SearchColumnTypes(w, r, params) })(w, r)
}

func (a api) ListDeletedColumnTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedColumnTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnType.ListDeletedColumnTypes(w, r, params) })(w, r)
}

func (a api) CreateColumnType(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.columnType.CreateColumnType)(w, r)
}

func (a api) GetColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnType.GetColumnTypeById(w, r, id) })(w, r)
}

func (a api) UpdateColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnType.UpdateColumnTypeById(w, r, id) })(w, r)
}

func (a api) DeleteColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnType.DeleteColumnTypeById(w, r, id) })(w, r)
}

func (a api) UndeleteColumnTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.ColumnTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.columnType.UndeleteColumnTypeById(w, r, id) })(w, r)
}

func (a api) ListDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListDatabaseCalculationsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.ListDatabaseCalculations(w, r, params)
	})(w, r)
}

func (a api) SearchDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDatabaseCalculationsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.SearchDatabaseCalculations(w, r, params)
	})(w, r)
}

func (a api) ListDeletedDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDatabaseCalculationsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.ListDeletedDatabaseCalculations(w, r, params)
	})(w, r)
}

func (a api) FilterDatabaseCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.FilterDatabaseCalculationsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.FilterDatabaseCalculations(w, r, params)
	})(w, r)
}

func (a api) CreateDatabaseCalculation(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.databaseCalculation.CreateDatabaseCalculation)(w, r)
}

func (a api) GetDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.GetDatabaseCalculationById(w, r, id)
	})(w, r)
}

func (a api) UpdateDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.UpdateDatabaseCalculationById(w, r, id)
	})(w, r)
}

func (a api) DeleteDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.DeleteDatabaseCalculationById(w, r, id)
	})(w, r)
}

func (a api) UndeleteDatabaseCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.DatabaseCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.databaseCalculation.UndeleteDatabaseCalculationById(w, r, id)
	})(w, r)
}

func (a api) ListDomainCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDomainCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.domainCat.ListDomainCats(w, r, params) })(w, r)
}

func (a api) SearchDomainCats(w http.ResponseWriter, r *http.Request, params openapiv1.SearchDomainCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.domainCat.SearchDomainCats(w, r, params) })(w, r)
}

func (a api) ListDeletedDomainCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedDomainCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.domainCat.ListDeletedDomainCats(w, r, params) })(w, r)
}

func (a api) CreateDomainCat(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.domainCat.CreateDomainCat)(w, r)
}

func (a api) GetDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.domainCat.GetDomainCatById(w, r, id) })(w, r)
}

func (a api) UpdateDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.domainCat.UpdateDomainCatById(w, r, id) })(w, r)
}

func (a api) DeleteDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.domainCat.DeleteDomainCatById(w, r, id) })(w, r)
}

func (a api) UndeleteDomainCatById(w http.ResponseWriter, r *http.Request, id openapiv1.DomainCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.domainCat.UndeleteDomainCatById(w, r, id) })(w, r)
}

func (a api) ListFollowingCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListFollowingCalculationsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.followingCalculation.ListFollowingCalculations(w, r, params)
	})(w, r)
}

func (a api) ListDeletedFollowingCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedFollowingCalculationsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.followingCalculation.ListDeletedFollowingCalculations(w, r, params)
	})(w, r)
}

func (a api) FilterFollowingCalculations(w http.ResponseWriter, r *http.Request, params openapiv1.FilterFollowingCalculationsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.followingCalculation.FilterFollowingCalculations(w, r, params)
	})(w, r)
}

func (a api) CreateFollowingCalculation(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.followingCalculation.CreateFollowingCalculation)(w, r)
}

func (a api) GetFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.followingCalculation.GetFollowingCalculationById(w, r, id)
	})(w, r)
}

func (a api) UpdateFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.followingCalculation.UpdateFollowingCalculationById(w, r, id)
	})(w, r)
}

func (a api) DeleteFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.followingCalculation.DeleteFollowingCalculationById(w, r, id)
	})(w, r)
}

func (a api) UndeleteFollowingCalculationById(w http.ResponseWriter, r *http.Request, id openapiv1.FollowingCalculationIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) {
		a.followingCalculation.UndeleteFollowingCalculationById(w, r, id)
	})(w, r)
}

func (a api) ListGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.ListGroupLevelsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.ListGroupLevels(w, r, params) })(w, r)
}

func (a api) SearchGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.SearchGroupLevelsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.SearchGroupLevels(w, r, params) })(w, r)
}

func (a api) ListDeletedGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedGroupLevelsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.ListDeletedGroupLevels(w, r, params) })(w, r)
}

func (a api) FilterGroupLevels(w http.ResponseWriter, r *http.Request, params openapiv1.FilterGroupLevelsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.FilterGroupLevels(w, r, params) })(w, r)
}

func (a api) CreateGroupLevels(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.groupLevels.CreateGroupLevels)(w, r)
}

func (a api) GetGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.GetGroupLevelsById(w, r, id) })(w, r)
}

func (a api) UpdateGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.UpdateGroupLevelsById(w, r, id) })(w, r)
}

func (a api) DeleteGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.DeleteGroupLevelsById(w, r, id) })(w, r)
}

func (a api) UndeleteGroupLevelsById(w http.ResponseWriter, r *http.Request, id openapiv1.GroupLevelsIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.groupLevels.UndeleteGroupLevelsById(w, r, id) })(w, r)
}

func (a api) ListHasToGroups(w http.ResponseWriter, r *http.Request, params openapiv1.ListHasToGroupsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.hasToGroup.ListHasToGroups(w, r, params) })(w, r)
}

func (a api) SearchHasToGroups(w http.ResponseWriter, r *http.Request, params openapiv1.SearchHasToGroupsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.hasToGroup.SearchHasToGroups(w, r, params) })(w, r)
}

func (a api) ListDeletedHasToGroups(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedHasToGroupsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.hasToGroup.ListDeletedHasToGroups(w, r, params) })(w, r)
}

func (a api) CreateHasToGroup(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.hasToGroup.CreateHasToGroup)(w, r)
}

func (a api) GetHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.hasToGroup.GetHasToGroupById(w, r, id) })(w, r)
}

func (a api) UpdateHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.hasToGroup.UpdateHasToGroupById(w, r, id) })(w, r)
}

func (a api) DeleteHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.hasToGroup.DeleteHasToGroupById(w, r, id) })(w, r)
}

func (a api) UndeleteHasToGroupById(w http.ResponseWriter, r *http.Request, id openapiv1.HasToGroupIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.hasToGroup.UndeleteHasToGroupById(w, r, id) })(w, r)
}

func (a api) ListSchemaCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListSchemaCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.schemaCat.ListSchemaCats(w, r, params) })(w, r)
}

func (a api) ListDeletedSchemaCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedSchemaCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.schemaCat.ListDeletedSchemaCats(w, r, params) })(w, r)
}

func (a api) CreateSchemaCat(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.schemaCat.CreateSchemaCat)(w, r)
}

func (a api) GetSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.schemaCat.GetSchemaCatById(w, r, id) })(w, r)
}

func (a api) UpdateSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.schemaCat.UpdateSchemaCatById(w, r, id) })(w, r)
}

func (a api) DeleteSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.schemaCat.DeleteSchemaCatById(w, r, id) })(w, r)
}

func (a api) UndeleteSchemaCatById(w http.ResponseWriter, r *http.Request, id openapiv1.SchemaCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.schemaCat.UndeleteSchemaCatById(w, r, id) })(w, r)
}

func (a api) ListTableCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListTableCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableCat.ListTableCats(w, r, params) })(w, r)
}

func (a api) ListDeletedTableCats(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedTableCatsParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableCat.ListDeletedTableCats(w, r, params) })(w, r)
}

func (a api) CreateTableCat(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.tableCat.CreateTableCat)(w, r)
}

func (a api) GetTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableCat.GetTableCatById(w, r, id) })(w, r)
}

func (a api) UpdateTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableCat.UpdateTableCatById(w, r, id) })(w, r)
}

func (a api) DeleteTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableCat.DeleteTableCatById(w, r, id) })(w, r)
}

func (a api) UndeleteTableCatById(w http.ResponseWriter, r *http.Request, id openapiv1.TableCatIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableCat.UndeleteTableCatById(w, r, id) })(w, r)
}

func (a api) ListTableTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListTableTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableType.ListTableTypes(w, r, params) })(w, r)
}

func (a api) ListDeletedTableTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListDeletedTableTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableType.ListDeletedTableTypes(w, r, params) })(w, r)
}

func (a api) CreateTableType(w http.ResponseWriter, r *http.Request) {
	requireCatalogRole(a.tableType.CreateTableType)(w, r)
}

func (a api) GetTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableType.GetTableTypeById(w, r, id) })(w, r)
}

func (a api) UpdateTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableType.UpdateTableTypeById(w, r, id) })(w, r)
}

func (a api) DeleteTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableType.DeleteTableTypeById(w, r, id) })(w, r)
}

func (a api) UndeleteTableTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.TableTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.tableType.UndeleteTableTypeById(w, r, id) })(w, r)
}

func (a api) ListCalculationTypes(w http.ResponseWriter, r *http.Request, params openapiv1.ListCalculationTypesParams) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.calculationType.ListCalculationTypes(w, r, params) })(w, r)
}

func (a api) GetCalculationTypeById(w http.ResponseWriter, r *http.Request, id openapiv1.CalculationTypeIdParam) {
	requireCatalogRole(func(w http.ResponseWriter, r *http.Request) { a.calculationType.GetCalculationTypeById(w, r, id) })(w, r)
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
