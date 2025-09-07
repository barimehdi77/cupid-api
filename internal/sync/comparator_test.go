package sync

import (
	"testing"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/stretchr/testify/assert"
)

// TestNewDataComparator tests the NewDataComparator function
func TestNewDataComparator(t *testing.T) {
	// Act
	comparator := NewDataComparator()

	// Assert
	assert.NotNil(t, comparator)
}

// TestDataComparator_ComparePropertyData tests the ComparePropertyData method
func TestDataComparator_ComparePropertyData(t *testing.T) {
	t.Run("IdenticalData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()

		// Act
		changes := comparator.ComparePropertyData(propertyData1, propertyData2)

		// Assert
		assert.NotNil(t, changes)
		assert.False(t, changes.HasChanges())
		assert.False(t, changes.PropertyChanged)
		assert.False(t, changes.ReviewsChanged)
		assert.False(t, changes.TranslationsChanged)
		assert.Empty(t, changes.Changes)
	})

	t.Run("DifferentPropertyData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()
		propertyData2.Property.HotelName = "Different Hotel Name"

		// Act
		changes := comparator.ComparePropertyData(propertyData1, propertyData2)

		// Assert
		assert.NotNil(t, changes)
		assert.True(t, changes.HasChanges())
		assert.True(t, changes.PropertyChanged)
		assert.False(t, changes.ReviewsChanged)
		assert.False(t, changes.TranslationsChanged)
		assert.Contains(t, changes.Changes, "property")
	})

	t.Run("DifferentReviews", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()
		propertyData2.Reviews = []cupid.Review{
			{
				ReviewID:     123,
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
		}

		// Act
		changes := comparator.ComparePropertyData(propertyData1, propertyData2)

		// Assert
		assert.NotNil(t, changes)
		assert.True(t, changes.HasChanges())
		assert.False(t, changes.PropertyChanged)
		assert.True(t, changes.ReviewsChanged)
		assert.False(t, changes.TranslationsChanged)
		assert.Contains(t, changes.Changes, "reviews")
	})

	t.Run("DifferentTranslations", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()
		propertyData2.Translations = map[string]*cupid.Property{
			"fr": {
				HotelID:   propertyData2.Property.HotelID,
				HotelName: "Hôtel de Luxe",
			},
		}

		// Act
		changes := comparator.ComparePropertyData(propertyData1, propertyData2)

		// Assert
		assert.NotNil(t, changes)
		assert.True(t, changes.HasChanges())
		assert.False(t, changes.PropertyChanged)
		assert.False(t, changes.ReviewsChanged)
		assert.True(t, changes.TranslationsChanged)
		assert.Contains(t, changes.Changes, "translations")
	})

	t.Run("AllDifferent", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()
		propertyData2.Property.HotelName = "Different Hotel"
		propertyData2.Reviews = []cupid.Review{
			{
				ReviewID:     123,
				AverageScore: 4,
				Name:         "John Doe",
			},
		}
		propertyData2.Translations = map[string]*cupid.Property{
			"fr": {
				HotelID:   propertyData2.Property.HotelID,
				HotelName: "Hôtel de Luxe",
			},
		}

		// Act
		changes := comparator.ComparePropertyData(propertyData1, propertyData2)

		// Assert
		assert.NotNil(t, changes)
		assert.True(t, changes.HasChanges())
		assert.True(t, changes.PropertyChanged)
		assert.True(t, changes.ReviewsChanged)
		assert.True(t, changes.TranslationsChanged)
		assert.Contains(t, changes.Changes, "property")
		assert.Contains(t, changes.Changes, "reviews")
		assert.Contains(t, changes.Changes, "translations")
	})
}

