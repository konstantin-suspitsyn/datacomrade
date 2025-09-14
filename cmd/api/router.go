package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/internal/healthcheck"
	"github.com/konstantin-suspitsyn/datacomrade/internal/services"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
)

func routes(services *services.ServiceLayer) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Timeout(time.Second * 60))

	// Sends 500 status if something goes south
	r.Use(middleware.Recoverer)

	r.Use(httprate.LimitByIP(100, time.Minute*1))
	r.Use(middleware.CleanPath)

	// Sends custom 404, 405 responces
	r.NotFound(custresponse.NotFoundResponse)
	r.MethodNotAllowed(custresponse.MethodNotAllowed)

	r.Use(GetAuthMiddlewareFunc(services))

	r.Get("/healthcheck", healthcheck.ReturnOk)

	r.Get("/refresh", services.UserService.GetAccessTokenByRefresh)

	r.Route(configs.USERS_V1, func(r chi.Router) {
		r.Post(configs.INDEX_PAGE_LINK, services.UserService.UserRegister)
		r.Get(configs.REFRESH_JWT_LINK, services.UserService.GetAccessTokenByRefresh)
		r.Delete(configs.REFRESH_JWT_LINK, services.UserService.UserLogout)
		r.Put(configs.ACTIVATE_LINK, services.UserService.UserActivate)
		r.Post(configs.ACTIVATE_LINK, services.UserService.UserLogin)
		r.With(IsAuthorized).Get("/me", services.UserService.Me)
	})

	r.Route(configs.DOMAIN_LINK, func(r chi.Router) {
		r.With(IsAdmin).Get(configs.GET_DOMAIN, services.SharedDataService.GetAllDomains)
	})

	r.Route(configs.ROLES_LINK, func(r chi.Router) {
		r.With(IsAdmin).Get(configs.GET_DOMAIN, services.RoleService.GetAllDomains)
	})
	return r
}
