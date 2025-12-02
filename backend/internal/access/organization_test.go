package access

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/authentication"
	"reece.start/internal/constants"

	"github.com/stretchr/testify/require"
)

// createTestContext creates an echo context with JWT claims
func createTestContext(t *testing.T, claims *authentication.JwtClaims) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("claims", claims)
	return c
}

func TestHasOrganizationAccess(t *testing.T) {
	t.Run("WithRequiredScopes", func(t *testing.T) {
		scopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
			constants.UserScopeOrganizationUpdate,
		}
		claims := &authentication.JwtClaims{
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		params := HasOrganizationAccessParams{
			OrganizationID: uuid.New(),
			Scopes: []constants.UserScope{
				constants.UserScopeOrganizationRead,
			},
		}

		err := HasOrganizationAccess(c, params)
		require.NoError(t, err)
	})

	t.Run("WithAllRequiredScopes", func(t *testing.T) {
		scopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
			constants.UserScopeOrganizationUpdate,
			constants.UserScopeOrganizationDelete,
		}
		claims := &authentication.JwtClaims{
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		params := HasOrganizationAccessParams{
			OrganizationID: uuid.New(),
			Scopes: []constants.UserScope{
				constants.UserScopeOrganizationRead,
				constants.UserScopeOrganizationUpdate,
			},
		}

		err := HasOrganizationAccess(c, params)
		require.NoError(t, err)
	})

	t.Run("MissingRequiredScope", func(t *testing.T) {
		scopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}
		claims := &authentication.JwtClaims{
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		params := HasOrganizationAccessParams{
			OrganizationID: uuid.New(),
			Scopes: []constants.UserScope{
				constants.UserScopeOrganizationRead,
				constants.UserScopeOrganizationUpdate, // Missing this scope
			},
		}

		err := HasOrganizationAccess(c, params)
		require.Error(t, err)
		require.Equal(t, api.ErrForbiddenNoAccess, err)
	})

	t.Run("NoScopes", func(t *testing.T) {
		scopes := []constants.UserScope{}
		claims := &authentication.JwtClaims{
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		params := HasOrganizationAccessParams{
			OrganizationID: uuid.New(),
			Scopes: []constants.UserScope{
				constants.UserScopeOrganizationRead,
			},
		}

		err := HasOrganizationAccess(c, params)
		require.Error(t, err)
		require.Equal(t, api.ErrForbiddenNoAccess, err)
	})

	t.Run("MissingScopesInClaims", func(t *testing.T) {
		claims := &authentication.JwtClaims{
			Scopes: nil,
		}
		c := createTestContext(t, claims)

		params := HasOrganizationAccessParams{
			OrganizationID: uuid.New(),
			Scopes: []constants.UserScope{
				constants.UserScopeOrganizationRead,
			},
		}

		err := HasOrganizationAccess(c, params)
		require.Error(t, err)
		// Should return error from GetScopesFromJWT when scopes are nil
		require.Contains(t, err.Error(), "scopes are not set")
	})

	t.Run("EmptyRequiredScopes", func(t *testing.T) {
		scopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}
		claims := &authentication.JwtClaims{
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		params := HasOrganizationAccessParams{
			OrganizationID: uuid.New(),
			Scopes:         []constants.UserScope{},
		}

		err := HasOrganizationAccess(c, params)
		require.NoError(t, err) // Empty scopes means no requirements, so access is granted
	})
}

func TestHasAdminAccess(t *testing.T) {
	t.Run("WithAdminRoleAndRequiredScopes", func(t *testing.T) {
		role := constants.UserRoleAdmin
		scopes := []constants.UserScope{
			constants.UserScopeAdmin,
			constants.UserScopeAdminUsersList,
			constants.UserScopeAdminUsersRead,
		}
		claims := &authentication.JwtClaims{
			Role:   &role,
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{
			constants.UserScopeAdminUsersList,
		}

		err := HasAdminAccess(c, requiredScopes)
		require.NoError(t, err)
	})

	t.Run("WithAdminRoleAndAllRequiredScopes", func(t *testing.T) {
		role := constants.UserRoleAdmin
		scopes := []constants.UserScope{
			constants.UserScopeAdmin,
			constants.UserScopeAdminUsersList,
			constants.UserScopeAdminUsersRead,
			constants.UserScopeAdminUsersImpersonate,
		}
		claims := &authentication.JwtClaims{
			Role:   &role,
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{
			constants.UserScopeAdminUsersList,
			constants.UserScopeAdminUsersRead,
		}

		err := HasAdminAccess(c, requiredScopes)
		require.NoError(t, err)
	})

	t.Run("NonAdminRole", func(t *testing.T) {
		role := constants.UserRoleDefault
		scopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}
		claims := &authentication.JwtClaims{
			Role:   &role,
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{
			constants.UserScopeAdminUsersList,
		}

		err := HasAdminAccess(c, requiredScopes)
		require.Error(t, err)
		require.Equal(t, api.ErrForbiddenNoAdminAccess, err)
	})

	t.Run("AdminRoleButMissingScopes", func(t *testing.T) {
		role := constants.UserRoleAdmin
		scopes := []constants.UserScope{
			constants.UserScopeAdmin,
		}
		claims := &authentication.JwtClaims{
			Role:   &role,
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{
			constants.UserScopeAdminUsersList, // Missing this scope
		}

		err := HasAdminAccess(c, requiredScopes)
		require.Error(t, err)
		require.Equal(t, api.ErrForbiddenNoAdminAccess, err)
	})

	t.Run("AdminRoleButNoScopes", func(t *testing.T) {
		role := constants.UserRoleAdmin
		scopes := []constants.UserScope{}
		claims := &authentication.JwtClaims{
			Role:   &role,
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{
			constants.UserScopeAdminUsersList,
		}

		err := HasAdminAccess(c, requiredScopes)
		require.Error(t, err)
		require.Equal(t, api.ErrForbiddenNoAdminAccess, err)
	})

	t.Run("MissingRoleInClaims", func(t *testing.T) {
		scopes := []constants.UserScope{
			constants.UserScopeAdmin,
			constants.UserScopeAdminUsersList,
		}
		claims := &authentication.JwtClaims{
			Role:   nil,
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{
			constants.UserScopeAdminUsersList,
		}

		err := HasAdminAccess(c, requiredScopes)
		require.Error(t, err)
		// Should return error from GetRoleFromJWT when role is nil
		require.Contains(t, err.Error(), "role is not set")
	})

	t.Run("MissingScopesInClaims", func(t *testing.T) {
		role := constants.UserRoleAdmin
		claims := &authentication.JwtClaims{
			Role:   &role,
			Scopes: nil,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{
			constants.UserScopeAdminUsersList,
		}

		err := HasAdminAccess(c, requiredScopes)
		require.Error(t, err)
		// Should return error from GetScopesFromJWT when scopes are nil
		require.Contains(t, err.Error(), "scopes are not set")
	})

	t.Run("EmptyRequiredScopes", func(t *testing.T) {
		role := constants.UserRoleAdmin
		scopes := []constants.UserScope{
			constants.UserScopeAdmin,
		}
		claims := &authentication.JwtClaims{
			Role:   &role,
			Scopes: &scopes,
		}
		c := createTestContext(t, claims)

		requiredScopes := []constants.UserScope{}

		err := HasAdminAccess(c, requiredScopes)
		require.NoError(t, err) // Empty scopes means no requirements, so access is granted
	})
}

func TestHasScopes(t *testing.T) {
	t.Run("AllScopesPresent", func(t *testing.T) {
		requiredScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
			constants.UserScopeOrganizationUpdate,
		}
		grantedScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
			constants.UserScopeOrganizationUpdate,
			constants.UserScopeOrganizationDelete,
		}

		result := hasScopes(requiredScopes, grantedScopes)
		require.True(t, result)
	})

	t.Run("SomeScopesPresent", func(t *testing.T) {
		requiredScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
			constants.UserScopeOrganizationUpdate,
		}
		grantedScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}

		result := hasScopes(requiredScopes, grantedScopes)
		require.False(t, result)
	})

	t.Run("NoScopesPresent", func(t *testing.T) {
		requiredScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
			constants.UserScopeOrganizationUpdate,
		}
		grantedScopes := []constants.UserScope{
			constants.UserScopeOrganizationDelete,
		}

		result := hasScopes(requiredScopes, grantedScopes)
		require.False(t, result)
	})

	t.Run("EmptyRequiredScopes", func(t *testing.T) {
		requiredScopes := []constants.UserScope{}
		grantedScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}

		result := hasScopes(requiredScopes, grantedScopes)
		require.True(t, result) // Empty required scopes means no requirements
	})

	t.Run("EmptyGrantedScopes", func(t *testing.T) {
		requiredScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}
		grantedScopes := []constants.UserScope{}

		result := hasScopes(requiredScopes, grantedScopes)
		require.False(t, result)
	})

	t.Run("BothEmpty", func(t *testing.T) {
		requiredScopes := []constants.UserScope{}
		grantedScopes := []constants.UserScope{}

		result := hasScopes(requiredScopes, grantedScopes)
		require.True(t, result) // Empty required scopes means no requirements
	})

	t.Run("SingleScopeMatch", func(t *testing.T) {
		requiredScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}
		grantedScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}

		result := hasScopes(requiredScopes, grantedScopes)
		require.True(t, result)
	})

	t.Run("SingleScopeMismatch", func(t *testing.T) {
		requiredScopes := []constants.UserScope{
			constants.UserScopeOrganizationRead,
		}
		grantedScopes := []constants.UserScope{
			constants.UserScopeOrganizationUpdate,
		}

		result := hasScopes(requiredScopes, grantedScopes)
		require.False(t, result)
	})
}

func TestHasRole(t *testing.T) {
	t.Run("MatchingRoles", func(t *testing.T) {
		role := constants.UserRoleAdmin
		grantedRole := constants.UserRoleAdmin

		result := hasRole(role, grantedRole)
		require.True(t, result)
	})

	t.Run("NonMatchingRoles", func(t *testing.T) {
		role := constants.UserRoleDefault
		grantedRole := constants.UserRoleAdmin

		result := hasRole(role, grantedRole)
		require.False(t, result)
	})

	t.Run("DefaultRole", func(t *testing.T) {
		role := constants.UserRoleDefault
		grantedRole := constants.UserRoleDefault

		result := hasRole(role, grantedRole)
		require.True(t, result)
	})
}
