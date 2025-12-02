package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationMembership struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	Role           string    `gorm:"not null;size:20;default:'member'"`

	// Relationships
	User         User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}
