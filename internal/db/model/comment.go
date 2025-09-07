package models

type Comment struct {
	Base
	Comment string `gorm:"type:text;not null" json:"comment"`
	UserID  string `gorm:"type:uuid;not null" json:"-"`
	CardID  string `gorm:"type:uuid;not null" json:"-"`
	User    *User  `gorm:"foreignKey:UserID" json:"user"`
	Card    *Card  `gorm:"foreignKey:CardID" json:"card,omitempty"`
}
