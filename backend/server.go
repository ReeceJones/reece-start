package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/resend/resend-go/v2"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
	"reece.start/internal/database"
	echoServer "reece.start/internal/echo"
	"reece.start/internal/jobs"
	appMiddleware "reece.start/internal/middleware"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	slog.Info("Database connected")

	// Run database migrations
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Error migrating database, %s", err)
	}

	slog.Info("Database migrated")

	// Create minio client (Storage)
	minioClient, err := minio.New(config.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.StorageAccessKeyId, config.StorageSecretAccessKey, ""),
		Secure: config.StorageUseSSL,
	})
	if err != nil {
		log.Fatalf("Error creating minio client, %s", err)
	}

	slog.Info("Minio client created")

	// Create resend client (Email)
	resendClient := resend.NewClient(config.ResendApiKey)

	slog.Info("Resend client created")

	// Create stripe client and configure the global API key
	stripeGo.Key = config.StripeSecretKey
	stripeClient := stripeGo.NewClient(config.StripeSecretKey)

	// Run River migrations to create River tables
	// Use rivermigrate package to run migrations programmatically
	// See: https://riverqueue.com/docs/migrations#go-migration-api
	ctx := context.Background()
	riverDriver := riverdatabasesql.New(conn)
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

	// Create and start river client (Background jobs)
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

	// Create Echo server with all middleware and routes
	e := echoServer.NewEcho(appMiddleware.AppDependencies{
		Config:       config,
		DB:           db,
		MinioClient:  minioClient,
		RiverClient:  riverClient,
		ResendClient: resendClient,
		StripeClient: stripeClient,
	})

	// Optional: Add body dump middleware for debugging (production only)
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody []byte, resBody []byte) {
		slog.Info("Body dump", "request", string(reqBody), "response", string(resBody))
	}))

	// Start http server
	e.Logger.Fatal(e.Start(":8080"))
}
