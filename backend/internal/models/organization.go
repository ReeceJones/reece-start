package models

import "gorm.io/gorm"

type OrganizationStripeAccount struct {
	// Account ID
	AccountID string `gorm:"index"`

	// Status flags
	AutomaticIndirectTaxStatus string
	CardPaymentsStatus string
	StripeBalancePayoutsStatus string
	StripeBalanceTransfersStatus string

	// Requirements
	HasPendingRequirements bool

	// Onboarding status
	OnboardingStatus string
}

type Organization struct {
	gorm.Model

	// Basic information
	Name string `gorm:"not null;size:100"`
	Description string `gorm:"size:255;default:null"`
	LogoFileStorageKey string
	Address Address `gorm:"embedded;embeddedPrefix:address_"`

	// Contact information
	ContactEmail string
	ContactPhone string
	ContactPhoneCountry string

	// Localization fields
	Currency string `gorm:"not null;size:3"`
	Locale string `gorm:"not null;size:5"`

	// Stripe fields
	Stripe OrganizationStripeAccount `gorm:"embedded;embeddedPrefix:stripe_"`

	// Onboarding status
	OnboardingStatus string
	
	// Relationships
	Memberships []OrganizationMembership `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}