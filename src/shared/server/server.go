package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mw "github.com/Menschomat/bly.li/shared/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config holds server configuration
type Config struct {
	ServerPort         string
	MetricsPort        string
	CorsAllowedOrigins []string
	CorsMaxAge         int
	Logger             *slog.Logger
}

// Server represents a generic HTTP server with common functionality
type Server struct {
	config Config
	router *chi.Mux
}

// NewServer creates a new server instance with the given configuration
func NewServer(config Config) *Server {
	return &Server{
		config: config,
		router: chi.NewRouter(),
	}
}

// Router returns the underlying chi router for custom route configuration
func (s *Server) Router() *chi.Mux {
	return s.router
}

// ConfigureCommonMiddleware sets up common middleware used across services
func (s *Server) ConfigureCommonMiddleware() {
	s.router.Use(mw.SlogLogger(s.config.Logger))
	s.router.Use(mw.InstrumentHandler)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.config.CorsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           s.config.CorsMaxAge,
	}))
}

// ServeMetrics starts the Prometheus metrics endpoint
func (s *Server) ServeMetrics(errChan chan<- error) {
	metricsRouter := chi.NewRouter()
	metricsRouter.Handle("/metrics", promhttp.Handler())

	s.config.Logger.Info("Prometheus metrics available on " + s.config.MetricsPort + "/metrics")
	errChan <- http.ListenAndServe(":"+s.config.MetricsPort, metricsRouter)
}

// ServeHTTP starts the main HTTP server
func (s *Server) ServeHTTP(errChan chan<- error) {
	s.config.Logger.Info("Backend runs on " + s.config.ServerPort)
	errChan <- http.ListenAndServe(":"+s.config.ServerPort, s.router)
}

// HandleShutdown waits for server errors or shutdown signals and handles shutdown logic
func (s *Server) HandleShutdown(ctx context.Context, serverErrChan <-chan error) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		if err != nil {
			s.config.Logger.Error("Server error", "error", err)
		}
	case <-stopChan:
		s.config.Logger.Info("Shutdown signal received. Stopping server...")
	case <-ctx.Done():
		s.config.Logger.Info("Context cancelled, shutting down.")
	}
}
