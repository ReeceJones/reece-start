package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"reece.start/internal/authentication"
	"reece.start/internal/configuration"
)

// JWT Authentication middleware
func JwtAuthMiddleware(config *configuration.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := getTokenFromRequest(c)

			if err != nil {
				return err
			}

			// Validate the token
			claims, err := authentication.ValidateJWT(config, tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid_token",
				})
			}

			// Store claims in context
			c.Set("claims", claims)
			return next(c)
		}
	}
}

// GetUserIDFromJWT extracts the user ID from the JWT claims in the context
func GetUserIDFromJWT(c echo.Context) (uint, error) {
	claims := c.Get("claims").(*authentication.JwtClaims)
	userID, err := strconv.ParseUint(claims.UserId, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}

func GetImpersonatingUserIDFromJWT(c echo.Context) (uint, error) {
	claims := c.Get("claims").(*authentication.JwtClaims)

	if claims.ImpersonatingUserId == nil {
		return 0, errors.New("impersonating user ID is not set")
	}

	impersonatingUserID, err := strconv.ParseUint(*claims.ImpersonatingUserId, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(impersonatingUserID), nil
}

func getTokenFromRequest(c echo.Context) (string, error) {
	tokenString, err := getTokenFromCookie(c)
	if err == nil {
		return tokenString, nil
	}
	return getTokenFromAuthorizationHeader(c)
}

func getTokenFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("app-session-token")
	if err != nil {
		return "", c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing_token_cookie",
		})
	}
	return cookie.Value, nil
}

func getTokenFromAuthorizationHeader(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing_authorization_header",
		})
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid_authorization_format",
		})
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	return tokenString, nil
}
