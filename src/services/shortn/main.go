package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	cfgUtils "github.com/Menschomat/bly.li/shared/utils/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	appConfig m.ShortnConfig
	start     int = 1
	end       int = 1
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

func (p *Server) DeleteShort(w http.ResponseWriter, r *http.Request, short string) {
	if redis.ShortExists(short) {
		err := redis.DeleteUrl(short)
		if err != nil {
			return
		}
	}
	if mongo.ShortExists(short) {
		err := mongo.DeleteShortURLByShort(short)
		if err != nil {
			return
		}
	}
}

func (p *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-Auth-User")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Hello")
	if err != nil {
		return
	}
	slog.Debug(user)
}

func (p *Server) PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortn m.ShortnReq
	err := json.NewDecoder(r.Body).Decode(&shortn)
	url, err := u.ParseUrl(shortn.Url)
	if err != nil {
		apiUtils.BadRequestError(w)
		return
	}
	short, err := u.GetSquidShort(uint64(start))
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
	start++
	redis.StoreUrl(short, url)
	payload, err := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if err != nil {
		apiUtils.InternalServerError(w)
		return
	}
	usrInfo, err := oidc.GetUsrInfoFromCtx(r.Context())
	if usrInfo != nil {
		_, err = mongo.StoreShortURL(m.ShortURL{URL: url, Short: short, Owner: usrInfo.Email})
	}
	if err != nil {
		_, err = mongo.StoreShortURL(m.ShortURL{URL: url, Short: short})
	}
	_, err = w.Write(payload)
	if start > end {
		slog.Info("Exiting... range exeeded")
		os.Exit(0)
	}
}

func main() {
	err := cfgUtils.FillEnvStruct(&appConfig)
	if err != nil {
		slog.Error("There's an error with the Config", "error", err)
	}
	// Set up OIDC provider and OAuth2 config
	slog.Info("*_-_-_-BlyLi-Shortn-_-_-_*")
	// Create new Chi-Router
	r := chi.NewRouter()
	// Add Middlewares to Router
	r.Use(middleware.Logger)
	r.Use(oidc.JWTVerifier)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	server := &Server{}
	api.HandlerFromMux(server, r)
	conn := u.CreateZkConnection()
	defer conn.Close()
	start, end, err = u.AllocateRange(conn)
	if err != nil {
		slog.Error("There's an error with the range", "error", err)
	}
	slog.Info("Range", "start", start, "end", end)
	err = http.ListenAndServe(":8082", r)
	if err != nil {
		slog.Error("There's an error with the server", "error", err)
	}
	defer mongo.CloseClientDB()
}
