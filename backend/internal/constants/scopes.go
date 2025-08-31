package constants

type OrganizationScope string

const (
	OrganizationScopeRead OrganizationScope = "organization:read"
	OrganizationScopeUpdate OrganizationScope = "organization:update"
	OrganizationScopeDelete OrganizationScope = "organization:delete"
	OrganizationMembershipsScopeList OrganizationScope = "organization:memberships:list"
	OrganizationMembershipsScopeRead OrganizationScope = "organization:memberships:read"
	OrganizationMembershipsScopeCreate OrganizationScope = "organization:memberships:create"
	OrganizationMembershipsScopeUpdate OrganizationScope = "organization:memberships:update"
	OrganizationMembershipsScopeDelete OrganizationScope = "organization:memberships:delete"
	OrganizationInvitationsScopeList OrganizationScope = "organization:invitations:list"
	OrganizationInvitationsScopeRead OrganizationScope = "organization:invitations:read"
	OrganizationInvitationsScopeCreate OrganizationScope = "organization:invitations:create"
	OrganizationInvitationsScopeUpdate OrganizationScope = "organization:invitations:update"
	OrganizationInvitationsScopeDelete OrganizationScope = "organization:invitations:delete"
)