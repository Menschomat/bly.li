package logging

import (
	"log/slog"
	"os"
	"sync"

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
		config, _ := loki.NewDefaultConfig("http://loki:3100/loki/api/v1/push")
		config.TenantID = "xyz"
		client, _ := loki.New(config)
		// Stdout handler
		stdoutHandler := slog.NewJSONHandler(os.Stdout, nil)
		// Loki handler
		lokiHandler := slogloki.Option{Level: slog.LevelDebug, Client: client}.NewLokiHandler()

		handler = slogmulti.Fanout(stdoutHandler, lokiHandler)
	})
	return handler
}

// NewLogger creates a logger with a service-name.
func NewLogger(serviceName string) *slog.Logger {
	return slog.New(getHandler()).With("service_name", serviceName)
}
