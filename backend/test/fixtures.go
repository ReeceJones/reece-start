package test

// Fixtures for integration tests
//
// All fixture functions auto-generate test data (names, emails, passwords, etc.) using UUIDs
// to ensure uniqueness across test runs. This makes tests cleaner and reduces boilerplate.
//
// Basic Usage:
//   // Create a simple user
//   user, password, token := test.CreateTestUser(t, tc)
//
//   // Create an authenticated user with organization (returns all you need for testing)
//   user, org, token := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)
//
// Custom Data:
//   // Override specific fields when needed
//   user, password, token := test.CreateTestUserWithOptions(t, tc, test.TestUserOptions{
//       Email: "specific@example.com",
//   })
//
// All fixture functions use the API endpoints (not direct DB access) to ensure true integration testing.

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

// generateTestEmail generates a unique test email address
func generateTestEmail() string {
	return fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8])
}

// generateTestName generates a test name
func generateTestName() string {
	return fmt.Sprintf("Test User %s", uuid.New().String()[:8])
}

// generateTestOrgName generates a test organization name
func generateTestOrgName() string {
	return fmt.Sprintf("Test Org %s", uuid.New().String()[:8])
}

// generateTestPassword returns a consistent test password
func generateTestPassword() string {
	return "testPassword123!"
}

// TestUserOptions allows customizing test user creation
type TestUserOptions struct {
	Name     string
	Email    string
	Password string
}

// CreateTestUser creates a test user via the API with auto-generated details
// Returns the user model, password (for login), and JWT token
func CreateTestUser(t *testing.T, tc *TestContext) (*models.User, string, string) {
	return CreateTestUserWithOptions(t, tc, TestUserOptions{})
}

// CreateTestUserWithOptions creates a test user via the API with custom options
// Returns the user model, password (for login), and JWT token
func CreateTestUserWithOptions(t *testing.T, tc *TestContext, opts TestUserOptions) (*models.User, string, string) {
	if opts.Name == "" {
		opts.Name = generateTestName()
	}
	if opts.Email == "" {
		opts.Email = generateTestEmail()
	}
	if opts.Password == "" {
		opts.Password = generateTestPassword()
	}

	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeUser,
			"attributes": map[string]interface{}{
				"name":     opts.Name,
				"email":    opts.Email,
				"password": opts.Password,
			},
		},
	}

	rec := tc.MakeRequest(http.MethodPost, "/users", reqBody, nil)
	require.Equal(t, http.StatusCreated, rec.Code, "Failed to create user")

	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].(map[string]interface{})
	userIDStr := data["id"].(string)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	require.NoError(t, err)

	meta := data["meta"].(map[string]interface{})
	token := meta["token"].(string)

	// Fetch the user from DB to return the full model
	var user models.User
	err = tc.DB.First(&user, uint(userID)).Error
	require.NoError(t, err)

	return &user, opts.Password, token
}

// LoginTestUser logs in a user via the API and returns the JWT token
func LoginTestUser(t *testing.T, tc *TestContext, email, password string) string {
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeUser,
			"attributes": map[string]interface{}{
				"email":    email,
				"password": password,
			},
		},
	}

	rec := tc.MakeRequest(http.MethodPost, "/users/login", reqBody, nil)
	require.Equal(t, http.StatusOK, rec.Code, "Failed to login user")

	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].(map[string]interface{})
	meta := data["meta"].(map[string]interface{})
	token := meta["token"].(string)

	return token
}

// TestOrganizationOptions allows customizing test organization creation
type TestOrganizationOptions struct {
	Name string
}

// CreateTestOrganization creates a test organization via the API with auto-generated details
func CreateTestOrganization(t *testing.T, tc *TestContext, token string) *models.Organization {
	return CreateTestOrganizationWithOptions(t, tc, token, TestOrganizationOptions{})
}

