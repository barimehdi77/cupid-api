package sync

import (
	"time"
)

// SyncStats represents synchronization statistics
type SyncStats struct {
	TotalProperties   int       `json:"total_properties"`
	UpdatedProperties int       `json:"updated_properties"`
	FailedProperties  int       `json:"failed_properties"`
	LastSync          time.Time `json:"last_sync"`
	LastError         error     `json:"last_error,omitempty"`
}

// SyncResult represents the result of a synchronization operation
type SyncResult struct {
	SyncID            string        `json:"sync_id"`
	Status            string        `json:"status"`
	StartTime         time.Time     `json:"start_time"`
	EndTime           time.Time     `json:"end_time"`
	Duration          time.Duration `json:"duration"`
	TotalProperties   int           `json:"total_properties"`
	UpdatedProperties int           `json:"updated_properties"`
	FailedProperties  int           `json:"failed_properties"`
	Error             error         `json:"error,omitempty"`
}

// SyncStatus represents the current status of the sync service
type SyncStatus struct {
	IsRunning         bool      `json:"is_running"`
	LastSync          time.Time `json:"last_sync"`
	NextSync          time.Time `json:"next_sync"`
	TotalProperties   int       `json:"total_properties"`
	UpdatedProperties int       `json:"updated_properties"`
	FailedProperties  int       `json:"failed_properties"`
	SyncInterval      string    `json:"sync_interval"`
	LastError         error     `json:"last_error,omitempty"`
}

// SyncLog represents a sync operation log entry
type SyncLog struct {
	ID                int        `json:"id"`
	SyncID            string     `json:"sync_id"`
	SyncType          string     `json:"sync_type"`
	Status            string     `json:"status"`
	StartedAt         time.Time  `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	TotalProperties   int        `json:"total_properties"`
	UpdatedProperties int        `json:"updated_properties"`
	FailedProperties  int        `json:"failed_properties"`
	ErrorMessage      string     `json:"error_message,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

// SyncSettings represents sync configuration settings
type SyncSettings struct {
	ID           int       `json:"id"`
	SettingKey   string    `json:"setting_key"`
	SettingValue string    `json:"setting_value"`
	Description  string    `json:"description"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetSuccessRate calculates the success rate of the sync operation
func (sr *SyncResult) GetSuccessRate() float64 {
	if sr.TotalProperties == 0 {
		return 0.0
	}
	return float64(sr.UpdatedProperties) / float64(sr.TotalProperties) * 100.0
}

// GetFailureRate calculates the failure rate of the sync operation
func (sr *SyncResult) GetFailureRate() float64 {
	if sr.TotalProperties == 0 {
		return 0.0
	}
	return float64(sr.FailedProperties) / float64(sr.TotalProperties) * 100.0
}

// IsSuccessful returns true if the sync operation was successful
func (sr *SyncResult) IsSuccessful() bool {
	return sr.Status == "completed" && sr.Error == nil
}

// GetDurationString returns the duration as a formatted string
func (sr *SyncResult) GetDurationString() string {
	if sr.Duration < time.Minute {
		return sr.Duration.Round(time.Second).String()
	}
	return sr.Duration.Round(time.Second).String()
}

// GetSyncAge returns the age of the last sync
func (ss *SyncStatus) GetSyncAge() time.Duration {
	if ss.LastSync.IsZero() {
		return 0
	}
	return time.Since(ss.LastSync)
}

// IsSyncOverdue returns true if the sync is overdue
func (ss *SyncStatus) IsSyncOverdue() bool {
	if !ss.IsRunning && !ss.LastSync.IsZero() {
		// If not running and last sync was more than 2x the interval ago
		interval, _ := time.ParseDuration(ss.SyncInterval)
		return time.Since(ss.LastSync) > interval*2
	}
	return false
}

// GetNextSyncIn returns the time until the next sync
func (ss *SyncStatus) GetNextSyncIn() time.Duration {
	if ss.NextSync.IsZero() {
		return 0
	}
	return time.Until(ss.NextSync)
}

// IsHealthy returns true if the sync service is healthy
func (ss *SyncStatus) IsHealthy() bool {
	// Service is healthy if:
	// 1. It's running, OR
	// 2. It's not running but last sync was recent (within 2x interval)
	return ss.IsRunning || !ss.IsSyncOverdue()
}

// GetUptime returns the uptime of the sync service
func (ss *SyncStatus) GetUptime() time.Duration {
	if ss.LastSync.IsZero() {
		return 0
	}
	return time.Since(ss.LastSync)
}

// GetSyncFrequency returns the sync frequency as a human-readable string
func (ss *SyncStatus) GetSyncFrequency() string {
	return ss.SyncInterval
}

// GetSyncSummary returns a summary of the sync status
func (ss *SyncStatus) GetSyncSummary() string {
	if ss.IsRunning {
		return "Sync service is running"
	}

	if ss.IsSyncOverdue() {
		return "Sync service is overdue"
	}

	if ss.LastSync.IsZero() {
		return "Sync service has never run"
	}

	return "Sync service is healthy"
}

// GetSyncMetrics returns key metrics for monitoring
func (ss *SyncStatus) GetSyncMetrics() map[string]interface{} {
	return map[string]interface{}{
		"is_running":         ss.IsRunning,
		"is_healthy":         ss.IsHealthy(),
		"is_overdue":         ss.IsSyncOverdue(),
		"last_sync_age":      ss.GetSyncAge().String(),
		"next_sync_in":       ss.GetNextSyncIn().String(),
		"total_properties":   ss.TotalProperties,
		"updated_properties": ss.UpdatedProperties,
		"failed_properties":  ss.FailedProperties,
		"sync_interval":      ss.SyncInterval,
		"summary":            ss.GetSyncSummary(),
	}
}

// GetSyncLogSummary returns a summary of a sync log entry
func (sl *SyncLog) GetSyncLogSummary() string {
	if sl.Status == "completed" {
		return "Sync completed successfully"
	}
	if sl.Status == "failed" {
		return "Sync failed"
	}
	if sl.Status == "running" {
		return "Sync in progress"
	}
	return "Sync status unknown"
}

// GetSyncLogDuration returns the duration of the sync operation
func (sl *SyncLog) GetSyncLogDuration() time.Duration {
	if sl.CompletedAt == nil {
		return time.Since(sl.StartedAt)
	}
	return sl.CompletedAt.Sub(sl.StartedAt)
}

// IsSyncLogSuccessful returns true if the sync log represents a successful operation
func (sl *SyncLog) IsSyncLogSuccessful() bool {
	return sl.Status == "completed" && sl.ErrorMessage == ""
}

// GetSyncLogSuccessRate calculates the success rate from the sync log
func (sl *SyncLog) GetSyncLogSuccessRate() float64 {
	if sl.TotalProperties == 0 {
		return 0.0
	}
	return float64(sl.UpdatedProperties) / float64(sl.TotalProperties) * 100.0
}

// GetSyncLogFailureRate calculates the failure rate from the sync log
func (sl *SyncLog) GetSyncLogFailureRate() float64 {
	if sl.TotalProperties == 0 {
		return 0.0
	}
	return float64(sl.FailedProperties) / float64(sl.TotalProperties) * 100.0
}
