package routes

import (
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service   *service.UserService
	validator *validators.Validator
}

func NewUserHandler(userService *service.UserService, validator *validators.Validator) *UserHandler {
	return &UserHandler{
		service:   userService,
		validator: validator,
	}
}

// UsersGET godoc
//
//	@Description 	List all users
//	@Summary 		List users
//	@Tags			users
//	@Success		200 {object} responses.OkCollectionResponse[auth.User]
//	@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/users [get]
//	@Produce		json
func (userHandler *UserHandler) GET(context *fiber.Ctx) error {
	users, err := userHandler.service.GetUsers()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, users))
}

// UsersPOST godoc
//
//	@Description 	Create user
//	@Summary 		Create user
//	@Tags			users
//	@Param			request body validators.CreateUserValidator false "request data"
//	@Success		201 {object} responses.SuccessResponse[auth.User]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/users [post]
//	@Accept			json
//	@Produce		json
func (userHandler *UserHandler) POST(context *fiber.Ctx) error {
	createValidator := new(validators.CreateUserValidator)

	if err := context.BodyParser(createValidator); err != nil {
		return err
	}

	if errs := userHandler.validator.Validate(createValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	user, err := userHandler.service.CreateUser(createValidator)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.OKResponse(fiber.StatusCreated, fiber.Map{
			"message": fmt.Sprintf("The user: %s has been successfully created", user.Username),
			"user":    user,
		}),
	)
}

// UserGET godoc
//
//	@Description 	List user by ID
//	@Summary 		List user
//	@Tags			users
//	@Param			id path string true "user ID"
//	@Success		200 {object} responses.OkResponse[auth.User]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/users/{id} [get]
//	@Produce		json
func (userHandler *UserHandler) GETByID(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := userHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	user, err := userHandler.service.GetUser(paramValidator.Id)
	if err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, user))
}

// UserPATCH godoc
//
//	@Description 	Update user by ID
//	@Summary 		Update user
//	@Tags			users
//	@Param			id path string true "user ID"
//	@Param			request body responses.SuccessResponse[auth.User] false "request data"
//	@Success		200 {object} responses.OkCollectionResponse[auth.User]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/users/{id} [patch]
//	@Accept			json
//	@Produce		json
func (userHandler *UserHandler) PATCH(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := userHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	user, err := userHandler.service.GetUser(paramValidator.Id)
	if err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	updateUserValidator := new(validators.UpdateUserValidator)

	if err := context.BodyParser(updateUserValidator); err != nil {
		return err
	}

	if errs := userHandler.validator.Validate(updateUserValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	updatedUser, err := userHandler.service.UpdateUser(user.ID, updateUserValidator)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, fiber.Map{
			"message": fmt.Sprintf("user: %s has been successfully updated", updatedUser.Username),
			"user":    updatedUser,
		}),
	)
}

// UserDELETE godoc
//
//	@Description 	Delete user by ID
//	@Summary 		Delete user
//	@Tags			users
//	@Param			id path string true "user ID"
//	@Success		200 {object} responses.OkCollectionResponse[auth.User]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		404,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/api/users/{id} [delete]
//	@Produce		json
func (userHandler *UserHandler) DELETE(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := userHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	user, err := userHandler.service.GetUser(paramValidator.Id)
	if err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	if err := userHandler.service.DeleteUser(user.ID); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: fmt.Sprintf("user: %s has been successfully deleted", user.Username)},
		),
	)
}
