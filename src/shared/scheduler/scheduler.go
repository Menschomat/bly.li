package scheduler

import (
	"context"
	"log/slog"
	"time"
)

// Job defines the contract for tasks that can be scheduled.
type Job interface {
	Run()
}
type FuncJob func()

func (f FuncJob) Run() { f() }

// Scheduler runs a Job periodically with proper resource management.
type Scheduler struct {
	ctx    context.Context
	name   string
	cancel context.CancelFunc
	ticker *time.Ticker
	job    Job
	logger *slog.Logger
}

// NewScheduler initializes a scheduler with a Job interface.
func NewScheduler(name string, interval time.Duration, job Job, logger *slog.Logger) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Scheduler{
		ctx:    ctx,
		name:   name,
		cancel: cancel,
		ticker: time.NewTicker(interval),
		job:    job,
		logger: logger,
	}
	s.logger.Info("Scheduler created", "name", s.name)
	go s.start()
	return s
}

// start runs the scheduled job at defined intervals.
func (s *Scheduler) start() {
	s.logger.Info("Scheduler starting", "name", s.name)
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Scheduler stopped", "name", s.name)
			return
		case <-s.ticker.C:
			go s.job.Run() // Run the Job in a separate goroutine
		}
	}
}

// Stop stops the scheduler and releases resources.
func (s *Scheduler) Stop() {
	s.cancel()
	s.ticker.Stop()
	s.logger.Info("Scheduler fully stopped", "name", s.name)
}
