package email

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
	"reece.start/internal/configuration"
)

type SendEmailParams struct {
	From    string
	To      []string
	Subject string
	Html    string
}

type SendEmailRequest struct {
	Params       SendEmailParams
	ResendClient *resend.Client
	Config       *configuration.Config
}

type SendEmailResponse struct {
	ID string
}

func SendEmail(request SendEmailRequest) (*SendEmailResponse, error) {

	if !request.Config.EnableEmail {
		// Print the email content
		fmt.Printf("From: %s\nTo: %s\nSubject: %s\nHtml: %s\n", request.Params.From, request.Params.To, request.Params.Subject, request.Params.Html)

		// if email is disabled, return a random id
		id := uuid.New().String()
		return &SendEmailResponse{
			ID: id,
		}, nil
	}

	resp, err := request.ResendClient.Emails.Send(&resend.SendEmailRequest{
		From:    request.Params.From,
		To:      request.Params.To,
		Subject: request.Params.Subject,
		Html:    request.Params.Html,
	})

	if err != nil {
		return nil, err
	}

	return &SendEmailResponse{
		ID: resp.Id,
	}, nil
}
