package sync

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/barimehdi77/cupid-api/internal/store"
	"go.uber.org/zap"
)

// SyncService manages data synchronization between Cupid API and database
type SyncService struct {
	cupidService *cupid.Service
	storage      store.Storage
	scheduler    *Scheduler
	config       *Config
	isRunning    bool
	lastSync     time.Time
	stats        *SyncStats
	mu           sync.RWMutex
}

// Config holds synchronization configuration
type Config struct {
	Interval        time.Duration
	BatchSize       int
	MaxConcurrent   int
	RetryAttempts   int
	RetryDelay      time.Duration
	RateLimitPerSec int
	EnableAuto      bool
}

// DefaultConfig returns default synchronization configuration
func DefaultConfig() *Config {
	return &Config{
		Interval:        12 * time.Hour,
		BatchSize:       10,
		MaxConcurrent:   5,
		RetryAttempts:   3,
		RetryDelay:      5 * time.Second,
		RateLimitPerSec: 10,
		EnableAuto:      true,
	}
}

// NewSyncService creates a new synchronization service
func NewSyncService(cupidService *cupid.Service, storage store.Storage, config *Config) *SyncService {
	if config == nil {
		config = DefaultConfig()
	}

	return &SyncService{
		cupidService: cupidService,
		storage:      storage,
		config:       config,
		stats:        &SyncStats{},
	}
}

// Start begins the automatic synchronization scheduler
func (s *SyncService) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("sync service is already running")
	}

	if !s.config.EnableAuto {
		logger.Info("Automatic sync is disabled")
		return nil
	}

	s.scheduler = NewScheduler(s.config.Interval, s.performSync)
	s.isRunning = true

	logger.LogStartup("Sync Service",
		zap.Duration("interval", s.config.Interval),
		zap.Int("batch_size", s.config.BatchSize),
		zap.Int("max_concurrent", s.config.MaxConcurrent),
	)

	// Start scheduler in background
	go s.scheduler.Start(ctx)

	return nil
}

// Stop stops the automatic synchronization scheduler
func (s *SyncService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("sync service is not running")
	}

	if s.scheduler != nil {
		s.scheduler.Stop()
	}

	s.isRunning = false
	logger.LogShutdown("Sync Service", zap.String("reason", "manual stop"))

	return nil
}

// SyncNow performs an immediate synchronization
func (s *SyncService) SyncNow(ctx context.Context) (*SyncResult, error) {
	logger.Info("Starting manual synchronization")

	result, err := s.performSync(ctx)
	if err != nil {
		logger.LogError("Manual sync failed", err)
		return result, err
	}

	logger.LogSuccess("Manual sync completed",
		zap.Int("total_properties", result.TotalProperties),
		zap.Int("updated_properties", result.UpdatedProperties),
		zap.Int("failed_properties", result.FailedProperties),
		zap.Duration("duration", result.Duration),
	)

	return result, nil
}

// GetStatus returns the current synchronization status
func (s *SyncService) GetStatus() *SyncStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nextSync := time.Time{}
	if s.scheduler != nil {
		nextSync = s.scheduler.GetNextRun()
	}

	return &SyncStatus{
		IsRunning:         s.isRunning,
		LastSync:          s.lastSync,
		NextSync:          nextSync,
		TotalProperties:   s.stats.TotalProperties,
		UpdatedProperties: s.stats.UpdatedProperties,
		FailedProperties:  s.stats.FailedProperties,
		SyncInterval:      s.config.Interval.String(),
		LastError:         s.stats.LastError,
	}
}

