package stripe

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v83/webhook"
	"reece.start/internal/access"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
)

// StripeSnapshotWebhookEndpoint handles incoming Stripe webhook events
func StripeSnapshotWebhookEndpoint(c echo.Context) error {
	// Get dependencies from context
	config := middleware.GetConfig(c)
	riverClient := middleware.GetRiverClient(c)

	// Check if webhook secret is configured
	if config.StripeWebhookSecret == "" {
		return api.ErrStripeWebhookSecretNotConfigured
	}

	// Read the request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return api.ErrStripeWebhookSignatureMissing
	}

	// Get the Stripe signature from headers
	stripeSignature := c.Request().Header.Get("Stripe-Signature")
	if stripeSignature == "" {
		return api.ErrStripeWebhookSignatureMissing
	}

	// Verify the webhook signature with ignoreAPIVersionMismatch set to true
	event, err := webhook.ConstructEventWithOptions(
		body,
		stripeSignature,
		config.StripeWebhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		return api.ErrStripeWebhookSignatureInvalid
	}

	// Log the webhook event for debugging
	log.Printf("Received Stripe webhook event: %s (ID: %s)", event.Type, event.ID)

	// Enqueue background job for processing
	err = enqueueSnapshotWebhookProcessing(EnqueueSnapshotWebhookProcessingServiceRequest{
		RiverClient: riverClient,
		Event:       &event,
		Context:     c.Request().Context(),
	})
	if err != nil {
		log.Printf("Failed to enqueue webhook processing: %v", err)
		return api.ErrStripeWebhookEventUnhandled
	}

	// Return success response immediately (webhook will be processed in background)
	return c.JSON(http.StatusOK, map[string]any{})
}

func StripeThinWebhookEndpoint(c echo.Context) error {
	// Get dependencies from context
	config := middleware.GetConfig(c)
	riverClient := middleware.GetRiverClient(c)
	stripeClient := middleware.GetStripeClient(c)

	// Check if webhook secret is configured
	if config.StripeWebhookSecret == "" {
		return api.ErrStripeWebhookSecretNotConfigured
	}

	// Read the request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return api.ErrStripeWebhookSignatureMissing
	}

	// Get the Stripe signature from headers
	stripeSignature := c.Request().Header.Get("Stripe-Signature")
	if stripeSignature == "" {
		return api.ErrStripeWebhookSignatureMissing
	}

	// Verify the webhook signature
	eventContainer, err := stripeClient.ParseEventNotification(body, stripeSignature, config.StripeWebhookSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		return api.ErrStripeWebhookSignatureInvalid
	}

	event := eventContainer.GetEventNotification()

	// Log the webhook event for debugging
	log.Printf("Received Stripe webhook event: %s (ID: %s)", event.Type, event.ID)

	// Enqueue background job for processing
	err = enqueueThinWebhookProcessing(EnqueueThinWebhookProcessingServiceRequest{
		RiverClient: riverClient,
		Event:       eventContainer,
		Context:     c.Request().Context(),
	})
	if err != nil {
		log.Printf("Failed to enqueue webhook processing: %v", err)
		return api.ErrStripeWebhookEventUnhandled
	}

	// Return success response immediately (webhook will be processed in background)
	return c.JSON(http.StatusOK, map[string]any{})
}

// CreateCheckoutSessionEndpoint creates a Stripe checkout session for upgrading to pro plan
func CreateCheckoutSessionEndpoint(c echo.Context, req CreateCheckoutSessionRequest) error {
	db := middleware.GetDB(c)
	config := middleware.GetConfig(c)
	stripeClient := middleware.GetStripeClient(c)
	ctx := c.Request().Context()

	// Get organization ID from path
	organizationIDStr := c.Param("id")
	organizationID, err := strconv.ParseUint(organizationIDStr, 10, 32)
	if err != nil {
		return api.ErrInvalidOrganizationID
	}

	// Check organization access
	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: uint(organizationID),
		Scopes:         []constants.UserScope{constants.UserScopeOrganizationBillingUpdate},
	}); err != nil {
		return err
	}

	// Create checkout session
	session, err := CreateCheckoutSession(CreateCheckoutSessionServiceRequest{
		Context:      ctx,
		Config:       config,
		StripeClient: stripeClient,
		DB:           db,
		Params: CreateCheckoutSessionParams{
			OrganizationID: uint(organizationID),
			SuccessURL:     req.SuccessURL,
			CancelURL:      req.CancelURL,
		},
	})

	if err != nil {
		log.Printf("Failed to create checkout session: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create checkout session")
	}

	return c.JSON(http.StatusOK, CheckoutSessionResponse{
		Data: CheckoutSessionData{
			Type: "checkout-session",
			ID:   session.ID,
			Attributes: CheckoutSessionAttributes{
				URL: session.URL,
			},
		},
	})
}

