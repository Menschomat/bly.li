package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	redis "github.com/Menschomat/bly.li/shared/redis"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func getLong(w http.ResponseWriter, r *http.Request) {
	var short string = chi.URLParam(r, "short")
	w.Header().Set("Location", redis.GetUrl(short))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	log.Println("*_-_-_-BlyLi-Blowup-_-_-_*")
	rand.Seed(time.Now().UnixNano())
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{short}", getLong)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
}
