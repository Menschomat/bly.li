package model

import "time"

type ShortURL struct {
	Short string `bson:"short"`
	URL   string `bson:"url"`
	Count int    `bson:"count"`
	Owner string `bson:"owner,omitempty"`
}

type ShortClick struct {
	Short     string
	Timestamp time.Time
	Ip        string
	UsrAgent  string
}

type ShortClickCount struct {
	Short     string    `bson:"short"`
	Timestamp time.Time `bson:"timestamp"`
	Count     int       `bson:"count"`
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
	ZookeeperHost        string `env:"ZOOKEEPER_HOST, default=zookeeper:2181"`
	ZookeeperCounterPath string `env:"ZOOKEEPER_COUNTER_PATH, default=/shortn-ranges"`
	ServerPort           string `env:"SERVER_PORT, default=:8082"`
	MetricsPort          string `env:"METRICS_PORT, default=:9082"`
	CorsAllowedOrigins   string `env:"CORS_ALLOWED_ORIGINS, default=https://*,http://*"`
	CorsMaxAge           int    `env:"CORS_MAX_AGE, default=300"`
}

type BlowupConfig struct {
	ServerPort   string `env:"SERVER_PORT, default=:8081"`
	MetricsPort  string `env:"METRICS_PORT, default=:9081"`
	RedirectCode int    `env:"REDIRECT_CODE, default=302"`
}

type DasherConfig struct {
	ServerPort  string `env:"SERVER_PORT, default=:8083"`
	MetricsPort string `env:"METRICS_PORT, default=:9083"`
}

type PersoConfig struct {
	ServerPort      string `env:"SERVER_PORT, default=:8084"`
	MetricsPort     string `env:"METRICS_PORT, default=:9084"`
	CleanupInterval string `env:"CLEANUP_INTERVAL, default=24h"`
}

type LoggingConfig struct {
	LokiUrl    string `env:"LOKI_URL, default=http://loki:3100/loki/api/v1/push"`
	LokiTenant string `env:"LOKI_TENANT, default=single"`
	LogLevel   string `env:"LOG_LEVEL, default=info"`
}
