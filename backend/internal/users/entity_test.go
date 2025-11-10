package users

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

func TestMapUserToResponse(t *testing.T) {
	t.Run("maps user with all fields", func(t *testing.T) {
		now := time.Now()
		user := &models.User{
			Model: gorm.Model{
				ID: 123,
			},
			Name:  "Test User",
			Email: "test@example.com",
			Revocation: models.UserTokenRevocation{
				LastValidIssuedAt: &now,
				CanRefresh:        true,
			},
		}

		dto := &UserDto{
			User:                user,
			Token:               "test-token",
			LogoDistributionUrl: "https://example.com/logo.png",
		}

		result := mapUserToResponse(dto)

		assert.Equal(t, strconv.FormatUint(123, 10), result.Data.Id)
		assert.Equal(t, constants.ApiTypeUser, result.Data.Type)
		assert.Equal(t, "Test User", result.Data.Attributes.Name)
		assert.Equal(t, "test@example.com", result.Data.Attributes.Email)
		assert.Equal(t, "test-token", result.Data.Meta.Token)
		assert.Equal(t, "https://example.com/logo.png", result.Data.Meta.LogoDistributionUrl)
		assert.Equal(t, &now, result.Data.Meta.TokenRevocation.LastIssuedAt)
		assert.True(t, result.Data.Meta.TokenRevocation.CanRefresh)
	})

	t.Run("maps user without optional fields", func(t *testing.T) {
		user := &models.User{
			Model: gorm.Model{
				ID: 456,
			},
			Name:  "Another User",
			Email: "another@example.com",
			Revocation: models.UserTokenRevocation{
				LastValidIssuedAt: nil,
				CanRefresh:        false,
			},
		}

		dto := &UserDto{
			User:                user,
			Token:               "",
			LogoDistributionUrl: "",
		}

		result := mapUserToResponse(dto)

		assert.Equal(t, strconv.FormatUint(456, 10), result.Data.Id)
		assert.Equal(t, "Another User", result.Data.Attributes.Name)
		assert.Equal(t, "another@example.com", result.Data.Attributes.Email)
		assert.Empty(t, result.Data.Meta.Token)
		assert.Empty(t, result.Data.Meta.LogoDistributionUrl)
		assert.Nil(t, result.Data.Meta.TokenRevocation.LastIssuedAt)
		assert.False(t, result.Data.Meta.TokenRevocation.CanRefresh)
	})
}

func TestMapCreateAuthenticatedUserTokenToResponse(t *testing.T) {
	t.Run("maps token response with organization", func(t *testing.T) {
		req := CreateAuthenticatedUserTokenRequest{
			Data: CreateAuthenticatedUserTokenData{
				Type: constants.ApiTypeToken,
				Relationships: CreateAuthenticatedUserTokenRelationships{
					Organization: &OrganizationRelationship{
						Data: OrganizationRelationshipData{
							Id:   "123",
							Type: "organization",
						},
					},
				},
				Meta: CreateAuthenticatedUserTokenRequestMeta{
					StopImpersonating: false,
				},
			},
		}

		token := "jwt-token-123"

		result := mapCreateAuthenticatedUserTokenToResponse(req, token)

		assert.Equal(t, constants.ApiTypeToken, result.Data.Type)
		assert.Equal(t, token, result.Data.Meta.Token)
		assert.NotNil(t, result.Data.Relationships.Organization)
		assert.Equal(t, "123", result.Data.Relationships.Organization.Data.Id)
		assert.Equal(t, "organization", result.Data.Relationships.Organization.Data.Type)
	})

	t.Run("maps token response with impersonated user", func(t *testing.T) {
		req := CreateAuthenticatedUserTokenRequest{
			Data: CreateAuthenticatedUserTokenData{
				Type: constants.ApiTypeToken,
				Relationships: CreateAuthenticatedUserTokenRelationships{
					ImpersonatedUser: &UserRelationship{
						Data: UserRelationshipData{
							Id:   "456",
							Type: "user",
						},
					},
				},
				Meta: CreateAuthenticatedUserTokenRequestMeta{
					StopImpersonating: true,
				},
			},
		}

		token := "jwt-token-456"

		result := mapCreateAuthenticatedUserTokenToResponse(req, token)

		assert.Equal(t, constants.ApiTypeToken, result.Data.Type)
		assert.Equal(t, token, result.Data.Meta.Token)
		assert.NotNil(t, result.Data.Relationships.ImpersonatedUser)
		assert.Equal(t, "456", result.Data.Relationships.ImpersonatedUser.Data.Id)
		assert.Equal(t, "user", result.Data.Relationships.ImpersonatedUser.Data.Type)
	})

	t.Run("maps token response with both relationships", func(t *testing.T) {
		req := CreateAuthenticatedUserTokenRequest{
			Data: CreateAuthenticatedUserTokenData{
				Type: constants.ApiTypeToken,
				Relationships: CreateAuthenticatedUserTokenRelationships{
					Organization: &OrganizationRelationship{
						Data: OrganizationRelationshipData{
							Id:   "789",
							Type: "organization",
						},
					},
					ImpersonatedUser: &UserRelationship{
						Data: UserRelationshipData{
							Id:   "101",
							Type: "user",
						},
					},
				},
				Meta: CreateAuthenticatedUserTokenRequestMeta{
					StopImpersonating: false,
				},
			},
		}

		token := "jwt-token-789"

		result := mapCreateAuthenticatedUserTokenToResponse(req, token)

		assert.Equal(t, token, result.Data.Meta.Token)
		assert.NotNil(t, result.Data.Relationships.Organization)
		assert.NotNil(t, result.Data.Relationships.ImpersonatedUser)
	})
}
