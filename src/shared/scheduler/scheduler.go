package scheduler

import (
	"context"
	"fmt"
	"time"
)

// Scheduler struct with a user-defined function
type Scheduler struct {
	ctx    context.Context
	cancel context.CancelFunc
	ticker *time.Ticker
	job    func() // Function to execute
}

// NewScheduler initializes a scheduler with a custom function
func NewScheduler(interval time.Duration, job func()) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Scheduler{
		ctx:    ctx,
		cancel: cancel,
		ticker: time.NewTicker(interval),
		job:    job, // Store the provided function
	}

	// Start scheduler in a separate goroutine
	go s.start()
	return s
}

// Start running tasks at intervals
func (s *Scheduler) start() {
	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("Scheduler stopped.")
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
	fmt.Println("Scheduler fully stopped.")
}
