-- +goose Up
-- +goose StatementBegin
CREATE TABLE property_details (
    property_id BIGINT PRIMARY KEY REFERENCES properties(hotel_id) ON DELETE CASCADE,
    address JSONB,
    checkin_info JSONB,
    facilities JSONB,
    policies JSONB,
    rooms JSONB,
    photos JSONB,
    contact_info JSONB,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create GIN indexes for JSONB fields to enable efficient queries
CREATE INDEX idx_property_details_address ON property_details USING GIN (address);
CREATE INDEX idx_property_details_facilities ON property_details USING GIN (facilities);
CREATE INDEX idx_property_details_policies ON property_details USING GIN (policies);
CREATE INDEX idx_property_details_rooms ON property_details USING GIN (rooms);
CREATE INDEX idx_property_details_photos ON property_details USING GIN (photos);
CREATE INDEX idx_property_details_contact_info ON property_details USING GIN (contact_info);
CREATE INDEX idx_property_details_metadata ON property_details USING GIN (metadata);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS property_details;
-- +goose StatementEnd