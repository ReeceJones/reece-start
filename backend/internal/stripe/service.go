package stripe

import (
	"fmt"
	"log/slog"
	"strconv"

	"encoding/json"
	"errors"
	"time"

	stripeGo "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/billingportal/session"
	checkoutSession "github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/subscription"

	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/internal/utils"
)

func CreateStripeConnectAccount(request CreateStripeAccountServiceRequest) (*stripeGo.V2CoreAccount, error) {
	stripeClient := request.StripeClient
	context := request.Context

	params := &stripeGo.V2CoreAccountCreateParams{
		Dashboard:   stripeGo.String(string(stripeGo.V2CoreAccountDashboardFull)),
		DisplayName: stripeGo.String(request.Params.DisplayName),
		Identity: &stripeGo.V2CoreAccountCreateIdentityParams{
			EntityType:      stripeGo.String(string(request.Params.Type)),
			Country:         stripeGo.String(request.Params.ResidingCountry),
			BusinessDetails: getBusinessDetails(request),
			Individual:      getIndividual(request),
		},
		Defaults: &stripeGo.V2CoreAccountCreateDefaultsParams{
			Currency: stripeGo.String(request.Params.Currency),
			Locales: []*string{
				stripeGo.String(request.Params.Locale),
			},
			Responsibilities: &stripeGo.V2CoreAccountCreateDefaultsResponsibilitiesParams{
				FeesCollector:   stripeGo.String("stripe"),
				LossesCollector: stripeGo.String("stripe"),
			},
		},
		Configuration: &stripeGo.V2CoreAccountCreateConfigurationParams{
			Customer: &stripeGo.V2CoreAccountCreateConfigurationCustomerParams{
				Capabilities: &stripeGo.V2CoreAccountCreateConfigurationCustomerCapabilitiesParams{
					AutomaticIndirectTax: &stripeGo.V2CoreAccountCreateConfigurationCustomerCapabilitiesAutomaticIndirectTaxParams{
						Requested: stripeGo.Bool(true),
					},
				},
			},
			Merchant: &stripeGo.V2CoreAccountCreateConfigurationMerchantParams{
				Capabilities: getMerchantCapabilities(request),
			},
			Recipient: &stripeGo.V2CoreAccountCreateConfigurationRecipientParams{
				Capabilities: &stripeGo.V2CoreAccountCreateConfigurationRecipientCapabilitiesParams{
					StripeBalance: &stripeGo.V2CoreAccountCreateConfigurationRecipientCapabilitiesStripeBalanceParams{
						StripeTransfers: &stripeGo.V2CoreAccountCreateConfigurationRecipientCapabilitiesStripeBalanceStripeTransfersParams{
							Requested: stripeGo.Bool(true),
						},
					},
				},
			},
		},
		Metadata: map[string]string{
			"organization_id": fmt.Sprintf("%d", request.Params.OrganizationID),
		},
		Include: []*string{
			stripeGo.String("configuration.customer"),
			stripeGo.String("configuration.merchant"),
			stripeGo.String("configuration.recipient"),
			stripeGo.String("requirements"),
		},
	}

	if request.Params.ContactEmail != "" {
		params.ContactEmail = stripeGo.String(request.Params.ContactEmail)
	}

	slog.Info("Creating stripe connect account with params", "params", params)

	account, err := stripeClient.V2CoreAccounts.Create(context, params)

	if err != nil {
		slog.Error("Failed to create stripe connect account", "error", err)
		return nil, err
	}

	return account, nil
}

func CreateOnboardingLink(request CreateOnboardingLinkServiceRequest) (*stripeGo.V2CoreAccountLink, error) {
	stripeClient := request.StripeClient
	context := request.Context
	params := request.Params

	link, err := stripeClient.V2CoreAccountLinks.Create(context, &stripeGo.V2CoreAccountLinkCreateParams{
		Account: stripeGo.String(params.AccountID),
		UseCase: &stripeGo.V2CoreAccountLinkCreateUseCaseParams{
			Type: stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseTypeAccountOnboarding)),
			AccountOnboarding: &stripeGo.V2CoreAccountLinkCreateUseCaseAccountOnboardingParams{
				CollectionOptions: &stripeGo.V2CoreAccountLinkCreateUseCaseAccountOnboardingCollectionOptionsParams{
					Fields:             stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingCollectionOptionsFieldsCurrentlyDue)),
					FutureRequirements: stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingCollectionOptionsFutureRequirementsInclude)),
				},
				Configurations: []*string{
					stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingConfigurationCustomer)),
					stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingConfigurationMerchant)),
					stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingConfigurationRecipient)),
				},
				RefreshURL: stripeGo.String(params.RefreshURL),
				ReturnURL:  stripeGo.String(params.ReturnURL),
			},
		},
	})

	if err != nil {
		slog.Error("Failed to create stripe onboarding link", "error", err)
		return nil, err
	}

	return link, err
}

