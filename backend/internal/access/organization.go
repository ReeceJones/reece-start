package access

import (
	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/middleware"
	"reece.start/internal/models"
)

func HasAccessToOrganization(c echo.Context, organizationID uint) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)

	var count int64
	err = db.Model(&models.OrganizationMembership{}).
		Where("user_id = ? AND organization_id = ?", userID, organizationID).
		Count(&count).Error

	if err != nil {
		return api.ErrForbiddenNoAccess
	}

	if count == 0 {
		return api.ErrForbiddenNoAccess
	}

	return nil
}

func HasAdminAccessToOrganization(c echo.Context, organizationID uint) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)

	var count int64
	err = db.Model(&models.OrganizationMembership{}).
		Where("user_id = ? AND organization_id = ?", userID, organizationID).
		Count(&count).Error
	if err != nil {
		return api.ErrForbiddenNoAccess
	}

	if count == 0 {
		return api.ErrForbiddenNoAdminAccess
	}

	return nil
}