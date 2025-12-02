package stripe

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reece.start/internal/constants"
)

func TestCheckoutSessionResponse(t *testing.T) {
	t.Run("maps checkout session to response correctly", func(t *testing.T) {
		sessionID := "cs_test_" + uuid.New().String()[:24]
		sessionURL := "https://checkout.stripe.com/test/" + sessionID

		response := CheckoutSessionResponse{
			Data: CheckoutSessionData{
				Type: "checkout-session",
				ID:   sessionID,
				Attributes: CheckoutSessionAttributes{
					URL: sessionURL,
				},
			},
		}

		assert.Equal(t, "checkout-session", response.Data.Type)
		assert.Equal(t, sessionID, response.Data.ID)
		assert.Equal(t, sessionURL, response.Data.Attributes.URL)
	})
}

func TestBillingPortalSessionResponse(t *testing.T) {
	t.Run("maps billing portal session to response correctly", func(t *testing.T) {
		sessionID := "bps_test_" + uuid.New().String()[:24]
		sessionURL := "https://billing.stripe.com/test/" + sessionID

		response := BillingPortalSessionResponse{
			Data: BillingPortalSessionData{
				Type: "billing-portal-session",
				ID:   sessionID,
				Attributes: BillingPortalSessionAttributes{
					URL: sessionURL,
				},
			},
		}

		assert.Equal(t, "billing-portal-session", response.Data.Type)
		assert.Equal(t, sessionID, response.Data.ID)
		assert.Equal(t, sessionURL, response.Data.Attributes.URL)
	})
}

func TestSubscriptionResponse(t *testing.T) {
	t.Run("maps subscription to response correctly", func(t *testing.T) {
		planPeriodID := uint(123)
		billingPeriodStart := time.Now()
		billingPeriodEnd := billingPeriodStart.AddDate(0, 1, 0)
		billingPeriodStartStr := billingPeriodStart.Format(time.RFC3339)
		billingPeriodEndStr := billingPeriodEnd.Format(time.RFC3339)

		response := SubscriptionResponse{
			Data: SubscriptionData{
				Type: "subscription",
				ID:   strconv.FormatUint(uint64(planPeriodID), 10),
				Attributes: SubscriptionAttributes{
					Plan:               string(constants.MembershipPlanPro),
					BillingPeriodStart: &billingPeriodStartStr,
					BillingPeriodEnd:   &billingPeriodEndStr,
					BillingAmount:      1000,
				},
			},
		}

		assert.Equal(t, "subscription", response.Data.Type)
		assert.Equal(t, strconv.FormatUint(uint64(planPeriodID), 10), response.Data.ID)
		assert.Equal(t, string(constants.MembershipPlanPro), response.Data.Attributes.Plan)
		assert.NotNil(t, response.Data.Attributes.BillingPeriodStart)
		assert.NotNil(t, response.Data.Attributes.BillingPeriodEnd)
		assert.Equal(t, 1000, response.Data.Attributes.BillingAmount)
	})

	t.Run("maps free plan subscription to response correctly", func(t *testing.T) {
		response := SubscriptionResponse{
			Data: SubscriptionData{
				Type: "subscription",
				Attributes: SubscriptionAttributes{
					Plan:               string(constants.MembershipPlanFree),
					BillingPeriodStart: nil,
					BillingPeriodEnd:   nil,
					BillingAmount:      0,
				},
			},
		}

		assert.Equal(t, "subscription", response.Data.Type)
		assert.Empty(t, response.Data.ID)
		assert.Equal(t, string(constants.MembershipPlanFree), response.Data.Attributes.Plan)
		assert.Nil(t, response.Data.Attributes.BillingPeriodStart)
		assert.Nil(t, response.Data.Attributes.BillingPeriodEnd)
		assert.Equal(t, 0, response.Data.Attributes.BillingAmount)
	})
}

