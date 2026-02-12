package models

import "time"

type EmailVerificationToken struct {
	BaseID
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	ExpiresOn time.Time `gorm:"not null" json:"expires_on"`
	UserID    string    `gorm:"type:uuid; not null" json:"-"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
