package oidc

import (
	"context"
	"net/http"
	"sync"
	"time"

	m "github.com/Menschomat/bly.li/shared/model"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
	cfgUtils "github.com/Menschomat/bly.li/shared/utils/config"
	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	appConfig    m.ShortnConfig
	oidcProvider *oidc.Provider
	providerOnce sync.Once
)

func GetOidcProvider() (*oidc.Provider, error) {
	err := cfgUtils.FillEnvStruct(&appConfig)
	if err != nil {
		panic(err)
	}
	providerOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		oidcProvider = initOidcProvider(ctx, appConfig.OidcConfig.OidcUrl)
	})
	return oidcProvider, nil
}

func JWTVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseCtx := context.Background()
		oidcProvider, err := GetOidcProvider()
		if err != nil {
			panic(err)
		}
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

func initOidcProvider(ctx context.Context, url string) *oidc.Provider {
	provider, oidcErr := oidc.NewProvider(ctx, url)
	if oidcErr != nil {
		panic(oidcErr)
	}
	return provider
}
