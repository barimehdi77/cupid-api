-- +goose Up
-- +goose StatementBegin

-- Add sync tracking columns to properties table
ALTER TABLE properties ADD COLUMN last_synced TIMESTAMP DEFAULT NOW();
ALTER TABLE properties ADD COLUMN sync_status VARCHAR(20) DEFAULT 'pending';
ALTER TABLE properties ADD COLUMN data_version INTEGER DEFAULT 1;
ALTER TABLE properties ADD COLUMN last_updated TIMESTAMP DEFAULT NOW();

-- Add indexes for performance
CREATE INDEX idx_properties_last_synced ON properties(last_synced);
CREATE INDEX idx_properties_sync_status ON properties(sync_status);
CREATE INDEX idx_properties_data_version ON properties(data_version);

-- Create sync_logs table for tracking sync operations
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

-- Add indexes for sync_logs
CREATE INDEX idx_sync_logs_sync_id ON sync_logs(sync_id);
CREATE INDEX idx_sync_logs_status ON sync_logs(status);
CREATE INDEX idx_sync_logs_started_at ON sync_logs(started_at);

-- Create sync_settings table for configuration
CREATE TABLE sync_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(50) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default sync settings
INSERT INTO sync_settings (setting_key, setting_value, description) VALUES
('sync_interval', '12h', 'Automatic sync interval'),
('sync_batch_size', '10', 'Number of properties to process in each batch'),
('sync_max_concurrent', '5', 'Maximum concurrent property fetches'),
('sync_retry_attempts', '3', 'Number of retry attempts for failed operations'),
('sync_enable_auto', 'true', 'Enable automatic synchronization'),
('sync_rate_limit', '10', 'API requests per second limit');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop sync settings table
DROP TABLE IF EXISTS sync_settings;

-- Drop sync logs table
DROP TABLE IF EXISTS sync_logs;

-- Drop indexes
DROP INDEX IF EXISTS idx_properties_data_version;
DROP INDEX IF EXISTS idx_properties_sync_status;
DROP INDEX IF EXISTS idx_properties_last_synced;

-- Drop sync tracking columns
ALTER TABLE properties DROP COLUMN IF EXISTS last_updated;
ALTER TABLE properties DROP COLUMN IF EXISTS data_version;
ALTER TABLE properties DROP COLUMN IF EXISTS sync_status;
ALTER TABLE properties DROP COLUMN IF EXISTS last_synced;

-- +goose StatementEnd