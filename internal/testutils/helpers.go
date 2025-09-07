package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barimehdi77/cupid-api/internal/api"
	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHelper provides common testing utilities
type TestHelper struct {
	td *TestData
}

// NewTestHelper creates a new TestHelper instance
func NewTestHelper() *TestHelper {
	return &TestHelper{
		td: NewTestData(),
	}
}

// SetupTestRouter creates a test router with handlers
func (th *TestHelper) SetupTestRouter(mockStorage store.Storage) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handlers := api.NewHandlers(mockStorage)

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

// CreateTestRequest creates a test HTTP request
func (th *TestHelper) CreateTestRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// CreateTestRequestWithQuery creates a test HTTP request with query parameters
func (th *TestHelper) CreateTestRequestWithQuery(method, url string, queryParams map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

// AssertAPIResponse asserts that the response matches expected values
func (th *TestHelper) AssertAPIResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedSuccess bool) {
	t.Helper()

	assert.Equal(t, expectedStatus, w.Code)

	var response api.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedSuccess, response.Success)
}

// AssertAPIResponseWithData asserts that the response matches expected values and has data
func (th *TestHelper) AssertAPIResponseWithData(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedSuccess bool, expectedData interface{}) {
	t.Helper()

	th.AssertAPIResponse(t, w, expectedStatus, expectedSuccess)

	var response api.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	if expectedData != nil {
		assert.NotNil(t, response.Data)
	} else {
		assert.Nil(t, response.Data)
	}
}

// AssertAPIResponseWithMeta asserts that the response has pagination metadata
func (th *TestHelper) AssertAPIResponseWithMeta(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedSuccess bool, expectedPage, expectedLimit, expectedTotal int) {
	t.Helper()

	th.AssertAPIResponse(t, w, expectedStatus, expectedSuccess)

	var response api.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response.Meta)
	assert.Equal(t, expectedPage, response.Meta.Page)
	assert.Equal(t, expectedLimit, response.Meta.Limit)
	assert.Equal(t, expectedTotal, response.Meta.Total)
	assert.Equal(t, expectedTotal, response.Meta.TotalItems)
}

// AssertErrorResponse asserts that the response is an error response
func (th *TestHelper) AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	t.Helper()

	th.AssertAPIResponse(t, w, expectedStatus, false)

	var response api.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedError, response.Error)
}

// AssertPropertyResponse asserts that the property response matches expected values
func (th *TestHelper) AssertPropertyResponse(t *testing.T, propertyResponse interface{}, expectedProperty *cupid.Property) {
	t.Helper()

	property, ok := propertyResponse.(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, float64(expectedProperty.HotelID), property["hotel_id"])
	assert.Equal(t, expectedProperty.HotelName, property["hotel_name"])
	assert.Equal(t, expectedProperty.HotelType, property["hotel_type"])
	assert.Equal(t, expectedProperty.Chain, property["chain"])
	assert.Equal(t, expectedProperty.Stars, int(property["stars"].(float64)))
	assert.Equal(t, expectedProperty.Rating, property["rating"])
	assert.Equal(t, float64(expectedProperty.ReviewCount), property["review_count"])
}

// AssertReviewResponse asserts that the review response matches expected values
func (th *TestHelper) AssertReviewResponse(t *testing.T, reviewResponse interface{}, expectedReview cupid.Review) {
	t.Helper()

	review, ok := reviewResponse.(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, float64(expectedReview.ReviewID), review["review_id"])
	assert.Equal(t, expectedReview.AverageScore, int(review["average_score"].(float64)))
	assert.Equal(t, expectedReview.Country, review["country"])
	assert.Equal(t, expectedReview.Name, review["name"])
	assert.Equal(t, expectedReview.Headline, review["headline"])
	assert.Equal(t, expectedReview.Pros, review["pros"])
	assert.Equal(t, expectedReview.Cons, review["cons"])
	assert.Equal(t, expectedReview.Date, review["date"])
	assert.Equal(t, expectedReview.Language, review["language"])
}

// AssertTranslationResponse asserts that the translation response matches expected values
func (th *TestHelper) AssertTranslationResponse(t *testing.T, translationResponse interface{}, expectedLanguage string, expectedProperty *cupid.Property) {
	t.Helper()

	translation, ok := translationResponse.(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, expectedLanguage, translation["language"])
	assert.Equal(t, float64(expectedProperty.HotelID), translation["hotel_id"])
	assert.Equal(t, expectedProperty.HotelName, translation["hotel_name"])
}

// AssertPaginationMetadata asserts that the pagination metadata is correct
func (th *TestHelper) AssertPaginationMetadata(t *testing.T, meta *api.Meta, expectedPage, expectedLimit, expectedTotal int) {
	t.Helper()

	assert.Equal(t, expectedPage, meta.Page)
	assert.Equal(t, expectedLimit, meta.Limit)
	assert.Equal(t, expectedTotal, meta.Total)
	assert.Equal(t, expectedTotal, meta.TotalItems)

	expectedTotalPages := (expectedTotal + expectedLimit - 1) / expectedLimit
	assert.Equal(t, expectedTotalPages, meta.TotalPages)
	assert.Equal(t, expectedPage < expectedTotalPages, meta.HasNext)
	assert.Equal(t, expectedPage > 1, meta.HasPrev)
}

