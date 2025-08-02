package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/konstantin-suspitsyn/datacomrade/internal/healthcheck"
	"github.com/konstantin-suspitsyn/datacomrade/internal/services"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
)

func routes(services *services.ServiceLayer) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(time.Second * 60))

	// Sends 500 status if something goes south
	r.Use(middleware.Recoverer)

	r.Use(httprate.LimitByIP(100, time.Minute*1))
	r.Use(middleware.CleanPath)

	// Sends custom 404, 405 responces
	r.NotFound(custresponse.NotFoundResponse)
	r.MethodNotAllowed(custresponse.MethodNotAllowed)

	r.Get("/healthcheck", healthcheck.ReturnOk)

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", services.UserService.UserRegister)
		r.Put("/activate", services.UserService.UserActivate)
	})
	return r
}
