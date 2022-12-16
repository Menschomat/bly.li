package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	m "github.com/Menschomat/bly.li/model"
	redis "github.com/Menschomat/bly.li/redis"
	u "github.com/Menschomat/bly.li/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func getLong(w http.ResponseWriter, r *http.Request) {
	var short string = chi.URLParam(r, "short")
	w.Header().Set("Location", redis.GetUrl(short))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortn m.ShortnReq
	json.NewDecoder(r.Body).Decode(&shortn)
	url, err := u.ParseUrl(shortn.Url)
	if err != nil {
		u.BadRequestError(w, r)
		return
	}
	short := u.GetUniqueShort()
	redis.StoreUrl(short, url)
	payload, err := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if err != nil {
		u.InternalServerError(w, r)
		return
	}
	w.Write(payload)
}

func main() {
	log.Println("*_-_-_-Bly.li-_-_-_*")
	rand.Seed(time.Now().UnixNano())
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{short}", getLong)
	r.Post("/", store)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
}
