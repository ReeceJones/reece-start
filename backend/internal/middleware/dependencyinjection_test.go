package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/resend/resend-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stripeGo "github.com/stripe/stripe-go/v83"
	"reece.start/internal/configuration"
	testdb "reece.start/test/db"
)

func createTestDependencies() AppDependencies {
	// Create test dependencies with minimal configuration
	// DB will be set in tests that need it using testdb.SetupDB()
	return AppDependencies{
		Config: &configuration.Config{
			Host:                   "localhost",
			Port:                   "8080",
			FrontendUrl:            "http://localhost:3000",
			DatabaseUri:            "postgres://test",
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
		},
		DB:           nil, // Will be set in tests that need it
		MinioClient:  nil,
		RiverClient:  nil,
		ResendClient: nil,
		StripeClient: nil,
	}
}

func TestDependencyInjectionMiddleware(t *testing.T) {
	t.Run("SetsAllDependencies", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		// Use shared test database connection
		db := testdb.SetupDB(t)
		deps.DB = db

		handler := func(c echo.Context) error {
			// Verify all dependencies are set
			config := GetConfig(c)
			assert.NotNil(t, config)
			assert.Equal(t, "localhost", config.Host)

			db := GetDB(c)
			assert.NotNil(t, db)

			minioClient := GetMinioClient(c)
			// Can be nil in tests
			_ = minioClient

			riverClient := GetRiverClient(c)
			// Can be nil in tests
			_ = riverClient

			resendClient := GetResendClient(c)
			// Can be nil in tests
			_ = resendClient

			stripeClient := GetStripeClient(c)
			// Can be nil in tests
			_ = stripeClient

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SetsConfig", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		handler := func(c echo.Context) error {
			config := GetConfig(c)
			require.NotNil(t, config)
			assert.Equal(t, "localhost", config.Host)
			assert.Equal(t, "8080", config.Port)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"host": config.Host,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SetsDB", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		// Use shared test database connection
		db := testdb.SetupDB(t)
		deps.DB = db

		handler := func(c echo.Context) error {
			db := GetDB(c)
			require.NotNil(t, db)
			// Verify it's the same instance
			assert.Equal(t, db, deps.DB)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SetsMinioClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		// Create a mock Minio client
		minioClient, err := minio.New("localhost:9000", &minio.Options{
			Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
			Secure: false,
		})
		if err != nil {
			t.Skip("Skipping test: Minio client creation failed")
		}
		deps.MinioClient = minioClient

		handler := func(c echo.Context) error {
			client := GetMinioClient(c)
			require.NotNil(t, client)
			assert.Equal(t, minioClient, client)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SetsNilMinioClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()
		deps.MinioClient = nil

		handler := func(c echo.Context) error {
			client := GetMinioClient(c)
			// Should handle nil gracefully
			assert.Nil(t, client)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SetsResendClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()
		deps.ResendClient = resend.NewClient("test-key")

		handler := func(c echo.Context) error {
			client := GetResendClient(c)
			require.NotNil(t, client)
			assert.Equal(t, deps.ResendClient, client)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SetsStripeClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()
		deps.StripeClient = stripeGo.NewClient("sk_test_mock")

		handler := func(c echo.Context) error {
			client := GetStripeClient(c)
			require.NotNil(t, client)
			assert.Equal(t, deps.StripeClient, client)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("ReturnsConfig", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		handler := func(c echo.Context) error {
			config := GetConfig(c)
			require.NotNil(t, config)
			assert.Equal(t, deps.Config, config)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"host": config.Host,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PanicsWhenNotSet", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())

		// Don't set config in context
		assert.Panics(t, func() {
			_ = GetConfig(c)
		})
	})
}

func TestGetDB(t *testing.T) {
	t.Run("ReturnsDB", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		// Use shared test database connection
		db := testdb.SetupDB(t)
		deps.DB = db

		handler := func(c echo.Context) error {
			retrievedDB := GetDB(c)
			require.NotNil(t, retrievedDB)
			assert.Equal(t, db, retrievedDB)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PanicsWhenNotSet", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())

		// Don't set DB in context
		assert.Panics(t, func() {
			_ = GetDB(c)
		})
	})
}

func TestGetMinioClient(t *testing.T) {
	t.Run("ReturnsClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		minioClient, err := minio.New("localhost:9000", &minio.Options{
			Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
			Secure: false,
		})
		if err != nil {
			t.Skip("Skipping test: Minio client creation failed")
		}
		deps.MinioClient = minioClient

		handler := func(c echo.Context) error {
			client := GetMinioClient(c)
			assert.Equal(t, minioClient, client)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PanicsWhenNotSet", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())

		// Don't set MinioClient in context
		assert.Panics(t, func() {
			_ = GetMinioClient(c)
		})
	})
}

func TestGetRiverClient(t *testing.T) {
	t.Run("ReturnsClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()

		// River client requires a database connection
		// For this test, we'll just verify the middleware sets it
		handler := func(c echo.Context) error {
			client := GetRiverClient(c)
			// Can be nil in tests
			_ = client
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PanicsWhenNotSet", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())

		// Don't set RiverClient in context
		assert.Panics(t, func() {
			_ = GetRiverClient(c)
		})
	})
}

func TestGetResendClient(t *testing.T) {
	t.Run("ReturnsClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()
		deps.ResendClient = resend.NewClient("test-key")

		handler := func(c echo.Context) error {
			client := GetResendClient(c)
			require.NotNil(t, client)
			assert.Equal(t, deps.ResendClient, client)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PanicsWhenNotSet", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())

		// Don't set ResendClient in context
		assert.Panics(t, func() {
			_ = GetResendClient(c)
		})
	})
}

func TestGetStripeClient(t *testing.T) {
	t.Run("ReturnsClient", func(t *testing.T) {
		e := echo.New()
		deps := createTestDependencies()
		deps.StripeClient = stripeGo.NewClient("sk_test_mock")

		handler := func(c echo.Context) error {
			client := GetStripeClient(c)
			require.NotNil(t, client)
			assert.Equal(t, deps.StripeClient, client)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := DependencyInjectionMiddleware(deps)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PanicsWhenNotSet", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())

		// Don't set StripeClient in context
		assert.Panics(t, func() {
			_ = GetStripeClient(c)
		})
	})
}
