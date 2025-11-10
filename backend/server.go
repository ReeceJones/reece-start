package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/resend/resend-go/v2"
	stripeGo "github.com/stripe/stripe-go/v83"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
	"reece.start/internal/database"
	"reece.start/internal/jobs"
	appMiddleware "reece.start/internal/middleware"
	"reece.start/internal/server"
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

	// Create and start river client (Background jobs)
	riverClient, err := jobs.NewRiverClient(context.Background(), jobs.RiverClientConfig{
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

	log.Printf("River client created and started\n")

	// Create Echo server with all middleware and routes
	e := server.NewEcho(appMiddleware.AppDependencies{
		Config:       config,
		DB:           db,
		MinioClient:  minioClient,
		RiverClient:  riverClient,
		ResendClient: resendClient,
		StripeClient: stripeClient,
	})

	// Optional: Add body dump middleware for debugging (production only)
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody []byte, resBody []byte) {
		log.Println(string(reqBody))
	}))

	// Start http server
	e.Logger.Fatal(e.Start(":8080"))
}
