package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"go.uber.org/zap"
)

// StoreProperty stores a complete property with all its data
func (s *storage) StoreProperty(ctx context.Context, propertyData *cupid.PropertyData) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Store main property
	if err := s.storeMainProperty(ctx, tx, &propertyData.Property); err != nil {
		return fmt.Errorf("failed to store main property: %w", err)
	}

	// Store property details (JSONB data)
	if err := s.storePropertyDetails(ctx, tx, propertyData); err != nil {
		return fmt.Errorf("failed to store property details: %w", err)
	}

	// Store reviews
	if err := s.storeReviews(ctx, tx, propertyData.Property.HotelID, propertyData.Reviews); err != nil {
		return fmt.Errorf("failed to store reviews: %w", err)
	}

	// Store translations
	if err := s.storeTranslations(ctx, tx, propertyData.Property.HotelID, propertyData.Translations); err != nil {
		return fmt.Errorf("failed to store translations: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Property stored successfully",
		zap.Int64("hotel_id", propertyData.Property.HotelID),
		zap.String("hotel_name", propertyData.Property.HotelName),
	)

	return nil
}

// storeMainProperty stores the main property data
func (s *storage) storeMainProperty(ctx context.Context, tx *sql.Tx, property *cupid.Property) error {
	query := `
		INSERT INTO properties (
			hotel_id, cupid_id, hotel_name, hotel_type, hotel_type_id,
			chain, chain_id, latitude, longitude, stars, rating, review_count,
			airport_code, city, state, country, postal_code, main_image_th
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) ON CONFLICT (hotel_id) DO UPDATE SET
			cupid_id = EXCLUDED.cupid_id,
			hotel_name = EXCLUDED.hotel_name,
			hotel_type = EXCLUDED.hotel_type,
			hotel_type_id = EXCLUDED.hotel_type_id,
			chain = EXCLUDED.chain,
			chain_id = EXCLUDED.chain_id,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			stars = EXCLUDED.stars,
			rating = EXCLUDED.rating,
			review_count = EXCLUDED.review_count,
			airport_code = EXCLUDED.airport_code,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			country = EXCLUDED.country,
			postal_code = EXCLUDED.postal_code,
			main_image_th = EXCLUDED.main_image_th,
			updated_at = NOW()
	`

	_, err := tx.ExecContext(ctx, query,
		property.HotelID, property.CupidID, property.HotelName, property.HotelType, property.HotelTypeID,
		property.Chain, property.ChainID, property.Latitude, property.Longitude, property.Stars,
		property.Rating, property.ReviewCount, property.AirportCode, property.Address.City,
		property.Address.State, property.Address.Country, property.Address.PostalCode, property.MainImageTh,
	)

	return err
}

// storePropertyDetails stores complex data as JSONB
func (s *storage) storePropertyDetails(ctx context.Context, tx *sql.Tx, propertyData *cupid.PropertyData) error {
	// Prepare JSONB data
	details := map[string]interface{}{
		"address":    propertyData.Property.Address,
		"checkin":    propertyData.Property.CheckIn,
		"facilities": propertyData.Property.Facilities,
		"policies":   propertyData.Property.Policies,
		"rooms":      propertyData.Property.Rooms,
		"photos":     propertyData.Property.Photos,
		"contact_info": map[string]interface{}{
			"phone": propertyData.Property.Phone,
			"email": propertyData.Property.Email,
			"fax":   propertyData.Property.Fax,
		},
		"metadata": map[string]interface{}{
			"parking":        propertyData.Property.Parking,
			"group_room_min": propertyData.Property.GroupRoomMin,
			"child_allowed":  propertyData.Property.ChildAllowed,
			"pets_allowed":   propertyData.Property.PetsAllowed,
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal property details: %w", err)
	}

	query := `
		INSERT INTO property_details (property_id, address, checkin_info, facilities, policies, rooms, photos, contact_info, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (property_id) DO UPDATE SET
			address = EXCLUDED.address,
			checkin_info = EXCLUDED.checkin_info,
			facilities = EXCLUDED.facilities,
			policies = EXCLUDED.policies,
			rooms = EXCLUDED.rooms,
			photos = EXCLUDED.photos,
			contact_info = EXCLUDED.contact_info,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
	`

	_, err = tx.ExecContext(ctx, query,
		propertyData.Property.HotelID,
		jsonData, // address
		jsonData, // checkin_info
		jsonData, // facilities
		jsonData, // policies
		jsonData, // rooms
		jsonData, // photos
		jsonData, // contact_info
		jsonData, // metadata
	)

	return err
}

// storeReviews stores property reviews
func (s *storage) storeReviews(ctx context.Context, tx *sql.Tx, hotelID int64, reviews []cupid.Review) error {
	if len(reviews) == 0 {
		return nil
	}

	// Delete existing reviews for this property
	_, err := tx.ExecContext(ctx, "DELETE FROM reviews WHERE property_id = $1", hotelID)
	if err != nil {
		return fmt.Errorf("failed to delete existing reviews: %w", err)
	}

	// Insert new reviews
	query := `
		INSERT INTO reviews (property_id, review_id, average_score, country, type, name, date, headline, language, pros, cons, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	for _, review := range reviews {
		_, err := tx.ExecContext(ctx, query,
			hotelID, review.ReviewID, review.AverageScore, review.Country, review.Type,
			review.Name, review.Date, review.Headline, review.Language, review.Pros,
			review.Cons, review.Source,
		)
		if err != nil {
			return fmt.Errorf("failed to insert review: %w", err)
		}
	}

	return nil
}

// storeTranslations stores property translations
func (s *storage) storeTranslations(ctx context.Context, tx *sql.Tx, hotelID int64, translations map[string]*cupid.Property) error {
	if len(translations) == 0 {
		return nil
	}

	// Delete existing translations for this property
	_, err := tx.ExecContext(ctx, "DELETE FROM translations WHERE property_id = $1", hotelID)
	if err != nil {
		return fmt.Errorf("failed to delete existing translations: %w", err)
	}

	// Insert new translations
	query := `
		INSERT INTO translations (property_id, language, hotel_name, description, markdown_description, important_info)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for lang, translation := range translations {
		_, err := tx.ExecContext(ctx, query,
			hotelID, lang, translation.HotelName, translation.Description,
			translation.MarkdownDescription, translation.ImportantInfo,
		)
		if err != nil {
			return fmt.Errorf("failed to insert translation: %w", err)
		}
	}

	return nil
}
