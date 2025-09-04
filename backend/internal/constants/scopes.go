package constants

type UserScope string

const (
	// Organization
	UserScopeOrganizationRead UserScope = "organization:read"
	UserScopeOrganizationUpdate UserScope = "organization:update"
	UserScopeOrganizationDelete UserScope = "organization:delete"
	UserScopeOrganizationMembershipsList UserScope = "organization:memberships:list"
	UserScopeOrganizationMembershipsRead UserScope = "organization:memberships:read"
	UserScopeOrganizationMembershipsCreate UserScope = "organization:memberships:create"
	UserScopeOrganizationMembershipsUpdate UserScope = "organization:memberships:update"
	UserScopeOrganizationMembershipsDelete UserScope = "organization:memberships:delete"
	UserScopeOrganizationInvitationsList UserScope = "organization:invitations:list"
	UserScopeOrganizationInvitationsRead UserScope = "organization:invitations:read"
	UserScopeOrganizationInvitationsCreate UserScope = "organization:invitations:create"
	UserScopeOrganizationInvitationsUpdate UserScope = "organization:invitations:update"
	UserScopeOrganizationInvitationsDelete UserScope = "organization:invitations:delete"

	// Admin
	UserScopeAdmin UserScope = "admin"
	UserScopeAdminUsersList UserScope = "admin:users:list"
	UserScopeAdminUsersRead UserScope = "admin:users:read"
	UserScopeAdminUsersImpersonate UserScope = "admin:users:impersonate"
)