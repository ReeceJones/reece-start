package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/resend/resend-go/v2"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
	"reece.start/internal/constants"
	"reece.start/internal/database"
	echoServer "reece.start/internal/echo"
	"reece.start/internal/jobs"
	appMiddleware "reece.start/internal/middleware"
	"reece.start/internal/posthog"
)

func main() {
	setupLogger()

	config, err := configuration.LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("Error loading environment variables, %s", err)
	}

	conn, db := createDatabaseConnectionPool(config)
	runDatabaseMigrations(db)

	minioClient := createMinioClient(config)
	initializeStorageBuckets(minioClient)

	resendClient := createResendClient(config)
	stripeClient := createStripeClient(config)
	posthogClient := createPostHogClient(config)

	ctx := context.Background()
	runRiverMigrations(ctx, conn)

	riverClient := createRiverClient(ctx, config, conn, db, resendClient, stripeClient)

	e := createEchoServer(config, db, minioClient, riverClient, resendClient, stripeClient, posthogClient)

	// Optional: Add body dump middleware for debugging (production only)
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody []byte, resBody []byte) {
		slog.Info("Body dump", "request", string(reqBody), "response", string(resBody))
	}))

	// Add sentry middleware
	if config.SentryDsn != "" {
		e.Use(sentryecho.New(sentryecho.Options{}))
	}

	// Start HTTP server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Set up signal handling for graceful shutdown
	sigintOrTerm := make(chan os.Signal, 1)
	signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either server error or shutdown signal
	select {
	case err := <-serverErr:
		slog.Error("Server error", "error", err)
		gracefulShutdown(ctx, e, riverClient, posthogClient, 10*time.Second)
	case <-sigintOrTerm:
		slog.Info("Received SIGINT/SIGTERM; initiating graceful shutdown")
		gracefulShutdown(ctx, e, riverClient, posthogClient, 10*time.Second)
	}
}

func setupLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func setupSentry(config *configuration.Config) {
	if config.SentryDsn == "" {
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.SentryDsn,
	})

	if err != nil {
		log.Fatalf("Error initializing sentry, %s", err)
	}
}

func createDatabaseConnectionPool(config *configuration.Config) (*pgxpool.Pool, *gorm.DB) {
	pool, err := pgxpool.New(context.Background(), config.DatabaseUri)
	if err != nil {
		log.Fatalf("Error opening database, %s", err)
	}

	// More info here: https://github.com/pitabwire/frame/blob/v1.63.1/datastore/pool/connection.go#L46-L91
	connector := stdlib.GetPoolConnector(pool)
	sqlConn := sql.OpenDB(connector)

	// don't hold idle connections
	sqlConn.SetMaxIdleConns(0)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlConn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
	}

	slog.Info("Database connected")
	return pool, db
}

func runDatabaseMigrations(db *gorm.DB) {
	err := database.Migrate(db)
	if err != nil {
		log.Fatalf("Error migrating database, %s", err)
	}

	slog.Info("Database migrated")
}

func createMinioClient(config *configuration.Config) *minio.Client {
	minioClient, err := minio.New(config.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.StorageAccessKeyId, config.StorageSecretAccessKey, ""),
		Secure: config.StorageUseSSL,
	})
	if err != nil {
		log.Fatalf("Error creating minio client, %s", err)
	}

	slog.Info("Minio client created")
	return minioClient
}

func initializeStorageBuckets(client *minio.Client) {
	ctx := context.Background()
	buckets := []constants.StorageBucket{
		constants.StorageBucketUserLogos,
		constants.StorageBucketOrganizationLogos,
	}

	for _, bucket := range buckets {
		bucketName := string(bucket)
		exists, err := client.BucketExists(ctx, bucketName)
		if err != nil {
			log.Fatalf("Error checking if bucket %s exists, %s", bucketName, err)
		}

		if !exists {
			err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if err != nil {
				log.Fatalf("Error creating bucket %s, %s", bucketName, err)
			}
			slog.Info("Created storage bucket", "bucket", bucketName)
		} else {
			slog.Info("Storage bucket already exists", "bucket", bucketName)
		}
	}

	slog.Info("Storage buckets initialized")
}

func createResendClient(config *configuration.Config) *resend.Client {
	resendClient := resend.NewClient(config.ResendApiKey)
	slog.Info("Resend client created")
	return resendClient
}

