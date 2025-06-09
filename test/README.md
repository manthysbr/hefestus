# Hefestus API Tests

This directory contains comprehensive tests for the Hefestus API, including unit tests for the refactored models and E2E tests for the Swagger UI.

## Test Structure

- **Unit Tests**: Located in `internal/models/models_test.go` - Tests the refactored data models
- **Integration Tests**: Located in `test/integration_test.go` - Tests API endpoints and Swagger UI functionality
- **E2E Tests**: Located in `test/e2e/swagger-ui.spec.js` - Playwright tests for full Swagger UI interaction

## Running Tests

### Unit Tests
```bash
# Run unit tests for models
go test ./internal/models -v

# Run with coverage
go test ./internal/models -cover
```

### Integration Tests
```bash
# First, start the server in one terminal:
go run cmd/server/main.go

# Then run integration tests in another terminal:
go test ./test -v
```

### E2E Tests (Playwright)
```bash
# Install dependencies (one time)
cd test
npm install
npx playwright install

# Run E2E tests (requires server to be running)
cd test
npm run test:e2e

# Run with browser visible
npm run test:e2e:headed

# Debug tests
npm run test:e2e:debug
```

## Test Coverage

### Models Package
- **100% code coverage** for all refactored models
- Tests JSON serialization/deserialization
- Tests constructor functions
- Tests field validation tags
- Tests backward compatibility of JSON field names

### Integration Tests
- Swagger UI accessibility
- Swagger JSON schema validation
- API endpoint functionality
- Model compatibility validation
- Constructor function behavior

### E2E Tests
- Complete Swagger UI interaction
- Model schema validation in browser
- API endpoint testing through Swagger UI
- Domain validation testing
- Response structure validation

## Refactoring Compliance

The tests ensure that the `models.go` refactoring follows "Effective Go" principles:

1. **Package Documentation**: ✅ Added comprehensive package documentation
2. **Naming Conventions**: ✅ Use English field names (Cause/Solution) with JSON tags for compatibility
3. **Constructor Functions**: ✅ Added `New*` functions for all major structs
4. **Grouping**: ✅ Logically grouped related structures
5. **Documentation**: ✅ Every exported type has proper documentation
6. **Backward Compatibility**: ✅ JSON field names remain in Portuguese for API compatibility

## Notes

- Integration tests automatically skip if the server is not running
- E2E tests require Playwright browsers to be installed
- All tests maintain backward compatibility with the existing API
- The refactored models use English field names internally but maintain Portuguese JSON field names