package models

import (
	"VEDA95/open_board/api/internal/auth"
	"time"

	"gorm.io/datatypes"
)

type Session struct {
	Base
	ExpiresOn        time.Time      `gorm:"not null" json:"expires_on"`
	RefreshExpiresOn *time.Time     `json:"refresh_expires_on"`
	Type             string         `gorm:"type:varchar(32);not null" json:"session_type"`
	RememberMe       bool           `gorm:"not null;default:false" json:"remember_me"`
	AccessToken      string         `gorm:"type:text;not null" json:"access_token"`
	RefreshToken     *string        `gorm:"type:text" json:"refresh_token"`
	IPAddress        string         `gorm:"type:varchar(255);not null" json:"ip_address"`
	UserAgent        string         `gorm:"type:varchar(255);not null" json:"user_agent"`
	AdditionalInfo   datatypes.JSON `gorm:"type:jsonb" json:"additional_info"`
	UserID           string         `gorm:"type:uuid;not null" json:"user_id"`
}

func (session *Session) GenerateTokens() error {
	token, err := auth.CreateSessionToken()
	if err != nil {
		return err
	}

	session.AccessToken = token

	if !session.RememberMe {
		return nil
	}

	refreshToken, err := auth.CreateSessionToken()
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
