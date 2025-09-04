package models

import (
	"time"

	"gorm.io/gorm"
)

// Simple struct to track token revocation and refreshability.
type UserTokenRevocation struct {
	// Any token with an `iat` older than this should be considered invalid
	LastValidIssuedAt *time.Time

	// If true, the frontend can automatically refresh the token.
	// If false, the frontend needs to flush the token from cookies and force the user to re-authenticate.
	CanRefresh bool `gorm:"not null;default:true"`
}

type User struct {
	gorm.Model
	Name           string `gorm:"not null"`
	Email          string `gorm:"index:idx_email,unique;not null"`
	HashedPassword []byte `gorm:"not null"`
	LogoFileStorageKey string

	// Control fields
	Revocation UserTokenRevocation `gorm:"embedded;embeddedPrefix:revocation_"`

	// Admin fields
	Role string `gorm:"not null;size:20;default:'default'"`
	
	// Relationships
	OrganizationMemberships []OrganizationMembership `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
