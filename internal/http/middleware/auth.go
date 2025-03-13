package middleware

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/log"
	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-sqlbuilder"
	"strings"
	"time"
)

func CheckUserAuthentication(context *fiber.Ctx) error {
	authHeader := context.Get("Authorization")
	isCookie := false

	if len(authHeader) == 0 {
		authCookie := context.Cookies("open_board_auth_session", "")

		if len(authCookie) == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		authHeader = authCookie
		isCookie = true
	}

	authHeaderSplit := strings.Split(authHeader, " ")

	if len(authHeaderSplit) < 2 {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var session auth.UserSession
	authToken := authHeaderSplit[1]
	sessionQuery := sqlbuilder.Select(
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
	).From("open_board_user_session")
	sessionQuery.
		Where(sessionQuery.Equal("access_token", authToken)).
		Join("open_board_user", "open_board_user_session.user_id = open_board_user.id")

	if err := db.Instance.One(sessionQuery, &session); err != nil {
		return err
	}

	if len(session.Id) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "token not found")
	}

	now := time.Now()

	log.Logger.Debug().Time("now", now).Msg("current time")
	log.Logger.Debug().Time("session_expires_on", session.ExpiresOn).Time("session_refresh_expires_on", *session.RefreshExpiresOn).Msg("session expiration date")
	log.Logger.Debug().Bool("is_after_expiration", now.After(session.ExpiresOn)).Bool("is_after_refresh_expiration", now.After(*session.RefreshExpiresOn)).Msg("is current time past the expiration date")

	if now.After(session.ExpiresOn) {
		if session.RefreshExpiresOn != nil && now.After(*session.RefreshExpiresOn) {
			deleteSessionQuery := sqlbuilder.DeleteFrom("open_board_user_session")
			deleteSessionQuery.Where(deleteSessionQuery.Equal("id", session.Id))

			if err := db.Instance.Exec(deleteSessionQuery); err != nil {
				return err
			}

			if isCookie {
				context.ClearCookie("open_board_session_remember_me")
			}
		}

		if isCookie {
			context.ClearCookie("open_board_auth_session")
		}

		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	context.Locals("auth_session", session)

	return context.Next()
}
