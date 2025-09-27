package models

import "gorm.io/gorm"

type OrganizationOnboardingData struct {
	StripeBusinessType string `gorm:"not null"`
	ResidingCountry string `gorm:"not null;size:2"`
	ContactEmail string
	ContactPhone string
	RegisteredBusinessName string
}

type Organization struct {
	gorm.Model

	// Basic information
	Name string `gorm:"not null;size:100"`
	Description string `gorm:"size:255;default:null"`
	LogoFileStorageKey string
	Address Address `gorm:"embedded;embeddedPrefix:address_"`

	// Localization fields
	Currency string `gorm:"not null;size:3"`
	Locale string `gorm:"not null;size:5"`

	// Onboarding fields
	OnboardingData OrganizationOnboardingData `gorm:"embedded;embeddedPrefix:onboarding_data_"`

	// Stripe fields
	StripeAccountID string
	
	// Relationships
	Memberships []OrganizationMembership `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}