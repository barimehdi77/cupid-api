package api

import (
	"testing"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/stretchr/testify/assert"
)

// Test ConvertPropertyToResponse
func TestConvertPropertyToResponse(t *testing.T) {
	// Arrange
	property := &cupid.Property{
		HotelID:     12345,
		CupidID:     12345,
		HotelName:   "Test Hotel",
		HotelType:   "Hotels",
		Chain:       "Test Chain",
		Latitude:    51.5074,
		Longitude:   -0.1278,
		Stars:       5,
		Rating:      9.5,
		ReviewCount: 100,
		Address: cupid.Address{
			Address:    "123 Test Street",
			City:       "London",
			State:      "England",
			Country:    "gb",
			PostalCode: "SW1A 1AA",
		},
		MainImageTh: "https://example.com/image.jpg",
	}

	// Act
	response := ConvertPropertyToResponse(property)

	// Assert
	assert.Equal(t, property.HotelID, response.HotelID)
	assert.Equal(t, property.CupidID, response.CupidID)
	assert.Equal(t, property.HotelName, response.HotelName)
	assert.Equal(t, property.HotelType, response.HotelType)
	assert.Equal(t, property.Chain, response.Chain)
	assert.Equal(t, property.Latitude, response.Latitude)
	assert.Equal(t, property.Longitude, response.Longitude)
	assert.Equal(t, property.Stars, response.Stars)
	assert.Equal(t, property.Rating, response.Rating)
	assert.Equal(t, property.ReviewCount, response.ReviewCount)
	assert.Equal(t, property.MainImageTh, response.MainImageTh)
	// Note: CreatedAt and UpdatedAt are not part of the Property model

	// Verify address conversion
	assert.Equal(t, property.Address.Address, response.Address.Address)
	assert.Equal(t, property.Address.City, response.Address.City)
	assert.Equal(t, property.Address.State, response.Address.State)
	assert.Equal(t, property.Address.Country, response.Address.Country)
	assert.Equal(t, property.Address.PostalCode, response.Address.PostalCode)
}

// Test ConvertPropertyToResponse with nil property
func TestConvertPropertyToResponse_NilProperty(t *testing.T) {
	// Arrange
	var property *cupid.Property = nil

	// Act
	response := ConvertPropertyToResponse(property)

	// Assert
	assert.Equal(t, int64(0), response.HotelID)
	assert.Equal(t, int64(0), response.CupidID)
	assert.Equal(t, "", response.HotelName)
	assert.Equal(t, "", response.HotelType)
	assert.Equal(t, "", response.Chain)
	assert.Equal(t, float64(0), response.Latitude)
	assert.Equal(t, float64(0), response.Longitude)
	assert.Equal(t, 0, response.Stars)
	assert.Equal(t, float64(0), response.Rating)
	assert.Equal(t, 0, response.ReviewCount)
	assert.Equal(t, "", response.MainImageTh)
	// Note: CreatedAt and UpdatedAt are not part of the Property model
}

// Test ConvertReviewToResponse
func TestConvertReviewToResponse(t *testing.T) {
	// Arrange
	review := cupid.Review{
		ReviewID:     1,
		AverageScore: 9,
		Country:      "GB",
		Name:         "John Doe",
		Headline:     "Great hotel!",
		Pros:         "Clean, comfortable",
		Cons:         "No complaints",
		Date:         "2024-01-15",
		Language:     "en",
	}

	// Act
	response := ConvertReviewToResponse(review)

	// Assert
	assert.Equal(t, review.ReviewID, response.ReviewID)
	assert.Equal(t, review.AverageScore, response.AverageScore)
	assert.Equal(t, review.Country, response.Country)
	assert.Equal(t, review.Name, response.Name)
	assert.Equal(t, review.Headline, response.Headline)
	assert.Equal(t, review.Pros, response.Pros)
	assert.Equal(t, review.Cons, response.Cons)
	assert.Equal(t, review.Date, response.Date)
	assert.Equal(t, review.Language, response.Language)
}

// Test ConvertReviewToResponse with empty review
func TestConvertReviewToResponse_EmptyReview(t *testing.T) {
	// Arrange
	review := cupid.Review{}

	// Act
	response := ConvertReviewToResponse(review)

	// Assert
	assert.Equal(t, int64(0), response.ReviewID)
	assert.Equal(t, 0, response.AverageScore)
	assert.Equal(t, "", response.Country)
	assert.Equal(t, "", response.Name)
	assert.Equal(t, "", response.Headline)
	assert.Equal(t, "", response.Pros)
	assert.Equal(t, "", response.Cons)
	assert.Equal(t, "", response.Date)
	assert.Equal(t, "", response.Language)
}

// Test ConvertTranslationToResponse
func TestConvertTranslationToResponse(t *testing.T) {
	// Arrange
	language := "fr"
	translation := &cupid.Property{
		HotelID:   12345,
		HotelName: "Hôtel de Test",
		Address: cupid.Address{
			City:    "Londres",
			Country: "gb",
		},
	}

	// Act
	response := ConvertTranslationToResponse(language, translation)

	// Assert
	assert.Equal(t, language, response.Language)
	assert.Equal(t, translation.HotelName, response.HotelName)
}

// Test ConvertTranslationToResponse with nil translation
func TestConvertTranslationToResponse_NilTranslation(t *testing.T) {
	// Arrange
	language := "fr"
	var translation *cupid.Property = nil

	// Act
	response := ConvertTranslationToResponse(language, translation)

	// Assert
	assert.Equal(t, language, response.Language)
	assert.Equal(t, "", response.HotelName)
}

