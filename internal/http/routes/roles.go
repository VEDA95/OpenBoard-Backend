package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/huandu/go-sqlbuilder"
	"slices"
)

func RolesGET(context *fiber.Ctx) error {
	roles, err := auth.GetRoles()

	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, roles))
}

func RolesPOST(context *fiber.Ctx) error {
	dataValidator := new(validators.CreateRoleValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	var output auth.Role
	insertRoleQuery := sqlbuilder.InsertInto("open_board_role").
		Cols("name").
		Values(dataValidator.Name).
		Returning("id AS role_identifier", "name AS role_name")

	if len(dataValidator.Permissions) > 0 {
		transaction, err := db.Instance.Begin()

		if err != nil {
			return err
		}

		var role auth.Role

		if err := transaction.One(insertRoleQuery, &role); err != nil {
			return err
		}

		rows := make([]map[string]interface{}, 0)
		rolesPermissionsQuery := sqlbuilder.InsertInto("open_board_role_permissions").Cols("role_id", "permission_id")

		for _, permission := range dataValidator.Permissions {
			rolesPermissionsQuery.Values(role.Id, permission)
		}

		roleQuery := sqlbuilder.Select(auth.RolePermissionsQueryColumns...).From("open_board_role_permissions")
		roleQuery.
			JoinWithOption(sqlbuilder.LeftJoin, "open_board_role", "open_board_role_permissions.role_id = open_board_role.id").
			JoinWithOption(sqlbuilder.LeftJoin, "open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
			Where(roleQuery.Equal("open_board_role_permissions.role_id", role.Id))

		if err := transaction.Exec(rolesPermissionsQuery); err != nil {
			return err
		}

		if err := transaction.Many(roleQuery, &rows); err != nil {
			return err
		}

		if err := transaction.Commit(); err != nil {
			return err
		}

		for _, row := range rows {
			if len(output.Id) == 0 {
				output = auth.Role{
					Id:   row["role_identifier"].(uuid.UUID).String(),
					Name: row["role_name"].(string),
				}
			}

			output.Permissions = append(output.Permissions, &auth.RolePermission{
				Id:   row["permission_identifier"].(uuid.UUID).String(),
				Path: row["permission_path"].(string),
			})
		}

	} else {
		roleQuery := insertRoleQuery

		if err := db.Instance.One(roleQuery, &output); err != nil {
			return err
		}
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.OKResponse(fiber.StatusCreated, fiber.Map{
			"message": fmt.Sprintf("Role: %s has be successfully created", output.Id),
			"role":    output,
		}),
	)
}

func RoleGET(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	role, err := auth.GetRole(paramValidator.Id)

	if err != nil {
		return err
	}

	if role == nil {
		return fiber.NewError(fiber.StatusNotFound, "role not found")
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, role))
}

func RolePATCH(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	role, err := auth.GetRole(paramValidator.Id)

	if err != nil {
		return err
	}

	if role == nil {
		return fiber.NewError(fiber.StatusNotFound, "role not found")
	}

	dataValidator := new(validators.UpdateRoleValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if dataValidator.Name == nil && dataValidator.Permissions == nil {
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, role))
	}

	transaction, err := db.Instance.Begin()

	if err != nil {
		return err
	}

	if dataValidator.Name != nil && *dataValidator.Name != role.Name {
		roleQuery := sqlbuilder.Update("open_board_role")
		roleQuery.
			Where(roleQuery.Equal("open_board_role.id", role.Id)).
			Set(roleQuery.Assign("name", dataValidator.Name))

		if err := transaction.Exec(roleQuery); err != nil {
			return err
		}

		role.Name = *dataValidator.Name
	}

	if dataValidator.Permissions != nil {
		permissionsToAdd := make([]interface{}, 0)
		permissionsToRemove := make([]interface{}, 0)

		for _, permission := range *dataValidator.Permissions {
			match := slices.ContainsFunc(role.Permissions, func(rolePermission *auth.RolePermission) bool {
				return permission == rolePermission.Id
			})

			if !match {
				permissionsToAdd = append(permissionsToAdd, permission)
			}
		}

		for _, permission := range role.Permissions {
			match := slices.ContainsFunc(*dataValidator.Permissions, func(rolePermission string) bool {
				return permission.Id == rolePermission
			})

			if !match {
				permissionsToRemove = append(permissionsToRemove, permission.Id)
			}
		}

		if len(permissionsToAdd) > 0 {
			addPermissionQuery := sqlbuilder.InsertInto("open_board_role_permissions").Cols("role_id", "permission_id")

			for _, permission := range permissionsToAdd {
				addPermissionQuery.Values(role.Id, permission)
			}

			if err := transaction.Exec(addPermissionQuery); err != nil {
				return err
			}
		}

		if len(permissionsToRemove) > 0 {
			removePermissionQuery := sqlbuilder.DeleteFrom("open_board_role_permissions")
			removePermissionQuery.Where(removePermissionQuery.In("permission_id", permissionsToRemove...))

			if err := transaction.Exec(removePermissionQuery); err != nil {
				return err
			}
		}

		permissions := make([]*auth.RolePermission, 0)
		permissionsQuery := sqlbuilder.Select(
			"open_board_role_permission.id AS permission_identifier",
			"open_board_role_permission.path AS permission_path",
		).From("open_board_role_permissions")
		permissionsQuery.
			Join("open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
			Where(permissionsQuery.Equal("open_board_role_permissions.role_id", role.Id))

		if err := transaction.Many(permissionsQuery, &permissions); err != nil {
			return err
		}

		role.Permissions = permissions
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, fiber.Map{
			"message": fmt.Sprintf("role: %s has been updated successfully", role.Id),
			"role":    role,
		}),
	)
}

func RoleDELETE(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	role, err := auth.GetRole(paramValidator.Id)

	if err != nil {
		return err
	}

	if role == nil {
		return fiber.NewError(fiber.StatusNotFound, "role not found")
	}

	deleteRoleQuery := sqlbuilder.DeleteFrom("open_board_role")
	deleteRoleQuery.Where(deleteRoleQuery.Equal("id", role.Id))

	if err := db.Instance.Exec(deleteRoleQuery); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: fmt.Sprintf("role: %s has been deleted successfully", role.Id)},
		),
	)
}
