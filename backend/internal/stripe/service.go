package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/riverqueue/river"
	stripeGo "github.com/stripe/stripe-go/v82"
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

	params := &stripeGo.V2CoreAccountCreateIdentityIndividualParams{
		Address: &stripeGo.V2CoreAccountCreateIdentityIndividualAddressParams{
			Line1: stripeGo.String(request.Params.Address.Line1),
			Line2: stripeGo.String(request.Params.Address.Line2),
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
	// if request.Params.ResidingCountry == "US" {
	// 	cap.ACHDebitPayments = &stripeGo.V2CoreAccountCreateConfigurationMerchantCapabilitiesACHDebitPaymentsParams{
	// 		Requested: stripeGo.Bool(true),
	// 	}
	// }

	return cap
}

// processWebhookEvent processes different types of Stripe webhook events
func processWebhookEvent(request ProcessWebhookEventServiceRequest) error {
	event := request.Event

	switch event.Type {
	// Stripe API events
	case "customer.created":
		return handleCustomerCreated(request)
	case "customer.updated":
		return handleCustomerUpdated(request)
	case "customer.deleted":
		return handleCustomerDeleted(request)
	case "invoice.payment_succeeded":
		return handleInvoicePaymentSucceeded(request)
	case "invoice.payment_failed":
		return handleInvoicePaymentFailed(request)
	case "customer.subscription.created":
		return handleSubscriptionCreated(request)
	case "customer.subscription.updated":
		return handleSubscriptionUpdated(request)
	case "customer.subscription.deleted":
		return handleSubscriptionDeleted(request)
	// Stripe Connect events
	case "account.updated":
		return handleAccountUpdated(request)
	case "account.application.deauthorized":
		return handleAccountApplicationDeauthorized(request)
	case "capability.updated":
		return handleCapabilityUpdated(request)
	case "person.created":
		return handlePersonCreated(request)
	case "person.updated":
		return handlePersonUpdated(request)
	case "person.deleted":
		return handlePersonDeleted(request)
	case "payout.created":
		return handlePayoutCreated(request)
	case "payout.updated":
		return handlePayoutUpdated(request)
	case "payout.paid":
		return handlePayoutPaid(request)
	case "payout.failed":
		return handlePayoutFailed(request)
	case "topup.created":
		return handleTopupCreated(request)
	case "topup.succeeded":
		return handleTopupSucceeded(request)
	case "topup.failed":
		return handleTopupFailed(request)
	case "transfer.created":
		return handleTransferCreated(request)
	case "transfer.updated":
		return handleTransferUpdated(request)
	case "application_fee.created":
		return handleApplicationFeeCreated(request)
	case "application_fee.refunded":
		return handleApplicationFeeRefunded(request)
	default:
		// Log unhandled events but don't fail
		fmt.Printf("Unhandled webhook event type: %s\n", event.Type)
		return nil
	}
}

// handleCustomerCreated processes customer.created events
func handleCustomerCreated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing customer.created event: %s\n", event.ID)
	
	// The customer object is already parsed in event.Data.Object
	customer := event.Data.Object
	_ = customer // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, create or update a customer record in your database
	// Example:
	// customerData := customer.(map[string]interface{})
	// stripeCustomerID := customerData["id"].(string)
	// email := customerData["email"].(string)
	// 
	// // Create customer in your database
	// if err := s.createCustomerRecord(stripeCustomerID, email); err != nil {
	//     return fmt.Errorf("failed to create customer record: %w", err)
	// }
	
	return nil
}

// handleCustomerUpdated processes customer.updated events
func handleCustomerUpdated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing customer.updated event: %s\n", event.ID)
	
	// The customer object is already parsed in event.Data.Object
	customer := event.Data.Object
	_ = customer // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update a customer record in your database
	
	return nil
}

// handleCustomerDeleted processes customer.deleted events
func handleCustomerDeleted(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing customer.deleted event: %s\n", event.ID)
	
	// The customer object is already parsed in event.Data.Object
	customer := event.Data.Object
	_ = customer // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, mark a customer record as deleted in your database
	
	return nil
}

