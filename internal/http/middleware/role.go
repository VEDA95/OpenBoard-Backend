package middleware

import (
	"VEDA95/open_board/api/internal/auth"
	"github.com/gofiber/fiber/v2"
	"slices"
)

func CheckUserPermissions(permissions []string) func(ctx *fiber.Ctx) error {
	return func(context *fiber.Ctx) error {
		sessionLocal := context.Locals("auth_session")

		if sessionLocal == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		session := sessionLocal.(auth.UserSession)
		userPermissions := make([]string, 0)

		for _, role := range session.User.Roles {
			for _, permission := range role.Permissions {
				userPermissions = append(userPermissions, permission.Path)
			}
		}

		if len(permissions) == 0 || slices.Contains(userPermissions, "auth:superuser") {
			return context.Next()
		}

		match := false

		for _, permission := range permissions {
			permissionMatch := slices.ContainsFunc(userPermissions, func(userPermission string) bool {
				return permission == userPermission
			})

			if permissionMatch {
				match = true
			}
		}

		if !match {
			return fiber.NewError(fiber.StatusForbidden, "unauthorized")
		}

		return context.Next()
	}
}
