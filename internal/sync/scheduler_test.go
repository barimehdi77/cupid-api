package sync

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSyncFunc is a mock sync function
type MockSyncFunc struct {
	mock.Mock
}

func (m *MockSyncFunc) Sync(ctx context.Context) (*SyncResult, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SyncResult), args.Error(1)
}

// TestNewScheduler tests the NewScheduler function
func TestNewScheduler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		interval := 1 * time.Hour
		mockSyncFunc := &MockSyncFunc{}

		// Act
		scheduler := NewScheduler(interval, mockSyncFunc.Sync)

		// Assert
		assert.NotNil(t, scheduler)
		assert.Equal(t, interval, scheduler.interval)
		assert.NotNil(t, scheduler.syncFunc)
		assert.NotNil(t, scheduler.stopChan)
		assert.False(t, scheduler.isRunning)
		assert.NotZero(t, scheduler.nextRun)
	})
}

// TestScheduler_IsRunning tests the IsRunning method
func TestScheduler_IsRunning(t *testing.T) {
	t.Run("InitialState", func(t *testing.T) {
		// Arrange
		interval := 1 * time.Hour
		mockSyncFunc := &MockSyncFunc{}
		scheduler := NewScheduler(interval, mockSyncFunc.Sync)

		// Act & Assert
		assert.False(t, scheduler.IsRunning())
	})
}

// TestScheduler_Stop tests the Stop method
func TestScheduler_Stop(t *testing.T) {
	t.Run("StopWhenNotRunning", func(t *testing.T) {
		// Arrange
		interval := 1 * time.Hour
		mockSyncFunc := &MockSyncFunc{}
		scheduler := NewScheduler(interval, mockSyncFunc.Sync)

		// Act
		scheduler.Stop()

		// Assert
		assert.False(t, scheduler.IsRunning())
	})
}

// TestScheduler_GetNextRun tests the GetNextRun method
func TestScheduler_GetNextRun(t *testing.T) {
	t.Run("GetNextRun", func(t *testing.T) {
		// Arrange
		interval := 1 * time.Hour
		mockSyncFunc := &MockSyncFunc{}
		scheduler := NewScheduler(interval, mockSyncFunc.Sync)

		// Act
		nextRun := scheduler.GetNextRun()

		// Assert
		assert.True(t, nextRun.After(time.Now()))
	})
}

// TestScheduler_Constructor tests the constructor with different intervals
func TestScheduler_Constructor(t *testing.T) {
	t.Run("WithDifferentIntervals", func(t *testing.T) {
		// Test with different intervals
		intervals := []time.Duration{
			1 * time.Hour,
			2 * time.Hour,
			30 * time.Minute,
			24 * time.Hour,
		}

		for _, interval := range intervals {
			// Arrange
			mockSyncFunc := &MockSyncFunc{}
			scheduler := NewScheduler(interval, mockSyncFunc.Sync)

			// Assert
			assert.NotNil(t, scheduler)
			assert.Equal(t, interval, scheduler.interval)
			assert.NotNil(t, scheduler.syncFunc)
			assert.NotNil(t, scheduler.stopChan)
			assert.False(t, scheduler.isRunning)
			assert.NotZero(t, scheduler.nextRun)
		}
	})
}
