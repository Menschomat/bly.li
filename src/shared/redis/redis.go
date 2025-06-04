package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Menschomat/bly.li/shared/model"
	"github.com/redis/go-redis/v9"
)

var (
	cacheClient *redis.Client
	targetTtl   = 1 * time.Minute
	ctx         = context.Background()
)

func GetRedisClient() *redis.Client {
	if cacheClient == nil {
		cacheClient = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})
	}
	return cacheClient
}

func StoreUrl(shortUrl model.ShortURL) error {
	slog.Info("Storing short " + shortUrl.Short)
	key := "url:" + shortUrl.Short
	pipe := GetRedisClient().Pipeline()

	// Set ShortURL properties
	pipe.HSet(ctx, key, "url", shortUrl.URL, "count", shortUrl.Count, "owner", shortUrl.Owner)
	pipe.HSet(ctx, key, "createdAt", shortUrl.CreatedAt.Format(time.RFC3339)) // Use RFC3339 for consistent formatting
	pipe.HSet(ctx, key, "updatedAt", shortUrl.UpdatedAt.Format(time.RFC3339)) // Use RFC3339 for consistent formatting

	// Set TTL
	pipe.Expire(ctx, key, targetTtl)

	_, err := pipe.Exec(ctx)
	return err
}

func GetShort(short string) (*model.ShortURL, error) {
	key := "url:" + short
	data, err := GetRedisClient().HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		slog.Info("Could not fetch url from redis!")
		return nil, err // Return error if fetching from Redis fails or no data found
	}

	if err := GetRedisClient().Expire(ctx, key, targetTtl).Err(); err != nil {
		slog.Error("Could not refresh TTL in redis!", "error", err)
	}

	// Convert count to an integer with error handling
	count, err := strconv.Atoi(data["count"])
	if err != nil {
		return nil, fmt.Errorf("invalid count value for short '%s': %v", short, err)
	}

	// Parse CreatedAt and UpdatedAt from string to time.Time
	createdAt, err := time.Parse(time.RFC3339, data["createdAt"])
	if err != nil {
		return nil, fmt.Errorf("invalid CreatedAt value for short '%s': %v", short, err)
	}

	updatedAt, err := time.Parse(time.RFC3339, data["updatedAt"])
	if err != nil {
		return nil, fmt.Errorf("invalid UpdatedAt value for short '%s': %v", short, err)
	}

	// Create and return ShortURL object
	return &model.ShortURL{
		Short:     short,
		URL:       data["url"],
		Count:     count,
		Owner:     data["owner"],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func GetUrl(short string) (string, error) {
	key := "url:" + short
	url, err := GetRedisClient().HGet(ctx, key, "url").Result()
	if err != nil || GetRedisClient().Expire(ctx, key, targetTtl).Err() != nil {
		slog.Warn("Could not fetch url from redis!")
		return "", err
	}
	return url, nil
}

func DeleteUrl(short string) error {
	key := "url:" + short
	if _, err := GetRedisClient().HDel(ctx, key, "url").Result(); err != nil {
		slog.Warn("Could not delete url from redis!")
		return err
	}
	RemoveUnsaved(short)
	return nil
}

func ShortExists(short string) bool {
	return GetRedisClient().Exists(ctx, "url:"+short).Val() > 0
}

func RegisterClick(click model.ShortClick) {
	_cache := GetRedisClient()
	out, err3 := json.Marshal(click)
	if err3 != nil {
		slog.Error("Error marshalling click", "error", err3)
		return
	}

	if err := _cache.Publish(ctx, "blowup_action_click", out).Err(); err != nil {
		slog.Error("Error publishing click", "error", err)
	}

	if eventID, err4 := _cache.XAdd(ctx, &redis.XAddArgs{
		Stream: "blowup_action_click",
		Values: map[string]interface{}{"data": string(out)},
	}).Result(); err4 != nil {
		slog.Error("Could not add event:", "error", err4)
	} else {
		slog.Debug("Click event added with ID: " + eventID)
	}
}
