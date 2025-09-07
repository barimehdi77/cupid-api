package store

import (
	"context"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/database"
)

// Storage interface defines all storage operations
type Storage interface {
	// Property operations
	StoreProperty(ctx context.Context, propertyData *cupid.PropertyData) error
	GetProperty(ctx context.Context, hotelID int64) (*cupid.PropertyData, error)
	ListProperties(ctx context.Context, limit, offset int, filters PropertyFilters) ([]*cupid.Property, error)
	CountProperties(ctx context.Context, filters PropertyFilters) (int, error)
	UpdateProperty(ctx context.Context, hotelID int64, propertyData *cupid.PropertyData) error
	DeleteProperty(ctx context.Context, hotelID int64) error

	// Review operations
	GetPropertyReviews(ctx context.Context, hotelID int64) ([]cupid.Review, error)
	GetReviewsByScore(ctx context.Context, minScore, maxScore int, limit, offset int) ([]cupid.Review, error)

	// Translation operations
	GetPropertyTranslations(ctx context.Context, hotelID int64) (map[string]*cupid.Property, error)
	GetTranslationByLanguage(ctx context.Context, hotelID int64, language string) (*cupid.Property, error)

	// Search operations
	SearchProperties(ctx context.Context, query string, limit, offset int) ([]*cupid.Property, error)
	CountSearchProperties(ctx context.Context, query string) (int, error)
	GetPropertiesByLocation(ctx context.Context, city, country string, limit, offset int) ([]*cupid.Property, error)
	CountPropertiesByLocation(ctx context.Context, city, country string) (int, error)
	GetPropertiesByRating(ctx context.Context, minRating float64, limit, offset int) ([]*cupid.Property, error)
	CountPropertiesByRating(ctx context.Context, minRating float64) (int, error)
}

// PropertyFilters contains filtering options for property queries
type PropertyFilters struct {
	City      string
	Country   string
	MinStars  int
	MaxStars  int
	MinRating float64
	MaxRating float64
	HotelType string
	Chain     string
}

// storage implements the Storage interface
type storage struct {
	db *database.DB
}

// NewStorage creates a new storage instance
func NewStorage(db *database.DB) Storage {
	return &storage{db: db}
}
