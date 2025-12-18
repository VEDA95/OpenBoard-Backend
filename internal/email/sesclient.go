package email

import (
	models "VEDA95/open_board/api/internal/db/model"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SESClient struct {
	client  *sesv2.Client
	context context.Context
}

func (mailClient *SESClient) Initialize(emailSettings *models.EmailSettings) error {
	if emailSettings.SESSecretAccessKey == nil || len(*emailSettings.SESSecretAccessKey) == 0 {
		return nil
	}

	if emailSettings.SESAccessKeyID == nil || len(*emailSettings.SESAccessKeyID) == 0 {
		return nil
	}

	creds := credentials.NewStaticCredentialsProvider(*emailSettings.SESAccessKeyID, *emailSettings.SESSecretAccessKey, "")
	options := [](func(*config.LoadOptions) error){
		config.WithCredentialsProvider(creds),
	}

	if len(emailSettings.SESRegion) > 0 {
		options = append(options, config.WithRegion(emailSettings.SESRegion))
	}

	mailClient.context = context.Background()
	conf, err := config.LoadDefaultConfig(mailClient.context, options...)
	if err != nil {
		return err
	}

	mailClient.client = sesv2.NewFromConfig(conf)

	return nil
}

func (mailClient *SESClient) Send(senderAddress string, recipient string, body string, subject *string, senderName *string) error {
	if mailClient.client == nil {
		return errors.New("ses client has not been initialized")
	}

	message := &sesv2.SendEmailInput{
		FromEmailAddress: &senderAddress,
		Destination:      &types.Destination{ToAddresses: []string{recipient}},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: subject},
				Body:    &types.Body{Html: &types.Content{Data: &body}},
			},
		},
	}

	if _, err := mailClient.client.SendEmail(mailClient.context, message); err != nil {
		return err
	}

	return nil
}
