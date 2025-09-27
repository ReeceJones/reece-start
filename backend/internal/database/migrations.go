package database

import (
	"gorm.io/gorm"
	"reece.start/internal/models"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.OrganizationMembership{},
		&models.OrganizationInvitation{},
		&models.OrganizationPlanPeriod{},
	)
}
