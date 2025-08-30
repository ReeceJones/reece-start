package organizations

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/internal/users"
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

type OrganizationMembershipDto struct {
	Membership          *models.OrganizationMembership
	User                *models.User
	UserLogoDistributionUrl string
	Organization        *models.Organization
}

// Organization Membership Service Types
type CreateOrganizationMembershipParams struct {
	UserID         uint
	OrganizationID uint
	Role           string
}

type CreateOrganizationMembershipServiceRequest struct {
	Params CreateOrganizationMembershipParams
	Tx     *gorm.DB
}

type GetOrganizationMembershipsServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
	MinioClient    *minio.Client
}

type GetOrganizationMembershipByIDServiceRequest struct {
	MembershipID uint
	Tx           *gorm.DB
	MinioClient  *minio.Client
}

type UpdateOrganizationMembershipParams struct {
	MembershipID uint
	Role         *string
}

type UpdateOrganizationMembershipServiceRequest struct {
	Params UpdateOrganizationMembershipParams
	Tx     *gorm.DB
}

type DeleteOrganizationMembershipServiceRequest struct {
	MembershipID uint
	Tx           *gorm.DB
}

// Organization Invitation Service Types
type CreateOrganizationInvitationParams struct {
	Email          string
	Role           string
	OrganizationID uint
	InvitingUserID uint
}

type CreateOrganizationInvitationServiceRequest struct {
	Params      CreateOrganizationInvitationParams
	Tx          *gorm.DB
	RiverClient *river.Client[*sql.Tx]
}

type OrganizationInvitationDto struct {
	Invitation *models.OrganizationInvitation
}

type GetOrganizationInvitationsServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
}

type GetOrganizationInvitationByIDServiceRequest struct {
	InvitationID uint
	Tx           *gorm.DB
}

