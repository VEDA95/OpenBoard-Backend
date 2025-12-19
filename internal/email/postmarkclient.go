package email

import (
	models "VEDA95/open_board/api/internal/db/model"
	"errors"
	"fmt"

	"github.com/keighl/postmark"
)

type PostMarkClient struct {
	client *postmark.Client
}

func (mailClient *PostMarkClient) Initialize(emailSettings *models.EmailSettings) error {
	if emailSettings.PostmarkAccountToken == nil || len(*emailSettings.PostmarkAccountToken) == 0 {
		return nil
	}

	if emailSettings.PostmarkServerToken == nil || len(*emailSettings.PostmarkServerToken) == 0 {
		return nil
	}

	mailClient.client = postmark.NewClient(*emailSettings.PostmarkServerToken, *emailSettings.PostmarkAccountToken)

	return nil
}

func (mailClient *PostMarkClient) Send(senderAddress string, recipient string, body string, subject *string, senderName *string) error {
	if mailClient.client == nil {
		return errors.New("postmark client is not initialized")
	}

	sender := ""

	if subject != nil && len(*subject) > 0 {
		sender = fmt.Sprintf("%s <%s>", senderAddress, *senderName)
	} else {
		sender = senderAddress
	}

	message := postmark.Email{
		From:     sender,
		To:       recipient,
		Subject:  *subject,
		HtmlBody: body,
	}

	if _, err := mailClient.client.SendEmail(message); err != nil {
		return err
	}

	return nil
}
