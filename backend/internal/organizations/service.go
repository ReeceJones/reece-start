package organizations

import (
	"errors"
	"log"

	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

// Service request/response types
type CreateOrganizationParams struct {
	Name   string
	Description string
	UserID uint
}

type CreateOrganizationServiceRequest struct {
	Params CreateOrganizationParams
	Tx     *gorm.DB
}

type GetOrganizationsByUserIDServiceRequest struct {
	UserID uint
	Tx     *gorm.DB
}

type GetOrganizationByIDServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
}

type UpdateOrganizationParams struct {
	OrganizationID uint
	Name           *string
	Description    *string
}

type UpdateOrganizationServiceRequest struct {
	Params UpdateOrganizationParams
	Tx     *gorm.DB
}

type DeleteOrganizationServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
}

type CheckUserOrganizationAccessServiceRequest struct {
	UserID         uint
	OrganizationID uint
	Tx             *gorm.DB
}

type CheckUserOrganizationAdminAccessServiceRequest struct {
	UserID         uint
	OrganizationID uint
	Tx             *gorm.DB
}

// Service functions
func createOrganization(request CreateOrganizationServiceRequest) (*models.Organization, error) {
	tx := request.Tx
	params := request.Params

	log.Printf("Creating organization: %+v", params)

	// Create the organization
	organization := &models.Organization{
		Name: params.Name,
		Description: params.Description,
	}

	err := tx.Create(&organization).Error
	if err != nil {
		return nil, err
	}

	// Create organization membership for the user who created it (as admin)
	membership := &models.OrganizationMembership{
		UserID:         params.UserID,
		OrganizationID: organization.ID,
		Role:           string(constants.OrganizationRoleAdmin),
	}

	err = tx.Create(&membership).Error
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func getOrganizationsByUserID(request GetOrganizationsByUserIDServiceRequest) ([]models.Organization, error) {
	tx := request.Tx
	userID := request.UserID

	var organizations []models.Organization
	err := tx.Model(&models.Organization{}).
		Joins("INNER JOIN organization_memberships ON organizations.id = organization_memberships.organization_id").
		Where("organization_memberships.user_id = ? AND organization_memberships.deleted_at IS NULL", userID).
		Find(&organizations).Error
	if err != nil {
		return nil, err
	}


	return organizations, nil
}

func getOrganizationByID(request GetOrganizationByIDServiceRequest) (*models.Organization, error) {
	tx := request.Tx
	organizationID := request.OrganizationID

	var organization models.Organization
	err := tx.First(&organization, organizationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}

	return &organization, nil
}

func updateOrganization(request UpdateOrganizationServiceRequest) (*models.Organization, error) {
	tx := request.Tx
	params := request.Params

	// Get the existing organization
	organization, err := getOrganizationByID(GetOrganizationByIDServiceRequest{
		OrganizationID: params.OrganizationID,
		Tx:             tx,
	})
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if params.Name != nil {
		organization.Name = *params.Name
	}

	if params.Description != nil {
		organization.Description = *params.Description
	}

	// Save the updated organization
	err = tx.Save(organization).Error
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func deleteOrganization(request DeleteOrganizationServiceRequest) error {
	tx := request.Tx
	organizationID := request.OrganizationID

	// Delete organization memberships first (due to foreign key constraints)
	err := tx.Where("organization_id = ?", organizationID).Delete(&models.OrganizationMembership{}).Error
	if err != nil {
		return err
	}

	// Delete the organization
	err = tx.Delete(&models.Organization{}, organizationID).Error
	if err != nil {
		return err
	}

	return nil
}

func checkUserOrganizationAccess(request CheckUserOrganizationAccessServiceRequest) (bool, error) {
	tx := request.Tx
	userID := request.UserID
	organizationID := request.OrganizationID

	var count int64
	err := tx.Model(&models.OrganizationMembership{}).
		Where("user_id = ? AND organization_id = ?", userID, organizationID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func checkUserOrganizationAdminAccess(request CheckUserOrganizationAdminAccessServiceRequest) (bool, error) {
	tx := request.Tx
	userID := request.UserID
	organizationID := request.OrganizationID

	var membership models.OrganizationMembership
	err := tx.Where("user_id = ? AND organization_id = ?", userID, organizationID).
		First(&membership).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return membership.Role == string(constants.OrganizationRoleAdmin), nil
}