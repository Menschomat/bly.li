package cleanup

import (
	"context"
	"time"

	"github.com/Menschomat/bly.li/services/perso/logging"
	r "github.com/Menschomat/bly.li/shared/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

var (
	cleanedStreamEvents = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "blyli_stream_events_cleaned_total",
		Help: "Total number of cleaned stream events",
	})
	logger = logging.GetLogger()
)

// InitMetrics registers Prometheus metrics for cleanup events.
func InitMetrics() {
	prometheus.MustRegister(cleanedStreamEvents)
}

// CleanupStream deletes fully acknowledged messages from the Redis stream.
// It works by determining a "safe" cutoff ID using XPENDING, then deleting messages older than this.
func CleanupStream() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := r.GetRedisClient()
	const (
		streamKey = "blowup_action_click"
		groupName = "blowup_action_click_group"
		batchSize = 100
	)

	trimUpToID, ok := findCleanupBoundary(ctx, client, streamKey, groupName)
	if !ok {
		return // Already logged; nothing to delete or error occurred.
	}

	deletedCount := deleteMessagesBeforeID(ctx, client, streamKey, trimUpToID, batchSize)
	cleanedStreamEvents.Add(float64(deletedCount))
	logger.Info("Cleanup completed", "deleted", deletedCount, "safeID", trimUpToID)
}

// findCleanupBoundary computes the "safe" ID up to which messages can be deleted.
func findCleanupBoundary(ctx context.Context, client *redis.Client, streamKey, groupName string) (string, bool) {
	pending, err := client.XPending(ctx, streamKey, groupName).Result()
	if err != nil {
		logger.Error("Error getting XPENDING info", "err", err)
		return "", false
	}
	if pending.Count > 0 {
		return pending.Lower, true
	}

	entries, err := client.XRevRangeN(ctx, streamKey, "+", "-", 1).Result()
	if err != nil {
		logger.Error("Error getting latest stream ID", "err", err)
		return "", false
	}
	if len(entries) == 0 {
		logger.Info("Stream already empty, nothing to clean.")
		return "", false
	}
	return entries[0].ID, true
}

// deleteMessagesBeforeID repeatedly deletes batches of messages before the provided ID.
func deleteMessagesBeforeID(ctx context.Context, client *redis.Client, streamKey, upToID string, batchSize int) int64 {
	var totalDeleted int64

	for {
		entries, err := client.XRangeN(ctx, streamKey, "-", upToID, int64(batchSize)).Result()
		if err != nil {
			logger.Error("Error reading XRANGE", "err", err)
			break
		}
		if len(entries) == 0 {
			break
		}

		ids := make([]string, len(entries))
		for i, entry := range entries {
			ids[i] = entry.ID
		}

		n, err := client.XDel(ctx, streamKey, ids...).Result()
		if err != nil {
			logger.Error("Error deleting messages", "err", err)
			break
		}
		totalDeleted += n

		if len(entries) < batchSize {
			break
		}
	}

	return totalDeleted
}
