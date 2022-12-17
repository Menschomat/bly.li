package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
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
		cacheClient.Set(ctx, "url:"+short, url, 0).Err(),
	}
	for _, err := range errs {
		if err != nil {
			log.Println("Error while storing url", err)
		}
	}
}
func GetUrl(short string) string {
	url, err := cacheClient.Get(ctx, "url:"+short).Result()
	if err != nil {
		log.Panicln(err)
	}
	return url
}
func ShortExists(short string) bool {
	exists := cacheClient.Exists(ctx, short)
	if exists.Val() > 0 {
		return true
	}
	return false
}
