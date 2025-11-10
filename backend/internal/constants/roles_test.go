package constants

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserRoleToScopes(t *testing.T) {
	t.Run("AdminRole", func(t *testing.T) {
		scopes, exists := UserRoleToScopes[UserRoleAdmin]
		require.True(t, exists, "UserRoleAdmin should exist in UserRoleToScopes")

		expectedScopes := []UserScope{
			UserScopeAdmin,
			UserScopeAdminUsersList,
			UserScopeAdminUsersRead,
			UserScopeAdminUsersImpersonate,
		}

		require.Equal(t, len(expectedScopes), len(scopes), "Admin role should have correct number of scopes")

		for _, expectedScope := range expectedScopes {
			require.Contains(t, scopes, expectedScope, "Admin role should contain scope: %s", expectedScope)
		}

		// Verify no unexpected scopes
		for _, scope := range scopes {
			require.Contains(t, expectedScopes, scope, "Admin role should not contain unexpected scope: %s", scope)
		}
	})

	t.Run("DefaultRole", func(t *testing.T) {
		scopes, exists := UserRoleToScopes[UserRoleDefault]
		require.True(t, exists, "UserRoleDefault should exist in UserRoleToScopes")
		require.Empty(t, scopes, "Default role should have no scopes")
	})

	t.Run("AllRolesCovered", func(t *testing.T) {
		allUserRoles := []UserRole{
			UserRoleAdmin,
			UserRoleDefault,
		}

		for _, role := range allUserRoles {
			_, exists := UserRoleToScopes[role]
			require.True(t, exists, "Role %s should exist in UserRoleToScopes", role)
		}
	})

	t.Run("NoExtraRoles", func(t *testing.T) {
		expectedRoles := map[UserRole]bool{
			UserRoleAdmin:   true,
			UserRoleDefault: true,
		}

		for role := range UserRoleToScopes {
			require.True(t, expectedRoles[role], "Unexpected role found in UserRoleToScopes: %s", role)
		}
	})

	t.Run("AdminHasNoOrganizationScopes", func(t *testing.T) {
		adminScopes := UserRoleToScopes[UserRoleAdmin]

		// User admin scopes should not contain organization scopes
		organizationScopes := []UserScope{
			UserScopeOrganizationRead,
			UserScopeOrganizationUpdate,
			UserScopeOrganizationDelete,
			UserScopeOrganizationMembershipsList,
			UserScopeOrganizationMembershipsRead,
			UserScopeOrganizationMembershipsCreate,
			UserScopeOrganizationMembershipsUpdate,
			UserScopeOrganizationMembershipsDelete,
			UserScopeOrganizationInvitationsList,
			UserScopeOrganizationInvitationsRead,
			UserScopeOrganizationInvitationsCreate,
			UserScopeOrganizationInvitationsUpdate,
			UserScopeOrganizationInvitationsDelete,
			UserScopeOrganizationStripeUpdate,
			UserScopeOrganizationBillingUpdate,
		}

		for _, orgScope := range organizationScopes {
			require.NotContains(t, adminScopes, orgScope, "User admin role should not contain organization scope: %s", orgScope)
		}
	})

	t.Run("NoDuplicates", func(t *testing.T) {
		for role, scopes := range UserRoleToScopes {
			seen := make(map[UserScope]bool)
			for _, scope := range scopes {
				require.False(t, seen[scope], "Role %s should not have duplicate scope: %s", role, scope)
				seen[scope] = true
			}
		}
	})
}

