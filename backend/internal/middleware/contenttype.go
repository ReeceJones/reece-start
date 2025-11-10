package middleware

import (
	"log"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
)

var allowedContentTypes = []string{"application/json", "application/json; charset=utf-8"}

// Content-Type middleware
func ContentTypeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") == "" {
			return c.JSON(http.StatusUnsupportedMediaType, api.ApiError{
				Message: "missing_content_type",
			})
		}

		contentType := c.Request().Header.Get("Content-Type")
		if !slices.Contains(allowedContentTypes, contentType) {
			log.Printf("Invalid content type: %s", c.Request().Header.Get("Content-Type"))
			return c.JSON(http.StatusUnsupportedMediaType, api.ApiError{
				Message: "invalid_content_type",
			})
		}

		return next(c)
	}
}
