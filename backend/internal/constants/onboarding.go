package constants

type StripeOnboardingStatus string
type OnboardingStatus string

const (
	StripeOnboardingStatusPending             StripeOnboardingStatus = "pending"
	StripeOnboardingStatusCompleted           StripeOnboardingStatus = "completed"
	StripeOnboardingStatusMissingRequirements StripeOnboardingStatus = "missing_requirements"
	StripeOnboardingStatusMissingCapabilities StripeOnboardingStatus = "missing_capabilities"
)

const (
	OnboardingStatusPending    OnboardingStatus = "pending"
	OnboardingStatusCompleted  OnboardingStatus = "completed"
	OnboardingStatusInProgress OnboardingStatus = "in_progress"
)
