package access

import (
	"slices"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/authentication"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
	"reece.start/internal/models"
)

type HasOrganizationAccessParams struct {
	OrganizationID uint
	Scopes []constants.UserScope
}

func HasOrganizationAccess(c echo.Context, params HasOrganizationAccessParams) error {
	db := middleware.GetDB(c)
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}
	
	// fetch the organization membership to get the role
	var membership models.OrganizationMembership
	err = db.Where("user_id = ? AND organization_id = ?", userID, params.OrganizationID).
		First(&membership).Error
	if err != nil {
		return api.ErrForbiddenNoAccess
	}
	
	// check if the role has the required scopes
	organization_scopes := constants.OrganizationRoleToScopes[constants.OrganizationRole(membership.Role)]
	for _, scope := range params.Scopes {
		if !slices.Contains(organization_scopes, scope) {
			return api.ErrForbiddenNoAccess
		}
	}

	return nil
}

// HasAdminAccess checks if the user has admin access based on their role and scopes
func HasAdminAccess(c echo.Context, scopes []constants.UserScope) error {
	claims := c.Get("claims").(*authentication.JwtClaims)
	
	// Check if user has admin role
	if claims.Role == nil || *claims.Role != constants.UserRoleAdmin {
		return api.ErrForbiddenNoAdminAccess
	}
	
	// Check if user has the required scopes
	if claims.Scopes == nil {
		return api.ErrForbiddenNoAdminAccess
	}
	
	for _, scope := range scopes {
		if !slices.Contains(*claims.Scopes, scope) {
			return api.ErrForbiddenNoAdminAccess
		}
	}

	return nil
}