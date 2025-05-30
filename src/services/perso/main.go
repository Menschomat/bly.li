package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Menschomat/bly.li/services/perso/cleanup"
	"github.com/Menschomat/bly.li/services/perso/logging"
	"github.com/Menschomat/bly.li/services/perso/persistence"
	"github.com/Menschomat/bly.li/services/perso/tracking"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/scheduler"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	logger = logging.GetLogger()
)

func main() {
	logger.Info("Starting")
	mongo.InitMongoPackage(logger)
	//Schedulers------------------------------
	unsavedScheduler := scheduler.NewScheduler("persistance", 10*time.Second, persistence.PersistUnsaved, logger)
	// New scheduler job to cleanup acknowledged stream messages every minute.
	cleanupScheduler := scheduler.NewScheduler("clean-up", 10*time.Second, cleanup.CleanupStream, logger)
	// Create the aggregator for click events.
	aggregator := tracking.NewClickAggregator()

	// New scheduler task to flush aggregated clicks every 5 minutes.
	clickFlushScheduler := scheduler.NewScheduler("click-handling", 30*time.Second, func() {
		aggregated := aggregator.Flush()
		if len(aggregated) > 0 {
			tracking.PersistAggregatedClicks(aggregated)
		} else {
			logger.Debug("No aggregated clicks to persist.")
		}
	}, logger)
	aggregatorScheduler := scheduler.NewScheduler("aggregation", 1*time.Minute, func() {
		tracking.AggregateClicks()
	}, logger)
	// Create a cancellable context for graceful shutdown of the consumer.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the Redis stream consumer in its own goroutine.
	go tracking.RunConsumer(ctx, aggregator)

	serverErrChan := make(chan error, 1)

	// --- Metrics router on a 2nd port ---
	go func() {
		//Metrics---------------------------------
		cleanup.InitMetrics()
		metricsRouter := chi.NewRouter()
		metricsRouter.Handle("/metrics", promhttp.Handler())
		logger.Info("Prometheus metrics available on :2115/metrics")
		serverErrChan <- http.ListenAndServe(":2115", metricsRouter)
	}()

	// Handle shutdown signals.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverErrChan:
		logger.Error("Server error", "error", err)
	case <-stopChan:
		logger.Info("Shutdown signal received. Stopping scheduler and consumer...")
	}

	// Stop the scheduler.
	unsavedScheduler.Stop()
	cleanupScheduler.Stop()
	clickFlushScheduler.Stop()
	aggregatorScheduler.Stop()

	// Cancel the consumer context to shut it down gracefully.
	cancel()

	// Optionally, wait a bit for cleanup.
	time.Sleep(5 * time.Second)
	logger.Info("Server shut down successfully.")
}
