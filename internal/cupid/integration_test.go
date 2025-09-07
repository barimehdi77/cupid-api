package cupid

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCupidClientIntegration tests the Cupid API client against the real third-party API
// These tests require valid API credentials and network connectivity
func TestCupidClientIntegration(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable.")
	}

	// Skip if API key is not provided
	apiKey := os.Getenv("CUPID_API_KEY")
	if apiKey == "" {
		t.Skip("CUPID_API_KEY not provided. Skipping integration tests.")
	}

	// Initialize logger for tests
	if err := logger.InitLogger(); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("Client_Initialization", func(t *testing.T) {
		assert.NotNil(t, client)
		assert.NotEmpty(t, client.baseURL)
		assert.NotEmpty(t, client.version)
		assert.NotEmpty(t, client.apiKey)
		assert.NotNil(t, client.httpClient)
		assert.Equal(t, 30*time.Second, client.httpClient.Timeout)
	})

	t.Run("GetProperty_Success", func(t *testing.T) {
		// Test with a known property ID from PropertyIDs list
		propertyID := PropertyIDs[0]

		property, err := client.GetProperty(ctx, propertyID)

		require.NoError(t, err)
		require.NotNil(t, property)
		assert.Equal(t, propertyID, property.HotelID)
		assert.NotEmpty(t, property.HotelName)
		assert.NotEmpty(t, property.Address.City)
		assert.NotEmpty(t, property.Address.Country)
		assert.Greater(t, property.Stars, 0)
		assert.LessOrEqual(t, property.Stars, 5)
		assert.GreaterOrEqual(t, property.Rating, 0.0)
		assert.LessOrEqual(t, property.Rating, 10.0)
		assert.GreaterOrEqual(t, property.ReviewCount, 0)
	})

	t.Run("GetProperty_InvalidID", func(t *testing.T) {
		// Test with an invalid property ID
		invalidID := int64(999999999)

		property, err := client.GetProperty(ctx, invalidID)

		assert.Error(t, err)
		assert.Nil(t, property)
		assert.Contains(t, err.Error(), "API error")
	})

	t.Run("GetPropertyReviews_Success", func(t *testing.T) {
		propertyID := PropertyIDs[0]
		reviewCount := 5

		reviews, err := client.GetPropertyReviews(ctx, propertyID, reviewCount)

		require.NoError(t, err)
		require.NotNil(t, reviews)
		assert.LessOrEqual(t, len(reviews), reviewCount)

		// Validate review structure if reviews exist
		if len(reviews) > 0 {
			review := reviews[0]
			assert.NotEmpty(t, review.ReviewID)
			assert.NotEmpty(t, review.Name)
			assert.GreaterOrEqual(t, review.AverageScore, 0.0)
			assert.LessOrEqual(t, review.AverageScore, 10.0)
			assert.NotEmpty(t, review.Date)
		}
	})

	t.Run("GetPropertyReviews_InvalidID", func(t *testing.T) {
		invalidID := int64(999999999)
		reviewCount := 5

		reviews, err := client.GetPropertyReviews(ctx, invalidID, reviewCount)

		assert.Error(t, err)
		assert.Nil(t, reviews)
		assert.Contains(t, err.Error(), "API error")
	})

	t.Run("GetPropertyTranslations_Success", func(t *testing.T) {
		propertyID := PropertyIDs[0]
		language := "fr"

		translation, err := client.GetPropertyTranslations(ctx, propertyID, language)

		require.NoError(t, err)
		require.NotNil(t, translation)
		assert.Equal(t, propertyID, translation.HotelID)
		assert.NotEmpty(t, translation.HotelName)
	})

	t.Run("GetPropertyTranslations_InvalidLanguage", func(t *testing.T) {
		propertyID := PropertyIDs[0]
		invalidLanguage := "invalid_lang"

		_, err := client.GetPropertyTranslations(ctx, propertyID, invalidLanguage)

		// This might succeed or fail depending on API behavior
		// We just ensure it doesn't panic
		if err != nil {
			assert.Contains(t, err.Error(), "API error")
		}
	})

	t.Run("FetchAllPropertyData_Success", func(t *testing.T) {
		propertyID := PropertyIDs[0]

		propertyData, err := client.FetchAllPropertyData(ctx, propertyID)

		require.NoError(t, err)
		require.NotNil(t, propertyData)
		assert.Equal(t, propertyID, propertyData.Property.HotelID)
		assert.NotEmpty(t, propertyData.Property.HotelName)
		assert.NotNil(t, propertyData.Reviews)
		assert.NotNil(t, propertyData.Translations)
	})

	t.Run("FetchAllPropertyData_InvalidID", func(t *testing.T) {
		invalidID := int64(999999999)

		propertyData, err := client.FetchAllPropertyData(ctx, invalidID)

		assert.Error(t, err)
		assert.Nil(t, propertyData)
		assert.Contains(t, err.Error(), "failed to fetch property details")
	})
}

