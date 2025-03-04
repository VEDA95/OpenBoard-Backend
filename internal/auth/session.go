package auth

import "time"

type UserSession struct {
	Id             string                 `db:"id"`
	DateCreated    time.Time              `db:"date_created"`
	DateUpdated    *time.Time             `db:"date_updated,omitempty"`
	ExpiresOn      time.Time              `db:"expires_on"`
	User           *User                  `db:"open_board_user"`
	SessionType    string                 `db:"session_type"`
	Remember       bool                   `db:"remember_me"`
	AccessToken    string                 `db:"access_token"`
	RefreshToken   *string                `db:"refresh_token,omitempty"`
	IPAddress      string                 `db:"ip_address"`
	UserAgent      string                 `db:"user_agent"`
	AdditionalInfo map[string]interface{} `db:"additional_info"`
}
