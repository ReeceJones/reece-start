package api

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ParseOrganizationIDFromParams parses organization ID from URL parameter
func ParseOrganizationIDFromParams(c echo.Context) (uint, error) {
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

// ParseMembershipIDFromParams parses membership ID from URL parameter
func ParseMembershipIDFromParams(c echo.Context) (uint, error) {
	paramMembershipID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, ErrInvalidMembershipID
	}
	return uint(paramMembershipID), nil
}

// ParseOrganizationInvitationIDFromParams parses invitation ID from URL parameter
func ParseOrganizationInvitationIDFromParams(c echo.Context) (uuid.UUID, error) {
	paramInvitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, ErrInvalidInvitationID
	}
	return paramInvitationID, nil
}
