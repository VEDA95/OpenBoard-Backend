package models

import (
	"time"

	"github.com/pquerna/otp/totp"
)

type PasswordResetToken struct {
	BaseID
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	ExpiresOn time.Time `gorm:"not null" json:"expires_on"`
	Type      string    `gorm:"type:varchar(16);not null" json:"type"`
	Token     string    `gorm:"type:varchar(6);not null" json:"token"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

func (passwordResetToken *PasswordResetToken) GenerateToken(email string) error {
	now := time.Now()
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "open_board",
		AccountName: email,
	})
	if err != nil {
		return err
	}

	token, err := totp.GenerateCode(key.Secret(), now)
	if err != nil {
		return err
	}

	passwordResetToken.Token = token
	passwordResetToken.ExpiresOn = now.Add(time.Minute * 15)

	return nil
}

func (passwordResetToken *PasswordResetToken) IsValid() bool {
	return time.Now().Before(passwordResetToken.ExpiresOn)
}