// TestDataComparator_ComparePropertyFields tests the ComparePropertyFields method
func TestDataComparator_ComparePropertyFields(t *testing.T) {
	t.Run("SameFields", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property
		fields := []string{"hotel_name", "rating", "stars"}

		// Act
		hasChanges := comparator.ComparePropertyFields(&property1, &property2, fields)

		// Assert
		assert.False(t, hasChanges)
	})

	t.Run("DifferentHotelName", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property
		property2.HotelName = "Different Name"
		fields := []string{"hotel_name"}

		// Act
		hasChanges := comparator.ComparePropertyFields(&property1, &property2, fields)

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("DifferentRating", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property
		property1.Rating = 4.5
		property2.Rating = 4.8
		fields := []string{"rating"}

		// Act
		hasChanges := comparator.ComparePropertyFields(&property1, &property2, fields)

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("DifferentStars", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property
		property1.Stars = 4
		property2.Stars = 5
		fields := []string{"stars"}

		// Act
		hasChanges := comparator.ComparePropertyFields(&property1, &property2, fields)

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("DifferentAddress", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property
		property2.Address.City = "Different City"
		fields := []string{"address"}

		// Act
		hasChanges := comparator.ComparePropertyFields(&property1, &property2, fields)

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("DifferentMainImage", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property
		property2.MainImageTh = "different-image.jpg"
		fields := []string{"main_image"}

		// Act
		hasChanges := comparator.ComparePropertyFields(&property1, &property2, fields)

		// Assert
		assert.True(t, hasChanges)
	})
}

// TestDataComparator_GetChangedFields tests the GetChangedFields method
func TestDataComparator_GetChangedFields(t *testing.T) {
	t.Run("NoChanges", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property

		// Act
		changedFields := comparator.GetChangedFields(&property1, &property2)

		// Assert
		assert.Empty(t, changedFields)
	})

	t.Run("SomeChanges", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		property1 := getSamplePropertyData().Property
		property2 := getSamplePropertyData().Property
		property1.HotelName = "Original Name"
		property1.Rating = 4.5
		property1.Stars = 4
		property2.HotelName = "Different Name"
		property2.Rating = 4.8
		property2.Stars = 5

		// Act
		changedFields := comparator.GetChangedFields(&property1, &property2)

		// Assert
		assert.NotEmpty(t, changedFields)
		assert.Contains(t, changedFields, "hotel_name")
		assert.Contains(t, changedFields, "rating")
		assert.Contains(t, changedFields, "stars")
	})
}

// TestDataComparator_CompareReviewsByScore tests the CompareReviewsByScore method
func TestDataComparator_CompareReviewsByScore(t *testing.T) {
	t.Run("SameReviews", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		reviews1 := []cupid.Review{
			{ReviewID: 1, AverageScore: 4},
			{ReviewID: 2, AverageScore: 3},
		}
		reviews2 := []cupid.Review{
			{ReviewID: 1, AverageScore: 4},
			{ReviewID: 2, AverageScore: 3},
		}

		// Act
		hasChanges := comparator.CompareReviewsByScore(reviews1, reviews2, 1, 5)

		// Assert
		assert.False(t, hasChanges)
	})

	t.Run("DifferentReviews", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		reviews1 := []cupid.Review{
			{ReviewID: 1, AverageScore: 4},
		}
		reviews2 := []cupid.Review{
			{ReviewID: 2, AverageScore: 3},
		}

		// Act
		hasChanges := comparator.CompareReviewsByScore(reviews1, reviews2, 1, 5)

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("DifferentScores", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		reviews1 := []cupid.Review{
			{ReviewID: 1, AverageScore: 4},
		}
		reviews2 := []cupid.Review{
			{ReviewID: 1, AverageScore: 5},
		}

		// Act
		hasChanges := comparator.CompareReviewsByScore(reviews1, reviews2, 1, 5)

		// Assert
		assert.True(t, hasChanges)
	})
}

// TestDataComparator_CompareTranslationsByLanguage tests the CompareTranslationsByLanguage method
func TestDataComparator_CompareTranslationsByLanguage(t *testing.T) {
	t.Run("SameTranslations", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		translations1 := map[string]*cupid.Property{
			"fr": {HotelID: 123, HotelName: "Hôtel de Luxe"},
		}
		translations2 := map[string]*cupid.Property{
			"fr": {HotelID: 123, HotelName: "Hôtel de Luxe"},
		}

		// Act
		hasChanges := comparator.CompareTranslationsByLanguage(translations1, translations2, "fr")

		// Assert
		assert.False(t, hasChanges)
	})

	t.Run("DifferentTranslations", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		translations1 := map[string]*cupid.Property{
			"fr": {HotelID: 123, HotelName: "Hôtel de Luxe"},
		}
		translations2 := map[string]*cupid.Property{
			"fr": {HotelID: 123, HotelName: "Hôtel Magnifique"},
		}

		// Act
		hasChanges := comparator.CompareTranslationsByLanguage(translations1, translations2, "fr")

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("MissingTranslation", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		translations1 := map[string]*cupid.Property{
			"fr": {HotelID: 123, HotelName: "Hôtel de Luxe"},
		}
		translations2 := map[string]*cupid.Property{}

		// Act
		hasChanges := comparator.CompareTranslationsByLanguage(translations1, translations2, "fr")

		// Assert
		assert.True(t, hasChanges)
	})
}