// handleInvoicePaymentSucceeded processes invoice.payment_succeeded events
func handleInvoicePaymentSucceeded(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing invoice.payment_succeeded event: %s\n", event.ID)
	
	// The invoice object is already parsed in event.Data.Object
	invoice := event.Data.Object
	_ = invoice // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update subscription status, send confirmation email, etc.
	
	return nil
}

// handleInvoicePaymentFailed processes invoice.payment_failed events
func handleInvoicePaymentFailed(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing invoice.payment_failed event: %s\n", event.ID)
	
	// The invoice object is already parsed in event.Data.Object
	invoice := event.Data.Object
	_ = invoice // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, handle failed payment, send notification, etc.
	
	return nil
}

// handleSubscriptionCreated processes customer.subscription.created events
func handleSubscriptionCreated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing customer.subscription.created event: %s\n", event.ID)
	
	// The subscription object is already parsed in event.Data.Object
	subscription := event.Data.Object
	_ = subscription // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, create a subscription record in your database
	
	return nil
}

// handleSubscriptionUpdated processes customer.subscription.updated events
func handleSubscriptionUpdated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing customer.subscription.updated event: %s\n", event.ID)
	
	// The subscription object is already parsed in event.Data.Object
	subscription := event.Data.Object
	_ = subscription // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update a subscription record in your database
	
	return nil
}

// handleSubscriptionDeleted processes customer.subscription.deleted events
func handleSubscriptionDeleted(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing customer.subscription.deleted event: %s\n", event.ID)
	
	// The subscription object is already parsed in event.Data.Object
	subscription := event.Data.Object
	_ = subscription // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, mark a subscription as cancelled in your database
	
	return nil
}

// enqueueWebhookProcessing enqueues a webhook event for background processing
func enqueueWebhookProcessing(request EnqueueWebhookProcessingServiceRequest) error {
	riverClient := request.RiverClient
	event := request.Event
	// Serialize the event data
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Create the job
	job := WebhookProcessingJob{
		EventID:   event.ID,
		EventType: string(event.Type),
		EventData: eventData,
	}

	// Enqueue the job with a delay to allow for proper processing
	_, err = riverClient.Insert(context.Background(), job, &river.InsertOpts{
		ScheduledAt: time.Now().Add(1 * time.Second), // Small delay to ensure webhook is fully processed
	})
	if err != nil {
		return fmt.Errorf("failed to enqueue webhook processing job: %w", err)
	}

	fmt.Printf("Enqueued webhook processing job for event %s (type: %s)\n", event.ID, event.Type)
	return nil
}

// Stripe Connect webhook handlers

// handleAccountUpdated processes account.updated events
func handleAccountUpdated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing account.updated event: %s\n", event.ID)
	
	// The account object is already parsed in event.Data.Object
	account := event.Data.Object
	_ = account // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update connected account information in your database
	
	return nil
}

// handleAccountApplicationDeauthorized processes account.application.deauthorized events
func handleAccountApplicationDeauthorized(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing account.application.deauthorized event: %s\n", event.ID)
	
	// The account object is already parsed in event.Data.Object
	account := event.Data.Object
	_ = account // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, remove connected account access or notify the merchant
	
	return nil
}

// handleCapabilityUpdated processes capability.updated events
func handleCapabilityUpdated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing capability.updated event: %s\n", event.ID)
	
	// The capability object is already parsed in event.Data.Object
	capability := event.Data.Object
	_ = capability // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update account capabilities status in your database
	
	return nil
}

// handlePersonCreated processes person.created events
func handlePersonCreated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing person.created event: %s\n", event.ID)
	
	// The person object is already parsed in event.Data.Object
	person := event.Data.Object
	_ = person // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, store person information for the connected account
	
	return nil
}

