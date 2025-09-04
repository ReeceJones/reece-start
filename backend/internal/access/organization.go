package access

import (
	"slices"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
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