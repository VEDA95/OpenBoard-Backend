package models

import (
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Role struct {
	ID          string       `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
}

func (role *Role) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	role.ID = id.String()
	return nil
}
