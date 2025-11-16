package organizations

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/internal/stripe"
	"reece.start/internal/users"
	"reece.start/internal/utils"
)

// Service functions
func createOrganization(request CreateOrganizationServiceRequest) (*OrganizationDto, error) {
	tx := request.Tx
	params := request.Params

	slog.Info("Creating organization", "params", params)

	localCurrency := utils.GetCurrencyForCountry(params.Address.Country)

	// Create the organization
	organization := &models.Organization{
		Name:                params.Name,
		Description:         params.Description,
		ContactEmail:        params.ContactEmail,
		ContactPhone:        params.ContactPhone,
		ContactPhoneCountry: params.ContactPhoneCountry,
		Currency:            string(localCurrency),
		Locale:              params.Locale,
		Address: models.Address{
			Line1:           params.Address.Line1,
			Line2:           params.Address.Line2,
			City:            params.Address.City,
			StateOrProvince: params.Address.StateOrProvince,
			Zip:             params.Address.Zip,
			Country:         params.Address.Country,
		},
		OnboardingStatus: string(constants.OnboardingStatusPending),
		Stripe: models.OrganizationStripeAccount{
			OnboardingStatus: string(constants.StripeOnboardingStatusPending),
		},
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

		slog.Info("Uploading logo for organization", "organizationID", organization.ID, "length", len(decodedImage))

		// Get the mime type from the image
		mimeType := http.DetectContentType(decodedImage)

		slog.Info("Detected logo mime type", "mimeType", mimeType)

		objectName := fmt.Sprintf("%d", organization.ID)

		// upload the image to minio
		_, err = request.MinioClient.PutObject(context.Background(), string(constants.StorageBucketOrganizationLogos), objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{
			ContentType: mimeType,
		})

		if err != nil {
			return nil, err
		}

		slog.Info("Uploaded logo for organization", "organizationID", organization.ID)

		organization.LogoFileStorageKey = objectName
	}

	// Create the Stripe connect account
	account, err := stripe.CreateStripeConnectAccount(stripe.CreateStripeAccountServiceRequest{
		Context:      request.Context,
		Config:       request.Config,
		StripeClient: request.StripeClient,
		Params: stripe.CreateStripeAccountParams{
			OrganizationID:  organization.ID,
			DisplayName:     organization.Name,
			Type:            stripeGo.AccountBusinessType(params.EntityType),
			ContactEmail:    params.ContactEmail,
			ContactPhone:    params.ContactPhone,
			Currency:        localCurrency,
			Locale:          params.Locale,
			ResidingCountry: params.Address.Country,
			Address: stripe.Address{
				Line1:           params.Address.Line1,
				Line2:           params.Address.Line2,
				City:            params.Address.City,
				StateOrProvince: params.Address.StateOrProvince,
				Zip:             params.Address.Zip,
				Country:         params.Address.Country,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	slog.Info("Created Stripe connect account", "accountID", account.ID, "organizationID", organization.ID)

	// Update the stripe information on the organization
	err = updateOrganizationStripeInformation(UpdateOrganizationStripeInformationServiceRequest{
		Organization:  organization,
		StripeAccount: *account,
	})
	if err != nil {
		return nil, err
	}

	err = updateOnboardingStatus(organization)
	if err != nil {
		return nil, err
	}

	// Save the updated organization
	err = tx.Save(organization).Error
	if err != nil {
		return nil, err
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

		slog.Info("Uploading logo for organization", "organizationID", organization.ID, "length", len(decodedImage))

		// Get the mime type from the image
		mimeType := http.DetectContentType(decodedImage)

		slog.Info("Detected logo mime type", "mimeType", mimeType)

		objectName := fmt.Sprintf("%d", organization.ID)

		// upload the image to minio
		_, err = minioClient.PutObject(context.Background(), string(constants.StorageBucketOrganizationLogos), objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{
			ContentType: mimeType,
		})

		if err != nil {
			return nil, err
		}

		slog.Info("Updated logo for organization", "organizationID", organization.ID)

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

	slog.Info("Presigned URL for organization logo", "organizationID", organizationID, "objectName", objectName, "url", presignedUrl.String())

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
			return nil, api.ErrMembershipNotFound
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

		// Also update the user's token revocation
		err = tx.Model(&models.User{}).Where("id = ?", membership.User.ID).Update("revocation_can_refresh", true).Update("revocation_last_valid_issued_at", time.Now()).Error
		if err != nil {
			return nil, err
		}
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

// Organization Invitation Service Functions
func createOrganizationInvitation(request CreateOrganizationInvitationServiceRequest) (*OrganizationInvitationDto, error) {
	tx := request.Tx
	params := request.Params
	riverClient := request.RiverClient

	// Check if there's already a pending invitation for this email and organization
	var existingInvitation models.OrganizationInvitation
	err := tx.Where("email = ? AND organization_id = ? AND status = ?", params.Email, params.OrganizationID, string(constants.OrganizationInvitationStatusPending)).
		First(&existingInvitation).Error

	if err == nil {
		return nil, api.ErrInvitationAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create the organization invitation
	invitation := &models.OrganizationInvitation{
		Email:          params.Email,
		OrganizationID: params.OrganizationID,
		InvitingUserID: params.InvitingUserID,
		Role:           params.Role,
		Status:         string(constants.OrganizationInvitationStatusPending),
	}

	err = tx.Create(&invitation).Error
	if err != nil {
		return nil, err
	}

	// Enqueue background job to send invitation email
	sqlTx := utils.GetGormSQLTx(tx)
	_, err = riverClient.InsertTx(tx.Statement.Context, sqlTx, OrganizationInvitationEmailJobArgs{
		InvitationId: invitation.ID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue invitation email job: %w", err)
	}

	slog.Info("Created organization invitation and enqueued email job", "invitationID", invitation.ID)

	return &OrganizationInvitationDto{
		Invitation:   invitation,
		Organization: nil, // Organization data not needed for creation
		InvitingUser: nil, // Inviting user data not needed for creation
	}, nil
}

func getOrganizationInvitations(request GetOrganizationInvitationsServiceRequest) ([]*OrganizationInvitationDto, error) {
	tx := request.Tx
	organizationID := request.OrganizationID

	var invitations []models.OrganizationInvitation
	err := tx.Where("organization_id = ? AND status = ?", organizationID, string(constants.OrganizationInvitationStatusPending)).Find(&invitations).Error
	if err != nil {
		return nil, err
	}

	var invitationDtos []*OrganizationInvitationDto
	for _, invitation := range invitations {
		invitationDtos = append(invitationDtos, &OrganizationInvitationDto{
			Invitation:   &invitation,
			Organization: nil, // Organization data not needed for list endpoint
			InvitingUser: nil, // Inviting user data not needed for list endpoint
		})
	}

	return invitationDtos, nil
}

func getOrganizationInvitationByID(request GetOrganizationInvitationByIDServiceRequest) (*OrganizationInvitationDto, error) {
	tx := request.Tx
	invitationID := request.InvitationID
	minioClient := request.MinioClient

	var invitation models.OrganizationInvitation
	err := tx.First(&invitation, invitationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrInvitationNotFound
		}
		return nil, err
	}

	// Get the organization data with logo distribution URL
	organizationDto, err := getOrganizationByID(GetOrganizationByIDServiceRequest{
		OrganizationID: invitation.OrganizationID,
		Tx:             tx,
		MinioClient:    minioClient,
	})
	if err != nil {
		return nil, err
	}

	// Get the inviting user data
	var invitingUser models.User
	err = tx.First(&invitingUser, invitation.InvitingUserID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrUserNotFound
		}
		return nil, err
	}

	// Get user logo distribution URL
	userLogoDistributionUrl, err := users.GetUserLogoDistributionUrl(users.GetUserLogoDistributionUrlServiceRequest{
		UserID:      invitingUser.ID,
		Tx:          tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	invitingUserDto := &InvitingUserDto{
		User:                    &invitingUser,
		UserLogoDistributionUrl: userLogoDistributionUrl,
	}

	return &OrganizationInvitationDto{
		Invitation:   &invitation,
		Organization: organizationDto,
		InvitingUser: invitingUserDto,
	}, nil
}

func deleteOrganizationInvitation(request DeleteOrganizationInvitationServiceRequest) error {
	tx := request.Tx
	invitationID := request.InvitationID

	// Mark the invitation as revoked
	err := tx.Model(&models.OrganizationInvitation{}).Where("id = ?", invitationID).Update("status", string(constants.OrganizationInvitationStatusRevoked)).Error
	if err != nil {
		return err
	}

	return nil
}

func acceptOrganizationInvitation(request AcceptOrganizationInvitationServiceRequest) (*OrganizationInvitationDto, error) {
	tx := request.Tx
	invitationID := request.InvitationID
	userID := request.UserID

	// First, get the invitation and verify it's pending
	var invitation models.OrganizationInvitation
	err := tx.First(&invitation, invitationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrInvitationNotFound
		}
		return nil, err
	}

	// Check if invitation is still pending
	if invitation.Status != string(constants.OrganizationInvitationStatusPending) {
		return nil, api.ErrInvitationNotPending
	}

	// Get the user to check their email matches the invitation
	var user models.User
	err = tx.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrUserNotFound
		}
		return nil, err
	}

	// Verify the user's email matches the invitation email
	if user.Email != invitation.Email {
		return nil, api.ErrInvitationEmailMismatch
	}

	// Check if user is already a member of the organization
	var existingMembership models.OrganizationMembership
	err = tx.Where("user_id = ? AND organization_id = ?", userID, invitation.OrganizationID).
		First(&existingMembership).Error

	if err == nil {
		return nil, api.ErrUserAlreadyMember
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create organization membership
	membership := &models.OrganizationMembership{
		UserID:         userID,
		OrganizationID: invitation.OrganizationID,
		Role:           invitation.Role,
	}

	err = tx.Create(&membership).Error
	if err != nil {
		return nil, err
	}

	// Update invitation status to accepted
	err = tx.Model(&invitation).Update("status", string(constants.OrganizationInvitationStatusAccepted)).Error
	if err != nil {
		return nil, err
	}

	// Get the updated invitation with organization and user data
	return getOrganizationInvitationByID(GetOrganizationInvitationByIDServiceRequest{
		InvitationID: invitationID,
		Tx:           tx,
		MinioClient:  request.MinioClient,
	})
}

func declineOrganizationInvitation(request DeclineOrganizationInvitationServiceRequest) (*OrganizationInvitationDto, error) {
	tx := request.Tx
	invitationID := request.InvitationID
	userID := request.UserID

	// First, get the invitation and verify it's pending
	var invitation models.OrganizationInvitation
	err := tx.First(&invitation, invitationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrInvitationNotFound
		}
		return nil, err
	}

	// Check if invitation is still pending
	if invitation.Status != string(constants.OrganizationInvitationStatusPending) {
		return nil, api.ErrInvitationNotPending
	}

	// Get the user to check their email matches the invitation
	var user models.User
	err = tx.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrUserNotFound
		}
		return nil, err
	}

	// Verify the user's email matches the invitation email
	if user.Email != invitation.Email {
		return nil, api.ErrInvitationEmailMismatch
	}

	// Update invitation status to declined
	err = tx.Model(&invitation).Update("status", string(constants.OrganizationInvitationStatusDeclined)).Error
	if err != nil {
		return nil, err
	}

	// Get the updated invitation with organization and user data
	return getOrganizationInvitationByID(GetOrganizationInvitationByIDServiceRequest{
		InvitationID: invitationID,
		Tx:           tx,
		MinioClient:  request.MinioClient,
	})
}

func updateOrganizationStripeInformation(request UpdateOrganizationStripeInformationServiceRequest) error {
	organization := request.Organization
	stripeAccount := request.StripeAccount
	requirements := stripeAccount.Requirements

	organization.Stripe.AccountID = stripeAccount.ID
	organization.Stripe.AutomaticIndirectTaxStatus = string(stripeAccount.Configuration.Customer.Capabilities.AutomaticIndirectTax.Status)
	organization.Stripe.CardPaymentsStatus = string(stripeAccount.Configuration.Merchant.Capabilities.CardPayments.Status)
	organization.Stripe.StripeBalancePayoutsStatus = string(stripeAccount.Configuration.Recipient.Capabilities.StripeBalance.Payouts.Status)
	organization.Stripe.StripeBalanceTransfersStatus = string(stripeAccount.Configuration.Recipient.Capabilities.StripeBalance.StripeTransfers.Status)

	if requirements != nil {
		organization.Stripe.HasPendingRequirements = len(requirements.Entries) > 0
	}

	organization.Stripe.OnboardingStatus = string(utils.DetermineStripeOnboardingStatus(organization))

	return nil
}

func updateOnboardingStatus(organization *models.Organization) error {
	if organization.Stripe.OnboardingStatus == string(constants.StripeOnboardingStatusCompleted) {
		organization.OnboardingStatus = string(constants.OnboardingStatusCompleted)
		return nil
	}

	organization.OnboardingStatus = string(constants.OnboardingStatusInProgress)

	return nil
}

// moved to utils.DetermineStripeOnboardingStatus

func createStripeOnboardingLink(request CreateStripeOnboardingLinkServiceRequest) (*stripeGo.V2CoreAccountLink, error) {
	tx := request.Db
	params := request.Params
	stripeClient := request.StripeClient
	context := request.Context
	config := request.Config

	var organization models.Organization
	if err := tx.First(&organization, params.OrganizationID).Error; err != nil {
		return nil, err
	}

	if organization.Stripe.AccountID == "" {
		return nil, fmt.Errorf("organization %d does not have a Stripe account", organization.ID)
	}

	refreshUrl := fmt.Sprintf("%s/app/%d/stripe-onboarding", config.FrontendUrl, params.OrganizationID)
	returnUrl := fmt.Sprintf("%s/app/%d", config.FrontendUrl, params.OrganizationID)

	link, err := stripe.CreateOnboardingLink(stripe.CreateOnboardingLinkServiceRequest{
		Context:      context,
		StripeClient: stripeClient,
		Params: stripe.CreateOnboardingLinkParams{
			AccountID:  organization.Stripe.AccountID,
			RefreshURL: refreshUrl,
			ReturnURL:  returnUrl,
		},
	})
	if err != nil {
		return nil, err
	}

	return link, nil
}

func createStripeDashboardLink(request CreateStripeDashboardLinkServiceRequest) (string, error) {
	tx := request.Db
	params := request.Params

	var organization models.Organization
	if err := tx.First(&organization, params.OrganizationID).Error; err != nil {
		return "", err
	}

	if organization.Stripe.AccountID == "" {
		return "", fmt.Errorf("organization %d does not have a Stripe account", organization.ID)
	}

	// For Standard Connect accounts, construct the dashboard URL directly
	dashboardURL := fmt.Sprintf("https://dashboard.stripe.com/b/%s", organization.Stripe.AccountID)

	return dashboardURL, nil
}
