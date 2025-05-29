package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mw "github.com/Menschomat/bly.li/shared/api/middleware"
	"github.com/Menschomat/bly.li/shared/data"
	"github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/utils"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Menschomat/bly.li/services/blowup/api"
	"github.com/Menschomat/bly.li/services/blowup/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	registedClicks = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "click_registered_total",
		Help: "Total number of clicks handled by blowup",
	})
	logger                     = logging.GetLogger()
	_      api.ServerInterface = (*Server)(nil)
)

type Server struct{}

// GetShort FindPets implements all the handlers in the ServerInterface
func (p *Server) GetShort(w http.ResponseWriter, r *http.Request, short string) {
	if utils.IsValidShort(short) {

		url := data.GetShort(short)
		if url == nil {
			apiUtils.BadRequestError(w)
			return
		}
		ip := apiUtils.ReadUserIP(r)
		userAgent := r.UserAgent()
		go redis.RegisterClick(model.ShortClick{Short: url.Short, Ip: ip, UsrAgent: userAgent, Timestamp: time.Now()})
		registedClicks.Inc() //Increment clickCounter
		w.Header().Set("Location", url.URL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}

}

func main() {
	logger.Info("Starting")

	// Initialize router
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

	// Graceful shutdown handling
	server := &Server{}
	api.HandlerFromMux(server, r)
	serverErrChan := make(chan error, 1)

	// --- Metrics router on a 2nd port ---
	go func() {
		prometheus.MustRegister(registedClicks)
		metricsRouter := chi.NewRouter()
		metricsRouter.Handle("/metrics", promhttp.Handler())
		logger.Info("Prometheus metrics available on :2114/metrics")
		serverErrChan <- http.ListenAndServe(":2114", metricsRouter)
	}()

	// --- Main server on 8083 ---
	// HTTP server in a goroutine
	go func() {
		logger.Info("Backend runs on :8081")
		serverErrChan <- http.ListenAndServe(":8081", r)
	}()

	// Handle shutdown signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		logger.Error("Server error", "error", err)
	case <-stopChan:
		logger.Info("Shutdown signal received. Stopping server...")
	}

	logger.Info("Server shut down successfully.")
}
