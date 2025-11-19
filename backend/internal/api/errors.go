package api

import (
	"errors"
	"strings"

	"github.com/lib/pq"
)

// Common business logic errors that can be returned from service layer transactions
// These errors can be checked using errors.Is() for better error handling
var (
	// Authentication/Authorization errors
	ErrForbiddenNoAccess                = errors.New("unauthorized")
	ErrForbiddenNoAdminAccess           = errors.New("you don't have admin access to this organization")
	ErrForbiddenOwnProfileOnly          = errors.New("you can only update your own profile")
	ErrUnauthorizedInvalidLogin         = errors.New("invalid email or password")
	ErrForbiddenImpersonationNotAllowed = errors.New("impersonation is not allowed")
	ErrMissingAuthorizationHeader       = errors.New("missing authorization header")
	ErrInvalidAuthorizationFormat       = errors.New("invalid authorization format")
	ErrInvalidToken                     = errors.New("invalid token")

	// Resource not found errors
	ErrMembershipNotFound      = errors.New("membership not found")
	ErrInvitationNotFound      = errors.New("invitation not found")
	ErrInvitationAlreadyExists = errors.New("an invitation already exists for this email and organization")
	ErrInvitationNotPending    = errors.New("invitation is no longer pending")
	ErrInvitationEmailMismatch = errors.New("invitation email does not match user email")
	ErrUserAlreadyMember       = errors.New("user is already a member of this organization")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserEmailAlreadyExists  = errors.New("a user with this email already exists")

	// Invalid ID errors
	ErrInvalidOrganizationID = errors.New("invalid organization id")
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrInvalidMembershipID   = errors.New("invalid membership id")
	ErrInvalidInvitationID   = errors.New("invalid invitation id")

	// Stripe webhook errors
	ErrStripeWebhookSecretNotConfigured = errors.New("stripe webhook secret not configured")
	ErrStripeWebhookSignatureMissing    = errors.New("stripe webhook signature missing")
	ErrStripeWebhookSignatureInvalid    = errors.New("stripe webhook signature invalid")
	ErrStripeWebhookEventInvalid        = errors.New("stripe webhook event invalid")
	ErrStripeWebhookEventUnhandled      = errors.New("stripe webhook event unhandled")
)

// IsUniqueConstraintViolation checks if an error is a PostgreSQL unique constraint violation
// It checks for PostgreSQL error code 23505 (unique_violation)
func IsUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}

	// Try to unwrap to pq.Error
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // unique_violation
	}

	// Fallback: check error message for common patterns
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "unique constraint") ||
		strings.Contains(errMsg, "violates unique constraint")
}
