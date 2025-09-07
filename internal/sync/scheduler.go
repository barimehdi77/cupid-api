package sync

import (
	"context"
	"sync"
	"time"

	"github.com/barimehdi77/cupid-api/internal/logger"
	"go.uber.org/zap"
)

// Scheduler manages automatic synchronization timing
type Scheduler struct {
	interval  time.Duration
	ticker    *time.Ticker
	stopChan  chan struct{}
	isRunning bool
	mu        sync.RWMutex
	nextRun   time.Time
	syncFunc  func(context.Context) (*SyncResult, error)
}

// NewScheduler creates a new scheduler
func NewScheduler(interval time.Duration, syncFunc func(context.Context) (*SyncResult, error)) *Scheduler {
	return &Scheduler{
		interval: interval,
		stopChan: make(chan struct{}),
		syncFunc: syncFunc,
		nextRun:  time.Now().Add(interval),
	}
}

// Start begins the scheduler
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return
	}
	s.isRunning = true
	s.mu.Unlock()

	s.ticker = time.NewTicker(s.interval)
	defer s.ticker.Stop()

	logger.Info("Scheduler started",
		zap.Duration("interval", s.interval),
		zap.Time("next_run", s.nextRun),
	)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Scheduler stopped due to context cancellation")
			return
		case <-s.stopChan:
			logger.Info("Scheduler stopped manually")
			return
		case <-s.ticker.C:
			s.runSync(ctx)
		}
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	close(s.stopChan)
	s.isRunning = false

	if s.ticker != nil {
		s.ticker.Stop()
	}

	logger.Info("Scheduler stopped")
}

// IsRunning returns whether the scheduler is running
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

// GetNextRun returns the next scheduled run time
func (s *Scheduler) GetNextRun() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nextRun
}

// runSync executes the synchronization function
func (s *Scheduler) runSync(ctx context.Context) {
	logger.Info("Starting scheduled synchronization")

	startTime := time.Now()
	result, err := s.syncFunc(ctx)
	duration := time.Since(startTime)

	if err != nil {
		logger.LogError("Scheduled sync failed", err,
			zap.Duration("duration", duration),
		)
	} else {
		logger.LogSuccess("Scheduled sync completed",
			zap.Int("total_properties", result.TotalProperties),
			zap.Int("updated_properties", result.UpdatedProperties),
			zap.Int("failed_properties", result.FailedProperties),
			zap.Duration("duration", duration),
		)
	}

	// Update next run time
	s.mu.Lock()
	s.nextRun = time.Now().Add(s.interval)
	s.mu.Unlock()

	logger.Debug("Next sync scheduled",
		zap.Time("next_run", s.nextRun),
	)
}
