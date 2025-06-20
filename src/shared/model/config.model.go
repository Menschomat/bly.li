package model

//Configuration-Model

type OidcConfig struct {
	OidcClientId string `env:"OIDC_CLIENT_ID, default=12345"`
	OidcUrl      string `env:"OIDC_URL, default=http://127.0.0.1"`
}

type MongoDdConfig struct {
	Database       string `env:"MONGO_DATABASE, default=short_url_db"`
	MongoServerUrl string `env:"MONGO_SERVER_URL, default=mongodb://mongodb:27017"`
}

// BaseServerConfig contains common server configuration fields
type BaseServerConfig struct {
	ServerPort         string `env:"SERVER_PORT"`
	MetricsPort        string `env:"METRICS_PORT"`
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS, default=https://*,http://*"`
	CorsMaxAge         int    `env:"CORS_MAX_AGE, default=300"`
}

type ShortnConfig struct {
	BaseServerConfig
	ZookeeperHost        string `env:"ZOOKEEPER_HOST, default=zookeeper:2181"`
	ZookeeperCounterPath string `env:"ZOOKEEPER_COUNTER_PATH, default=/shortn-ranges"`
}

func (c *ShortnConfig) SetDefaults() {
	c.ServerPort = map[string]string{"": "8082"}[c.ServerPort]
	c.MetricsPort = map[string]string{"": "9082"}[c.MetricsPort]
}

type BlowupConfig struct {
	BaseServerConfig
	RedirectCode int `env:"REDIRECT_CODE, default=302"`
}

func (c *BlowupConfig) SetDefaults() {
	c.ServerPort = map[string]string{"": "8081"}[c.ServerPort]
	c.MetricsPort = map[string]string{"": "9081"}[c.MetricsPort]
}

type DasherConfig struct {
	BaseServerConfig
}

func (c *DasherConfig) SetDefaults() {
	c.ServerPort = map[string]string{"": "8083"}[c.ServerPort]
	c.MetricsPort = map[string]string{"": "9083"}[c.MetricsPort]
}

type PersoConfig struct {
	BaseServerConfig
	CleanupInterval string `env:"CLEANUP_INTERVAL, default=1m"`
}

func (c *PersoConfig) SetDefaults() {
	c.ServerPort = map[string]string{"": "8084"}[c.ServerPort]
	c.MetricsPort = map[string]string{"": "9084"}[c.MetricsPort]
}

type LoggingConfig struct {
	LokiUrl    string `env:"LOKI_URL, default=http://loki:3100/loki/api/v1/push"`
	LokiTenant string `env:"LOKI_TENANT, default=single"`
	LogLevel   string `env:"LOG_LEVEL, default=info"`
}
