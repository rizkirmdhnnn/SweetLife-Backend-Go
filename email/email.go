package email

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// EmailClient is an interface to send email
type EmailClient interface {
	SendEmail(toEmail, subject, body string) error
}

// emailClient is an implementation of EmailClient
type emailClient struct {
	domain        string
	privateAPIKey string
	senderEmail   string
	mgClient      *mailgun.MailgunImpl
}

// NewEmailClient creates a new email client
func NewEmailClient(domain, privateAPIKey, senderEmail string) *emailClient {
	mg := mailgun.NewMailgun(domain, privateAPIKey)
	return &emailClient{
		domain:        domain,
		privateAPIKey: privateAPIKey,
		senderEmail:   senderEmail,
		mgClient:      mg,
	}
}

// SendEmail sends an email
func (e *emailClient) SendEmail(toEmail, subject, body string) error {
	// Create a new message
	message := e.mgClient.NewMessage(
		e.senderEmail,
		subject,
		body,
		toEmail,
	)

	// Set HTML body
	message.SetHtml(body)

	// Send email with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send email
	_, _, err := e.mgClient.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
