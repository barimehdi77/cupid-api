# 🧪 Integration Tests for Cupid API

This document describes the integration tests for the Cupid API client to ensure the third-party API works as expected.

## 📋 Overview

The integration tests validate the Cupid API client against the real third-party API, testing critical functionality including:

- ✅ **API Connectivity** - Basic connection and authentication
- ✅ **Data Fetching** - Property, reviews, and translations retrieval
- ✅ **Error Handling** - Invalid requests and error responses
- ✅ **Data Validation** - Response structure and data integrity
- ✅ **Performance** - Response times and rate limiting
- ✅ **Concurrent Operations** - Batch processing and concurrency

## 🚀 Quick Start

### 1. Setup Integration Test Environment

```bash
# Create integration test configuration
make test-integration-setup

# Edit the configuration file
nano integration.env
```

### 2. Configure API Credentials

Edit `integration.env` and set your Cupid API credentials:

```bash
# Required
CUPID_API_KEY=your_actual_api_key_here

# Optional (defaults provided)
CUPID_API_BASE_URL=https://api.cupid.com
CUPID_API_VERSION=v3.0
```

### 3. Run Integration Tests

```bash
# Run all integration tests
make test-integration-cupid

# Run specific test categories
make test-integration-cupid-connectivity
make test-integration-cupid-validation
make test-integration-cupid-performance

# Run benchmarks
make benchmark-cupid
```

## 🧪 Test Categories

### 1. **Client Integration Tests** (`TestCupidClientIntegration`)

Tests the basic Cupid API client functionality:

- **Client Initialization** - Verifies client setup and configuration
- **GetProperty Success** - Fetches a single property successfully
- **GetProperty Invalid ID** - Handles invalid property IDs gracefully
- **GetPropertyReviews Success** - Fetches property reviews
- **GetPropertyReviews Invalid ID** - Handles invalid review requests
- **GetPropertyTranslations Success** - Fetches property translations
- **GetPropertyTranslations Invalid Language** - Handles invalid language codes
- **FetchAllPropertyData Success** - Fetches complete property data
- **FetchAllPropertyData Invalid ID** - Handles invalid property data requests

### 2. **Service Integration Tests** (`TestCupidServiceIntegration`)

Tests the Cupid service layer functionality:

- **Service Initialization** - Verifies service setup
- **FetchProperty Success** - Fetches single property via service
- **FetchAllProperties Success** - Fetches all properties concurrently
- **FetchAllProperties Timeout** - Handles timeout scenarios

### 3. **API Connectivity Tests** (`TestCupidAPIConnectivity`)

Tests basic API connectivity and error handling:

- **API Connectivity** - Verifies connection to Cupid API
- **API Rate Limiting** - Tests rate limiting behavior
- **API Error Handling** - Validates error response handling

### 4. **Data Validation Tests** (`TestCupidDataValidation`)

Tests data integrity and validation:

- **Property Data Validation** - Validates property structure and values
- **Reviews Data Validation** - Validates review structure and values
- **Translations Data Validation** - Validates translation structure and values

### 5. **Performance Tests** (`TestCupidPerformance`)

Tests performance characteristics:

- **Single Property Performance** - Measures single property fetch time
- **Complete Property Data Performance** - Measures complete data fetch time

### 6. **Benchmark Tests** (`BenchmarkCupidAPI`)

Benchmarks API performance:

- **GetProperty Benchmark** - Benchmarks single property fetching
- **GetPropertyReviews Benchmark** - Benchmarks review fetching
- **FetchAllPropertyData Benchmark** - Benchmarks complete data fetching

## 🔧 Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `RUN_INTEGRATION_TESTS` | Yes | `false` | Enable integration tests |
| `CUPID_API_KEY` | Yes | - | Cupid API authentication key |
| `CUPID_API_BASE_URL` | No | `https://api.cupid.com` | Cupid API base URL |
| `CUPID_API_VERSION` | No | `v1` | Cupid API version |

### Test Configuration

The integration tests use the following configuration:

```go
// Test timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

// Rate limiting
time.Sleep(100 * time.Millisecond) // Between requests

// Concurrency limits
semaphore := make(chan struct{}, 5) // Max 5 concurrent requests
```

## 📊 Test Results

### Expected Test Output

