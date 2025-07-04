# E2E Tests for OGP Verification Service

This directory contains end-to-end tests for the OGP Verification Service API.

## Overview

The E2E tests verify the complete functionality of the service by making real HTTP requests to a running instance of the API and validating the responses.

## Test Coverage

### Core Functionality
- ✅ Health endpoint verification
- ✅ OGP verification with complete data (GitHub)
- ✅ OGP verification with good data (Wikipedia)
- ✅ OGP verification with minimal data (Example.com)

### Error Handling
- ✅ Empty URL validation
- ✅ Invalid JSON handling
- ✅ Private IP blocking

### Security & Performance
- ✅ Rate limiting (10 requests/minute per IP)
- ✅ CORS headers validation
- ✅ Method not allowed responses
- ✅ API latency testing (< 10 seconds)
- ✅ Concurrent request handling

## Running the Tests

### Prerequisites
- Docker and Docker Compose
- Network access to test websites (GitHub, Wikipedia, Example.com)

### Run All E2E Tests

```bash
cd backend/tests/e2e
./run-e2e.sh
```

### Manual Testing

1. Start the test environment:
```bash
docker-compose -f docker-compose.e2e.yml up --build
```

2. Run specific tests:
```bash
docker-compose -f docker-compose.e2e.yml exec e2e-tests go test -v -run TestHealthEndpoint
```

3. Cleanup:
```bash
docker-compose -f docker-compose.e2e.yml down --volumes --remove-orphans
```

## Test Structure

### Test Files
- `main_test.go` - Core E2E test implementation
- `docker-compose.e2e.yml` - Test environment configuration
- `Dockerfile.e2e` - Test runner container
- `run-e2e.sh` - Test execution script

### Test Environment
- **Backend Service**: Runs on port 8080 with health checks
- **Test Runner**: Go test environment with curl and jq tools
- **Network**: Isolated Docker network for testing

## Test Scenarios

### 1. Health Check (`TestHealthEndpoint`)
Verifies the `/health` endpoint returns 200 OK.

### 2. GitHub OGP Test (`TestOGPVerifyEndpoint_GitHub`)
Tests with a website that has complete OGP metadata:
- Validates all required OGP fields are present
- Checks platform-specific previews
- Verifies validation results

### 3. Wikipedia OGP Test (`TestOGPVerifyEndpoint_Wikipedia`)
Tests with a website that has good OGP metadata:
- Validates title and description presence
- Checks validation flags

### 4. Example.com Test (`TestOGPVerifyEndpoint_ExampleCom`)
Tests with a minimal OGP implementation:
- Expects warnings but still valid response
- Validates error handling for missing data

### 5. Error Cases (`TestOGPVerifyEndpoint_ErrorCases`)
Tests various error conditions:
- Empty URL (400 Bad Request)
- Invalid JSON (400 Bad Request)
- Private IP (500 Internal Server Error)

### 6. Rate Limiting (`TestRateLimiting`)
Validates the 10 requests/minute rate limit:
- Makes 10 successful requests
- Verifies 11th request returns 429

### 7. CORS Headers (`TestCORSHeaders`)
Validates CORS preflight handling:
- Tests OPTIONS method
- Checks required CORS headers

### 8. Method Validation (`TestMethodNotAllowed`)
Tests unsupported HTTP methods return 405.

### 9. Performance (`TestAPILatency`)
Validates API response time < 10 seconds.

### 10. Concurrency (`TestConcurrentRequests`)
Tests handling of multiple simultaneous requests.

## Expected Results

All tests should pass with:
- ✅ 10/10 tests passing
- No test failures or errors
- Performance within acceptable limits
- Proper error handling

## Troubleshooting

### Common Issues

1. **Network connectivity issues**
   - Ensure Docker has internet access
   - Check if test websites are accessible

2. **Rate limiting during development**
   - Wait 1 minute between test runs
   - Or restart the backend service

3. **Health check failures**
   - Increase health check timeout in docker-compose.e2e.yml
   - Check backend logs for startup issues

### Debug Commands

```bash
# View backend logs
docker-compose -f docker-compose.e2e.yml logs backend-e2e

# View test logs
docker-compose -f docker-compose.e2e.yml logs e2e-tests

# Test specific endpoint manually
docker-compose -f docker-compose.e2e.yml exec e2e-tests curl -X POST http://backend-e2e:8080/api/v1/ogp/verify -H "Content-Type: application/json" -d '{"url":"https://github.com"}'
```

## Integration with CI/CD

These E2E tests can be integrated into GitHub Actions or other CI/CD pipelines:

```yaml
- name: Run E2E Tests
  run: |
    cd backend/tests/e2e
    ./run-e2e.sh
```