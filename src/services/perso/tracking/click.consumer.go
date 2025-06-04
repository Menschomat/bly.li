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
	agg.groups[click.Short] = append(agg.groups[click.Short], click)
}

// Flush retrieves and resets the current groups atomically.
func (agg *ClickAggregator) Flush() map[string][]model.ShortClick {
	agg.mu.Lock()
	defer agg.mu.Unlock()
	flushed := agg.groups
	agg.groups = make(map[string][]model.ShortClick)
	return flushed
}

// PersistAggregatedClicks processes the aggregated clicks per short.
// For each short, it persists the click count and updates necessary stores.
func PersistAggregatedClicks(aggregated map[string][]model.ShortClick) {
	var allClicks []model.ShortClick
	for short, clicks := range aggregated {
		count := len(clicks)
		logger.Info("Persisting clicks for short", "clicks", count, "short", short)
		// Update the short's total click count in both redis and mark as unsaved
		s := data.GetShort(short)
		if s != nil {
			s.Count += count
			if err := r.StoreUrl(*s); err != nil {
				l.LogRedisError(err)
			}
			r.MarkUnsaved(s.Short)
		}
		allClicks = append(allClicks, clicks...)
	}
	mongo.InsetTimeseriesData(mongo.CollectionClicks, allClicks)
}

// RunConsumer continuously reads click events from a Redis stream and aggregates them.
func RunConsumer(ctx context.Context, aggregator *ClickAggregator) {
	client := r.GetRedisClient()
	const (
		streamKey = "blowup_action_click"
		groupName = "blowup_action_click_group"
		blockTime = 5 * time.Second
		batchSize = 1
	)
	consumerName := uuid.NewString() // Unique consumer name for this instance

	if err := createConsumerGroup(ctx, client, streamKey, groupName); err != nil {
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
			processStreamMessages(ctx, client, aggregator, streamKey, groupName, consumerName, batchSize, blockTime)
		}
	}
}

// createConsumerGroup ensures the Redis consumer group exists.
func createConsumerGroup(ctx context.Context, client *redis.Client, streamKey, groupName string) error {
	err := client.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

// processStreamMessages reads and acknowledges messages from the Redis stream.
func processStreamMessages(
	ctx context.Context,
	client *redis.Client,
	aggregator *ClickAggregator,
	streamKey, groupName, consumerName string,
	batchSize int,
	block time.Duration,
) {
	streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamKey, ">"},
		Count:    int64(batchSize),
		Block:    block,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return // No messages available, continue loop
		}
		logger.Error("Error reading from stream", "err", err)
		return
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			clickEvent, ok := parseClickEvent(message)
			if !ok {
				continue
			}
			aggregator.Add(clickEvent)
			acknowledgeMessage(ctx, client, streamKey, groupName, message.ID)
		}
	}
}

// parseClickEvent safely extracts and unmarshals a ShortClick event from a Redis message.
func parseClickEvent(message redis.XMessage) (model.ShortClick, bool) {
	dataStr, ok := message.Values["data"].(string)
	if !ok {
		logger.Warn("Message missing 'data' field or it is not a string", "message_id", message.ID)
		return model.ShortClick{}, false
	}
	var clickEvent model.ShortClick
	if err := json.Unmarshal([]byte(dataStr), &clickEvent); err != nil {
		logger.Error("Error unmarshaling message", "message_id", message.ID, "err", err)
		return model.ShortClick{}, false
	}
	return clickEvent, true
}

// acknowledgeMessage acknowledges the consumption of a message in the Redis stream.
func acknowledgeMessage(ctx context.Context, client *redis.Client, streamKey, groupName, messageID string) {
	if _, err := client.XAck(ctx, streamKey, groupName, messageID).Result(); err != nil {
		logger.Error("Error acknowledging message", "message_id", messageID, "err", err)
	}
}
