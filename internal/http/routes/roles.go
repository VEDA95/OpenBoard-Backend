package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-sqlbuilder"
	"slices"
	"time"
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
		Returning("id AS role_identifier", "date_created AS role_date_created", "name AS role_name")

	if len(dataValidator.Permissions) > 0 {
		rows := make([]map[string]interface{}, 0)
		rolesPermissionsQuery := sqlbuilder.InsertInto("open_board_role_permissions").Cols("role_id", "permission_id")

		for _, permission := range dataValidator.Permissions {
			rolesPermissionsQuery.Values("inserted_role.role_identifier", permission)
		}

		roleQuery := sqlbuilder.With(
			sqlbuilder.CTETable("inserted_role", "role_identifier", "role_date_created", "role_name").As(insertRoleQuery),
			sqlbuilder.CTEQuery("inserted_role_permissions").As(rolesPermissionsQuery),
		)
		roleQuery.
			Select(auth.RolePermissionsQueryColumns...).
			From("open_board_role_permissions").
			Join("open_board_role", "open_board_role_permissions.role_id = open_board_role.id").
			Join("open_board_role_permission", "open_board_role_permissions.permission_id = open_board_role_permission.id").
			Where("open_board_role_permissions.role_id = inserted_role.role_identifier")

		if err := db.Instance.Many(roleQuery, &rows); err != nil {
			return err
		}

		for _, row := range rows {
			if len(output.Id) == 0 {
				output = auth.Role{
					Id:          row["role_identifier"].(string),
					DateCreated: row["role_date_created"].(time.Time),
					Name:        row["role_name"].(string),
				}
			}

			output.Permissions = append(output.Permissions, &auth.RolePermission{
				Id:          row["permission_identifier"].(string),
				DateCreated: row["permission_date_created"].(time.Time),
				Path:        row["permission_path"].(string),
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
		roleQuery.Set(roleQuery.Assign("name", dataValidator.Name))

		if err := transaction.Exec(roleQuery); err != nil {
			return err
		}

		role.Name = *dataValidator.Name
	}

	if dataValidator.Permissions != nil {
		permissionsToAdd := make([]string, len(role.Permissions))
		permissionsToRemove := make([]string, len(role.Permissions))

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

		cteTables := make([]*sqlbuilder.CTEQueryBuilder, 2)

		if len(permissionsToAdd) > 0 {
			addPermissionQuery := sqlbuilder.InsertInto("open_board_role_permissions").Cols("role_id", "permission_id")

			for _, permission := range permissionsToAdd {
				addPermissionQuery.Values(role.Id, permission)
			}

			cteTables = append(cteTables, sqlbuilder.CTEQuery("add_permissions").As(addPermissionQuery))
		}

		if len(permissionsToRemove) > 0 {
			removePermissionQuery := sqlbuilder.DeleteFrom("open_board_role_permissions")
			removePermissionQuery.Where(removePermissionQuery.In("permission_id", permissionsToRemove))
			cteTables = append(cteTables, sqlbuilder.CTEQuery("remove_permissions").As(removePermissionQuery))
		}

		if len(cteTables) > 0 {
			permissions := make([]*auth.RolePermission, 0)
			permissionUpdateQuery := sqlbuilder.With(cteTables...)
			permissionUpdateQuery.
				Select(
					"open_board_role_permission.id AS permission_identifier",
					"open_board_role_permission.date_created AS permission_date_created",
					"open_board_role_permission.path AS permission_path",
				).
				From("open_board_role_permissions").
				Join("open_board_role_permission", "open_board_role_permissions.permission_id = permission_identifier").
				Where(sqlbuilder.NewCond().Equal("open_board_role_permissions.role_id", role.Id))

			if err := transaction.Many(permissionUpdateQuery, &permissions); err != nil {
				return err
			}

			role.Permissions = permissions
		}
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