func createStripeClient(config *configuration.Config) *stripeGo.Client {
	stripeGo.Key = config.StripeSecretKey
	stripeClient := stripeGo.NewClient(config.StripeSecretKey)
	slog.Info("Stripe client created")
	return stripeClient
}

func runRiverMigrations(ctx context.Context, conn *pgxpool.Pool) {
	riverDriver := riverpgxv5.New(conn)
	migrator, err := rivermigrate.New(riverDriver, nil)
	if err != nil {
		log.Fatalf("Error creating River migrator, %s", err)
	}

	// Migrate up to the latest version
	// Empty MigrateOpts migrates all the way up to the latest version
	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		log.Fatalf("Error running River migrations, %s", err)
	}

	slog.Info("River migrations completed")
}

func createRiverClient(
	ctx context.Context,
	config *configuration.Config,
	conn *pgxpool.Pool,
	db *gorm.DB,
	resendClient *resend.Client,
	stripeClient *stripeGo.Client,
) *river.Client[pgx.Tx] {
	riverClient, err := jobs.NewRiverClient(ctx, jobs.RiverClientConfig{
		SQLConn:      conn,
		DB:           db,
		Config:       config,
		ResendClient: resendClient,
		StripeClient: stripeClient,
		StartWorkers: true, // Start workers in production
	})
	if err != nil {
		log.Fatalf("Error creating/starting river client, %s", err)
	}

	slog.Info("River client created and started")
	return riverClient
}

func createPostHogClient(config *configuration.Config) *posthog.Client {
	client := posthog.NewClient(config)
	slog.Info("PostHog client created")
	return client
}

func createEchoServer(
	config *configuration.Config,
	db *gorm.DB,
	minioClient *minio.Client,
	riverClient *river.Client[pgx.Tx],
	resendClient *resend.Client,
	stripeClient *stripeGo.Client,
	posthogClient *posthog.Client,
) *echo.Echo {
	e := echoServer.NewEcho(appMiddleware.AppDependencies{
		Config:        config,
		DB:            db,
		MinioClient:   minioClient,
		RiverClient:   riverClient,
		ResendClient:  resendClient,
		StripeClient:  stripeClient,
		PostHogClient: posthogClient,
	})

	return e
}

// gracefulShutdown handles graceful shutdown of all services
func gracefulShutdown(
	ctx context.Context,
	e *echo.Echo,
	riverClient *river.Client[pgx.Tx],
	posthogClient *posthog.Client,
	timeout time.Duration,
) {
	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Try soft stop for River (wait for jobs to finish)
	slog.Info("Initiating soft stop for River (waiting for jobs to finish)")
	softStopCtx, softStopCancel := context.WithTimeout(shutdownCtx, timeout)
	defer softStopCancel()

	riverStopped := make(chan error, 1)
	go func() {
		err := riverClient.Stop(softStopCtx)
		riverStopped <- err
	}()

	softStopSucceeded := false
	select {
	case err := <-riverStopped:
		if err == nil {
			slog.Info("River soft stop succeeded")
			softStopSucceeded = true
		} else if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			slog.Warn("River soft stop timeout; initiating hard stop")
		} else {
			slog.Warn("Error during River soft stop", "error", err)
		}
	case <-softStopCtx.Done():
		slog.Warn("River soft stop timeout; initiating hard stop")
	}

	// If soft stop didn't work, try hard stop
	if !softStopSucceeded {
		// Check if already stopped (non-blocking)
		select {
		case <-riverClient.Stopped():
			slog.Info("River already stopped")
		default:
			hardStopCtx, hardStopCancel := context.WithTimeout(shutdownCtx, timeout)
			defer hardStopCancel()

			slog.Info("Initiating hard stop for River (cancelling all jobs)")
			err := riverClient.StopAndCancel(hardStopCtx)
			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				slog.Warn("Error during River hard stop", "error", err)
			} else if errors.Is(err, context.DeadlineExceeded) {
				slog.Warn("River hard stop timeout; proceeding with shutdown")
			} else {
				slog.Info("River hard stop succeeded")
			}
		}
	}

	// Close PostHog client
	slog.Info("Closing PostHog client")
	if err := posthogClient.Close(); err != nil {
		slog.Warn("Error closing PostHog client", "error", err)
	} else {
		slog.Info("PostHog client closed")
	}

	// Shutdown Echo server
	slog.Info("Shutting down HTTP server")
	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Warn("Error shutting down HTTP server", "error", err)
	} else {
		slog.Info("HTTP server shut down successfully")
	}
}
