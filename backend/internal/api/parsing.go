package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ParseOrganizationIDFromParams parses organization ID from URL parameter
func ParseOrganizationIDFromParams(c echo.Context) (uuid.UUID, error) {
	paramOrgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, ErrInvalidOrganizationID
	}
	return paramOrgID, nil
}

// ParseOrganizationIDFromString parses organization ID from string
func ParseOrganizationIDFromString(idStr string) (uuid.UUID, error) {
	paramOrgID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, ErrInvalidOrganizationID
	}
	return paramOrgID, nil
}

// ParseUserIDFromString parses user ID from string
func ParseUserIDFromString(idStr string) (uuid.UUID, error) {
	paramUserID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, ErrInvalidUserID
	}
	return paramUserID, nil
}

// ParseMembershipIDFromParams parses membership ID from URL parameter
func ParseMembershipIDFromParams(c echo.Context) (uuid.UUID, error) {
	paramMembershipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, ErrInvalidMembershipID
	}
	return paramMembershipID, nil
}

// ParseOrganizationInvitationIDFromParams parses invitation ID from URL parameter
func ParseOrganizationInvitationIDFromParams(c echo.Context) (uuid.UUID, error) {
	paramInvitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, ErrInvalidInvitationID
	}
	return paramInvitationID, nil
}
