package stripe

import (
	"context"

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
	FirstName string
	LastName string
}

type CompanyAccount struct {
	RegisteredName string
}

type CreateStripeAccountParams struct {
	OrganizationID uint
	Type stripe.AccountBusinessType
	DisplayName string
	ContactEmail string
	ContactPhone string
	Currency string
	Locale string
	ResidingCountry string
	Address Address
	Individual IndividualAccount
	Company CompanyAccount
}

type CreateStripeAccountServiceRequest struct {
	Context context.Context
	Config *configuration.Config
	StripeClient *Client
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