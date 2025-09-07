package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/barimehdi77/cupid-api/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage implements the store.Storage interface for testing
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) StoreProperty(ctx context.Context, propertyData *cupid.PropertyData) error {
	args := m.Called(ctx, propertyData)
	return args.Error(0)
}

func (m *MockStorage) GetProperty(ctx context.Context, hotelID int64) (*cupid.PropertyData, error) {
	args := m.Called(ctx, hotelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cupid.PropertyData), args.Error(1)
}

func (m *MockStorage) ListProperties(ctx context.Context, limit, offset int, filters store.PropertyFilters) ([]*cupid.Property, error) {
	args := m.Called(ctx, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*cupid.Property), args.Error(1)
}

func (m *MockStorage) CountProperties(ctx context.Context, filters store.PropertyFilters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) UpdateProperty(ctx context.Context, hotelID int64, propertyData *cupid.PropertyData) error {
	args := m.Called(ctx, hotelID, propertyData)
	return args.Error(0)
}

func (m *MockStorage) DeleteProperty(ctx context.Context, hotelID int64) error {
	args := m.Called(ctx, hotelID)
	return args.Error(0)
}

func (m *MockStorage) GetPropertyReviews(ctx context.Context, hotelID int64) ([]cupid.Review, error) {
	args := m.Called(ctx, hotelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cupid.Review), args.Error(1)
}

func (m *MockStorage) GetReviewsByScore(ctx context.Context, minScore, maxScore int, limit, offset int) ([]cupid.Review, error) {
	args := m.Called(ctx, minScore, maxScore, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cupid.Review), args.Error(1)
}

func (m *MockStorage) GetPropertyTranslations(ctx context.Context, hotelID int64) (map[string]*cupid.Property, error) {
	args := m.Called(ctx, hotelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*cupid.Property), args.Error(1)
}

func (m *MockStorage) GetTranslationByLanguage(ctx context.Context, hotelID int64, language string) (*cupid.Property, error) {
	args := m.Called(ctx, hotelID, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cupid.Property), args.Error(1)
}

func (m *MockStorage) SearchProperties(ctx context.Context, query string, limit, offset int) ([]*cupid.Property, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*cupid.Property), args.Error(1)
}

func (m *MockStorage) CountSearchProperties(ctx context.Context, query string) (int, error) {
	args := m.Called(ctx, query)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) GetPropertiesByLocation(ctx context.Context, city, country string, limit, offset int) ([]*cupid.Property, error) {
	args := m.Called(ctx, city, country, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*cupid.Property), args.Error(1)
}

func (m *MockStorage) CountPropertiesByLocation(ctx context.Context, city, country string) (int, error) {
	args := m.Called(ctx, city, country)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) GetPropertiesByRating(ctx context.Context, minRating float64, limit, offset int) ([]*cupid.Property, error) {
	args := m.Called(ctx, minRating, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*cupid.Property), args.Error(1)
}

func (m *MockStorage) CountPropertiesByRating(ctx context.Context, minRating float64) (int, error) {
	args := m.Called(ctx, minRating)
	return args.Int(0), args.Error(1)
}

// Test data fixtures
func createTestProperty() *cupid.Property {
	return &cupid.Property{
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
	}
}

func createTestPropertyData() *cupid.PropertyData {
	return &cupid.PropertyData{
		Property: *createTestProperty(),
		Reviews: []cupid.Review{
			{
				ReviewID:     1,
				AverageScore: 9,
				Country:      "GB",
				Name:         "John Doe",
				Headline:     "Great hotel!",
				Pros:         "Clean, comfortable",
				Cons:         "No complaints",
				Date:         "2024-01-15",
				Language:     "en",
			},
		},
		Translations: map[string]*cupid.Property{
			"fr": {
				HotelID:   12345,
				HotelName: "Hôtel de Test",
				Address: cupid.Address{
					City:    "Londres",
					Country: "gb",
				},
			},
		},
	}
}

