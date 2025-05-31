package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Menschomat/bly.li/services/shortn/api"
	"github.com/Menschomat/bly.li/services/shortn/logging"
	u "github.com/Menschomat/bly.li/services/shortn/utils"
	"github.com/Menschomat/bly.li/shared/config"
	m "github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/oidc"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/server"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
)

var (
	logger = logging.GetLogger()
	cfg    = config.ShortnConfig()
)

/* -------------------------------------------------------------------- */
/*  Server                                                              */
/* -------------------------------------------------------------------- */

// ShortnServer handles URL shortening and range allocation
type ShortnServer struct {
	mu    sync.Mutex // guards start & end
	start int        // next id to hand out
	end   int        // inclusive upper bound of current range
}

/* ------------------------- range management ------------------------- */

func (s *ShortnServer) allocateRangeLocked() {
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

func (s *ShortnServer) nextShort() (string, error) {
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

var _ api.ServerInterface = (*ShortnServer)(nil)

func (s *ShortnServer) PostStore(w http.ResponseWriter, r *http.Request) {
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

	srv := server.NewServer(server.Config{
		ServerPort:         cfg.ServerPort,
		MetricsPort:        cfg.MetricsPort,
		CorsAllowedOrigins: strings.Split(cfg.CorsAllowedOrigins, ","),
		CorsMaxAge:         cfg.CorsMaxAge,
		Logger:             logger,
	})

	srv.ConfigureCommonMiddleware()
	srv.Router().Use(oidc.JWTVerifier)

	apiServer := &ShortnServer{}
	apiServer.allocateRangeLocked() // grab the first block before serving
	api.HandlerFromMux(apiServer, srv.Router())

	serverErrChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.ServeMetrics(serverErrChan)
	go srv.ServeHTTP(serverErrChan)

	srv.HandleShutdown(ctx, serverErrChan)
	logger.Info("Server shut down successfully.")
}
