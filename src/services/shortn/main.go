package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/Menschomat/bly.li/services/shortn/api"
	u "github.com/Menschomat/bly.li/services/shortn/utils"
	m "github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/oidc"
	"github.com/Menschomat/bly.li/shared/redis"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
		slog.Error("range allocation failed", "error", err)
		slog.Info("Exiting… range exceeded")
		os.Exit(1)
	}
	slog.Info("Range allocated", "start", _start, "end", _end)
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
		slog.Error("invalid request payload", "error", err)
		apiUtils.BadRequestError(w)
		return
	}

	url, err := u.ParseUrl(shortn.Url)
	if err != nil {
		slog.Warn("invalid url in request", "url", shortn.Url, "error", err)
		apiUtils.BadRequestError(w)
		return
	}

	/* ----------- generate short code (thread-safe!) ------------------------ */

	short, err := s.nextShort()
	if err != nil {
		slog.Error("failed to generate short url", "error", err)
		apiUtils.InternalServerError(w)
		return
	}

	/* ----------- persist --------------------------------------------------- */

	if err := redis.StoreUrl(short, url, 0); err != nil {
		slog.Error("failed to store url in redis", "short", short, "url", url, "error", err)
		apiUtils.InternalServerError(w)
		return
	}

	usrInfo, _ := oidc.GetUsrInfoFromCtx(r.Context()) // ignore “no user” error
	shortURL := m.ShortURL{URL: url, Short: short}
	if usrInfo != nil {
		shortURL.Owner = usrInfo.Email
	}
	if _, err := mongo.StoreShortURL(shortURL); err != nil {
		slog.Error("database error storing short url", "short", short, "url", url, "error", err)
	}

	/* ----------- respond --------------------------------------------------- */

	payload, _ := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if _, err := w.Write(payload); err != nil {
		slog.Error("failed to write HTTP response", "error", err)
	}
}

/* -------------------------------------------------------------------- */
/*  main                                                                */
/* -------------------------------------------------------------------- */

func main() {
	slog.Info("*_-_-_-BlyLi-Shortn-_-_-_*")

	defer mongo.CloseClientDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(oidc.JWTVerifier)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	server := &Server{}
	server.allocateRangeLocked() // grab the first block before serving

	api.HandlerFromMux(server, r)

	if err := http.ListenAndServe(":8082", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("API server exited with error", "error", err)
		os.Exit(1)
	}
}
