package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	genericError "errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-sqlbuilder"
	"os"
	"strconv"
	"time"
)

var userQueryColumns = []string{
	"id",
	"date_created",
	"date_updated",
	"last_login",
	"username",
	"email",
	"first_name",
	"last_name",
	"hashed_password",
	"enabled",
	"email_verified",
}

func LocalLogin(context *fiber.Ctx) error {
	expiresInEnv := os.Getenv("AUTH_SESSION_EXPIRES_IN")
	refreshExpiresInEnv := os.Getenv("AUTH_SESSION_REFRESH_EXPIRES_IN")

	if len(expiresInEnv) == 0 || len(refreshExpiresInEnv) == 0 {
		return genericError.New("AUTH_SESSION_EXPIRES_IN and/or AUTH_SESSION_REFRESH_EXPIRES_IN environment variable(s) was not set")
	}

	expiresIn, err := strconv.Atoi(expiresInEnv)

	if err != nil {
		return err
	}

	refreshExpiresIn, err := strconv.Atoi(refreshExpiresInEnv)

	if err != nil {
		return err
	}

	dataValidator := new(validators.LocalLoginValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	userQuery := sqlbuilder.Select(userQueryColumns...).From("user")
	userQuery.Where(userQuery.Equal("username", dataValidator.Username))
	user := new(auth.User)

	if err := db.Instance.One(userQuery, user); err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	if !auth.CheckPasswordHash(dataValidator.Password, user.HashedPassword) {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authorized")
	}

	token, err := auth.CreateSessionToken()

	if err != nil {
		return err
	}

	now := time.Now()
	queryColumns := []string{"user_id", "expires_on", "session_type", "access_token", "ip_address", "user_agent"}
	queryValues := []interface{}{
		user.Id,
		now.Local().Add(time.Duration(expiresIn)),
		dataValidator.ReturnType,
		token,
		context.IP(),
		context.Get("User-Agent"),
	}

	if dataValidator.Remember {
		refreshToken, err := auth.CreateSessionToken()

		if err != nil {
			return err
		}

		queryColumns = append(queryColumns, "remember_me", "refresh_expires_in", "refresh_token")
		queryValues = append(
			queryValues,
			true,
			now.Local().Add(time.Duration(refreshExpiresIn)),
			refreshToken,
		)
	}

	sessionQuery := sqlbuilder.InsertInto("open_board_user_session").Cols(queryColumns...).Values(queryValues...)

	if err := db.Instance.Exec(sessionQuery); err != nil {
		return err
	}

	if dataValidator.ReturnType == "session" {
		context.Status(fiber.StatusCreated)
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session",
			Value:    token,
			Expires:  queryValues[1].(time.Time),
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   "localhost:8080",
		})

		if dataValidator.Remember {
			context.Cookie(&fiber.Cookie{
				Name:     "open_board_session_remember_me",
				Value:    queryValues[len(queryValues)-1].(string),
				Expires:  queryValues[len(queryValues)-2].(time.Time),
				HTTPOnly: true,
				Secure:   false,
				Path:     "/",
				Domain:   "localhost:8080",
			})
		}

		return nil
	}

	responseMap := fiber.Map{
		"message":      fmt.Sprintf("%s has been successfully logged in", user.Username),
		"user":         user,
		"access_token": token,
		"expires_in":   expiresIn,
	}

	if dataValidator.Remember {
		responseMap["remember_me"] = true
		responseMap["refresh_expires_in"] = refreshExpiresIn
		responseMap["refresh_token"] = queryValues[len(queryValues)-1]
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.OKResponse(fiber.StatusOK, responseMap),
	)
}

func LocalLogout(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(auth.UserSession)
	deleteSessionQuery := sqlbuilder.DeleteFrom("open_board_user_session")
	deleteSessionQuery.Where(deleteSessionQuery.Equal("id", session.Id))

	if err := db.Instance.Exec(deleteSessionQuery); err != nil {
		return err
	}

	if session.SessionType == "session" {
		context.Status(fiber.StatusOK)
		context.ClearCookie("open_board_session")

		if session.Remember {
			context.ClearCookie("open_board_session_remember_me")
		}

		return nil
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.GenericMessage{Message: fmt.Sprintf("%s logged out successfully", session.User.Username)},
	)
}

func LocalRefresh(context *fiber.Ctx) error {
	return nil
}

func UserInfoGET(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(auth.UserSession)

	return responses.JSONResponse(context, fiber.StatusOK, session.User)
}

func UsersGET(context *fiber.Ctx) error {
	users := make([]auth.User, 0)
	userQuery := sqlbuilder.Select(userQueryColumns...).From("open_board_user")

	if err := db.Instance.Many(userQuery, &users); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, users)
}

func UsersPOST(context *fiber.Ctx) error {
	createValidator := new(validators.CreateUserValidator)

	if err := context.BodyParser(createValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(createValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	queryColumns := []string{"username", "email", "hashed_password"}
	queryValues := []interface{}{
		createValidator.Username,
		createValidator.Email,
		auth.HashPassword(createValidator.Password),
	}

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
		fiber.Map{
			"message": fmt.Sprintf("The user: %s has been successfully created!", user.Username),
			"user":    user,
		},
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

	return responses.JSONResponse(context, fiber.StatusOK, user)
}

func UserPATCH(context *fiber.Ctx) error {
	return nil
}

func UserDelete(context *fiber.Ctx) error {
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
		responses.GenericMessage{Message: fmt.Sprintf("user: %s has been successfully deleted!", user.Username)},
	)
}
