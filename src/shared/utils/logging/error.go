package logging

import (
	"log/slog"
)

func LogMongoError(err error) {
	if err != nil {
		slog.Error("MongoDB error occurred", "error", err)
	}
}
func LogRedisError(err error) {
	if err != nil {
		slog.Error("Redis error occurred", "error", err)
	}
}
