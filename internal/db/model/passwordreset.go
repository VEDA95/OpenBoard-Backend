package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	ExpiresOn time.Time `gorm:"not null" json:"expires_on"`
	Type      string    `gorm:"type:varchar(16);not null" json:"type"`
	Token     string    `gorm:"type:varchar(6);not null" json:"token"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
}

func (resetToken *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	resetToken.ID = id.String()
	return nil
}
