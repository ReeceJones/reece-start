package middleware

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Common parsing error variables for HTTP layer
var (
	ErrInvalidUserToken      = errors.New("invalid user token")
	ErrInvalidOrganizationID = errors.New("invalid organization ID")
	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrInvalidMembershipID   = errors.New("invalid membership ID")
	ErrInvalidInvitationID   = errors.New("invalid invitation ID")
)

// Helper functions for controllers to use instead of direct parsing

// HandleJWTError checks JWT validation and returns appropriate error
func HandleJWTError(c echo.Context) (uint, error) {
	userID, err := GetUserIDFromJWT(c)
	if err != nil {
		return 0, ErrInvalidUserToken
	}
	return userID, nil
}

// ParseOrganizationID parses organization ID from URL parameter
func ParseOrganizationID(c echo.Context) (uint, error) {
	paramOrgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, ErrInvalidOrganizationID
	}
	return uint(paramOrgID), nil
}

// ParseOrganizationIDFromString parses organization ID from string
func ParseOrganizationIDFromString(idStr string) (uint, error) {
	paramOrgID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, ErrInvalidOrganizationID
	}
	return uint(paramOrgID), nil
}

// ParseUserIDFromString parses user ID from string
func ParseUserIDFromString(idStr string) (uint, error) {
	paramUserID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, ErrInvalidUserID
	}
	return uint(paramUserID), nil
}

// ParseMembershipID parses membership ID from URL parameter
func ParseMembershipID(c echo.Context) (uint, error) {
	paramMembershipID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, ErrInvalidMembershipID
	}
	return uint(paramMembershipID), nil
}

// ParseOrganizationInvitationID parses invitation ID from URL parameter
func ParseOrganizationInvitationID(c echo.Context) (uuid.UUID, error) {
	paramInvitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, ErrInvalidInvitationID
	}
	return paramInvitationID, nil
}
