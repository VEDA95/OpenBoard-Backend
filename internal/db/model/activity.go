package models

type CardActivity struct {
	Base
	Activity string `gorm:"type:varchar(255);not null" json:"activity"`
	UserID   string `gorm:"type:uuid;not null" json:"-"`
	CardID   string `gorm:"type:uuid;not null" json:"-"`
	User     *User  `gorm:"foreignKey:UserID" json:"user"`
	Card     *Card  `gorm:"foreignKey:CardID" json:"card,omitempty"`
}
