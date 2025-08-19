package service

import (
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/email"
	"bytes"
	"errors"
	"os"
	"strconv"
	"text/template"

	"github.com/wneessen/go-mail"
)

type EmailService struct {
	repo   *repository.EmailRepository
	client *mail.Client
}

func NewEmailService(emailRepo *repository.EmailRepository) (*EmailService, error) {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpTLS := os.Getenv("SMTP_TLS")
	emailSender := os.Getenv("SENDER_EMAIL")

	if len(smtpServer) == 0 || len(smtpPort) == 0 || len(smtpUsername) == 0 || len(smtpPassword) == 0 || len(emailSender) == 0 || len(smtpTLS) == 0 {
		return nil, errors.New("required email settings are missing")
	}

	parsedTLSBool, err := strconv.ParseBool(smtpTLS)
	if err != nil {
		return nil, err
	}

	parsedPort, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, err
	}

	client, err := email.CreateEmailClient(smtpServer, parsedPort, mail.SMTPAuthPlain, parsedTLSBool, smtpUsername, smtpPassword)
	if err != nil {
		return nil, err
	}

	emailRepo.SetSender(emailSender)

	return &EmailService{
		repo:   emailRepo,
		client: client,
	}, nil
}

func (emailService *EmailService) RenderTemplate(name string, variables any) string {
	if emailService.repo.GetTemplateCount() == 0 {
		return ""
	}

	rawTemplate := emailService.repo.GetTemplate(name)

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

func (mailService *EmailService) SendMessage(subject string, recipient string, templateName string, templatePayload any) error {
	message := mail.NewMsg()

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, mailService.RenderTemplate(templateName, templatePayload))

	if err := message.From(mailService.repo.GetSender()); err != nil {
		return err
	}

	if err := message.To(recipient); err != nil {
		return err
	}

	if err := mailService.client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
