package utils

import (
	stripeGo "github.com/stripe/stripe-go/v83"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

// DetermineStripeOnboardingStatus computes the Stripe onboarding status for the given organization
// based on its stored capability statuses and requirements.
func DetermineStripeOnboardingStatus(organization *models.Organization) constants.StripeOnboardingStatus {
    if organization.Stripe.HasPendingRequirements {
        return constants.StripeOnboardingStatusMissingRequirements
    }

    if organization.Stripe.AutomaticIndirectTaxStatus != string(stripeGo.AccountCapabilityStatusActive) ||
        organization.Stripe.CardPaymentsStatus != string(stripeGo.AccountCapabilityStatusActive) ||
        organization.Stripe.StripeBalancePayoutsStatus != string(stripeGo.AccountCapabilityStatusActive) ||
        organization.Stripe.StripeBalanceTransfersStatus != string(stripeGo.AccountCapabilityStatusActive) {
        return constants.StripeOnboardingStatusMissingCapabilities
    }

    return constants.StripeOnboardingStatusCompleted
}

// ApplyStripeAccountToOrganization updates the organization's embedded Stripe fields
// based on the latest Stripe V2 account object and recomputes onboarding status.
func ApplyStripeAccountToOrganization(organization *models.Organization, account *stripeGo.V2CoreAccount) {
    organization.Stripe.AccountID = account.ID
    if account.Configuration != nil {
        if account.Configuration.Customer != nil && account.Configuration.Customer.Capabilities != nil && account.Configuration.Customer.Capabilities.AutomaticIndirectTax != nil {
            organization.Stripe.AutomaticIndirectTaxStatus = string(account.Configuration.Customer.Capabilities.AutomaticIndirectTax.Status)
        }

        if account.Configuration.Merchant != nil && account.Configuration.Merchant.Capabilities != nil && account.Configuration.Merchant.Capabilities.CardPayments != nil {
            organization.Stripe.CardPaymentsStatus = string(account.Configuration.Merchant.Capabilities.CardPayments.Status)
        }

        if account.Configuration.Recipient != nil && account.Configuration.Recipient.Capabilities != nil && account.Configuration.Recipient.Capabilities.StripeBalance != nil {
            if account.Configuration.Recipient.Capabilities.StripeBalance.Payouts != nil {
                organization.Stripe.StripeBalancePayoutsStatus = string(account.Configuration.Recipient.Capabilities.StripeBalance.Payouts.Status)
            }
            if account.Configuration.Recipient.Capabilities.StripeBalance.StripeTransfers != nil {
                organization.Stripe.StripeBalanceTransfersStatus = string(account.Configuration.Recipient.Capabilities.StripeBalance.StripeTransfers.Status)
            }
        }
    }

    if account.Requirements != nil {
        organization.Stripe.HasPendingRequirements = len(account.Requirements.Entries) > 0
    }

    organization.Stripe.OnboardingStatus = string(DetermineStripeOnboardingStatus(organization))
}


