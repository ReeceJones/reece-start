package email

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"reece.start/internal/models"
)

const templatePath = "internal/email/templates/*.html"

type HtmlTemplateParams struct {
	Template string
	Params   interface{}
}

type OrganizationInvitationEmailTemplateParams struct {
	InvitingUser       models.User
	Organization       models.Organization
	Invitation         models.OrganizationInvitation
	FrontendUrl        string
	ServiceName        string
	ServiceDescription string
}

func (params OrganizationInvitationEmailTemplateParams) ApplyHtmlTemplate() (string, error) {
	return applyHtmlTemplate(HtmlTemplateParams{
		Template: "organizationInvitationEmail",
		Params:   params,
	})
}

func applyHtmlTemplate(params HtmlTemplateParams) (string, error) {
	// Resolve template path relative to backend directory
	// This ensures templates can be found regardless of the current working directory
	templatePathResolved := resolveTemplatePath(templatePath)

	tmpl, err := template.ParseGlob(templatePathResolved)
	if err != nil {
		return "", err
	}

	// Execute the template
	buffer := bytes.Buffer{}
	err = tmpl.Execute(&buffer, params.Params)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// resolveTemplatePath finds the backend directory and resolves the template path relative to it
// This ensures templates can be found when running tests or from different working directories
func resolveTemplatePath(relativePath string) string {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		// If we can't get working directory, return original path
		return relativePath
	}

	// Find the backend directory by looking for go.mod
	backendDir := wd
	for {
		if _, err := os.Stat(filepath.Join(backendDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(backendDir)
		if parent == backendDir {
			// If we can't find go.mod, return original path
			return relativePath
		}
		backendDir = parent
	}

	// Return absolute path relative to backend directory
	return filepath.Join(backendDir, relativePath)
}
