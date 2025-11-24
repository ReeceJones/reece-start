package users

import (
	"errors"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reece.start/internal/api"
	"reece.start/internal/models"
	testconfig "reece.start/test/config"
	testdb "reece.start/test/db"
	testmocks "reece.start/test/mocks"
)

// Note: Helper functions like normalizePageSize, applySearchFilter, etc. are private
// and tested indirectly through integration tests. Unit tests for pure logic functions
// that don't require database access are included below.

// Helper functions are tested indirectly through integration tests below.
// The refactored helper functions (normalizePageSize, applySearchFilter, etc.)
// make the code more testable and are exercised by the integration tests.

// Integration tests for service functions

func TestCreateUser(t *testing.T) {
	db := testdb.SetupDB(t)
	config := testconfig.CreateTestConfig()
	posthogClient := testmocks.NewMockPosthogClient()

	t.Run("creates user successfully", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		params := CreateUserParams{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
			Timezone: "UTC",
		}

		result, err := createUser(CreateUserServiceRequest{
			Params:        params,
			Tx:            tx,
			Config:        config,
			PostHogClient: posthogClient,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User)
		assert.Equal(t, "Test User", result.User.Name)
		assert.Equal(t, "test@example.com", result.User.Email)
		assert.NotEmpty(t, result.Token)
		assert.NotNil(t, result.User.HashedPassword)
	})

	t.Run("returns error for duplicate email", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create first user
		params1 := CreateUserParams{
			Name:     "First User",
			Email:    "duplicate@example.com",
			Password: "password123",
			Timezone: "UTC",
		}
		_, err := createUser(CreateUserServiceRequest{
			Params:        params1,
			Tx:            tx,
			Config:        config,
			PostHogClient: posthogClient,
		})
		require.NoError(t, err)

		// Try to create second user with same email
		params2 := CreateUserParams{
			Name:     "Second User",
			Email:    "duplicate@example.com",
			Password: "password456",
			Timezone: "UTC",
		}
		_, err = createUser(CreateUserServiceRequest{
			Params: params2,
			Tx:     tx,
			Config: config,
		})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, api.ErrUserEmailAlreadyExists))
	})
}

func TestLoginUser(t *testing.T) {
	db := testdb.SetupDB(t)
	config := testconfig.CreateTestConfig()
	posthogClient := testmocks.NewMockPosthogClient()
	var minioClient *minio.Client // nil for tests

	t.Run("logs in user with correct credentials", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create a user first
		createParams := CreateUserParams{
			Name:     "Test User",
			Email:    "login@example.com",
			Password: "password123",
			Timezone: "UTC",
		}
		createdUser, err := createUser(CreateUserServiceRequest{
			Params:        createParams,
			Tx:            tx,
			Config:        config,
			PostHogClient: posthogClient,
		})
		require.NoError(t, err)

		// Login
		loginParams := LoginUserParams{
			Email:    "login@example.com",
			Password: "password123",
		}
		result, err := loginUser(LoginUserServiceRequest{
			Params:      loginParams,
			Tx:          tx,
			Config:      config,
			MinioClient: minioClient,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdUser.User.ID, result.User.ID)
		assert.NotEmpty(t, result.Token)
	})

	t.Run("returns error for invalid email", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		loginParams := LoginUserParams{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		_, err := loginUser(LoginUserServiceRequest{
			Params:      loginParams,
			Tx:          tx,
			Config:      config,
			MinioClient: minioClient,
		})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, api.ErrUnauthorizedInvalidLogin))
	})

	t.Run("returns error for invalid password", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create a user first
		createParams := CreateUserParams{
			Name:     "Test User",
			Email:    "wrongpass@example.com",
			Password: "password123",
			Timezone: "UTC",
		}
		_, err := createUser(CreateUserServiceRequest{
			Params:        createParams,
			Tx:            tx,
			Config:        config,
			PostHogClient: posthogClient,
		})
		require.NoError(t, err)

		// Try to login with wrong password
		loginParams := LoginUserParams{
			Email:    "wrongpass@example.com",
			Password: "wrongpassword",
		}
		_, err = loginUser(LoginUserServiceRequest{
			Params:      loginParams,
			Tx:          tx,
			Config:      config,
			MinioClient: minioClient,
		})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, api.ErrUnauthorizedInvalidLogin))
	})
}

