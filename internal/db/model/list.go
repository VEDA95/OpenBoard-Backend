package models

type List struct {
	Base
	Name     string  `gorm:"type:varchar(255);not null" json:"name"`
	Color    *string `gorm:"type:varchar(7)" json:"color"`
	Position int     `gorm:"not null" json:"position"`
	BoardID  string  `gorm:"type:uuid;not null" json:"-"`
	Board    *Board  `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	Cards    []*Card `gorm:"foreignKey:ListID" json:"cards"`
}
