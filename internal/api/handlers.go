package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/barimehdi77/cupid-api/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers contains all API handlers
type Handlers struct {
	storage      store.Storage
	syncHandlers *SyncHandlers
}

// NewHandlers creates a new handlers instance
func NewHandlers(storage store.Storage) *Handlers {
	return &Handlers{storage: storage}
}

// SetSyncHandlers sets the sync handlers
func (h *Handlers) SetSyncHandlers(syncHandlers *SyncHandlers) {
	h.syncHandlers = syncHandlers
}

// HealthCheckHandler handles health check requests
// @Summary Health check
// @Description Check if the API is running and database is connected
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=HealthResponse}
// @Router /health [get]
func (h *Handlers) HealthCheckHandler(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Database:  "connected",
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// ListPropertiesHandler handles listing properties with filtering and pagination
// @Summary List properties
// @Description Get a paginated list of properties with optional filtering
// @Tags properties
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param city query string false "Filter by city"
// @Param country query string false "Filter by country"
// @Param min_stars query int false "Minimum stars" minimum(1) maximum(5)
// @Param max_stars query int false "Maximum stars" minimum(1) maximum(5)
// @Param min_rating query number false "Minimum rating" minimum(0) maximum(10)
// @Param max_rating query number false "Maximum rating" minimum(0) maximum(10)
// @Param hotel_type query string false "Filter by hotel type"
// @Param chain query string false "Filter by chain"
// @Param search query string false "Search in hotel name, city, country"
// @Success 200 {object} APIResponse{data=[]PropertyResponse,meta=Meta}
// @Router /properties [get]
func (h *Handlers) ListPropertiesHandler(c *gin.Context) {
	var req PropertyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	// Convert to storage filters
	filters := store.PropertyFilters{
		City:      req.City,
		Country:   req.Country,
		MinStars:  req.MinStars,
		MaxStars:  req.MaxStars,
		MinRating: req.MinRating,
		MaxRating: req.MaxRating,
		HotelType: req.HotelType,
		Chain:     req.Chain,
	}

	offset := (req.Page - 1) * req.Limit

	var properties []*cupid.Property
	var err error

	if req.Search != "" {
		properties, err = h.storage.SearchProperties(c.Request.Context(), req.Search, req.Limit, offset)
	} else {
		properties, err = h.storage.ListProperties(c.Request.Context(), req.Limit, offset, filters)
	}

	if err != nil {
		logger.LogError("Failed to list properties", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch properties",
		})
		return
	}

	// Get total count for pagination
	totalCount, err := h.storage.CountProperties(c.Request.Context(), filters)
	if err != nil {
		logger.LogError("Failed to count properties", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to count properties",
		})
		return
	}

	// Convert to response format
	var response []PropertyResponse
	for _, property := range properties {
		response = append(response, ConvertPropertyToResponse(property))
	}

	// Calculate pagination metadata
	totalPages := (totalCount + req.Limit - 1) / req.Limit
	meta := &Meta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      totalCount,
		TotalItems: totalCount,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
		Meta:    meta,
	})
}

// GetPropertyHandler handles getting a single property by ID
// @Summary Get property by ID
// @Description Get detailed information about a specific property including reviews and translations
// @Tags properties
// @Accept json
// @Produce json
// @Param id path int true "Property ID"
// @Success 200 {object} APIResponse{data=PropertyWithDetailsResponse}
// @Failure 404 {object} APIResponse
// @Router /properties/{id} [get]
func (h *Handlers) GetPropertyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid property ID",
		})
		return
	}

	propertyData, err := h.storage.GetProperty(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "property not found" {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error:   "Property not found",
			})
			return
		}

		logger.LogError("Failed to get property", err, zap.Int64("property_id", id))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch property",
		})
		return
	}

	// Convert to response format
	propertyResponse := ConvertPropertyToResponse(&propertyData.Property)

	// Convert reviews
	var reviews []ReviewResponse
	for _, review := range propertyData.Reviews {
		reviews = append(reviews, ConvertReviewToResponse(review))
	}

	// Convert translations
	translations := make(map[string]TranslationResponse)
	for lang, translation := range propertyData.Translations {
		translations[lang] = ConvertTranslationToResponse(lang, translation)
	}

	response := PropertyWithDetailsResponse{
		Property:     propertyResponse,
		Reviews:      reviews,
		Translations: translations,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetPropertyReviewsHandler handles getting reviews for a specific property
// @Summary Get property reviews
// @Description Get all reviews for a specific property
// @Tags properties
// @Accept json
// @Produce json
// @Param id path int true "Property ID"
// @Success 200 {object} APIResponse{data=[]ReviewResponse}
// @Failure 404 {object} APIResponse
// @Router /properties/{id}/reviews [get]
func (h *Handlers) GetPropertyReviewsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid property ID",
		})
		return
	}

	reviews, err := h.storage.GetPropertyReviews(c.Request.Context(), id)
	if err != nil {
		logger.LogError("Failed to get property reviews", err, zap.Int64("property_id", id))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch reviews",
		})
		return
	}

	// Convert to response format
	var response []ReviewResponse
	for _, review := range reviews {
		response = append(response, ConvertReviewToResponse(review))
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetPropertyTranslationsHandler handles getting translations for a specific property
// @Summary Get property translations
// @Description Get all translations for a specific property
// @Tags properties
// @Accept json
// @Produce json
// @Param id path int true "Property ID"
// @Success 200 {object} APIResponse{data=map[string]TranslationResponse}
// @Failure 404 {object} APIResponse
// @Router /properties/{id}/translations [get]
func (h *Handlers) GetPropertyTranslationsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid property ID",
		})
		return
	}

	translations, err := h.storage.GetPropertyTranslations(c.Request.Context(), id)
	if err != nil {
		logger.LogError("Failed to get property translations", err, zap.Int64("property_id", id))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch translations",
		})
		return
	}

	// Convert to response format
	response := make(map[string]TranslationResponse)
	for lang, translation := range translations {
		response[lang] = ConvertTranslationToResponse(lang, translation)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// SearchPropertiesHandler handles searching properties
// @Summary Search properties
// @Description Search properties by name, city, or country
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} APIResponse{data=[]PropertyResponse,meta=Meta}
// @Router /search [get]
func (h *Handlers) SearchPropertiesHandler(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	properties, err := h.storage.SearchProperties(c.Request.Context(), req.Query, req.Limit, offset)
	if err != nil {
		logger.LogError("Failed to search properties", err, zap.String("query", req.Query))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to search properties",
		})
		return
	}

	// Get total count for pagination
	totalCount, err := h.storage.CountSearchProperties(c.Request.Context(), req.Query)
	if err != nil {
		logger.LogError("Failed to count search properties", err, zap.String("query", req.Query))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to count search results",
		})
		return
	}

	// Convert to response format
	var response []PropertyResponse
	for _, property := range properties {
		response = append(response, ConvertPropertyToResponse(property))
	}

	// Calculate pagination metadata
	totalPages := (totalCount + req.Limit - 1) / req.Limit
	meta := &Meta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      totalCount,
		TotalItems: totalCount,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
		Meta:    meta,
	})
}

