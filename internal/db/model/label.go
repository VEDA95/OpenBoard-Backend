package models

type Label struct {
	Base
	Name    string `gorm:"type:varchar(255);not null" json:"name"`
	Color   string `gorm:"type:varchar(7);not null" json:"color"`
	BoardID string `gorm:"type:uuid;not null" json:"-"`
	Board   *Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}
