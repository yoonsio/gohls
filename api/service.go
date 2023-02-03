package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	// DefaultPort used if port is not configured
	DefaultPort = 3000
)

// Service represents HTTP API Service
type Service struct {
	chi.Router
	Port int
}

// NewService returns new HTTP API Service
func NewService() *Service {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// http handlers
	r.Get("/", index)
	r.Get("/healthz", healthz)

	svc := &Service{
		Router: r,
	}
	return svc
}

// Run runs the API Service
func (s *Service) Run() error {

	// validate arguments
	if s.Port == 0 {
		s.Port = DefaultPort
	}

	// create new server
	addr := fmt.Sprintf("0.0.0.0:%d", s.Port)
	svr := &http.Server{
		Addr:    addr,
		Handler: s,
	}

	// start server
	log.Printf("starting server at %s\n", addr)
	return svr.ListenAndServe()
}
