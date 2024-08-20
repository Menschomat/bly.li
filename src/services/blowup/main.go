package main

import (
	"log"
	"net/http"

	"github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/utils"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"

	"github.com/Menschomat/bly.li/services/blowup/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

// GetShort FindPets implements all the handlers in the ServerInterface
func (p *Server) GetShort(w http.ResponseWriter, r *http.Request, short string) {
	if utils.IsValidShort(short) {
		url, err := redis.GetUrl(short)
		if err != nil || len(url) == 0 {
			var shortInfo *model.ShortURL
			shortInfo, err = mongo.GetShortURLByShort(short)
			if err == nil {
				redis.StoreUrl(shortInfo.Short, shortInfo.URL)
				url = shortInfo.URL
			}
		}
		if err == nil && len(url) > 0 {
			go redis.RegisterClick(short)
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
	apiUtils.BadRequestError(w)
}

func main() {
	log.Println("*_-_-_-BlyLi-Blowup-_-_-_*")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	server := &Server{}
	api.HandlerFromMux(server, r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
}
