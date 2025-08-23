package auth

import (
	models "VEDA95/open_board/api/internal/db/model"
	"time"
)

type LocalUserLogin struct {
	User             *models.User `json:"user"`
	AccessToken      string       `json:"access_token"`
	ExpiresIn        int          `json:"expires_in"`
	RefreshToken     *string      `json:"refresh_token,omitempty" extensions:"x-nullable,x-omitempty"`
	RefreshExpiresIn *int         `json:"refresh_expires_in,omitempty" extensions:"x-nullable,x-omitempty"`
}

type LoginResponse struct {
	Session          *models.Session
	ExpiresIn        int
	RefreshExpiresIn *int
}

type PasswordResetToken struct {
	Id          string    `db:"id"`
	DateCreated time.Time `db:"date_created"`
	Type        string    `db:"type"`
	UserId      string    `db:"user_id"`
	ExpiresOn   time.Time `db:"expires_on"`
	Token       string    `db:"token"`
}

type AuthenticatedPasswordResetResponse struct {
	Token string `json:"token"`
}

type PasswordResetEmailVariables struct {
	Token    string
	Email    string
	Username string
}
