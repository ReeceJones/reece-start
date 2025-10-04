package access

import (
	"slices"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/authentication"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
)

type HasOrganizationAccessParams struct {
	OrganizationID uint
	Scopes []constants.UserScope
}

func HasOrganizationAccess(c echo.Context, params HasOrganizationAccessParams) error {
	// Checks scopes in the JWT. This has multiple implications:
	// 1. If the user is updated then their token needs to be re-issued
	// 2. If a user is deleted or their role is downgraded, then their token needs to be re-issued or revoked
	// For both of the above situations, this will happen higher in the stack
	scopes, err := middleware.GetScopesFromJWT(c)
	if err != nil {
		return err
	}

	for _, scope := range params.Scopes {
		if !slices.Contains(scopes, scope) {
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