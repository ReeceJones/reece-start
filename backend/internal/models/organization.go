package models

import "gorm.io/gorm"

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
	WebsiteUrl string

	// Localization fields
	Currency string `gorm:"not null;size:3"`
	Locale string `gorm:"not null;size:5"`

	// Stripe fields
	StripeAccountID string `gorm:"index"`
	
	// Relationships
	Memberships []OrganizationMembership `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}