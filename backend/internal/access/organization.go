package access

import (
	"slices"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
)

type HasOrganizationAccessParams struct {
	OrganizationID uuid.UUID
	Scopes         []constants.UserScope
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

	if !hasScopes(params.Scopes, scopes) {
		return api.ErrForbiddenNoAccess
	}

	return nil
}

// HasAdminAccess checks if the user has admin access based on their role and scopes
func HasAdminAccess(c echo.Context, scopes []constants.UserScope) error {
	role, err := middleware.GetRoleFromJWT(c)
	if err != nil {
		return err
	}

	// Check if user has admin role
	if !hasRole(role, constants.UserRoleAdmin) {
		return api.ErrForbiddenNoAdminAccess
	}

	grantedScopes, err := middleware.GetScopesFromJWT(c)
	if err != nil {
		return err
	}

	// Check if user has the required scopes
	if len(grantedScopes) == 0 {
		return api.ErrForbiddenNoAdminAccess
	}

	if !hasScopes(scopes, grantedScopes) {
		return api.ErrForbiddenNoAdminAccess
	}

	return nil
}

func hasScopes(scopes []constants.UserScope, grantedScopes []constants.UserScope) bool {
	for _, scope := range scopes {
		if !slices.Contains(grantedScopes, scope) {
			return false
		}
	}
	return true
}

func hasRole(role constants.UserRole, grantedRole constants.UserRole) bool {
	return role == grantedRole
}
