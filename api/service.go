package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// DefaultAddress is default listening address if address is not configured
	DefaultAddress = "0.0.0.0:3000"
)

// Service represents HTTP API Service
type Service struct {
	http.Handler

	// config variables
	Addr string
}

// NewService returns new HTTP API Service
func NewService() *Service {
	svc := &Service{
		Handler: newRouter(),
	}
	return svc
}

// Run runs the API Service
func (s *Service) Run(ctx context.Context) {

	// validate arguments
	addr := s.Addr
	if addr == "" {
		addr = DefaultAddress
	}

	// create new server
	server := &http.Server{
		Addr:    addr,
		Handler: s.Handler,
	}

	serverCtx, serverStop := context.WithCancel(ctx)

	// listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		log.Printf("server shutting down...")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStop()
	}()

	// start server
	log.Printf("starting server at %s\n", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// wait for server to shut down
	<-serverCtx.Done()
}