// handlePersonUpdated processes person.updated events
func handlePersonUpdated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing person.updated event: %s\n", event.ID)
	
	// The person object is already parsed in event.Data.Object
	person := event.Data.Object
	_ = person // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update person information in your database
	
	return nil
}

// handlePersonDeleted processes person.deleted events
func handlePersonDeleted(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing person.deleted event: %s\n", event.ID)
	
	// The person object is already parsed in event.Data.Object
	person := event.Data.Object
	_ = person // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, remove person information from your database
	
	return nil
}

// handlePayoutCreated processes payout.created events
func handlePayoutCreated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing payout.created event: %s\n", event.ID)
	
	// The payout object is already parsed in event.Data.Object
	payout := event.Data.Object
	_ = payout // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, record the payout in your database
	
	return nil
}

// handlePayoutUpdated processes payout.updated events
func handlePayoutUpdated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing payout.updated event: %s\n", event.ID)
	
	// The payout object is already parsed in event.Data.Object
	payout := event.Data.Object
	_ = payout // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update payout status in your database
	
	return nil
}

// handlePayoutPaid processes payout.paid events
func handlePayoutPaid(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing payout.paid event: %s\n", event.ID)
	
	// The payout object is already parsed in event.Data.Object
	payout := event.Data.Object
	_ = payout // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, mark payout as completed and notify the connected account
	
	return nil
}

// handlePayoutFailed processes payout.failed events
func handlePayoutFailed(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing payout.failed event: %s\n", event.ID)
	
	// The payout object is already parsed in event.Data.Object
	payout := event.Data.Object
	_ = payout // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, handle failed payout and notify the connected account
	
	return nil
}

// handleTopupCreated processes topup.created events
func handleTopupCreated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing topup.created event: %s\n", event.ID)
	
	// The topup object is already parsed in event.Data.Object
	topup := event.Data.Object
	_ = topup // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, record the topup in your database
	
	return nil
}

// handleTopupSucceeded processes topup.succeeded events
func handleTopupSucceeded(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing topup.succeeded event: %s\n", event.ID)
	
	// The topup object is already parsed in event.Data.Object
	topup := event.Data.Object
	_ = topup // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update account balance information
	
	return nil
}

// handleTopupFailed processes topup.failed events
func handleTopupFailed(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing topup.failed event: %s\n", event.ID)
	
	// The topup object is already parsed in event.Data.Object
	topup := event.Data.Object
	_ = topup // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, handle failed topup and notify appropriate parties
	
	return nil
}

// handleTransferCreated processes transfer.created events
func handleTransferCreated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing transfer.created event: %s\n", event.ID)
	
	// The transfer object is already parsed in event.Data.Object
	transfer := event.Data.Object
	_ = transfer // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, record the transfer to a connected account
	
	return nil
}

// handleTransferUpdated processes transfer.updated events
func handleTransferUpdated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing transfer.updated event: %s\n", event.ID)
	
	// The transfer object is already parsed in event.Data.Object
	transfer := event.Data.Object
	_ = transfer // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, update transfer status in your database
	
	return nil
}

// handleApplicationFeeCreated processes application_fee.created events
func handleApplicationFeeCreated(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing application_fee.created event: %s\n", event.ID)
	
	// The application fee object is already parsed in event.Data.Object
	applicationFee := event.Data.Object
	_ = applicationFee // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, record application fee collection
	
	return nil
}

// handleApplicationFeeRefunded processes application_fee.refunded events
func handleApplicationFeeRefunded(request ProcessWebhookEventServiceRequest) error {
	event := request.Event
	db := request.DB
	config := request.Config
	_ = db     // Suppress unused variable warning
	_ = config // Suppress unused variable warning
	
	fmt.Printf("Processing application_fee.refunded event: %s\n", event.ID)
	
	// The application fee object is already parsed in event.Data.Object
	applicationFee := event.Data.Object
	_ = applicationFee // Suppress unused variable warning
	
	// TODO: Implement your business logic here
	// For example, handle application fee refund
	
	return nil
}
