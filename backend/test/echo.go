package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/resend/resend-go/v2"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
	echoServer "reece.start/internal/echo"
	"reece.start/internal/jobs"
	appMiddleware "reece.start/internal/middleware"
	testdb "reece.start/test/db"
	"reece.start/test/mocks"
)

// TestContext holds all the testing infrastructure
type TestContext struct {
	T            *testing.T
	Echo         *echo.Echo
	DB           *gorm.DB
	SQLConn      *sql.DB
	Config       *configuration.Config
	MinioClient  *minio.Client
	RiverClient  *river.Client[*sql.Tx]
	ResendClient *resend.Client
	StripeClient *stripeGo.Client
}

// SetupEchoTest sets up all of the testing infrastructure to run integration tests against echo handlers.
// It creates a PostgreSQL testcontainer, initializes all dependencies, and creates an Echo server
// with the same middleware and routes as the production server.
func SetupEchoTest(t *testing.T) *TestContext {
	// Setup Postgres
	sqlConn, db, connStr := testdb.SetupPostgresContainer(t)

	// Create test config
	config := &configuration.Config{
		Host:                   "localhost",
		Port:                   "8080",
		FrontendUrl:            "http://localhost:3000",
		DatabaseUri:            connStr,
		JwtSecret:              "test-secret",
		JwtIssuer:              "test-issuer",
		JwtAudience:            "test-audience",
		JwtExpirationTime:      86400,
		StorageEndpoint:        "localhost:9000",
		StorageAccessKeyId:     "minioadmin",
		StorageSecretAccessKey: "minioadmin",
		StorageUseSSL:          false,
		EnableEmail:            false,
		ResendApiKey:           "test-key",
	}

	// Replace the default HTTP transport FIRST to intercept all external API calls
	// This prevents Stripe, Resend, and any other external services from making real HTTP requests
	// Uses the shared testmocks package to avoid import cycles
	mocks.ReplaceDefaultTransportWithCleanup(t)

	// Create mock clients (for non-critical services)
	var minioClient *minio.Client // nil for now - tests can mock as needed

	// Create Resend client - HTTP calls will be intercepted by MockHTTPTransport
	// Also, EnableEmail is false in tests as a double safeguard
	resendClient := resend.NewClient("test-key")

	// Create mock Stripe client - HTTP calls will be intercepted by MockHTTPTransport
	stripeClient := NewMockStripeClient()

	// Create River client for background jobs (workers registered but NOT started in tests)
	// River tables are already created during initial migration in setupSharedPostgresContainer
	riverClient, err := jobs.NewRiverClient(t.Context(), jobs.RiverClientConfig{
		SQLConn:      sqlConn,
		DB:           db,
		Config:       config,
		ResendClient: resendClient,
		StripeClient: stripeClient,
		StartWorkers: false, // Don't start workers in tests
	})
	require.NoError(t, err)

	// Create Echo server with all middleware and routes (same as production)
	e := echoServer.NewEcho(appMiddleware.AppDependencies{
		Config:       config,
		DB:           db,
		MinioClient:  minioClient,
		RiverClient:  riverClient,
		ResendClient: resendClient,
		StripeClient: stripeClient,
	})

	return &TestContext{
		T:            t,
		Echo:         e,
		DB:           db,
		SQLConn:      sqlConn,
		Config:       config,
		MinioClient:  minioClient,
		RiverClient:  riverClient,
		ResendClient: resendClient,
		StripeClient: stripeClient,
	}
}

// MakeRequest makes an HTTP request to the test server and returns the response
func (tc *TestContext) MakeRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(tc.T, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	rec := httptest.NewRecorder()
	tc.Echo.ServeHTTP(rec, req)

	return rec
}

// MakeAuthenticatedRequest makes an authenticated HTTP request using a JWT token
func (tc *TestContext) MakeAuthenticatedRequest(method string, path string, body interface{}, token string) *httptest.ResponseRecorder {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	return tc.MakeRequest(method, path, body, headers)
}

// UnmarshalResponse unmarshals the response body into the target struct
func (tc *TestContext) UnmarshalResponse(rec *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(rec.Body.Bytes(), target)
	require.NoError(tc.T, err)
}