// CreateBillingPortalSessionEndpoint creates a Stripe billing portal session
func CreateBillingPortalSessionEndpoint(c echo.Context, req CreateBillingPortalSessionRequest) error {
	db := middleware.GetDB(c)
	config := middleware.GetConfig(c)
	stripeClient := middleware.GetStripeClient(c)
	ctx := c.Request().Context()

	// Get organization ID from path
	organizationIDStr := c.Param("id")
	organizationID, err := strconv.ParseUint(organizationIDStr, 10, 32)
	if err != nil {
		return api.ErrInvalidOrganizationID
	}

	// Check organization access
	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: uint(organizationID),
		Scopes:         []constants.UserScope{constants.UserScopeOrganizationBillingUpdate},
	}); err != nil {
		return err
	}

	// Create billing portal session
	session, err := CreateBillingPortalSession(CreateBillingPortalSessionServiceRequest{
		Context:      ctx,
		Config:       config,
		StripeClient: stripeClient,
		DB:           db,
		Params: CreateBillingPortalSessionParams{
			OrganizationID: uint(organizationID),
			ReturnURL:      req.ReturnURL,
		},
	})

	if err != nil {
		log.Printf("Failed to create billing portal session: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create billing portal session")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"data": map[string]any{
			"type": "billing-portal-session",
			"id":   session.ID,
			"attributes": map[string]any{
				"url": session.URL,
			},
		},
	})
}

// GetSubscriptionEndpoint retrieves the current subscription for an organization
func GetSubscriptionEndpoint(c echo.Context) error {
	db := middleware.GetDB(c)
	ctx := c.Request().Context()

	// Get organization ID from path
	organizationIDStr := c.Param("id")
	organizationID, err := strconv.ParseUint(organizationIDStr, 10, 32)
	if err != nil {
		return api.ErrInvalidOrganizationID
	}

	// Check organization access (any member can view)
	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: uint(organizationID),
		Scopes:         []constants.UserScope{constants.UserScopeOrganizationRead},
	}); err != nil {
		return err
	}

	// Get subscription
	subscription, err := GetSubscription(GetSubscriptionServiceRequest{
		Context:        ctx,
		DB:             db,
		OrganizationID: uint(organizationID),
	})

	if err != nil {
		log.Printf("Failed to get subscription: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get subscription")
	}

	// If no subscription found, return free plan
	if subscription == nil {
		return c.JSON(http.StatusOK, SubscriptionResponse{
			Data: SubscriptionData{
				Type: "subscription",
				Attributes: SubscriptionAttributes{
					Plan:               string(constants.MembershipPlanFree),
					BillingPeriodStart: nil,
					BillingPeriodEnd:   nil,
					BillingAmount:      0,
				},
			},
		})
	}

	billingPeriodStart := subscription.BillingPeriodStart.Format(time.RFC3339)
	billingPeriodEnd := subscription.BillingPeriodEnd.Format(time.RFC3339)

	return c.JSON(http.StatusOK, SubscriptionResponse{
		Data: SubscriptionData{
			Type: "subscription",
			ID:   strconv.FormatUint(uint64(subscription.ID), 10),
			Attributes: SubscriptionAttributes{
				Plan:               string(subscription.Plan),
				BillingPeriodStart: &billingPeriodStart,
				BillingPeriodEnd:   &billingPeriodEnd,
				BillingAmount:      subscription.BillingPeriodAmount,
			},
		},
	})
}
