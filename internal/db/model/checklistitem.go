package models

type CheckListItem struct {
	Base
	Name      string `gorm:"text;not null" json:"name"`
	IsChecked bool   `gorm:"not null;default:true" json:"is_checked" default:"true"`
	Position  int    `gorm:"not null" json:"position"`
	CardID    string `gorm:"type:uuid;not null" json:"-"`
	Card      *Card  `gorm:"foreignKey:CardID" json:"card,omitempty"`
}
