package models

import "gorm.io/gorm"

type OrganizationInvitation struct {
	gorm.Model
	Email string `gorm:"not null"`
	InvitationToken string `gorm:"not null"`
	OrganizationID uint `gorm:"not null"`
	InvitingUserID uint `gorm:"not null"`

	InvitingUser User `gorm:"foreignKey:InvitingUserID;constraint:OnDelete:CASCADE"`
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}