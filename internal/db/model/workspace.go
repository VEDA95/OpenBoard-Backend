package models

type Worksapce struct {
	Base
	Name        string        `gorm:"varchar(255);not null" json:"name"`
	Description *string       `gorm:"text" json:"description"`
	IsPublic    bool          `gorm:"not null;default:true" json:"is_public" default:"true"`
	UserID      string        `gorm:"type:uuid;not null" json:"-"`
	User        *User         `gorm:"foreignKey:UserID" json:"user"`
	Permissions []*Permission `gorm:"many2many:workspace_permissions" json:"permissions"`
	Boards      []*Board      `gorm:"foreignKey:WorkspaceID" json:"boards"`
}
