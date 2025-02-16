package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
				redis.StoreUrl(shortInfo.Short, shortInfo.URL, shortInfo.Count)
				url = shortInfo.URL
			}
		}
		if err == nil && len(url) > 0 {
			ip := readUserIP(r)
			userAgent := r.UserAgent()
			go redis.RegisterClick(model.ShortClick{Short: short, Ip: ip, UsrAgent: userAgent, Timestamp: time.Now()})
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
	apiUtils.BadRequestError(w)
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func main() {
	log.Println("*_-_-_-BlyLi-Blowup-_-_-_*")

	// Initialize router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Graceful shutdown handling
	server := &Server{}
	api.HandlerFromMux(server, r)

	// HTTP server in a goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- http.ListenAndServe(":8081", r)
	}()

	// Handle shutdown signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		log.Fatalln("Server error:", err)
	case <-stopChan:
		log.Println("Shutdown signal received. Stopping server...")
	}

	log.Println("Server shut down successfully.")
}
