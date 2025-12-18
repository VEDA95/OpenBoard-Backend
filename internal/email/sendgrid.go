package email

import (
	models "VEDA95/open_board/api/internal/db/model"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridClient struct {
	client *sendgrid.Client
}

func (mailClient *SendGridClient) Initialize(emailSettings *models.EmailSettings) error {
	if emailSettings.SendGridAPIKey == nil || len(*emailSettings.SendGridAPIKey) == 0 {
		return nil
	}

	mailClient.client = sendgrid.NewSendClient(*emailSettings.SendGridAPIKey)

	return nil
}

func (mailClient *SendGridClient) Send(senderAddress string, recipient string, body string, subject *string, senderName *string) error {
	from := mail.NewEmail(*senderName, senderAddress)
	to := mail.NewEmail("", recipient)
	message := mail.NewSingleEmail(from, *subject, to, "", body)

	if _, err := mailClient.client.Send(message); err != nil {
		return err
	}

	return nil
}
