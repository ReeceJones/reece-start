package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Content-Type middleware
func ContentTypeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") == "" {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{
				"error": "missing_content_type",
			})
		}

		if c.Request().Header.Get("Content-Type") != "application/json" {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{
				"error": "invalid_content_type",
			})
		}

		return next(c)
	}
}