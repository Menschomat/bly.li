package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Menschomat/bly.li/services/dasher/api"
	"github.com/Menschomat/bly.li/services/dasher/logging"
	"github.com/Menschomat/bly.li/shared/config"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/oidc"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/server"
)

var (
	logger                     = logging.GetLogger()
	cfg                        = config.DasherConfig()
	_      api.ServerInterface = (*DasherServer)(nil)
)

// DasherServer implements the API interface.
type DasherServer struct{}

// getShortURLOwner attempts to find the owner of the short url in Redis and MongoDB.
func getShortURLOwner(ctx context.Context, short string) (string, error) {
	if u, err := redis.GetShort(short); err == nil && u != nil && u.Owner != "" {
		return u.Owner, nil
	}
	if u, err := mongo.GetShortURLByShort(short); err == nil && u != nil && u.Owner != "" {
		return u.Owner, nil
	}
	return "", errors.New("owner not found")
}

// DeleteShortShort deletes a short URL if the current user is the owner.
func (s *DasherServer) DeleteShortShort(w http.ResponseWriter, r *http.Request, short string) {
	usrInfo, err := oidc.GetUsrInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("Failed to get user info from context", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	owner, err := getShortURLOwner(r.Context(), short)
	if err != nil || owner != usrInfo.Email {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	if err := deleteShortFromStores(short); err != nil {
		logger.Error("Failed to delete short URL", "short", short, "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteShortFromStores attempts to remove the short URL from both Redis and MongoDB.
func deleteShortFromStores(short string) error {
	// Attempt deletion from Redis (ignore not found)
	if err := redis.DeleteUrl(short); err != nil && redis.ShortExists(short) {
		return errors.New("failed to delete from Redis")
	}

	// Attempt deletion from MongoDB (ignore not found)
	if err := mongo.DeleteShortURLByShort(short); err != nil && mongo.ShortExists(short) {
		return errors.New("failed to delete from MongoDB")
	}

	return nil
}

// GetShortAll returns all short URLs owned by the authenticated user.
func (s *DasherServer) GetShortAll(w http.ResponseWriter, r *http.Request) {
	usrInfo, err := oidc.GetUsrInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("Failed to get user info from context", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	shorts := mongo.GetShortsByOwner(usrInfo.Email)
	responseJSON(w, shorts, http.StatusOK)
}

// responseJSON marshals data to JSON and writes it to the response.
func responseJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If JSON encoding fails, log and return 500.
		logger.Error("Failed to encode response JSON", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	logger.Info("Starting dasher service")
	mongo.InitMongoPackage(logger)
	defer mongo.CloseClientDB()

	srv := server.NewServer(server.Config{
		ServerPort:         cfg.ServerPort,
		MetricsPort:        cfg.MetricsPort,
		CorsAllowedOrigins: strings.Split(cfg.CorsAllowedOrigins, ","),
		CorsMaxAge:         cfg.CorsMaxAge,
		Logger:             logger,
	})

	srv.ConfigureCommonMiddleware()
	srv.Router().Use(oidc.JWTVerifier)

	apiHandler := &DasherServer{}
	api.HandlerFromMux(apiHandler, srv.Router())

	serverErrChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.ServeMetrics(serverErrChan)
	go srv.ServeHTTP(serverErrChan)

	srv.HandleShutdown(ctx, serverErrChan)
	logger.Info("Server shut down successfully.")
}
