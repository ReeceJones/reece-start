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

func processWebhookEvent(request ProcessWebhookEventServiceRequest) error {
	// IMPORTANT: do not rely on the data in the event object, as it may not be update to date or delivered out of order.
	// Always refetch the object from the stripe API.
	// V2 events are actually "thin" objects and already operate under this assumed behavior.
	event := request.Event

	switch event.Type {
	case "v2.core.account.updated":
		return handleAccountUpdated(request)
	case "v2.core.account.closed":
		return handleAccountClosed(request)
	case "capability.updated":
		return handleAccountCapabilityStatusUpdated(request)
	case "v2.core.account[configuration.customer].capability_status_updated":
		return handleAccountCapabilityStatusUpdated(request)
	case "v2.core.account[configuration.merchant].capability_status_updated":
		return handleAccountCapabilityStatusUpdated(request)
	case "v2.core.account[configuration.recipient].capability_status_updated":
		return handleAccountCapabilityStatusUpdated(request)
	case "v2.core.account[requirements].updated":
		return handleAccountRequirementsUpdated(request)
	default:
		// Log unhandled events but don't fail
		fmt.Printf("Unhandled webhook event type: %s\n", event.Type)
		return nil
	}
}

func handleAccountUpdated(request ProcessWebhookEventServiceRequest) error {
	return nil
}

// enqueueWebhookProcessing persists the event in the background job queue for async processing
func enqueueWebhookProcessing(request EnqueueWebhookProcessingServiceRequest) error {
    // Serialize the event for storage in the job
    payload, err := json.Marshal(request.Event)
    if err != nil {
        return err
    }

    // River client expects generic args; we defined a dedicated job type in webhook_processing_job.go
    _, err = request.RiverClient.Insert(request.Context, WebhookProcessingJob{
        EventID:   request.Event.ID,
        EventType: string(request.Event.Type),
        EventData: payload,
    }, nil)
    return err
}

// handleAccountClosed handles the case where a Stripe account is closed.
// It removes the Stripe account ID and resets onboarding and Stripe statuses.
func handleAccountClosed(request ProcessWebhookEventServiceRequest) error {
    accountID, err := extractAccountIDFromEvent(request)
    if err != nil {
        return nil // ignore if we cannot determine, to avoid retries
    }

    return request.DB.Transaction(func(tx *gorm.DB) error {
        var org models.Organization
        if err := tx.Where("stripe_account_id = ?", accountID).First(&org).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
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
}

// handleAccountCapabilityStatusUpdated handles capability status changes.
func handleAccountCapabilityStatusUpdated(request ProcessWebhookEventServiceRequest) error {
    accountID, err := extractAccountIDFromEvent(request)
    if err != nil {
        return nil
    }

    acc, err := request.StripeClient.V2CoreAccounts.Retrieve(request.Context, accountID, &stripeGo.V2CoreAccountRetrieveParams{
		Include: []*string{
			stripeGo.String("configuration.customer"),
			stripeGo.String("configuration.merchant"),
			stripeGo.String("configuration.recipient"),
			stripeGo.String("requirements"),
		},
	})
    if err != nil {
        return err
    }

    return request.DB.Transaction(func(tx *gorm.DB) error {
        var org models.Organization
        if err := tx.Where("stripe_account_id = ?", accountID).First(&org).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return nil
            }
            return err
        }

        // Apply all capability and requirements fields from the v2 core account
        utils.ApplyStripeAccountToOrganization(&org, acc)
        return tx.Save(&org).Error
    })
}

// handleAccountRequirementsUpdated handles requirements updates similarly to capability updates.
func handleAccountRequirementsUpdated(request ProcessWebhookEventServiceRequest) error {
    accountID, err := extractAccountIDFromEvent(request)
    if err != nil {
        return nil
    }

    acc, err := request.StripeClient.V2CoreAccounts.Retrieve(request.Context, accountID, &stripeGo.V2CoreAccountRetrieveParams{
        Include: []*string{
            stripeGo.String("configuration.customer"),
            stripeGo.String("configuration.merchant"),
            stripeGo.String("configuration.recipient"),
            stripeGo.String("requirements"),
        },
    })
    if err != nil {
        return err
    }

    return request.DB.Transaction(func(tx *gorm.DB) error {
        var org models.Organization
        if err := tx.Where("stripe_account_id = ?", accountID).First(&org).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return nil
            }
            return err
        }

        utils.ApplyStripeAccountToOrganization(&org, acc)
        return tx.Save(&org).Error
    })
}

// extractAccountIDFromEvent attempts to extract the account id from the event's account field or metadata
func extractAccountIDFromEvent(request ProcessWebhookEventServiceRequest) (string, error) {
    if request.Event.Account != "" {
        return request.Event.Account, nil
    }
    // fallback: some events may include in data.object.id for account.updated
    if request.Event.Data.Object != nil {
        if idRaw, ok := request.Event.Data.Object["id"]; ok {
            if idStr, ok2 := idRaw.(string); ok2 {
                return idStr, nil
            }
        }
    }
    return "", fmt.Errorf("account id not found in event")
}
