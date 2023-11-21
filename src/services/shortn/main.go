package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	u "github.com/Menschomat/bly.li/services/shortn/utils"
	m "github.com/Menschomat/bly.li/shared/model"
	redis "github.com/Menschomat/bly.li/shared/redis"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortn m.ShortnReq
	json.NewDecoder(r.Body).Decode(&shortn)
	url, err := u.ParseUrl(shortn.Url)
	if err != nil {
		apiUtils.BadRequestError(w, r)
		return
	}
	short := u.GetUniqueShort()
	redis.StoreUrl(short, url)
	payload, err := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if err != nil {
		apiUtils.InternalServerError(w, r)
		return
	}
	w.Write(payload)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-Auth-User")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello")
	log.Println(user)
	log.Println(r.Header)
	log.Println("HURZ")
}

func main() {
	log.Println("*_-_-_-BlyLi-Shortn-_-_-_*")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", store)
	r.Get("/", getAll)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
}
