package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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
		AllowCredentials: false,
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

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", services.UserService.UserRegister)
		r.Get("/refresh", services.UserService.GetAccessTokenByRefresh)
		r.Put("/activate", services.UserService.UserActivate)
		r.Post("/login", services.UserService.UserLogin)
		r.With(IsAuthorized).Get("/me", services.UserService.Me)
	})
	return r
}
