# 🔄 Data Synchronization System

This document describes the comprehensive data synchronization system implemented for the Cupid API project.

## 📋 Overview

The sync system provides both **automatic scheduling** (12h or 24h intervals) and **manual triggering** via API endpoints to keep property data up-to-date with the Cupid API.

## 🏗️ Architecture

### Core Components

1. **SyncService** (`internal/sync/sync.go`) - Main orchestration service
2. **Scheduler** (`internal/sync/scheduler.go`) - Automatic timing management
3. **DataComparator** (`internal/sync/comparator.go`) - Smart data comparison
4. **SyncHandlers** (`internal/api/sync_handlers.go`) - API endpoints
5. **Database Schema** - Sync tracking tables and columns

### Key Features

✅ **Automatic Scheduling**: 12h or 24h intervals  
✅ **Manual API Control**: POST endpoints for immediate sync  
✅ **Smart Comparison**: Only updates changed data  
✅ **Background Processing**: Non-blocking operations  
✅ **Comprehensive Monitoring**: Status, health, logs, statistics  
✅ **Error Handling**: Robust retry and recovery mechanisms  
✅ **Rate Limiting**: Respects API limits  

## 🚀 Usage

### Automatic Mode

The sync service starts automatically when the API server starts:

```bash
# Start API server (sync starts automatically with 12h interval)
make run

# Check sync status
make sync-status

# Check sync health
make sync-health
```

### Manual Mode

```bash
# Trigger immediate sync
make sync-now

# Start sync with 12h interval
make sync-start

# Start sync with 24h interval
make sync-start-24h

# Stop sync service
make sync-stop

# Get sync logs
make sync-logs

# Get sync settings
make sync-settings
```

### API Endpoints

#### Sync Management
```bash
# Trigger immediate sync
curl -X POST http://localhost:8080/api/v1/admin/sync

# Get sync status
curl http://localhost:8080/api/v1/admin/sync/status

# Start sync with custom interval
curl -X POST "http://localhost:8080/api/v1/admin/sync/start?interval=24h"

# Stop sync service
curl -X POST http://localhost:8080/api/v1/admin/sync/stop

# Get sync health
curl http://localhost:8080/api/v1/admin/sync/health

# Get sync logs
curl "http://localhost:8080/api/v1/admin/sync/logs?limit=10&offset=0"

# Get sync settings
curl http://localhost:8080/api/v1/admin/sync/settings

# Update sync settings
curl -X PUT http://localhost:8080/api/v1/admin/sync/settings \
  -H "Content-Type: application/json" \
  -d '[{"setting_key": "sync_interval", "setting_value": "6h"}]'
```

## 📊 Response Examples

### Sync Status Response
```json
{
  "success": true,
  "data": {
    "is_running": true,
    "last_sync": "2024-01-15T10:30:00Z",
    "next_sync": "2024-01-15T22:30:00Z",
    "total_properties": 70,
    "updated_properties": 15,
    "failed_properties": 0,
    "sync_interval": "12h",
    "last_error": null
  }
}
```

