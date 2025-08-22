package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"
	"fmt"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-sqlbuilder"
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

// UserInfoGET godoc
//
//	@Description 	Retrieve own user details. IMPORTANT: The actual endpoint is '/auth/@me' (with @ symbol)
//	@Summary 		Retrieve user details
//	@Tags			authentication
//	@Success		200 {object} responses.OkResponse[auth.User]
//	@Failure		401,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/auth/me [get]
//	@Produce		json
func UserInfoGET(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(auth.UserSession)

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, session.User))
}

// UserInfoPATCH godoc
//
//	@Description 	Update own user details. IMPORTANT: The actual endpoint is '/auth/@me' (with @ symbol)
//	@Summary 		Update user details
//	@Tags			authentication
//	@Param			request body validators.UpdateUserValidator false "request Data"
//	@Success		200 {object} responses.SuccessResponse[auth.User]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		401,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/auth/me [patch]
//	@Accept			json
//	@Produce		json
func UserInfoPATCH(context *fiber.Ctx) error {
	updateUserValidator := new(validators.UpdateUserValidator)

	if err := context.BodyParser(updateUserValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(updateUserValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	now := time.Now().Local()
	session := context.Locals("auth_session").(auth.UserSession)
	updateUserQuery := sqlbuilder.Update("open_board_user")
	updateUserQuery.Where(updateUserQuery.Equal("id", session.User.Id)).Set(updateUserQuery.Assign("date_updated", now))

	if (updateUserValidator.Username != nil && len(*updateUserValidator.Username) > 0) && *updateUserValidator.Username != session.User.Username {
		updateUserQuery.SetMore(updateUserQuery.Assign("username", updateUserValidator.Username))
	}

	if (updateUserValidator.Email != nil && len(*updateUserValidator.Email) > 0) && *updateUserValidator.Email != session.User.Email {
		updateUserQuery.SetMore(updateUserQuery.Assign("email", session.User.Email))
	}

	if updateUserValidator.FirstName != nil && updateUserValidator.FirstName != session.User.FirstName {
		if len(*updateUserValidator.FirstName) == 0 && session.User.FirstName != nil {
			updateUserQuery.SetMore(updateUserQuery.Assign("first_name", nil))
		} else {
			updateUserQuery.SetMore(updateUserQuery.Assign("first_name", updateUserValidator.FirstName))
		}
	}

	if updateUserValidator.LastName != nil && updateUserValidator.LastName != session.User.LastName {
		if len(*updateUserValidator.LastName) == 0 && session.User.LastName != nil {
			updateUserQuery.SetMore(updateUserQuery.Assign("last_name", nil))
		} else {
			updateUserQuery.SetMore(updateUserQuery.Assign("last_name", updateUserValidator.LastName))
		}
	}

	if updateUserValidator.Roles != nil {
		transaction, err := db.Instance.Begin()
		if err != nil {
			return err
		}

		rolesToAdd := make([]interface{}, 0)
		rolesToRemove := make([]interface{}, 0)

		for _, role := range *updateUserValidator.Roles {
			match := slices.ContainsFunc(session.User.Roles, func(userRole *auth.Role) bool {
				return role == userRole.Id
			})

			if !match {
				rolesToAdd = append(rolesToAdd, role)
			}
		}

		for _, role := range session.User.Roles {
			match := slices.ContainsFunc(*updateUserValidator.Roles, func(userRole string) bool {
				return role.Id == userRole
			})

			if !match {
				rolesToRemove = append(rolesToRemove, role)
			}
		}

		if len(rolesToAdd) > 0 {
			addRolesQuery := sqlbuilder.InsertInto("open_board_user_roles").Cols("user_id", "role_id")

			for _, role := range rolesToAdd {
				addRolesQuery.Values(session.User.Id, role)
			}

			if err := transaction.Exec(addRolesQuery); err != nil {
				return err
			}
		}

		if len(rolesToRemove) > 0 {
			removeRolesQuery := sqlbuilder.DeleteFrom("open_board_user_roles")
			removeRolesQuery.Where(removeRolesQuery.In("role_id", rolesToRemove...))

			if err := transaction.Exec(removeRolesQuery); err != nil {
				return err
			}
		}

		if err := transaction.Exec(updateUserQuery); err != nil {
			return err
		}

		if err := transaction.Commit(); err != nil {
			return err
		}
	} else {
		if err := db.Instance.Exec(updateUserQuery); err != nil {
			return err
		}
	}

	user, err := auth.GetUser(session.User.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, fiber.Map{
			"message": fmt.Sprintf("user: %s has been successfully updated", user.Username),
			"user":    user,
		}),
	)
}

// UserInfoDELETE godoc
//
//	@Description 	Delete own user account. IMPORTANT: The actual endpoint is '/auth/@me' (with @ symbol)
//	@Summary 		Delete user account
//	@Tags			authentication
//	@Success		200 {object} responses.OkResponse[responses.GenericMessage]
//	@Failure		401,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/auth/me [delete]
//	@Produce		json
func UserInfoDELETE(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(auth.UserSession)
	deleteUserQuery := sqlbuilder.DeleteFrom("open_board_user")
	deleteUserQuery.Where(deleteUserQuery.Equal("id", session.User.Id))

	if err := db.Instance.Exec(deleteUserQuery); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: fmt.Sprintf("user: %s has been successfully deleted", session.User.Username)},
		),
	)
}

// UserSessionsGET godoc
//
//	@Description 	List own active auth sessions. IMPORTANT: The actual endpoint is '/auth/@me/sessions' (with @ symbol)
//	@Summary 		List sessions
//	@Tags			authentication
//	@Success		200 {object} responses.OkCollectionResponse[auth.UserSessionReadOnly]
//	@Failure		401,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/auth/me/sessions [get]
//	@Produce		json
func UserSessionsGET(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(auth.UserSession)
	userSessions := make([]auth.UserSessionReadOnly, 0)
	userSessionQuery := sqlbuilder.Select("id", "date_created", "date_updated", "user_agent", "ip_address").From("open_board_user_session")
	userSessionQuery.Where(userSessionQuery.Equal("user_id", session.User.Id))

	if err := db.Instance.Many(userSessionQuery, &userSessions); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, userSessions))
}
