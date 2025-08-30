package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/constants"
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
				Code:    constants.ErrorCodeForbidden,
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrForbiddenOwnProfileOnly) {
			return c.JSON(http.StatusForbidden, api.ApiError{
				Code:    constants.ErrorCodeForbidden,
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrUnauthorizedInvalidLogin) {
			return c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:    constants.ErrorCodeUnauthorized,
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrMembershipNotFound) {
			return c.JSON(http.StatusNotFound, api.ApiError{
				Code:    constants.ErrorCodeNotFound,
				Message: err.Error(),
			})
		}

		if errors.Is(err, api.ErrInvitationAlreadyExists) {
			return c.JSON(http.StatusConflict, api.ApiError{
				Code:    constants.ErrorCodeConflict,
				Message: err.Error(),
			})
		}

		// Handle HTTP layer errors
		if errors.Is(err, ErrInvalidUserToken) {
			return c.JSON(http.StatusUnauthorized, api.ApiError{
				Code:    constants.ErrorCodeUnauthorized,
				Message: "Invalid user token",
			})
		}

		if errors.Is(err, ErrInvalidOrganizationID) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Code:    constants.ErrorCodeBadRequest,
				Message: "Invalid organization ID",
			})
		}

		if errors.Is(err, ErrInvalidUserID) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Code:    constants.ErrorCodeBadRequest,
				Message: "Invalid user ID",
			})
		}

		if errors.Is(err, ErrInvalidMembershipID) {
			return c.JSON(http.StatusBadRequest, api.ApiError{
				Code:    constants.ErrorCodeBadRequest,
				Message: "Invalid membership ID",
			})
		}

		// Handle GORM errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, api.ApiError{
				Code:    constants.ErrorCodeNotFound,
				Message: "Resource not found",
			})
		}

		// Handle Echo HTTP errors (if they bubble up)
		if he, ok := err.(*echo.HTTPError); ok {
			return c.JSON(he.Code, api.ApiError{
				Code:    constants.ErrorCodeBadRequest,
				Message: he.Message.(string),
			})
		}

		// Default to internal server error for unknown errors
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
	}
}
