package cupid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/barimehdi77/cupid-api/internal/env"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"go.uber.org/zap"
)

// Client represents the Cupid API client
type Client struct {
	baseURL    string
	version    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Cupid API client
func NewClient() *Client {
	return &Client{
		baseURL: env.GetEnvString("CUPID_API_BASE_URL", "https://api.cupid.com"),
		version: env.GetEnvString("CUPID_API_VERSION", "v1"),
		apiKey:  env.GetEnvString("CUPID_API_KEY", ""),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs HTTP request with retry logic
func (c *Client) doRequest(ctx context.Context, method, endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	logger.Debug("Making API request",
		zap.String("method", method),
		zap.String("url", url),
	)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CupidAPI-Client/1.0")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	logger.Debug("Making API request",
		zap.String("method", method),
		zap.String("url", url),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	return resp, nil
}

// GetProperty fetches a single property by ID
func (c *Client) GetProperty(ctx context.Context, propertyID int64) (*Property, error) {
	endpoint := fmt.Sprintf("/%s/property/%d", c.version, propertyID)

	resp, err := c.doRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch property %d: %w", propertyID, err)
	}
	defer resp.Body.Close()

	var property Property
	if err := json.NewDecoder(resp.Body).Decode(&property); err != nil {
		return nil, fmt.Errorf("failed to decode property response: %w", err)
	}

	logger.Info("Fetched property successfully",
		zap.Int64("property_id", propertyID),
		zap.String("name", property.HotelName),
	)

	return &property, nil
}

// GetPropertyReviews fetches reviews for a property
func (c *Client) GetPropertyReviews(ctx context.Context, propertyID int64, reviewCount int) ([]Review, error) {
	endpoint := fmt.Sprintf("/%s/property/reviews/%d/%d", c.version, propertyID, reviewCount)

	resp, err := c.doRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews for property %d: %w", propertyID, err)
	}
	defer resp.Body.Close()

	var reviews []Review
	if err := json.NewDecoder(resp.Body).Decode(&reviews); err != nil {
		return nil, fmt.Errorf("failed to decode reviews response: %w", err)
	}

	logger.Info("Fetched reviews successfully",
		zap.Int64("property_id", propertyID),
		zap.Int("review_count", len(reviews)),
	)

	return reviews, nil
}

// GetPropertyTranslations fetches translations for a property
func (c *Client) GetPropertyTranslations(ctx context.Context, propertyID int64, language string) (*Property, error) {
	endpoint := fmt.Sprintf("/%s/property/%d/lang/%s", c.version, propertyID, language)

	resp, err := c.doRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch translations for property %d in %s: %w", propertyID, language, err)
	}
	defer resp.Body.Close()

	var translationResponse TranslationResponse
	if err := json.NewDecoder(resp.Body).Decode(&translationResponse); err != nil {
		return nil, fmt.Errorf("failed to decode translation response: %w", err)
	}

	logger.Info("Fetched translation successfully",
		zap.Int64("property_id", propertyID),
		zap.String("language", language),
	)

	return &translationResponse.Data, nil
}

// FetchAllPropertyData fetches complete data for a property (details + reviews + translations)
func (c *Client) FetchAllPropertyData(ctx context.Context, propertyID int64) (*PropertyData, error) {
	logger.LogProgress("Fetching complete property data",
		zap.Int64("property_id", propertyID),
	)

	// Fetch property details
	property, err := c.GetProperty(ctx, propertyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch property details: %w", err)
	}

	// Fetch reviews using the review count from the property
	var reviews []Review
	if property.ReviewCount > 0 {
		reviews, err = c.GetPropertyReviews(ctx, propertyID, property.ReviewCount)
		if err != nil {
			logger.Warn("Failed to fetch reviews, continuing without them",
				zap.Int64("property_id", propertyID),
				zap.Int("review_count", property.ReviewCount),
				zap.Error(err),
			)
			reviews = []Review{} // Continue without reviews
		}
	} else {
		logger.Debug("No reviews available for property",
			zap.Int64("property_id", propertyID),
		)
		reviews = []Review{}
	}

	// Fetch translations (French and Spanish)
	translations := make(map[string]*Property)
	for _, lang := range []string{"fr", "es"} {
		translation, err := c.GetPropertyTranslations(ctx, propertyID, lang)
		if err != nil {
			logger.Warn("Failed to fetch translation, continuing without it",
				zap.Int64("property_id", propertyID),
				zap.String("language", lang),
				zap.Error(err),
			)
			continue
		}
		translations[lang] = translation
	}

	propertyData := &PropertyData{
		Property:     *property,
		Reviews:      reviews,
		Translations: translations,
	}

	logger.LogSuccess("Complete property data fetched",
		zap.Int64("property_id", propertyID),
		zap.Int("review_count", len(reviews)),
		zap.Int("translation_count", len(translations)),
	)

	return propertyData, nil
}
