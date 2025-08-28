package models

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	Name string `gorm:"not null;size:100"`
	Description string `gorm:"size:255;default:null"`
	LogoFileStorageKey string
	
	// Relationships
	Memberships []OrganizationMembership `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}