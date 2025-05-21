package email

import (
	"errors"
	"github.com/wneessen/go-mail"
	"os"
	"strconv"
)

var MailClient *Client

type Client struct {
	client *mail.Client
	sender string
}

func CreateEmailClient(server string, port int, authMethod mail.SMTPAuthType, tls bool, username string, password string) (*mail.Client, error) {
	clientOptions := []mail.Option{
		mail.WithPort(port),
		mail.WithSMTPAuth(authMethod),
	}

	if tls {
		clientOptions = append(clientOptions, mail.WithTLSPolicy(mail.TLSMandatory))
	}

	if authMethod == mail.SMTPAuthPlain {
		clientOptions = append(clientOptions, mail.WithUsername(username), mail.WithPassword(password))
	}

	mailClient, err := mail.NewClient(server, clientOptions...)

	if err != nil {
		return nil, err
	}

	return mailClient, nil
}

func InitializeEmailClient() error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpTLS := os.Getenv("SMTP_TLS")
	emailSender := os.Getenv("SENDER_EMAIL")

	if len(smtpServer) == 0 || len(smtpPort) == 0 || len(smtpUsername) == 0 || len(smtpPassword) == 0 || len(emailSender) == 0 || len(smtpTLS) == 0 {
		return errors.New("required email settings are missing")
	}

	parsedTLSBool, err := strconv.ParseBool(smtpTLS)

	if err != nil {
		return err
	}

	parsedPort, err := strconv.Atoi(smtpPort)

	if err != nil {
		return err
	}

	client, err := CreateEmailClient(smtpServer, parsedPort, mail.SMTPAuthPlain, parsedTLSBool, smtpUsername, smtpPassword)

	if err != nil {
		return err
	}

	MailClient = &Client{client: client, sender: emailSender}

	return nil
}

func (mailClient *Client) SendMessage(subject string, recipient string, contentType mail.ContentType, messageText string) error {
	if mailClient.client == nil {
		return errors.New("email client is nil. skipping sending email")
	}

	message := mail.NewMsg()

	message.Subject(subject)
	message.SetBodyString(contentType, messageText)

	if err := message.From(mailClient.sender); err != nil {
		return err
	}

	if err := message.To(recipient); err != nil {
		return err
	}

	if err := mailClient.client.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
