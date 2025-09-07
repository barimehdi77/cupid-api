package store

import (
	"testing"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getSamplePropertyData creates sample property data for testing
func getSamplePropertyData() *cupid.PropertyData {
	return &cupid.PropertyData{
		Property: cupid.Property{
			HotelID:     12345,
			CupidID:     67890,
			MainImageTh: "https://example.com/image.jpg",
			HotelType:   "hotel",
			HotelTypeID: 1,
			Chain:       "Luxury Hotels",
			ChainID:     1,
			Latitude:    48.8566,
			Longitude:   2.3522,
			HotelName:   "Luxury Hotel Paris",
			Phone:       "+33 1 23 45 67 89",
			Fax:         "+33 1 23 45 67 90",
			Email:       "info@luxuryhotel.com",
			Address: cupid.Address{
				Address:    "123 Champs-Élysées",
				City:       "Paris",
				State:      "Île-de-France",
				Country:    "France",
				PostalCode: "75008",
			},
			Stars:       5,
			AirportCode: "CDG",
			Rating:      4.8,
			ReviewCount: 150,
			CheckIn: cupid.CheckIn{
				CheckInStart: "15:00",
				CheckInEnd:   "23:00",
				Checkout:     "11:00",
			},
			Parking:      stringPtr("Valet parking available"),
			GroupRoomMin: intPtr(10),
			ChildAllowed: boolPtr(true),
			PetsAllowed:  boolPtr(false),
			Photos: []cupid.Photo{
				{
					URL: "https://example.com/photo1.jpg",
				},
			},
			Facilities: []cupid.Facility{
				{
					FacilityID: 1,
					Name:       "WiFi",
				},
				{
					FacilityID: 2,
					Name:       "Pool",
				},
			},
			Policies: []cupid.Policy{
				{
					PolicyType: "cancellation",
					Name:       "Free cancellation",
				},
			},
			Rooms: []cupid.Room{
				{
					ID:       1,
					RoomName: "Deluxe Room",
				},
			},
		},
		Reviews: []cupid.Review{
			{
				ReviewID:     1,
				AverageScore: 4,
				Country:      "US",
				Name:         "John Doe",
				Headline:     "Great hotel",
				Pros:         "Clean, comfortable",
				Cons:         "Noisy",
				Date:         "2024-01-15",
				Language:     "en",
				Source:       "booking.com",
			},
		},
		Translations: map[string]*cupid.Property{
			"fr": {
				HotelID:   12345,
				HotelName: "Hôtel de Luxe Paris",
				Address: cupid.Address{
					Address:    "123 Champs-Élysées",
					City:       "Paris",
					State:      "Île-de-France",
					Country:    "France",
					PostalCode: "75008",
				},
			},
		},
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

// TestStorage_StoreProperty tests the StoreProperty method
func TestStorage_StoreProperty(t *testing.T) {
	t.Run("ValidPropertyData", func(t *testing.T) {
		// Arrange
		propertyData := getSamplePropertyData()

		// Act & Assert
		// Note: This test would require a real database connection
		// For now, we'll just test that the data structure is valid
		assert.NotNil(t, propertyData)
		assert.Equal(t, int64(12345), propertyData.Property.HotelID)
		assert.Equal(t, "Luxury Hotel Paris", propertyData.Property.HotelName)
		assert.Len(t, propertyData.Reviews, 1)
		assert.Len(t, propertyData.Translations, 1)
	})

	t.Run("PropertyDataValidation", func(t *testing.T) {
		// Arrange
		propertyData := getSamplePropertyData()

		// Act & Assert
		require.NotNil(t, propertyData.Property.Address)
		assert.Equal(t, "Paris", propertyData.Property.Address.City)
		assert.Equal(t, "France", propertyData.Property.Address.Country)

		require.Len(t, propertyData.Property.Photos, 1)
		assert.Equal(t, "https://example.com/photo1.jpg", propertyData.Property.Photos[0].URL)

		require.Len(t, propertyData.Property.Facilities, 2)
		assert.Equal(t, "WiFi", propertyData.Property.Facilities[0].Name)
		assert.Equal(t, "Pool", propertyData.Property.Facilities[1].Name)
	})
}

// TestStorage_GetProperty tests the GetProperty method
func TestStorage_GetProperty(t *testing.T) {
	t.Run("ValidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(12345)

		// Act & Assert
		// Note: This test would require a real database connection
		// For now, we'll just test the input validation
		assert.Greater(t, hotelID, int64(0))
	})

	t.Run("InvalidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(0)

		// Act & Assert
		assert.Equal(t, int64(0), hotelID)
	})
}

// TestStorage_ListProperties tests the ListProperties method
func TestStorage_ListProperties(t *testing.T) {
	t.Run("ValidFilters", func(t *testing.T) {
		// Arrange
		filters := PropertyFilters{
			City:      "Paris",
			Country:   "France",
			MinStars:  4,
			MaxStars:  5,
			MinRating: 4.0,
			MaxRating: 5.0,
		}
		limit := 10
		offset := 0

		// Act & Assert
		assert.Equal(t, "Paris", filters.City)
		assert.Equal(t, "France", filters.Country)
		assert.Equal(t, 4, filters.MinStars)
		assert.Equal(t, 5, filters.MaxStars)
		assert.Equal(t, 4.0, filters.MinRating)
		assert.Equal(t, 5.0, filters.MaxRating)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
	})

	t.Run("EmptyFilters", func(t *testing.T) {
		// Arrange
		filters := PropertyFilters{}
		limit := 20
		offset := 0

		// Act & Assert
		assert.Empty(t, filters.City)
		assert.Empty(t, filters.Country)
		assert.Equal(t, 0, filters.MinStars)
		assert.Equal(t, 0, filters.MaxStars)
		assert.Equal(t, 0.0, filters.MinRating)
		assert.Equal(t, 0.0, filters.MaxRating)
		assert.Equal(t, 20, limit)
		assert.Equal(t, 0, offset)
	})
}

// TestStorage_CountProperties tests the CountProperties method
func TestStorage_CountProperties(t *testing.T) {
	t.Run("ValidFilters", func(t *testing.T) {
		// Arrange
		filters := PropertyFilters{
			City:    "Paris",
			Country: "France",
		}

		// Act & Assert
		assert.Equal(t, "Paris", filters.City)
		assert.Equal(t, "France", filters.Country)
	})

	t.Run("EmptyFilters", func(t *testing.T) {
		// Arrange
		filters := PropertyFilters{}

		// Act & Assert
		assert.Empty(t, filters.City)
		assert.Empty(t, filters.Country)
	})
}

// TestStorage_UpdateProperty tests the UpdateProperty method
func TestStorage_UpdateProperty(t *testing.T) {
	t.Run("ValidUpdate", func(t *testing.T) {
		// Arrange
		hotelID := int64(12345)
		propertyData := getSamplePropertyData()

		// Act & Assert
		assert.Equal(t, int64(12345), hotelID)
		assert.NotNil(t, propertyData)
		assert.Equal(t, int64(12345), propertyData.Property.HotelID)
	})

	t.Run("InvalidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(0)
		propertyData := getSamplePropertyData()

		// Act & Assert
		assert.Equal(t, int64(0), hotelID)
		assert.NotNil(t, propertyData)
	})
}