// CreateTestOrganizationWithOptions creates a test organization via the API with custom options
func CreateTestOrganizationWithOptions(t *testing.T, tc *TestContext, token string, opts TestOrganizationOptions) *models.Organization {
	if opts.Name == "" {
		opts.Name = generateTestOrgName()
	}

	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganization,
			"attributes": map[string]interface{}{
				"name":       opts.Name,
				"locale":     "en-US",
				"entityType": "llc",
				"address": map[string]interface{}{
					"line1":           "123 Test St",
					"city":            "Test City",
					"stateOrProvince": "CA",
					"zip":             "12345",
					"country":         "US",
				},
			},
		},
	}

	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organizations", reqBody, token)
	require.Equal(t, http.StatusCreated, rec.Code, "Failed to create organization")

	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].(map[string]interface{})
	orgIDStr := data["id"].(string)
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	require.NoError(t, err)

	// Fetch the organization from DB to return the full model
	var org models.Organization
	err = tc.DB.First(&org, uint(orgID)).Error
	require.NoError(t, err)

	return &org
}

// CreateTestOrganizationMembership creates a test organization membership via the API and returns it
func CreateTestOrganizationMembership(t *testing.T, tc *TestContext, userID, orgID uint, role constants.OrganizationRole, adminToken string) *models.OrganizationMembership {
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationMembership,
			"attributes": map[string]interface{}{
				"role": string(role),
			},
			"relationships": map[string]interface{}{
				"user": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(userID), 10),
						"type": constants.ApiTypeUser,
					},
				},
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(orgID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}

	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-memberships", reqBody, adminToken)
	require.Equal(t, http.StatusCreated, rec.Code, "Failed to create organization membership")

	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].(map[string]interface{})
	membershipIDStr := data["id"].(string)
	membershipID, err := strconv.ParseUint(membershipIDStr, 10, 32)
	require.NoError(t, err)

	// Fetch the membership from DB to return the full model
	var membership models.OrganizationMembership
	err = tc.DB.First(&membership, uint(membershipID)).Error
	require.NoError(t, err)

	return &membership
}

// TestUserWithOrgOptions allows customizing test user and organization creation
type TestUserWithOrgOptions struct {
	UserName     string
	UserEmail    string
	UserPassword string
	OrgName      string
}

// CreateTestUserWithOrganization creates a test user with an organization via the API
// The user will be automatically assigned as admin when creating the organization
// Returns user, organization, membership, password (for re-login), and JWT token
func CreateTestUserWithOrganization(t *testing.T, tc *TestContext) (*models.User, *models.Organization, *models.OrganizationMembership, string, string) {
	return CreateTestUserWithOrganizationWithOptions(t, tc, TestUserWithOrgOptions{})
}

// CreateTestUserWithOrganizationWithOptions creates a test user with an organization via the API with custom options
// Returns user, organization, membership, password (for re-login), and JWT token
func CreateTestUserWithOrganizationWithOptions(t *testing.T, tc *TestContext, opts TestUserWithOrgOptions) (*models.User, *models.Organization, *models.OrganizationMembership, string, string) {
	// Create user and get initial token
	user, password, token := CreateTestUserWithOptions(t, tc, TestUserOptions{
		Name:     opts.UserName,
		Email:    opts.UserEmail,
		Password: opts.UserPassword,
	})

	// Create organization (this automatically creates admin membership)
	org := CreateTestOrganizationWithOptions(t, tc, token, TestOrganizationOptions{
		Name: opts.OrgName,
	})

	// Fetch the membership that was created
	var membership models.OrganizationMembership
	err := tc.DB.Where("user_id = ? AND organization_id = ?", user.ID, org.ID).First(&membership).Error
	require.NoError(t, err)

	// Get a new token with organization context
	token = LoginTestUser(t, tc, user.Email, password)

	return user, org, &membership, password, token
}

// CreateAuthenticatedTestUser creates a test user with an organization via the API
// Returns the user, organization, and JWT token
// Currently only supports admin role (the default when creating an organization)
func CreateAuthenticatedTestUser(t *testing.T, tc *TestContext, role constants.OrganizationRole) (*models.User, *models.Organization, string) {
	return CreateAuthenticatedTestUserWithOptions(t, tc, role, TestUserWithOrgOptions{})
}

