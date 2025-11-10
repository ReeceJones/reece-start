# Testing Guide

## Overview

This project uses integration tests that run against real infrastructure (PostgreSQL database) using testcontainers. The testing infrastructure is designed to be fast and efficient while providing a realistic testing environment.

## Integration Tests

### When to Use Integration Tests

**Use integration tests sparingly.** They are significantly slower than unit tests because they:

- Start and manage Docker containers
- Execute real database operations
- Run against actual infrastructure

Integration tests should be used for:

- Testing HTTP endpoints end-to-end
- Testing database operations and migrations
- Testing complex workflows that require multiple components
- Testing authentication and authorization flows

For simple logic, data transformations, or isolated functions, prefer unit tests with mocks.

### Shared Container Architecture

The test infrastructure uses an optimized shared container pattern to minimize test execution time:

1. **Single Container**: A single PostgreSQL container is started once and reused across all tests
2. **Table Cleanup**: Between each test, all tables are truncated (not dropped) to ensure a clean state
3. **Connection Management**: Each test gets its own database connection, which is closed after the test completes
4. **Automatic Cleanup**: The container is automatically cleaned up when the test process exits

This approach provides significant performance improvements:

- **Before**: Each test started its own container (~5-10 seconds per test)
- **After**: Container starts once (~5-10 seconds total), tests run much faster

### How It Works

The shared container pattern is implemented in `test/postgres.go`:

- `setupSharedPostgresContainer()`: Uses `sync.Once` to ensure the container is only started once
- `SetupPostgresContainer()`: Called by each test to get a fresh database connection
- `CleanAllTables()`: Truncates all tables between tests to ensure isolation

```go
func TestMyEndpoint(t *testing.T) {
    // This will reuse the shared container if it already exists
    tc := test.SetupEchoTest(t)

    // Test implementation...
}
```

### Test Isolation

Each test is isolated through:

- **Fresh Database Connection**: Each test gets its own connection
- **Table Truncation**: All tables are truncated before each test runs
- **Sequence Reset**: Auto-increment sequences are reset (`RESTART IDENTITY`)
- **Foreign Key Handling**: CASCADE ensures all related data is cleaned up

### Running Tests

Run all tests:

```bash
go test ./...
```

Run tests in a specific package:

```bash
go test ./internal/users/...
```

Run a specific test:

```bash
go test -run TestUpdateUserEndpoint ./internal/users/
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Test Infrastructure

### SetupEchoTest

The `test.SetupEchoTest()` function provides a complete test environment:

- PostgreSQL database (shared container)
- Echo HTTP server with all middleware and routes
- Mock clients for external services (Stripe, Resend, etc.)
- Test context with helper methods

```go
func TestMyHandler(t *testing.T) {
    tc := test.SetupEchoTest(t)

    // Make authenticated request
    rec := tc.MakeAuthenticatedRequest("GET", "/api/users", nil, token)

    // Unmarshal response
    var response MyResponse
    tc.UnmarshalResponse(rec, &response)
}
```

### Helper Functions

The `test` package provides several helper functions:

- `SetupEchoTest()`: Complete test environment setup
- `MakeRequest()`: Make HTTP requests
- `MakeAuthenticatedRequest()`: Make authenticated HTTP requests
- `UnmarshalResponse()`: Parse JSON responses
- `CreateTestUser()`: Create test users
- `CreateTestOrganization()`: Create test organizations
- And more...

See `test/README.md` for complete documentation of available helpers.

## Best Practices

1. **Prefer Unit Tests**: Use integration tests only when necessary
2. **Keep Tests Fast**: Minimize database operations and external calls
3. **Clean State**: Don't rely on data from previous tests
4. **Isolation**: Each test should be independent and runnable in isolation
5. **Use Fixtures**: Use helper functions to create test data consistently
6. **Use Subtests**: Use `t.Run()` for test scenarios instead of `TestFoo_TestScenario` naming convention

### Subtests with t.Run

When testing multiple scenarios for the same functionality, use `t.Run()` to create subtests rather than creating separate test functions with underscore naming:

**Preferred:**
```go
func TestCreateUser(t *testing.T) {
    t.Run("creates user successfully", func(t *testing.T) {
        // Test implementation...
    })
    
    t.Run("returns error for duplicate email", func(t *testing.T) {
        // Test implementation...
    })
}
```

**Avoid:**
```go
func TestCreateUser_Success(t *testing.T) {
    // Test implementation...
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
    // Test implementation...
}
```

Benefits of using `t.Run()`:
- Better test organization and readability
- Easier to run specific scenarios: `go test -run TestCreateUser/creates_user_successfully`
- Shared setup/teardown code can be reused
- Clearer test output showing the hierarchy

## Performance Considerations

- Integration tests are slower than unit tests - use them judiciously
- The shared container pattern significantly reduces startup overhead
- Table truncation is fast but still adds some overhead per test
- Consider running integration tests separately from unit tests in CI/CD
