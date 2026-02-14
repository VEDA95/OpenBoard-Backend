package routes

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var _ = models.Permission{}

type PermissionHandler struct {
	service   *service.PermissionService
	validator *validators.Validator
}

func NewPermissionHandler(permissionService *service.PermissionService, validator *validators.Validator) *PermissionHandler {
	return &PermissionHandler{
		service:   permissionService,
		validator: validator,
	}
}

// PermissionsGET godoc
//
//	@Description 	List all permissions
//	@Summary 		List permissions
//	@Tags			authorization
//	@Success		200 {object} responses.OkCollectionResponse[models.Permission]
//	@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions [get]
//	@Produce		json
func (permissionHandler *PermissionHandler) GET(context *fiber.Ctx) error {
	permissions, err := permissionHandler.service.GetPermissions()
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
//	@Success		201 {object} responses.SuccessResponse[models.Permission]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions [post]
//	@Accept			json
//	@Produce		json
func (permissionHandler *PermissionHandler) POST(context *fiber.Ctx) error {
	dataValidator := new(validators.CreatePermissionValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := permissionHandler.validator.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	permission, err := permissionHandler.service.CreatePermission(dataValidator)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(fiber.StatusOK, fmt.Sprintf("permission: %s was created suucessfully", permission.ID), permission),
	)
}

// PermissionGET godoc
//
//	@Description 	List permission by ID
//	@Summary 		List permission
//	@Tags			authorization
//	@Param			id path string true "Permission ID"
//	@Success		200 {object} responses.OkResponse[models.Permission]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions/{id} [get]
//	@Produce		json
func (permissionHandler *PermissionHandler) GETByID(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := permissionHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	permission, err := permissionHandler.service.GetPermission(paramValidator.Id)
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
//	@Success		200 {object} responses.SuccessResponse[models.Permission]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/permissions/{id} [patch]
//	@Accept			json
//	@Produce		json
func (permissionHandler *PermissionHandler) PATCH(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := permissionHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	permission, err := permissionHandler.service.GetPermission(paramValidator.Id)
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

	if errs := permissionHandler.validator.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if dataValidator.Path == nil {
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, permission))
	}

	updatedPermission, err := permissionHandler.service.UpdatePermission(paramValidator.Id, dataValidator)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			fmt.Sprintf("permission: %s was updated sucessfully", updatedPermission.ID),
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
func (permissionHandler *PermissionHandler) DELETE(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := permissionHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	permission, err := permissionHandler.service.GetPermission(paramValidator.Id)
	if err != nil {
		return err
	}

	if permission == nil {
		return fiber.NewError(fiber.StatusNotFound, "permission not found")
	}

	if err := permissionHandler.service.DeletePermission(permission.ID); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
			Message: fmt.Sprintf("permission: %s was deleted sucessfully", permission.ID),
		}),
	)
}
