package sync

import (
	"reflect"
	"strings"

	"github.com/barimehdi77/cupid-api/internal/cupid"
)

// PropertyChanges represents changes detected in property data
type PropertyChanges struct {
	PropertyChanged     bool
	ReviewsChanged      bool
	TranslationsChanged bool
	Changes             []string
}

// HasChanges returns true if any changes were detected
func (pc *PropertyChanges) HasChanges() bool {
	return pc.PropertyChanged || pc.ReviewsChanged || pc.TranslationsChanged
}

// DataComparator handles comparison of property data
type DataComparator struct{}

// NewDataComparator creates a new data comparator
func NewDataComparator() *DataComparator {
	return &DataComparator{}
}

// ComparePropertyData compares fetched property data with stored data
func (dc *DataComparator) ComparePropertyData(fetched, stored *cupid.PropertyData) *PropertyChanges {
	changes := &PropertyChanges{
		Changes: make([]string, 0),
	}

	// Compare main property data
	if dc.compareProperty(&fetched.Property, &stored.Property) {
		changes.PropertyChanged = true
		changes.Changes = append(changes.Changes, "property")
	}

	// Compare reviews
	if dc.compareReviews(fetched.Reviews, stored.Reviews) {
		changes.ReviewsChanged = true
		changes.Changes = append(changes.Changes, "reviews")
	}

	// Compare translations
	if dc.compareTranslations(fetched.Translations, stored.Translations) {
		changes.TranslationsChanged = true
		changes.Changes = append(changes.Changes, "translations")
	}

	return changes
}

// compareProperty compares two property objects
func (dc *DataComparator) compareProperty(fetched, stored *cupid.Property) bool {
	// Compare basic fields
	if fetched.HotelID != stored.HotelID ||
		fetched.CupidID != stored.CupidID ||
		fetched.HotelName != stored.HotelName ||
		fetched.HotelType != stored.HotelType ||
		fetched.Chain != stored.Chain ||
		fetched.Stars != stored.Stars ||
		fetched.Rating != stored.Rating ||
		fetched.ReviewCount != stored.ReviewCount ||
		fetched.MainImageTh != stored.MainImageTh {
		return true
	}

	// Compare coordinates (with small tolerance for floating point)
	if !dc.compareFloat64(fetched.Latitude, stored.Latitude) ||
		!dc.compareFloat64(fetched.Longitude, stored.Longitude) {
		return true
	}

	// Compare address
	if dc.compareAddress(&fetched.Address, &stored.Address) {
		return true
	}

	return false
}

// compareAddress compares two address objects
func (dc *DataComparator) compareAddress(fetched, stored *cupid.Address) bool {
	return fetched.Address != stored.Address ||
		fetched.City != stored.City ||
		fetched.State != stored.State ||
		fetched.Country != stored.Country ||
		fetched.PostalCode != stored.PostalCode
}

// compareReviews compares two review slices
func (dc *DataComparator) compareReviews(fetched, stored []cupid.Review) bool {
	if len(fetched) != len(stored) {
		return true
	}

	// Create maps for easier comparison
	fetchedMap := make(map[int64]cupid.Review)
	for _, review := range fetched {
		fetchedMap[review.ReviewID] = review
	}

	storedMap := make(map[int64]cupid.Review)
	for _, review := range stored {
		storedMap[review.ReviewID] = review
	}

	// Check if all fetched reviews exist in stored and are the same
	for id, fetchedReview := range fetchedMap {
		storedReview, exists := storedMap[id]
		if !exists || dc.compareReview(&fetchedReview, &storedReview) {
			return true
		}
	}

	// Check if all stored reviews exist in fetched
	for id := range storedMap {
		if _, exists := fetchedMap[id]; !exists {
			return true
		}
	}

	return false
}

// compareReview compares two review objects
func (dc *DataComparator) compareReview(fetched, stored *cupid.Review) bool {
	return fetched.ReviewID != stored.ReviewID ||
		fetched.AverageScore != stored.AverageScore ||
		fetched.Country != stored.Country ||
		fetched.Name != stored.Name ||
		fetched.Headline != stored.Headline ||
		fetched.Pros != stored.Pros ||
		fetched.Cons != stored.Cons ||
		fetched.Date != stored.Date ||
		fetched.Language != stored.Language ||
		fetched.Source != stored.Source
}

// compareTranslations compares two translation maps
func (dc *DataComparator) compareTranslations(fetched, stored map[string]*cupid.Property) bool {
	if len(fetched) != len(stored) {
		return true
	}

	// Compare each language
	for lang, fetchedProp := range fetched {
		storedProp, exists := stored[lang]
		if !exists || dc.compareProperty(fetchedProp, storedProp) {
			return true
		}
	}

	// Check if all stored languages exist in fetched
	for lang := range stored {
		if _, exists := fetched[lang]; !exists {
			return true
		}
	}

	return false
}

