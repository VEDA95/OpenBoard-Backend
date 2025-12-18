package email

import (
	models "VEDA95/open_board/api/internal/db/model"
)

type EmailClient interface {
	Initialize(emailSettings *models.EmailSettings) error
	Send(senderAddress string, recipient string, body string, subject *string, senderName *string) error
}
