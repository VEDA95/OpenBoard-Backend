package models

type Role struct {
	BaseID
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
}
