package models

import "time"

type Card struct {
	Base
	Name               string          `gorm:"type:varchar(255);not null" json:"name"`
	Description        *string         `gorm:"type:text" json:"description"`
	Color              *string         `gorm:"type:varchar(7)" json:"color"`
	Position           int             `gorm:"not null" json:"position"`
	ReminderDate       *time.Time      `json:"reminder_date"`
	DueDate            *time.Time      `json:"due_date"`
	TimeSpent          *int            `json:"time_spent"`
	EstimatedTimeSpent *int            `json:"estimated_time_spent"`
	IsActive           bool            `gorm:"not null;default:true" json:"is_active" default:"true"`
	ListID             string          `gorm:"type:uuid;not null" json:"-"`
	List               *List           `gorm:"foreignKey:ListID" json:"list,omitempty"`
	Comments           []*Comment      `gorm:"foreignKey:CardID" json:"comments"`
	Labels             []*Label        `gorm:"many2many:card_labels" json:"labels"`
	Activities         []*CardActivity `gorm:"foreignKey:CardID" json:"activities"`
	Attachments        []*FileUpload   `gorm:"many2many:card_attachments" json:"attachments"`
}