// CreateAuthenticatedTestUserWithOptions creates a test user with an organization via the API with custom options
// Returns the user, organization, and JWT token
func CreateAuthenticatedTestUserWithOptions(t *testing.T, tc *TestContext, role constants.OrganizationRole, opts TestUserWithOrgOptions) (*models.User, *models.Organization, string) {
	user, org, _, _, token := CreateTestUserWithOrganizationWithOptions(t, tc, opts)

	// If the desired role is not admin, we need to update the membership
	// For now, this assumes the role is admin since organization creation makes the user an admin
	// If you need a non-admin role, you would need to create another admin user first
	if role != constants.OrganizationRoleAdmin {
		t.Fatalf("CreateAuthenticatedTestUser currently only supports admin role. Use CreateTestUserWithOrganization and CreateTestOrganizationMembership for other roles.")
	}

	return user, org, token
}

// RunPendingRiverJobs starts the River client temporarily and polls until all jobs are completed.
// This is useful for testing job enqueueing and execution.
// Jobs will run only once (max_attempts set to 1) and will not retry on failure.
// It will poll the job queue every pollInterval until no jobs are pending/running or maxWaitTime is reached.
func RunPendingRiverJobs(t *testing.T, db *gorm.DB, riverClient *river.Client[*sql.Tx], maxWaitTime time.Duration, pollInterval time.Duration) {
	ctx := context.Background()

	// Set all pending jobs to run only once (no retries)
	err := db.Exec(`
		UPDATE river_job 
		SET max_attempts = 1 
		WHERE state IN ('available', 'scheduled', 'retryable')
	`).Error
	require.NoError(t, err)

	// Start the River client to process jobs
	err = riverClient.Start(ctx)
	require.NoError(t, err)

	// Ensure we stop the client when done
	defer func() {
		err := riverClient.Stop(ctx)
		require.NoError(t, err)
	}()

	startTime := time.Now()

	// Poll until all jobs are completed or timeout
	for {
		// Check if we've exceeded max wait time
		if time.Since(startTime) > maxWaitTime {
			t.Fatalf("Timeout waiting for jobs to complete after %v", maxWaitTime)
			return
		}

		// Query the database directly to check for pending/running jobs
		// Jobs in 'available', 'running', 'retryable', or 'scheduled' states are still in progress
		var pendingCount int64
		err := db.Raw(`
			SELECT COUNT(*) 
			FROM river_job 
			WHERE state IN ('available', 'running', 'retryable', 'scheduled')
		`).Scan(&pendingCount).Error
		require.NoError(t, err)

		if pendingCount == 0 {
			// All jobs are completed or discarded - check if any failed
			var results struct {
				Completed int64
				Discarded int64
				Cancelled int64
			}
			err := db.Raw(`
				SELECT 
					COUNT(CASE WHEN state = 'completed' THEN 1 END) as completed,
					COUNT(CASE WHEN state = 'discarded' THEN 1 END) as discarded,
					COUNT(CASE WHEN state = 'cancelled' THEN 1 END) as cancelled
				FROM river_job
			`).Scan(&results).Error
			require.NoError(t, err)

			t.Logf("All jobs completed in %v (completed: %d, discarded: %d, cancelled: %d)",
				time.Since(startTime), results.Completed, results.Discarded, results.Cancelled)

			// Fail test if any jobs were discarded (failed)
			if results.Discarded > 0 {
				t.Fatalf("%d job(s) failed and were discarded", results.Discarded)
			}

			return
		}

		t.Logf("Waiting for %d jobs to complete...", pendingCount)
		time.Sleep(pollInterval)
	}
}

// RunAllPendingRiverJobs processes all pending River jobs with sensible defaults.
// It will wait up to 10 seconds and poll every 100ms.
// This is a convenience wrapper around RunPendingRiverJobs.
func RunAllPendingRiverJobs(t *testing.T, db *gorm.DB, riverClient *river.Client[*sql.Tx]) {
	RunPendingRiverJobs(t, db, riverClient, 10*time.Second, 100*time.Millisecond)
}
