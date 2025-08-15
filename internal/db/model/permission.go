package models

import (
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Permission struct {
	ID   string `gorm:"type:uuid;primaryKey" json:"id"`
	Path string
}

func (permission *Permission) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	permission.ID = id.String()
	return nil
}
