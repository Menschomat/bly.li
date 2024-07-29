package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Menschomat/bly.li/services/shortn/api"
	u "github.com/Menschomat/bly.li/services/shortn/utils"
	m "github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	redis "github.com/Menschomat/bly.li/shared/redis"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
	cfgUtils "github.com/Menschomat/bly.li/shared/utils/config"
	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	appConfig    m.ShortnConfig
	oidcProvider *oidc.Provider
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

func (p *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-Auth-User")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello")
	log.Println(user)
	log.Println(r.Header)
}

func (p *Server) PostStore(w http.ResponseWriter, r *http.Request) {
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
	mongo.StoreShortURL(m.ShortURL{URL: url, Short: short})
	w.Write(payload)
}

func JWTVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base_ctx := context.Background()
		var verifier = oidcProvider.Verifier(&oidc.Config{ClientID: appConfig.OidcConfig.OidcClientId})
		token, err := verifier.Verify(base_ctx, apiUtils.TokenFromHeader(r))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func InitOidcProvicer(ctx context.Context, url string) *oidc.Provider {
	provider, oidc_err := oidc.NewProvider(ctx, url)
	if oidc_err != nil {
		panic(oidc_err)
	}
	return provider
}

func main() {
	cfgUtils.FillEnvStruct(&appConfig)
	// Set up OIDC provider and OAuth2 config
	oidcProvider = InitOidcProvicer(context.Background(), appConfig.OidcConfig.OidcUrl)
	log.Println("*_-_-_-BlyLi-Shortn-_-_-_*")
	// Create new Chi-Router
	r := chi.NewRouter()
	// Add Middlewares to Router
	r.Use(middleware.Logger)
	r.Use(JWTVerifier)
	server := &Server{}
	api.HandlerFromMux(server, r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
	defer mongo.CloseClientDB()
}
