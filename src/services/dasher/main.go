package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Menschomat/bly.li/services/dasher/api"
	"github.com/Menschomat/bly.li/services/dasher/logging"
	mw "github.com/Menschomat/bly.li/shared/api/middleware"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/oidc"
	"github.com/Menschomat/bly.li/shared/redis"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var (
	logger                     = logging.GetLogger()
	_      api.ServerInterface = (*Server)(nil)
)

type Server struct{}

func checkOwner(short string) (string, error) {
	u, err := redis.GetShort(short)
	if err == nil && u != nil && u.Owner != "" {
		return u.Owner, nil
	}
	u, err = mongo.GetShortURLByShort(short)
	if err == nil && u != nil && u.Owner != "" {
		return u.Owner, nil
	}
	return "", errors.New("owner not found")
}

func (p *Server) DeleteShortShort(w http.ResponseWriter, r *http.Request, short string) {
	usrInfo, err := oidc.GetUsrInfoFromCtx(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	owner, err := checkOwner(short)
	if err != nil || owner != usrInfo.Email {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := redis.DeleteUrl(short); err != nil && redis.ShortExists(short) {
		http.Error(w, "Failed to delete from Redis", http.StatusInternalServerError)
		return
	}
	if err := mongo.DeleteShortURLByShort(short); err != nil && mongo.ShortExists(short) {
		http.Error(w, "Failed to delete from Mongo", http.StatusInternalServerError)
		return
	}
}
func (p *Server) GetShortAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usrInfo, err := oidc.GetUsrInfoFromCtx(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	shorts := mongo.GetShortsByOwner(usrInfo.Email)
	payload, err := json.Marshal(shorts)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func main() {
	logger.Info("Starting")
	mongo.InitMongoPackage(logger)
	r := chi.NewRouter()
	r.Use(mw.SlogLogger(logger))
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
	if err := http.ListenAndServe(":8083", r); err != nil {
		logger.Error("There's an error with the server", "error", err)
		return
	}
	defer mongo.CloseClientDB()
}
