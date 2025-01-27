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

var _ api.ServerInterface = (*Server)(nil)

type Server struct{}

func (p *Server) DeleteShort(w http.ResponseWriter, r *http.Request, short string) {
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

func (p *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-Auth-User")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Hello")
	if err != nil {
		return
	}
	//TODO remove this logging
	log.Println(user)
	log.Println(r.Header)
}

func (p *Server) PostStore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var shortn m.ShortnReq
	err := json.NewDecoder(r.Body).Decode(&shortn)
	url, err := u.ParseUrl(shortn.Url)
	if err != nil {
		apiUtils.BadRequestError(w)
		return
	}
	short := u.GetUniqueShort()
	redis.StoreUrl(short, url)
	payload, err := json.Marshal(m.ShortnRes{Url: url, Short: short})
	if err != nil {
		apiUtils.InternalServerError(w)
		return
	}
	usrInfo, err := GetUsrInfoFromCtx(r.Context())
	if usrInfo != nil {
		_, err = mongo.StoreShortURL(m.ShortURL{URL: url, Short: short, Owner: usrInfo.Email})
	}
	if err != nil {
		_, err = mongo.StoreShortURL(m.ShortURL{URL: url, Short: short})
	}
	_, err = w.Write(payload)
}
func GetUsrInfoFromCtx(ctx context.Context) (*struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}, error) {
	var claims struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}
	if err := ctx.Value("token").(*oidc.IDToken).Claims(&claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func JWTVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseCtx := context.Background()
		var verifier = oidcProvider.Verifier(&oidc.Config{ClientID: appConfig.OidcConfig.OidcClientId})
		token, err := verifier.Verify(baseCtx, apiUtils.TokenFromHeader(r))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func InitOidcProvider(ctx context.Context, url string) *oidc.Provider {
	provider, oidcErr := oidc.NewProvider(ctx, url)
	if oidcErr != nil {
		panic(oidcErr)
	}
	return provider
}

func main() {
	err := cfgUtils.FillEnvStruct(&appConfig)
	// Set up OIDC provider and OAuth2 config
	oidcProvider = InitOidcProvider(context.Background(), appConfig.OidcConfig.OidcUrl)
	log.Println("*_-_-_-BlyLi-Shortn-_-_-_*")
	// Create new Chi-Router
	r := chi.NewRouter()
	// Add Middlewares to Router
	r.Use(middleware.Logger)
	r.Use(JWTVerifier)
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
	err = http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
	defer mongo.CloseClientDB()
}
