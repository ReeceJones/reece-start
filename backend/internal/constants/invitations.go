package constants

type OrganizationInvitationStatus string

const (
	OrganizationInvitationStatusPending OrganizationInvitationStatus = "pending"
	OrganizationInvitationStatusAccepted OrganizationInvitationStatus = "accepted"
	OrganizationInvitationStatusDeclined OrganizationInvitationStatus = "declined"
	OrganizationInvitationStatusExpired OrganizationInvitationStatus = "expired"
	OrganizationInvitationStatusRevoked OrganizationInvitationStatus = "revoked"
)