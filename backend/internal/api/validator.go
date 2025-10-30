package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

// https://robinverton.de/blog/go-echo-generic-validation/
func Validated[T any](h func(c echo.Context, t T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		var t T
		if err := c.Bind(&t); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		if err := validate.Struct(t); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		return h(c, t)
	}
}

// ValidatedWithQuery validates both request body and query parameters
func ValidatedWithQuery[T any, Q any](h func(c echo.Context, t T, q Q) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		var t T
		var q Q

		// Bind and validate request body
		if err := c.Bind(&t); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		if err := validate.Struct(t); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		// Bind and validate query parameters
		if err := c.Bind(&q); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		if err := validate.Struct(q); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		return h(c, t, q)
	}
}

// ValidatedQuery validates only query parameters for GET requests
func ValidatedQuery[Q any](h func(c echo.Context, q Q) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		var q Q

		// Bind and validate query parameters
		if err := c.Bind(&q); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		if err := validate.Struct(q); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Message: err.Error(),
			})
		}

		return h(c, q)
	}
}
