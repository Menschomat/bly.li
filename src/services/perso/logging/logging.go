package logging

import (
	"log/slog"
	"sync"

	l "github.com/Menschomat/bly.li/shared/utils/logging"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func GetLogger() *slog.Logger {
	once.Do(func() {
		logger = l.NewLogger("perso")
	})
	return logger
}
