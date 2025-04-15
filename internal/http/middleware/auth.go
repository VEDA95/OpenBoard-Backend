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
	sessionQuery := auth.GetSessionQuery()
	sessionQuery.Where(sessionQuery.Equal("access_token", authToken))

	if err := db.Instance.One(sessionQuery, &session); err != nil {
		return err
	}

	if len(session.Id) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "token not found")
	}

	now := time.Now()

	log.Logger.Debug().Time("now", now).Msg("current time")
	log.Logger.Debug().
		Time("session_expires_on", session.ExpiresOn).
		Bool("is_after_expiration", now.After(session.ExpiresOn)).
		Msg("session expiration date")

	if session.RefreshExpiresOn != nil {
		log.Logger.Debug().
			Time("session_refresh_expires_on", *session.RefreshExpiresOn).
			Bool("is_after_refresh_expiration", now.After(*session.RefreshExpiresOn)).
			Msg("refresh session expiration date")
	}

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

	rows := make([]map[string]interface{}, 0)
	usersRolesQuery := auth.UsersRolesQuery()
	usersRolesQuery.Where(usersRolesQuery.Equal("open_board_user_roles.user_id", session.User.Id))

	if err := db.Instance.Many(usersRolesQuery, &rows); err != nil {
		return err
	}

	auth.AppendRolesToUser(rows, session.User)

	context.Locals("auth_session", session)

	return context.Next()
}
