package main

import (
	"fmt"

	"log/slog"
	"net/http"

	"github.com/Menschomat/bly.li/services/dasher/api"

	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/oidc"
	"github.com/Menschomat/bly.li/shared/redis"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

func (p *Server) DeleteShortShort(w http.ResponseWriter, r *http.Request, short string) {
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

func (p *Server) GetShortAll(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-Auth-User")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Hello")
	if err != nil {
		return
	}
	slog.Debug(user)
}

func main() {

	// Set up OIDC provider and OAuth2 config
	slog.Info("*_-_-_-BlyLi-Dasher-_-_-_*")
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
	err := http.ListenAndServe(":8083", r)
	if err != nil {
		slog.Error("There's an error with the server", "error", err)
	}
	defer mongo.CloseClientDB()
}
