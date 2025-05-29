package logging

import (
	"log/slog"
)

func LogMongoError(err error) {
	if err != nil {
		slog.Error("MongoDB-Error occured", "error", err)
	}
}
func LogRedisError(err error) {
	if err != nil {
		slog.Error("Redis-Error occured", "error", err)
	}
}
