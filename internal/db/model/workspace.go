package models

type Workspace struct {
	Base
	Name        string        `gorm:"varchar(255);not null" json:"name"`
	Description *string       `gorm:"text" json:"description"`
	IsPublic    bool          `gorm:"not null" json:"is_public"`
	UserID      string        `gorm:"type:uuid;not null" json:"-"`
	User        *User         `gorm:"foreignKey:UserID" json:"user"`
	Permissions []*Permission `gorm:"many2many:workspace_permissions" json:"permissions"`
	Boards      []*Board      `gorm:"foreignKey:WorkspaceID" json:"boards"`
}
