package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/configuration"
	"reece.start/internal/database"
	appMiddleware "reece.start/internal/middleware"
	"reece.start/internal/organizations"
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

	// Run database migrations
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Error migrating database, %s", err)
	}

	// Create minio client (Storage)
	log.Printf("Using storage access key id: %s", config.StorageAccessKeyId)

	minioClient, err := minio.New(config.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.StorageAccessKeyId, config.StorageSecretAccessKey, ""),
		Secure: config.StorageUseSSL,
	})
	if err != nil {
		log.Fatalf("Error creating minio client, %s", err)
	}

	// Create river client (Background jobs)
	workers := river.NewWorkers()
	// river.AddWorker(workers, &invoices.InvoiceUploadJobWorker{
	// 	DB:          db,
	// 	Config:      config,
	// 	MinioClient: minioClient,
	// })

	riverClient, err := river.NewClient(riverdatabasesql.New(conn), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})
	if err != nil {
		log.Fatalf("Error creating river client, %s", err)
	}

	// err = riverClient.Start(context.Background())
	// if err != nil {
	// 	log.Fatalf("Error starting river client, %s", err)
	// }

	// log.Printf("River client started\n")

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
	}))

	// health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Public user routes (no authentication required)
	e.POST("/users", api.Validated(users.CreateUserEndpoint))
	e.POST("/users/login", api.Validated(users.LoginEndpoint))

	// Protected user routes (authentication required)
	protected := e.Group("")
	protected.Use(appMiddleware.JwtAuthMiddleware(config))

	protected.GET("/users/me", users.GetAuthenticatedUserEndpoint)
	protected.PUT("/users/:id", api.Validated(users.UpdateUserEndpoint))

	protected.GET("/organizations", organizations.GetOrganizationsEndpoint)
	protected.POST("/organizations", api.Validated(organizations.CreateOrganizationEndpoint))
	protected.GET("/organizations/:id", organizations.GetOrganizationEndpoint)
	protected.PUT("/organizations/:id", api.Validated(organizations.UpdateOrganizationEndpoint))
	protected.DELETE("/organizations/:id", organizations.DeleteOrganizationEndpoint)

	// Start http server
	e.Logger.Fatal(e.Start(":8080"))
}
