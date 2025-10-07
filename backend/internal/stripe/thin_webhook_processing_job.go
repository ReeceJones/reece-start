package stripe

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
)

// SnapshotWebhookProcessingJob represents a background job for processing Stripe webhooks
type ThinWebhookProcessingJob struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	EventData []byte `json:"event_data"`
}

// Kind returns the job kind for River
func (ThinWebhookProcessingJob) Kind() string {
	return "thin_stripe_webhook_processing"
}

// SnapshotWebhookProcessingJobWorker handles the background processing of Stripe webhook events
type ThinWebhookProcessingJobWorker struct {
	river.WorkerDefaults[ThinWebhookProcessingJob]
	DB     *gorm.DB
	Config *configuration.Config
    StripeClient *stripeGo.Client
}

// Work processes the webhook event in the background
func (w *ThinWebhookProcessingJobWorker) Work(ctx context.Context, job *river.Job[ThinWebhookProcessingJob]) error {
	eventContainer, err := stripeGo.EventNotificationFromJSON(job.Args.EventData, *w.StripeClient)
	if err != nil {
		return fmt.Errorf("failed to parse event data: %w", err)
	}

	// Process the webhook event
	return processThinWebhookEvent(ProcessThinWebhookEventServiceRequest{
		Event:  eventContainer,
		DB:     w.DB,
		Config: w.Config,
        StripeClient: w.StripeClient,
        Context: ctx,
	})
}