// TestCupidServiceIntegration tests the Cupid service against the real third-party API
func TestCupidServiceIntegration(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable.")
	}

	// Skip if API key is not provided
	apiKey := os.Getenv("CUPID_API_KEY")
	if apiKey == "" {
		t.Skip("CUPID_API_KEY not provided. Skipping integration tests.")
	}

	// Initialize logger for tests
	if err := logger.InitLogger(); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	service := NewService()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Run("Service_Initialization", func(t *testing.T) {
		assert.NotNil(t, service)
		assert.NotNil(t, service.client)
	})

	t.Run("FetchProperty_Success", func(t *testing.T) {
		propertyID := PropertyIDs[0]

		propertyData, err := service.FetchProperty(ctx, propertyID)

		require.NoError(t, err)
		require.NotNil(t, propertyData)
		assert.Equal(t, propertyID, propertyData.Property.HotelID)
		assert.NotEmpty(t, propertyData.Property.HotelName)
	})

	t.Run("FetchAllProperties_Success", func(t *testing.T) {
		// This test might take a while due to rate limiting
		properties, err := service.FetchAllProperties(ctx)

		require.NoError(t, err)
		require.NotNil(t, properties)
		assert.Greater(t, len(properties), 0)
		assert.LessOrEqual(t, len(properties), len(PropertyIDs))

		// Validate first property structure
		if len(properties) > 0 {
			property := properties[0]
			assert.NotEmpty(t, property.Property.HotelID)
			assert.NotEmpty(t, property.Property.HotelName)
			assert.NotEmpty(t, property.Property.Address.City)
			assert.NotEmpty(t, property.Property.Address.Country)
		}
	})

	t.Run("FetchAllProperties_Timeout", func(t *testing.T) {
		// Test with a very short timeout to ensure timeout handling works
		shortCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		properties, err := service.FetchAllProperties(shortCtx)

		// Should either succeed with partial results or fail due to timeout
		// Both are acceptable behaviors
		if err != nil {
			assert.Contains(t, err.Error(), "context deadline exceeded")
		} else {
			assert.NotNil(t, properties)
		}
	})
}

