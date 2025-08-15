package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Base struct {
	ID        string     `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	base.ID = id.String()
	return nil
}
