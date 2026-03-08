package models

import (
	"time"

	"gorm.io/datatypes"
)

type PendingMultiAuthMethod struct {
	Base
	Type        string         `gorm:"type:varchar(255); not null" json:"type"`
	UserID      string         `gorm:"type:uuid" json:"-"`
	ExpiresAt   time.Time      `gorm:"not null" json:"expires_at"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	Credentials datatypes.JSON `gorm:"jsonb" json:"-"`
}