// performSync performs the actual synchronization work
func (s *SyncService) performSync(ctx context.Context) (*SyncResult, error) {
	startTime := time.Now()
	syncID := fmt.Sprintf("sync_%s", startTime.Format("20060102_150405"))

	// Create sync log entry
	if err := s.createSyncLog(ctx, syncID, "running"); err != nil {
		logger.Warn("Failed to create sync log", zap.Error(err))
	}

	result := &SyncResult{
		SyncID:    syncID,
		StartTime: startTime,
		Status:    "running",
	}

	// Fetch all properties from Cupid API
	logger.Info("Fetching properties from Cupid API")
	properties, err := s.cupidService.FetchAllProperties(ctx)
	if err != nil {
		result.Status = "failed"
		result.Error = err
		s.updateSyncLog(ctx, syncID, "failed", err)
		return result, fmt.Errorf("failed to fetch properties: %w", err)
	}

	result.TotalProperties = len(properties)
	logger.Info("Fetched properties from API",
		zap.Int("count", len(properties)),
	)

	// Process properties in batches
	updatedCount := 0
	failedCount := 0

	for i := 0; i < len(properties); i += s.config.BatchSize {
		end := i + s.config.BatchSize
		if end > len(properties) {
			end = len(properties)
		}

		batch := properties[i:end]
		batchUpdated, batchFailed, err := s.processBatch(ctx, batch)
		if err != nil {
			logger.LogError("Failed to process batch", err,
				zap.Int("batch_start", i),
				zap.Int("batch_size", len(batch)),
			)
			failedCount += len(batch)
		} else {
			updatedCount += batchUpdated
			failedCount += batchFailed
		}
	}

	// Update result
	result.UpdatedProperties = updatedCount
	result.FailedProperties = failedCount
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Status = "completed"

	// Update sync log
	s.updateSyncLog(ctx, syncID, "completed", nil)

	// Update stats
	s.mu.Lock()
	s.lastSync = result.EndTime
	s.stats = &SyncStats{
		TotalProperties:   result.TotalProperties,
		UpdatedProperties: result.UpdatedProperties,
		FailedProperties:  result.FailedProperties,
		LastSync:          result.EndTime,
		LastError:         nil,
	}
	s.mu.Unlock()

	return result, nil
}

// processBatch processes a batch of properties
func (s *SyncService) processBatch(ctx context.Context, properties []*cupid.PropertyData) (int, int, error) {
	semaphore := make(chan struct{}, s.config.MaxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex

	updatedCount := 0
	failedCount := 0

	for _, propertyData := range properties {
		wg.Add(1)
		go func(pd *cupid.PropertyData) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			// Add rate limiting
			time.Sleep(time.Duration(1000/s.config.RateLimitPerSec) * time.Millisecond)

			// Compare and update property
			updated, err := s.compareAndUpdateProperty(ctx, pd)

			mu.Lock()
			if err != nil {
				failedCount++
				logger.LogError("Failed to update property", err,
					zap.Int64("property_id", pd.Property.HotelID),
				)
			} else if updated {
				updatedCount++
			}
			mu.Unlock()
		}(propertyData)
	}

	wg.Wait()
	return updatedCount, failedCount, nil
}

// compareAndUpdateProperty compares fetched data with stored data and updates if different
func (s *SyncService) compareAndUpdateProperty(ctx context.Context, fetchedData *cupid.PropertyData) (bool, error) {
	// Get stored property data
	storedData, err := s.storage.GetProperty(ctx, fetchedData.Property.HotelID)
	if err != nil {
		// Property doesn't exist, store it
		if err := s.storage.StoreProperty(ctx, fetchedData); err != nil {
			return false, fmt.Errorf("failed to store new property: %w", err)
		}
		return true, nil
	}

	// Compare data
	comparator := NewDataComparator()
	changes := comparator.ComparePropertyData(fetchedData, storedData)
	if !changes.HasChanges() {
		// No changes, just update sync timestamp
		return false, s.updateSyncTimestamp(ctx, fetchedData.Property.HotelID)
	}

	// Update property with changes
	if err := s.storage.StoreProperty(ctx, fetchedData); err != nil {
		return false, fmt.Errorf("failed to update property: %w", err)
	}

	logger.Debug("Property updated",
		zap.Int64("property_id", fetchedData.Property.HotelID),
		zap.Strings("changes", changes.Changes),
	)

	return true, nil
}

// updateSyncTimestamp updates the last_synced timestamp for a property
func (s *SyncService) updateSyncTimestamp(ctx context.Context, hotelID int64) error {
	// This would be implemented in the storage layer
	// For now, we'll just log it
	logger.Debug("Updating sync timestamp",
		zap.Int64("property_id", hotelID),
	)
	return nil
}

// createSyncLog creates a new sync log entry
func (s *SyncService) createSyncLog(ctx context.Context, syncID, status string) error {
	// This would be implemented in the storage layer
	logger.Debug("Creating sync log",
		zap.String("sync_id", syncID),
		zap.String("status", status),
	)
	return nil
}

// updateSyncLog updates a sync log entry
func (s *SyncService) updateSyncLog(ctx context.Context, syncID, status string, err error) {
	// This would be implemented in the storage layer
	logger.Debug("Updating sync log",
		zap.String("sync_id", syncID),
		zap.String("status", status),
		zap.Error(err),
	)
}