// compareFloat64 compares two float64 values with small tolerance
func (dc *DataComparator) compareFloat64(a, b float64) bool {
	const tolerance = 0.0001
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}

// ComparePropertyFields compares specific fields of two properties
func (dc *DataComparator) ComparePropertyFields(fetched, stored *cupid.Property, fields []string) bool {
	for _, field := range fields {
		switch field {
		case "hotel_name":
			if fetched.HotelName != stored.HotelName {
				return true
			}
		case "rating":
			if !dc.compareFloat64(fetched.Rating, stored.Rating) {
				return true
			}
		case "review_count":
			if fetched.ReviewCount != stored.ReviewCount {
				return true
			}
		case "stars":
			if fetched.Stars != stored.Stars {
				return true
			}
		case "address":
			if dc.compareAddress(&fetched.Address, &stored.Address) {
				return true
			}
		case "main_image":
			if fetched.MainImageTh != stored.MainImageTh {
				return true
			}
		}
	}
	return false
}

// GetChangedFields returns a list of changed fields between two properties
func (dc *DataComparator) GetChangedFields(fetched, stored *cupid.Property) []string {
	changedFields := make([]string, 0)

	allFields := []string{
		"hotel_name", "rating", "review_count", "stars", "address", "main_image",
		"hotel_type", "chain", "latitude", "longitude",
	}

	for _, field := range allFields {
		if dc.ComparePropertyFields(fetched, stored, []string{field}) {
			changedFields = append(changedFields, field)
		}
	}

	return changedFields
}

// CompareReviewsByScore compares reviews by score range
func (dc *DataComparator) CompareReviewsByScore(fetched, stored []cupid.Review, minScore, maxScore int) bool {
	fetchedFiltered := dc.filterReviewsByScore(fetched, minScore, maxScore)
	storedFiltered := dc.filterReviewsByScore(stored, minScore, maxScore)

	return dc.compareReviews(fetchedFiltered, storedFiltered)
}

// filterReviewsByScore filters reviews by score range
func (dc *DataComparator) filterReviewsByScore(reviews []cupid.Review, minScore, maxScore int) []cupid.Review {
	filtered := make([]cupid.Review, 0)
	for _, review := range reviews {
		if review.AverageScore >= minScore && review.AverageScore <= maxScore {
			filtered = append(filtered, review)
		}
	}
	return filtered
}

// CompareTranslationsByLanguage compares translations for a specific language
func (dc *DataComparator) CompareTranslationsByLanguage(fetched, stored map[string]*cupid.Property, language string) bool {
	fetchedTrans, fetchedExists := fetched[language]
	storedTrans, storedExists := stored[language]

	if !fetchedExists && !storedExists {
		return false
	}

	if fetchedExists != storedExists {
		return true
	}

	return dc.compareProperty(fetchedTrans, storedTrans)
}

// GetTranslationLanguages returns all languages present in both maps
func (dc *DataComparator) GetTranslationLanguages(fetched, stored map[string]*cupid.Property) []string {
	languages := make(map[string]bool)

	for lang := range fetched {
		languages[lang] = true
	}

	for lang := range stored {
		languages[lang] = true
	}

	result := make([]string, 0, len(languages))
	for lang := range languages {
		result = append(result, lang)
	}

	return result
}

// ComparePropertyDataDeep performs a deep comparison using reflection
func (dc *DataComparator) ComparePropertyDataDeep(fetched, stored *cupid.PropertyData) *PropertyChanges {
	changes := &PropertyChanges{
		Changes: make([]string, 0),
	}

	// Deep compare property
	if !reflect.DeepEqual(fetched.Property, stored.Property) {
		changes.PropertyChanged = true
		changes.Changes = append(changes.Changes, "property")
	}

	// Deep compare reviews
	if !reflect.DeepEqual(fetched.Reviews, stored.Reviews) {
		changes.ReviewsChanged = true
		changes.Changes = append(changes.Changes, "reviews")
	}

	// Deep compare translations
	if !reflect.DeepEqual(fetched.Translations, stored.Translations) {
		changes.TranslationsChanged = true
		changes.Changes = append(changes.Changes, "translations")
	}

	return changes
}

// GetPropertyDataHash returns a hash-like string for quick comparison
func (dc *DataComparator) GetPropertyDataHash(data *cupid.PropertyData) string {
	// Simple hash based on key fields
	hash := strings.Builder{}
	hash.WriteString(data.Property.HotelName)
	hash.WriteString(data.Property.HotelType)
	hash.WriteString(data.Property.Chain)
	hash.WriteString(data.Property.Address.City)
	hash.WriteString(data.Property.Address.Country)
	hash.WriteString(data.Property.MainImageTh)

	// Add review count
	hash.WriteString(string(rune(data.Property.ReviewCount)))

	// Add rating (rounded to 2 decimal places)
	rating := int(data.Property.Rating * 100)
	hash.WriteString(string(rune(rating)))

	return hash.String()
}
