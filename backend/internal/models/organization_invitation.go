package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationInvitation struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email string `gorm:"not null"`
	Role string `gorm:"not null"`
	OrganizationID uint `gorm:"not null"`
	InvitingUserID uint `gorm:"not null"`
	Status string `gorm:"not null"`

	InvitingUser User `gorm:"foreignKey:InvitingUserID;constraint:OnDelete:CASCADE"`
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}