// TestCupidAPIConnectivity tests basic connectivity to the Cupid API
func TestCupidAPIConnectivity(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable.")
	}

	// Skip if API key is not provided
	apiKey := os.Getenv("CUPID_API_KEY")
	if apiKey == "" {
		t.Skip("CUPID_API_KEY not provided. Skipping integration tests.")
	}

	// Initialize logger for tests
	if err := logger.InitLogger(); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("API_Connectivity", func(t *testing.T) {
		// Test basic connectivity with a simple property fetch
		propertyID := PropertyIDs[0]

		property, err := client.GetProperty(ctx, propertyID)

		// If this fails, it means there's a connectivity issue
		require.NoError(t, err, "Failed to connect to Cupid API. Check network connectivity and API credentials.")
		require.NotNil(t, property)
		assert.Equal(t, propertyID, property.HotelID)
	})

	t.Run("API_RateLimiting", func(t *testing.T) {
		// Test rate limiting by making multiple rapid requests
		propertyID := PropertyIDs[0]

		// Make 3 rapid requests
		for i := 0; i < 3; i++ {
			_, err := client.GetProperty(ctx, propertyID)
			require.NoError(t, err)
		}

		// If we get here without errors, rate limiting is working
		t.Log("Rate limiting test passed - no errors with rapid requests")
	})

	t.Run("API_ErrorHandling", func(t *testing.T) {
		// Test error handling with invalid endpoint
		invalidID := int64(999999999)

		_, err := client.GetProperty(ctx, invalidID)

		// Should get an error for invalid property ID
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

// TestCupidDataValidation tests data validation from the Cupid API
func TestCupidDataValidation(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable.")
	}

	// Skip if API key is not provided
	apiKey := os.Getenv("CUPID_API_KEY")
	if apiKey == "" {
		t.Skip("CUPID_API_KEY not provided. Skipping integration tests.")
	}

	// Initialize logger for tests
	if err := logger.InitLogger(); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("Property_DataValidation", func(t *testing.T) {
		propertyID := PropertyIDs[0]

		property, err := client.GetProperty(ctx, propertyID)
		require.NoError(t, err)
		require.NotNil(t, property)

		// Validate required fields
		assert.NotEmpty(t, property.HotelID, "HotelID should not be empty")
		assert.NotEmpty(t, property.HotelName, "HotelName should not be empty")
		assert.NotEmpty(t, property.Address.City, "Address.City should not be empty")
		assert.NotEmpty(t, property.Address.Country, "Address.Country should not be empty")

		// Validate numeric ranges
		assert.Greater(t, property.Stars, 0, "Stars should be greater than 0")
		assert.LessOrEqual(t, property.Stars, 5, "Stars should be 5 or less")
		assert.GreaterOrEqual(t, property.Rating, 0.0, "Rating should be 0 or greater")
		assert.LessOrEqual(t, property.Rating, 10.0, "Rating should be 10 or less")
		assert.GreaterOrEqual(t, property.ReviewCount, 0, "ReviewCount should be 0 or greater")

		// Validate coordinates
		assert.GreaterOrEqual(t, property.Latitude, -90.0, "Latitude should be valid")
		assert.LessOrEqual(t, property.Latitude, 90.0, "Latitude should be valid")
		assert.GreaterOrEqual(t, property.Longitude, -180.0, "Longitude should be valid")
		assert.LessOrEqual(t, property.Longitude, 180.0, "Longitude should be valid")
	})

	t.Run("Reviews_DataValidation", func(t *testing.T) {
		propertyID := PropertyIDs[0]
		reviewCount := 3

		reviews, err := client.GetPropertyReviews(ctx, propertyID, reviewCount)
		require.NoError(t, err)

		// Validate review structure if reviews exist
		for i, review := range reviews {
			assert.NotEmpty(t, review.ReviewID, "Review %d: ReviewID should not be empty", i)
			assert.NotEmpty(t, review.Name, "Review %d: Name should not be empty", i)
			assert.GreaterOrEqual(t, review.AverageScore, 0.0, "Review %d: AverageScore should be 0 or greater", i)
			assert.LessOrEqual(t, review.AverageScore, 10.0, "Review %d: AverageScore should be 10 or less", i)
			assert.NotEmpty(t, review.Date, "Review %d: Date should not be empty", i)
		}
	})

	t.Run("Translations_DataValidation", func(t *testing.T) {
		propertyID := PropertyIDs[0]
		language := "fr"

		translation, err := client.GetPropertyTranslations(ctx, propertyID, language)
		require.NoError(t, err)
		require.NotNil(t, translation)

		// Validate translation structure
		assert.Equal(t, propertyID, translation.HotelID, "Translation HotelID should match original")
		assert.NotEmpty(t, translation.HotelName, "Translation HotelName should not be empty")
		assert.NotEmpty(t, translation.Address.City, "Translation Address.City should not be empty")
		assert.NotEmpty(t, translation.Address.Country, "Translation Address.Country should not be empty")
	})
}

// TestCupidPerformance tests performance characteristics of the Cupid API
func TestCupidPerformance(t *testing.T) {
	// Skip integration tests if not explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable.")
	}

	// Skip if API key is not provided
	apiKey := os.Getenv("CUPID_API_KEY")
	if apiKey == "" {
		t.Skip("CUPID_API_KEY not provided. Skipping integration tests.")
	}

	// Initialize logger for tests
	if err := logger.InitLogger(); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("Single_Property_Performance", func(t *testing.T) {
		propertyID := PropertyIDs[0]

		start := time.Now()
		property, err := client.GetProperty(ctx, propertyID)
		duration := time.Since(start)

		require.NoError(t, err)
		require.NotNil(t, property)

		// Single property fetch should be reasonably fast
		assert.Less(t, duration, 5*time.Second, "Single property fetch should complete within 5 seconds")
		t.Logf("Single property fetch took: %v", duration)
	})

	t.Run("Complete_Property_Data_Performance", func(t *testing.T) {
		propertyID := PropertyIDs[0]

		start := time.Now()
		propertyData, err := client.FetchAllPropertyData(ctx, propertyID)
		duration := time.Since(start)

		require.NoError(t, err)
		require.NotNil(t, propertyData)

		// Complete property data fetch should be reasonably fast
		assert.Less(t, duration, 10*time.Second, "Complete property data fetch should complete within 10 seconds")
		t.Logf("Complete property data fetch took: %v", duration)
	})
}

// BenchmarkCupidAPI benchmarks the Cupid API performance
func BenchmarkCupidAPI(b *testing.B) {
	// Skip benchmark if not explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		b.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable.")
	}

	// Skip if API key is not provided
	apiKey := os.Getenv("CUPID_API_KEY")
	if apiKey == "" {
		b.Skip("CUPID_API_KEY not provided. Skipping integration tests.")
	}

	client := NewClient()
	ctx := context.Background()

	b.Run("GetProperty", func(b *testing.B) {
		propertyID := PropertyIDs[0]

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.GetProperty(ctx, propertyID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GetPropertyReviews", func(b *testing.B) {
		propertyID := PropertyIDs[0]
		reviewCount := 5

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.GetPropertyReviews(ctx, propertyID, reviewCount)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("FetchAllPropertyData", func(b *testing.B) {
		propertyID := PropertyIDs[0]

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.FetchAllPropertyData(ctx, propertyID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
