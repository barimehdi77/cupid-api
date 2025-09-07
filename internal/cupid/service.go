package cupid

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/barimehdi77/cupid-api/internal/logger"
	"go.uber.org/zap"
)

// Service handles batch operations and business logic
type Service struct {
	client *Client
}

// NewService creates a new Cupid service
func NewService() *Service {
	return &Service{
		client: NewClient(),
	}
}

// fetchResult represents the aggregated results from concurrent property fetching operations.
// It contains all successfully fetched properties, any errors that occurred during fetching,
// and the total duration of the operation for performance tracking.
type fetchResult struct {
	// properties contains all successfully fetched property data
	properties []*PropertyData
	// fetchErrors contains all errors that occurred during individual property fetches
	fetchErrors []error
	// duration represents the total time taken for the entire fetch operation
	duration time.Duration
}

// FetchAllProperties fetches all properties from the predefined PropertyIDs list using concurrent processing.
// This is the main entry point for bulk property data retrieval.
//
// The function orchestrates the entire fetching process by:
//  1. Logging the start of the operation
//  2. Processing all properties concurrently with rate limiting
//  3. Collecting and aggregating results
//  4. Logging completion metrics and any errors
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - []*PropertyData: Slice of successfully fetched property data
//   - error: Always returns nil (errors are logged but don't fail the operation)
//
// Note: Individual property fetch failures are logged but don't cause the entire operation to fail.
// This ensures maximum data retrieval even when some properties are unavailable.
func (s *Service) FetchAllProperties(ctx context.Context) ([]*PropertyData, error) {
	s.logFetchStart()

	start := time.Now()
	result := s.processConcurrentFetches(ctx)
	result.duration = time.Since(start)

	s.logFetchResults(result)
	s.logFetchErrors(result.fetchErrors)

	return result.properties, nil
}

// logFetchStart logs the initiation of the property fetching operation.
// This provides visibility into when bulk fetching begins and how many properties
// are being processed, which is useful for monitoring and debugging.
func (s *Service) logFetchStart() {
	logger.LogStartup("Property data fetching",
		zap.Int("total_properties", len(PropertyIDs)),
	)
}

// processConcurrentFetches orchestrates the concurrent fetching of all properties.
// This function sets up the necessary concurrency infrastructure including:
//   - Result and error channels for goroutine communication
//   - WaitGroup for synchronization
//   - Semaphore for rate limiting (max 5 concurrent requests)
//
// The function launches worker goroutines for each property ID and then
// collects all results before returning them in an aggregated format.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - *fetchResult: Aggregated results containing properties, errors, and metadata
func (s *Service) processConcurrentFetches(ctx context.Context) *fetchResult {
	// Channel for results
	results := make(chan *PropertyData, len(PropertyIDs))
	errors := make(chan error, len(PropertyIDs))

	// WaitGroup for concurrency
	var wg sync.WaitGroup

	// Semaphore to limit concurrent requests (avoid rate limiting)
	semaphore := make(chan struct{}, 5) // Max 5 concurrent requests

	// Launch worker goroutines
	s.launchWorkerGoroutines(ctx, &wg, semaphore, results, errors)

	// Close channels when done
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect and return results
	return s.collectFetchResults(results, errors)
}

// launchWorkerGoroutines creates and starts a worker goroutine for each property ID.
// Each goroutine will independently fetch one property's data while respecting
// the concurrency limits imposed by the semaphore.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - wg: WaitGroup to track completion of all workers
//   - semaphore: Channel used as a semaphore to limit concurrent requests
//   - results: Channel for sending successfully fetched property data
//   - errors: Channel for sending any errors that occur during fetching
func (s *Service) launchWorkerGoroutines(ctx context.Context, wg *sync.WaitGroup, semaphore chan struct{}, results chan *PropertyData, errors chan error) {
	for _, propertyID := range PropertyIDs {
		wg.Add(1)
		go s.fetchPropertyWorker(ctx, propertyID, wg, semaphore, results, errors)
	}
}

