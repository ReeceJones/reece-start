package middleware

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/authentication"
	"reece.start/internal/configuration"
	"reece.start/internal/constants"
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
				return api.ErrInvalidToken
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

func GetRoleFromJWT(c echo.Context) (constants.UserRole, error) {
	claims := c.Get("claims").(*authentication.JwtClaims)
	if claims.Role == nil {
		return constants.UserRole(""), errors.New("role is not set")
	}
	return *claims.Role, nil
}

func GetScopesFromJWT(c echo.Context) ([]constants.UserScope, error) {
	claims := c.Get("claims").(*authentication.JwtClaims)
	if claims.Scopes == nil {
		return []constants.UserScope{}, errors.New("scopes are not set")
	}
	return *claims.Scopes, nil
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
	if err == nil && tokenString != "" {
		return tokenString, nil
	}
	return getTokenFromAuthorizationHeader(c)
}

func getTokenFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("app-session-token")
	if err != nil {
		// Return error without writing response - let getTokenFromRequest handle fallback
		return "", err
	}
	return cookie.Value, nil
}

func getTokenFromAuthorizationHeader(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", api.ErrMissingAuthorizationHeader
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", api.ErrInvalidAuthorizationFormat
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	return tokenString, nil
}
