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
	"time"
)

func UserInfoGET(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(auth.UserSession)

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, session.User))
}

func UsersGET(context *fiber.Ctx) error {
	users := make([]auth.User, 0)
	userQuery := sqlbuilder.Select(userQueryColumns...).From("open_board_user")

	if err := db.Instance.Many(userQuery, &users); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, users))
}

func UsersPOST(context *fiber.Ctx) error {
	createValidator := new(validators.CreateUserValidator)

	if err := context.BodyParser(createValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(createValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	hashedPassword, err := auth.HashPassword(createValidator.Password)

	if err != nil {
		return err
	}

	queryColumns := []string{"username", "email", "hashed_password"}
	queryValues := []interface{}{createValidator.Username, createValidator.Email, hashedPassword}

	if createValidator.FirstName != nil && len(*createValidator.FirstName) > 0 {
		queryColumns = append(queryColumns, "first_name")
		queryValues = append(queryValues, *createValidator.FirstName)
	}

	if createValidator.LastName != nil && len(*createValidator.LastName) > 0 {
		queryColumns = append(queryColumns, "last_name")
		queryValues = append(queryValues, *createValidator.LastName)
	}

	var userId string
	createUserQuery := sqlbuilder.InsertInto("open_board_user").Cols(queryColumns...).Values(queryValues...).Returning("id")

	if err := db.Instance.One(createUserQuery, &userId); err != nil {
		return err
	}

	user := new(auth.User)
	userQuery := sqlbuilder.Select(userQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("id", userId))

	if err := db.Instance.One(userQuery, user); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.OKResponse(fiber.StatusCreated, fiber.Map{
			"message": fmt.Sprintf("The user: %s has been successfully created!", user.Username),
			"user":    user,
		}),
	)
}

func UserGET(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	user := new(auth.User)
	userQuery := sqlbuilder.Select(userQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("id", paramValidator.Id))

	if err := db.Instance.One(userQuery, user); err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, user))
}

func UserPATCH(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	user := new(auth.User)
	userQuery := sqlbuilder.Select(userQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("id", paramValidator.Id))

	if err := db.Instance.One(userQuery, user); err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	updateUserValidator := new(validators.UpdateUserValidator)

	if err := context.BodyParser(updateUserValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(updateUserValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	updateUserQuery := sqlbuilder.Update("open_board_user")
	updateUserQuery.
		Where(updateUserQuery.Equal("id", paramValidator.Id)).
		Set(updateUserQuery.Equal("date_updated", time.Now().Local()))

	if updateUserValidator.Username != nil && len(*updateUserValidator.Username) > 0 && *updateUserValidator.Username != user.Username {
		updateUserQuery.SetMore(updateUserQuery.Equal("username", *updateUserValidator.Username))
	}

	if updateUserValidator.Email != nil && len(*updateUserValidator.Email) > 0 && *updateUserValidator.Email != user.Email {
		updateUserQuery.SetMore(updateUserQuery.Equal("email", updateUserValidator.Email))
	}

	if updateUserValidator.FirstName != nil && updateUserValidator.FirstName != user.FirstName {
		if len(*updateUserValidator.FirstName) == 0 && user.FirstName != nil {
			updateUserQuery.SetMore(updateUserQuery.Equal("first_name", nil))

		} else {
			updateUserQuery.SetMore(updateUserQuery.Equal("first_name", updateUserValidator.FirstName))
		}
	}

	if updateUserValidator.LastName != nil && updateUserValidator.LastName != user.LastName {
		if len(*updateUserValidator.LastName) == 0 && user.LastName != nil {
			updateUserQuery.SetMore(updateUserQuery.Equal("last_name", nil))
		} else {
			updateUserQuery.SetMore(updateUserQuery.Equal("last_name", updateUserValidator.LastName))
		}
	}

	if err := db.Instance.Exec(updateUserQuery); err != nil {
		return err
	}

	userQuery = sqlbuilder.Select(userQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("id", paramValidator.Id))

	if err := db.Instance.One(userQuery, user); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, fiber.Map{
			"message": fmt.Sprintf("user: %s has been successfully updated!", user.Username),
			"user":    user,
		}),
	)
}

func UserDELETE(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	user := new(auth.User)
	userQuery := sqlbuilder.Select(userQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("id", paramValidator.Id))

	if err := db.Instance.One(userQuery, user); err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	deleteUserQuery := sqlbuilder.DeleteFrom("open_board_user")
	deleteUserQuery.Where(deleteUserQuery.Equal("id", paramValidator.Id))

	if err := db.Instance.Exec(deleteUserQuery); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: fmt.Sprintf("user: %s has been successfully deleted!", user.Username)},
		),
	)
}
