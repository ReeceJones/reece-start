package stripe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/riverqueue/river"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
)

// WebhookProcessingJob represents a background job for processing Stripe webhooks
type WebhookProcessingJob struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	EventData []byte `json:"event_data"`
}

// Kind returns the job kind for River
func (WebhookProcessingJob) Kind() string {
	return "stripe_webhook_processing"
}

// WebhookProcessingJobWorker handles the background processing of Stripe webhook events
type WebhookProcessingJobWorker struct {
	river.WorkerDefaults[WebhookProcessingJob]
	DB     *gorm.DB
	Config *configuration.Config
    StripeClient *Client
}

// Work processes the webhook event in the background
func (w *WebhookProcessingJobWorker) Work(ctx context.Context, job *river.Job[WebhookProcessingJob]) error {
	// Parse the event data back into a Stripe event
    var event stripeGo.Event
	if err := json.Unmarshal(job.Args.EventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	// Process the webhook event
	return processWebhookEvent(ProcessWebhookEventServiceRequest{
		Event:  &event,
		DB:     w.DB,
		Config: w.Config,
        StripeClient: w.StripeClient,
        Context: ctx,
	})
}
