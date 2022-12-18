package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	redis "github.com/Menschomat/bly.li/shared/redis"
	utils "github.com/Menschomat/bly.li/shared/utils"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func blowUpShortToUrl(w http.ResponseWriter, r *http.Request) {
	var short string = chi.URLParam(r, "short")
	if utils.IsValidShort(short) {
		url := redis.GetUrl(short)
		if len(url) > 0 {
			go redis.RegisterClick(short)
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
	apiUtils.BadRequestError(w, r)

}

func main() {
	log.Println("*_-_-_-BlyLi-Blowup-_-_-_*")
	rand.Seed(time.Now().UnixNano())
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{short}", blowUpShortToUrl)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
}