func getBusinessDetails(request CreateStripeAccountServiceRequest) *stripeGo.V2CoreAccountCreateIdentityBusinessDetailsParams {
	if request.Params.Type != stripeGo.AccountBusinessTypeCompany {
		return nil
	}

	return &stripeGo.V2CoreAccountCreateIdentityBusinessDetailsParams{
		// RegisteredName: stripeGo.String(request.Params.Company.RegisteredName),
		Phone: stripeGo.String(request.Params.ContactPhone),
	}
}

func getIndividual(request CreateStripeAccountServiceRequest) *stripeGo.V2CoreAccountCreateIdentityIndividualParams {
	if request.Params.Type != stripeGo.AccountBusinessTypeIndividual {
		return nil
	}

	var line2 *string
	if request.Params.Address.Line2 != "" {
		line2 = stripeGo.String(request.Params.Address.Line2)
	}

	params := &stripeGo.V2CoreAccountCreateIdentityIndividualParams{
		Address: &stripeGo.V2CoreAccountCreateIdentityIndividualAddressParams{
			Line1:      stripeGo.String(request.Params.Address.Line1),
			Line2:      line2,
			City:       stripeGo.String(request.Params.Address.City),
			State:      stripeGo.String(request.Params.Address.StateOrProvince),
			PostalCode: stripeGo.String(request.Params.Address.Zip),
			Country:    stripeGo.String(request.Params.Address.Country),
		},
	}

	if request.Params.ContactEmail != "" {
		params.Email = stripeGo.String(request.Params.ContactEmail)
	}

	if request.Params.ContactPhone != "" {
		params.Phone = stripeGo.String(request.Params.ContactPhone)
	}

	return params
}

func getMerchantCapabilities(request CreateStripeAccountServiceRequest) *stripeGo.V2CoreAccountCreateConfigurationMerchantCapabilitiesParams {
	cap := &stripeGo.V2CoreAccountCreateConfigurationMerchantCapabilitiesParams{
		CardPayments: &stripeGo.V2CoreAccountCreateConfigurationMerchantCapabilitiesCardPaymentsParams{
			Requested: stripeGo.Bool(true),
		},
	}

	// Enable ACH debit payments for US businesses
	if request.Params.ResidingCountry == "US" && request.Config.StripeEnableACHDebitPayments {
		cap.ACHDebitPayments = &stripeGo.V2CoreAccountCreateConfigurationMerchantCapabilitiesACHDebitPaymentsParams{
			Requested: stripeGo.Bool(true),
		}
	}

	return cap
}

func processSnapshotWebhookEvent(request ProcessSnapshotWebhookEventServiceRequest) error {
	// IMPORTANT: do not rely on the data in the event object, as it may not be update to date or delivered out of order.
	// Always refetch the object from the stripe API.
	event := request.Event

	switch event.Type {
	case "customer.subscription.created", "customer.subscription.updated":
		return handleSubscriptionCreatedOrUpdated(request)
	case "customer.subscription.deleted":
		return handleSubscriptionDeleted(request)
	default:
		// Log unhandled events but don't fail
		fmt.Printf("Unhandled webhook event type (snapshot): %s\n", event.Type)
		return nil
	}
}

