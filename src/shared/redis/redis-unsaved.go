package redis

import (
	"log"
	"log/slog"
)

// MarkUnsaved adds an entry to the Redis set "unsaved"
func MarkUnsaved(short string) {
	slog.Info("Try to mark " + short + " as unsaved!")
	_cache := getRedisClient()
	err := _cache.SAdd(ctx, "unsaved", short).Err()
	if err != nil {
		log.Println("Error adding to unsaved set:", err)
	}
}

// GetUnsavedKeys retrieves all items from the Redis set "unsaved"
func GetUnsavedKeys() ([]string, error) {
	_cache := getRedisClient()
	keys, err := _cache.SMembers(ctx, "unsaved").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// RemoveUnsaved removes an entry from the "unsaved" set
func RemoveUnsaved(short string) error {
	_cache := getRedisClient()
	err := _cache.SRem(ctx, "unsaved", short).Err()
	if err != nil {
		log.Println("Error removing from unsaved set:", err)
		return err
	}
	return nil
}
