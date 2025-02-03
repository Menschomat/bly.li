package redis

import (
	"context"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/Menschomat/bly.li/shared/model"
	"github.com/redis/go-redis/v9"
)

var cacheClient *redis.Client
var targetTtl = 1 * time.Minute
var ctx = context.Background()

func getRedisClient() *redis.Client {
	if cacheClient != nil {
		return cacheClient
	}
	cacheClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	return cacheClient
}

func StoreUrl(short string, url string, count int) {
	_cache := getRedisClient()
	key := "url:" + short
	errs := [...]error{
		_cache.HSet(ctx, key, "url", url, "count", count).Err(),
		_cache.Expire(ctx, key, targetTtl).Err(),
	}
	for _, err := range errs {
		if err != nil {
			log.Println("Error while storing url", err)
		}
	}
}

func GetShort(short string) (u model.ShortURL, e error) {
	_cache := getRedisClient()
	key := "url:" + short

	// Fetch all fields of the ShortURL struct
	data, err := _cache.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		slog.Warn("Could not fetch url from redis!")
		return u, err // return empty struct and error
	}

	// Set the expiration time (refresh TTL)
	if err := _cache.Expire(ctx, key, targetTtl).Err(); err != nil {
		slog.Warn("Could not refresh TTL in redis!", "error", err)
	}

	// Convert count from string to int
	count, _ := strconv.Atoi(data["count"]) // Ignore error, assume default 0

	// Map Redis data to the model.ShortURL struct
	u = model.ShortURL{
		Short: short,       // assuming "short" is a field
		URL:   data["url"], // assuming "url" is a field
		Count: count,       // assuming "count" is a field
	}

	return u, nil
}

func GetUrl(short string) (u string, e error) {
	_cache := getRedisClient()
	key := "url:" + short
	url, err := _cache.HGet(ctx, key, "url").Result()
	if err != nil || _cache.Expire(ctx, key, targetTtl).Err() != nil {
		slog.Warn("Could not fetch url from redis!")
		return "", err
	}
	return url, nil
}

func DeleteUrl(short string) (e error) {
	_cache := getRedisClient()
	key := "url:" + short
	_, err := _cache.HDel(ctx, key, "url").Result()
	if err != nil || _cache.Expire(ctx, key, targetTtl).Err() != nil {
		slog.Warn("Could not delete url from redis!")
		return err
	}
	MarkToDel(short)
	RemoveUnsaved(short)
	return nil
}

func ShortExists(short string) bool {
	_cache := getRedisClient()
	exists := _cache.Exists(ctx, "url:"+short)
	return exists.Val() > 0
}

func RegisterClick(short string) {
	_cache := getRedisClient()
	count, err1 := strconv.Atoi(_cache.HGet(ctx, "url:"+short, "count").Val())
	if err1 != nil {
		log.Println(_cache.HGet(ctx, "url:"+short, "url").Val(), err1)
	}
	err2 := _cache.HSet(ctx, "url:"+short, "count", count+1).Err()
	if err2 != nil {
		log.Println(err2)
	}
	MarkUnsaved(short)
}
