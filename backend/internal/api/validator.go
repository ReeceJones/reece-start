package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"reece.start/internal/constants"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

// https://robinverton.de/blog/go-echo-generic-validation/
func Validated[T any](h func(c echo.Context, t T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		var t T
		if err := c.Bind(&t); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Code:    constants.ErrorCodeInvalidRequest,
				Message: err.Error(),
			})
		}

		if err := validate.Struct(t); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, ApiError{
				Code:    constants.ErrorCodeValidationFailed,
				Message: err.Error(),
			})
		}

		return h(c, t)
	}
}