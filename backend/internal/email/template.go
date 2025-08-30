package email

import (
	"bytes"
	"html/template"

	"reece.start/internal/models"
)

const templatePath = "internal/email/templates/*.html"

type HtmlTemplateParams struct {
	Template string
	Params interface{}
}

type OrganizationInvitationEmailTemplateParams struct {
	InvitingUser models.User
	Organization models.Organization
	Invitation models.OrganizationInvitation
	FrontendUrl string
	ServiceName string
	ServiceDescription string
}

func (params OrganizationInvitationEmailTemplateParams) ApplyHtmlTemplate() (string, error) {
	return applyHtmlTemplate(HtmlTemplateParams{
		Template: "organizationInvitationEmail",
		Params: params,
	})
}


func applyHtmlTemplate(params HtmlTemplateParams) (string, error) {
	tmpl, err := template.ParseGlob(templatePath)
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