package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	stripeGo "github.com/stripe/stripe-go/v83"
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
			return respondWithError(c, http.StatusForbidden, err)
		}

		if errors.Is(err, api.ErrForbiddenOwnProfileOnly) {
			return respondWithError(c, http.StatusForbidden, err)
		}

		if errors.Is(err, api.ErrUnauthorizedInvalidLogin) {
			return respondWithError(c, http.StatusUnauthorized, err)
		}

		if errors.Is(err, api.ErrMissingAuthorizationHeader) {
			return respondWithError(c, http.StatusUnauthorized, err)
		}

		if errors.Is(err, api.ErrInvalidAuthorizationFormat) {
			return respondWithError(c, http.StatusUnauthorized, err)
		}

		if errors.Is(err, api.ErrInvalidToken) {
			return respondWithError(c, http.StatusUnauthorized, err)
		}

		if errors.Is(err, api.ErrMembershipNotFound) {
			return respondWithError(c, http.StatusNotFound, err)
		}

		if errors.Is(err, api.ErrInvitationNotFound) {
			return respondWithError(c, http.StatusNotFound, err)
		}

		if errors.Is(err, api.ErrInvitationAlreadyExists) {
			return respondWithError(c, http.StatusConflict, err)
		}

		if errors.Is(err, api.ErrUserEmailAlreadyExists) {
			return respondWithError(c, http.StatusConflict, err)
		}

		// Handle HTTP layer errors
		if errors.Is(err, api.ErrForbiddenNoAccess) {
			return respondWithError(c, http.StatusForbidden, err)
		}

		if errors.Is(err, api.ErrInvalidOrganizationID) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrInvalidUserID) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrInvalidMembershipID) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrInvalidInvitationID) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrStripeWebhookSecretNotConfigured) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrStripeWebhookSignatureMissing) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrStripeWebhookSignatureInvalid) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrStripeWebhookEventInvalid) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		if errors.Is(err, api.ErrStripeWebhookEventUnhandled) {
			return respondWithError(c, http.StatusBadRequest, err)
		}

		// Handle any other ApiError
		if he, ok := err.(*api.ApiError); ok {
			return respondWithError(c, http.StatusInternalServerError, he)
		}

		// Handle Stripe errors
		var stripeErr *stripeGo.Error
		if errors.As(err, &stripeErr) {
			// Use the user-facing message from Stripe error
			message := stripeErr.Msg
			if message == "" {
				// Fallback to Error() if Msg is empty
				message = err.Error()
			}
			// Determine appropriate HTTP status code based on Stripe error type
			statusCode := http.StatusBadRequest
			if stripeErr.HTTPStatusCode > 0 {
				statusCode = stripeErr.HTTPStatusCode
			}
			c.JSON(statusCode, api.ApiError{
				Message: message,
			})
			return nil
		}

		// Handle GORM errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respondWithError(c, http.StatusNotFound, errors.New("resource not found"))
		}

		// Handle Echo HTTP errors (if they bubble up)
		if he, ok := err.(*echo.HTTPError); ok {
			// Sometimes api.ApiError shows up in http error message
			if ae, ok := he.Message.(api.ApiError); ok {
				c.JSON(he.Code, api.ApiError{
					Message: ae.Message,
				})
				return nil
			}

			c.JSON(he.Code, api.ApiError{
				Message: he.Message.(string),
			})
			return nil
		}

		slog.Error("Unhandled error", "error", err)

		// Default to internal server error for unknown errors
		return respondWithError(c, http.StatusInternalServerError, err)
	}
}

func respondWithError(c echo.Context, statusCode int, err error) error {
	c.JSON(statusCode, api.ApiError{
		Message: err.Error(),
	})
	return nil
}
