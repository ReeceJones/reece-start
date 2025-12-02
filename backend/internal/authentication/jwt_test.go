package authentication

import (
	"testing"
	"time"

	"reece.start/internal/constants"
	testconfig "reece.start/test/config"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateJWT(t *testing.T) {
	t.Run("MinimalOptions", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId: uuid.New(),
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Validate the token
		claims, err := ValidateJWT(config, token)
		require.NoError(t, err)
		require.Equal(t, options.UserId.String(), claims.UserId)
		require.Nil(t, claims.OrganizationId)
		require.Nil(t, claims.OrganizationRole)
		require.Nil(t, claims.Scopes)
		require.Nil(t, claims.Role)
		require.Nil(t, claims.IsImpersonating)
		require.Nil(t, claims.ImpersonatingUserId)
		require.Equal(t, config.JwtIssuer, claims.Issuer)
		require.Equal(t, config.JwtAudience, claims.Audience[0])
		require.Equal(t, options.UserId.String(), claims.Subject)
	})

	t.Run("WithOrganization", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		orgId := uuid.New()
		orgRole := constants.OrganizationRoleAdmin
		options := JwtOptions{
			UserId:           uuid.New(),
			OrganizationId:   &orgId,
			OrganizationRole: &orgRole,
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := ValidateJWT(config, token)
		require.NoError(t, err)
		require.Equal(t, options.UserId.String(), claims.UserId)
		require.NotNil(t, claims.OrganizationId)
		require.Equal(t, orgId.String(), *claims.OrganizationId)
		require.NotNil(t, claims.OrganizationRole)
		require.Equal(t, constants.OrganizationRoleAdmin, *claims.OrganizationRole)
	})

	t.Run("WithAllOptions", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		orgId := uuid.New()
		orgRole := constants.OrganizationRoleAdmin
		userRole := constants.UserRoleAdmin
		scopes := []constants.UserScope{
			constants.UserScopeAdmin,
			constants.UserScopeOrganizationRead,
		}
		isImpersonating := true
		impersonatingUserId := uuid.New()

		options := JwtOptions{
			UserId:              uuid.New(),
			OrganizationId:      &orgId,
			OrganizationRole:    &orgRole,
			Role:                &userRole,
			Scopes:              &scopes,
			IsImpersonating:     &isImpersonating,
			ImpersonatingUserId: &impersonatingUserId,
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := ValidateJWT(config, token)
		require.NoError(t, err)
		require.Equal(t, options.UserId.String(), claims.UserId)
		require.Equal(t, orgId.String(), *claims.OrganizationId)
		require.Equal(t, constants.OrganizationRoleAdmin, *claims.OrganizationRole)
		require.Equal(t, constants.UserRoleAdmin, *claims.Role)
		require.NotNil(t, claims.Scopes)
		require.Len(t, *claims.Scopes, 2)
		require.Equal(t, constants.UserScopeAdmin, (*claims.Scopes)[0])
		require.Equal(t, constants.UserScopeOrganizationRead, (*claims.Scopes)[1])
		require.NotNil(t, claims.IsImpersonating)
		require.True(t, *claims.IsImpersonating)
		require.Equal(t, impersonatingUserId.String(), *claims.ImpersonatingUserId)
	})

	t.Run("WithCustomExpiry", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		customExpiry := time.Now().Add(2 * time.Hour)
		options := JwtOptions{
			UserId:       uuid.New(),
			CustomExpiry: &customExpiry,
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := ValidateJWT(config, token)
		require.NoError(t, err)
		require.Equal(t, options.UserId.String(), claims.UserId)
		require.NotNil(t, claims.ExpiresAt)
		// Check that expiry is approximately 2 hours from now (within 1 minute tolerance)
		expectedExpiry := customExpiry.Unix()
		actualExpiry := claims.ExpiresAt.Unix()
		require.InDelta(t, expectedExpiry, actualExpiry, 60) // 60 second tolerance
	})

	t.Run("WithDefaultExpiry", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId: uuid.New(),
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := ValidateJWT(config, token)
		require.NoError(t, err)
		require.NotNil(t, claims.ExpiresAt)

		// Check that expiry is approximately 1 hour from now (within 1 minute tolerance)
		expectedExpiry := time.Now().Add(time.Duration(config.JwtExpirationTime) * time.Second).Unix()
		actualExpiry := claims.ExpiresAt.Unix()
		require.InDelta(t, expectedExpiry, actualExpiry, 60) // 60 second tolerance
	})

	t.Run("RoundTripWithValidateJWT", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		orgId := uuid.New()
		orgRole := constants.OrganizationRoleMember
		userRole := constants.UserRoleDefault
		scopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}
		isImpersonating := false

		options := JwtOptions{
			UserId:           uuid.New(),
			OrganizationId:   &orgId,
			OrganizationRole: &orgRole,
			Role:             &userRole,
			Scopes:           &scopes,
			IsImpersonating:  &isImpersonating,
		}

		// Create token
		token, err := CreateJWT(config, options)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Validate token
		claims, err := ValidateJWT(config, token)
		require.NoError(t, err)
		require.NotNil(t, claims)

		// Verify all fields
		require.Equal(t, options.UserId.String(), claims.UserId)
		require.NotNil(t, claims.OrganizationId)
		require.Equal(t, orgId.String(), *claims.OrganizationId)
		require.NotNil(t, claims.OrganizationRole)
		require.Equal(t, constants.OrganizationRoleMember, *claims.OrganizationRole)
		require.NotNil(t, claims.Role)
		require.Equal(t, constants.UserRoleDefault, *claims.Role)
		require.NotNil(t, claims.Scopes)
		require.Len(t, *claims.Scopes, 1)
		require.Equal(t, constants.UserScopeOrganizationRead, (*claims.Scopes)[0])
		require.NotNil(t, claims.IsImpersonating)
		require.False(t, *claims.IsImpersonating)
		require.Nil(t, claims.ImpersonatingUserId)
	})
}

