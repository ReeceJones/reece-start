package models

import "gorm.io/gorm"

type OrganizationMembership struct {
	gorm.Model
	UserID         uint `gorm:"not null;index"`
	OrganizationID uint `gorm:"not null;index"`
	Role           string `gorm:"not null;size:20;default:'member'"`
	
	// Relationships
	User         User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}