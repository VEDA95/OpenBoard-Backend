package middleware

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
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

	session := new(auth.UserSession)
	authToken := authHeaderSplit[1]
	sessionQuery := sqlbuilder.Select(
		"id",
		"date_created",
		"date_updated",
		"expires_on",
		"refresh_expires_on",
		"session_type",
		"access_token",
		"refresh_token",
		"ip_address",
		"user_agent",
		"additional_info",
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

	if err := db.Instance.One(sessionQuery, session); err != nil {
		return err
	}

	if session == nil {
		return fiber.NewError(fiber.StatusNotFound, "token not found")
	}

	if session.ExpiresOn.After(time.Now().Local()) {
		if session.RefreshExpiresOn.After(time.Now().Local()) {
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
