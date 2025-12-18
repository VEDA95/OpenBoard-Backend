package email

import (
	models "VEDA95/open_board/api/internal/db/model"

	"github.com/wneessen/go-mail"
)

type GoMailClient struct {
	client *mail.Client
}

func (mailClient *GoMailClient) Initialize(emailSettings *models.EmailSettings) error {
	if emailSettings.SMTPHost == nil || len(*emailSettings.SMTPHost) == 0 {
		return nil
	}

	if emailSettings.EmailFromAddress == nil || len(*emailSettings.EmailFromAddress) == 0 {
		return nil
	}

	if emailSettings.SMTPAuthMethod == nil || len(*emailSettings.SMTPAuthMethod) == 0 {
		return nil
	}

	if emailSettings.SMTPEncryption == nil || len(*emailSettings.SMTPEncryption) == 0 {
		return nil
	}

	clientOptions := []mail.Option{
		mail.WithPort(emailSettings.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthType(*emailSettings.SMTPAuthMethod)),
	}

	if *emailSettings.SMTPEncryption == "ssl" {
		clientOptions = append(clientOptions, mail.WithSSL())
	}

	if *emailSettings.SMTPEncryption == "tls" {
		clientOptions = append(clientOptions, mail.WithTLSPolicy(mail.TLSMandatory))
	}

	if mail.SMTPAuthType(*emailSettings.SMTPAuthMethod) == mail.SMTPAuthPlain {
		clientOptions = append(clientOptions, mail.WithUsername(*emailSettings.SMTPUsername), mail.WithPassword(*emailSettings.SMTPPassword))
	}

	client, err := mail.NewClient(*emailSettings.SMTPHost, clientOptions...)
	if err != nil {
		return err
	}

	mailClient.client = client

	return nil
}

func (mailClient *GoMailClient) Send(senderAddress string, recipient string, body string, subject *string, senderName *string) error {
	message := mail.NewMsg()

	if subject != nil && len(*subject) > 0 {
		message.Subject(*subject)
	}
	message.SetBodyString(mail.TypeTextHTML, body)

	if senderName != nil && len(*senderName) > 0 {
		if err := message.FromFormat(*senderName, senderAddress); err != nil {
			return err
		}
	} else {
		if err := message.From(senderAddress); err != nil {
			return err
		}
	}

	if err := message.To(recipient); err != nil {
		return err
	}

	if err := mailClient.client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
