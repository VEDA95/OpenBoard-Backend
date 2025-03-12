package auth

import "time"

type UserSession struct {
	Id               string                 `db:"session_id"`
	DateCreated      time.Time              `db:"session_date_created"`
	DateUpdated      *time.Time             `db:"session_date_updated,omitempty"`
	ExpiresOn        time.Time              `db:"session_expires_on"`
	RefreshExpiresOn *time.Time             `db:"session_refresh_expires_on,omitempty"`
	User             *User                  `db:""`
	SessionType      string                 `db:"user_session_type"`
	AccessToken      string                 `db:"session_access_token"`
	RefreshToken     *string                `db:"session_refresh_token,omitempty"`
	IPAddress        string                 `db:"session_ip_address"`
	UserAgent        string                 `db:"session_user_agent"`
	AdditionalInfo   map[string]interface{} `db:"session_additional_info"`
}
