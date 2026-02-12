package models

import "gorm.io/datatypes"

type FileUpload struct {
	Base
	Name              string         `gorm:"type:varchar(255);not null" json:"description"`
	Extension         string         `gorm:"type:varchar(8);not null" json:"extension"`
	Type              string         `gorm:"type:varchar(64);not null" json:"type"`
	Path              string         `gorm:"type:text;not null" json:"path"`
	Size              int            `gorm:"not null" json:"size"`
	AdditionalDetails datatypes.JSON `json:"additional_details"`
	UserID            string         `gorm:"type:uuid;not null" json:"-"`
	User              *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