type DeleteOrganizationInvitationServiceRequest struct {
	InvitationID uint
	Tx           *gorm.DB
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

// Organization Membership Service Functions
func createOrganizationMembership(request CreateOrganizationMembershipServiceRequest) (*OrganizationMembershipDto, error) {
	tx := request.Tx
	params := request.Params

	// Check if membership already exists
	var existingMembership models.OrganizationMembership
	err := tx.Where("user_id = ? AND organization_id = ?", params.UserID, params.OrganizationID).
		First(&existingMembership).Error
	
	if err == nil {
		return nil, api.ErrInvitationAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create the organization membership
	membership := &models.OrganizationMembership{
		UserID:         params.UserID,
		OrganizationID: params.OrganizationID,
		Role:           params.Role,
	}

	err = tx.Create(&membership).Error
	if err != nil {
		return nil, err
	}

	// Reload with preloaded relationships
	err = tx.Preload("User").Preload("Organization").First(&membership, membership.ID).Error
	if err != nil {
		return nil, err
	}

	return &OrganizationMembershipDto{
		Membership:   membership,
		User:         &membership.User,
		Organization: &membership.Organization,
	}, nil
}

func getOrganizationMemberships(request GetOrganizationMembershipsServiceRequest) ([]*OrganizationMembershipDto, error) {
	tx := request.Tx
	organizationID := request.OrganizationID
	minioClient := request.MinioClient

	var memberships []models.OrganizationMembership
	err := tx.Preload("User").Preload("Organization").Where("organization_id = ?", organizationID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	var membershipDtos []*OrganizationMembershipDto
	for _, membership := range memberships {
		// Get user logo distribution URL
		userLogoDistributionUrl, err := users.GetUserLogoDistributionUrl(users.GetUserLogoDistributionUrlServiceRequest{
			UserID:      membership.User.ID,
			Tx:          tx,
			MinioClient: minioClient,
		})
		if err != nil {
			return nil, err
		}

		membershipDtos = append(membershipDtos, &OrganizationMembershipDto{
			Membership:              &membership,
			User:                    &membership.User,
			UserLogoDistributionUrl: userLogoDistributionUrl,
			Organization:            &membership.Organization,
		})
	}

	return membershipDtos, nil
}

func getOrganizationMembershipByID(request GetOrganizationMembershipByIDServiceRequest) (*OrganizationMembershipDto, error) {
	tx := request.Tx
	membershipID := request.MembershipID
	minioClient := request.MinioClient

	var membership models.OrganizationMembership
	err := tx.Preload("User").Preload("Organization").First(&membership, membershipID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("membership not found")
		}
		return nil, err
	}

	// Get user logo distribution URL if MinioClient is provided
	var userLogoDistributionUrl string
	if minioClient != nil {
		userLogoDistributionUrl, err = users.GetUserLogoDistributionUrl(users.GetUserLogoDistributionUrlServiceRequest{
			UserID:      membership.User.ID,
			Tx:          tx,
			MinioClient: minioClient,
		})
		if err != nil {
			return nil, err
		}
	}

	return &OrganizationMembershipDto{
		Membership:              &membership,
		User:                    &membership.User,
		UserLogoDistributionUrl: userLogoDistributionUrl,
		Organization:            &membership.Organization,
	}, nil
}

func updateOrganizationMembership(request UpdateOrganizationMembershipServiceRequest) (*OrganizationMembershipDto, error) {
	tx := request.Tx
	params := request.Params

	// Get the existing membership
	membershipDto, err := getOrganizationMembershipByID(GetOrganizationMembershipByIDServiceRequest{
		MembershipID: params.MembershipID,
		Tx:           tx,
	})
	if err != nil {
		return nil, err
	}

	membership := membershipDto.Membership

	// Update fields if provided
	if params.Role != nil {
		membership.Role = *params.Role
	}

	// Save the updated membership
	err = tx.Save(membership).Error
	if err != nil {
		return nil, err
	}

	// Return updated DTO
	return &OrganizationMembershipDto{
		Membership:   membership,
		User:         membershipDto.User,
		Organization: membershipDto.Organization,
	}, nil
}

func deleteOrganizationMembership(request DeleteOrganizationMembershipServiceRequest) error {
	tx := request.Tx
	membershipID := request.MembershipID

	// Delete the membership
	err := tx.Delete(&models.OrganizationMembership{}, membershipID).Error
	if err != nil {
		return err
	}

	return nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Organization Invitation Service Functions
func createOrganizationInvitation(request CreateOrganizationInvitationServiceRequest) (*OrganizationInvitationDto, error) {
	tx := request.Tx
	params := request.Params
	// riverClient := request.RiverClient

	// Generate a secure random invitation token
	invitationToken, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}

	// Check if there's already a pending invitation for this email and organization
	var existingInvitation models.OrganizationInvitation
	err = tx.Where("email = ? AND organization_id = ?", params.Email, params.OrganizationID).
		First(&existingInvitation).Error
	
	if err == nil {
		return nil, api.ErrInvitationAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create the organization invitation
	invitation := &models.OrganizationInvitation{
		Email:           params.Email,
		InvitationToken: invitationToken,
		OrganizationID:  params.OrganizationID,
		InvitingUserID:  params.InvitingUserID,
		Role:            params.Role,
	}

	err = tx.Create(&invitation).Error
	if err != nil {
		return nil, err
	}

	// Enqueue background job to send invitation email
	// sqlTx := utils.GetGormSQLTx(tx)
	// _, err = riverClient.InsertTx(tx.Statement.Context, sqlTx, OrganizationInvitationEmailJobArgs{
	// 	InvitationId: invitation.ID,
	// }, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to enqueue invitation email job: %w", err)
	// }

	log.Printf("Created organization invitation %d and enqueued email job", invitation.ID)

	return &OrganizationInvitationDto{
		Invitation: invitation,
	}, nil
}

func getOrganizationInvitations(request GetOrganizationInvitationsServiceRequest) ([]*OrganizationInvitationDto, error) {
	tx := request.Tx
	organizationID := request.OrganizationID

	var invitations []models.OrganizationInvitation
	err := tx.Where("organization_id = ?", organizationID).Find(&invitations).Error
	if err != nil {
		return nil, err
	}

	var invitationDtos []*OrganizationInvitationDto
	for _, invitation := range invitations {
		invitationDtos = append(invitationDtos, &OrganizationInvitationDto{
			Invitation: &invitation,
		})
	}

	return invitationDtos, nil
}

func getOrganizationInvitationByID(request GetOrganizationInvitationByIDServiceRequest) (*OrganizationInvitationDto, error) {
	tx := request.Tx
	invitationID := request.InvitationID

	var invitation models.OrganizationInvitation
	err := tx.First(&invitation, invitationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrInvitationNotFound
		}
		return nil, err
	}

	return &OrganizationInvitationDto{
		Invitation: &invitation,
	}, nil
}

func deleteOrganizationInvitation(request DeleteOrganizationInvitationServiceRequest) error {
	tx := request.Tx
	invitationID := request.InvitationID

	// Delete the invitation
	err := tx.Delete(&models.OrganizationInvitation{}, invitationID).Error
	if err != nil {
		return err
	}

	return nil
}
