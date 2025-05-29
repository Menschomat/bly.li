package scheduler

import (
	"context"
	"log/slog"
	"time"
)

// Scheduler struct with a user-defined function
type Scheduler struct {
	ctx    context.Context
	name   string
	cancel context.CancelFunc
	ticker *time.Ticker
	job    func() // Function to execute
	logger *slog.Logger
}

// NewScheduler initializes a scheduler with a custom function
func NewScheduler(name string, interval time.Duration, job func(), logger *slog.Logger) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Scheduler{
		ctx:    ctx,
		name:   name,
		cancel: cancel,
		ticker: time.NewTicker(interval),
		job:    job, // Store the provided function
		logger: logger,
	}
	s.logger.Info("Scheduler created", "name", s.name)
	// Start scheduler in a separate goroutine
	go s.start()
	return s
}

// Start running tasks at intervals
func (s *Scheduler) start() {
	s.logger.Info("Scheduler starting", "name", s.name)
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Scheduler stopped", "name", s.name)
			return
		case <-s.ticker.C:
			// Run the user-defined function in a separate goroutine
			go s.job()
		}
	}
}

// Stop the scheduler
func (s *Scheduler) Stop() {
	s.cancel() // Cancel context
	s.ticker.Stop()
	s.logger.Info("Scheduler fully stopped", "name", s.name)
}