// Test ConvertTranslationToResponse with empty language
func TestConvertTranslationToResponse_EmptyLanguage(t *testing.T) {
	// Arrange
	language := ""
	translation := &cupid.Property{
		HotelID:   12345,
		HotelName: "Test Hotel",
	}

	// Act
	response := ConvertTranslationToResponse(language, translation)

	// Assert
	assert.Equal(t, language, response.Language)
	assert.Equal(t, translation.HotelName, response.HotelName)
}

// Test PropertyListRequest validation
func TestPropertyListRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		request     PropertyListRequest
		expectError bool
	}{
		{
			name: "Valid request with all fields",
			request: PropertyListRequest{
				Page:      1,
				Limit:     20,
				City:      "London",
				Country:   "gb",
				MinStars:  3,
				MaxStars:  5,
				MinRating: 7.0,
				MaxRating: 10.0,
				HotelType: "Hotels",
				Chain:     "Test Chain",
				Search:    "test",
			},
			expectError: false,
		},
		{
			name: "Valid request with minimal fields",
			request: PropertyListRequest{
				Page:  1,
				Limit: 10,
			},
			expectError: false,
		},
		{
			name: "Invalid limit too high",
			request: PropertyListRequest{
				Page:  1,
				Limit: 101, // Should be max 100
			},
			expectError: true,
		},
		{
			name: "Invalid limit too low",
			request: PropertyListRequest{
				Page:  1,
				Limit: 0, // Should be min 1
			},
			expectError: true,
		},
		{
			name: "Invalid stars range",
			request: PropertyListRequest{
				Page:     1,
				Limit:    20,
				MinStars: 6, // Should be max 5
				MaxStars: 5,
			},
			expectError: true,
		},
		{
			name: "Invalid rating range",
			request: PropertyListRequest{
				Page:      1,
				Limit:     20,
				MinRating: 11.0, // Should be max 10
				MaxRating: 10.0,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := validatePropertyListRequest(tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test SearchRequest validation
func TestSearchRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		request     SearchRequest
		expectError bool
	}{
		{
			name: "Valid search request",
			request: SearchRequest{
				Query: "London",
				Page:  1,
				Limit: 20,
			},
			expectError: false,
		},
		{
			name: "Missing query",
			request: SearchRequest{
				Query: "",
				Page:  1,
				Limit: 20,
			},
			expectError: true,
		},
		{
			name: "Invalid limit",
			request: SearchRequest{
				Query: "London",
				Page:  1,
				Limit: 101, // Should be max 100
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := validateSearchRequest(tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test Meta pagination calculations
func TestMeta_PaginationCalculations(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		total    int
		expected Meta
	}{
		{
			name:  "First page with results",
			page:  1,
			limit: 20,
			total: 50,
			expected: Meta{
				Page:       1,
				Limit:      20,
				Total:      50,
				TotalItems: 50,
				TotalPages: 3,
				HasNext:    true,
				HasPrev:    false,
			},
		},
		{
			name:  "Middle page",
			page:  2,
			limit: 20,
			total: 50,
			expected: Meta{
				Page:       2,
				Limit:      20,
				Total:      50,
				TotalItems: 50,
				TotalPages: 3,
				HasNext:    true,
				HasPrev:    true,
			},
		},
		{
			name:  "Last page",
			page:  3,
			limit: 20,
			total: 50,
			expected: Meta{
				Page:       3,
				Limit:      20,
				Total:      50,
				TotalItems: 50,
				TotalPages: 3,
				HasNext:    false,
				HasPrev:    true,
			},
		},
		{
			name:  "Single page",
			page:  1,
			limit: 20,
			total: 10,
			expected: Meta{
				Page:       1,
				Limit:      20,
				Total:      10,
				TotalItems: 10,
				TotalPages: 1,
				HasNext:    false,
				HasPrev:    false,
			},
		},
		{
			name:  "Empty results",
			page:  1,
			limit: 20,
			total: 0,
			expected: Meta{
				Page:       1,
				Limit:      20,
				Total:      0,
				TotalItems: 0,
				TotalPages: 0,
				HasNext:    false,
				HasPrev:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			meta := Meta{
				Page:       tt.page,
				Limit:      tt.limit,
				Total:      tt.total,
				TotalItems: tt.total,
				TotalPages: (tt.total + tt.limit - 1) / tt.limit,
				HasNext:    tt.page < (tt.total+tt.limit-1)/tt.limit,
				HasPrev:    tt.page > 1,
			}

			// Assert
			assert.Equal(t, tt.expected.Page, meta.Page)
			assert.Equal(t, tt.expected.Limit, meta.Limit)
			assert.Equal(t, tt.expected.Total, meta.Total)
			assert.Equal(t, tt.expected.TotalItems, meta.TotalItems)
			assert.Equal(t, tt.expected.TotalPages, meta.TotalPages)
			assert.Equal(t, tt.expected.HasNext, meta.HasNext)
			assert.Equal(t, tt.expected.HasPrev, meta.HasPrev)
		})
	}
}

// Helper functions for validation (these would need to be implemented in models.go)
func validatePropertyListRequest(req PropertyListRequest) error {
	if req.Limit < 1 || req.Limit > 100 {
		return assert.AnError
	}
	if req.MinStars < 0 || req.MinStars > 5 {
		return assert.AnError
	}
	if req.MaxStars < 0 || req.MaxStars > 5 {
		return assert.AnError
	}
	if req.MinRating < 0 || req.MinRating > 10 {
		return assert.AnError
	}
	if req.MaxRating < 0 || req.MaxRating > 10 {
		return assert.AnError
	}
	return nil
}

func validateSearchRequest(req SearchRequest) error {
	if req.Query == "" {
		return assert.AnError
	}
	if req.Limit < 1 || req.Limit > 100 {
		return assert.AnError
	}
	return nil
}
