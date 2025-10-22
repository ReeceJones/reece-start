package constants

type UserRole string

type OrganizationRole string


const (
	UserRoleAdmin UserRole = "admin"
	UserRoleDefault UserRole = "default"
)

const (
	OrganizationRoleAdmin  OrganizationRole = "admin"
	OrganizationRoleMember OrganizationRole = "member"
)

// Mapping for role -> scopes
var UserRoleToScopes = map[UserRole][]UserScope{
	UserRoleAdmin: {
		UserScopeAdmin,
		UserScopeAdminUsersList,
		UserScopeAdminUsersRead,
		UserScopeAdminUsersImpersonate,
	},
	UserRoleDefault: {},
}

// Mapping for role -> scopes
var OrganizationRoleToScopes = map[OrganizationRole][]UserScope{
	// Grant all roles to admin
	OrganizationRoleAdmin: {
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
	},

	// Grant limited (mostly read scopes) to the member
	OrganizationRoleMember: {
		UserScopeOrganizationRead,
		UserScopeOrganizationMembershipsList,
		UserScopeOrganizationMembershipsRead,
		UserScopeOrganizationInvitationsList,
		UserScopeOrganizationInvitationsRead,
	},
}
