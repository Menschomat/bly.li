package persistence

import (
	"log"

	"github.com/Menschomat/bly.li/services/perso/logging"
	"github.com/Menschomat/bly.li/shared/data"
	"github.com/Menschomat/bly.li/shared/mongo"
	r "github.com/Menschomat/bly.li/shared/redis"
)

var (
	logger = logging.GetLogger()
)

func PersistUnsaved() {
	keys, err := r.GetUnsavedKeys()
	if err != nil {
		log.Fatalln("There's an error with the server:", err)
	}
	for _, key := range keys {
		short := data.GetShort(key)
		if short == nil {
			logger.Error("Short not found in Redis!", "short", key, "error", err)
			continue
		}
		logger.Info("Storing changed short: " + short.Short)
		mongo.UpdateShortUrl(*short)
		r.RemoveUnsaved(short.Short)
	}
}
