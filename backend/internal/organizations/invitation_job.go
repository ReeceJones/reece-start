package organizations

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/resend/resend-go/v2"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
	"reece.start/internal/constants"
	"reece.start/internal/email"
	"reece.start/internal/models"
)


type OrganizationInvitationEmailJobArgs struct {
	InvitationId uint `json:"invitationId"`
}

func (OrganizationInvitationEmailJobArgs) Kind() string {
	return string(constants.JobKindOrganizationInvitationEmail)
}

type OrganizationInvitationEmailJobWorker struct {
	river.WorkerDefaults[OrganizationInvitationEmailJobArgs]
	DB          *gorm.DB
	Config      *configuration.Config
	ResendClient *resend.Client
}

type OrganizationInvitationHtmlTemplateParams struct {
	InvitingUser models.User
	Organization models.Organization
	Invitation models.OrganizationInvitation
	FrontendUrl string
}


func (w *OrganizationInvitationEmailJobWorker) Work(ctx context.Context, job *river.Job[OrganizationInvitationEmailJobArgs]) error {
	log.Printf("Sending organization invitation email %d", job.Args.InvitationId)

	// Get the inviting user
	var invitation models.OrganizationInvitation
	err := w.DB.Model(&models.OrganizationInvitation{}).Preload("InvitingUser").Preload("Organization").Where("id = ?", job.Args.InvitationId).First(&invitation).Error
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s invited you to join %s", invitation.InvitingUser.Name, invitation.Organization.Name)
	html, err := email.OrganizationInvitationEmailTemplateParams{
		InvitingUser: invitation.InvitingUser,
		Organization: invitation.Organization,
		Invitation: invitation,
		FrontendUrl: w.Config.FrontendUrl,
		ServiceName: constants.ServiceName,
		ServiceDescription: constants.ServiceDescription,
	}.ApplyHtmlTemplate()

	if err != nil {
		return err
	}
	
	w.ResendClient.Emails.Send(&resend.SendEmailRequest{
		From: string(constants.EmailSenderDefault),
		To: []string{invitation.Email},
		Subject: subject,
		Html: html,
	})

	return nil
}

func (w *OrganizationInvitationEmailJobWorker) Timeout(*river.Job[OrganizationInvitationEmailJobArgs]) time.Duration {
	return 180 * time.Second
}
