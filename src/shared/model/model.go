package model

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

//Configuration-Model

type OidcConfig struct {
	OidcClientId string `env:"OIDC_CLIENT_ID, default=12345"`
	OidcUrl      string `env:"OIDC_URL, default=http://127.0.0.1"`
}

type MongoDdConfig struct {
	Database       string `env:"MONGO_DATABASE, default=short_url_db"`
	MongoServerUrl string `env:"MONGN_SERVER_URL, default=mongodb://mongodb:27017"`
}

type ShortnConfig struct {
	ZookeeperUrl         string `env:"ZOOKEEPER_URL, default=http://localhost"`
	ZookeeperCounterPath string `env:"ZOOKEEPER_COUNTER_PATH, default=/counter"`
}
