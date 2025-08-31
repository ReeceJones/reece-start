package access

import (
	"slices"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
	"reece.start/internal/models"
)

// func HasAccessToOrganization(c echo.Context, organizationID uint) error {
// 	userID, err := middleware.GetUserIDFromJWT(c)
// 	if err != nil {
// 		return err
// 	}

// 	db := middleware.GetDB(c)

// 	var count int64
// 	err = db.Model(&models.OrganizationMembership{}).
// 		Where("user_id = ? AND organization_id = ?", userID, organizationID).
// 		Count(&count).Error

// 	if err != nil {
// 		return api.ErrForbiddenNoAccess
// 	}

// 	if count == 0 {
// 		return api.ErrForbiddenNoAccess
// 	}

// 	return nil
// }

// func HasAdminAccessToOrganization(c echo.Context, organizationID uint) error {
// 	userID, err := middleware.GetUserIDFromJWT(c)
// 	if err != nil {
// 		return err
// 	}

// 	db := middleware.GetDB(c)

// 	var count int64
// 	err = db.Model(&models.OrganizationMembership{}).
// 		Where("user_id = ? AND organization_id = ?", userID, organizationID).
// 		Count(&count).Error
// 	if err != nil {
// 		return api.ErrForbiddenNoAccess
// 	}

// 	if count == 0 {
// 		return api.ErrForbiddenNoAdminAccess
// 	}

// 	return nil
// }

type HasOrganizationAccessParams struct {
	OrganizationID uint
	Scopes []constants.OrganizationScope
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
	scopes := constants.OrganizationRoleToScopes[constants.OrganizationRole(membership.Role)]
	for _, scope := range params.Scopes {
		if !slices.Contains(scopes, scope) {
			return api.ErrForbiddenNoAccess
		}
	}

	return nil
}