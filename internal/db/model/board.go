package models

type Board struct {
	Base
	Name        string        `gorm:"type:varchar(255);not null" json:"name"`
	IsPublic    bool          `gorm:"not null;default:true" json:"is_public"`
	UserID      string        `gorm:"type:uuid;not null" json:"-"`
	WorkspaceID string        `gorm:"type:uuid;not null" json:"-"`
	User        *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Workspace   *Worksapce    `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
	Permissions []*Permission `gorm:"many2many:board_permissions" json:"permissions"`
	Lists       []*List       `gorm:"foreignKey:BoardID" json:"lists"`
}
