package auth

import (
	"github.com/huandu/go-sqlbuilder"
	"time"
)

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

var SessionQueryColumns = []string{
	"open_board_user_session.id AS session_id",
	"open_board_user_session.date_created AS session_date_created",
	"open_board_user_session.date_updated AS session_date_updated",
	"open_board_user_session.expires_on AS session_expires_on",
	"open_board_user_session.refresh_expires_on AS session_refresh_expires_on",
	"open_board_user_session.session_type AS user_session_type",
	"open_board_user_session.access_token AS session_access_token",
	"open_board_user_session.refresh_token AS session_refresh_token",
	"open_board_user_session.ip_address AS session_ip_address",
	"open_board_user_session.user_agent AS session_user_agent",
	"open_board_user_session.additional_info AS session_additional_info",
	"open_board_user.id",
	"open_board_user.date_created",
	"open_board_user.date_updated",
	"open_board_user.last_login",
	"open_board_user.username",
	"open_board_user.email",
	"open_board_user.first_name",
	"open_board_user.last_name",
	"open_board_user.hashed_password",
	"open_board_user.enabled",
	"open_board_user.email_verified",
}

func GetSessionQuery() *sqlbuilder.SelectBuilder {
	return sqlbuilder.Select(SessionQueryColumns...).
		From("open_board_user_session").
		Join("open_board_user", "open_board_user_session.user_id = open_board_user.id")
}
