package stripe

import (
	"io"
	"net/http"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v83/webhook"
	"reece.start/internal/api"
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

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(body, stripeSignature, config.StripeWebhookSecret)
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
