package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/resend/resend-go/v2"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/configuration"
	"reece.start/internal/database"
	appMiddleware "reece.start/internal/middleware"
	"reece.start/internal/organizations"
	"reece.start/internal/stripe"
	"reece.start/internal/users"
)




func main() {
	// Load environment variables
	config, err := configuration.LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("Error loading environment variables, %s", err)
	}

	// Create database connection
	conn, err := sql.Open("pgx", config.DatabaseUri)
	if err != nil {
		log.Fatalf("Error opening database, %s", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
	}

	log.Printf("Database connected\n")

	// Run database migrations
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Error migrating database, %s", err)
	}

	log.Printf("Database migrated\n")

	// Create minio client (Storage)
	minioClient, err := minio.New(config.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.StorageAccessKeyId, config.StorageSecretAccessKey, ""),
		Secure: config.StorageUseSSL,
	})
	if err != nil {
		log.Fatalf("Error creating minio client, %s", err)
	}

	log.Printf("Minio client created\n")

	// Create resend client (Email)
    resendClient := resend.NewClient(config.ResendApiKey)

	log.Printf("Resend client created\n")

	// Create stripe client and configure the global API key
	stripeGo.Key = config.StripeSecretKey
	stripeClient := stripeGo.NewClient(config.StripeSecretKey)

	// Create river client (Background jobs)
	workers := river.NewWorkers()
	river.AddWorker(workers, &organizations.OrganizationInvitationEmailJobWorker{
		DB:          db,
		Config:      config,
		ResendClient: resendClient,
	})
	river.AddWorker(workers, &stripe.SnapshotWebhookProcessingJobWorker{
		DB:     db,
		Config: config,
		StripeClient: stripeClient,
	})
	river.AddWorker(workers, &stripe.ThinWebhookProcessingJobWorker{
		DB:     db,
		Config: config,
		StripeClient: stripeClient,
	})

	riverClient, err := river.NewClient(riverdatabasesql.New(conn), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})
	if err != nil {
		log.Fatalf("Error creating river client, %s", err)
	}

	log.Printf("River client created\n")

	err = riverClient.Start(context.Background())
	if err != nil {
		log.Fatalf("Error starting river client, %s", err)
	}

	log.Printf("River client started\n")

	e := echo.New()

	// Accept application/json requests (JSON)
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody []byte, resBody []byte) {
		log.Println(string(reqBody))
	}))

	// Add logging middleware
	e.Use(middleware.Logger())

	// Add CORS middleware
	e.Use(middleware.CORS())

	// Add dependency injection middleware
	e.Use(appMiddleware.ContentTypeMiddleware)
	e.Use(appMiddleware.DependencyInjectionMiddleware(appMiddleware.AppDependencies{
		Config: config,
		DB: db,
		MinioClient: minioClient,
		RiverClient: riverClient,
		ResendClient: resendClient,
		StripeClient: stripeClient,
	}))
	
	// Add error handling middleware
	e.Use(appMiddleware.ErrorHandlingMiddleware)

	// health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Public user routes (no authentication required)
	e.POST("/users", api.Validated(users.CreateUserEndpoint))
	e.POST("/users/login", api.Validated(users.LoginEndpoint))

	// Public OAuth routes (no authentication required)
	e.POST("/oauth/google/callback", api.Validated(users.GoogleOAuthCallbackEndpoint))

	// Webhook routes
	e.POST("/webhooks/stripe/snapshot", stripe.StripeSnapshotWebhookEndpoint)
	e.POST("/webhooks/stripe/thin", stripe.StripeThinWebhookEndpoint)

	// Protected user routes (authentication required)
	protected := e.Group("")
	protected.Use(appMiddleware.JwtAuthMiddleware(config))

	protected.GET("/users/me", users.GetAuthenticatedUserEndpoint)
	protected.GET("/users", api.ValidatedQuery(users.GetUsersEndpoint))
	protected.POST("/users/me/token", api.Validated(users.CreateAuthenticatedUserTokenEndpoint))
	protected.PATCH("/users/:id", api.Validated(users.UpdateUserEndpoint))

	protected.GET("/organizations", organizations.GetOrganizationsEndpoint)
	protected.POST("/organizations", api.Validated(organizations.CreateOrganizationEndpoint))
	protected.GET("/organizations/:id", organizations.GetOrganizationEndpoint)
	protected.PATCH("/organizations/:id", api.Validated(organizations.UpdateOrganizationEndpoint))
	protected.DELETE("/organizations/:id", organizations.DeleteOrganizationEndpoint)
	protected.POST("/organizations/:id/stripe-onboarding-link", organizations.CreateStripeOnboardingLinkEndpoint)
	protected.POST("/organizations/:id/stripe-dashboard-link", organizations.CreateStripeDashboardLinkEndpoint)
	protected.GET("/organizations/:id/subscription", stripe.GetSubscriptionEndpoint)
	protected.POST("/organizations/:id/checkout-session", api.Validated(stripe.CreateCheckoutSessionEndpoint))
	protected.POST("/organizations/:id/billing-portal-session", api.Validated(stripe.CreateBillingPortalSessionEndpoint))

	protected.GET("/organization-memberships", api.ValidatedQuery(organizations.GetOrganizationMembershipsEndpoint))
	protected.GET("/organization-memberships/:id", organizations.GetOrganizationMembershipEndpoint)
	protected.POST("/organization-memberships", api.Validated(organizations.CreateOrganizationMembershipEndpoint))
	protected.PATCH("/organization-memberships/:id", api.Validated(organizations.UpdateOrganizationMembershipEndpoint))
	protected.DELETE("/organization-memberships/:id", organizations.DeleteOrganizationMembershipEndpoint)

	protected.POST("/organization-invitations", api.Validated(organizations.InviteToOrganizationEndpoint))
	protected.POST("/organization-invitations/:id/accept", api.Validated(organizations.AcceptOrganizationInvitationEndpoint))
	protected.POST("/organization-invitations/:id/decline", api.Validated(organizations.DeclineOrganizationInvitationEndpoint))
	protected.GET("/organization-invitations", api.ValidatedQuery(organizations.GetOrganizationInvitationsEndpoint))
	protected.GET("/organization-invitations/:id", organizations.GetOrganizationInvitationEndpoint)
	protected.DELETE("/organization-invitations/:id", organizations.DeleteOrganizationInvitationEndpoint)

	// Start http server
	e.Logger.Fatal(e.Start(":8080"))
}
