package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailSettings struct {
	ID                       int        `gorm:"primaryKey;check:check_single_row_email,id = 1" json:"id"`
	UpdatedAt                *time.Time `json:"updated_at"`
	EmailProvider            string     `gorm:"type:varchar(20);default:'smtp'" json:"email_provider"` // smtp, sendgrid, mailgun, ses, postmark
	EmailEnabled             bool       `gorm:"default:true" json:"email_enabled"`
	EmailFromAddress         *string    `gorm:"type:varchar(255)" json:"email_from_address"`
	EmailFromName            *string    `gorm:"type:varchar(100)" json:"email_from_name"`
	SMTPHost                 *string    `gorm:"type:varchar(255)" json:"smtp_host"`
	SMTPPort                 int        `gorm:"default:587" json:"smtp_port"`
	SMTPUsername             *string    `gorm:"type:varchar(255)" json:"smtp_username"`
	SMTPPassword             *string    `gorm:"type:varchar(255)" json:"smtp_password"`
	SMTPEncryption           *string    `gorm:"type:varchar(10);default:'tls'" json:"smtp_encryption"`    // tls, ssl, none
	SMTPAuthMethod           *string    `gorm:"type:varchar(20);default:'plain'" json:"smtp_auth_method"` // plain, login, cram-md5
	SMTPVerifySSL            bool       `gorm:"default:true" json:"smtp_verify_ssl"`
	SendGridAPIKey           *string    `gorm:"type:varchar(255)" json:"sendgrid_api_key"`
	MailgunAPIKey            *string    `gorm:"type:varchar(255)" json:"mailgun_api_key"`
	MailgunDomain            *string    `gorm:"type:varchar(255)" json:"mailgun_domain"`
	SESAccessKeyID           *string    `gorm:"type:varchar(255)" json:"ses_access_key_id"`
	SESSecretAccessKey       *string    `gorm:"type:varchar(255)" json:"ses_secret_access_key"`
	SESRegion                string     `gorm:"type:varchar(50);default:'us-east-1'" json:"ses_region"`
	PostmarkServerToken      *string    `gorm:"type:varchar(255)" json:"postmark_server_token"`
	PostmarkAccountToken     *string    `gorm:"type:varchar(255)" json:"postmark_account_token"`
	EmailFooterText          *string    `gorm:"type:text" json:"email_footer_text"`
	EnableEmailNotifications bool       `gorm:"default:true" json:"enable_email_notifications"`
	SendPasswordResetEmail   bool       `gorm:"default:true" json:"send_password_reset_email"`
	NotifyOnCardAssigned     bool       `gorm:"default:true" json:"notify_on_card_assigned"`
	NotifyOnCardComment      bool       `gorm:"default:true" json:"notify_on_card_comment"`
	NotifyOnCardDue          bool       `gorm:"default:true" json:"notify_on_card_due"`
	NotifyOnBoardInvite      bool       `gorm:"default:true" json:"notify_on_board_invite"`
	NotifyOnMention          bool       `gorm:"default:true" json:"notify_on_mention"`
}

func (e *EmailSettings) BeforeCreate(tx *gorm.DB) error {
	e.ID = 1
	return nil
}
