package model

type OidcConfig struct {
	OidcClientId string `env:"CLIENT_ID, default=12345"`
	OidcUrl      string `env:"URL, default=http://127.0.0.1"`
}

type ShortURL struct {
	Short string `bson:"short"`
	URL   string `bson:"url"`
	Owner string `bson:"owner,omitempty"`
}

type ShortnReq struct {
	Url string
	//Short string
}

type ShortnRes struct {
	Url   string
	Short string
}
type ShortnConfig struct {
	// OidcConfig will process values from $OIDC_* and
	// $OIDC_CLIENT_ID respectively.
	OidcConfig *OidcConfig `env:", prefix=OIDC_"`
}
