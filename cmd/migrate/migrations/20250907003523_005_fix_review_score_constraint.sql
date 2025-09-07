-- +goose Up
-- +goose StatementBegin
-- Drop the existing constraint
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_average_score_check;

-- Add new constraint for 10-point scale (1-10)
ALTER TABLE reviews ADD CONSTRAINT reviews_average_score_check CHECK (average_score >= 1 AND average_score <= 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the 10-point constraint
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_average_score_check;

-- Restore the original 5-point constraint
ALTER TABLE reviews ADD CONSTRAINT reviews_average_score_check CHECK (average_score >= 1 AND average_score <= 5);
-- +goose StatementEnd