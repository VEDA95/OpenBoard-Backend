package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"gorm.io/datatypes"
)

type Session struct {
	Base
	ExpiresOn        time.Time      `gorm:"not null" json:"-"`
	RefreshExpiresOn *time.Time     `json:"-"`
	Type             string         `gorm:"type:varchar(32);not null" json:"-"`
	RememberMe       bool           `gorm:"not null;default:false" json:"-"`
	AccessToken      string         `gorm:"type:text;not null" json:"-"`
	RefreshToken     *string        `gorm:"type:text" json:"-"`
	IPAddress        string         `gorm:"type:varchar(255);not null" json:"ip_address"`
	UserAgent        string         `gorm:"type:varchar(255);not null" json:"user_agent"`
	AdditionalInfo   datatypes.JSON `gorm:"type:jsonb" json:"-"`
	UserID           string         `gorm:"type:uuid;not null" json:"-"`
	User             *User          `gorm:"foreignKey:UserID" json:"-"`
}

func (*Session) createSessionToken() (string, error) {
	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func (session *Session) GenerateTokens() error {
	token, err := session.createSessionToken()
	if err != nil {
		return err
	}

	session.AccessToken = token

	if !session.RememberMe {
		return nil
	}

	refreshToken, err := session.createSessionToken()
	if err != nil {
		return err
	}

	session.RefreshToken = &refreshToken

	return nil
}

func (session *Session) IsValid() bool {
	return time.Now().Before(session.ExpiresOn)
}

func (session *Session) IsRefreshValid() bool {
	return session.RememberMe && time.Now().Before(*session.RefreshExpiresOn)
}
