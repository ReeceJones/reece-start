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
	StripeClient *stripeGo.Client
	Params CreateStripeAccountParams
}

// ProcessSnapshotWebhookEventServiceRequest contains parameters for processing webhook events
type ProcessSnapshotWebhookEventServiceRequest struct {
    Event  *stripeGo.Event
	DB     *gorm.DB
	Config *configuration.Config
    StripeClient *stripeGo.Client
    Context context.Context
}

type ProcessThinWebhookEventServiceRequest struct {
    Event  stripeGo.EventNotificationContainer
	DB     *gorm.DB
	Config *configuration.Config
    StripeClient *stripeGo.Client
    Context context.Context
}

type FetchAndUpdateAccountServiceRequest struct {
    AccountID string
	DB     *gorm.DB
	Config *configuration.Config
    StripeClient *stripeGo.Client
    Context context.Context
}

// EnqueueSnapshotWebhookProcessingServiceRequest contains parameters for enqueueing webhook processing
type EnqueueSnapshotWebhookProcessingServiceRequest struct {
    RiverClient *river.Client[*sql.Tx]
    Event       *stripeGo.Event
    Context     context.Context
}

type EnqueueThinWebhookProcessingServiceRequest struct {
    RiverClient *river.Client[*sql.Tx]
    Event       stripeGo.EventNotificationContainer
    Context     context.Context
}

type CreateOnboardingLinkServiceRequest struct {
	Context context.Context
	StripeClient *stripeGo.Client
    Db *gorm.DB
	Params CreateOnboardingLinkParams
}

type CreateOnboardingLinkParams struct {
	AccountID string
	RefreshURL string
	ReturnURL string
}