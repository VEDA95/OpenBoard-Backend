package models

import (
	"time"

	"gorm.io/gorm"
)

type GeneralSettings struct {
	ID                     int         `gorm:"primaryKey;check:check_single_row_general,id = 1" json:"-"`
	UpdatedAt              *time.Time  `json:"updated_at"`
	AppName                string      `gorm:"type:varchar(100);not null;default:'Kanban Board'" json:"app_name"`
	AppURL                 string      `gorm:"type:varchar(255);not null;default:'http://localhost:3000'" json:"app_url"`
	AppLogoID              *string     `gorm:"type:uuid" json:"-"`
	AppFaviconID           *string     `gorm:"type:uuid" json:"-"`
	AppDescription         *string     `gorm:"type:text" json:"app_description"`
	ShowAnnouncementBanner bool        `gorm:"default:false" json:"show_announcement_banner"`
	AnnouncementMessage    *string     `gorm:"type:text" json:"announcement_message"`
	AnnouncementType       *string     `gorm:"type:varchar(20);default:'info'" json:"announcement_type"`
	DefaultLanguage        string      `gorm:"type:varchar(10);default:'en'" json:"default_language"`
	DefaultTimezone        string      `gorm:"type:varchar(50);default:'UTC'" json:"default_timezone"`
	DefaultItemsPerPage    int         `gorm:"default:25" json:"default_items_per_page"`
	MaxFileSize            int         `gorm:"default:1024" json:"max_file_size"`
	AppLogo                *FileUpload `gorm:"foreignKey:AppLogoID" json:"app_logo"`
	AppFavicon             *FileUpload `gorm:"foreignKey:AppFaviconID" json:"app_favicon"`
}

func (g *GeneralSettings) BeforeCreate(tx *gorm.DB) error {
	g.ID = 1

	return nil
}
