package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/barimehdi77/cupid-api/internal/cupid"
)

// GetReviewsByScore retrieves reviews within a score range
func (s *storage) GetReviewsByScore(ctx context.Context, minScore, maxScore int, limit, offset int) ([]cupid.Review, error) {
	query := `
		SELECT r.review_id, r.average_score, r.country, r.type, r.name, r.date, r.headline, r.language, r.pros, r.cons, r.source
		FROM reviews r
		WHERE r.average_score >= $1 AND r.average_score <= $2
		ORDER BY r.average_score DESC, r.date DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := s.db.QueryContext(ctx, query, minScore, maxScore, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []cupid.Review
	for rows.Next() {
		var review cupid.Review
		err := rows.Scan(
			&review.ReviewID, &review.AverageScore, &review.Country, &review.Type,
			&review.Name, &review.Date, &review.Headline, &review.Language,
			&review.Pros, &review.Cons, &review.Source,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// GetTranslationByLanguage retrieves a specific translation
func (s *storage) GetTranslationByLanguage(ctx context.Context, hotelID int64, language string) (*cupid.Property, error) {
	query := `
		SELECT hotel_name, description, markdown_description, important_info
		FROM translations
		WHERE property_id = $1 AND language = $2
	`

	var translation cupid.Property
	err := s.db.QueryRowContext(ctx, query, hotelID, language).Scan(
		&translation.HotelName, &translation.Description,
		&translation.MarkdownDescription, &translation.ImportantInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("translation not found")
		}
		return nil, err
	}

	return &translation, nil
}

// SearchProperties performs a text search on properties
func (s *storage) SearchProperties(ctx context.Context, query string, limit, offset int) ([]*cupid.Property, error) {
	searchQuery := `
		SELECT hotel_id, cupid_id, hotel_name, hotel_type, hotel_type_id,
			   chain, chain_id, latitude, longitude, stars, rating, review_count,
			   airport_code, city, state, country, postal_code, main_image_th
		FROM properties
		WHERE hotel_name ILIKE $1 OR city ILIKE $1 OR country ILIKE $1
		ORDER BY rating DESC, review_count DESC
		LIMIT $2 OFFSET $3
	`

	searchTerm := "%" + query + "%"
	rows, err := s.db.QueryContext(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []*cupid.Property
	for rows.Next() {
		var property cupid.Property
		err := rows.Scan(
			&property.HotelID, &property.CupidID, &property.HotelName, &property.HotelType, &property.HotelTypeID,
			&property.Chain, &property.ChainID, &property.Latitude, &property.Longitude, &property.Stars,
			&property.Rating, &property.ReviewCount, &property.AirportCode, &property.Address.City,
			&property.Address.State, &property.Address.Country, &property.Address.PostalCode, &property.MainImageTh,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, &property)
	}

	return properties, nil
}

// GetPropertiesByLocation retrieves properties by location
func (s *storage) GetPropertiesByLocation(ctx context.Context, city, country string, limit, offset int) ([]*cupid.Property, error) {
	filters := PropertyFilters{
		City:    city,
		Country: country,
	}
	return s.ListProperties(ctx, limit, offset, filters)
}

// GetPropertiesByRating retrieves properties by minimum rating
func (s *storage) GetPropertiesByRating(ctx context.Context, minRating float64, limit, offset int) ([]*cupid.Property, error) {
	filters := PropertyFilters{
		MinRating: minRating,
	}
	return s.ListProperties(ctx, limit, offset, filters)
}
