package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

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

var (
	start int = 1
	end   int = 1
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

func (p *Server) PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	short, err := u.GetSquidShort(uint64(start))
	if err != nil {
		slog.Error("failed to generate short url", "start", start, "error", err)
		apiUtils.InternalServerError(w)
		return
	}
	start++

	if err := redis.StoreUrl(short, url, 0); err != nil {
		slog.Error("failed to store url in redis", "short", short, "url", url, "error", err)
		apiUtils.InternalServerError(w)
		return
	}

	payload, err := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		apiUtils.InternalServerError(w)
		return
	}

	usrInfo, err := oidc.GetUsrInfoFromCtx(r.Context())
	if err != nil {
		slog.Warn("Failed to get UserInfo", "error", err)
	}

	shortURL := m.ShortURL{URL: url, Short: short}
	if usrInfo != nil {
		shortURL.Owner = usrInfo.Email
	}

	if _, storeErr := mongo.StoreShortURL(shortURL); storeErr != nil {
		slog.Error("database error storing short url", "short", short, "url", url, "error", storeErr)
	}

	if _, err := w.Write(payload); err != nil {
		slog.Error("failed to write HTTP response", "error", err)
	}

	if start > end {
		allocateRange()
	}
}

func allocateRange() {
	_start, _end, err := u.AllocateRange()
	if err != nil {
		slog.Error("range allocation failed", "error", err)
		slog.Info("Exiting... range exceeded")
		os.Exit(1)
	}
	slog.Info("Range allocated", "start", _start, "end", _end)
	start = _start
	end = _end
}

func main() {
	slog.Info("*_-_-_-BlyLi-Shortn-_-_-_*")

	// Set up database and other resources before listen
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
	api.HandlerFromMux(server, r)

	allocateRange()

	if err := http.ListenAndServe(":8082", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("API server exited with error", "error", err)
		os.Exit(1)
	}
}
