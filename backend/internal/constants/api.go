package constants

type ApiType string

const (
	ApiTypeUser ApiType = "user"
	ApiTypeOrganization ApiType = "organization"
	ApiTypeOrganizationMembership ApiType = "organization-membership"
	ApiTypeOrganizationInvitation ApiType = "organization-invitation"
)