func handleSubscriptionCreatedOrUpdated(request ProcessSnapshotWebhookEventServiceRequest) error {
	var sub stripeGo.Subscription
	err := json.Unmarshal(request.Event.Data.Raw, &sub)
	if err != nil {
		slog.Error("Failed to unmarshal subscription", "error", err)
		return err
	}

	// Fetch the subscription from Stripe to get the latest data
	fetchedSub, err := subscription.Get(sub.ID, nil)
	if err != nil {
		slog.Error("Failed to fetch subscription from Stripe", "error", err)
		return err
	}

	// Get organization from subscription metadata
	// We store the organization_id in metadata when creating the checkout session
	var org models.Organization

	orgIDStr, ok := fetchedSub.Metadata["organization_id"]
	if !ok {
		slog.Error("No organization_id found in subscription metadata for subscription", "subscriptionID", fetchedSub.ID)
		return nil
	}

	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		slog.Error("Failed to parse organization ID from metadata", "error", err)
		return err
	}

	err = request.DB.WithContext(request.Context).First(&org, uint(orgID)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("Organization not found for subscription", "organizationID", orgID, "subscriptionID", fetchedSub.ID)
			return nil
		}
		return err
	}

	slog.Info("Found organization for subscription", "organizationID", org.ID, "subscriptionID", fetchedSub.ID)

	// Only process active or trialing subscriptions
	if fetchedSub.Status != stripeGo.SubscriptionStatusActive &&
		fetchedSub.Status != stripeGo.SubscriptionStatusTrialing {
		slog.Info("Subscription is not active or trialing, skipping", "subscriptionID", fetchedSub.ID)
		return nil
	}

	// Determine the plan based on the product ID
	var plan constants.MembershipPlan
	for _, item := range fetchedSub.Items.Data {
		slog.Info("Subscription item", "productID", item.Price.Product.ID, "subscriptionID", fetchedSub.ID)
		if item.Price.Product.ID == request.Config.StripeProPlanProductId {
			plan = constants.MembershipPlanPro
			break
		}
	}

	if plan == "" {
		slog.Info("Unknown product in subscription, defaulting to free plan")
		plan = constants.MembershipPlanFree
	}

	// Create or update the plan period
	planPeriod := models.OrganizationPlanPeriod{
		OrganizationID:       org.ID,
		Plan:                 plan,
		StripeSubscriptionID: fetchedSub.ID,
		BillingPeriodStart:   time.Unix(fetchedSub.BillingCycleAnchor, 0),
		BillingPeriodEnd:     time.Unix(fetchedSub.BillingCycleAnchor, 0).AddDate(0, 1, 0), // Assuming monthly subscription
		BillingPeriodAmount:  int(fetchedSub.Items.Data[0].Price.UnitAmount),
	}

	// Check if a plan period already exists for this subscription
	var existingPlanPeriod models.OrganizationPlanPeriod
	err = request.DB.WithContext(request.Context).
		Where("stripe_subscription_id = ?", fetchedSub.ID).
		First(&existingPlanPeriod).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new plan period
		if err := request.DB.WithContext(request.Context).Create(&planPeriod).Error; err != nil {
			return err
		}
		slog.Info("Created new plan period for organization", "organizationID", org.ID)
	} else {
		// Update existing plan period
		existingPlanPeriod.Plan = planPeriod.Plan
		existingPlanPeriod.BillingPeriodStart = planPeriod.BillingPeriodStart
		existingPlanPeriod.BillingPeriodEnd = planPeriod.BillingPeriodEnd
		existingPlanPeriod.BillingPeriodAmount = planPeriod.BillingPeriodAmount
		if err := request.DB.WithContext(request.Context).Save(&existingPlanPeriod).Error; err != nil {
			return err
		}
		slog.Info("Updated plan period for organization", "organizationID", org.ID)
	}

	return nil
}

func handleSubscriptionDeleted(request ProcessSnapshotWebhookEventServiceRequest) error {
	var sub stripeGo.Subscription
	err := json.Unmarshal(request.Event.Data.Raw, &sub)
	if err != nil {
		slog.Error("Failed to unmarshal subscription", "error", err)
		return err
	}

	// Delete the plan period for this subscription
	err = request.DB.WithContext(request.Context).
		Where("stripe_subscription_id = ?", sub.ID).
		Delete(&models.OrganizationPlanPeriod{}).Error

	if err != nil {
		slog.Error("Failed to delete plan period", "error", err)
		return err
	}

	slog.Info("Deleted plan period for subscription", "subscriptionID", sub.ID)
	return nil
}

