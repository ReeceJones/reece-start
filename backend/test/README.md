# Integration Testing Infrastructure

This package provides comprehensive testing infrastructure for running integration tests against Echo HTTP handlers using testcontainers.

## Overview

The test infrastructure sets up a complete test environment that mirrors the production server setup:

- **PostgreSQL Database**: Spun up using testcontainers
- **Echo HTTP Server**: Configured with all middleware and routes
- **Dependencies**: All application dependencies (DB, config, clients) properly initialized
- **Fixtures**: Helper functions for creating test data
- **HTTP Helpers**: Utilities for making requests and assertions

## Files

### `postgres.go`

Sets up a PostgreSQL testcontainer with:

- Automatic connection string generation
- GORM initialization
- Database migrations
- Automatic cleanup

### `echo.go`

Sets up the Echo HTTP server with:

- Uses `server.NewEcho()` - the same initialization as production
- Test context holding all dependencies
- HTTP request helpers (`MakeRequest`, `MakeAuthenticatedRequest`)
- Response unmarshaling utilities

### `fixtures.go`

Provides helper functions for creating test data:

- `CreateTestUser`: Create a basic user
- `CreateTestAdminUser`: Create an admin user
- `CreateTestOrganization`: Create an organization
- `CreateTestOrganizationMembership`: Create a membership
- `CreateTestJWT`: Create a JWT token
- `CreateTestUserWithOrganization`: Create a user with organization
- `CreateAuthenticatedTestUser`: Create a user with organization and JWT token

## Shared Initialization

Several key components are shared between production and tests:

### Echo Server

- **`internal/server/server.go`**: `NewEcho(deps)` - creates Echo with all middleware and routes
- **`internal/routes/routes.go`**: `Register(e, config)` - registers all application routes

### River Client (Background Jobs)

- **`internal/jobs/river.go`**: `NewRiverClient(ctx, cfg)` - creates River client with all workers
- **Production**: Workers are started (`StartWorkers: true`)
- **Tests**: Workers registered but NOT started (`StartWorkers: false`) to prevent job processing errors

Both `backend/server.go` (production) and `backend/test/echo.go` (tests) use these shared functions.
This ensures tests use the **exact same configuration** as production.

## Usage

### Basic Test Setup

```go
func TestMyEndpoint(t *testing.T) {
    // Setup test context with all dependencies
    tc := test.SetupEchoTest(t)

    // Make a request
    rec := tc.MakeRequest(http.MethodGet, "/my-endpoint", nil, nil)

    // Assert response
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

### Testing Public Endpoints

```go
func TestCreateUserEndpoint(t *testing.T) {
    tc := test.SetupEchoTest(t)

    reqBody := map[string]interface{}{
        "data": map[string]interface{}{
            "type": constants.ApiTypeUser,
            "attributes": map[string]interface{}{
                "name":     "Test User",
                "email":    "test@example.com",
                "password": "password123",
            },
        },
    }

    rec := tc.MakeRequest(http.MethodPost, "/users", reqBody, nil)

    assert.Equal(t, http.StatusCreated, rec.Code)
}
```

### Testing Authenticated Endpoints

```go
func TestAuthenticatedEndpoint(t *testing.T) {
    tc := test.SetupEchoTest(t)

    // Create authenticated user with organization
    user, org, token := test.CreateAuthenticatedTestUser(
        t, tc.DB, tc.Config,
        "Test User", "test@example.com", "password123",
        "Test Org", constants.OrganizationRoleAdmin,
    )

    // Make authenticated request
    rec := tc.MakeAuthenticatedRequest(
        http.MethodGet,
        "/users/me",
        nil,
        token,
    )

    assert.Equal(t, http.StatusOK, rec.Code)
}
```

### Creating Test Data

```go
func TestWithExistingData(t *testing.T) {
    tc := test.SetupEchoTest(t)

    // Create a test user
    user := test.CreateTestUser(
        t, tc.DB,
        "Test User",
        "test@example.com",
        "password123",
    )

    // Create an organization
    org := test.CreateTestOrganization(t, tc.DB, "Test Org")

    // Create membership
    membership := test.CreateTestOrganizationMembership(
        t, tc.DB,
        user.ID,
        org.ID,
        constants.OrganizationRoleAdmin,
    )

    // ... rest of test
}
```

### Parsing JSON Responses

```go
func TestResponseParsing(t *testing.T) {
    tc := test.SetupEchoTest(t)

    rec := tc.MakeRequest(http.MethodGet, "/users/me", nil, nil)

    // Parse as generic JSON
    var response map[string]interface{}
    tc.UnmarshalResponse(rec, &response)

    data := response["data"].(map[string]interface{})
    attributes := data["attributes"].(map[string]interface{})

    assert.Equal(t, "Test User", attributes["name"])
}
```

### Testing Background Jobs

```go
func TestJobEnqueueing(t *testing.T) {
    tc := test.SetupEchoTest(t)

    // Make request that enqueues a background job
    rec := tc.MakeRequest(http.MethodPost, "/some-endpoint", reqBody, nil)
    assert.Equal(t, http.StatusOK, rec.Code)

    // Process the enqueued jobs (polls until complete)
    test.RunAllPendingRiverJobs(t, tc.DB, tc.RiverClient)

    // Verify the side effects of the job
    // (e.g., check database, verify email was sent, etc.)
}
```

You can also specify custom timeout and poll interval:

```go
// Wait up to 30 seconds, polling every 200ms
test.RunPendingRiverJobs(t, tc.DB, tc.RiverClient, 30*time.Second, 200*time.Millisecond)
```

## Testing Best Practices

1. **Isolation**: Each test gets a fresh database (via testcontainers)
2. **Cleanup**: Testcontainers automatically cleans up after each test
3. **Fixtures**: Use fixture functions to create consistent test data
4. **Assertions**: Use testify/assert for clear assertions
5. **Error Messages**: Add descriptive error messages to assertions

## Example Test File

See `backend/internal/users/http_test.go` for comprehensive examples of:

- Testing user creation
- Testing login
- Testing duplicate emails
- Testing authenticated endpoints
- Testing updates
- Database verification

## Running Tests

```bash
# Run all tests
go test ./...

# Run specific test file
go test ./internal/users/http_test.go

# Run specific test
go test -run TestCreateUserEndpoint_ValidUser ./internal/users

# Run with verbose output
go test -v ./...
```

## Requirements

- Docker (for testcontainers)
- Go 1.21+
- PostgreSQL testcontainer image will be pulled automatically

## Performance

- Each test creates a fresh PostgreSQL container
- Tests can run in parallel (each gets its own container)
- Cleanup is automatic via `t.Cleanup()`
- Typical test startup: 2-5 seconds (container spin-up time)