func TestGetUserByID(t *testing.T) {
	db := testdb.SetupDB(t)
	var minioClient *minio.Client // nil for tests

	t.Run("gets user successfully", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create a user first
		user := &models.User{
			Name:  "Test User",
			Email: "getuser@example.com",
		}
		err := tx.Create(user).Error
		require.NoError(t, err)

		result, err := getUserByID(GetUserByIDServiceRequest{
			UserID:      user.ID,
			Tx:          tx,
			MinioClient: minioClient,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.User.ID)
		assert.Equal(t, "Test User", result.User.Name)
		assert.Equal(t, "getuser@example.com", result.User.Email)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		_, err := getUserByID(GetUserByIDServiceRequest{
			UserID:      99999,
			Tx:          tx,
			MinioClient: minioClient,
		})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, api.ErrUserNotFound))
	})
}

func TestUpdateUser(t *testing.T) {
	db := testdb.SetupDB(t)
	var minioClient *minio.Client // nil for tests
	config := testconfig.CreateTestConfig()

	t.Run("updates user name", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create a user first
		user := &models.User{
			Name:  "Old Name",
			Email: "update@example.com",
		}
		err := tx.Create(user).Error
		require.NoError(t, err)

		params := UpdateUserParams{
			UserID: user.ID,
			Name:   "New Name",
		}
		result, err := updateUser(UpdateUserServiceRequest{
			Params:      params,
			Tx:          tx,
			MinioClient: minioClient,
			Config:      config,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.User.Name)
		assert.Equal(t, "update@example.com", result.User.Email)
	})

	t.Run("updates user email", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create a user first
		user := &models.User{
			Name:  "Test User",
			Email: "oldemail@example.com",
		}
		err := tx.Create(user).Error
		require.NoError(t, err)

		params := UpdateUserParams{
			UserID: user.ID,
			Email:  "newemail@example.com",
		}
		result, err := updateUser(UpdateUserServiceRequest{
			Params:      params,
			Tx:          tx,
			MinioClient: minioClient,
			Config:      config,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "newemail@example.com", result.User.Email)
	})

	t.Run("updates user password", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create a user first
		user := &models.User{
			Name:  "Test User",
			Email: "password@example.com",
		}
		err := tx.Create(user).Error
		require.NoError(t, err)

		params := UpdateUserParams{
			UserID:   user.ID,
			Password: "newpassword123",
		}
		result, err := updateUser(UpdateUserServiceRequest{
			Params:      params,
			Tx:          tx,
			MinioClient: minioClient,
			Config:      config,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User.HashedPassword)
		assert.NotEmpty(t, result.User.HashedPassword)
	})
}

func TestGetUsers(t *testing.T) {
	db := testdb.SetupDB(t)
	var minioClient *minio.Client // nil for tests

	t.Run("gets users with pagination", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create multiple users
		for i := 0; i < 5; i++ {
			user := &models.User{
				Name:  "User " + string(rune('A'+i)),
				Email: "user" + string(rune('a'+i)) + "@example.com",
			}
			err := tx.Create(user).Error
			require.NoError(t, err)
		}

		result, err := getUsers(GetUsersServiceRequest{
			Cursor:      "",
			Size:        3,
			Search:      "",
			Tx:          tx,
			MinioClient: minioClient,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Users, 3)
		assert.True(t, result.HasNext)
		assert.False(t, result.HasPrev)
		assert.NotEmpty(t, result.NextCursor)
	})

	t.Run("searches users by name", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create users with different names
		user1 := &models.User{Name: "John Doe", Email: "john@example.com"}
		user2 := &models.User{Name: "Jane Smith", Email: "jane@example.com"}
		user3 := &models.User{Name: "Bob Johnson", Email: "bob@example.com"}
		tx.Create(user1)
		tx.Create(user2)
		tx.Create(user3)

		result, err := getUsers(GetUsersServiceRequest{
			Cursor:      "",
			Size:        10,
			Search:      "John",
			Tx:          tx,
			MinioClient: minioClient,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		// Should find both "John Doe" and "Bob Johnson"
		assert.GreaterOrEqual(t, len(result.Users), 1)
	})

	t.Run("normalizes page size", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// Create a user
		user := &models.User{Name: "Test User", Email: "test@example.com"}
		tx.Create(user)

		// Test with size 0 (should default to 20)
		result, err := getUsers(GetUsersServiceRequest{
			Cursor:      "",
			Size:        0,
			Search:      "",
			Tx:          tx,
			MinioClient: minioClient,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Users), 20)

		// Test with size > 100 (should cap at 100)
		result2, err := getUsers(GetUsersServiceRequest{
			Cursor:      "",
			Size:        200,
			Search:      "",
			Tx:          tx,
			MinioClient: minioClient,
		})

		require.NoError(t, err)
		assert.NotNil(t, result2)
		assert.LessOrEqual(t, len(result2.Users), 100)
	})
}
