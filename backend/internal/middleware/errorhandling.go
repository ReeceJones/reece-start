package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"reece.start/internal/api"
)

// ErrorHandlingMiddleware provides a centralized way to handle common errors
// It wraps handler functions and automatically converts known errors to appropriate HTTP responses
func ErrorHandlingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		// Handle specific business logic errors first
		if errors.Is(err, api.ErrForbiddenNoAdminAccess) {
			return c.JSON(http.StatusForbidden, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrForbiddenOwnProfileOnly) {
			return c.JSON(http.StatusForbidden, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrUnauthorizedInvalidLogin) {
			return c.JSON(http.StatusUnauthorized, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrMembershipNotFound) {
			return c.JSON(http.StatusNotFound, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrInvitationNotFound) {
			return c.JSON(http.StatusNotFound, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrInvitationAlreadyExists) {
			return c.JSON(http.StatusConflict, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrUserEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, api.ApiError{
				Message: err.Error(),
			})
		}

		// Handle HTTP layer errors
		if errors.Is(err, api.ErrForbiddenNoAccess) {
			return c.JSON(http.StatusForbidden, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrInvalidOrganizationID) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrInvalidMembershipID) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrInvalidInvitationID) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrStripeWebhookSecretNotConfigured) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrStripeWebhookSignatureMissing) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrStripeWebhookSignatureInvalid) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrStripeWebhookEventInvalid) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrStripeWebhookEventUnhandled) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Message: err.Error(),
			})
		}

		// Handle any other ApiError
		if he, ok := err.(*api.ApiError); ok {
			return c.JSON(http.StatusInternalServerError, api.ApiError{
				Message: he.Message,
			})
		}

		// Handle GORM errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, api.ApiError{
				Message: "Resource not found",
			})
		}

		// Handle Echo HTTP errors (if they bubble up)
		if he, ok := err.(*echo.HTTPError); ok {
			// Sometimes api.ApiError shows up in http error message
			if ae, ok := he.Message.(api.ApiError); ok {
				return c.JSON(he.Code, api.ApiError{
					Message: ae.Message,
				})
			}

			return c.JSON(he.Code, api.ApiError{
				Message: he.Message.(string),
			})
		}

		log.Printf("Unhandled error: %v", err)

		// Default to internal server error for unknown errors
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Message: err.Error(),
		})
	}
}
