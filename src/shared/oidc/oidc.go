package oidc

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Menschomat/bly.li/shared/config"
	apiUtils "github.com/Menschomat/bly.li/shared/utils/api"
	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	oidcProvider *oidc.Provider
	providerOnce sync.Once
)

func GetOidcProvider() (*oidc.Provider, error) {
	providerOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		oidcProvider = initOidcProvider(ctx, config.OidcConfig().OidcUrl)
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
		var verifier = oidcProvider.Verifier(&oidc.Config{ClientID: config.OidcConfig().OidcClientId})
		token, err := verifier.Verify(baseCtx, apiUtils.TokenFromHeader(r))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SubjectFromCtx(ctx context.Context) string {
	return ctx.Value("token").(*oidc.IDToken).Subject
}

func initOidcProvider(ctx context.Context, url string) *oidc.Provider {
	provider, oidcErr := oidc.NewProvider(ctx, url)
	if oidcErr != nil {
		panic(oidcErr)
	}
	return provider
}
