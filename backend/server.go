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
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
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
)

func main() {
	setupLogger()

	config, err := configuration.LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("Error loading environment variables, %s", err)
	}

	conn, db := createDatabaseConnection(config)
	runDatabaseMigrations(db)

	minioClient := createMinioClient(config)
	initializeStorageBuckets(minioClient)

	resendClient := createResendClient(config)
	stripeClient := createStripeClient(config)

	ctx := context.Background()
	runRiverMigrations(ctx, conn)

	riverClient := createRiverClient(ctx, config, conn, db, resendClient, stripeClient)

	e := createEchoServer(config, db, minioClient, riverClient, resendClient, stripeClient)

	// Optional: Add body dump middleware for debugging (production only)
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody []byte, resBody []byte) {
		slog.Info("Body dump", "request", string(reqBody), "response", string(resBody))
	}))

	// Start http server
	e.Logger.Fatal(e.Start(":8080"))
}

func setupLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func createDatabaseConnection(config *configuration.Config) (*sql.DB, *gorm.DB) {
	conn, err := sql.Open("pgx", config.DatabaseUri)
	if err != nil {
		log.Fatalf("Error opening database, %s", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
	}

	slog.Info("Database connected")
	return conn, db
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

func runRiverMigrations(ctx context.Context, conn *sql.DB) {
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
}

func createRiverClient(
	ctx context.Context,
	config *configuration.Config,
	conn *sql.DB,
	db *gorm.DB,
	resendClient *resend.Client,
	stripeClient *stripeGo.Client,
) *river.Client[*sql.Tx] {
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

func createEchoServer(
	config *configuration.Config,
	db *gorm.DB,
	minioClient *minio.Client,
	riverClient *river.Client[*sql.Tx],
	resendClient *resend.Client,
	stripeClient *stripeGo.Client,
) *echo.Echo {
	e := echoServer.NewEcho(appMiddleware.AppDependencies{
		Config:       config,
		DB:           db,
		MinioClient:  minioClient,
		RiverClient:  riverClient,
		ResendClient: resendClient,
		StripeClient: stripeClient,
	})

	return e
}
