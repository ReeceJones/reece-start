package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reece.start/internal/database"
)

var (
	sharedPostgresContainer testcontainers.Container
	sharedPostgresOnce      sync.Once
)

// getSharedConnStr returns the connection string for the shared PostgreSQL container
func getSharedConnStr() (string, error) {
	if sharedPostgresContainer == nil {
		return "", fmt.Errorf("shared postgres container not initialized")
	}
	ctx := context.Background()
	endpoint, err := sharedPostgresContainer.Endpoint(ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get endpoint: %w", err)
	}
	return fmt.Sprintf("postgres://test:test@%s/testdb?sslmode=disable", endpoint), nil
}

// setupSharedPostgresContainer starts a single PostgreSQL testcontainer that will be reused across all tests
func setupSharedPostgresContainer(t *testing.T) {
	sharedPostgresOnce.Do(func() {
		ctx := context.Background()

		postgresC, err := testcontainers.Run(
			ctx, "postgres:latest",
			testcontainers.WithEnv(map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "testdb",
			}),
			testcontainers.WithExposedPorts("5432/tcp"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			),
		)
		require.NoError(t, err)

		sharedPostgresContainer = postgresC

		// Get the connection string endpoint
		connStr, err := getSharedConnStr()
		require.NoError(t, err)

		// Create an initial connection to run migrations
		conn, err := sql.Open("pgx", connStr)
		require.NoError(t, err)

		// Verify the connection works by pinging
		err = conn.Ping()
		require.NoError(t, err)

		db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})
		require.NoError(t, err)

		// Run migrations once
		err = database.Migrate(db)
		require.NoError(t, err)

		// Run River migrations to create River tables
		// Use rivermigrate package to run migrations programmatically
		// See: https://riverqueue.com/docs/migrations#go-migration-api
		riverDriver := riverdatabasesql.New(conn)
		migrator, err := rivermigrate.New(riverDriver, nil)
		require.NoError(t, err, "Failed to create River migrator")

		// Migrate up to the latest version
		// Empty MigrateOpts migrates all the way up to the latest version
		_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
		require.NoError(t, err, "Failed to run River migrations")

		// Close the initial connection
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	})
}

// CleanAllTables truncates all tables in the database to reset state between tests
// It excludes River tables (river_*) which should persist across tests
func CleanAllTables(db *gorm.DB) error {
	// Get all table names, excluding River tables
	var tables []string
	if err := db.Raw(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename NOT LIKE 'river_%'
	`).Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	if len(tables) == 0 {
		return nil
	}

	// Truncate all tables using CASCADE to handle foreign key constraints
	// RESTART IDENTITY resets auto-increment sequences
	// Quote table names to handle any special characters safely
	query := fmt.Sprintf(`TRUNCATE TABLE "%s"`, tables[0])
	for _, table := range tables[1:] {
		query += fmt.Sprintf(`, "%s"`, table)
	}
	query += ` RESTART IDENTITY CASCADE`

	if err := db.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	return nil
}

// SetupDB returns a GORM database connection for service-level tests.
// This uses a shared PostgreSQL container that is reused across all tests.
// The database is automatically cleaned (tables truncated) before each test.
func SetupDB(t *testing.T) *gorm.DB {
	// Ensure the shared container is started
	setupSharedPostgresContainer(t)

	// Get the connection string (fresh each time to ensure it's current)
	connStr, err := getSharedConnStr()
	require.NoError(t, err)

	// Create a new connection from the shared container
	conn, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	// Verify the connection works by pinging
	err = conn.Ping()
	require.NoError(t, err)

	// Create GORM connection
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})
	require.NoError(t, err)

	// Clean all tables before each test
	err = CleanAllTables(db)
	require.NoError(t, err)

	// Register cleanup to close this connection
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	})

	return db
}

// SetupPostgresContainer returns a new database connection using the shared PostgreSQL container
// This is used by HTTP tests that need both sql.DB and gorm.DB connections
func SetupPostgresContainer(t *testing.T) (*sql.DB, *gorm.DB, string) {
	// Ensure the shared container is started
	setupSharedPostgresContainer(t)

	// Get the connection string (fresh each time to ensure it's current)
	connStr, err := getSharedConnStr()
	require.NoError(t, err)

	// Create a new connection from the shared container
	pool, err := pgxpool.New(context.Background(), connStr)
	require.NoError(t, err)

	connector := stdlib.GetPoolConnector(pool)
	sqlConn := sql.OpenDB(connector)

	sqlConn.SetMaxIdleConns(0)

	// Create the gorm connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlConn,
	}), &gorm.Config{})
	require.NoError(t, err)

	// Clean all tables before each test
	err = CleanAllTables(db)
	require.NoError(t, err)

	// Clean River jobs between tests (but keep River tables)
	// This ensures test isolation while preserving River table structure
	err = CleanRiverJobs(db)
	require.NoError(t, err)

	// Register cleanup to close this connection
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	})

	// Note: We don't register cleanup for the shared container here because
	// t.Cleanup() runs when each test finishes, not when all tests finish.
	// The container will be automatically cleaned up by testcontainers when
	// the test process exits, or we can rely on Docker's container lifecycle.

	return sqlConn, db, connStr
}

// CleanRiverJobs deletes all River jobs between tests to ensure test isolation
// This is separate from CleanAllTables because we want to keep River tables but clean jobs
func CleanRiverJobs(db *gorm.DB) error {
	// Delete all jobs from river_job table
	err := db.Exec(`DELETE FROM river_job`).Error
	if err != nil {
		return fmt.Errorf("failed to clean River jobs: %w", err)
	}
	return nil
}