func setupTestRouter(handlers *Handlers) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Initialize logger for testing
	logger.InitLogger()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthCheckHandler)
		v1.GET("/properties", handlers.ListPropertiesHandler)
		v1.GET("/properties/:id", handlers.GetPropertyHandler)
		v1.GET("/properties/:id/reviews", handlers.GetPropertyReviewsHandler)
		v1.GET("/properties/:id/translations", handlers.GetPropertyTranslationsHandler)
		v1.GET("/properties/location", handlers.GetPropertiesByLocationHandler)
		v1.GET("/properties/rating", handlers.GetPropertiesByRatingHandler)
		v1.GET("/search", handlers.SearchPropertiesHandler)
	}

	return router
}

// Test HealthCheckHandler
func TestHealthCheckHandler(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	// Verify health response structure
	healthData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "healthy", healthData["status"])
	assert.Equal(t, "1.0.0", healthData["version"])
	assert.Equal(t, "connected", healthData["database"])
}

// Test ListPropertiesHandler - Success Case
func TestListPropertiesHandler_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testProperties := []*cupid.Property{createTestProperty()}
	testFilters := store.PropertyFilters{}

	mockStorage.On("ListProperties", mock.Anything, 20, 0, testFilters).Return(testProperties, nil)
	mockStorage.On("CountProperties", mock.Anything, testFilters).Return(1, nil)

	req, _ := http.NewRequest("GET", "/api/v1/properties?limit=20&page=1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Meta)

	// Verify pagination metadata
	assert.Equal(t, 1, response.Meta.Page)
	assert.Equal(t, 20, response.Meta.Limit)
	assert.Equal(t, 1, response.Meta.Total)
	assert.Equal(t, 1, response.Meta.TotalItems)
	assert.Equal(t, 1, response.Meta.TotalPages)
	assert.False(t, response.Meta.HasNext)
	assert.False(t, response.Meta.HasPrev)

	// Verify property data
	properties, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, properties, 1)

	mockStorage.AssertExpectations(t)
}

// Test ListPropertiesHandler - Database Error
func TestListPropertiesHandler_DatabaseError(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testFilters := store.PropertyFilters{}

	mockStorage.On("ListProperties", mock.Anything, 20, 0, testFilters).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/properties", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Failed to fetch properties", response.Error)

	mockStorage.AssertExpectations(t)
}

// Test ListPropertiesHandler - Invalid Query Parameters
func TestListPropertiesHandler_InvalidQueryParams(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	req, _ := http.NewRequest("GET", "/api/v1/properties?limit=invalid", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Invalid query parameters")
}

// Test GetPropertyHandler - Success Case
func TestGetPropertyHandler_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testPropertyData := createTestPropertyData()

	mockStorage.On("GetProperty", mock.Anything, int64(12345)).Return(testPropertyData, nil)

	req, _ := http.NewRequest("GET", "/api/v1/properties/12345", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	// Verify property with details structure
	propertyData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, propertyData["property"])
	assert.NotNil(t, propertyData["reviews"])
	assert.NotNil(t, propertyData["translations"])

	mockStorage.AssertExpectations(t)
}

// Test GetPropertyHandler - Property Not Found
func TestGetPropertyHandler_NotFound(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	mockStorage.On("GetProperty", mock.Anything, int64(99999)).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/properties/99999", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Failed to fetch property", response.Error)

	mockStorage.AssertExpectations(t)
}

// Test GetPropertyHandler - Invalid Property ID
func TestGetPropertyHandler_InvalidID(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	req, _ := http.NewRequest("GET", "/api/v1/properties/invalid", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid property ID", response.Error)
}

// Test SearchPropertiesHandler - Success Case
func TestSearchPropertiesHandler_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testProperties := []*cupid.Property{createTestProperty()}
	searchQuery := "London"

	mockStorage.On("SearchProperties", mock.Anything, searchQuery, 20, 0).Return(testProperties, nil)
	mockStorage.On("CountSearchProperties", mock.Anything, searchQuery).Return(1, nil)

	req, _ := http.NewRequest("GET", "/api/v1/search?q=London&limit=20&page=1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Meta)

	// Verify search results
	properties, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, properties, 1)

	mockStorage.AssertExpectations(t)
}

