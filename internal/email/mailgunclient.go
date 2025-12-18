package email

import (
	models "VEDA95/open_board/api/internal/db/model"
	"context"
	"errors"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

type MainlGunClient struct {
	client *mailgun.Client
	domain string
}

func (mailClient *MainlGunClient) Initialize(emailSettings *models.EmailSettings) error {
	if emailSettings.MailgunAPIKey == nil || len(*emailSettings.MailgunAPIKey) == 0 {
		return nil
	}

	if emailSettings.MailgunDomain == nil || len(*emailSettings.MailgunDomain) == 0 {
		return nil
	}

	mailClient.domain = *emailSettings.MailgunDomain
	mailClient.client = mailgun.NewMailgun(*emailSettings.MailgunAPIKey)

	return nil
}

func (mailClient *MainlGunClient) Send(senderAddress string, recipient string, body string, subject *string, senderName *string) error {
	if mailClient.client == nil {
		return errors.New("mailgun client is not initialized")
	}

	message := mailgun.NewMessage(mailClient.domain, senderAddress, *subject, "", recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	message.SetHTML(body)

	_, err := mailClient.client.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
