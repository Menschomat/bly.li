package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Menschomat/bly.li/services/shortn/api"
	"github.com/Menschomat/bly.li/services/shortn/logging"
	u "github.com/Menschomat/bly.li/services/shortn/utils"
	mw "github.com/Menschomat/bly.li/shared/api/middleware"
	"github.com/Menschomat/bly.li/shared/config"
	m "github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/oidc"
	"github.com/Menschomat/bly.li/shared/redis"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	logger = logging.GetLogger()
	cfg    = config.ShortnConfig()
)

/* -------------------------------------------------------------------- */
/*  Server                                                              */
/* -------------------------------------------------------------------- */

type Server struct {
	mu    sync.Mutex // guards start & end
	start int        // next id to hand out
	end   int        // inclusive upper bound of current range
}

/* ------------------------- range management ------------------------- */

func (s *Server) allocateRangeLocked() {
	// called with s.mu held
	_start, _end, err := u.AllocateRange()
	if err != nil {
		logger.Error("range allocation failed", "error", err)
		logger.Info("Exiting… range exceeded")
		os.Exit(1)
	}
	logger.Info("Range allocated", "start", _start, "end", _end)
	s.start = _start
	s.end = _end
}

func (s *Server) nextShort() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.start > s.end { // out of numbers? → fetch a new block
		s.allocateRangeLocked()
	}
	id := s.start
	s.start++

	return u.GetSquidShort(uint64(id))
}

/* ---------------------------- handlers ------------------------------ */

var _ api.ServerInterface = (*Server)(nil)

func (s *Server) PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	/* ----------- parse body ------------------------------------------------ */

	var shortn m.ShortnReq
	if err := json.NewDecoder(r.Body).Decode(&shortn); err != nil {
		logger.Error("invalid request payload", "error", err)
		apiUtils.BadRequestError(w)
		return
	}

	url, err := u.ParseUrl(shortn.Url)
	if err != nil {
		logger.Warn("invalid url in request", "url", shortn.Url, "error", err)
		apiUtils.BadRequestError(w)
		return
	}

	/* ----------- generate short code (thread-safe!) ------------------------ */

	short, err := s.nextShort()
	if err != nil {
		logger.Error("failed to generate short url", "error", err)
		apiUtils.InternalServerError(w)
		return
	}

	/* ----------- persist --------------------------------------------------- */

	usrInfo, _ := oidc.GetUsrInfoFromCtx(r.Context()) // ignore "no user" error
	shortURL := m.ShortURL{URL: url, Short: short, Owner: "", Count: 0}
	if usrInfo != nil {
		shortURL.Owner = usrInfo.Email
	}

	if err := redis.StoreUrl(shortURL.Short, shortURL.URL, shortURL.Count, shortURL.Owner); err != nil {
		logger.Error("failed to store url in redis", "short", short, "url", url, "error", err)
		apiUtils.InternalServerError(w)
		return
	}
	if _, err := mongo.StoreShortURL(shortURL); err != nil {
		logger.Error("database error storing short url", "short", short, "url", url, "error", err)
	}

	/* ----------- respond --------------------------------------------------- */

	payload, _ := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if _, err := w.Write(payload); err != nil {
		logger.Error("failed to write HTTP response", "error", err)
	}
}

/* -------------------------------------------------------------------- */
/*  main                                                                */
/* -------------------------------------------------------------------- */

func main() {
	logger.Info("Starting shortn service")
	mongo.InitMongoPackage(logger)
	defer mongo.CloseClientDB()

	mainRouter := configureMainRouter()
	server := &Server{}
	server.allocateRangeLocked() // grab the first block before serving
	api.HandlerFromMux(server, mainRouter)

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
	r.Use(mw.SlogLogger(logger))
	r.Use(mw.InstrumentHandler)
	r.Use(oidc.JWTVerifier)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(cfg.CorsAllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           cfg.CorsMaxAge,
	}))
	return r
}

// serveMetrics starts the Prometheus metrics endpoint.
func serveMetrics(errChan chan<- error) {
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
