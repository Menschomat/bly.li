package logging

import (
	"log/slog"
	"os"
	"sync"
)

var (
	handler slog.Handler
	once    sync.Once
)

func getHandler() slog.Handler {
	once.Do(func() {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	})
	return handler
}

// NewLogger creates a logger with a service-name.
func NewLogger(serviceName string) *slog.Logger {
	return slog.New(getHandler()).With("service-name", serviceName)
}
