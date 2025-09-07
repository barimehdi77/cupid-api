package api

import (
	"time"

	"github.com/barimehdi77/cupid-api/internal/cupid"
)

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination and metadata information
type Meta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// PropertyListRequest represents query parameters for listing properties
type PropertyListRequest struct {
	Page      int     `form:"page"`
	Limit     int     `form:"limit"`
	City      string  `form:"city"`
	Country   string  `form:"country"`
	MinStars  int     `form:"min_stars"`
	MaxStars  int     `form:"max_stars"`
	MinRating float64 `form:"min_rating"`
	MaxRating float64 `form:"max_rating"`
	HotelType string  `form:"hotel_type"`
	Chain     string  `form:"chain"`
	Search    string  `form:"search"`
}

// PropertyResponse represents a property in API responses
type PropertyResponse struct {
	HotelID     int64                    `json:"hotel_id"`
	CupidID     int64                    `json:"cupid_id"`
	HotelName   string                   `json:"hotel_name"`
	HotelType   string                   `json:"hotel_type"`
	Chain       string                   `json:"chain"`
	Latitude    float64                  `json:"latitude"`
	Longitude   float64                  `json:"longitude"`
	Stars       int                      `json:"stars"`
	Rating      float64                  `json:"rating"`
	ReviewCount int                      `json:"review_count"`
	AirportCode string                   `json:"airport_code"`
	Address     AddressResponse          `json:"address"`
	MainImageTh string                   `json:"main_image_th"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Details     *PropertyDetailsResponse `json:"details,omitempty"`
}

// AddressResponse represents address information in API responses
type AddressResponse struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// PropertyDetailsResponse represents complex property details
type PropertyDetailsResponse struct {
	Address     interface{} `json:"address,omitempty"`
	CheckIn     interface{} `json:"checkin,omitempty"`
	Facilities  interface{} `json:"facilities,omitempty"`
	Policies    interface{} `json:"policies,omitempty"`
	Rooms       interface{} `json:"rooms,omitempty"`
	Photos      interface{} `json:"photos,omitempty"`
	ContactInfo interface{} `json:"contact_info,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

// ReviewResponse represents a review in API responses
type ReviewResponse struct {
	ID           int64     `json:"id"`
	ReviewID     int64     `json:"review_id"`
	AverageScore int       `json:"average_score"`
	Country      string    `json:"country"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Date         string    `json:"date"`
	Headline     string    `json:"headline"`
	Language     string    `json:"language"`
	Pros         string    `json:"pros"`
	Cons         string    `json:"cons"`
	Source       string    `json:"source"`
	CreatedAt    time.Time `json:"created_at"`
}

// TranslationResponse represents a translation in API responses
type TranslationResponse struct {
	Language            string    `json:"language"`
	HotelName           string    `json:"hotel_name"`
	Description         string    `json:"description"`
	MarkdownDescription string    `json:"markdown_description"`
	ImportantInfo       string    `json:"important_info"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// PropertyWithDetailsResponse represents a complete property with all details
type PropertyWithDetailsResponse struct {
	Property     PropertyResponse               `json:"property"`
	Reviews      []ReviewResponse               `json:"reviews"`
	Translations map[string]TranslationResponse `json:"translations"`
}

// SearchRequest represents search query parameters
type SearchRequest struct {
	Query string `form:"q" binding:"required"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}

// ReviewListRequest represents query parameters for listing reviews
type ReviewListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	Limit    int    `form:"limit" binding:"min=1,max=100"`
	MinScore int    `form:"min_score" binding:"min=1,max=10"`
	MaxScore int    `form:"max_score" binding:"min=1,max=10"`
	Country  string `form:"country"`
	Language string `form:"language"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
}

// ConvertPropertyToResponse converts a cupid.Property to PropertyResponse
func ConvertPropertyToResponse(property *cupid.Property) PropertyResponse {
	if property == nil {
		return PropertyResponse{}
	}

	return PropertyResponse{
		HotelID:     property.HotelID,
		CupidID:     property.CupidID,
		HotelName:   property.HotelName,
		HotelType:   property.HotelType,
		Chain:       property.Chain,
		Latitude:    property.Latitude,
		Longitude:   property.Longitude,
		Stars:       property.Stars,
		Rating:      property.Rating,
		ReviewCount: property.ReviewCount,
		AirportCode: property.AirportCode,
		Address: AddressResponse{
			Address:    property.Address.Address,
			City:       property.Address.City,
			State:      property.Address.State,
			Country:    property.Address.Country,
			PostalCode: property.Address.PostalCode,
		},
		MainImageTh: property.MainImageTh,
	}
}

// ConvertReviewToResponse converts a cupid.Review to ReviewResponse
func ConvertReviewToResponse(review cupid.Review) ReviewResponse {
	return ReviewResponse{
		ReviewID:     review.ReviewID,
		AverageScore: review.AverageScore,
		Country:      review.Country,
		Type:         review.Type,
		Name:         review.Name,
		Date:         review.Date,
		Headline:     review.Headline,
		Language:     review.Language,
		Pros:         review.Pros,
		Cons:         review.Cons,
		Source:       review.Source,
	}
}

// ConvertTranslationToResponse converts a cupid.Property to TranslationResponse
func ConvertTranslationToResponse(language string, translation *cupid.Property) TranslationResponse {
	if translation == nil {
		return TranslationResponse{
			Language: language,
		}
	}

	return TranslationResponse{
		Language:            language,
		HotelName:           translation.HotelName,
		Description:         translation.Description,
		MarkdownDescription: translation.MarkdownDescription,
		ImportantInfo:       translation.ImportantInfo,
	}
}
