package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"reece.start/internal/constants"
)

// A period of time that an organization is subscribed to a plan
type OrganizationPlanPeriod struct {
	gorm.Model
	OrganizationID       uuid.UUID                `gorm:"not null;type:uuid"`
	Plan                 constants.MembershipPlan `gorm:"not null;index"`
	StripeSubscriptionID string                   `gorm:"not null"`
	BillingPeriodStart   time.Time                `gorm:"not null;index"`
	BillingPeriodEnd     time.Time                `gorm:"not null;index"`
	BillingPeriodAmount  int                      `gorm:"not null"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}
