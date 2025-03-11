package auth

import "time"

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
