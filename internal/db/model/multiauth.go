package models

import "gorm.io/datatypes"

type MultiAuthMethod struct {
	Base
	Name        string         `gorm:"type:varchar(255); not null" json:"name"`
	Type        string         `gorm:"type:varchar(255); not null" json:"type"`
	Credentials datatypes.JSON `gorm:"jsonb" json:"-"`
	UserID      string         `gorm:"type:uuid" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
}
