-- +goose Up
-- +goose StatementBegin
CREATE TABLE properties (
    hotel_id BIGINT PRIMARY KEY,
    cupid_id BIGINT,
    hotel_name VARCHAR(255) NOT NULL,
    hotel_type VARCHAR(100),
    hotel_type_id INTEGER,
    chain VARCHAR(100),
    chain_id INTEGER,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    stars INTEGER,
    rating DECIMAL(3, 2),
    review_count INTEGER DEFAULT 0,
    airport_code VARCHAR(10),
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    main_image_th TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for common queries
CREATE INDEX idx_properties_city ON properties(city);
CREATE INDEX idx_properties_country ON properties(country);
CREATE INDEX idx_properties_stars ON properties(stars);
CREATE INDEX idx_properties_rating ON properties(rating);
CREATE INDEX idx_properties_hotel_type ON properties(hotel_type);
CREATE INDEX idx_properties_chain ON properties(chain);
CREATE INDEX idx_properties_location ON properties(latitude, longitude);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS properties;
-- +goose StatementEnd