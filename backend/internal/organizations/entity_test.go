package organizations

import (
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

func TestMapOrganizationToResponse(t *testing.T) {
	t.Run("maps organization with all fields", func(t *testing.T) {
		organization := &models.Organization{
			Model: gorm.Model{
				ID: 123,
			},
			Name:                "Test Organization",
			Description:         "Test Description",
			ContactEmail:        "contact@example.com",
			ContactPhone:        "1234567890",
			ContactPhoneCountry: "US",
			Locale:              "en-US",
			Currency:            "USD",
			OnboardingStatus:    string(constants.OnboardingStatusCompleted),
			Address: models.Address{
				Line1:           "123 Test St",
				Line2:           "Apt 4",
				City:            "Test City",
				StateOrProvince: "CA",
				Zip:             "12345",
				Country:         "US",
			},
			Stripe: models.OrganizationStripeAccount{
				HasPendingRequirements: false,
				OnboardingStatus:       string(constants.StripeOnboardingStatusCompleted),
			},
		}

		dto := &OrganizationDto{
			Organization:        organization,
			LogoDistributionUrl: "https://example.com/logo.png",
		}

		result := mapOrganizationToResponse(dto)

		assert.Equal(t, strconv.FormatUint(123, 10), result.Id)
		assert.Equal(t, constants.ApiTypeOrganization, result.Type)
		assert.Equal(t, "Test Organization", result.Attributes.Name)
		assert.Equal(t, "Test Description", result.Attributes.Description)
		assert.Equal(t, "contact@example.com", result.Attributes.ContactEmail)
		assert.Equal(t, "1234567890", result.Attributes.ContactPhone)
		assert.Equal(t, "US", result.Attributes.ContactPhoneCountry)
		assert.Equal(t, "en-US", result.Attributes.Locale)
		assert.Equal(t, "https://example.com/logo.png", result.Meta.LogoDistributionUrl)
		assert.Equal(t, string(constants.OnboardingStatusCompleted), result.Meta.OnboardingStatus)
		assert.False(t, result.Meta.Stripe.HasPendingRequirements)
		assert.Equal(t, string(constants.StripeOnboardingStatusCompleted), result.Meta.Stripe.OnboardingStatus)
		assert.Equal(t, "123 Test St", result.Attributes.Address.Line1)
		assert.Equal(t, "Apt 4", result.Attributes.Address.Line2)
		assert.Equal(t, "Test City", result.Attributes.Address.City)
		assert.Equal(t, "CA", result.Attributes.Address.StateOrProvince)
		assert.Equal(t, "12345", result.Attributes.Address.Zip)
		assert.Equal(t, "US", result.Attributes.Address.Country)
	})

	t.Run("maps organization without optional fields", func(t *testing.T) {
		organization := &models.Organization{
			Model: gorm.Model{
				ID: 456,
			},
			Name:             "Another Organization",
			Description:      "",
			ContactEmail:     "",
			ContactPhone:     "",
			Locale:           "en-US",
			OnboardingStatus: string(constants.OnboardingStatusPending),
			Address: models.Address{
				Line1:   "456 Test Ave",
				City:    "Another City",
				Country: "US",
			},
			Stripe: models.OrganizationStripeAccount{
				HasPendingRequirements: true,
				OnboardingStatus:       string(constants.StripeOnboardingStatusPending),
			},
		}

		dto := &OrganizationDto{
			Organization:        organization,
			LogoDistributionUrl: "",
		}

		result := mapOrganizationToResponse(dto)

		assert.Equal(t, strconv.FormatUint(456, 10), result.Id)
		assert.Equal(t, "Another Organization", result.Attributes.Name)
		assert.Empty(t, result.Attributes.Description)
		assert.Empty(t, result.Attributes.ContactEmail)
		assert.Empty(t, result.Attributes.ContactPhone)
		assert.Empty(t, result.Meta.LogoDistributionUrl)
		assert.Equal(t, string(constants.OnboardingStatusPending), result.Meta.OnboardingStatus)
		assert.True(t, result.Meta.Stripe.HasPendingRequirements)
	})
}

func TestOrganizationsToResponse(t *testing.T) {
	t.Run("maps multiple organizations", func(t *testing.T) {
		org1 := &OrganizationDto{
			Organization: &models.Organization{
				Model: gorm.Model{ID: 1},
				Name:  "Org 1",
			},
			LogoDistributionUrl: "https://example.com/logo1.png",
		}

		org2 := &OrganizationDto{
			Organization: &models.Organization{
				Model: gorm.Model{ID: 2},
				Name:  "Org 2",
			},
			LogoDistributionUrl: "https://example.com/logo2.png",
		}

		result := organizationsToResponse([]*OrganizationDto{org1, org2})

		assert.Len(t, result.Data, 2)
		assert.Equal(t, "Org 1", result.Data[0].Attributes.Name)
		assert.Equal(t, "Org 2", result.Data[1].Attributes.Name)
	})

	t.Run("maps empty list", func(t *testing.T) {
		result := organizationsToResponse([]*OrganizationDto{})

		assert.Len(t, result.Data, 0)
	})
}

func TestMapMembershipToResponse(t *testing.T) {
	t.Run("maps membership with all fields", func(t *testing.T) {
		user := &models.User{
			Model: gorm.Model{ID: 100},
			Name:  "Test User",
			Email: "user@example.com",
		}

		organization := &models.Organization{
			Model: gorm.Model{ID: 200},
			Name:  "Test Org",
		}

		membership := &models.OrganizationMembership{
			Model:          gorm.Model{ID: 50},
			UserID:         100,
			OrganizationID: 200,
			Role:           string(constants.OrganizationRoleAdmin),
		}

		dto := &OrganizationMembershipDto{
			Membership:   membership,
			User:         user,
			Organization: organization,
		}

		result := mapMembershipToResponse(dto)

		assert.Equal(t, strconv.FormatUint(50, 10), result.Id)
		assert.Equal(t, constants.ApiTypeOrganizationMembership, result.Type)
		assert.Equal(t, string(constants.OrganizationRoleAdmin), result.Attributes.Role)
		assert.Equal(t, strconv.FormatUint(100, 10), result.Relationships.User.Data.Id)
		assert.Equal(t, constants.ApiTypeUser, result.Relationships.User.Data.Type)
		assert.Equal(t, strconv.FormatUint(200, 10), result.Relationships.Organization.Data.Id)
		assert.Equal(t, constants.ApiTypeOrganization, result.Relationships.Organization.Data.Type)
	})

	t.Run("maps membership with member role", func(t *testing.T) {
		user := &models.User{
			Model: gorm.Model{ID: 101},
		}

		organization := &models.Organization{
			Model: gorm.Model{ID: 201},
		}

		membership := &models.OrganizationMembership{
			Model:          gorm.Model{ID: 51},
			UserID:         101,
			OrganizationID: 201,
			Role:           string(constants.OrganizationRoleMember),
		}

		dto := &OrganizationMembershipDto{
			Membership:   membership,
			User:         user,
			Organization: organization,
		}

		result := mapMembershipToResponse(dto)

		assert.Equal(t, string(constants.OrganizationRoleMember), result.Attributes.Role)
	})
}

func TestMapUserToIncludedData(t *testing.T) {
	t.Run("maps user to included data", func(t *testing.T) {
		user := &models.User{
			Model: gorm.Model{ID: 100},
			Name:  "Test User",
			Email: "user@example.com",
		}

		organization := &models.Organization{
			Model: gorm.Model{ID: 200},
		}

		membership := &models.OrganizationMembership{
			Model:          gorm.Model{ID: 50},
			UserID:         100,
			OrganizationID: 200,
		}

		dto := &OrganizationMembershipDto{
			Membership:              membership,
			User:                    user,
			UserLogoDistributionUrl: "https://example.com/user-logo.png",
			Organization:            organization,
		}

		result := mapUserToIncludedData(dto)

		assert.Equal(t, strconv.FormatUint(100, 10), result.Id)
		assert.Equal(t, constants.ApiTypeUser, result.Type)
		assert.Equal(t, "Test User", result.Attributes.Name)
		assert.Equal(t, "user@example.com", result.Attributes.Email)
		assert.Equal(t, "https://example.com/user-logo.png", result.Meta.LogoDistributionUrl)
	})
}

func TestMapOrganizationToIncludedData(t *testing.T) {
	t.Run("maps organization to included data", func(t *testing.T) {
		organization := &models.Organization{
			Model: gorm.Model{
				ID: 200,
			},
			Name:        "Test Org",
			Description: "Test Description",
		}

		dto := &OrganizationDto{
			Organization:        organization,
			LogoDistributionUrl: "https://example.com/org-logo.png",
		}

		result := mapOrganizationToIncludedData(dto)

		assert.Equal(t, strconv.FormatUint(200, 10), result.Id)
		assert.Equal(t, constants.ApiTypeOrganization, result.Type)
		assert.Equal(t, "Test Org", result.Attributes.Name)
		assert.Equal(t, "Test Description", result.Attributes.Description)
		assert.Equal(t, "https://example.com/org-logo.png", result.Meta.LogoDistributionUrl)
	})
}

func TestMapInvitingUserToIncludedData(t *testing.T) {
	t.Run("maps inviting user to included data", func(t *testing.T) {
		user := &models.User{
			Model: gorm.Model{ID: 300},
			Name:  "Inviting User",
			Email: "inviting@example.com",
		}

		dto := &InvitingUserDto{
			User:                    user,
			UserLogoDistributionUrl: "https://example.com/inviting-logo.png",
		}

		result := mapInvitingUserToIncludedData(dto)

		assert.Equal(t, strconv.FormatUint(300, 10), result.Id)
		assert.Equal(t, constants.ApiTypeUser, result.Type)
		assert.Equal(t, "Inviting User", result.Attributes.Name)
		assert.Equal(t, "inviting@example.com", result.Attributes.Email)
		assert.Equal(t, "https://example.com/inviting-logo.png", result.Meta.LogoDistributionUrl)
	})
}

func TestMapMembershipsToResponseWithIncluded(t *testing.T) {
	t.Run("maps memberships with included users", func(t *testing.T) {
		user1 := &models.User{
			Model: gorm.Model{ID: 100},
			Name:  "User 1",
		}

		user2 := &models.User{
			Model: gorm.Model{ID: 101},
			Name:  "User 2",
		}

		organization := &models.Organization{
			Model: gorm.Model{ID: 200},
		}

		membership1 := &models.OrganizationMembership{
			Model:          gorm.Model{ID: 50},
			UserID:         100,
			OrganizationID: 200,
			Role:           string(constants.OrganizationRoleAdmin),
		}

		membership2 := &models.OrganizationMembership{
			Model:          gorm.Model{ID: 51},
			UserID:         101,
			OrganizationID: 200,
			Role:           string(constants.OrganizationRoleMember),
		}

		dto1 := &OrganizationMembershipDto{
			Membership:              membership1,
			User:                    user1,
			UserLogoDistributionUrl: "https://example.com/user1.png",
			Organization:            organization,
		}

		dto2 := &OrganizationMembershipDto{
			Membership:              membership2,
			User:                    user2,
			UserLogoDistributionUrl: "https://example.com/user2.png",
			Organization:            organization,
		}

		result := mapMembershipsToResponseWithIncluded([]*OrganizationMembershipDto{dto1, dto2})

		assert.Len(t, result.Data, 2)
		assert.Len(t, result.Included, 2)
		assert.Equal(t, strconv.FormatUint(50, 10), result.Data[0].Id)
		assert.Equal(t, strconv.FormatUint(51, 10), result.Data[1].Id)
	})

	t.Run("deduplicates users in included section", func(t *testing.T) {
		user := &models.User{
			Model: gorm.Model{ID: 100},
			Name:  "Same User",
		}

		organization := &models.Organization{
			Model: gorm.Model{ID: 200},
		}

		membership1 := &models.OrganizationMembership{
			Model:          gorm.Model{ID: 50},
			UserID:         100,
			OrganizationID: 200,
		}

		membership2 := &models.OrganizationMembership{
			Model:          gorm.Model{ID: 51},
			UserID:         100, // Same user
			OrganizationID: 200,
		}

		dto1 := &OrganizationMembershipDto{
			Membership:              membership1,
			User:                    user,
			UserLogoDistributionUrl: "",
			Organization:            organization,
		}

		dto2 := &OrganizationMembershipDto{
			Membership:              membership2,
			User:                    user,
			UserLogoDistributionUrl: "",
			Organization:            organization,
		}

		result := mapMembershipsToResponseWithIncluded([]*OrganizationMembershipDto{dto1, dto2})

		assert.Len(t, result.Data, 2)
		assert.Len(t, result.Included, 1) // Should deduplicate
	})
}

func TestMapInvitationToResponse(t *testing.T) {
	t.Run("maps invitation with all fields", func(t *testing.T) {
		invitationID := uuid.New()
		invitation := &models.OrganizationInvitation{
			Model:          gorm.Model{ID: 1},
			ID:             invitationID,
			Email:          "invitee@example.com",
			OrganizationID: 200,
			InvitingUserID: 300,
			Role:           string(constants.OrganizationRoleAdmin),
			Status:         string(constants.OrganizationInvitationStatusPending),
		}

		dto := &OrganizationInvitationDto{
			Invitation: invitation,
		}

		result := mapInvitationToResponse(dto)

		assert.Equal(t, invitationID.String(), result.Id)
		assert.Equal(t, constants.ApiTypeOrganizationInvitation, result.Type)
		assert.Equal(t, "invitee@example.com", result.Attributes.Email)
		assert.Equal(t, string(constants.OrganizationRoleAdmin), result.Attributes.Role)
		assert.Equal(t, string(constants.OrganizationInvitationStatusPending), result.Attributes.Status)
		assert.Equal(t, strconv.FormatUint(200, 10), result.Relationships.Organization.Data.Id)
		assert.Equal(t, constants.ApiTypeOrganization, result.Relationships.Organization.Data.Type)
		assert.Equal(t, strconv.FormatUint(300, 10), result.Relationships.InvitingUser.Data.Id)
		assert.Equal(t, constants.ApiTypeUser, result.Relationships.InvitingUser.Data.Type)
	})

	t.Run("maps invitation with member role", func(t *testing.T) {
		invitationID := uuid.New()
		invitation := &models.OrganizationInvitation{
			Model:  gorm.Model{ID: 2},
			ID:     invitationID,
			Email:  "member@example.com",
			Role:   string(constants.OrganizationRoleMember),
			Status: string(constants.OrganizationInvitationStatusAccepted),
		}

		dto := &OrganizationInvitationDto{
			Invitation: invitation,
		}

		result := mapInvitationToResponse(dto)

		assert.Equal(t, string(constants.OrganizationRoleMember), result.Attributes.Role)
		assert.Equal(t, string(constants.OrganizationInvitationStatusAccepted), result.Attributes.Status)
	})
}

func TestMapInvitationsToResponse(t *testing.T) {
	t.Run("maps multiple invitations", func(t *testing.T) {
		invitation1 := &models.OrganizationInvitation{
			Model:  gorm.Model{ID: 1},
			ID:     uuid.New(),
			Email:  "invitee1@example.com",
			Status: string(constants.OrganizationInvitationStatusPending),
		}

		invitation2 := &models.OrganizationInvitation{
			Model:  gorm.Model{ID: 2},
			ID:     uuid.New(),
			Email:  "invitee2@example.com",
			Status: string(constants.OrganizationInvitationStatusPending),
		}

		dto1 := &OrganizationInvitationDto{
			Invitation: invitation1,
		}

		dto2 := &OrganizationInvitationDto{
			Invitation: invitation2,
		}

		result := mapInvitationsToResponse([]*OrganizationInvitationDto{dto1, dto2})

		assert.Len(t, result.Data, 2)
		assert.Equal(t, "invitee1@example.com", result.Data[0].Attributes.Email)
		assert.Equal(t, "invitee2@example.com", result.Data[1].Attributes.Email)
	})

	t.Run("maps empty list", func(t *testing.T) {
		result := mapInvitationsToResponse([]*OrganizationInvitationDto{})

		assert.Len(t, result.Data, 0)
	})
}
