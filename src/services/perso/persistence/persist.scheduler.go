package persistence

import (
	"os"

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
		logger.Error("There's an error with the server", "err", err)
		os.Exit(1)
	}
	for _, key := range keys {
		short := data.GetShort(key)
		if short == nil {
			logger.Warn("Unsaved short not found in Redis", "short", key)
			r.RemoveUnsaved(key)
			continue
		}
		logger.Info("Storing changed short: " + short.Short)
		mongo.UpdateShortUrl(*short)
		r.RemoveUnsaved(key)
	}

}
