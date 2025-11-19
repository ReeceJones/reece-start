package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reece.start/internal/api"
	"reece.start/internal/authentication"
	"reece.start/internal/constants"
	testconfig "reece.start/test/config"
)

func TestJwtAuthMiddleware(t *testing.T) {
	t.Run("ValidTokenFromHeader", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Create a valid token
		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
		})
		require.NoError(t, err)

		// Create a handler that checks if claims are set
		handler := func(c echo.Context) error {
			claims := c.Get("claims")
			require.NotNil(t, claims)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		// Setup middleware
		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		// Make request with Authorization header
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Assert success
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, true, response["success"])
	})

	t.Run("ValidTokenFromCookie", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Create a valid token
		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
		})
		require.NoError(t, err)

		// Create a handler that checks if claims are set
		handler := func(c echo.Context) error {
			claims := c.Get("claims")
			require.NotNil(t, claims)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		// Setup middleware
		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		// Make request with cookie
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "app-session-token",
			Value: token,
		})
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Assert success
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, true, response["success"])
	})

	t.Run("CookieTakesPrecedenceOverHeader", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Create two tokens
		cookieToken, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
		})
		require.NoError(t, err)

		headerToken, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 456,
		})
		require.NoError(t, err)

		// Create a handler that checks user ID from claims
		handler := func(c echo.Context) error {
			claims := c.Get("claims").(*authentication.JwtClaims)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"user_id": claims.UserId,
			})
		}

		// Setup middleware
		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		// Make request with both cookie and header
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "app-session-token",
			Value: cookieToken,
		})
		req.Header.Set("Authorization", "Bearer "+headerToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Assert cookie token was used (user ID 123)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "123", response["user_id"])
	})

	t.Run("InvalidToken", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Add error handling middleware
		e.Use(ErrorHandlingMiddleware)

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Assert unauthorized
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "invalid token", apiErr.Message)
	})

	t.Run("MissingToken", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Add error handling middleware
		e.Use(ErrorHandlingMiddleware)

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Assert unauthorized
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse response body
		body := rec.Body.Bytes()
		require.NotEmpty(t, body, "Response body should not be empty")

		var apiErr api.ApiError
		err := json.Unmarshal(body, &apiErr)
		require.NoError(t, err, "Should be able to unmarshal response")
		assert.Equal(t, "missing authorization header", apiErr.Message)
	})

	t.Run("InvalidAuthorizationFormat", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Add error handling middleware
		e.Use(ErrorHandlingMiddleware)

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Assert unauthorized
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse response body
		body := rec.Body.Bytes()
		require.NotEmpty(t, body, "Response body should not be empty")

		var apiErr api.ApiError
		err := json.Unmarshal(body, &apiErr)
		require.NoError(t, err, "Should be able to unmarshal response")
		assert.Equal(t, "invalid authorization format", apiErr.Message)
	})
}

func TestGetUserIDFromJWT(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Create token with user ID
		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		userID, err := GetUserIDFromJWT(c)
		require.NoError(t, err)
		assert.Equal(t, uint(123), userID)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		e := echo.New()

		// Create claims with invalid user ID
		claims := &authentication.JwtClaims{
			UserId: "invalid",
		}

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		userID, err := GetUserIDFromJWT(c)
		require.Error(t, err)
		assert.Equal(t, uint(0), userID)
	})

	t.Run("MissingClaims", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())

		// Don't set claims - should panic or return error
		assert.Panics(t, func() {
			_, _ = GetUserIDFromJWT(c)
		})
	})
}

func TestGetRoleFromJWT(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		userRole := constants.UserRoleAdmin
		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
			Role:   &userRole,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		role, err := GetRoleFromJWT(c)
		require.NoError(t, err)
		assert.Equal(t, constants.UserRoleAdmin, role)
	})

	t.Run("MissingRole", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
			Role:   nil,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		role, err := GetRoleFromJWT(c)
		require.Error(t, err)
		assert.Equal(t, constants.UserRole(""), role)
		assert.Contains(t, err.Error(), "role is not set")
	})
}

func TestGetScopesFromJWT(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		scopes := []constants.UserScope{
			constants.UserScopeAdmin,
			constants.UserScopeOrganizationRead,
		}
		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
			Scopes: &scopes,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		retrievedScopes, err := GetScopesFromJWT(c)
		require.NoError(t, err)
		assert.Len(t, retrievedScopes, 2)
		assert.Equal(t, constants.UserScopeAdmin, retrievedScopes[0])
		assert.Equal(t, constants.UserScopeOrganizationRead, retrievedScopes[1])
	})

	t.Run("MissingScopes", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
			Scopes: nil,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		scopes, err := GetScopesFromJWT(c)
		require.Error(t, err)
		assert.Empty(t, scopes)
		assert.Contains(t, err.Error(), "scopes are not set")
	})
}

func TestGetImpersonatingUserIDFromJWT(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		impersonatingUserID := "789"
		isImpersonating := true
		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId:              123,
			IsImpersonating:     &isImpersonating,
			ImpersonatingUserId: &impersonatingUserID,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		userID, err := GetImpersonatingUserIDFromJWT(c)
		require.NoError(t, err)
		assert.Equal(t, uint(789), userID)
	})

	t.Run("MissingImpersonatingUserID", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		userID, err := GetImpersonatingUserIDFromJWT(c)
		require.Error(t, err)
		assert.Equal(t, uint(0), userID)
		assert.Contains(t, err.Error(), "impersonating user ID is not set")
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		impersonatingUserID := "invalid"
		isImpersonating := true
		token, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId:              123,
			IsImpersonating:     &isImpersonating,
			ImpersonatingUserId: &impersonatingUserID,
		})
		require.NoError(t, err)

		claims, err := authentication.ValidateJWT(config, token)
		require.NoError(t, err)

		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/test", nil), httptest.NewRecorder())
		c.Set("claims", claims)

		userID, err := GetImpersonatingUserIDFromJWT(c)
		require.Error(t, err)
		assert.Equal(t, uint(0), userID)
	})
}

// Test private functions by testing through public API
func TestGetTokenFromRequest(t *testing.T) {
	t.Run("CookieFirst", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		// Create two tokens
		cookieToken, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 123,
		})
		require.NoError(t, err)

		headerToken, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 456,
		})
		require.NoError(t, err)

		// Create a handler that extracts user ID to verify which token was used
		handler := func(c echo.Context) error {
			claims := c.Get("claims").(*authentication.JwtClaims)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"user_id": claims.UserId,
			})
		}

		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "app-session-token",
			Value: cookieToken,
		})
		req.Header.Set("Authorization", "Bearer "+headerToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Cookie token should be used (user ID 123)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "123", response["user_id"])
	})

	t.Run("HeaderFallback", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		e := echo.New()

		headerToken, err := authentication.CreateJWT(config, authentication.JwtOptions{
			UserId: 456,
		})
		require.NoError(t, err)

		handler := func(c echo.Context) error {
			claims := c.Get("claims").(*authentication.JwtClaims)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"user_id": claims.UserId,
			})
		}

		middleware := JwtAuthMiddleware(config)
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		// No cookie, only header
		req.Header.Set("Authorization", "Bearer "+headerToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Header token should be used (user ID 456)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "456", response["user_id"])
	})
}
