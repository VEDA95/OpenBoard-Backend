package auth

import (
	"VEDA95/open_board/api/internal/db"
	"errors"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

type User struct {
	Id             string     `json:"id" db:"id"`
	DateCreated    time.Time  `json:"date_created" db:"date_created"`
	DateUpdated    *time.Time `json:"date_updated" db:"date_updated,omitempty"`
	LastLogin      *time.Time `json:"last_login" db:"last_login,omitempty"`
	Username       string     `json:"username" db:"username"`
	Email          string     `json:"email_address" db:"email"`
	FirstName      *string    `json:"first_name" db:"first_name,omitempty"`
	LastName       *string    `json:"last_name" db:"last_name,omitempty"`
	Enabled        bool       `json:"enabled" db:"enabled"`
	EmailVerified  bool       `json:"email_verified" db:"email_verified"`
	HashedPassword string     `json:"-" db:"hashed_password"`
}

var UserQueryColumns = []string{
	"id",
	"date_created",
	"date_updated",
	"last_login",
	"username",
	"email",
	"first_name",
	"last_name",
	"hashed_password",
	"enabled",
	"email_verified",
}

func GetUsers() ([]User, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	output := make([]User, 0)
	usersQuery := sqlbuilder.Select(UserQueryColumns...).From("open_board_user")

	if err := db.Instance.Many(usersQuery, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func GetUser(id string) (*User, error) {
	if db.Instance == nil {
		return nil, errors.New("database not initialized")
	}

	var output User
	userQuery := sqlbuilder.Select(UserQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("id", id))

	if err := db.Instance.One(userQuery, &output); err != nil {
		return nil, err
	}
	
	return &output, nil
}
