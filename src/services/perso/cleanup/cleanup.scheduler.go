package cleanup

import (
	"context"
	"time"

	"github.com/Menschomat/bly.li/services/perso/logging"
	r "github.com/Menschomat/bly.li/shared/redis"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cleanedStreamEvents = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "stream_events_cleaned_total",
		Help: "Total number of clicks handled by blowup",
	})
	logger = logging.GetLogger()
)

func InitMetrics() {
	prometheus.MustRegister(cleanedStreamEvents)
}

// cleanupStream is a scheduled job that deletes messages that have been fully acknowledged.
// It determines a "safe" cutoff using XPENDING and then deletes messages with IDs less than that cutoff.
func CleanupStream() {
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
		logger.Error("Error getting XPENDING info", "err", err)
		return
	}

	var trimUpToID string
	if pending.Count > 0 {
		trimUpToID = pending.Lower
	} else {
		// Get the ID of the newest message in the stream
		entries, err := client.XRevRangeN(ctx, streamKey, "+", "-", 1).Result()
		if err != nil {
			logger.Error("Error getting latest stream ID", "err", err)
			return
		}
		if len(entries) == 0 {
			// Stream is already empty
			logger.Info("Stream already empty, nothing to clean.")
			return
		}
		trimUpToID = entries[0].ID
	}

	// Delete messages with IDs less than safeID in batches.
	var totalDeleted int64
	for {
		// Fetch a batch of up to 100 messages older than safeID.
		entries, err := client.XRangeN(ctx, streamKey, "-", trimUpToID, 100).Result()
		if err != nil {
			logger.Error("Error reading XRANGE", "err", err)
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
			logger.Error("Error deleting messages", "err", err)
			break
		}
		totalDeleted += n

		// If fewer than 100 messages were returned, we assume we're done.
		if len(entries) < 100 {
			break
		}
	}
	cleanedStreamEvents.Add(float64(totalDeleted))
	logger.Info("Cleanup completed",
		"deleted", totalDeleted,
		"safeID", trimUpToID,
	)
}