func processThinWebhookEvent(request ProcessThinWebhookEventServiceRequest) error {
	// IMPORTANT: do not rely on the data in the event object, as it may not be update to date or delivered out of order.
	// Always refetch the object from the stripe API.
	eventContainer := request.Event

	switch evt := eventContainer.(type) {
	case *stripeGo.V2CoreAccountClosedEventNotification:
		return handleAccountClosed(request, evt)
	case *stripeGo.V2CoreAccountUpdatedEventNotification:
		return handleAccountUpdated(request, evt)
	case *stripeGo.V2CoreAccountIncludingConfigurationCustomerCapabilityStatusUpdatedEventNotification:
		return handleAccountCustomerCapabilityStatusUpdated(request, evt)
	case *stripeGo.V2CoreAccountIncludingConfigurationMerchantCapabilityStatusUpdatedEventNotification:
		return handleAccountMerchantCapabilityStatusUpdated(request, evt)
	case *stripeGo.V2CoreAccountIncludingConfigurationRecipientCapabilityStatusUpdatedEventNotification:
		return handleAccountRecipientCapabilityStatusUpdated(request, evt)
	case *stripeGo.V2CoreAccountIncludingRequirementsUpdatedEventNotification:
		return handleAccountRequirementsUpdated(request, evt)
	case *stripeGo.V2CoreAccountIncludingIdentityUpdatedEventNotification:
		return handleAccountIdentityUpdated(request, evt)
	default:
		slog.Info("Unhandled webhook event type (thin)", "type", evt.GetEventNotification().Type)
		return nil
	}
}

// enqueueSnapshotWebhookProcessing persists the event in the background job queue for async processing
func enqueueSnapshotWebhookProcessing(request EnqueueSnapshotWebhookProcessingServiceRequest) error {
	// Serialize the event for storage in the job
	payload, err := json.Marshal(request.Event)
	if err != nil {
		return err
	}

	// River client expects generic args; we defined a dedicated job type in webhook_processing_job.go
	_, err = request.RiverClient.Insert(request.Context, SnapshotWebhookProcessingJob{
		EventID:   request.Event.ID,
		EventType: string(request.Event.Type),
		EventData: payload,
	}, nil)
	return err
}

func enqueueThinWebhookProcessing(request EnqueueThinWebhookProcessingServiceRequest) error {
	// need ot manually serialize these as the correct type so that the related_object field is serialized correctly
	var eventData []byte
	var err error

	switch request.Event.(type) {
	case *stripeGo.UnknownEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.UnknownEventNotification))
	case *stripeGo.V2CoreAccountUpdatedEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.V2CoreAccountUpdatedEventNotification))
	case *stripeGo.V2CoreAccountClosedEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.V2CoreAccountClosedEventNotification))
	case *stripeGo.V2CoreAccountIncludingConfigurationCustomerCapabilityStatusUpdatedEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.V2CoreAccountIncludingConfigurationCustomerCapabilityStatusUpdatedEventNotification))
	case *stripeGo.V2CoreAccountIncludingConfigurationMerchantCapabilityStatusUpdatedEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.V2CoreAccountIncludingConfigurationMerchantCapabilityStatusUpdatedEventNotification))
	case *stripeGo.V2CoreAccountIncludingConfigurationRecipientCapabilityStatusUpdatedEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.V2CoreAccountIncludingConfigurationRecipientCapabilityStatusUpdatedEventNotification))
	case *stripeGo.V2CoreAccountIncludingRequirementsUpdatedEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.V2CoreAccountIncludingRequirementsUpdatedEventNotification))
	case *stripeGo.V2CoreAccountIncludingIdentityUpdatedEventNotification:
		eventData, err = json.Marshal(request.Event.(*stripeGo.V2CoreAccountIncludingIdentityUpdatedEventNotification))
	default:
		return fmt.Errorf("unhandled event type (thin): %s", request.Event.GetEventNotification().Type)
	}

	if err != nil {
		return err
	}

	_, err = request.RiverClient.Insert(request.Context, ThinWebhookProcessingJob{
		EventID:   request.Event.GetEventNotification().ID,
		EventType: string(request.Event.GetEventNotification().Type),
		EventData: eventData,
	}, nil)
	return err
}

func handleAccountUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountUpdatedEventNotification) error {
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID:    event.RelatedObject.ID,
		DB:           request.DB,
		Config:       request.Config,
		StripeClient: request.StripeClient,
		Context:      request.Context,
	})

	if err != nil {
		return err
	}

	slog.Info("Account updated", "accountID", event.RelatedObject.ID)

	return nil
}

