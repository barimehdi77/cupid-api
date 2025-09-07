-- +goose Up
-- +goose StatementBegin
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    property_id BIGINT NOT NULL REFERENCES properties(hotel_id) ON DELETE CASCADE,
    review_id BIGINT NOT NULL,
    average_score INTEGER NOT NULL CHECK (average_score >= 1 AND average_score <= 5),
    country VARCHAR(100),
    type VARCHAR(50),
    name VARCHAR(255),
    date DATE,
    headline TEXT,
    language VARCHAR(10),
    pros TEXT,
    cons TEXT,
    source VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for common queries
CREATE INDEX idx_reviews_property_id ON reviews(property_id);
CREATE INDEX idx_reviews_average_score ON reviews(average_score);
CREATE INDEX idx_reviews_country ON reviews(country);
CREATE INDEX idx_reviews_language ON reviews(language);
CREATE INDEX idx_reviews_date ON reviews(date);
CREATE UNIQUE INDEX idx_reviews_property_review_id ON reviews(property_id, review_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reviews;
-- +goose StatementEnd