// GetPropertiesByLocationHandler handles getting properties by location
// @Summary Get properties by location
// @Description Get properties filtered by city and/or country
// @Tags properties
// @Accept json
// @Produce json
// @Param city query string false "City name"
// @Param country query string false "Country name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} APIResponse{data=[]PropertyResponse,meta=Meta}
// @Router /properties/location [get]
func (h *Handlers) GetPropertiesByLocationHandler(c *gin.Context) {
	city := c.Query("city")
	country := c.Query("country")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	properties, err := h.storage.GetPropertiesByLocation(c.Request.Context(), city, country, limit, offset)
	if err != nil {
		logger.LogError("Failed to get properties by location", err, zap.String("city", city), zap.String("country", country))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch properties",
		})
		return
	}

	// Get total count for pagination
	totalCount, err := h.storage.CountPropertiesByLocation(c.Request.Context(), city, country)
	if err != nil {
		logger.LogError("Failed to count properties by location", err, zap.String("city", city), zap.String("country", country))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to count properties",
		})
		return
	}

	// Convert to response format
	var response []PropertyResponse
	for _, property := range properties {
		response = append(response, ConvertPropertyToResponse(property))
	}

	// Calculate pagination metadata
	totalPages := (totalCount + limit - 1) / limit
	meta := &Meta{
		Page:       page,
		Limit:      limit,
		Total:      totalCount,
		TotalItems: totalCount,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
		Meta:    meta,
	})
}

// GetPropertiesByRatingHandler handles getting properties by minimum rating
// @Summary Get properties by rating
// @Description Get properties with a minimum rating
// @Tags properties
// @Accept json
// @Produce json
// @Param min_rating query number true "Minimum rating" minimum(0) maximum(10)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} APIResponse{data=[]PropertyResponse,meta=Meta}
// @Router /properties/rating [get]
func (h *Handlers) GetPropertiesByRatingHandler(c *gin.Context) {
	minRatingStr := c.Query("min_rating")
	if minRatingStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "min_rating parameter is required",
		})
		return
	}

	minRating, err := strconv.ParseFloat(minRatingStr, 64)
	if err != nil || minRating < 0 || minRating > 10 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid min_rating parameter",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	properties, err := h.storage.GetPropertiesByRating(c.Request.Context(), minRating, limit, offset)
	if err != nil {
		logger.LogError("Failed to get properties by rating", err, zap.Float64("min_rating", minRating))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to fetch properties",
		})
		return
	}

	// Get total count for pagination
	totalCount, err := h.storage.CountPropertiesByRating(c.Request.Context(), minRating)
	if err != nil {
		logger.LogError("Failed to count properties by rating", err, zap.Float64("min_rating", minRating))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to count properties",
		})
		return
	}

	// Convert to response format
	var response []PropertyResponse
	for _, property := range properties {
		response = append(response, ConvertPropertyToResponse(property))
	}

	// Calculate pagination metadata
	totalPages := (totalCount + limit - 1) / limit
	meta := &Meta{
		Page:       page,
		Limit:      limit,
		Total:      totalCount,
		TotalItems: totalCount,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
		Meta:    meta,
	})
}