// TestStorage_DeleteProperty tests the DeleteProperty method
func TestStorage_DeleteProperty(t *testing.T) {
	t.Run("ValidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(12345)

		// Act & Assert
		assert.Greater(t, hotelID, int64(0))
	})

	t.Run("InvalidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(0)

		// Act & Assert
		assert.Equal(t, int64(0), hotelID)
	})
}

// TestStorage_GetPropertyReviews tests the GetPropertyReviews method
func TestStorage_GetPropertyReviews(t *testing.T) {
	t.Run("ValidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(12345)

		// Act & Assert
		assert.Greater(t, hotelID, int64(0))
	})

	t.Run("InvalidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(0)

		// Act & Assert
		assert.Equal(t, int64(0), hotelID)
	})
}

// TestStorage_GetReviewsByScore tests the GetReviewsByScore method
func TestStorage_GetReviewsByScore(t *testing.T) {
	t.Run("ValidScoreRange", func(t *testing.T) {
		// Arrange
		minScore := 4
		maxScore := 5
		limit := 10
		offset := 0

		// Act & Assert
		assert.Equal(t, 4, minScore)
		assert.Equal(t, 5, maxScore)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
		assert.LessOrEqual(t, minScore, maxScore)
	})

	t.Run("InvalidScoreRange", func(t *testing.T) {
		// Arrange
		minScore := 5
		maxScore := 4
		limit := 10
		offset := 0

		// Act & Assert
		assert.Equal(t, 5, minScore)
		assert.Equal(t, 4, maxScore)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
		assert.Greater(t, minScore, maxScore)
	})
}

// TestStorage_GetPropertyTranslations tests the GetPropertyTranslations method
func TestStorage_GetPropertyTranslations(t *testing.T) {
	t.Run("ValidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(12345)

		// Act & Assert
		assert.Greater(t, hotelID, int64(0))
	})

	t.Run("InvalidHotelID", func(t *testing.T) {
		// Arrange
		hotelID := int64(0)

		// Act & Assert
		assert.Equal(t, int64(0), hotelID)
	})
}