// fetchPropertyWorker is the worker function that fetches data for a single property.
// This function runs in its own goroutine and handles:
//   - Semaphore acquisition for rate limiting
//   - Rate limiting delay to avoid overwhelming the external API
//   - Actual property data fetching via the client
//   - Error handling and logging
//   - Result communication via channels
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - propertyID: The unique identifier of the property to fetch
//   - wg: WaitGroup to signal completion
//   - semaphore: Channel used as a semaphore to limit concurrent requests
//   - results: Channel for sending successfully fetched property data
//   - errors: Channel for sending any errors that occur during fetching
//
// The function implements a "fail-fast" approach where individual errors don't
// block other workers, ensuring maximum throughput even with partial failures.
func (s *Service) fetchPropertyWorker(ctx context.Context, propertyID int64, wg *sync.WaitGroup, semaphore chan struct{}, results chan *PropertyData, errors chan error) {
	defer wg.Done()

	// Acquire semaphore
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	// Add small delay to avoid rate limiting
	time.Sleep(100 * time.Millisecond)

	propertyData, err := s.client.FetchAllPropertyData(ctx, propertyID)
	if err != nil {
		logger.LogError("Property fetch failed", err,
			zap.Int64("property_id", propertyID),
		)
		errors <- fmt.Errorf("property %d: %w", propertyID, err)
		return
	}

	results <- propertyData
}

// collectFetchResults aggregates all results from the worker goroutines.
// This function reads from both the results and errors channels until both are closed,
// collecting all successful property data and any errors that occurred.
//
// The function uses a select statement to read from both channels concurrently,
// ensuring that neither successful results nor errors block the collection process.
//
// Parameters:
//   - results: Channel containing successfully fetched property data
//   - errors: Channel containing any errors from failed fetch attempts
//
// Returns:
//   - *fetchResult: Aggregated results containing all properties and errors
//
// Note: This function blocks until both channels are closed by the goroutine
// that waits for all workers to complete.
func (s *Service) collectFetchResults(results chan *PropertyData, errors chan error) *fetchResult {
	var properties []*PropertyData
	var fetchErrors []error

	for {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
			} else {
				properties = append(properties, result)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
			} else {
				fetchErrors = append(fetchErrors, err)
			}
		}

		if results == nil && errors == nil {
			break
		}
	}

	return &fetchResult{
		properties:  properties,
		fetchErrors: fetchErrors,
	}
}

// logFetchResults logs comprehensive metrics about the completed fetch operation.
// This provides valuable insights for monitoring, performance analysis, and debugging.
//
// The logged metrics include:
//   - Number of successful property fetches
//   - Number of failed attempts
//   - Total operation duration
//   - Throughput (properties per second)
//
// Parameters:
//   - result: The aggregated results containing metrics to log
//
// This function uses structured logging to ensure metrics can be easily
// parsed and analyzed by monitoring systems.
func (s *Service) logFetchResults(result *fetchResult) {
	logger.LogSuccess("Property data fetching completed",
		zap.Int("successful", len(result.properties)),
		zap.Int("failed", len(result.fetchErrors)),
		zap.Duration("duration", result.duration),
		zap.Float64("properties_per_second", float64(len(result.properties))/result.duration.Seconds()),
	)
}

// logFetchErrors logs detailed information about any errors that occurred during fetching.
// To prevent log spam, this function limits the number of individual errors logged
// while still providing visibility into the overall error rate.
//
// The function:
//   - Logs a summary of the total error count
//   - Logs details for the first 5 errors (configurable via maxErrorsToLog)
//   - Skips logging if no errors occurred
//
// Parameters:
//   - fetchErrors: Slice of all errors that occurred during the fetch operation
//
// This approach balances the need for error visibility with log management,
// preventing excessive log output while ensuring critical error information is captured.
func (s *Service) logFetchErrors(fetchErrors []error) {
	if len(fetchErrors) == 0 {
		return
	}

	logger.Warn("Some properties failed to fetch",
		zap.Int("error_count", len(fetchErrors)),
	)

	// Log first few errors for debugging
	maxErrorsToLog := 5
	for i, err := range fetchErrors {
		if i >= maxErrorsToLog {
			break
		}
		logger.Error("Fetch error", zap.Error(err))
	}
}

// FetchProperty fetches data for a single property by its ID.
// This is a simpler alternative to FetchAllProperties when only one property is needed.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - propertyID: The unique identifier of the property to fetch
//
// Returns:
//   - *PropertyData: The fetched property data, or nil if an error occurred
//   - error: Any error that occurred during the fetch operation
//
// Unlike FetchAllProperties, this function directly returns any errors that occur
// rather than logging them and continuing with partial results.
func (s *Service) FetchProperty(ctx context.Context, propertyID int64) (*PropertyData, error) {
	return s.client.FetchAllPropertyData(ctx, propertyID)
}
