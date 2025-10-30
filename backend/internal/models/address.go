package models

// This is an embedded model for an address. Do not use this as a standalone model.
type Address struct {
	Line1           string `gorm:"not null"`
	Line2           string
	City            string `gorm:"not null"`
	StateOrProvince string `gorm:"not null"`
	Zip             string `gorm:"not null"`
	Country         string `gorm:"not null"` // ISO 3166-1 alpha-2 (2 letter country code)
}
