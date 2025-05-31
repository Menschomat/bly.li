package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Menschomat/bly.li/services/blowup/api"
	"github.com/Menschomat/bly.li/services/blowup/logging"
	"github.com/Menschomat/bly.li/shared/config"
	"github.com/Menschomat/bly.li/shared/data"
	"github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/server"
	"github.com/Menschomat/bly.li/shared/utils"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
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
	_      api.ServerInterface = (*BlowupServer)(nil)
)

// BlowupServer handles URL redirection and click tracking
type BlowupServer struct{}

func (s *BlowupServer) GetShort(w http.ResponseWriter, r *http.Request, short string) {
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

	prometheus.MustRegister(clicksRegistered)
	prometheus.MustRegister(shortNotFound)

	srv := server.NewServer(server.Config{
		ServerPort:         cfg.ServerPort,
		MetricsPort:        cfg.MetricsPort,
		CorsAllowedOrigins: strings.Split(cfg.CorsAllowedOrigins, ","),
		CorsMaxAge:         cfg.CorsMaxAge,
		Logger:             logger,
	})

	srv.ConfigureCommonMiddleware()
	srv.Router().Use(middleware.Recoverer)

	apiHandler := &BlowupServer{}
	api.HandlerFromMux(apiHandler, srv.Router())

	serverErrChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.ServeMetrics(serverErrChan)
	go srv.ServeHTTP(serverErrChan)

	srv.HandleShutdown(ctx, serverErrChan)
	logger.Info("Server shut down successfully.")
}
