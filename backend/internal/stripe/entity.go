package stripe

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"
	stripeGo "github.com/stripe/stripe-go/v83"
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
    Type stripeGo.AccountBusinessType
	DisplayName string
	ContactEmail string
	ContactPhone string
    Currency stripeGo.Currency
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
    Event  *stripeGo.Event
	DB     *gorm.DB
	Config *configuration.Config
    StripeClient *Client
    Context context.Context
}

// EnqueueWebhookProcessingServiceRequest contains parameters for enqueueing webhook processing
type EnqueueWebhookProcessingServiceRequest struct {
    RiverClient *river.Client[*sql.Tx]
    Event       *stripeGo.Event
    Context     context.Context
}

type CreateOnboardingLinkServiceRequest struct {
	Context context.Context
	StripeClient *Client
    Db *gorm.DB
	Params CreateOnboardingLinkParams
}

type CreateOnboardingLinkParams struct {
	AccountID string
	RefreshURL string
	ReturnURL string
}