func TestCreateStripeAccountParams(t *testing.T) {
	t.Run("creates params correctly", func(t *testing.T) {
		params := CreateStripeAccountParams{
			OrganizationID:  uuid.New(),
			Type:            "individual",
			DisplayName:     "Test Account",
			ContactEmail:    "test@example.com",
			ContactPhone:    "1234567890",
			Currency:        "usd",
			Locale:          "en-US",
			ResidingCountry: "US",
			Address: Address{
				Line1:           "123 Test St",
				Line2:           "Apt 4",
				City:            "Test City",
				StateOrProvince: "CA",
				Zip:             "12345",
				Country:         "US",
			},
			Individual: IndividualAccount{
				FirstName: "John",
				LastName:  "Doe",
			},
		}

		assert.Equal(t, params.OrganizationID, params.OrganizationID)
		assert.Equal(t, "individual", string(params.Type))
		assert.Equal(t, "Test Account", params.DisplayName)
		assert.Equal(t, "test@example.com", params.ContactEmail)
		assert.Equal(t, "123 Test St", params.Address.Line1)
		assert.Equal(t, "John", params.Individual.FirstName)
	})
}

func TestCreateOnboardingLinkParams(t *testing.T) {
	t.Run("creates params correctly", func(t *testing.T) {
		accountID := "acct_test_" + uuid.New().String()[:24]
		params := CreateOnboardingLinkParams{
			AccountID:  accountID,
			RefreshURL: "https://example.com/refresh",
			ReturnURL:  "https://example.com/return",
		}

		assert.Equal(t, accountID, params.AccountID)
		assert.Equal(t, "https://example.com/refresh", params.RefreshURL)
		assert.Equal(t, "https://example.com/return", params.ReturnURL)
	})
}

func TestCreateCheckoutSessionParams(t *testing.T) {
	t.Run("creates params correctly", func(t *testing.T) {
		params := CreateCheckoutSessionParams{
			OrganizationID: uuid.New(),
			SuccessURL:     "https://example.com/success",
			CancelURL:      "https://example.com/cancel",
		}

		assert.Equal(t, params.OrganizationID, params.OrganizationID)
		assert.Equal(t, "https://example.com/success", params.SuccessURL)
		assert.Equal(t, "https://example.com/cancel", params.CancelURL)
	})
}

func TestCreateBillingPortalSessionParams(t *testing.T) {
	t.Run("creates params correctly", func(t *testing.T) {
		params := CreateBillingPortalSessionParams{
			OrganizationID: uuid.New(),
			ReturnURL:      "https://example.com/return",
		}

		assert.Equal(t, params.OrganizationID, params.OrganizationID)
		assert.Equal(t, "https://example.com/return", params.ReturnURL)
	})
}

func TestAddress(t *testing.T) {
	t.Run("creates address correctly", func(t *testing.T) {
		address := Address{
			Line1:           "123 Test St",
			Line2:           "Apt 4",
			City:            "Test City",
			StateOrProvince: "CA",
			Zip:             "12345",
			Country:         "US",
		}

		assert.Equal(t, "123 Test St", address.Line1)
		assert.Equal(t, "Apt 4", address.Line2)
		assert.Equal(t, "Test City", address.City)
		assert.Equal(t, "CA", address.StateOrProvince)
		assert.Equal(t, "12345", address.Zip)
		assert.Equal(t, "US", address.Country)
	})
}

func TestIndividualAccount(t *testing.T) {
	t.Run("creates individual account correctly", func(t *testing.T) {
		account := IndividualAccount{
			FirstName: "John",
			LastName:  "Doe",
		}

		assert.Equal(t, "John", account.FirstName)
		assert.Equal(t, "Doe", account.LastName)
	})
}

func TestCompanyAccount(t *testing.T) {
	t.Run("creates company account correctly", func(t *testing.T) {
		account := CompanyAccount{
			RegisteredName: "Test Company LLC",
		}

		assert.Equal(t, "Test Company LLC", account.RegisteredName)
	})
}
