-- +goose Up
-- +goose StatementBegin
CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    property_id BIGINT NOT NULL REFERENCES properties(hotel_id) ON DELETE CASCADE,
    language VARCHAR(10) NOT NULL,
    hotel_name VARCHAR(255),
    description TEXT,
    markdown_description TEXT,
    important_info TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for common queries
CREATE INDEX idx_translations_property_id ON translations(property_id);
CREATE INDEX idx_translations_language ON translations(language);
CREATE UNIQUE INDEX idx_translations_property_language ON translations(property_id, language);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS translations;
-- +goose StatementEnd