func TestValidateJWT(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId: uuid.New(),
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)

		claims, err := ValidateJWT(config, token)
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.Equal(t, options.UserId.String(), claims.UserId)
	})

	t.Run("InvalidSecret", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId: uuid.New(),
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)

		// Use different secret for validation
		invalidConfig := testconfig.CreateTestConfig()
		invalidConfig.JwtSecret = "different-secret"

		claims, err := ValidateJWT(invalidConfig, token)
		require.Error(t, err)
		require.Nil(t, claims)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		pastTime := time.Now().Add(-1 * time.Hour)
		options := JwtOptions{
			UserId:       uuid.New(),
			CustomExpiry: &pastTime,
		}

		token, err := CreateJWT(config, options)
		require.NoError(t, err)

		claims, err := ValidateJWT(config, token)
		require.Error(t, err)
		require.Nil(t, claims)

		// Check that error message contains expiration-related text
		require.Contains(t, err.Error(), "expired")
	})

	t.Run("MalformedToken", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		malformedToken := "not.a.valid.jwt.token"

		claims, err := ValidateJWT(config, malformedToken)
		require.Error(t, err)
		require.Nil(t, claims)
	})

	t.Run("EmptyToken", func(t *testing.T) {
		config := testconfig.CreateTestConfig()

		claims, err := ValidateJWT(config, "")
		require.Error(t, err)
		require.Nil(t, claims)
	})

	t.Run("WrongIssuer", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId: uuid.New(),
		}

		// Create token with different issuer
		wrongIssuerConfig := testconfig.CreateTestConfig()
		wrongIssuerConfig.JwtIssuer = "wrong-issuer"
		wrongIssuerToken, err := CreateJWT(wrongIssuerConfig, options)
		require.NoError(t, err)

		// Validate with original config
		// Note: jwt-go v5 doesn't validate issuer by default, so this should still succeed
		claims, err := ValidateJWT(config, wrongIssuerToken)
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.Equal(t, "wrong-issuer", claims.Issuer)
	})

	t.Run("WrongAudience", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId: uuid.New(),
		}

		// Create token with different audience
		wrongAudienceConfig := testconfig.CreateTestConfig()
		wrongAudienceConfig.JwtAudience = "wrong-audience"
		wrongAudienceToken, err := CreateJWT(wrongAudienceConfig, options)
		require.NoError(t, err)

		// Validate with original config
		// Note: jwt-go v5 doesn't validate audience by default, so this should still succeed
		claims, err := ValidateJWT(config, wrongAudienceToken)
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.Contains(t, claims.Audience, "wrong-audience")
	})
}

func TestGetActiveOrganizationIdFromOptions(t *testing.T) {
	t.Run("WithOrganizationId", func(t *testing.T) {
		orgId := uuid.New()
		options := JwtOptions{
			UserId:         uuid.New(),
			OrganizationId: &orgId,
		}

		result := getActiveOrganizationIdFromOptions(options)
		require.NotNil(t, result)
		require.Equal(t, orgId.String(), *result)
	})

	t.Run("WithoutOrganizationId", func(t *testing.T) {
		options := JwtOptions{
			UserId: uuid.New(),
		}

		result := getActiveOrganizationIdFromOptions(options)
		require.Nil(t, result)
	})

	t.Run("NilOrganizationId", func(t *testing.T) {
		options := JwtOptions{
			UserId:         uuid.New(),
			OrganizationId: nil,
		}

		result := getActiveOrganizationIdFromOptions(options)
		require.Nil(t, result)
	})
}

func TestGetExpiryFromOptions(t *testing.T) {
	t.Run("WithCustomExpiry", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		customExpiry := time.Now().Add(2 * time.Hour)
		options := JwtOptions{
			UserId:       uuid.New(),
			CustomExpiry: &customExpiry,
		}

		result := getExpiryFromOptions(config, options)
		require.NotNil(t, result)

		expectedExpiry := customExpiry.Unix()
		actualExpiry := result.Unix()
		require.Equal(t, expectedExpiry, actualExpiry)
	})

	t.Run("WithoutCustomExpiry", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId: uuid.New(),
		}

		result := getExpiryFromOptions(config, options)
		require.NotNil(t, result)

		// Check that expiry is approximately 1 hour from now (within 1 minute tolerance)
		expectedExpiry := time.Now().Add(time.Duration(config.JwtExpirationTime) * time.Second).Unix()
		actualExpiry := result.Unix()
		require.InDelta(t, expectedExpiry, actualExpiry, 60) // 60 second tolerance
	})

	t.Run("NilCustomExpiry", func(t *testing.T) {
		config := testconfig.CreateTestConfig()
		options := JwtOptions{
			UserId:       uuid.New(),
			CustomExpiry: nil,
		}

		result := getExpiryFromOptions(config, options)
		require.NotNil(t, result)

		// Check that expiry is approximately 1 hour from now (within 1 minute tolerance)
		expectedExpiry := time.Now().Add(time.Duration(config.JwtExpirationTime) * time.Second).Unix()
		actualExpiry := result.Unix()
		require.InDelta(t, expectedExpiry, actualExpiry, 60) // 60 second tolerance
	})
}
