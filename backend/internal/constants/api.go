package constants

type ApiType string

const (
	ApiTypeUser ApiType = "user"
	ApiTypeOrganization ApiType = "organization"
	ApiTypeToken ApiType = "token"
	ApiTypeOrganizationMembership ApiType = "organization-membership"
	ApiTypeOrganizationInvitation ApiType = "organization-invitation"
	ApiTypeStripeAccountLink ApiType = "stripe-account-link"
)
