package constants

type OrganizationRole string

const (
	OrganizationRoleAdmin  OrganizationRole = "admin"
	OrganizationRoleMember OrganizationRole = "member"
)

// Mapping for role -> scopes
var OrganizationRoleToScopes = map[OrganizationRole][]OrganizationScope{
	// Grant all roles to admin
	OrganizationRoleAdmin: {
		OrganizationScopeRead,
		OrganizationScopeUpdate,
		OrganizationScopeDelete,
		OrganizationMembershipsScopeList,
		OrganizationMembershipsScopeRead,
		OrganizationMembershipsScopeCreate,
		OrganizationMembershipsScopeUpdate,
		OrganizationMembershipsScopeDelete,
		OrganizationInvitationsScopeList,
		OrganizationInvitationsScopeRead,
		OrganizationInvitationsScopeCreate,
		OrganizationInvitationsScopeUpdate,
		OrganizationInvitationsScopeDelete,
	},

	// Grant limited (mostly read scopes) to the member
	OrganizationRoleMember: {
		OrganizationScopeRead,
		OrganizationMembershipsScopeList,
		OrganizationMembershipsScopeRead,
		OrganizationInvitationsScopeList,
		OrganizationInvitationsScopeRead,
	},
}
