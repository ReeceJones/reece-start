package jobs

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
	"reece.start/internal/organizations"
	"reece.start/internal/stripe"

	"github.com/resend/resend-go/v2"
)

type RiverClientConfig struct {
	SQLConn      *sql.DB
	DB           *gorm.DB
	Config       *configuration.Config
	ResendClient *resend.Client
	StripeClient *stripeGo.Client
	StartWorkers bool // If true, starts the worker listener
}

// NewRiverClient creates and optionally starts a River client with all workers registered
func NewRiverClient(ctx context.Context, cfg RiverClientConfig) (*river.Client[*sql.Tx], error) {
	// Register all workers
	workers := river.NewWorkers()

	addWorkers(workers, cfg)

	// Create River client
	riverClient, err := river.NewClient(riverdatabasesql.New(cfg.SQLConn), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, err
	}

	// Start workers if requested (production mode)
	if cfg.StartWorkers {
		err = riverClient.Start(ctx)
		if err != nil {
			return nil, err
		}
	}

	return riverClient, nil
}

func addWorkers(workers *river.Workers, cfg RiverClientConfig) {
	river.AddWorker(workers, &organizations.OrganizationInvitationEmailJobWorker{
		DB:           cfg.DB,
		Config:       cfg.Config,
		ResendClient: cfg.ResendClient,
	})
	river.AddWorker(workers, &stripe.SnapshotWebhookProcessingJobWorker{
		DB:           cfg.DB,
		Config:       cfg.Config,
		StripeClient: cfg.StripeClient,
	})
	river.AddWorker(workers, &stripe.ThinWebhookProcessingJobWorker{
		DB:           cfg.DB,
		Config:       cfg.Config,
		StripeClient: cfg.StripeClient,
	})
}
