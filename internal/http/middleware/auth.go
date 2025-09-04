package middleware

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	service *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		service: authService,
	}
}

func (authMiddleware *AuthMiddleware) RequireAuthentication() fiber.Handler {
	return func(context *fiber.Ctx) error {
		var token string
		authHeader := context.Get("Authorization")
		isCookie := false

		if len(authHeader) > 0 {
			authSplit := strings.Split(authHeader, " ")
			if len(authSplit) == 2 && authSplit[0] == "Bearer" {
				token = authSplit[1]
			}
		} else {
			token = context.Cookies("open_board_auth_session", "")
		}

		if len(token) == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		session, err := authMiddleware.service.ValidateSession(token)
		if err != nil {
			errString := err.Error()
			isInvalidCredentials := errString == "invalid credentials"
			isRefreshRequired := errString == "refresh required"

			if isInvalidCredentials && isCookie {
				context.ClearCookie("open_board_session_remember_me")
				context.ClearCookie("open_board_auth_session")
			}

			if isRefreshRequired && isCookie {
				context.ClearCookie("open_board_auth_session")
			}

			if isInvalidCredentials || isRefreshRequired {
				return fiber.NewError(fiber.StatusUnauthorized, errString)
			}

			return err
		}

		context.Locals("auth_session", *session)

		return context.Next()
	}
}

func (*AuthMiddleware) RequireAuthorization(permissions ...string) fiber.Handler {
	return func(context *fiber.Ctx) error {
		if len(permissions) == 0 {
			return context.Next()
		}

		session, ok := context.Locals("auth_session").(models.Session)

		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		permissionMap := make(map[string]bool)

		for _, role := range session.User.Roles {
			for _, permission := range role.Permissions {
				permissionMap[permission.Path] = true
			}
		}

		if permissionMap["auth:superuser"] {
			return context.Next()
		}

		for _, permission := range permissions {
			if !permissionMap[permission] {
				return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
			}
		}

		return context.Next()
	}
}
