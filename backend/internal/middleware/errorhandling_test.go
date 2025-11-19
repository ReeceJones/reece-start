package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/gorm"
	"reece.start/internal/api"
)

func TestErrorHandlingMiddleware(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, true, response["success"])
	})

	t.Run("ErrForbiddenNoAdminAccess", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrForbiddenNoAdminAccess
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrForbiddenNoAdminAccess.Error(), apiErr.Message)
	})

	t.Run("ErrForbiddenOwnProfileOnly", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrForbiddenOwnProfileOnly
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrForbiddenOwnProfileOnly.Error(), apiErr.Message)
	})

	t.Run("ErrUnauthorizedInvalidLogin", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrUnauthorizedInvalidLogin
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrUnauthorizedInvalidLogin.Error(), apiErr.Message)
	})

	t.Run("ErrMembershipNotFound", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrMembershipNotFound
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrMembershipNotFound.Error(), apiErr.Message)
	})

	t.Run("ErrInvitationNotFound", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrInvitationNotFound
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrInvitationNotFound.Error(), apiErr.Message)
	})

	t.Run("ErrInvitationAlreadyExists", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrInvitationAlreadyExists
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrInvitationAlreadyExists.Error(), apiErr.Message)
	})

	t.Run("ErrUserEmailAlreadyExists", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrUserEmailAlreadyExists
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrUserEmailAlreadyExists.Error(), apiErr.Message)
	})

	t.Run("ErrForbiddenNoAccess", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrForbiddenNoAccess
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrForbiddenNoAccess.Error(), apiErr.Message)
	})

	t.Run("ErrInvalidOrganizationID", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrInvalidOrganizationID
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrInvalidOrganizationID.Error(), apiErr.Message)
	})

	t.Run("ErrInvalidUserID", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrInvalidUserID
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrInvalidUserID.Error(), apiErr.Message)
	})

	t.Run("ErrInvalidMembershipID", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrInvalidMembershipID
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrInvalidMembershipID.Error(), apiErr.Message)
	})

	t.Run("ErrInvalidInvitationID", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrInvalidInvitationID
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrInvalidInvitationID.Error(), apiErr.Message)
	})

	t.Run("ErrStripeWebhookSecretNotConfigured", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrStripeWebhookSecretNotConfigured
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrStripeWebhookSecretNotConfigured.Error(), apiErr.Message)
	})

	t.Run("ErrStripeWebhookSignatureMissing", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrStripeWebhookSignatureMissing
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrStripeWebhookSignatureMissing.Error(), apiErr.Message)
	})

	t.Run("ErrStripeWebhookSignatureInvalid", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrStripeWebhookSignatureInvalid
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrStripeWebhookSignatureInvalid.Error(), apiErr.Message)
	})

	t.Run("ErrStripeWebhookEventInvalid", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrStripeWebhookEventInvalid
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrStripeWebhookEventInvalid.Error(), apiErr.Message)
	})

	t.Run("ErrStripeWebhookEventUnhandled", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return api.ErrStripeWebhookEventUnhandled
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, api.ErrStripeWebhookEventUnhandled.Error(), apiErr.Message)
	})

	t.Run("ApiError", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return &api.ApiError{
				Message: "custom error message",
			}
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "custom error message", apiErr.Message)
	})

	t.Run("GormRecordNotFound", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return gorm.ErrRecordNotFound
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "resource not found", apiErr.Message)
	})

	t.Run("EchoHTTPError", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "bad request", apiErr.Message)
	})

	t.Run("EchoHTTPErrorWithApiError", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusBadRequest, api.ApiError{
				Message: "nested api error",
			})
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "nested api error", apiErr.Message)
	})

	t.Run("UnknownError", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return errors.New("unknown error")
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "unknown error", apiErr.Message)
	})

	t.Run("WrappedError", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return errors.Join(api.ErrForbiddenNoAdminAccess, errors.New("additional context"))
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// errors.Is() works with joined errors, so it should match ErrForbiddenNoAdminAccess
		assert.Equal(t, http.StatusForbidden, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		// The error message will be the combined message from errors.Join()
		assert.Contains(t, apiErr.Message, api.ErrForbiddenNoAdminAccess.Error())
		assert.Contains(t, apiErr.Message, "additional context")
	})

	t.Run("StripeError", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return &stripeGo.Error{
				Msg:            "Your card was declined.",
				HTTPStatusCode: http.StatusPaymentRequired,
			}
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusPaymentRequired, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "Your card was declined.", apiErr.Message)
	})

	t.Run("StripeErrorWithEmptyMsg", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return &stripeGo.Error{
				Msg:            "",
				HTTPStatusCode: http.StatusBadRequest,
			}
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		// Should fallback to Error() method when Msg is empty
		assert.NotEmpty(t, apiErr.Message)
	})

	t.Run("StripeErrorWithoutStatusCode", func(t *testing.T) {
		e := echo.New()

		handler := func(c echo.Context) error {
			return &stripeGo.Error{
				Msg:            "Invalid API key provided.",
				HTTPStatusCode: 0, // No status code set
			}
		}

		middleware := ErrorHandlingMiddleware
		e.GET("/test", handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// Should default to BadRequest when HTTPStatusCode is 0
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var apiErr api.ApiError
		err := json.Unmarshal(rec.Body.Bytes(), &apiErr)
		require.NoError(t, err)
		assert.Equal(t, "Invalid API key provided.", apiErr.Message)
	})
}
