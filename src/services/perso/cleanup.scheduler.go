package main

import (
	"context"
	"log"
	"time"

	r "github.com/Menschomat/bly.li/shared/redis"
)

// cleanupStream is a scheduled job that deletes messages that have been fully acknowledged.
// It determines a "safe" cutoff using XPENDING and then deletes messages with IDs less than that cutoff.
func cleanupStream() {
	// Create a short-lived context for cleanup.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new Redis client.
	client := r.GetRedisClient()
	streamKey := "blowup_action_click"
	groupName := "blowup_action_click_group"

	// Get pending message info for the consumer group.
	pending, err := client.XPending(ctx, streamKey, groupName).Result()
	if err != nil {
		log.Printf("Error getting XPENDING info: %v", err)
		return
	}

	var safeID string
	if pending.Count > 0 {
		// There are pending messages; use the smallest pending ID as the cutoff.
		safeID = pending.Lower
	} else {
		// No pending messages. Use the first message in the stream.
		entries, err := client.XRange(ctx, streamKey, "-", "+").Result()
		if err != nil {
			log.Printf("Error getting XRANGE for stream: %v", err)
			return
		}
		if len(entries) == 0 {
			// Stream is empty.
			return
		}
		safeID = entries[0].ID
	}

	// Delete messages with IDs less than safeID in batches.
	var totalDeleted int64
	for {
		// Fetch a batch of up to 100 messages older than safeID.
		entries, err := client.XRangeN(ctx, streamKey, "-", safeID, 100).Result()
		if err != nil {
			log.Printf("Error reading XRANGE: %v", err)
			break
		}
		if len(entries) == 0 {
			break
		}

		ids := make([]string, 0, len(entries))
		for _, entry := range entries {
			ids = append(ids, entry.ID)
		}

		// Delete the batch of messages.
		n, err := client.XDel(ctx, streamKey, ids...).Result()
		if err != nil {
			log.Printf("Error deleting messages: %v", err)
			break
		}
		totalDeleted += n

		// If fewer than 100 messages were returned, we assume we're done.
		if len(entries) < 100 {
			break
		}
	}
	log.Printf("Cleanup completed: deleted %d messages older than safeID %s", totalDeleted, safeID)
}
