package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Menschomat/bly.li/services/shortn/api"
	u "github.com/Menschomat/bly.li/services/shortn/utils"
	m "github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
	cfgUtils "github.com/Menschomat/bly.li/shared/utils/config"
	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	appConfig    m.ShortnConfig
	oidcProvider *oidc.Provider
)

// Ensure Server implements the ServerInterface
var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

// Utility function to delete short URLs from both Redis and MongoDB
func deleteShortURL(short string) error {
	if redis.ShortExists(short) {
		if err := redis.DeleteUrl(short); err != nil {
			return err
		}
	}
	if mongo.ShortExists(short) {
		if err := mongo.DeleteShortURLByShort(short); err != nil {
			return err
		}
	}
	return nil
}

func (p *Server) DeleteShort(w http.ResponseWriter, _ *http.Request, short string) {
	if err := deleteShortURL(short); err != nil {
		apiUtils.InternalServerError(w)
		log.Printf("Error deleting short URL: %v", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (p *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo, err := GetUserInfoFromCtx(r.Context())
	if err != nil {
		apiUtils.InternalServerError(w)
		log.Printf("Error retrieving user info: %v", err)
		return
	}

	if userInfo != nil {
		shorts, err := mongo.GetShortsByUsr(userInfo.Sub)
		if err != nil {
			apiUtils.InternalServerError(w)
			log.Printf("Error fetching user shorts: %v", err)
			return
		}

		payload, _ := json.Marshal(shorts)
		_, err = w.Write(payload)
		if err != nil {
			return
		}
	}
}

func (p *Server) PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortReq m.ShortnReq

	if err := json.NewDecoder(r.Body).Decode(&shortReq); err != nil {
		apiUtils.BadRequestError(w)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	url, err := u.ParseUrl(shortReq.Url)
	if err != nil {
		apiUtils.BadRequestError(w)
		log.Printf("Invalid URL: %v", err)
		return
	}

	short := u.GetUniqueShort()

	redis.StoreUrl(short, url)

	payload, err := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if err != nil {
		apiUtils.InternalServerError(w)
		log.Printf("Error marshalling response payload: %v", err)
		return
	}

	userInfo, err := GetUserInfoFromCtx(r.Context())
	if userInfo != nil {
		_, err = mongo.StoreShortURL(m.ShortURL{URL: url, Short: short, Owner: userInfo.Sub})
	} else {
		_, err = mongo.StoreShortURL(m.ShortURL{URL: url, Short: short})
	}

	if err != nil {
		apiUtils.InternalServerError(w)
		log.Printf("Error storing short URL in MongoDB: %v", err)
		return
	}

	_, err = w.Write(payload)
	if err != nil {
		return
	}
}

// GetUserInfoFromCtx Move utility function to a shared utils package if necessary
func GetUserInfoFromCtx(ctx context.Context) (*struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Sub      string `json:"sub"`
}, error) {
	var claims struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Sub      string `json:"sub"`
	}

	token := ctx.Value("token").(*oidc.IDToken)
	if err := token.Claims(&claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// JWTVerifier Middleware to verify JWT tokens
func JWTVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseCtx := context.Background()
		verifier := oidcProvider.Verifier(&oidc.Config{ClientID: appConfig.OidcConfig.OidcClientId})

		token, err := verifier.Verify(baseCtx, apiUtils.TokenFromHeader(r))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Printf("JWT verification failed: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func InitOidcProvider(ctx context.Context, url string) *oidc.Provider {
	provider, oidcErr := oidc.NewProvider(ctx, url)
	if oidcErr != nil {
		log.Fatalf("Error initializing OIDC provider: %v", oidcErr)
	}
	return provider
}

func main() {
	// Load application config
	if err := cfgUtils.FillEnvStruct(&appConfig); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize OIDC provider
	oidcProvider = InitOidcProvider(context.Background(), appConfig.OidcConfig.OidcUrl)
	log.Println("*_-_-_-BlyLi-Shortn-_-_-_*")

	// Create a new Chi Router
	r := chi.NewRouter()

	// Add Middlewares
	r.Use(middleware.Logger)
	r.Use(JWTVerifier)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	server := &Server{}
	api.HandlerFromMux(server, r)

	// Start server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	// Close MongoDB client on exit
	defer mongo.CloseClientDB()
}