### Sync Health Response
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "is_running": true,
    "is_healthy": true,
    "is_overdue": false,
    "last_sync_age": "2h30m",
    "next_sync_in": "9h30m",
    "sync_interval": "12h",
    "summary": "Sync service is running",
    "checked_at": "2024-01-15T13:00:00Z"
  }
}
```

### Manual Sync Response
```json
{
  "success": true,
  "data": {
    "status": "running",
    "message": "Synchronization started in background",
    "estimated_duration": "5-10 minutes",
    "triggered_at": "2024-01-15T13:00:00Z"
  }
}
```

## ⚙️ Configuration

### Environment Variables
```bash
# Sync configuration (optional - defaults provided)
SYNC_INTERVAL=12h                    # or 24h
SYNC_BATCH_SIZE=10
SYNC_MAX_CONCURRENT=5
SYNC_RETRY_ATTEMPTS=3
SYNC_ENABLE_AUTO=true
SYNC_RATE_LIMIT=10
```

### Database Schema

#### Properties Table (Updated)
```sql
-- New columns added for sync tracking
ALTER TABLE properties ADD COLUMN last_synced TIMESTAMP DEFAULT NOW();
ALTER TABLE properties ADD COLUMN sync_status VARCHAR(20) DEFAULT 'pending';
ALTER TABLE properties ADD COLUMN data_version INTEGER DEFAULT 1;
ALTER TABLE properties ADD COLUMN last_updated TIMESTAMP DEFAULT NOW();
```

#### Sync Logs Table
```sql
CREATE TABLE sync_logs (
    id SERIAL PRIMARY KEY,
    sync_id VARCHAR(50) UNIQUE NOT NULL,
    sync_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    total_properties INTEGER DEFAULT 0,
    updated_properties INTEGER DEFAULT 0,
    failed_properties INTEGER DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### Sync Settings Table
```sql
CREATE TABLE sync_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(50) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 🔍 How It Works

### 1. Data Comparison Process

The system uses intelligent comparison to detect changes:

```go
// Compare property data
comparator := NewDataComparator()
changes := comparator.ComparePropertyData(fetchedData, storedData)

if changes.HasChanges() {
    // Update only if changes detected
    updateProperty(fetchedData)
}
```

### 2. Batch Processing

Properties are processed in configurable batches:

```go
// Process 10 properties at a time with max 5 concurrent
for i := 0; i < len(properties); i += batchSize {
    batch := properties[i:i+batchSize]
    processBatch(batch)
}
```

### 3. Rate Limiting

API calls are rate-limited to respect external API limits:

```go
// 10 requests per second
time.Sleep(time.Duration(1000/rateLimit) * time.Millisecond)
```

### 4. Error Handling

Robust error handling with retry mechanisms:

```go
// Retry failed operations up to 3 times
for attempt := 1; attempt <= retryAttempts; attempt++ {
    if err := operation(); err == nil {
        break
    }
    time.Sleep(retryDelay)
}
```

## 📈 Monitoring & Statistics

### Sync Metrics
- **Total Properties**: Number of properties processed
- **Updated Properties**: Number of properties with changes
- **Failed Properties**: Number of properties that failed to update
- **Success Rate**: Percentage of successful updates
- **Duration**: Time taken for sync operations

### Health Indicators
- **Is Running**: Whether sync service is active
- **Is Healthy**: Overall health status
- **Is Overdue**: Whether sync is behind schedule
- **Last Sync Age**: Time since last successful sync
- **Next Sync In**: Time until next scheduled sync

## 🛠️ Development

### Running Tests
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific sync tests
go test ./internal/sync/...
```

### Adding New Features

1. **New Comparison Logic**: Add methods to `DataComparator`
2. **New API Endpoints**: Add handlers to `SyncHandlers`
3. **New Configuration**: Add settings to `Config` struct
4. **New Monitoring**: Add metrics to `SyncStats`

## 🔧 Troubleshooting

### Common Issues

1. **Sync Not Starting**
   ```bash
   # Check if sync is enabled
   curl http://localhost:8080/api/v1/admin/sync/status
   
   # Check logs
   make logs
   ```

2. **High Failure Rate**
   ```bash
   # Check sync health
   make sync-health
   
   # Check API connectivity
   curl http://localhost:8080/api/v1/health
   ```

3. **Sync Overdue**
   ```bash
   # Check if service is running
   make sync-status
   
   # Restart sync service
   make sync-stop
   make sync-start
   ```

### Debug Mode

Enable debug logging for detailed sync information:

```bash
# Set log level to debug
export LOG_LEVEL=debug
make run
```

## 📚 API Documentation

Full API documentation is available at:
- **Swagger UI**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/docs/swagger.json

## 🎯 Benefits

✅ **Efficient**: Only updates changed data  
✅ **Reliable**: Robust error handling and retry logic  
✅ **Scalable**: Handles growing datasets efficiently  
✅ **Monitorable**: Comprehensive status and health monitoring  
✅ **Flexible**: Both automatic and manual control  
✅ **Non-intrusive**: Doesn't affect existing functionality  

## 🚀 Next Steps

1. **Enhanced Monitoring**: Add Prometheus metrics
2. **Webhook Notifications**: Notify on sync completion
3. **Data Validation**: Add data quality checks
4. **Performance Optimization**: Implement caching strategies
5. **Multi-tenant Support**: Per-tenant sync configurations

---

**Ready to sync! 🎉** The system is now fully operational with both automatic scheduling and manual control capabilities.
