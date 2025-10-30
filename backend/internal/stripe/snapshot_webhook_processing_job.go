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

// SnapshotWebhookProcessingJob represents a background job for processing Stripe webhooks
type SnapshotWebhookProcessingJob struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	EventData []byte `json:"event_data"`
}

// Kind returns the job kind for River
func (SnapshotWebhookProcessingJob) Kind() string {
	return "snapshot_stripe_webhook_processing"
}

// SnapshotWebhookProcessingJobWorker handles the background processing of Stripe webhook events
type SnapshotWebhookProcessingJobWorker struct {
	river.WorkerDefaults[SnapshotWebhookProcessingJob]
	DB           *gorm.DB
	Config       *configuration.Config
	StripeClient *stripeGo.Client
}

// Work processes the webhook event in the background
func (w *SnapshotWebhookProcessingJobWorker) Work(ctx context.Context, job *river.Job[SnapshotWebhookProcessingJob]) error {
	// Parse the event data back into a Stripe event
	var event stripeGo.Event
	if err := json.Unmarshal(job.Args.EventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	// Process the webhook event
	return processSnapshotWebhookEvent(ProcessSnapshotWebhookEventServiceRequest{
		Event:        &event,
		DB:           w.DB,
		Config:       w.Config,
		StripeClient: w.StripeClient,
		Context:      ctx,
	})
}
