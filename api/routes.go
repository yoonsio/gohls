package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func newRouter() *chi.Mux {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// top level http handlers
	r.Get("/", index)
	r.Get("/healthz", healthz)

	// api http handlers
	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/v1", func(apiV1Router chi.Router) {
			apiV1Router.Use(apiVersionCtx("v1"))
			r.Get("/", index)
		})
	})

	return r
}
