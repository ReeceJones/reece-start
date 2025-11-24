package mocks

import (
	"github.com/google/uuid"
	stripeGo "github.com/stripe/stripe-go/v83"
)

// NewMockStripeClient creates a new Stripe client that won't make actual API calls
// Note: ReplaceDefaultTransport() must be called before creating clients for this to work
func NewMockStripeClient() *stripeGo.Client {
	// Create a test API key (Stripe test keys start with sk_test_)
	// Using a clearly fake test key that won't work with real Stripe API
	testKey := "sk_test_mock_" + uuid.New().String()[:32]

	// Set the global key for package-level functions (like checkoutSession.New)
	stripeGo.Key = testKey

	// Create Stripe client - HTTP calls will be intercepted by MockHTTPTransport
	// (which should already be set up via ReplaceDefaultTransport())
	client := stripeGo.NewClient(testKey)

	return client
}

// CreateMockStripeAccount creates a mock Stripe Connect account for testing
func CreateMockStripeAccount(displayName string, metadata map[string]string) *stripeGo.V2CoreAccount {
	accountID := "acct_" + uuid.New().String()[:24]

	account := &stripeGo.V2CoreAccount{
		ID:          accountID,
		DisplayName: displayName,
		Metadata:    metadata,
		Identity: &stripeGo.V2CoreAccountIdentity{
			EntityType: stripeGo.V2CoreAccountIdentityEntityTypeIndividual,
			Country:    "US",
		},
		Defaults: &stripeGo.V2CoreAccountDefaults{
			Currency: stripeGo.CurrencyUSD,
		},
		Configuration: &stripeGo.V2CoreAccountConfiguration{
			Customer: &stripeGo.V2CoreAccountConfigurationCustomer{
				Capabilities: &stripeGo.V2CoreAccountConfigurationCustomerCapabilities{
					AutomaticIndirectTax: &stripeGo.V2CoreAccountConfigurationCustomerCapabilitiesAutomaticIndirectTax{
						Status: stripeGo.V2CoreAccountConfigurationCustomerCapabilitiesAutomaticIndirectTaxStatusActive,
					},
				},
			},
			Merchant: &stripeGo.V2CoreAccountConfigurationMerchant{
				Capabilities: &stripeGo.V2CoreAccountConfigurationMerchantCapabilities{
					CardPayments: &stripeGo.V2CoreAccountConfigurationMerchantCapabilitiesCardPayments{
						Status: stripeGo.V2CoreAccountConfigurationMerchantCapabilitiesCardPaymentsStatusActive,
					},
				},
			},
			Recipient: &stripeGo.V2CoreAccountConfigurationRecipient{
				Capabilities: &stripeGo.V2CoreAccountConfigurationRecipientCapabilities{
					StripeBalance: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalance{
						StripeTransfers: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfers{
							Status: stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfersStatusActive,
						},
					},
				},
			},
		},
		Requirements: &stripeGo.V2CoreAccountRequirements{},
	}

	return account
}
