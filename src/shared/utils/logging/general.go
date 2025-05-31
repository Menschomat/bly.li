package logging

import (
	"log/slog"
	"os"
	"sync"

	"github.com/Menschomat/bly.li/shared/config"
	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
	slogmulti "github.com/samber/slog-multi"
)

var (
	handler slog.Handler
	once    sync.Once
)

func getHandler() slog.Handler {
	once.Do(func() {
		cfg := config.LoggingConfig()

		config, _ := loki.NewDefaultConfig(cfg.LokiUrl)
		config.TenantID = cfg.LokiTenant
		client, _ := loki.New(config)

		// Stdout handler
		stdoutHandler := slog.NewJSONHandler(os.Stdout, nil)

		// Parse log level
		var logLevel slog.Level
		switch cfg.LogLevel {
		case "debug":
			logLevel = slog.LevelDebug
		case "warn":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		default:
			logLevel = slog.LevelInfo
		}

		// Loki handler
		lokiHandler := slogloki.Option{
			Level:  logLevel,
			Client: client,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove all user fields (let only Labels map set the Loki labels)
				// This results in user fields being kept in the entry's body, not as Loki labels.
				if a.Key != "service_name" && a.Key != "instance" && a.Key != "level" {
					// Return an Attr that will not be interpreted as a label
					a.Key = "" // or skip it in a more advanced way
				}
				return a
			},
		}.NewLokiHandler()

		handler = slogmulti.Fanout(stdoutHandler, lokiHandler)
	})
	return handler
}

// NewLogger creates a logger with a service-name.
func NewLogger(serviceName string) *slog.Logger {
	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		instanceID, _ = os.Hostname()
	}
	return slog.New(getHandler()).With("service_name", serviceName).With("instance", instanceID)
}
