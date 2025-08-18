package models

import (
	"VEDA95/open_board/api/internal/auth"
	"time"
)

type User struct {
	Base
	LastLogin           *time.Time            `json:"last_login"`
	Username            string                `gorm:"type:varchar(255);not null;unique;index" json:"username"`
	Email               string                `gorm:"type:varchar(255);not null;unique;index" json:"email"`
	FirstName           *string               `gorm:"type:varchar(255)" json:"first_name"`
	LastName            *string               `gorm:"type:varchar(255)" json:"last_name"`
	HashedPassword      string                `gorm:"type:text;not null" json:"-"`
	Enabled             bool                  `gorm:"not null" json:"enabled" default:"true"`
	EmailVerified       bool                  `gorm:"not null" json:"email_verified" default:"false"`
	Roles               []*Role               `gorm:"many2many:user_roles" json:"roles"`
	PasswordResetTokens []*PasswordResetToken `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Sessions            []*Session            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}

func (user *User) HashPassword(password string) error {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil
	}

	user.HashedPassword = hashedPassword

	return nil
}

func (user *User) ValidPassword(password string) bool {
	return auth.CheckPasswordHash(password, user.HashedPassword)
}
