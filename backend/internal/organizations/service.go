package organizations

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

// Service request/response types
type CreateOrganizationParams struct {
	Name   string
	Description string
	UserID uint
	Logo   string
}

type CreateOrganizationServiceRequest struct {
	Params CreateOrganizationParams
	Tx     *gorm.DB
	MinioClient *minio.Client
}

type GetOrganizationsByUserIDServiceRequest struct {
	UserID uint
	Tx     *gorm.DB
	MinioClient *minio.Client
}

type GetOrganizationByIDServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
	MinioClient    *minio.Client
}

type UpdateOrganizationParams struct {
	OrganizationID uint
	Name           *string
	Description    *string
	Logo           *string
}

type UpdateOrganizationServiceRequest struct {
	Params UpdateOrganizationParams
	Tx     *gorm.DB
	MinioClient *minio.Client
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

type GetOrganizationLogoDistributionUrlServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
	MinioClient    *minio.Client
}

type OrganizationDto struct {
	Organization        *models.Organization
	LogoDistributionUrl string
}

// Service functions
func createOrganization(request CreateOrganizationServiceRequest) (*OrganizationDto, error) {
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

	// Handle logo upload if provided
	if params.Logo != "" {
		// decode the image from base64 to a binary file
		decodedImage, err := base64.StdEncoding.DecodeString(params.Logo)
		if err != nil {
			return nil, err
		}

		log.Printf("Uploading logo for organization %d of length %d\n", organization.ID, len(decodedImage))

		// Get the mime type from the image
		mimeType := http.DetectContentType(decodedImage)

		log.Printf("Detected logo mime type: %s\n", mimeType)

		objectName := fmt.Sprintf("%d", organization.ID)

		// upload the image to minio
		_, err = request.MinioClient.PutObject(context.Background(), string(constants.StorageBucketOrganizationLogos), objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{
			ContentType: mimeType,
		})

		if err != nil {
			return nil, err
		}

		log.Printf("Uploaded logo for organization %d\n", organization.ID)

		organization.LogoFileStorageKey = objectName

		// Save the updated organization with logo key
		err = tx.Save(organization).Error
		if err != nil {
			return nil, err
		}
	}

	// Get the logo distribution URL for the new organization
	logoDistributionUrl, err := getOrganizationLogoDistributionUrl(GetOrganizationLogoDistributionUrlServiceRequest{
		OrganizationID: organization.ID,
		Tx:             tx,
		MinioClient:    request.MinioClient,
	})
	if err != nil {
		return nil, err
	}

	return &OrganizationDto{
		Organization:        organization,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}

func getOrganizationsByUserID(request GetOrganizationsByUserIDServiceRequest) ([]*OrganizationDto, error) {
	tx := request.Tx
	userID := request.UserID
	minioClient := request.MinioClient

	var organizations []models.Organization
	err := tx.Model(&models.Organization{}).
		Joins("INNER JOIN organization_memberships ON organizations.id = organization_memberships.organization_id").
		Where("organization_memberships.user_id = ? AND organization_memberships.deleted_at IS NULL", userID).
		Find(&organizations).Error
	if err != nil {
		return nil, err
	}

	// Convert to DTOs with logo distribution URLs
	var orgDtos []*OrganizationDto
	for _, org := range organizations {
		logoDistributionUrl, err := getOrganizationLogoDistributionUrl(GetOrganizationLogoDistributionUrlServiceRequest{
			OrganizationID: org.ID,
			Tx:             tx,
			MinioClient:    minioClient,
		})
		if err != nil {
			return nil, err
		}

		orgDtos = append(orgDtos, &OrganizationDto{
			Organization:        &org,
			LogoDistributionUrl: logoDistributionUrl,
		})
	}

	return orgDtos, nil
}

func getOrganizationByID(request GetOrganizationByIDServiceRequest) (*OrganizationDto, error) {
	tx := request.Tx
	organizationID := request.OrganizationID
	minioClient := request.MinioClient

	var organization models.Organization
	err := tx.First(&organization, organizationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}

	// Get the logo distribution URL for the organization
	logoDistributionUrl, err := getOrganizationLogoDistributionUrl(GetOrganizationLogoDistributionUrlServiceRequest{
		OrganizationID: organization.ID,
		Tx:             tx,
		MinioClient:    minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &OrganizationDto{
		Organization:        &organization,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}

func updateOrganization(request UpdateOrganizationServiceRequest) (*OrganizationDto, error) {
	tx := request.Tx
	params := request.Params
	minioClient := request.MinioClient

	// Get the existing organization
	orgDto, err := getOrganizationByID(GetOrganizationByIDServiceRequest{
		OrganizationID: params.OrganizationID,
		Tx:             tx,
		MinioClient:    minioClient,
	})
	if err != nil {
		return nil, err
	}

	organization := orgDto.Organization

	// Update fields if provided
	if params.Name != nil {
		organization.Name = *params.Name
	}

	if params.Description != nil {
		organization.Description = *params.Description
	}

	if params.Logo != nil && *params.Logo != "" {
		// decode the image from base64 to a binary file
		decodedImage, err := base64.StdEncoding.DecodeString(*params.Logo)
		if err != nil {
			return nil, err
		}

		log.Printf("Uploading logo for organization %d of length %d\n", organization.ID, len(decodedImage))

		// Get the mime type from the image
		mimeType := http.DetectContentType(decodedImage)

		log.Printf("Detected logo mime type: %s\n", mimeType)

		objectName := fmt.Sprintf("%d", organization.ID)

		// upload the image to minio
		_, err = minioClient.PutObject(context.Background(), string(constants.StorageBucketOrganizationLogos), objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{
			ContentType: mimeType,
		})

		if err != nil {
			return nil, err
		}

		log.Printf("Updated logo for organization %d\n", organization.ID)

		organization.LogoFileStorageKey = objectName
	}

	// Save the updated organization
	err = tx.Save(organization).Error
	if err != nil {
		return nil, err
	}

	// Get the logo distribution URL for the updated organization
	logoDistributionUrl, err := getOrganizationLogoDistributionUrl(GetOrganizationLogoDistributionUrlServiceRequest{
		OrganizationID: organization.ID,
		Tx:             tx,
		MinioClient:    minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &OrganizationDto{
		Organization:        organization,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
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

func getOrganizationLogoDistributionUrl(request GetOrganizationLogoDistributionUrlServiceRequest) (string, error) {
	tx := request.Tx
	minioClient := request.MinioClient
	organizationID := request.OrganizationID

	var organization models.Organization
	err := tx.First(&organization, organizationID).Error
	if err != nil {
		return "", err
	}

	objectName := organization.LogoFileStorageKey
	if objectName == "" {
		return "", nil
	}

	presignedUrl, err := minioClient.PresignedGetObject(context.Background(), string(constants.StorageBucketOrganizationLogos), objectName, time.Hour*24, url.Values{})
	if err != nil {
		return "", err
	}

	return presignedUrl.String(), nil
}