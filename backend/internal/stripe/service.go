package stripe

import (
	"fmt"
	"log"

	"encoding/json"
	"errors"

	stripeGo "github.com/stripe/stripe-go/v83"

	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/internal/utils"
)


func CreateStripeConnectAccount(request CreateStripeAccountServiceRequest) (*stripeGo.V2CoreAccount, error) {
	stripeClient := request.StripeClient
	context := request.Context

	params := &stripeGo.V2CoreAccountCreateParams{
		Dashboard: stripeGo.String(string(stripeGo.V2CoreAccountDashboardFull)),
		DisplayName: stripeGo.String(request.Params.DisplayName),
		Identity: &stripeGo.V2CoreAccountCreateIdentityParams{
			EntityType: stripeGo.String(string(request.Params.Type)),
			Country: stripeGo.String(request.Params.ResidingCountry),
			BusinessDetails: getBusinessDetails(request),
			Individual: getIndividual(request),
		},
		Defaults: &stripeGo.V2CoreAccountCreateDefaultsParams{
			Currency: stripeGo.String(request.Params.Currency),
			Locales: []*string{
				stripeGo.String(request.Params.Locale),
			},
			Responsibilities: &stripeGo.V2CoreAccountCreateDefaultsResponsibilitiesParams{
				FeesCollector: stripeGo.String("stripe"),
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


	log.Printf("Creating stripe connect account with params: %+v", params)

	account, err := stripeClient.V2CoreAccounts.Create(context, params)

	if err != nil {
		log.Printf("failed to create stripe connect account: %v", err)
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
					Fields: stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingCollectionOptionsFieldsCurrentlyDue)),
					FutureRequirements: stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingCollectionOptionsFutureRequirementsInclude)),
				},
				Configurations: []*string{
					stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingConfigurationCustomer)),
					stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingConfigurationMerchant)),
					stripeGo.String(string(stripeGo.V2CoreAccountLinkUseCaseAccountOnboardingConfigurationRecipient)),
				},
				RefreshURL: stripeGo.String(params.RefreshURL),
				ReturnURL: stripeGo.String(params.ReturnURL),
			},
		},
	})

	if err != nil {
		log.Printf("failed to create stripe onboarding link: %v", err)
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
			Line1: stripeGo.String(request.Params.Address.Line1),
			Line2: line2,
			City: stripeGo.String(request.Params.Address.City),
			State: stripeGo.String(request.Params.Address.StateOrProvince),
			PostalCode: stripeGo.String(request.Params.Address.Zip),
			Country: stripeGo.String(request.Params.Address.Country),
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
	default:
		// Log unhandled events but don't fail
		fmt.Printf("Unhandled webhook event type (snapshot): %s\n", event.Type)
		return nil
	}
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
			log.Printf("Unhandled webhook event type (thin): %s\n", evt.GetEventNotification().Type)
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
		AccountID: event.RelatedObject.ID,
		DB: request.DB,
		Config: request.Config,
		StripeClient: request.StripeClient,
		Context: request.Context,
	})

	if err != nil {
		return err
	}

	log.Printf("Account updated: %s", event.RelatedObject.ID)

	return nil
}

func handleAccountClosed(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountClosedEventNotification) error {	
	accountID := event.RelatedObject.ID

	log.Printf("Account closed: %s. Reverting stripe information.", accountID)

	// clear out all stripe information on the account
	err := request.DB.WithContext(request.Context).Transaction(func(tx *gorm.DB) error {
		var org models.Organization
		if err := tx.Where("stripe_account_id = ?", accountID).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Organization not found for stripe_account_id: %s during account closure", accountID)
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

	log.Printf("Account closed: %s. Stripe information reverted.", accountID)

	return nil
}

func handleAccountCustomerCapabilityStatusUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingConfigurationCustomerCapabilityStatusUpdatedEventNotification) error {
	log.Printf("Account customer capability status updated: %s", event.RelatedObject.ID)
	
	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID: event.RelatedObject.ID,
		DB: request.DB,
		Config: request.Config,
		StripeClient: request.StripeClient,
		Context: request.Context,
	})

	if err != nil {
		return err
	}

	log.Printf("Account customer capability status updated: %s", event.RelatedObject.ID)

	return nil
}

func handleAccountMerchantCapabilityStatusUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingConfigurationMerchantCapabilityStatusUpdatedEventNotification) error {
	log.Printf("Account merchant capability status updated: %s", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID: event.RelatedObject.ID,
		DB: request.DB,
		Config: request.Config,
		StripeClient: request.StripeClient,
		Context: request.Context,
	})

	if err != nil {
		return err
	}

	log.Printf("Account merchant capability status updated: %s", event.RelatedObject.ID)

	return nil
}

func handleAccountRecipientCapabilityStatusUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingConfigurationRecipientCapabilityStatusUpdatedEventNotification) error {
	log.Printf("Account recipient capability status updated: %s", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID: event.RelatedObject.ID,
		DB: request.DB,
		Config: request.Config,
		StripeClient: request.StripeClient,
		Context: request.Context,
	})

	if err != nil {
		return err
	}

	log.Printf("Account recipient capability status updated: %s", event.RelatedObject.ID)

	return nil
}

func handleAccountRequirementsUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingRequirementsUpdatedEventNotification) error {
	log.Printf("Account requirements updated: %s", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID: event.RelatedObject.ID,
		DB: request.DB,
		Config: request.Config,
		StripeClient: request.StripeClient,
		Context: request.Context,
	})

	if err != nil {
		return err
	}

	log.Printf("Account requirements updated: %s", event.RelatedObject.ID)

	return nil
}

func handleAccountIdentityUpdated(request ProcessThinWebhookEventServiceRequest, event *stripeGo.V2CoreAccountIncludingIdentityUpdatedEventNotification) error {
	log.Printf("Account identity updated: %s", event.RelatedObject.ID)

	// RelatedObject is v2.core.account
	err := fetchAndUpdateAccount(FetchAndUpdateAccountServiceRequest{
		AccountID: event.RelatedObject.ID,
		DB: request.DB,
		Config: request.Config,
		StripeClient: request.StripeClient,
		Context: request.Context,
	})

	if err != nil {
		return err
	}

	log.Printf("Account identity updated: %s", event.RelatedObject.ID)

	return nil
}

// fetchAndUpdateAccount handles capability status changes.
func fetchAndUpdateAccount(request FetchAndUpdateAccountServiceRequest) error {
	stripeClient := request.StripeClient
	accountID := request.AccountID
	context := request.Context

	log.Printf("Updating account capability status for account %s", accountID)

	params := &stripeGo.V2CoreAccountRetrieveParams{}
	params.AddExtra("include", "configuration.customer")
	params.AddExtra("include", "configuration.merchant")
	params.AddExtra("include", "configuration.recipient")
	params.AddExtra("include", "requirements")

	account, err := stripeClient.V2CoreAccounts.Retrieve(context, accountID, params)
	if err != nil {
		log.Printf("failed to fetch account: %v", err)
		return err
	}

	accountJson, _ := json.Marshal(account)
	log.Printf("Account fetched: %s", string(accountJson))
	configurationJson, _ := json.Marshal(account.Configuration)
	log.Printf("Account configuration: %s", string(configurationJson))
	requirementsJson, _ := json.Marshal(account.Requirements)
	log.Printf("Account requirements: %s", string(requirementsJson))
	identityJson, _ := json.Marshal(account.Identity)
	log.Printf("Account identity: %s", string(identityJson))

	err = request.DB.WithContext(context).Transaction(func(tx *gorm.DB) error {
		var org models.Organization
		if err := tx.Where("stripe_account_id = ?", accountID).First(&org).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Organization not found for stripe_account_id: %s", accountID)
				return nil
			}
			return err
		}
		utils.ApplyStripeAccountToOrganization(&org, account)
		return tx.Save(&org).Error
	})

	if err != nil {
		log.Printf("failed to update account: %v", err)
		return err
	}

	log.Printf("Account updated: %s", accountID)

	return nil
}
