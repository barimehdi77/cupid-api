package sync

import (
	"context"
	"testing"
	"time"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

// MockCupidService is a mock implementation of the Cupid service
type MockCupidService struct {
	mock.Mock
}

func (m *MockCupidService) FetchAllProperties(ctx context.Context) ([]*cupid.PropertyData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*cupid.PropertyData), args.Error(1)
}

func (m *MockCupidService) FetchProperty(ctx context.Context, propertyID int64) (*cupid.PropertyData, error) {
	args := m.Called(ctx, propertyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cupid.PropertyData), args.Error(1)
}

// MockStorage is a mock implementation of the Storage interface
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

// TestConfig tests the configuration structure
func TestConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		// Act
		config := DefaultConfig()

		// Assert
		assert.NotNil(t, config)
		assert.Equal(t, 12*time.Hour, config.Interval)
		assert.Equal(t, 10, config.BatchSize)
		assert.Equal(t, 5, config.MaxConcurrent)
		assert.Equal(t, 3, config.RetryAttempts)
		assert.Equal(t, 5*time.Second, config.RetryDelay)
		assert.Equal(t, 10, config.RateLimitPerSec)
		assert.True(t, config.EnableAuto)
	})

	t.Run("CustomConfig", func(t *testing.T) {
		// Arrange
		config := &Config{
			Interval:        6 * time.Hour,
			BatchSize:       5,
			MaxConcurrent:   3,
			RetryAttempts:   2,
			RetryDelay:      2 * time.Second,
			RateLimitPerSec: 5,
			EnableAuto:      false,
		}

		// Act & Assert
		assert.NotNil(t, config)
		assert.Equal(t, 6*time.Hour, config.Interval)
		assert.Equal(t, 5, config.BatchSize)
		assert.Equal(t, 3, config.MaxConcurrent)
		assert.Equal(t, 2, config.RetryAttempts)
		assert.Equal(t, 2*time.Second, config.RetryDelay)
		assert.Equal(t, 5, config.RateLimitPerSec)
		assert.False(t, config.EnableAuto)
	})
}

// TestSyncStats tests the SyncStats structure
func TestSyncStats(t *testing.T) {
	t.Run("InitialStats", func(t *testing.T) {
		// Act
		stats := &SyncStats{}

		// Assert
		assert.NotNil(t, stats)
		assert.Equal(t, 0, stats.TotalProperties)
		assert.Equal(t, 0, stats.UpdatedProperties)
		assert.Equal(t, 0, stats.FailedProperties)
		assert.Equal(t, time.Time{}, stats.LastSync)
		assert.Nil(t, stats.LastError)
	})

	t.Run("StatsWithData", func(t *testing.T) {
		// Arrange
		stats := &SyncStats{
			TotalProperties:   100,
			UpdatedProperties: 50,
			FailedProperties:  5,
			LastSync:          time.Now(),
			LastError:         nil,
		}

		// Act & Assert
		assert.NotNil(t, stats)
		assert.Equal(t, 100, stats.TotalProperties)
		assert.Equal(t, 50, stats.UpdatedProperties)
		assert.Equal(t, 5, stats.FailedProperties)
		assert.False(t, stats.LastSync.IsZero())
		assert.Nil(t, stats.LastError)
	})
}

// TestMockCupidService tests the mock Cupid service
func TestMockCupidService(t *testing.T) {
	t.Run("FetchAllProperties", func(t *testing.T) {
		// Arrange
		mockService := &MockCupidService{}
		expectedProperties := []*cupid.PropertyData{getSamplePropertyData()}
		mockService.On("FetchAllProperties", mock.Anything).Return(expectedProperties, nil)

		// Act
		properties, err := mockService.FetchAllProperties(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.Len(t, properties, 1)
		assert.Equal(t, expectedProperties[0], properties[0])
		mockService.AssertExpectations(t)
	})

	t.Run("FetchProperty", func(t *testing.T) {
		// Arrange
		mockService := &MockCupidService{}
		expectedProperty := getSamplePropertyData()
		propertyID := int64(12345)
		mockService.On("FetchProperty", mock.Anything, propertyID).Return(expectedProperty, nil)

		// Act
		property, err := mockService.FetchProperty(context.Background(), propertyID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperty, property)
		mockService.AssertExpectations(t)
	})
}

// TestMockStorage tests the mock storage
func TestMockStorage(t *testing.T) {
	t.Run("StoreProperty", func(t *testing.T) {
		// Arrange
		mockStorage := &MockStorage{}
		propertyData := getSamplePropertyData()
		mockStorage.On("StoreProperty", mock.Anything, propertyData).Return(nil)

		// Act
		err := mockStorage.StoreProperty(context.Background(), propertyData)

		// Assert
		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("GetProperty", func(t *testing.T) {
		// Arrange
		mockStorage := &MockStorage{}
		expectedProperty := getSamplePropertyData()
		hotelID := int64(12345)
		mockStorage.On("GetProperty", mock.Anything, hotelID).Return(expectedProperty, nil)

		// Act
		property, err := mockStorage.GetProperty(context.Background(), hotelID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperty, property)
		mockStorage.AssertExpectations(t)
	})

	t.Run("CountProperties", func(t *testing.T) {
		// Arrange
		mockStorage := &MockStorage{}
		filters := store.PropertyFilters{City: "Paris"}
		expectedCount := 10
		mockStorage.On("CountProperties", mock.Anything, filters).Return(expectedCount, nil)

		// Act
		count, err := mockStorage.CountProperties(context.Background(), filters)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		mockStorage.AssertExpectations(t)
	})
}

// TestDataComparator tests the data comparator
func TestDataComparator(t *testing.T) {
	t.Run("NewDataComparator", func(t *testing.T) {
		// Act
		comparator := NewDataComparator()

		// Assert
		assert.NotNil(t, comparator)
	})

	t.Run("ComparePropertyData_SameData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData()
		property2 := getSamplePropertyData()

		// Act
		changes := comparator.ComparePropertyData(property1, property2)

		// Assert
		assert.NotNil(t, changes)
		assert.False(t, changes.HasChanges())
	})

	t.Run("ComparePropertyData_DifferentData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData()
		property2 := getSamplePropertyData()
		property2.Property.HotelName = "Different Hotel"

		// Act
		changes := comparator.ComparePropertyData(property1, property2)

		// Assert
		assert.NotNil(t, changes)
		assert.True(t, changes.HasChanges())
	})
}

// TestScheduler tests the scheduler
func TestScheduler(t *testing.T) {
	t.Run("NewScheduler", func(t *testing.T) {
		// Arrange
		interval := 12 * time.Hour
		syncFunc := func(ctx context.Context) (*SyncResult, error) {
			return &SyncResult{}, nil
		}

		// Act
		scheduler := NewScheduler(interval, syncFunc)

		// Assert
		assert.NotNil(t, scheduler)
		assert.False(t, scheduler.isRunning)
		assert.Equal(t, interval, scheduler.interval)
	})

	t.Run("SchedulerWithCustomInterval", func(t *testing.T) {
		// Arrange
		interval := 6 * time.Hour
		syncFunc := func(ctx context.Context) (*SyncResult, error) {
			return &SyncResult{}, nil
		}

		// Act
		scheduler := NewScheduler(interval, syncFunc)

		// Assert
		assert.NotNil(t, scheduler)
		assert.Equal(t, interval, scheduler.interval)
	})
}
