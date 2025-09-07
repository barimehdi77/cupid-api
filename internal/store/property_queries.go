package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/barimehdi77/cupid-api/internal/cupid"
)

// GetProperty retrieves a complete property with all its data
func (s *storage) GetProperty(ctx context.Context, hotelID int64) (*cupid.PropertyData, error) {
	// Get main property
	property, err := s.getMainProperty(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	// Get reviews
	reviews, err := s.GetPropertyReviews(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	// Get translations
	translations, err := s.GetPropertyTranslations(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	return &cupid.PropertyData{
		Property:     *property,
		Reviews:      reviews,
		Translations: translations,
	}, nil
}

// getMainProperty retrieves the main property data
func (s *storage) getMainProperty(ctx context.Context, hotelID int64) (*cupid.Property, error) {
	query := `
		SELECT hotel_id, cupid_id, hotel_name, hotel_type, hotel_type_id,
			   chain, chain_id, latitude, longitude, stars, rating, review_count,
			   airport_code, city, state, country, postal_code, main_image_th
		FROM properties
		WHERE hotel_id = $1
	`

	var property cupid.Property
	err := s.db.QueryRowContext(ctx, query, hotelID).Scan(
		&property.HotelID, &property.CupidID, &property.HotelName, &property.HotelType, &property.HotelTypeID,
		&property.Chain, &property.ChainID, &property.Latitude, &property.Longitude, &property.Stars,
		&property.Rating, &property.ReviewCount, &property.AirportCode, &property.Address.City,
		&property.Address.State, &property.Address.Country, &property.Address.PostalCode, &property.MainImageTh,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("property not found")
		}
		return nil, err
	}

	return &property, nil
}

// ListProperties retrieves a list of properties with optional filtering
func (s *storage) ListProperties(ctx context.Context, limit, offset int, filters PropertyFilters) ([]*cupid.Property, error) {
	query := `
		SELECT hotel_id, cupid_id, hotel_name, hotel_type, hotel_type_id,
			   chain, chain_id, latitude, longitude, stars, rating, review_count,
			   airport_code, city, state, country, postal_code, main_image_th
		FROM properties
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filters.City != "" {
		query += fmt.Sprintf(" AND city ILIKE $%d", argIndex)
		args = append(args, "%"+filters.City+"%")
		argIndex++
	}

	if filters.Country != "" {
		query += fmt.Sprintf(" AND country ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Country+"%")
		argIndex++
	}

	if filters.MinStars > 0 {
		query += fmt.Sprintf(" AND stars >= $%d", argIndex)
		args = append(args, filters.MinStars)
		argIndex++
	}

	if filters.MaxStars > 0 {
		query += fmt.Sprintf(" AND stars <= $%d", argIndex)
		args = append(args, filters.MaxStars)
		argIndex++
	}

	if filters.MinRating > 0 {
		query += fmt.Sprintf(" AND rating >= $%d", argIndex)
		args = append(args, filters.MinRating)
		argIndex++
	}

	if filters.MaxRating > 0 {
		query += fmt.Sprintf(" AND rating <= $%d", argIndex)
		args = append(args, filters.MaxRating)
		argIndex++
	}

	if filters.HotelType != "" {
		query += fmt.Sprintf(" AND hotel_type ILIKE $%d", argIndex)
		args = append(args, "%"+filters.HotelType+"%")
		argIndex++
	}

	if filters.Chain != "" {
		query += fmt.Sprintf(" AND chain ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Chain+"%")
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY rating DESC, review_count DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
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

// CountProperties counts the total number of properties matching the given filters
func (s *storage) CountProperties(ctx context.Context, filters PropertyFilters) (int, error) {
	query := "SELECT COUNT(*) FROM properties WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Add filters
	if filters.City != "" {
		query += fmt.Sprintf(" AND city ILIKE $%d", argIndex)
		args = append(args, "%"+filters.City+"%")
		argIndex++
	}

	if filters.Country != "" {
		query += fmt.Sprintf(" AND country ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Country+"%")
		argIndex++
	}

	if filters.MinStars > 0 {
		query += fmt.Sprintf(" AND stars >= $%d", argIndex)
		args = append(args, filters.MinStars)
		argIndex++
	}

	if filters.MaxStars > 0 {
		query += fmt.Sprintf(" AND stars <= $%d", argIndex)
		args = append(args, filters.MaxStars)
		argIndex++
	}

	if filters.MinRating > 0 {
		query += fmt.Sprintf(" AND rating >= $%d", argIndex)
		args = append(args, filters.MinRating)
		argIndex++
	}

	if filters.MaxRating > 0 {
		query += fmt.Sprintf(" AND rating <= $%d", argIndex)
		args = append(args, filters.MaxRating)
		argIndex++
	}

	if filters.HotelType != "" {
		query += fmt.Sprintf(" AND hotel_type ILIKE $%d", argIndex)
		args = append(args, "%"+filters.HotelType+"%")
		argIndex++
	}

	if filters.Chain != "" {
		query += fmt.Sprintf(" AND chain ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Chain+"%")
		argIndex++
	}

	var count int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count properties: %w", err)
	}

	return count, nil
}

// GetPropertyReviews retrieves reviews for a specific property
func (s *storage) GetPropertyReviews(ctx context.Context, hotelID int64) ([]cupid.Review, error) {
	query := `
		SELECT review_id, average_score, country, type, name, date, headline, language, pros, cons, source
		FROM reviews
		WHERE property_id = $1
		ORDER BY date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, hotelID)
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

// GetPropertyTranslations retrieves all translations for a specific property
func (s *storage) GetPropertyTranslations(ctx context.Context, hotelID int64) (map[string]*cupid.Property, error) {
	query := `
		SELECT language, hotel_name, description, markdown_description, important_info
		FROM translations
		WHERE property_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	translations := make(map[string]*cupid.Property)
	for rows.Next() {
		var lang string
		var translation cupid.Property
		err := rows.Scan(
			&lang, &translation.HotelName, &translation.Description,
			&translation.MarkdownDescription, &translation.ImportantInfo,
		)
		if err != nil {
			return nil, err
		}
		translations[lang] = &translation
	}

	return translations, nil
}

// UpdateProperty updates an existing property
func (s *storage) UpdateProperty(ctx context.Context, hotelID int64, propertyData *cupid.PropertyData) error {
	return s.StoreProperty(ctx, propertyData)
}

// DeleteProperty deletes a property and all its related data
func (s *storage) DeleteProperty(ctx context.Context, hotelID int64) error {
	query := "DELETE FROM properties WHERE hotel_id = $1"
	_, err := s.db.ExecContext(ctx, query, hotelID)
	return err
}
