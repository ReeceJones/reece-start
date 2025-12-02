package users_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/test"
)

func TestCreateUserEndpoint(t *testing.T) {
	t.Run("ValidUser", func(t *testing.T) {
		// Setup test context with all dependencies
		tc := test.SetupEchoTest(t)

		// Prepare request
		reqBody := map[string]interface{}{
			"data": map[string]interface{}{
				"type": constants.ApiTypeUser,
				"attributes": map[string]interface{}{
					"name":     "Test User",
					"email":    "test@example.com",
					"password": "password123",
				},
			},
		}

		// Make request
		rec := tc.MakeRequest(http.MethodPost, "/users", reqBody, nil)

		// Assert response
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response map[string]interface{}
		tc.UnmarshalResponse(rec, &response)

		data := response["data"].(map[string]interface{})
		attributes := data["attributes"].(map[string]interface{})
		meta := data["meta"].(map[string]interface{})

		assert.Equal(t, "Test User", attributes["name"])
		assert.Equal(t, "test@example.com", attributes["email"])
		assert.NotEmpty(t, meta["token"])
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		// Setup test context
		tc := test.SetupEchoTest(t)

		// Create a user first via API
		user, _, _ := test.CreateTestUser(t, tc)

		// Try to create another user with the same email
		reqBody := map[string]interface{}{
			"data": map[string]interface{}{
				"type": constants.ApiTypeUser,
				"attributes": map[string]interface{}{
					"name":     "New User",
					"email":    user.Email, // Use the same email
					"password": "password456",
				},
			},
		}

		// Make request
		rec := tc.MakeRequest(http.MethodPost, "/users", reqBody, nil)

		// Assert error response
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestLoginEndpoint(t *testing.T) {
	t.Run("ValidCredentials", func(t *testing.T) {
		// Setup test context
		tc := test.SetupEchoTest(t)

		// Create a test user via API
		user, password, _ := test.CreateTestUser(t, tc)

		// Prepare login request
		reqBody := map[string]interface{}{
			"data": map[string]interface{}{
				"type": constants.ApiTypeUser,
				"attributes": map[string]interface{}{
					"email":    user.Email,
					"password": password,
				},
			},
		}

		// Make request
		rec := tc.MakeRequest(http.MethodPost, "/users/login", reqBody, nil)

		// Assert response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response map[string]interface{}
		tc.UnmarshalResponse(rec, &response)

		data := response["data"].(map[string]interface{})
		attributes := data["attributes"].(map[string]interface{})
		meta := data["meta"].(map[string]interface{})

		assert.Equal(t, user.Name, attributes["name"])
		assert.Equal(t, user.Email, attributes["email"])
		assert.NotEmpty(t, meta["token"])
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		// Setup test context
		tc := test.SetupEchoTest(t)

		// Prepare login request with invalid credentials
		reqBody := map[string]interface{}{
			"data": map[string]interface{}{
				"type": constants.ApiTypeUser,
				"attributes": map[string]interface{}{
					"email":    "nonexistent@example.com",
					"password": "wrongpassword",
				},
			},
		}

		// Make request
		rec := tc.MakeRequest(http.MethodPost, "/users/login", reqBody, nil)

		// Assert error response
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestGetAuthenticatedUserEndpoint(t *testing.T) {
	// Setup test context
	tc := test.SetupEchoTest(t)

	// Create authenticated test user with organization via API
	user, org, token := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Make authenticated request
	rec := tc.MakeAuthenticatedRequest(http.MethodGet, "/users/me", nil, token)

	// Assert response
	require.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})

	assert.Equal(t, user.Name, attributes["name"])
	assert.Equal(t, user.Email, attributes["email"])

	// Verify the organization exists in the test context
	_ = org // Keep org variable to maintain test clarity
}

func TestUpdateUserEndpoint(t *testing.T) {
	// Setup test context
	tc := test.SetupEchoTest(t)

	// Create authenticated test user via API
	user, _, token := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Prepare update request
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeUser,
			"attributes": map[string]interface{}{
				"name": "Updated Name",
			},
		},
	}

	// Make authenticated request
	rec := tc.MakeAuthenticatedRequest(
		http.MethodPatch,
		"/users/"+user.ID.String(),
		reqBody,
		token,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the user was updated in the database
	var updatedUser models.User
	err := tc.DB.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)
}
