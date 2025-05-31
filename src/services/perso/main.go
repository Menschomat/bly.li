package main

import (
	"context"
	"strings"
	"time"

	"github.com/Menschomat/bly.li/services/perso/cleanup"
	"github.com/Menschomat/bly.li/services/perso/logging"
	"github.com/Menschomat/bly.li/services/perso/persistence"
	"github.com/Menschomat/bly.li/services/perso/tracking"
	"github.com/Menschomat/bly.li/shared/config"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/scheduler"
	"github.com/Menschomat/bly.li/shared/server"
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
		cleanupInterval = 1 * time.Minute
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

	srv := server.NewServer(server.Config{
		ServerPort:         cfg.ServerPort,
		MetricsPort:        cfg.MetricsPort,
		CorsAllowedOrigins: strings.Split(cfg.CorsAllowedOrigins, ","),
		CorsMaxAge:         cfg.CorsMaxAge,
		Logger:             logger,
	})

	cleanup.InitMetrics()
	serverErrChan := make(chan error, 1)

	go srv.ServeMetrics(serverErrChan)
	go srv.ServeHTTP(serverErrChan)

	// Handle shutdown signals.
	go func() {
		srv.HandleShutdown(ctx, serverErrChan)
		// Stop all schedulers gracefully.
		for _, s := range schedulers {
			s.Stop()
		}
		// Allow time for cleanup if needed.
		time.Sleep(5 * time.Second)
	}()

	logger.Info("Server shut down successfully.")
}
