package redis

import "log"

// MarkUnsaved adds an entry to the Redis set "unsaved"
func MarkToDel(short string) {
	_cache := getRedisClient()
	err := _cache.SAdd(ctx, "unsaved", short).Err()
	if err != nil {
		log.Println("Error adding to unsaved set:", err)
	}
}

// GetUnsavedKeys retrieves all items from the Redis set "unsaved"
func GetToDelKeys() ([]string, error) {
	_cache := getRedisClient()
	keys, err := _cache.SMembers(ctx, "unsaved").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// RemoveUnsaved removes an entry from the "unsaved" set
func RemoveToDel(short string) error {
	_cache := getRedisClient()
	err := _cache.SRem(ctx, "unsaved", short).Err()
	if err != nil {
		log.Println("Error removing from unsaved set:", err)
		return err
	}
	return nil
}