// TestDataComparator_GetTranslationLanguages tests the GetTranslationLanguages method
func TestDataComparator_GetTranslationLanguages(t *testing.T) {
	t.Run("EmptyTranslations", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		translations1 := map[string]*cupid.Property{}
		translations2 := map[string]*cupid.Property{}

		// Act
		languages := comparator.GetTranslationLanguages(translations1, translations2)

		// Assert
		assert.Empty(t, languages)
	})

	t.Run("SameLanguages", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		translations1 := map[string]*cupid.Property{
			"fr": {HotelID: 123},
			"es": {HotelID: 123},
		}
		translations2 := map[string]*cupid.Property{
			"fr": {HotelID: 123},
			"es": {HotelID: 123},
		}

		// Act
		languages := comparator.GetTranslationLanguages(translations1, translations2)

		// Assert
		assert.Len(t, languages, 2)
		assert.Contains(t, languages, "fr")
		assert.Contains(t, languages, "es")
	})

	t.Run("DifferentLanguages", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		translations1 := map[string]*cupid.Property{
			"fr": {HotelID: 123},
		}
		translations2 := map[string]*cupid.Property{
			"es": {HotelID: 123},
		}

		// Act
		languages := comparator.GetTranslationLanguages(translations1, translations2)

		// Assert
		assert.Len(t, languages, 2)
		assert.Contains(t, languages, "fr")
		assert.Contains(t, languages, "es")
	})
}

// TestDataComparator_ComparePropertyDataDeep tests the ComparePropertyDataDeep method
func TestDataComparator_ComparePropertyDataDeep(t *testing.T) {
	t.Run("IdenticalData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()

		// Act
		changes := comparator.ComparePropertyDataDeep(propertyData1, propertyData2)

		// Assert
		assert.NotNil(t, changes)
		assert.False(t, changes.HasChanges())
	})

	t.Run("DifferentData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()
		propertyData2.Property.HotelName = "Different Name"

		// Act
		changes := comparator.ComparePropertyDataDeep(propertyData1, propertyData2)

		// Assert
		assert.NotNil(t, changes)
		assert.True(t, changes.HasChanges())
	})
}

// TestDataComparator_GetPropertyDataHash tests the GetPropertyDataHash method
func TestDataComparator_GetPropertyDataHash(t *testing.T) {
	t.Run("SameData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()

		// Act
		hash1 := comparator.GetPropertyDataHash(propertyData1)
		hash2 := comparator.GetPropertyDataHash(propertyData2)

		// Assert
		assert.Equal(t, hash1, hash2)
		assert.NotEmpty(t, hash1)
	})

	t.Run("DifferentData", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()
		propertyData1 := getSamplePropertyData()
		propertyData2 := getSamplePropertyData()
		propertyData2.Property.HotelName = "Different Name"

		// Act
		hash1 := comparator.GetPropertyDataHash(propertyData1)
		hash2 := comparator.GetPropertyDataHash(propertyData2)

		// Assert
		assert.NotEqual(t, hash1, hash2)
		assert.NotEmpty(t, hash1)
		assert.NotEmpty(t, hash2)
	})
}

// TestDataComparator_compareFloat64 tests the compareFloat64 method
func TestDataComparator_compareFloat64(t *testing.T) {
	t.Run("SameValues", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()

		// Act
		same := comparator.compareFloat64(3.14159, 3.14159)

		// Assert
		assert.True(t, same)
	})

	t.Run("ValuesWithinTolerance", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()

		// Act
		same := comparator.compareFloat64(3.14159, 3.14160)

		// Assert
		assert.True(t, same)
	})

	t.Run("ValuesOutsideTolerance", func(t *testing.T) {
		// Arrange
		comparator := NewDataComparator()

		// Act
		same := comparator.compareFloat64(3.14159, 3.14259)

		// Assert
		assert.False(t, same)
	})
}

// TestPropertyChanges_HasChanges tests the HasChanges method
func TestPropertyChanges_HasChanges(t *testing.T) {
	t.Run("NoChanges", func(t *testing.T) {
		// Arrange
		changes := &PropertyChanges{}

		// Act
		hasChanges := changes.HasChanges()

		// Assert
		assert.False(t, hasChanges)
	})

	t.Run("PropertyChanged", func(t *testing.T) {
		// Arrange
		changes := &PropertyChanges{
			PropertyChanged: true,
		}

		// Act
		hasChanges := changes.HasChanges()

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("ReviewsChanged", func(t *testing.T) {
		// Arrange
		changes := &PropertyChanges{
			ReviewsChanged: true,
		}

		// Act
		hasChanges := changes.HasChanges()

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("TranslationsChanged", func(t *testing.T) {
		// Arrange
		changes := &PropertyChanges{
			TranslationsChanged: true,
		}

		// Act
		hasChanges := changes.HasChanges()

		// Assert
		assert.True(t, hasChanges)
	})

	t.Run("AllChanged", func(t *testing.T) {
		// Arrange
		changes := &PropertyChanges{
			PropertyChanged:     true,
			ReviewsChanged:      true,
			TranslationsChanged: true,
		}

		// Act
		hasChanges := changes.HasChanges()

		// Assert
		assert.True(t, hasChanges)
	})
}
