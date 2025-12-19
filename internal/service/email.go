package service

import (
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/email"
	"bytes"
	"errors"
	"text/template"

	"gorm.io/gorm"
)

type EmailService struct {
	emailRepo         *repository.EmailRepository
	emailSettingsRepo *repository.EmailSettingsRepository
	client            email.EmailClient
}

func NewEmailService(emailRepo *repository.EmailRepository, emailSettingsRepo *repository.EmailSettingsRepository) (*EmailService, error) {
	emailService := &EmailService{
		emailRepo:         emailRepo,
		emailSettingsRepo: emailSettingsRepo,
	}

	if err := emailService.InitializeClient(); err != nil {
		return nil, err
	}

	return emailService, nil
}

func (emailService *EmailService) InitializeClient() error {
	emailSettings, err := emailService.emailSettingsRepo.Find()

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if emailSettings == nil {
		return errors.New("unable to fetch email settings")
	}

	if !emailSettings.EmailEnabled {
		if emailService.client != nil {
			emailService.client = nil
		}

		return nil
	}

	if emailSettings.EmailFromAddress == nil || len(*emailSettings.EmailFromAddress) == 0 {
		return errors.New("valid email address for the sender is required")
	}

	if emailSettings.EmailFromName != nil && len(*emailSettings.EmailFromName) > 0 {
		emailService.emailRepo.SetSenderName(emailSettings.EmailFromName)
	}

	emailService.emailRepo.SetSender(*emailSettings.EmailFromAddress)

	if emailSettings.EmailProvider == "smtp" {
		emailService.client = &email.GoMailClient{}
	}

	if emailSettings.EmailProvider == "sendgrid" {
		emailService.client = &email.SendGridClient{}
	}

	if emailSettings.EmailProvider == "mailgun" {
		emailService.client = &email.MailGunClient{}
	}

	if emailSettings.EmailProvider == "ses" {
		emailService.client = &email.SESClient{}
	}

	if emailSettings.EmailProvider == "postmark" {
		emailService.client = &email.PostMarkClient{}
	}

	if err := emailService.client.Initialize(emailSettings); err != nil {
		return err
	}

	return nil
}

func (emailService *EmailService) RenderTemplate(name string, variables any) string {
	if emailService.emailRepo.GetTemplateCount() == 0 {
		return ""
	}

	rawTemplate := emailService.emailRepo.GetTemplate(name)

	if len(rawTemplate) == 0 {
		return rawTemplate
	}

	templateInstance, err := template.New(name).Parse(rawTemplate)
	if err != nil {
		return ""
	}

	buffer := new(bytes.Buffer)

	if err := templateInstance.Execute(buffer, variables); err != nil {
		return ""
	}

	return buffer.String()
}

func (emailService *EmailService) SendMessage(recipient string, templateName string, templatePayload any, subject *string) error {
	if emailService.client == nil {
		return errors.New("email client has not been initialized")
	}

	err := emailService.client.Send(
		emailService.emailRepo.GetSender(),
		recipient,
		emailService.RenderTemplate(templateName, templatePayload),
		subject,
		emailService.emailRepo.GetSenderName(),
	)
	if err != nil {
		return err
	}

	return nil
}