// Note: MockStorage implementation would be in the handlers_test.go file
// This file provides general test utilities

// TestCase represents a test case structure
type TestCase struct {
	Name            string
	Setup           func() (store.Storage, *gin.Engine)
	Request         func() (*http.Request, error)
	ExpectedStatus  int
	ExpectedSuccess bool
	ExpectedError   string
	ExpectedData    interface{}
	ExpectedMeta    *PaginationMeta
}

// PaginationMeta represents expected pagination metadata
type PaginationMeta struct {
	Page       int
	Limit      int
	Total      int
	TotalItems int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

// RunTestCase runs a test case
func (th *TestHelper) RunTestCase(t *testing.T, tc TestCase) {
	t.Helper()

	t.Run(tc.Name, func(t *testing.T) {
		// Setup
		storage, router := tc.Setup()

		// Create request
		req, err := tc.Request()
		assert.NoError(t, err)

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert response
		th.AssertAPIResponse(t, w, tc.ExpectedStatus, tc.ExpectedSuccess)

		if tc.ExpectedError != "" {
			th.AssertErrorResponse(t, w, tc.ExpectedStatus, tc.ExpectedError)
		}

		if tc.ExpectedData != nil {
			th.AssertAPIResponseWithData(t, w, tc.ExpectedStatus, tc.ExpectedSuccess, tc.ExpectedData)
		}

		if tc.ExpectedMeta != nil {
			th.AssertAPIResponseWithMeta(t, w, tc.ExpectedStatus, tc.ExpectedSuccess, tc.ExpectedMeta.Page, tc.ExpectedMeta.Limit, tc.ExpectedMeta.Total)
		}

		// Note: Mock assertions would be handled by the specific test
		_ = storage // Avoid unused variable warning
	})
}

// CreateTestContext creates a test context
func (th *TestHelper) CreateTestContext() context.Context {
	return context.Background()
}

// AssertJSONResponse asserts that the response is valid JSON
func (th *TestHelper) AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	var response interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
}

// AssertResponseContains asserts that the response contains expected text
func (th *TestHelper) AssertResponseContains(t *testing.T, w *httptest.ResponseRecorder, expectedText string) {
	t.Helper()

	assert.Contains(t, w.Body.String(), expectedText)
}

// AssertResponseNotContains asserts that the response does not contain unexpected text
func (th *TestHelper) AssertResponseNotContains(t *testing.T, w *httptest.ResponseRecorder, unexpectedText string) {
	t.Helper()

	assert.NotContains(t, w.Body.String(), unexpectedText)
}

// CreateTestServer creates a test server
func (th *TestHelper) CreateTestServer(router *gin.Engine) *httptest.Server {
	return httptest.NewServer(router)
}

// CloseTestServer closes a test server
func (th *TestHelper) CloseTestServer(server *httptest.Server) {
	server.Close()
}

// AssertStatusCode asserts that the response has the expected status code
func (th *TestHelper) AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()

	assert.Equal(t, expectedStatus, w.Code, "Expected status code %d, got %d", expectedStatus, w.Code)
}

// AssertContentType asserts that the response has the expected content type
func (th *TestHelper) AssertContentType(t *testing.T, w *httptest.ResponseRecorder, expectedContentType string) {
	t.Helper()

	assert.Equal(t, expectedContentType, w.Header().Get("Content-Type"))
}

// AssertHeader asserts that the response has the expected header
func (th *TestHelper) AssertHeader(t *testing.T, w *httptest.ResponseRecorder, headerName, expectedValue string) {
	t.Helper()

	assert.Equal(t, expectedValue, w.Header().Get(headerName))
}

// AssertResponseSize asserts that the response has the expected size
func (th *TestHelper) AssertResponseSize(t *testing.T, w *httptest.ResponseRecorder, expectedSize int) {
	t.Helper()

	assert.Equal(t, expectedSize, w.Body.Len())
}

// AssertResponseNotEmpty asserts that the response is not empty
func (th *TestHelper) AssertResponseNotEmpty(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	assert.NotEmpty(t, w.Body.String())
}

// AssertResponseEmpty asserts that the response is empty
func (th *TestHelper) AssertResponseEmpty(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	assert.Empty(t, w.Body.String())
}

// PrintResponse prints the response for debugging
func (th *TestHelper) PrintResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	t.Logf("Response Status: %d", w.Code)
	t.Logf("Response Headers: %v", w.Header())
	t.Logf("Response Body: %s", w.Body.String())
}

// PrintRequest prints the request for debugging
func (th *TestHelper) PrintRequest(t *testing.T, req *http.Request) {
	t.Helper()

	t.Logf("Request Method: %s", req.Method)
	t.Logf("Request URL: %s", req.URL.String())
	t.Logf("Request Headers: %v", req.Header)
}
