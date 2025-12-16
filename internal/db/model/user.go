package models

import (
	"time"

	"github.com/alexedwards/argon2id"
)

type User struct {
	Base
	LastLogin      *time.Time    `json:"last_login"`
	Username       string        `gorm:"type:varchar(255);not null;unique;index" json:"username"`
	Email          string        `gorm:"type:varchar(255);not null;unique;index" json:"email"`
	FirstName      *string       `gorm:"type:varchar(255)" json:"first_name"`
	LastName       *string       `gorm:"type:varchar(255)" json:"last_name"`
	HashedPassword string        `gorm:"type:text;not null" json:"-"`
	ThumbnailID    *string       `gorm:"type:uuid;" json:"-"`
	Enabled        bool          `gorm:"not null;default:true" json:"enabled" default:"true"`
	EmailVerified  bool          `gorm:"not null;default:false" json:"email_verified" default:"false"`
	Thumbnail      *FileUpload   `gorm:"foreignKey:ThumbnailID" json:"thumbnail"`
	Roles          []*Role       `gorm:"many2many:user_roles" json:"roles"`
	Sessions       []*Session    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Files          []*FileUpload `gorm:"foreignKey:UserID" json:"files,omitempty"`
}

func (user *User) HashPassword(password string) error {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	user.HashedPassword = hashedPassword

	return nil
}

func (user *User) ValidPassword(password string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, user.HashedPassword)
	if err != nil {
		return false
	}

	return match
}

func (user *User) permissionPaths() map[string]bool {
	permissionMap := make(map[string]bool)

	if len(user.Roles) == 0 {
		return permissionMap
	}

	for _, role := range user.Roles {
		if len(role.Permissions) == 0 {
			continue
		}

		for _, permission := range role.Permissions {
			permissionMap[permission.Path] = true
		}
	}

	return permissionMap
}

func (user *User) IsAuthorized(paths ...string) bool {
	if user.IsSuperuser() {
		return true
	}

	permissions := user.permissionPaths()

	if (len(permissions) == 0 && len(paths) == 0) || (len(permissions) > 0 && len(paths) == 0) {
		return true
	}

	if len(permissions) == 0 && len(paths) > 0 {
		return false
	}

	for _, permission := range paths {
		if !permissions[permission] {
			return false
		}
	}

	return true
}

func (user *User) IsAuthorizedPartial(paths ...string) bool {
	if user.IsSuperuser() {
		return true
	}

	permissions := user.permissionPaths()

	if (len(permissions) == 0 && len(paths) == 0) || (len(permissions) > 0 && len(paths) == 0) {
		return true
	}

	if len(permissions) == 0 && len(paths) > 0 {
		return false
	}

	partial := false

	for _, permission := range paths {
		if permissions[permission] {
			partial = true
			break
		}
	}

	return partial
}

func (user *User) IsSuperuser() bool {
	permissions := user.permissionPaths()

	return permissions["auth:superuser"]
}
