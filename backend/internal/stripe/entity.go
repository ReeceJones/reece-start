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
	Line1           string
	Line2           string
	City            string
	StateOrProvince string
	Zip             string
	Country         string
}

type IndividualAccount struct {
	FirstName string
	LastName  string
}

type CompanyAccount struct {
	RegisteredName string
}

type CreateStripeAccountParams struct {
	OrganizationID  uint
	Type            stripeGo.AccountBusinessType
	DisplayName     string
	ContactEmail    string
	ContactPhone    string
	Currency        stripeGo.Currency
	Locale          string
	ResidingCountry string
	Address         Address
	Individual      IndividualAccount
	Company         CompanyAccount
}

type CreateStripeAccountServiceRequest struct {
	Context      context.Context
	Config       *configuration.Config
	StripeClient *stripeGo.Client
	Params       CreateStripeAccountParams
}

// ProcessSnapshotWebhookEventServiceRequest contains parameters for processing webhook events
type ProcessSnapshotWebhookEventServiceRequest struct {
	Event        *stripeGo.Event
	DB           *gorm.DB
	Config       *configuration.Config
	StripeClient *stripeGo.Client
	Context      context.Context
}

type ProcessThinWebhookEventServiceRequest struct {
	Event        stripeGo.EventNotificationContainer
	DB           *gorm.DB
	Config       *configuration.Config
	StripeClient *stripeGo.Client
	Context      context.Context
}

type FetchAndUpdateAccountServiceRequest struct {
	AccountID    string
	DB           *gorm.DB
	Config       *configuration.Config
	StripeClient *stripeGo.Client
	Context      context.Context
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
	Context      context.Context
	StripeClient *stripeGo.Client
	Db           *gorm.DB
	Params       CreateOnboardingLinkParams
}

type CreateOnboardingLinkParams struct {
	AccountID  string
	RefreshURL string
	ReturnURL  string
}

type CreateCheckoutSessionServiceRequest struct {
	Context      context.Context
	Config       *configuration.Config
	StripeClient *stripeGo.Client
	DB           *gorm.DB
	Params       CreateCheckoutSessionParams
}

type CreateCheckoutSessionParams struct {
	OrganizationID uint
	SuccessURL     string
	CancelURL      string
}

type CreateBillingPortalSessionServiceRequest struct {
	Context      context.Context
	Config       *configuration.Config
	StripeClient *stripeGo.Client
	DB           *gorm.DB
	Params       CreateBillingPortalSessionParams
}

type CreateBillingPortalSessionParams struct {
	OrganizationID uint
	ReturnURL      string
}

type GetSubscriptionServiceRequest struct {
	Context        context.Context
	DB             *gorm.DB
	OrganizationID uint
}

// Response structs for HTTP endpoints
type CheckoutSessionResponse struct {
	Data CheckoutSessionData `json:"data"`
}

type CheckoutSessionData struct {
	Type       string                    `json:"type"`
	ID         string                    `json:"id"`
	Attributes CheckoutSessionAttributes `json:"attributes"`
}

type CheckoutSessionAttributes struct {
	URL string `json:"url"`
}

type BillingPortalSessionResponse struct {
	Data BillingPortalSessionData `json:"data"`
}

type BillingPortalSessionData struct {
	Type       string                         `json:"type"`
	ID         string                         `json:"id"`
	Attributes BillingPortalSessionAttributes `json:"attributes"`
}

type BillingPortalSessionAttributes struct {
	URL string `json:"url"`
}

type SubscriptionResponse struct {
	Data SubscriptionData `json:"data"`
}

type SubscriptionData struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id,omitempty"`
	Attributes SubscriptionAttributes `json:"attributes"`
}

type SubscriptionAttributes struct {
	Plan               string  `json:"plan"`
	BillingPeriodStart *string `json:"billingPeriodStart"`
	BillingPeriodEnd   *string `json:"billingPeriodEnd"`
	BillingAmount      int     `json:"billingAmount"`
}

// Request structs for HTTP endpoints
type CreateCheckoutSessionRequest struct {
	SuccessURL string `json:"successUrl" validate:"required,url"`
	CancelURL  string `json:"cancelUrl" validate:"required,url"`
}

type CreateBillingPortalSessionRequest struct {
	ReturnURL string `json:"returnUrl" validate:"required,url"`
}
