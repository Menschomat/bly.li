package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Menschomat/bly.li/services/blowup/api"
	"github.com/Menschomat/bly.li/services/blowup/logging"
	mw "github.com/Menschomat/bly.li/shared/api/middleware"
	"github.com/Menschomat/bly.li/shared/config"
	"github.com/Menschomat/bly.li/shared/data"
	"github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/utils"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	clicksRegistered = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "blyli_click_registered_total",
		Help: "Total number of clicks handled by blowup",
	})
	shortNotFound = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "blyli_short_not_found_total",
		Help: "Total number of not found short requests",
	})

	logger                     = logging.GetLogger()
	cfg                        = config.BlowupConfig()
	_      api.ServerInterface = (*server)(nil)
)

// server implements api.ServerInterface.
type server struct{}

// getShort handles short code redirection logic.
func (s *server) GetShort(w http.ResponseWriter, r *http.Request, short string) {
	if !utils.IsValidShort(short) {
		apiUtils.BadRequestError(w)
		return
	}

	shortURL := data.GetShort(short)
	if shortURL == nil {
		shortNotFound.Inc()
		apiUtils.BadRequestError(w)
		return
	}

	ip := apiUtils.ReadUserIP(r)
	userAgent := r.UserAgent()
	go redis.RegisterClick(model.ShortClick{
		Short:     shortURL.Short,
		Ip:        ip,
		UsrAgent:  userAgent,
		Timestamp: time.Now(),
	})

	clicksRegistered.Inc()
	http.Redirect(w, r, shortURL.URL, cfg.RedirectCode)
}

func main() {
	logger.Info("Starting service")
	mongo.InitMongoPackage(logger)

	mainRouter := configureMainRouter()
	apiHandler := &server{}
	api.HandlerFromMux(apiHandler, mainRouter)

	serverErrChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go serveMetrics(serverErrChan)
	go serveMainHTTP(serverErrChan, mainRouter)

	handleShutdown(ctx, serverErrChan)
	logger.Info("Server shut down successfully.")
}

// configureMainRouter initialises and configures the HTTP router.
func configureMainRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(mw.SlogLogger(logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	return r
}

// serveMetrics starts the Prometheus metrics endpoint.
func serveMetrics(errChan chan<- error) {
	prometheus.MustRegister(clicksRegistered)
	prometheus.MustRegister(shortNotFound)

	metricsRouter := chi.NewRouter()
	metricsRouter.Handle("/metrics", promhttp.Handler())

	logger.Info("Prometheus metrics available on " + cfg.MetricsPort + "/metrics")
	errChan <- http.ListenAndServe(cfg.MetricsPort, metricsRouter)
}

// serveMainHTTP starts the main HTTP API.
func serveMainHTTP(errChan chan<- error, handler http.Handler) {
	logger.Info("Backend runs on " + cfg.ServerPort)
	errChan <- http.ListenAndServe(cfg.ServerPort, handler)
}

// handleShutdown waits for server errors or shutdown signals and handles shutdown logic.
func handleShutdown(ctx context.Context, serverErrChan <-chan error) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		if err != nil {
			logger.Error("Server error", "error", err)
		}
	case <-stopChan:
		logger.Info("Shutdown signal received. Stopping server...")
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down.")
	}
}
