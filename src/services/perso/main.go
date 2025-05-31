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
	"github.com/Menschomat/bly.li/shared/config"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/scheduler"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	logger = logging.GetLogger()
	cfg    = config.PersoConfig()
)

// schedulerJob wraps configuration for background scheduler tasks.
type schedulerJob struct {
	name     string
	interval time.Duration
	task     func()
}

func main() {
	logger.Info("Starting perso service")
	mongo.InitMongoPackage(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cleanupInterval, err := time.ParseDuration(cfg.CleanupInterval)
	if err != nil {
		logger.Error("Invalid cleanup interval", "error", err)
		cleanupInterval = 24 * time.Hour // fallback to default
	}

	// Set up all scheduler jobs in a slice for consistency and easier maintenance.
	aggregator := tracking.NewClickAggregator()
	jobs := []struct {
		name     string
		interval time.Duration
		job      scheduler.Job
	}{
		{
			name:     "persistance",
			interval: 10 * time.Second,
			job:      scheduler.FuncJob(persistence.PersistUnsaved),
		},
		{
			name:     "clean-up",
			interval: cleanupInterval,
			job:      scheduler.FuncJob(cleanup.CleanupStream),
		},
		{
			name:     "click-handling",
			interval: 30 * time.Second,
			job: scheduler.FuncJob(func() {
				aggregated := aggregator.Flush()
				if len(aggregated) > 0 {
					tracking.PersistAggregatedClicks(aggregated)
				} else {
					logger.Debug("No aggregated clicks to persist.")
				}
			}),
		},
		{
			name:     "aggregation",
			interval: 1 * time.Minute,
			job:      scheduler.FuncJob(tracking.AggregateClicks),
		},
	}

	// Start all schedulers and track them for shutdown.
	var schedulers []*scheduler.Scheduler
	for _, job := range jobs {
		s := scheduler.NewScheduler(job.name, job.interval, job.job, logger)
		schedulers = append(schedulers, s)
	}

	// Start the Redis stream consumer in its own goroutine.
	go tracking.RunConsumer(ctx, aggregator)

	serverErrChan := make(chan error, 1)

	// Start the main HTTP server
	go startMainServer(serverErrChan)
	// Start the metrics server
	go startMetricsServer(serverErrChan)

	// Handle shutdown signals.
	waitForShutdown(serverErrChan, schedulers, cancel)
	logger.Info("Server shut down successfully.")
}

// startMainServer runs the main HTTP server
func startMainServer(errChan chan<- error) {
	mainRouter := chi.NewRouter()
	// Add your HTTP endpoints here
	logger.Info("Main server running on " + cfg.ServerPort)
	errChan <- http.ListenAndServe(cfg.ServerPort, mainRouter)
}

// startMetricsServer runs the Prometheus metrics endpoint.
func startMetricsServer(errChan chan<- error) {
	cleanup.InitMetrics()
	metricsRouter := chi.NewRouter()
	metricsRouter.Handle("/metrics", promhttp.Handler())
	logger.Info("Prometheus metrics available on " + cfg.MetricsPort + "/metrics")
	errChan <- http.ListenAndServe(cfg.MetricsPort, metricsRouter)
}

// waitForShutdown handles graceful shutdown and resource cleanup.
func waitForShutdown(
	serverErrChan <-chan error,
	schedulers []*scheduler.Scheduler,
	cancel context.CancelFunc,
) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		logger.Error("Server error", "error", err)
	case <-stopChan:
		logger.Info("Shutdown signal received. Stopping scheduler and consumer...")
	}

	// Stop all schedulers gracefully.
	for _, s := range schedulers {
		s.Stop()
	}

	// Cancel the consumer context.
	cancel()
	time.Sleep(5 * time.Second) // Allow time for cleanup if needed.
}
