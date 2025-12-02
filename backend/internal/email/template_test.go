package email

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

func setupTemplateTest(t *testing.T) func() {
	// Get the current working directory
	wd, err := os.Getwd()
	require.NoError(t, err)

	// Find the backend directory by looking for go.mod
	backendDir := wd
	for {
		if _, err := os.Stat(filepath.Join(backendDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(backendDir)
		if parent == backendDir {
			t.Fatalf("Could not find backend directory (go.mod)")
		}
		backendDir = parent
	}

	// Change to backend directory
	err = os.Chdir(backendDir)
	require.NoError(t, err)

	// Return cleanup function to restore original directory
	return func() {
		os.Chdir(wd)
	}
}

func TestOrganizationInvitationEmailTemplateParams(t *testing.T) {
	t.Run("ApplyHtmlTemplate", func(t *testing.T) {
		cleanup := setupTemplateTest(t)
		defer cleanup()
		// Create test data
		invitingUser := models.User{
			Model: gorm.Model{
				ID: 1,
			},
			Name:  "John Doe",
			Email: "john@example.com",
		}

		organization := models.Organization{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Test Organization",
		}

		invitationID := uuid.New()
		invitation := models.OrganizationInvitation{
			Model: gorm.Model{
				ID: 1,
			},
			ID:             invitationID,
			Email:          "invitee@example.com",
			Role:           string(constants.OrganizationRoleAdmin),
			OrganizationID: uuid.New(),
			InvitingUserID: uuid.New(),
			Status:         string(constants.OrganizationInvitationStatusPending),
		}

		params := OrganizationInvitationEmailTemplateParams{
			InvitingUser:       invitingUser,
			Organization:       organization,
			Invitation:         invitation,
			FrontendUrl:        "http://localhost:3000",
			ServiceName:        constants.ServiceName,
			ServiceDescription: constants.ServiceDescription,
		}

		// Apply template
		html, err := params.ApplyHtmlTemplate()

		// Assertions
		require.NoError(t, err)
		require.NotEmpty(t, html)

		// Verify template content
		assert.Contains(t, html, invitingUser.Name)
		assert.Contains(t, html, organization.Name)
		assert.Contains(t, html, constants.ServiceName)
		assert.Contains(t, html, constants.ServiceDescription)
		assert.Contains(t, html, params.FrontendUrl)
		assert.Contains(t, html, invitationID.String())
	})
}

func TestApplyHtmlTemplate(t *testing.T) {
	t.Run("ValidTemplate", func(t *testing.T) {
		cleanup := setupTemplateTest(t)
		defer cleanup()
		// Create test data
		invitingUser := models.User{
			Model: gorm.Model{
				ID: 1,
			},
			Name:  "Jane Smith",
			Email: "jane@example.com",
		}

		organization := models.Organization{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "Acme Corp",
		}

		invitationID := uuid.New()
		invitation := models.OrganizationInvitation{
			Model: gorm.Model{
				ID: 1,
			},
			ID:             invitationID,
			Email:          "newuser@example.com",
			Role:           string(constants.OrganizationRoleMember),
			OrganizationID: uuid.New(),
			InvitingUserID: uuid.New(),
			Status:         string(constants.OrganizationInvitationStatusPending),
		}

		templateParams := OrganizationInvitationEmailTemplateParams{
			InvitingUser:       invitingUser,
			Organization:       organization,
			Invitation:         invitation,
			FrontendUrl:        "https://example.com",
			ServiceName:        constants.ServiceName,
			ServiceDescription: constants.ServiceDescription,
		}

		// Apply template using the method
		html, err := templateParams.ApplyHtmlTemplate()

		// Assertions
		require.NoError(t, err)
		require.NotEmpty(t, html)

		// Verify all expected content is present
		assert.Contains(t, html, invitingUser.Name)
		assert.Contains(t, html, organization.Name)
		assert.Contains(t, html, constants.ServiceName)
		assert.Contains(t, html, constants.ServiceDescription)
		assert.Contains(t, html, templateParams.FrontendUrl)
		assert.Contains(t, html, invitationID.String())

		// Verify HTML structure (should contain anchor tag)
		assert.Contains(t, html, "<a")
		assert.Contains(t, html, "</a>")
		assert.Contains(t, html, "href=")
	})
}
