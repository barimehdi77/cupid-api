package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/barimehdi77/cupid-api/internal/sync"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SyncHandlers contains sync-related API handlers
type SyncHandlers struct {
	syncService *sync.SyncService
}

// NewSyncHandlers creates a new sync handlers instance
func NewSyncHandlers(syncService *sync.SyncService) *SyncHandlers {
	return &SyncHandlers{
		syncService: syncService,
	}
}

// TriggerSyncHandler handles manual sync trigger requests
// @Summary Trigger manual synchronization
// @Description Manually trigger a synchronization operation
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=SyncResult}
// @Failure 500 {object} APIResponse
// @Router /admin/sync [post]
func (h *SyncHandlers) TriggerSyncHandler(c *gin.Context) {
	logger.Info("Manual sync triggered via API")

	// Trigger sync in background
	go func() {
		ctx := c.Request.Context()
		result, err := h.syncService.SyncNow(ctx)
		if err != nil {
			logger.LogError("Manual sync failed", err)
		} else {
			logger.LogSuccess("Manual sync completed",
				zap.String("sync_id", result.SyncID),
				zap.Int("total_properties", result.TotalProperties),
				zap.Int("updated_properties", result.UpdatedProperties),
				zap.Duration("duration", result.Duration),
			)
		}
	}()

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":             "running",
			"message":            "Synchronization started in background",
			"estimated_duration": "5-10 minutes",
			"triggered_at":       time.Now(),
		},
	})
}

// GetSyncStatusHandler handles sync status requests
// @Summary Get sync status
// @Description Get the current status of the synchronization service
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=SyncStatus}
// @Router /admin/sync/status [get]
func (h *SyncHandlers) GetSyncStatusHandler(c *gin.Context) {
	status := h.syncService.GetStatus()

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    status,
	})
}

// StopSyncHandler handles sync stop requests
// @Summary Stop sync service
// @Description Stop the automatic synchronization service
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /admin/sync/stop [post]
func (h *SyncHandlers) StopSyncHandler(c *gin.Context) {
	logger.Info("Sync stop requested via API")

	err := h.syncService.Stop()
	if err != nil {
		logger.LogError("Failed to stop sync service", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to stop sync service",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":    "Sync service stopped successfully",
			"stopped_at": time.Now(),
			"status":     "stopped",
		},
	})
}

// StartSyncHandler handles sync start requests
// @Summary Start sync service
// @Description Start the automatic synchronization service
// @Tags admin
// @Accept json
// @Produce json
// @Param interval query string false "Sync interval (e.g., 12h, 24h)" default(12h)
// @Success 200 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /admin/sync/start [post]
func (h *SyncHandlers) StartSyncHandler(c *gin.Context) {
	intervalStr := c.DefaultQuery("interval", "12h")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid interval format. Use format like '12h' or '24h'",
		})
		return
	}

	logger.Info("Sync start requested via API",
		zap.String("interval", interval.String()),
	)

	ctx := c.Request.Context()
	err = h.syncService.Start(ctx)
	if err != nil {
		logger.LogError("Failed to start sync service", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to start sync service",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":    "Sync service started successfully",
			"started_at": time.Now(),
			"interval":   interval.String(),
			"status":     "running",
			"next_sync":  time.Now().Add(interval),
		},
	})
}

// GetSyncLogsHandler handles sync logs requests
// @Summary Get sync logs
// @Description Get synchronization operation logs
// @Tags admin
// @Accept json
// @Produce json
// @Param limit query int false "Number of logs to return" default(10)
// @Param offset query int false "Number of logs to skip" default(0)
// @Success 200 {object} APIResponse{data=[]SyncLog}
// @Router /admin/sync/logs [get]
func (h *SyncHandlers) GetSyncLogsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid limit. Must be between 1 and 100",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid offset. Must be >= 0",
		})
		return
	}

	// For now, return empty logs since we haven't implemented the storage layer
	// This would be implemented to fetch from sync_logs table
	logs := []sync.SyncLog{}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    logs,
		Meta: &Meta{
			Page:  (offset / limit) + 1,
			Limit: limit,
			Total: len(logs),
		},
	})
}

// GetSyncSettingsHandler handles sync settings requests
// @Summary Get sync settings
// @Description Get current synchronization settings
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=[]SyncSettings}
// @Router /admin/sync/settings [get]
func (h *SyncHandlers) GetSyncSettingsHandler(c *gin.Context) {
	// For now, return default settings
	// This would be implemented to fetch from sync_settings table
	settings := []sync.SyncSettings{
		{
			ID:           1,
			SettingKey:   "sync_interval",
			SettingValue: "12h",
			Description:  "Automatic sync interval",
		},
		{
			ID:           2,
			SettingKey:   "sync_batch_size",
			SettingValue: "10",
			Description:  "Number of properties to process in each batch",
		},
		{
			ID:           3,
			SettingKey:   "sync_max_concurrent",
			SettingValue: "5",
			Description:  "Maximum concurrent property fetches",
		},
		{
			ID:           4,
			SettingKey:   "sync_enable_auto",
			SettingValue: "true",
			Description:  "Enable automatic synchronization",
		},
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    settings,
	})
}

// UpdateSyncSettingsHandler handles sync settings update requests
// @Summary Update sync settings
// @Description Update synchronization settings
// @Tags admin
// @Accept json
// @Produce json
// @Param settings body []SyncSettings true "Sync settings to update"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Router /admin/sync/settings [put]
func (h *SyncHandlers) UpdateSyncSettingsHandler(c *gin.Context) {
	var settings []sync.SyncSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// For now, just log the settings update
	// This would be implemented to update sync_settings table
	logger.Info("Sync settings update requested",
		zap.Int("settings_count", len(settings)),
	)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":    "Sync settings updated successfully",
			"updated_at": time.Now(),
			"settings":   settings,
		},
	})
}

// GetSyncHealthHandler handles sync health check requests
// @Summary Get sync health
// @Description Get the health status of the synchronization service
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=map[string]interface{}}
// @Router /admin/sync/health [get]
func (h *SyncHandlers) GetSyncHealthHandler(c *gin.Context) {
	status := h.syncService.GetStatus()

	health := map[string]interface{}{
		"status":        "healthy",
		"is_running":    status.IsRunning,
		"is_healthy":    status.IsHealthy(),
		"is_overdue":    status.IsSyncOverdue(),
		"last_sync_age": status.GetSyncAge().String(),
		"next_sync_in":  status.GetNextSyncIn().String(),
		"sync_interval": status.SyncInterval,
		"summary":       status.GetSyncSummary(),
		"checked_at":    time.Now(),
	}

	// Determine overall health status
	if !status.IsHealthy() {
		health["status"] = "unhealthy"
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    health,
	})
}
