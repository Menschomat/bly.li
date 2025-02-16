package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Menschomat/bly.li/shared/scheduler"
)

func main() {
	log.Println("*_-_-_-BlyLi-Perso-_-_-_*")

	unsavedScheduler := scheduler.NewScheduler(10*time.Second, persistUnsaved)
	// New scheduler job to cleanup acknowledged stream messages every minute.
	cleanupScheduler := scheduler.NewScheduler(1*time.Minute, cleanupStream)
	// Create the aggregator for click events.
	aggregator := NewClickAggregator()

	// New scheduler task to flush aggregated clicks every 5 minutes.
	aggregatorFlushScheduler := scheduler.NewScheduler(5*time.Minute, func() {
		aggregated := aggregator.Flush()
		if len(aggregated) > 0 {
			persistAggregatedClicks(aggregated)
		} else {
			log.Println("No aggregated clicks to persist this 5-minute window.")
		}
	})
	// Create a cancellable context for graceful shutdown of the consumer.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the Redis stream consumer in its own goroutine.
	go runConsumer(ctx, aggregator)

	// Handle shutdown signals.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a shutdown signal.
	<-stopChan
	log.Println("Shutdown signal received. Stopping scheduler and consumer...")

	// Stop the scheduler.
	unsavedScheduler.Stop()
	cleanupScheduler.Stop()
	aggregatorFlushScheduler.Stop()

	// Cancel the consumer context to shut it down gracefully.
	cancel()

	// Optionally, wait a bit for cleanup.
	time.Sleep(5 * time.Second)
	log.Println("Server shut down successfully.")
}
