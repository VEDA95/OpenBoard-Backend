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
)

// PermissionsGET godoc
//
//	@Description 	List all permissions
//	@Summary 		List permissions
//	@Tags			authorization
//	@Success		200 {object} responses.OkCollectionResponse[auth.Permission]
//	@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions [get]
//	@Produce		json
func PermissionsGET(context *fiber.Ctx) error {
	permissions, err := auth.GetPermissions()

	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, permissions))
}

// PermissionsPOST godoc
//
//	@Description 	Create permission
//	@Summary 		Create permission
//	@Tags			authorization
//	@Param			request body validators.CreatePermissionValidator false "Request Data"
//	@Success		201 {object} responses.SuccessResponse[auth.Permission]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions [post]
//	@Accept			json
//	@Produce		json
func PermissionsPOST(context *fiber.Ctx) error {
	dataValidator := new(validators.CreatePermissionValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	var output auth.Permission
	createPermissionQuery := sqlbuilder.
		InsertInto("open_board_role_permission").
		Cols("path").
		Values(dataValidator.Path).
		Returning("*")

	if err := db.Instance.One(createPermissionQuery, &output); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(fiber.StatusOK, fmt.Sprintf("permission: %s was created suucessfully", output.Id), output),
	)
}

// PermissionGET godoc
//
//	@Description 	List permission by ID
//	@Summary 		List permission
//	@Tags			authorization
//	@Param			id path string true "Permission ID"
//	@Success		200 {object} responses.OkResponse[auth.Permission]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions/{id} [get]
//	@Produce		json
func PermissionGET(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	permission, err := auth.GetPermission(paramValidator.Id)

	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, permission))
}

// PermissionPATCH godoc
//
//	@Description 	Update permission by ID
//	@Summary 		Update permission
//	@Tags			authorization
//	@Param			id path string true "Permission ID"
//	@Param			request body validators.UpdatePermissionValidator false "Request Data"
//	@Success		200 {object} responses.SuccessResponse[auth.Permission]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions/{id} [patch]
//	@Accept			json
//	@Produce		json
func PermissionPATCH(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	permission, err := auth.GetPermission(paramValidator.Id)

	if err != nil {
		return err
	}

	if permission == nil {
		return fiber.NewError(fiber.StatusNotFound, "permission not found")
	}

	dataValidator := new(validators.UpdatePermissionValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if dataValidator.Path == nil {
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, permission))
	}

	var updatedPermission auth.Permission
	updatePermissionQuery := sqlbuilder.Update("open_board_role_permission")
	updatePermissionQuery.
		Where(updatePermissionQuery.Equal("id", permission.Id)).
		Set(updatePermissionQuery.Assign("path", dataValidator.Path))
	withPermissionQuery := sqlbuilder.With(
		sqlbuilder.CTEQuery("role_permission_update").As(updatePermissionQuery),
	).Select("*").From("open_board_role_permission")

	if err := db.Instance.One(withPermissionQuery, &updatedPermission); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			fmt.Sprintf("permission: %s was updated sucessfully", updatedPermission.Id),
			updatedPermission,
		),
	)
}

// PermissionDELETE godoc
//
//	@Description 	Delete permission by ID
//	@Summary 		Delete permission
//	@Tags			authorization
//	@Param			id path string true "Permission ID"
//	@Success		200 {object} responses.OkResponse[responses.GenericMessage]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions/{id} [delete]
//	@Produce		json
func PermissionDELETE(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	permission, err := auth.GetPermission(paramValidator.Id)

	if err != nil {
		return err
	}

	if permission == nil {
		return fiber.NewError(fiber.StatusNotFound, "permission not found")
	}

	deletePermissionQuery := sqlbuilder.DeleteFrom("open_board_role_permission")
	deletePermissionQuery.Where(deletePermissionQuery.Equal("id", permission.Id))

	if err := db.Instance.Exec(deletePermissionQuery); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
			Message: fmt.Sprintf("permission: %s was deleted sucessfully", permission.Id),
		}),
	)
}
