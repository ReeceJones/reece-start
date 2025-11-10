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
)

func TestContentTypeMiddleware(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, true, response["success"])
	})

	t.Run("ValidJSONWithCharset", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, true, response["success"])
	})

	t.Run("MissingContentType", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		// Don't set Content-Type header
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "missing_content_type", apiErr.Message)
	})

	t.Run("InvalidContentType", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "invalid_content_type", apiErr.Message)
	})

	t.Run("InvalidContentTypeXML", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "application/xml")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "invalid_content_type", apiErr.Message)
	})

	t.Run("InvalidContentTypeFormData", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "multipart/form-data")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "invalid_content_type", apiErr.Message)
	})

	t.Run("InvalidContentTypeWithCharset", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "invalid_content_type", apiErr.Message)
	})

	t.Run("EmptyContentType", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.POST("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "missing_content_type", apiErr.Message)
	})

	t.Run("GETRequest", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ContentTypeMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		// GET requests typically don't have Content-Type
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Should fail because Content-Type is missing
		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "missing_content_type", apiErr.Message)
	})
}
