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

	// http handlers
	r.Get("/", index)
	r.Get("/healthz", healthz)

	return r
}
