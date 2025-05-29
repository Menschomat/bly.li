package tracking

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/Menschomat/bly.li/services/perso/logging"
	"github.com/Menschomat/bly.li/shared/data"
	"github.com/Menschomat/bly.li/shared/model"
	"github.com/Menschomat/bly.li/shared/mongo"
	r "github.com/Menschomat/bly.li/shared/redis"
	l "github.com/Menschomat/bly.li/shared/utils/logging"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	logger = logging.GetLogger()
)

// ClickAggregator groups click events by the Short field.
type ClickAggregator struct {
	mu     sync.Mutex
	groups map[string][]model.ShortClick
}

// NewClickAggregator creates a new aggregator instance.
func NewClickAggregator() *ClickAggregator {
	return &ClickAggregator{
		groups: make(map[string][]model.ShortClick),
	}
}

// Add groups a click event by its Short value.
func (agg *ClickAggregator) Add(click model.ShortClick) {
	agg.mu.Lock()
	defer agg.mu.Unlock()

	// Append click to the slice for the given short.
	agg.groups[click.Short] = append(agg.groups[click.Short], click)
}

// Flush retrieves and resets the current groups.
func (agg *ClickAggregator) Flush() map[string][]model.ShortClick {
	agg.mu.Lock()
	defer agg.mu.Unlock()

	// Copy current groups and reset.
	flushed := agg.groups
	agg.groups = make(map[string][]model.ShortClick)
	return flushed
}

// persistAggregatedClicks processes the aggregated clicks.
// For each short, it calls your persist function (e.g., saving to MongoDB).
func PersistAggregatedClicks(aggregated map[string][]model.ShortClick) {
	summedClicks := []model.ShortClick{}
	for short, clicks := range aggregated {
		increment := len(clicks)
		logger.Info("Persisting clicks for short",
			"clicks", increment,
			"short", short,
		)
		mongo.InsetTimeseriesDoc(short, increment, time.Now())
		s := data.GetShort(short)
		if s != nil {
			s.Count += increment
			err := r.StoreUrl(s.Short, s.URL, s.Count, s.Owner)
			l.LogRedisError(err)
			r.MarkUnsaved(s.Short)
		}
		summedClicks = append(summedClicks, clicks...)
	}
	mongo.InsetTimeseriesData("clicks", summedClicks)
}

// runConsumer starts a Redis stream consumer that aggregates click events.
func RunConsumer(ctx context.Context, aggregator *ClickAggregator) {
	client := r.GetRedisClient()

	streamKey := "blowup_action_click"
	groupName := "blowup_action_click_group"
	consumerName := uuid.NewString() // Unique consumer name

	// Create the consumer group if it doesn't exist.
	err := client.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		logger.Error("Could not create consumer group", "err", err)
		os.Exit(1)
	}
	logger.Info("Consumer started. Waiting for messages...")
	for {
		select {
		case <-ctx.Done():
			logger.Info("Consumer shutting down...")
			return
		default:
			// XREADGROUP blocks for 5 seconds if no messages are available.
			streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{streamKey, ">"},
				Count:    1,
				Block:    5 * time.Second,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue // No messages available.
				}
				logger.Error("Error reading from stream", "err", err)
				continue
			}
			// Process messages.
			for _, stream := range streams {
				for _, message := range stream.Messages {
					dataStr, ok := message.Values["data"].(string)
					if !ok {
						logger.Warn("Message missing 'data' field or it is not a string", "message_id", message.ID)
						continue
					}
					var clickEvent model.ShortClick
					if err := json.Unmarshal([]byte(dataStr), &clickEvent); err != nil {
						logger.Error("Error unmarshaling message", "message_id", message.ID, "err", err)
						continue
					}
					// Add the click event to the aggregator.
					aggregator.Add(clickEvent)
					// Acknowledge the message.
					if _, err := client.XAck(ctx, streamKey, groupName, message.ID).Result(); err != nil {
						logger.Error("Error acknowledging message", "message_id", message.ID, "err", err)
					}
				}
			}
		}
	}
}
