package routes

import (
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	service   *service.RoleService
	validator *validators.Validator
}

func NewRoleHandler(roleService *service.RoleService, validator *validators.Validator) *RoleHandler {
	return &RoleHandler{
		service:   roleService,
		validator: validator,
	}
}

// RolesGET godoc
//
//	@Description	Get all roles
//	@Summary 		List roles
//	@Tags 			authorization
//	@Success		200	{object} responses.OkCollectionResponse[auth.Role]
//	@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/roles [get]
//	@Produce		json
func (roleHandler *RoleHandler) GET(context *fiber.Ctx) error {
	roles, err := roleHandler.service.GetRoles()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, roles))
}

// RolesPOST godoc
//
//		@Description	Create Role
//		@Summary 		Create role
//		@Tags 			authorization
//		@param			request body validators.CreateRoleValidator false "Request Data"
//	 	@Success		201	{object} responses.SuccessResponse[auth.Role]
//		@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//		@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Router			/api/roles [post]
//		@Accept			json
//		@Produce		json
func (roleHandler *RoleHandler) POST(context *fiber.Ctx) error {
	dataValidator := new(validators.CreateRoleValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := roleHandler.validator.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	role, err := roleHandler.service.CreateRole(dataValidator)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.CreateSuccessResponse(
			fiber.StatusCreated,
			fmt.Sprintf("Role: %s has be successfully created", role.ID),
			role,
		),
	)
}

// RoleGET godoc
//
//	@Description	Get role by ID
//	@Summary 		List role
//	@Tags 			authorization
//	@Param			id path string true "Role ID"
//	@Success		200	{object} responses.OkResponse[auth.Role]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/roles/{id} [get]
//	@Produce		json
func (roleHandler *RoleHandler) GETByID(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := roleHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	role, err := roleHandler.service.GetRole(paramValidator.Id)
	if err != nil {
		return err
	}

	if role == nil {
		return fiber.NewError(fiber.StatusNotFound, "role not found")
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, role))
}

// RolePATCH godoc
//
//	@Description	Update role by ID
//	@Summary 		Update role
//	@Tags 			authorization
//	@Param			id path string true "Role ID"
//	@Param			request body validators.UpdateRoleValidator false "Request Data"
//	@Success		200	{object} responses.SuccessResponse[auth.Role]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/roles/{id} [patch]
//	@Accept			json
//	@Produce		json
func (roleHandler *RoleHandler) PATCH(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := roleHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	role, err := roleHandler.service.GetRole(paramValidator.Id)
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

	if errs := roleHandler.validator.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	updatedRole, err := roleHandler.service.UpdateRole(paramValidator.Id, dataValidator)

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(fiber.StatusOK, fmt.Sprintf("role: %s has been updated successfully", updatedRole.ID), updatedRole.ID),
	)
}

// RoleDELETE godoc
//
//	@Description	Delete Role by ID
//	@Summary 		Delete role
//	@Tags 			authorization
//	@Param			id path string true "Role ID"
//	@Success		200	{object} responses.OkResponse[responses.GenericMessage]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/roles/{id} [delete]
//	@Accept			json
//	@Produce		json
func (roleHandler *RoleHandler) DELETE(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := roleHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	role, err := roleHandler.service.GetRole(paramValidator.Id)
	if err != nil {
		return err
	}

	if role == nil {
		return fiber.NewError(fiber.StatusNotFound, "role not found")
	}

	if err := roleHandler.service.DeleteRole(role.ID); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: fmt.Sprintf("role: %s has been deleted successfully", role.ID)},
		),
	)
}
