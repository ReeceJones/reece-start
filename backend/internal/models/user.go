package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `gorm:"not null"`
	Email          string `gorm:"index:idx_email,unique;not null"`
	HashedPassword []byte `gorm:"not null"`
	
	// Relationships
	OrganizationMemberships []OrganizationMembership `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
