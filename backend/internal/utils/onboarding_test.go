package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	stripeGo "github.com/stripe/stripe-go/v83"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

func TestDetermineStripeOnboardingStatus(t *testing.T) {
	t.Run("MissingRequirements", func(t *testing.T) {
		organization := &models.Organization{
			Stripe: models.OrganizationStripeAccount{
				HasPendingRequirements:       true,
				AutomaticIndirectTaxStatus:   string(stripeGo.AccountCapabilityStatusActive),
				CardPaymentsStatus:           string(stripeGo.AccountCapabilityStatusActive),
				StripeBalancePayoutsStatus:   string(stripeGo.AccountCapabilityStatusActive),
				StripeBalanceTransfersStatus: string(stripeGo.AccountCapabilityStatusActive),
			},
		}

		status := DetermineStripeOnboardingStatus(organization)
		assert.Equal(t, constants.StripeOnboardingStatusMissingRequirements, status)
	})

	t.Run("MissingCapabilities", func(t *testing.T) {
		tests := []struct {
			name         string
			organization *models.Organization
		}{
			{
				name: "missing automatic indirect tax",
				organization: &models.Organization{
					Stripe: models.OrganizationStripeAccount{
						HasPendingRequirements:       false,
						AutomaticIndirectTaxStatus:   string(stripeGo.AccountCapabilityStatusPending),
						CardPaymentsStatus:           string(stripeGo.AccountCapabilityStatusActive),
						StripeBalancePayoutsStatus:   string(stripeGo.AccountCapabilityStatusActive),
						StripeBalanceTransfersStatus: string(stripeGo.AccountCapabilityStatusActive),
					},
				},
			},
			{
				name: "missing card payments",
				organization: &models.Organization{
					Stripe: models.OrganizationStripeAccount{
						HasPendingRequirements:       false,
						AutomaticIndirectTaxStatus:   string(stripeGo.AccountCapabilityStatusActive),
						CardPaymentsStatus:           string(stripeGo.AccountCapabilityStatusPending),
						StripeBalancePayoutsStatus:   string(stripeGo.AccountCapabilityStatusActive),
						StripeBalanceTransfersStatus: string(stripeGo.AccountCapabilityStatusActive),
					},
				},
			},
			{
				name: "missing stripe balance payouts",
				organization: &models.Organization{
					Stripe: models.OrganizationStripeAccount{
						HasPendingRequirements:       false,
						AutomaticIndirectTaxStatus:   string(stripeGo.AccountCapabilityStatusActive),
						CardPaymentsStatus:           string(stripeGo.AccountCapabilityStatusActive),
						StripeBalancePayoutsStatus:   string(stripeGo.AccountCapabilityStatusPending),
						StripeBalanceTransfersStatus: string(stripeGo.AccountCapabilityStatusActive),
					},
				},
			},
			{
				name: "missing stripe balance transfers",
				organization: &models.Organization{
					Stripe: models.OrganizationStripeAccount{
						HasPendingRequirements:       false,
						AutomaticIndirectTaxStatus:   string(stripeGo.AccountCapabilityStatusActive),
						CardPaymentsStatus:           string(stripeGo.AccountCapabilityStatusActive),
						StripeBalancePayoutsStatus:   string(stripeGo.AccountCapabilityStatusActive),
						StripeBalanceTransfersStatus: string(stripeGo.AccountCapabilityStatusPending),
					},
				},
			},
			{
				name: "missing multiple capabilities",
				organization: &models.Organization{
					Stripe: models.OrganizationStripeAccount{
						HasPendingRequirements:       false,
						AutomaticIndirectTaxStatus:   string(stripeGo.AccountCapabilityStatusPending),
						CardPaymentsStatus:           string(stripeGo.AccountCapabilityStatusPending),
						StripeBalancePayoutsStatus:   string(stripeGo.AccountCapabilityStatusActive),
						StripeBalanceTransfersStatus: string(stripeGo.AccountCapabilityStatusActive),
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				status := DetermineStripeOnboardingStatus(tt.organization)
				assert.Equal(t, constants.StripeOnboardingStatusMissingCapabilities, status)
			})
		}
	})

	t.Run("Completed", func(t *testing.T) {
		organization := &models.Organization{
			Stripe: models.OrganizationStripeAccount{
				HasPendingRequirements:       false,
				AutomaticIndirectTaxStatus:   string(stripeGo.AccountCapabilityStatusActive),
				CardPaymentsStatus:           string(stripeGo.AccountCapabilityStatusActive),
				StripeBalancePayoutsStatus:   string(stripeGo.AccountCapabilityStatusActive),
				StripeBalanceTransfersStatus: string(stripeGo.AccountCapabilityStatusActive),
			},
		}

		status := DetermineStripeOnboardingStatus(organization)
		assert.Equal(t, constants.StripeOnboardingStatusCompleted, status)
	})
}

func TestApplyStripeAccountToOrganization(t *testing.T) {
	t.Run("CompleteAccount", func(t *testing.T) {
		organization := &models.Organization{
			Stripe: models.OrganizationStripeAccount{},
		}

		account := &stripeGo.V2CoreAccount{
			ID: "acct_test123",
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
							Payouts: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalancePayouts{
								Status: stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalancePayoutsStatusActive,
							},
							StripeTransfers: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfers{
								Status: stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfersStatusActive,
							},
						},
					},
				},
			},
			Requirements: &stripeGo.V2CoreAccountRequirements{
				Entries: []*stripeGo.V2CoreAccountRequirementsEntry{},
			},
		}

		ApplyStripeAccountToOrganization(organization, account)

		assert.Equal(t, "acct_test123", organization.Stripe.AccountID)
		assert.Equal(t, string(stripeGo.AccountCapabilityStatusActive), organization.Stripe.AutomaticIndirectTaxStatus)
		assert.Equal(t, string(stripeGo.AccountCapabilityStatusActive), organization.Stripe.CardPaymentsStatus)
		assert.Equal(t, string(stripeGo.AccountCapabilityStatusActive), organization.Stripe.StripeBalancePayoutsStatus)
		assert.Equal(t, string(stripeGo.AccountCapabilityStatusActive), organization.Stripe.StripeBalanceTransfersStatus)
		assert.False(t, organization.Stripe.HasPendingRequirements)
		assert.Equal(t, string(constants.StripeOnboardingStatusCompleted), organization.Stripe.OnboardingStatus)
	})

	t.Run("WithPendingRequirements", func(t *testing.T) {
		organization := &models.Organization{
			Stripe: models.OrganizationStripeAccount{},
		}

		account := &stripeGo.V2CoreAccount{
			ID: "acct_test456",
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
							Payouts: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalancePayouts{
								Status: stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalancePayoutsStatusActive,
							},
							StripeTransfers: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfers{
								Status: stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfersStatusActive,
							},
						},
					},
				},
			},
			Requirements: &stripeGo.V2CoreAccountRequirements{
				Entries: []*stripeGo.V2CoreAccountRequirementsEntry{
					{Description: "individual.dob"},
				},
			},
		}

		ApplyStripeAccountToOrganization(organization, account)

		assert.Equal(t, "acct_test456", organization.Stripe.AccountID)
		assert.True(t, organization.Stripe.HasPendingRequirements)
		assert.Equal(t, string(constants.StripeOnboardingStatusMissingRequirements), organization.Stripe.OnboardingStatus)
	})

	t.Run("WithMissingCapabilities", func(t *testing.T) {
		organization := &models.Organization{
			Stripe: models.OrganizationStripeAccount{},
		}

		account := &stripeGo.V2CoreAccount{
			ID: "acct_test789",
			Configuration: &stripeGo.V2CoreAccountConfiguration{
				Customer: &stripeGo.V2CoreAccountConfigurationCustomer{
					Capabilities: &stripeGo.V2CoreAccountConfigurationCustomerCapabilities{
						AutomaticIndirectTax: &stripeGo.V2CoreAccountConfigurationCustomerCapabilitiesAutomaticIndirectTax{
							Status: stripeGo.V2CoreAccountConfigurationCustomerCapabilitiesAutomaticIndirectTaxStatusPending,
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
							Payouts: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalancePayouts{
								Status: stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalancePayoutsStatusActive,
							},
							StripeTransfers: &stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfers{
								Status: stripeGo.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfersStatusActive,
							},
						},
					},
				},
			},
			Requirements: &stripeGo.V2CoreAccountRequirements{
				Entries: []*stripeGo.V2CoreAccountRequirementsEntry{},
			},
		}

		ApplyStripeAccountToOrganization(organization, account)

		assert.Equal(t, "acct_test789", organization.Stripe.AccountID)
		assert.Equal(t, string(stripeGo.AccountCapabilityStatusPending), organization.Stripe.AutomaticIndirectTaxStatus)
		assert.False(t, organization.Stripe.HasPendingRequirements)
		assert.Equal(t, string(constants.StripeOnboardingStatusMissingCapabilities), organization.Stripe.OnboardingStatus)
	})

	t.Run("NilConfiguration", func(t *testing.T) {
		organization := &models.Organization{
			Stripe: models.OrganizationStripeAccount{
				AccountID: "acct_old",
			},
		}

		account := &stripeGo.V2CoreAccount{
			ID:            "acct_new",
			Configuration: nil,
			Requirements: &stripeGo.V2CoreAccountRequirements{
				Entries: []*stripeGo.V2CoreAccountRequirementsEntry{},
			},
		}

		ApplyStripeAccountToOrganization(organization, account)

		assert.Equal(t, "acct_new", organization.Stripe.AccountID)
		assert.False(t, organization.Stripe.HasPendingRequirements)
		// Capabilities should remain empty/default since Configuration is nil
		assert.Equal(t, string(constants.StripeOnboardingStatusMissingCapabilities), organization.Stripe.OnboardingStatus)
	})

	t.Run("NilRequirements", func(t *testing.T) {
		organization := &models.Organization{
			Stripe: models.OrganizationStripeAccount{
				HasPendingRequirements: true,
			},
		}

		account := &stripeGo.V2CoreAccount{
			ID:            "acct_test",
			Configuration: &stripeGo.V2CoreAccountConfiguration{},
			Requirements:  nil,
		}

		ApplyStripeAccountToOrganization(organization, account)

		assert.Equal(t, "acct_test", organization.Stripe.AccountID)
		// HasPendingRequirements should remain unchanged when Requirements is nil
		assert.True(t, organization.Stripe.HasPendingRequirements)
	})
}