func handleAccountClosed(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountClosedEventNotification) error {
	accountID := event.RelatedObject.ID

	slog.Info("Account closed", "accountID", accountID, "reverting stripe information")

	// clear out all stripe information on the account
	err := request.DB.WithContext(request.Context).Transaction(func(tx *gorm.DB) error {
		var org models.Organization
		if err := tx.Where("stripe_account_id = ?", accountID).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Error("Organization not found for stripe_account_id", "accountID", accountID, "during account closure")
				return nil
			}
			return err
		}

		org.Stripe.AccountID = ""
		org.Stripe.AutomaticIndirectTaxStatus = ""
		org.Stripe.CardPaymentsStatus = ""
		org.Stripe.StripeBalancePayoutsStatus = ""
		org.Stripe.StripeBalanceTransfersStatus = ""
		org.Stripe.HasPendingRequirements = false
		org.Stripe.OnboardingStatus = string(utils.DetermineStripeOnboardingStatus(&org))
		// Also reset top-level onboarding to in_progress since connect is gone
		org.OnboardingStatus = string(constants.OnboardingStatusInProgress)

		return tx.Save(&org).Error
	})

	if err != nil {
		return err
	}

	slog.Info("Account closed", "accountID", accountID, "stripe information reverted")

	return nil
}

func handleAccountCustomerCapabilityStatusUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingConfigurationCustomerCapabilityStatusUpdatedEventNotification) error {
	slog.Info("Account customer capability status updated", "accountID", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID:    event.RelatedObject.ID,
		DB:           request.DB,
		Config:       request.Config,
		StripeClient: request.StripeClient,
		Context:      request.Context,
	})

	if err != nil {
		return err
	}

	slog.Info("Account customer capability status updated", "accountID", event.RelatedObject.ID)

	return nil
}

func handleAccountMerchantCapabilityStatusUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingConfigurationMerchantCapabilityStatusUpdatedEventNotification) error {
	slog.Info("Account merchant capability status updated", "accountID", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID:    event.RelatedObject.ID,
		DB:           request.DB,
		Config:       request.Config,
		StripeClient: request.StripeClient,
		Context:      request.Context,
	})

	if err != nil {
		return err
	}

	slog.Info("Account merchant capability status updated", "accountID", event.RelatedObject.ID)

	return nil
}

func handleAccountRecipientCapabilityStatusUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingConfigurationRecipientCapabilityStatusUpdatedEventNotification) error {
	slog.Info("Account recipient capability status updated", "accountID", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID:    event.RelatedObject.ID,
		DB:           request.DB,
		Config:       request.Config,
		StripeClient: request.StripeClient,
		Context:      request.Context,
	})

	if err != nil {
		return err
	}

	slog.Info("Account recipient capability status updated", "accountID", event.RelatedObject.ID)

	return nil
}

func handleAccountRequirementsUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingRequirementsUpdatedEventNotification) error {
	slog.Info("Account requirements updated", "accountID", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID:    event.RelatedObject.ID,
		DB:           request.DB,
		Config:       request.Config,
		StripeClient: request.StripeClient,
		Context:      request.Context,
	})

	if err != nil {
		return err
	}

	slog.Info("Account requirements updated", "accountID", event.RelatedObject.ID)

	return nil
}

func handleAccountIdentityUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingIdentityUpdatedEventNotification) error {
	slog.Info("Account identity updated", "accountID", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID:    event.RelatedObject.ID,
		DB:           request.DB,
		Config:       request.Config,
		StripeClient: request.StripeClient,
		Context:      request.Context,
	})

	if err != nil {
		return err
	}

	slog.Info("Account identity updated", "accountID", event.RelatedObject.ID)

	return nil
}

func CreateCheckoutSession(request CreateCheckoutSessionServiceRequest) (*stripeGo.CheckoutSession, error) {
	db := request.DB
	context := request.Context
	params := request.Params
	config := request.Config
	// Note: Package-level functions use the API key configured when stripeClient was created in server.go

	if config.StripeProPlanPriceId == "" {
		return nil, errors.New("stripe pro plan price ID is not configured")
	}

	if config.StripeProPlanProductId == "" {
		return nil, errors.New("stripe pro plan product ID is not configured")
	}

	// Get the organization
	var org models.Organization
	if err := db.WithContext(context).First(&org, params.OrganizationID).Error; err != nil {
		return nil, err
	}

	// Check if organization has a Stripe Connect account
	if org.Stripe.AccountID == "" {
		return nil, errors.New("organization does not have a Stripe Connect account")
	}

	// In Accounts v2, the account ID is used as the customer ID via customer_account parameter
	// Create checkout session
	sessionParams := &stripeGo.CheckoutSessionParams{
		Mode: stripeGo.String(string(stripeGo.CheckoutSessionModeSubscription)),
		LineItems: []*stripeGo.CheckoutSessionLineItemParams{
			{
				Price:    stripeGo.String(config.StripeProPlanPriceId),
				Quantity: stripeGo.Int64(1),
			},
		},
		SuccessURL: stripeGo.String(params.SuccessURL),
		CancelURL:  stripeGo.String(params.CancelURL),
		Metadata: map[string]string{
			"organization_id": fmt.Sprintf("%d", org.ID),
		},
	}

	// Use customer_account instead of customer for Accounts v2
	sessionParams.AddExtra("customer_account", org.Stripe.AccountID)

	// Create the session (uses the API key configured in stripeClient)
	sess, err := checkoutSession.New(sessionParams)
	if err != nil {
		slog.Error("Failed to create checkout session", "error", err)
		return nil, err
	}

	return sess, nil
}

