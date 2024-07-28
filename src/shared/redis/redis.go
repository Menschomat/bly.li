package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var cacheClient = getRedisClient()
var ctx = context.Background()

func getRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func StoreUrl(short string, url string) {
	errs := [...]error{
		cacheClient.HSet(ctx, "url:"+short, "url", url, "count", 0).Err(),
	}
	for _, err := range errs {
		if err != nil {
			log.Println("Error while storing url", err)
		}
	}
}
func GetUrl(short string) (u string, e error) {
	url, err := cacheClient.HGet(ctx, "url:"+short, "url").Result()
	if err != nil {
		log.Println("Warning: Could not fetch url from redis!")
		return "", err
	}
	return url, nil
}
func ShortExists(short string) bool {
	exists := cacheClient.Exists(ctx, "url:"+short)
	if exists.Val() > 0 {
		return true
	}
	return false
}

func RegisterClick(short string) {
	count, err1 := strconv.Atoi(cacheClient.HGet(ctx, "url:"+short, "count").Val())
	if err1 != nil {
		log.Println(cacheClient.HGet(ctx, "url:"+short, "url").Val(), err1)
	}
	err2 := cacheClient.HSet(ctx, "url:"+short, "count", count+1).Err()
	if err2 != nil {
		log.Println(err2)
	}
}
