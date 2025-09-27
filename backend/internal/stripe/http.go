package stripe

import (
	"io"
	"net/http"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
	"github.com/stripe/stripe-go/v82/webhook"
	"reece.start/internal/api"
	"reece.start/internal/configuration"
)

// StripeWebhookEndpoint handles incoming Stripe webhook events
func StripeWebhookEndpoint(c echo.Context) error {
	// Get dependencies from context
	config := c.Get("config").(*configuration.Config)
	riverClient := c.Get("riverClient").(*river.Client[river.JobArgs])
	
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
	err = enqueueWebhookProcessing(EnqueueWebhookProcessingServiceRequest{
		RiverClient: riverClient,
		Event:       &event,
	})
	if err != nil {
		log.Printf("Failed to enqueue webhook processing: %v", err)
		return api.ErrStripeWebhookEventUnhandled
	}

	// Return success response immediately (webhook will be processed in background)
	return c.JSON(http.StatusOK, map[string]any{})
}