func TestOrganizationRoleToScopes(t *testing.T) {
	t.Run("AdminRole", func(t *testing.T) {
		scopes, exists := OrganizationRoleToScopes[OrganizationRoleAdmin]
		require.True(t, exists, "OrganizationRoleAdmin should exist in OrganizationRoleToScopes")

		expectedScopes := []UserScope{
			UserScopeOrganizationRead,
			UserScopeOrganizationUpdate,
			UserScopeOrganizationDelete,
			UserScopeOrganizationMembershipsList,
			UserScopeOrganizationMembershipsRead,
			UserScopeOrganizationMembershipsCreate,
			UserScopeOrganizationMembershipsUpdate,
			UserScopeOrganizationMembershipsDelete,
			UserScopeOrganizationInvitationsList,
			UserScopeOrganizationInvitationsRead,
			UserScopeOrganizationInvitationsCreate,
			UserScopeOrganizationInvitationsUpdate,
			UserScopeOrganizationInvitationsDelete,
			UserScopeOrganizationStripeUpdate,
			UserScopeOrganizationBillingUpdate,
		}

		require.Equal(t, len(expectedScopes), len(scopes), "Organization admin role should have correct number of scopes")

		for _, expectedScope := range expectedScopes {
			require.Contains(t, scopes, expectedScope, "Organization admin role should contain scope: %s", expectedScope)
		}

		// Verify no unexpected scopes
		for _, scope := range scopes {
			require.Contains(t, expectedScopes, scope, "Organization admin role should not contain unexpected scope: %s", scope)
		}
	})

	t.Run("MemberRole", func(t *testing.T) {
		scopes, exists := OrganizationRoleToScopes[OrganizationRoleMember]
		require.True(t, exists, "OrganizationRoleMember should exist in OrganizationRoleToScopes")

		expectedScopes := []UserScope{
			UserScopeOrganizationRead,
			UserScopeOrganizationMembershipsList,
			UserScopeOrganizationMembershipsRead,
			UserScopeOrganizationInvitationsList,
			UserScopeOrganizationInvitationsRead,
		}

		require.Equal(t, len(expectedScopes), len(scopes), "Organization member role should have correct number of scopes")

		for _, expectedScope := range expectedScopes {
			require.Contains(t, scopes, expectedScope, "Organization member role should contain scope: %s", expectedScope)
		}

		// Verify no unexpected scopes
		for _, scope := range scopes {
			require.Contains(t, expectedScopes, scope, "Organization member role should not contain unexpected scope: %s", scope)
		}
	})

	t.Run("MemberRoleHasOnlyReadScopes", func(t *testing.T) {
		scopes, exists := OrganizationRoleToScopes[OrganizationRoleMember]
		require.True(t, exists)

		// Member should NOT have write/update/delete/create scopes
		writeScopes := []UserScope{
			UserScopeOrganizationUpdate,
			UserScopeOrganizationDelete,
			UserScopeOrganizationMembershipsCreate,
			UserScopeOrganizationMembershipsUpdate,
			UserScopeOrganizationMembershipsDelete,
			UserScopeOrganizationInvitationsCreate,
			UserScopeOrganizationInvitationsUpdate,
			UserScopeOrganizationInvitationsDelete,
			UserScopeOrganizationStripeUpdate,
			UserScopeOrganizationBillingUpdate,
		}

		for _, writeScope := range writeScopes {
			require.NotContains(t, scopes, writeScope, "Organization member role should not have write scope: %s", writeScope)
		}
	})

	t.Run("AllRolesCovered", func(t *testing.T) {
		allOrganizationRoles := []OrganizationRole{
			OrganizationRoleAdmin,
			OrganizationRoleMember,
		}

		for _, role := range allOrganizationRoles {
			_, exists := OrganizationRoleToScopes[role]
			require.True(t, exists, "Role %s should exist in OrganizationRoleToScopes", role)
		}
	})

	t.Run("NoExtraRoles", func(t *testing.T) {
		expectedRoles := map[OrganizationRole]bool{
			OrganizationRoleAdmin:  true,
			OrganizationRoleMember: true,
		}

		for role := range OrganizationRoleToScopes {
			require.True(t, expectedRoles[role], "Unexpected role found in OrganizationRoleToScopes: %s", role)
		}
	})

	t.Run("AdminHasAllScopes", func(t *testing.T) {
		adminScopes := OrganizationRoleToScopes[OrganizationRoleAdmin]
		memberScopes := OrganizationRoleToScopes[OrganizationRoleMember]

		// Admin should have all scopes that member has, plus more
		for _, memberScope := range memberScopes {
			require.Contains(t, adminScopes, memberScope, "Admin role should have all scopes that member role has: %s", memberScope)
		}

		// Admin should have more scopes than member
		require.Greater(t, len(adminScopes), len(memberScopes), "Admin role should have more scopes than member role")
	})

	t.Run("NoUserAdminScopes", func(t *testing.T) {
		adminScopes := OrganizationRoleToScopes[OrganizationRoleAdmin]
		memberScopes := OrganizationRoleToScopes[OrganizationRoleMember]

		// Organization roles should not contain user admin scopes
		userAdminScopes := []UserScope{
			UserScopeAdmin,
			UserScopeAdminUsersList,
			UserScopeAdminUsersRead,
			UserScopeAdminUsersImpersonate,
		}

		for _, userAdminScope := range userAdminScopes {
			require.NotContains(t, adminScopes, userAdminScope, "Organization admin role should not contain user admin scope: %s", userAdminScope)
			require.NotContains(t, memberScopes, userAdminScope, "Organization member role should not contain user admin scope: %s", userAdminScope)
		}
	})

	t.Run("NoDuplicates", func(t *testing.T) {
		for role, scopes := range OrganizationRoleToScopes {
			seen := make(map[UserScope]bool)
			for _, scope := range scopes {
				require.False(t, seen[scope], "Role %s should not have duplicate scope: %s", role, scope)
				seen[scope] = true
			}
		}
	})

	t.Run("MemberScopesAreSubsetOfAdminScopes", func(t *testing.T) {
		adminScopes := OrganizationRoleToScopes[OrganizationRoleAdmin]
		memberScopes := OrganizationRoleToScopes[OrganizationRoleMember]

		// All member scopes should be present in admin scopes
		for _, memberScope := range memberScopes {
			require.True(t, slices.Contains(adminScopes, memberScope),
				"Member scope %s should be present in admin scopes", memberScope)
		}
	})
}
