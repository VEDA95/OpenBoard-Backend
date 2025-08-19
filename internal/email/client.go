package email

import (
	"github.com/wneessen/go-mail"
)

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
