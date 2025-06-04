package data

import (
	"github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/utils"
	"github.com/Menschomat/bly.li/shared/utils/logging"
)

func GetShort(short string) *model.ShortURL {
	if utils.IsValidShort(short) {
		url, err := redis.GetShort(short)
		logging.LogRedisError(err)
		if err != nil || url == nil {
			url, err = mongo.GetShortURLByShort(short)
			if err == nil {
				redis.StoreUrl(*url)
				return url
			}
			logging.LogMongoError(err)
		}
		return url
	}
	return nil
}
