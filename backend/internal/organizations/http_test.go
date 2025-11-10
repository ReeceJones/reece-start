package organizations_test

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/test"
)

// createTokenWithOrganizationContext creates a token with organization context for testing
func createTokenWithOrganizationContext(t *testing.T, tc *test.TestContext, initialToken string, orgID uint) string {
	tokenReqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeToken,
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(orgID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}
	tokenRec := tc.MakeAuthenticatedRequest(http.MethodPost, "/users/me/token", tokenReqBody, initialToken)
	require.Equal(t, http.StatusOK, tokenRec.Code)

	var tokenResponse map[string]interface{}
	tc.UnmarshalResponse(tokenRec, &tokenResponse)
	tokenData := tokenResponse["data"].(map[string]interface{})
	tokenMeta := tokenData["meta"].(map[string]interface{})
	return tokenMeta["token"].(string)
}

func TestCreateOrganizationEndpoint(t *testing.T) {
	t.Run("ValidOrganization", func(t *testing.T) {
		tc := test.SetupEchoTest(t)

		// Create authenticated user
		user, _, token := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

		// Prepare request
		reqBody := map[string]interface{}{
			"data": map[string]interface{}{
				"type": constants.ApiTypeOrganization,
				"attributes": map[string]interface{}{
					"name":        "New Organization",
					"description": "Test Description",
					"locale":      "en-US",
					"entityType":  "llc",
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

		// Make request
		rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organizations", reqBody, token)

		// Assert response
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse response
		var response map[string]interface{}
		tc.UnmarshalResponse(rec, &response)

		data := response["data"].(map[string]interface{})
		attributes := data["attributes"].(map[string]interface{})
		meta := data["meta"].(map[string]interface{})

		assert.Equal(t, "New Organization", attributes["name"])
		assert.Equal(t, "Test Description", attributes["description"])
		assert.NotNil(t, meta["onboardingStatus"])

		// Verify organization was created in database
		var org models.Organization
		err := tc.DB.Where("name = ?", "New Organization").First(&org).Error
		require.NoError(t, err)

		// Verify membership was created for the user
		var membership models.OrganizationMembership
		err = tc.DB.Where("user_id = ? AND organization_id = ?", user.ID, org.ID).First(&membership).Error
		require.NoError(t, err)
		assert.Equal(t, string(constants.OrganizationRoleAdmin), membership.Role)
	})
}

func TestGetOrganizationsEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, token := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create another organization for the same user
	org2 := test.CreateTestOrganization(t, tc, token)

	// Make request
	rec := tc.MakeAuthenticatedRequest(http.MethodGet, "/organizations", nil, token)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2)

	// Verify both organizations are in the response
	orgIDs := make(map[uint]bool)
	for _, item := range data {
		orgData := item.(map[string]interface{})
		orgIDStr := orgData["id"].(string)
		orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
		require.NoError(t, err)
		orgIDs[uint(orgID)] = true
	}

	assert.True(t, orgIDs[org.ID])
	assert.True(t, orgIDs[org2.ID])
}

func TestGetOrganizationEndpoint(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tc := test.SetupEchoTest(t)

		// Create authenticated user with organization
		user, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

		// Create a token with organization context
		token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

		// Make request
		rec := tc.MakeAuthenticatedRequest(http.MethodGet, "/organizations/"+strconv.FormatUint(uint64(org.ID), 10), nil, token)

		// Assert response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response map[string]interface{}
		tc.UnmarshalResponse(rec, &response)

		data := response["data"].(map[string]interface{})
		attributes := data["attributes"].(map[string]interface{})

		assert.Equal(t, strconv.FormatUint(uint64(org.ID), 10), data["id"])
		assert.Equal(t, org.Name, attributes["name"])

		// Keep user variable to avoid unused variable warning
		_ = user
	})

	t.Run("Unauthorized", func(t *testing.T) {
		tc := test.SetupEchoTest(t)

		// Create two users with separate organizations
		_, org1, _ := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)
		_, _, token2 := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

		// Try to access org1 with token2 (should fail)
		rec := tc.MakeAuthenticatedRequest(http.MethodGet, "/organizations/"+strconv.FormatUint(uint64(org1.ID), 10), nil, token2)

		// Assert unauthorized response
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestUpdateOrganizationEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

	// Prepare update request
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganization,
			"attributes": map[string]interface{}{
				"name":        "Updated Name",
				"description": "Updated Description",
			},
		},
	}

	// Make request
	rec := tc.MakeAuthenticatedRequest(
		http.MethodPatch,
		"/organizations/"+strconv.FormatUint(uint64(org.ID), 10),
		reqBody,
		token,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the organization was updated in the database
	var updatedOrg models.Organization
	err := tc.DB.First(&updatedOrg, org.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedOrg.Name)
	assert.Equal(t, "Updated Description", updatedOrg.Description)
}

func TestDeleteOrganizationEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	user, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

	// Make request
	rec := tc.MakeAuthenticatedRequest(
		http.MethodDelete,
		"/organizations/"+strconv.FormatUint(uint64(org.ID), 10),
		nil,
		token,
	)

	// Assert response
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify organization is deleted
	var deletedOrg models.Organization
	err := tc.DB.First(&deletedOrg, org.ID).Error
	require.Error(t, err)

	// Verify membership is also deleted
	var membership models.OrganizationMembership
	err = tc.DB.Where("user_id = ? AND organization_id = ?", user.ID, org.ID).First(&membership).Error
	require.Error(t, err)
}

func TestGetOrganizationMembershipsEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	user1, org, initialToken1 := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token1 := createTokenWithOrganizationContext(t, tc, initialToken1, org.ID)

	// Create another user
	user2, _, _ := test.CreateTestUser(t, tc)

	// Add user2 to the organization
	test.CreateTestOrganizationMembership(t, tc, user2.ID, org.ID, constants.OrganizationRoleMember, token1)

	// Make request
	rec := tc.MakeAuthenticatedRequest(
		http.MethodGet,
		"/organization-memberships?organizationId="+strconv.FormatUint(uint64(org.ID), 10),
		nil,
		token1,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	// Verify both users are in the response
	userIDs := make(map[uint]bool)
	included := response["included"].([]interface{})
	for _, item := range included {
		userData := item.(map[string]interface{})
		userIDStr := userData["id"].(string)
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		require.NoError(t, err)
		userIDs[uint(userID)] = true
	}

	assert.True(t, userIDs[user1.ID])
	assert.True(t, userIDs[user2.ID])
}

func TestGetOrganizationMembershipEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	user, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

	// Get the membership ID
	var membership models.OrganizationMembership
	err := tc.DB.Where("user_id = ? AND organization_id = ?", user.ID, org.ID).First(&membership).Error
	require.NoError(t, err)

	// Make request
	rec := tc.MakeAuthenticatedRequest(
		http.MethodGet,
		"/organization-memberships/"+strconv.FormatUint(uint64(membership.ID), 10),
		nil,
		token,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})

	assert.Equal(t, strconv.FormatUint(uint64(membership.ID), 10), data["id"])
	assert.Equal(t, string(constants.OrganizationRoleAdmin), attributes["role"])
}

func TestCreateOrganizationMembershipEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialAdminToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	adminToken := createTokenWithOrganizationContext(t, tc, initialAdminToken, org.ID)

	// Create another user
	user2, _, _ := test.CreateTestUser(t, tc)

	// Prepare request
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationMembership,
			"attributes": map[string]interface{}{
				"role": string(constants.OrganizationRoleMember),
			},
			"relationships": map[string]interface{}{
				"user": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(user2.ID), 10),
						"type": constants.ApiTypeUser,
					},
				},
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}

	// Make request
	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-memberships", reqBody, adminToken)

	// Assert response
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verify membership was created
	var membership models.OrganizationMembership
	err := tc.DB.Where("user_id = ? AND organization_id = ?", user2.ID, org.ID).First(&membership).Error
	require.NoError(t, err)
	assert.Equal(t, string(constants.OrganizationRoleMember), membership.Role)
}

func TestUpdateOrganizationMembershipEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialAdminToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	adminToken := createTokenWithOrganizationContext(t, tc, initialAdminToken, org.ID)

	// Create another user and add them as member
	user2, _, _ := test.CreateTestUser(t, tc)
	membership := test.CreateTestOrganizationMembership(t, tc, user2.ID, org.ID, constants.OrganizationRoleMember, adminToken)

	// Prepare update request
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationMembership,
			"attributes": map[string]interface{}{
				"role": string(constants.OrganizationRoleAdmin),
			},
		},
	}

	// Make request
	rec := tc.MakeAuthenticatedRequest(
		http.MethodPatch,
		"/organization-memberships/"+strconv.FormatUint(uint64(membership.ID), 10),
		reqBody,
		adminToken,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify membership was updated
	var updatedMembership models.OrganizationMembership
	err := tc.DB.First(&updatedMembership, membership.ID).Error
	require.NoError(t, err)
	assert.Equal(t, string(constants.OrganizationRoleAdmin), updatedMembership.Role)
}

func TestDeleteOrganizationMembershipEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialAdminToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	adminToken := createTokenWithOrganizationContext(t, tc, initialAdminToken, org.ID)

	// Create another user and add them as member
	user2, _, _ := test.CreateTestUser(t, tc)
	membership := test.CreateTestOrganizationMembership(t, tc, user2.ID, org.ID, constants.OrganizationRoleMember, adminToken)

	// Make request
	rec := tc.MakeAuthenticatedRequest(
		http.MethodDelete,
		"/organization-memberships/"+strconv.FormatUint(uint64(membership.ID), 10),
		nil,
		adminToken,
	)

	// Assert response
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify membership is deleted
	var deletedMembership models.OrganizationMembership
	err := tc.DB.First(&deletedMembership, membership.ID).Error
	require.Error(t, err)
}

func TestInviteToOrganizationEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

	// Prepare request
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationInvitation,
			"attributes": map[string]interface{}{
				"email": "invitee@example.com",
				"role":  string(constants.OrganizationRoleAdmin),
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}

	// Make request
	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-invitations", reqBody, token)

	// Assert response
	require.Equal(t, http.StatusCreated, rec.Code, "Expected 201 Created, got %d: %s", rec.Code, rec.Body.String())

	// Parse response
	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	require.NotNil(t, response["data"], "Response data should not be nil")
	data := response["data"].(map[string]interface{})
	require.NotNil(t, data["attributes"], "Response attributes should not be nil")
	attributes := data["attributes"].(map[string]interface{})

	assert.Equal(t, "invitee@example.com", attributes["email"])
	assert.Equal(t, string(constants.OrganizationRoleAdmin), attributes["role"])
	assert.Equal(t, string(constants.OrganizationInvitationStatusPending), attributes["status"])

	// Verify invitation was created in database
	var invitation models.OrganizationInvitation
	err := tc.DB.Where("email = ? AND organization_id = ?", "invitee@example.com", org.ID).First(&invitation).Error
	require.NoError(t, err)
}

func TestGetOrganizationInvitationsEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

	// Create invitations via API
	reqBody1 := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationInvitation,
			"attributes": map[string]interface{}{
				"email": "invitee1@example.com",
				"role":  string(constants.OrganizationRoleAdmin),
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}
	tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-invitations", reqBody1, token)

	reqBody2 := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationInvitation,
			"attributes": map[string]interface{}{
				"email": "invitee2@example.com",
				"role":  string(constants.OrganizationRoleMember),
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}
	tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-invitations", reqBody2, token)

	// Make request to get invitations
	rec := tc.MakeAuthenticatedRequest(
		http.MethodGet,
		"/organization-invitations?organizationId="+strconv.FormatUint(uint64(org.ID), 10),
		nil,
		token,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2)
}

func TestGetOrganizationInvitationEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

	// Create invitation via API
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationInvitation,
			"attributes": map[string]interface{}{
				"email": "invitee@example.com",
				"role":  string(constants.OrganizationRoleAdmin),
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}
	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-invitations", reqBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createResponse map[string]interface{}
	tc.UnmarshalResponse(rec, &createResponse)
	invitationID := createResponse["data"].(map[string]interface{})["id"].(string)

	// Make request to get invitation
	rec = tc.MakeAuthenticatedRequest(
		http.MethodGet,
		"/organization-invitations/"+invitationID,
		nil,
		token,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response map[string]interface{}
	tc.UnmarshalResponse(rec, &response)

	data := response["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})

	assert.Equal(t, invitationID, data["id"])
	assert.Equal(t, "invitee@example.com", attributes["email"])
}

func TestDeleteOrganizationInvitationEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	token := createTokenWithOrganizationContext(t, tc, initialToken, org.ID)

	// Create invitation via API
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationInvitation,
			"attributes": map[string]interface{}{
				"email": "invitee@example.com",
				"role":  string(constants.OrganizationRoleAdmin),
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}
	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-invitations", reqBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createResponse map[string]interface{}
	tc.UnmarshalResponse(rec, &createResponse)
	invitationID := createResponse["data"].(map[string]interface{})["id"].(string)

	// Make request to delete invitation
	rec = tc.MakeAuthenticatedRequest(
		http.MethodDelete,
		"/organization-invitations/"+invitationID,
		nil,
		token,
	)

	// Assert response
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify invitation status is revoked
	var invitation models.OrganizationInvitation
	err := tc.DB.Where("id = ?", invitationID).First(&invitation).Error
	require.NoError(t, err)
	assert.Equal(t, string(constants.OrganizationInvitationStatusRevoked), invitation.Status)
}

func TestAcceptOrganizationInvitationEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialAdminToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	adminToken := createTokenWithOrganizationContext(t, tc, initialAdminToken, org.ID)

	// Create invitation for a new user
	inviteeEmail := "invitee@example.com"
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationInvitation,
			"attributes": map[string]interface{}{
				"email": inviteeEmail,
				"role":  string(constants.OrganizationRoleMember),
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}
	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-invitations", reqBody, adminToken)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createResponse map[string]interface{}
	tc.UnmarshalResponse(rec, &createResponse)
	invitationID := createResponse["data"].(map[string]interface{})["id"].(string)

	// Create the invitee user
	inviteeUser, _, inviteeToken := test.CreateTestUserWithOptions(t, tc, test.TestUserOptions{
		Email: inviteeEmail,
	})

	// Make request to accept invitation
	acceptReqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"id":   invitationID,
			"type": constants.ApiTypeOrganizationInvitation,
		},
	}
	rec = tc.MakeAuthenticatedRequest(
		http.MethodPost,
		"/organization-invitations/"+invitationID+"/accept",
		acceptReqBody,
		inviteeToken,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify invitation status is accepted
	var invitation models.OrganizationInvitation
	err := tc.DB.Where("id = ?", invitationID).First(&invitation).Error
	require.NoError(t, err)
	assert.Equal(t, string(constants.OrganizationInvitationStatusAccepted), invitation.Status)

	// Verify membership was created
	var membership models.OrganizationMembership
	err = tc.DB.Where("user_id = ? AND organization_id = ?", inviteeUser.ID, org.ID).First(&membership).Error
	require.NoError(t, err)
	assert.Equal(t, string(constants.OrganizationRoleMember), membership.Role)
}

func TestDeclineOrganizationInvitationEndpoint(t *testing.T) {
	tc := test.SetupEchoTest(t)

	// Create authenticated user with organization
	_, org, initialAdminToken := test.CreateAuthenticatedTestUser(t, tc, constants.OrganizationRoleAdmin)

	// Create a token with organization context
	adminToken := createTokenWithOrganizationContext(t, tc, initialAdminToken, org.ID)

	// Create invitation for a new user
	inviteeEmail := "invitee@example.com"
	reqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": constants.ApiTypeOrganizationInvitation,
			"attributes": map[string]interface{}{
				"email": inviteeEmail,
				"role":  string(constants.OrganizationRoleMember),
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   strconv.FormatUint(uint64(org.ID), 10),
						"type": constants.ApiTypeOrganization,
					},
				},
			},
		},
	}
	rec := tc.MakeAuthenticatedRequest(http.MethodPost, "/organization-invitations", reqBody, adminToken)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createResponse map[string]interface{}
	tc.UnmarshalResponse(rec, &createResponse)
	invitationID := createResponse["data"].(map[string]interface{})["id"].(string)

	// Create the invitee user
	_, _, inviteeToken := test.CreateTestUserWithOptions(t, tc, test.TestUserOptions{
		Email: inviteeEmail,
	})

	// Make request to decline invitation
	declineReqBody := map[string]interface{}{
		"data": map[string]interface{}{
			"id":   invitationID,
			"type": constants.ApiTypeOrganizationInvitation,
		},
	}
	rec = tc.MakeAuthenticatedRequest(
		http.MethodPost,
		"/organization-invitations/"+invitationID+"/decline",
		declineReqBody,
		inviteeToken,
	)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify invitation status is declined
	var invitation models.OrganizationInvitation
	err := tc.DB.Where("id = ?", invitationID).First(&invitation).Error
	require.NoError(t, err)
	assert.Equal(t, string(constants.OrganizationInvitationStatusDeclined), invitation.Status)

	// Verify membership was NOT created
	// Get the invitee user ID from the invitation email
	var inviteeUser models.User
	err = tc.DB.Where("email = ?", inviteeEmail).First(&inviteeUser).Error
	require.NoError(t, err)

	var membership models.OrganizationMembership
	err = tc.DB.Where("organization_id = ? AND user_id = ?", org.ID, inviteeUser.ID).First(&membership).Error
	require.Error(t, err)
}
