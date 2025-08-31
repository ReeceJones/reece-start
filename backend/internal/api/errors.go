package api

import "errors"

// Common business logic errors that can be returned from service layer transactions
// These errors can be checked using errors.Is() for better error handling
var (
	// Authentication/Authorization errors
	ErrForbiddenNoAccess        = errors.New("Unauthorized")
	ErrForbiddenNoAdminAccess   = errors.New("You don't have admin access to this organization")
	ErrForbiddenOwnProfileOnly  = errors.New("You can only update your own profile")
	ErrUnauthorizedInvalidLogin = errors.New("Invalid email or password")

	// Resource not found errors
	ErrMembershipNotFound      = errors.New("Membership not found")
	ErrInvitationNotFound      = errors.New("Invitation not found")
	ErrInvitationAlreadyExists = errors.New("An invitation already exists for this email and organization")
	ErrInvitationNotPending    = errors.New("Invitation is no longer pending")
	ErrInvitationEmailMismatch = errors.New("Invitation email does not match user email")
	ErrUserAlreadyMember       = errors.New("User is already a member of this organization")
	ErrUserNotFound            = errors.New("User not found")

	// Invalid ID errors
	ErrInvalidOrganizationID = errors.New("invalid organization ID")
	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrInvalidMembershipID   = errors.New("invalid membership ID")
	ErrInvalidInvitationID   = errors.New("invalid invitation ID")
)