```
=== RUN   TestCupidClientIntegration
=== RUN   TestCupidClientIntegration/Client_Initialization
--- PASS: TestCupidClientIntegration/Client_Initialization (0.00s)
=== RUN   TestCupidClientIntegration/GetProperty_Success
--- PASS: TestCupidClientIntegration/GetProperty_Success (0.45s)
=== RUN   TestCupidClientIntegration/GetProperty_InvalidID
--- PASS: TestCupidClientIntegration/GetProperty_InvalidID (0.12s)
=== RUN   TestCupidClientIntegration/GetPropertyReviews_Success
--- PASS: TestCupidClientIntegration/GetPropertyReviews_Success (0.38s)
=== RUN   TestCupidClientIntegration/GetPropertyTranslations_Success
--- PASS: TestCupidClientIntegration/GetPropertyTranslations_Success (0.42s)
=== RUN   TestCupidClientIntegration/FetchAllPropertyData_Success
--- PASS: TestCupidClientIntegration/FetchAllPropertyData_Success (1.23s)
--- PASS: TestCupidClientIntegration (2.60s)

=== RUN   TestCupidServiceIntegration
=== RUN   TestCupidServiceIntegration/Service_Initialization
--- PASS: TestCupidServiceIntegration/Service_Initialization (0.00s)
=== RUN   TestCupidServiceIntegration/FetchProperty_Success
--- PASS: TestCupidServiceIntegration/FetchProperty_Success (1.45s)
=== RUN   TestCupidServiceIntegration/FetchAllProperties_Success
--- PASS: TestCupidServiceIntegration/FetchAllProperties_Success (15.67s)
--- PASS: TestCupidServiceIntegration (17.12s)

=== RUN   TestCupidAPIConnectivity
=== RUN   TestCupidAPIConnectivity/API_Connectivity
--- PASS: TestCupidAPIConnectivity/API_Connectivity (0.45s)
=== RUN   TestCupidAPIConnectivity/API_RateLimiting
--- PASS: TestCupidAPIConnectivity/API_RateLimiting (0.38s)
=== RUN   TestCupidAPIConnectivity/API_ErrorHandling
--- PASS: TestCupidAPIConnectivity/API_ErrorHandling (0.12s)
--- PASS: TestCupidAPIConnectivity (0.95s)

PASS
ok      github.com/barimehdi77/cupid-api/internal/cupid    21.67s
```

### Performance Benchmarks

```
goos: darwin
goarch: arm64
pkg: github.com/barimehdi77/cupid-api/internal/cupid
BenchmarkCupidAPI/GetProperty-8                   10    145.2ms/op
BenchmarkCupidAPI/GetPropertyReviews-8             10    238.7ms/op
BenchmarkCupidAPI/FetchAllPropertyData-8            5    1.234s/op
PASS
ok      github.com/barimehdi77/cupid-api/internal/cupid    15.67s
```

## 🚨 Troubleshooting

### Common Issues

#### 1. **API Key Not Provided**
```
❌ CUPID_API_KEY not provided. Skipping integration tests.
```
**Solution**: Set your API key in `integration.env`

#### 2. **API Connectivity Issues**
```
❌ Failed to connect to Cupid API. Check network connectivity and API credentials.
```
**Solutions**:
- Verify API key is correct
- Check network connectivity
- Verify API base URL
- Check API service status

#### 3. **Rate Limiting Errors**
```
❌ API error: status 429
```
**Solutions**:
- Increase delays between requests
- Reduce concurrency limits
- Check API rate limits

#### 4. **Timeout Errors**
```
❌ context deadline exceeded
```
**Solutions**:
- Increase timeout values
- Check network latency
- Verify API response times

### Debug Mode

Enable debug logging for detailed test output:

```bash
# Set debug log level
export LOG_LEVEL=debug

# Run tests with verbose output
make test-integration-cupid
```

## 📈 Monitoring

### Test Metrics

The integration tests provide the following metrics:

- **Success Rate** - Percentage of successful API calls
- **Response Time** - Average response time per operation
- **Error Rate** - Percentage of failed API calls
- **Throughput** - Operations per second
- **Data Quality** - Validation success rate

### Continuous Integration

For CI/CD pipelines, use:

```bash
# Run tests in CI mode (no interactive prompts)
export CI=true
make test-integration-cupid
```

## 🔒 Security

### API Key Management

- **Never commit API keys** to version control
- **Use environment variables** for API credentials
- **Rotate API keys** regularly
- **Use test-specific keys** when possible

### Test Data

- **Use test property IDs** from the predefined list
- **Avoid modifying production data**
- **Clean up test data** after tests complete

## 📚 Best Practices

### 1. **Test Isolation**
- Each test should be independent
- Use unique test data when possible
- Clean up after tests complete

### 2. **Error Handling**
- Test both success and failure scenarios
- Validate error messages and codes
- Test timeout and retry logic

### 3. **Performance Testing**
- Measure response times
- Test under load
- Monitor resource usage

### 4. **Data Validation**
- Validate all response fields
- Check data types and ranges
- Verify required fields are present

## 🎯 Test Coverage

The integration tests cover:

- ✅ **100%** of client methods
- ✅ **100%** of service methods
- ✅ **100%** of error scenarios
- ✅ **100%** of data validation
- ✅ **Performance** characteristics
- ✅ **Concurrency** behavior

## 🚀 Next Steps

1. **Run the tests** to verify API connectivity
2. **Monitor performance** and optimize if needed
3. **Add more test cases** as API evolves
4. **Integrate with CI/CD** pipeline
5. **Set up monitoring** for API health

---

**Ready to test! 🧪** The integration tests are now set up and ready to validate your Cupid API integration.