func CreateBillingPortalSession(request CreateBillingPortalSessionServiceRequest) (*stripeGo.BillingPortalSession, error) {
	db := request.DB
	context := request.Context
	params := request.Params
	config := request.Config
	// Note: Package-level functions use the API key configured when stripeClient was created in server.go

	// Get the organization
	var org models.Organization
	if err := db.WithContext(context).First(&org, params.OrganizationID).Error; err != nil {
		return nil, err
	}

	if org.Stripe.AccountID == "" {
		return nil, errors.New("organization does not have a Stripe Connect account")
	}

	// Create billing portal session using customer_account for Accounts v2
	sessionParams := &stripeGo.BillingPortalSessionParams{
		ReturnURL: stripeGo.String(params.ReturnURL),
	}

	// Add billing portal configuration if provided
	if config.StripeBillingPortalConfigurationId != "" {
		sessionParams.Configuration = stripeGo.String(config.StripeBillingPortalConfigurationId)
	}

	// Use customer_account instead of customer for Accounts v2
	sessionParams.AddExtra("customer_account", org.Stripe.AccountID)

	// Create the session (uses the API key configured in stripeClient)
	sess, err := session.New(sessionParams)
	if err != nil {
		slog.Error("Failed to create billing portal session", "error", err)
		return nil, err
	}

	return sess, nil
}

func GetSubscription(request GetSubscriptionServiceRequest) (*models.OrganizationPlanPeriod, error) {
	db := request.DB
	context := request.Context

	// Get the most recent active subscription for the organization
	var planPeriod models.OrganizationPlanPeriod
	err := db.WithContext(context).
		Where("organization_id = ?", request.OrganizationID).
		Where("billing_period_end > ?", time.Now()).
		Order("billing_period_end DESC").
		First(&planPeriod).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No active subscription found
			return nil, nil
		}
		return nil, err
	}

	return &planPeriod, nil
}

// fetchAndUpdateAccount handles capability status changes.
func fetchAndUpdateAccount(request FetchAndUpdateAccountServiceRequest) error {
	stripeClient := request.StripeClient
	accountID := request.AccountID
	context := request.Context

	slog.Info("Updating account capability status for account", "accountID", accountID)

	params := &stripeGo.V2CoreAccountRetrieveParams{}
	params.AddExtra("include", "configuration.customer")
	params.AddExtra("include", "configuration.merchant")
	params.AddExtra("include", "configuration.recipient")
	params.AddExtra("include", "requirements")

	account, err := stripeClient.V2CoreAccounts.Retrieve(context, accountID, params)
	if err != nil {
		slog.Error("Failed to fetch account", "error", err)
		return err
	}

	accountJson, _ := json.Marshal(account)
	slog.Info("Account fetched", "accountID", accountID, "account", string(accountJson))
	configurationJson, _ := json.Marshal(account.Configuration)
	slog.Info("Account configuration", "accountID", accountID, "configuration", string(configurationJson))
	requirementsJson, _ := json.Marshal(account.Requirements)
	slog.Info("Account requirements", "accountID", accountID, "requirements", string(requirementsJson))
	identityJson, _ := json.Marshal(account.Identity)
	slog.Info("Account identity", "accountID", accountID, "identity", string(identityJson))

	err = request.DB.WithContext(context).Transaction(func(tx *gorm.DB) error {
		var org models.Organization
		if err := tx.Where("stripe_account_id = ?", accountID).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Error("Organization not found for stripe_account_id", "accountID", accountID)
				return nil
			}
			return err
		}
		utils.ApplyStripeAccountToOrganization(&org, account)
		return tx.Save(&org).Error
	})

	if err != nil {
		slog.Error("Failed to update account", "error", err)
		return err
	}

	slog.Info("Account updated", "accountID", accountID)

	return nil
}