// Test SearchPropertiesHandler - Missing Query Parameter
func TestSearchPropertiesHandler_MissingQuery(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	req, _ := http.NewRequest("GET", "/api/v1/search", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Invalid query parameters")
}

// Test GetPropertiesByRatingHandler - Success Case
func TestGetPropertiesByRatingHandler_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testProperties := []*cupid.Property{createTestProperty()}
	minRating := 9.0

	mockStorage.On("GetPropertiesByRating", mock.Anything, minRating, 20, 0).Return(testProperties, nil)
	mockStorage.On("CountPropertiesByRating", mock.Anything, minRating).Return(1, nil)

	req, _ := http.NewRequest("GET", "/api/v1/properties/rating?min_rating=9.0&limit=20&page=1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Meta)

	mockStorage.AssertExpectations(t)
}

// Test GetPropertiesByRatingHandler - Missing Rating Parameter
func TestGetPropertiesByRatingHandler_MissingRating(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	req, _ := http.NewRequest("GET", "/api/v1/properties/rating", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "min_rating parameter is required", response.Error)
}

// Test GetPropertiesByRatingHandler - Invalid Rating Parameter
func TestGetPropertiesByRatingHandler_InvalidRating(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	req, _ := http.NewRequest("GET", "/api/v1/properties/rating?min_rating=invalid", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid min_rating parameter", response.Error)
}

// Test GetPropertiesByLocationHandler - Success Case
func TestGetPropertiesByLocationHandler_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testProperties := []*cupid.Property{createTestProperty()}
	city := "London"
	country := "gb"

	mockStorage.On("GetPropertiesByLocation", mock.Anything, city, country, 20, 0).Return(testProperties, nil)
	mockStorage.On("CountPropertiesByLocation", mock.Anything, city, country).Return(1, nil)

	req, _ := http.NewRequest("GET", "/api/v1/properties/location?city=London&country=gb&limit=20&page=1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Meta)

	mockStorage.AssertExpectations(t)
}

// Test GetPropertyReviewsHandler - Success Case
func TestGetPropertyReviewsHandler_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testReviews := []cupid.Review{
		{
			ReviewID:     1,
			AverageScore: 9,
			Country:      "GB",
			Name:         "John Doe",
			Headline:     "Great hotel!",
			Pros:         "Clean, comfortable",
			Cons:         "No complaints",
			Date:         "2024-01-15",
			Language:     "en",
		},
	}

	mockStorage.On("GetPropertyReviews", mock.Anything, int64(12345)).Return(testReviews, nil)

	req, _ := http.NewRequest("GET", "/api/v1/properties/12345/reviews", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	// Verify reviews data
	reviews, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, reviews, 1)

	mockStorage.AssertExpectations(t)
}

// Test GetPropertyTranslationsHandler - Success Case
func TestGetPropertyTranslationsHandler_Success(t *testing.T) {
	// Arrange
	mockStorage := new(MockStorage)
	handlers := NewHandlers(mockStorage)
	router := setupTestRouter(handlers)

	testTranslations := map[string]*cupid.Property{
		"fr": {
			HotelID:   12345,
			HotelName: "Hôtel de Test",
			Address: cupid.Address{
				City:    "Londres",
				Country: "gb",
			},
		},
	}

	mockStorage.On("GetPropertyTranslations", mock.Anything, int64(12345)).Return(testTranslations, nil)

	req, _ := http.NewRequest("GET", "/api/v1/properties/12345/translations", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	// Verify translations data
	translations, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, translations, "fr")

	mockStorage.AssertExpectations(t)
}
