package redis

import (
	"context"
	"encoding/json"
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

func StoreUrl(short, url string, count int, owner string) error {
	slog.Info("STORING")
	key := "url:" + short

	pipe := GetRedisClient().Pipeline()
	pipe.HSet(ctx, key, "url", url, "count", count, "owner", owner)
	pipe.Expire(ctx, key, targetTtl)

	_, err := pipe.Exec(ctx)
	return err
}

func GetShort(short string) (*model.ShortURL, error) {
	key := "url:" + short
	data, err := GetRedisClient().HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		slog.Info("Could not fetch url from redis!")
		return nil, err
	}
	if err := GetRedisClient().Expire(ctx, key, targetTtl).Err(); err != nil {
		slog.Error("Could not refresh TTL in redis!", "error", err)
	}

	count, _ := strconv.Atoi(data["count"])
	return &model.ShortURL{
		Short: short,
		URL:   data["url"],
		Count: count,
		Owner: data["owner"],
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
	//short := click.Short
	_cache := GetRedisClient()

	// Concurrently fetch count and URL
	//countResult := _cache.HGet(ctx, "url:"+short, "count")
	//count, err1 := strconv.Atoi(countResult.Val())
	//if err1 == nil && _cache.HSet(ctx, "url:"+short, "count", count+1).Err() != nil {
	//	slog.Error("Error updating count")
	//}
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
	//MarkUnsaved(short)
}
