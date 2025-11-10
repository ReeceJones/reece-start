package email

import (
	"bytes"
	"os"
	"testing"

	"github.com/resend/resend-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reece.start/internal/configuration"
	"reece.start/testmocks"
)

func TestSendEmail(t *testing.T) {
	t.Run("EmailDisabled", func(t *testing.T) {
		config := &configuration.Config{
			EnableEmail: false,
		}

		resendClient := resend.NewClient("test-key")

		params := SendEmailParams{
			From:    "test@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<p>Test HTML</p>",
		}

		// Capture stdout to verify email content is printed
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Send email
		response, err := SendEmail(SendEmailRequest{
			Params:       params,
			ResendClient: resendClient,
			Config:       config,
		})

		// Restore stdout
		w.Close()
		os.Stdout = originalStdout
		buf := &bytes.Buffer{}
		buf.ReadFrom(r)

		// Assertions
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.NotEmpty(t, response.ID)
		// Verify email content was printed
		output := buf.String()
		assert.Contains(t, output, "From: test@example.com")
		assert.Contains(t, output, "To: [recipient@example.com]")
		assert.Contains(t, output, "Subject: Test Subject")
		assert.Contains(t, output, "Html: <p>Test HTML</p>")
	})

	t.Run("EmailEnabled_Success", func(t *testing.T) {
		// Setup: Replace default transport to intercept Resend API calls
		// This uses the common fixture from testmocks/transport.go
		testmocks.ReplaceDefaultTransportWithCleanup(t)

		config := &configuration.Config{
			EnableEmail:  true,
			ResendApiKey: "test-key",
		}

		resendClient := resend.NewClient("test-key")

		params := SendEmailParams{
			From:    "test@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test Subject",
			Html:    "<p>Test HTML</p>",
		}

		// Send email - should be intercepted by MockHTTPTransport from testmocks/transport.go
		response, err := SendEmail(SendEmailRequest{
			Params:       params,
			ResendClient: resendClient,
			Config:       config,
		})

		// Assertions
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.NotEmpty(t, response.ID)
	})

	t.Run("EmailEnabled_MultipleRecipients", func(t *testing.T) {
		// Setup: Replace default transport to intercept Resend API calls
		// This uses the common fixture from testmocks/transport.go
		testmocks.ReplaceDefaultTransportWithCleanup(t)

		config := &configuration.Config{
			EnableEmail:  true,
			ResendApiKey: "test-key",
		}

		resendClient := resend.NewClient("test-key")

		params := SendEmailParams{
			From:    "test@example.com",
			To:      []string{"recipient1@example.com", "recipient2@example.com"},
			Subject: "Test Subject",
			Html:    "<p>Test HTML</p>",
		}

		// Send email - should be intercepted by MockHTTPTransport from testmocks/transport.go
		response, err := SendEmail(SendEmailRequest{
			Params:       params,
			ResendClient: resendClient,
			Config:       config,
		})

		// Assertions
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.NotEmpty(t, response.ID)
	})
}
