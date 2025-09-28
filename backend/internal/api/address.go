package api

type Address struct {
	Line1 string `json:"line1" validate:"required"`
	Line2 string `json:"line2" validate:"omitempty"`
	City string `json:"city" validate:"required"`
	StateOrProvince string `json:"stateOrProvince" validate:"required"`
	Zip string `json:"zip" validate:"required"`
	Country string `json:"country" validate:"required"`
}