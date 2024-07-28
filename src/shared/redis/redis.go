package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var cacheClient *redis.Client
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

func StoreUrl(short string, url string) {
	_cache := getRedisClient()
	key := "url:" + short
	errs := [...]error{
		_cache.HSet(ctx, key, "url", url, "count", 0).Err(),
		_cache.Expire(ctx, key, 60*60*60*24).Err(),
	}
	for _, err := range errs {
		if err != nil {
			log.Println("Error while storing url", err)
		}
	}
}
func GetUrl(short string) (u string, e error) {
	_cache := getRedisClient()
	url, err := _cache.HGet(ctx, "url:"+short, "url").Result()
	if err != nil {
		log.Println("Warning: Could not fetch url from redis!")
		return "", err
	}
	return url, nil
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
}