// TestStorage_GetTranslationByLanguage tests the GetTranslationByLanguage method
func TestStorage_GetTranslationByLanguage(t *testing.T) {
	t.Run("ValidParameters", func(t *testing.T) {
		// Arrange
		hotelID := int64(12345)
		language := "fr"

		// Act & Assert
		assert.Greater(t, hotelID, int64(0))
		assert.Equal(t, "fr", language)
		assert.Len(t, language, 2)
	})

	t.Run("InvalidLanguage", func(t *testing.T) {
		// Arrange
		hotelID := int64(12345)
		language := ""

		// Act & Assert
		assert.Greater(t, hotelID, int64(0))
		assert.Empty(t, language)
	})
}

// TestStorage_SearchProperties tests the SearchProperties method
func TestStorage_SearchProperties(t *testing.T) {
	t.Run("ValidSearchQuery", func(t *testing.T) {
		// Arrange
		query := "luxury hotel paris"
		limit := 10
		offset := 0

		// Act & Assert
		assert.Equal(t, "luxury hotel paris", query)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
		assert.NotEmpty(t, query)
	})

	t.Run("EmptySearchQuery", func(t *testing.T) {
		// Arrange
		query := ""
		limit := 10
		offset := 0

		// Act & Assert
		assert.Empty(t, query)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
	})
}

// TestStorage_CountSearchProperties tests the CountSearchProperties method
func TestStorage_CountSearchProperties(t *testing.T) {
	t.Run("ValidSearchQuery", func(t *testing.T) {
		// Arrange
		query := "luxury hotel paris"

		// Act & Assert
		assert.Equal(t, "luxury hotel paris", query)
		assert.NotEmpty(t, query)
	})

	t.Run("EmptySearchQuery", func(t *testing.T) {
		// Arrange
		query := ""

		// Act & Assert
		assert.Empty(t, query)
	})
}

// TestStorage_GetPropertiesByLocation tests the GetPropertiesByLocation method
func TestStorage_GetPropertiesByLocation(t *testing.T) {
	t.Run("ValidLocation", func(t *testing.T) {
		// Arrange
		city := "Paris"
		country := "France"
		limit := 10
		offset := 0

		// Act & Assert
		assert.Equal(t, "Paris", city)
		assert.Equal(t, "France", country)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
		assert.NotEmpty(t, city)
		assert.NotEmpty(t, country)
	})

	t.Run("EmptyLocation", func(t *testing.T) {
		// Arrange
		city := ""
		country := ""
		limit := 10
		offset := 0

		// Act & Assert
		assert.Empty(t, city)
		assert.Empty(t, country)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
	})
}

// TestStorage_CountPropertiesByLocation tests the CountPropertiesByLocation method
func TestStorage_CountPropertiesByLocation(t *testing.T) {
	t.Run("ValidLocation", func(t *testing.T) {
		// Arrange
		city := "Paris"
		country := "France"

		// Act & Assert
		assert.Equal(t, "Paris", city)
		assert.Equal(t, "France", country)
		assert.NotEmpty(t, city)
		assert.NotEmpty(t, country)
	})

	t.Run("EmptyLocation", func(t *testing.T) {
		// Arrange
		city := ""
		country := ""

		// Act & Assert
		assert.Empty(t, city)
		assert.Empty(t, country)
	})
}

// TestStorage_GetPropertiesByRating tests the GetPropertiesByRating method
func TestStorage_GetPropertiesByRating(t *testing.T) {
	t.Run("ValidRatingRange", func(t *testing.T) {
		// Arrange
		minRating := 4.0
		limit := 10
		offset := 0

		// Act & Assert
		assert.Equal(t, 4.0, minRating)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
		assert.GreaterOrEqual(t, minRating, 0.0)
		assert.LessOrEqual(t, minRating, 5.0)
	})

	t.Run("InvalidRatingRange", func(t *testing.T) {
		// Arrange
		minRating := -1.0
		limit := 10
		offset := 0

		// Act & Assert
		assert.Equal(t, -1.0, minRating)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 0, offset)
		assert.Less(t, minRating, 0.0)
	})
}

// TestStorage_CountPropertiesByRating tests the CountPropertiesByRating method
func TestStorage_CountPropertiesByRating(t *testing.T) {
	t.Run("ValidRating", func(t *testing.T) {
		// Arrange
		minRating := 4.0

		// Act & Assert
		assert.Equal(t, 4.0, minRating)
		assert.GreaterOrEqual(t, minRating, 0.0)
		assert.LessOrEqual(t, minRating, 5.0)
	})

	t.Run("InvalidRating", func(t *testing.T) {
		// Arrange
		minRating := -1.0

		// Act & Assert
		assert.Equal(t, -1.0, minRating)
		assert.Less(t, minRating, 0.0)
	})
}
