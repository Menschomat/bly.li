package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Menschomat/bly.li/shared/model"
	r "github.com/Menschomat/bly.li/shared/redis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
func persistAggregatedClicks(aggregated map[string][]model.ShortClick) {
	for short, clicks := range aggregated {
		log.Printf("Persisting %d clicks for short: %s", len(clicks), short)
		// Call your persistence logic here.
		// Example: mongo.PersistClicksForShort(short, clicks)
	}
}

// runConsumer starts a Redis stream consumer that aggregates click events.
func runConsumer(ctx context.Context, aggregator *ClickAggregator) {
	client := r.GetRedisClient()

	streamKey := "blowup_action_click"
	groupName := "blowup_action_click_group"
	consumerName := uuid.NewString() // Unique consumer name

	// Create the consumer group if it doesn't exist.
	err := client.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Could not create consumer group: %v", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer shutting down...")
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
				log.Printf("Error reading from stream: %v", err)
				continue
			}

			// Process messages.
			for _, stream := range streams {
				for _, message := range stream.Messages {
					dataStr, ok := message.Values["data"].(string)
					if !ok {
						log.Printf("Message %s missing 'data' field or it is not a string", message.ID)
						continue
					}

					var clickEvent model.ShortClick
					if err := json.Unmarshal([]byte(dataStr), &clickEvent); err != nil {
						log.Printf("Error unmarshaling message %s: %v", message.ID, err)
						continue
					}

					// Add the click event to the aggregator.
					aggregator.Add(clickEvent)
					log.Printf("Aggregated click event for short %s", clickEvent.Short)

					// Acknowledge the message.
					if _, err := client.XAck(ctx, streamKey, groupName, message.ID).Result(); err != nil {
						log.Printf("Error acknowledging message %s: %v", message.ID, err)
					}
				}
			}
		}
	}
}
