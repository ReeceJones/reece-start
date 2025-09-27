package stripe

import (
	"github.com/riverqueue/river"
	"github.com/stripe/stripe-go/v82"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
)

type Address struct {
	Line1 string
	Line2 string
	City string
	StateOrProvince string
	Zip string
	Country string
}

type IndividualAccount struct {
	Email string
	Phone string
	FirstName string
	LastName string
}

type CompanyAccount struct {
	Phone string
	RegisteredName string
	Url string
	Structure string
}

type CreateStripeAccountParams struct {
	OrganizationID uint
	Type stripe.AccountBusinessType
	DisplayName string
	ContactEmail string
	Currency string
	Locale string
	ResidingCountry string
	Address Address
	Individual IndividualAccount
	Company CompanyAccount
}

type CreateStripeAccountServiceRequest struct {
	Config *configuration.Config
	StripeClient *stripe.Client
	Params CreateStripeAccountParams
}

// ProcessWebhookEventServiceRequest contains parameters for processing webhook events
type ProcessWebhookEventServiceRequest struct {
	Event  *stripe.Event
	DB     *gorm.DB
	Config *configuration.Config
}

// EnqueueWebhookProcessingServiceRequest contains parameters for enqueueing webhook processing
type EnqueueWebhookProcessingServiceRequest struct {
	RiverClient *river.Client[river.JobArgs]
	Event       *stripe.Event
}