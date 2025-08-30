package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type BaseID struct {
	ID string `gorm:"type:uuid;primaryKey" json:"id"`
}

type Base struct {
	BaseID
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (base *BaseID) BeforeCreate(tx *gorm.DB) error {
	if len(base.ID) > 0 {
		return nil
	}

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	base.ID = id.String()
	return nil
}
