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
	"strings"
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

	userQuery := sqlbuilder.Select(userQueryColumns...).From("open_board_user")
	userQuery.Where(userQuery.Equal("username", dataValidator.Username))
	user := new(auth.User)

	if err := db.Instance.One(userQuery, user); err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	if !auth.CheckPasswordHash(dataValidator.Password, user.HashedPassword) {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	token, err := auth.CreateSessionToken()

	if err != nil {
		return err
	}

	now := time.Now()
	queryColumns := []string{"user_id", "expires_on", "session_type", "access_token", "ip_address", "user_agent"}
	queryValues := []interface{}{
		user.Id,
		now.Local().Add(time.Second * time.Duration(expiresIn)),
		"local",
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
			now.Local().Add(time.Second*time.Duration(refreshExpiresIn)),
			refreshToken,
		)
	}

	transaction, err := db.Instance.Begin()

	if err != nil {
		return err
	}

	sessionQuery := sqlbuilder.InsertInto("open_board_user_session").Cols(queryColumns...).Values(queryValues...)
	updateUserQuery := sqlbuilder.Update("open_board_user")
	updateUserQuery.
		Where(updateUserQuery.Equal("id", user.Id)).
		Set(updateUserQuery.Equal("last_login", now))

	if err := transaction.Exec(sessionQuery); err != nil {
		return err
	}

	if err := transaction.Exec(updateUserQuery); err != nil {
		return err
	}

	if err := transaction.Commit(); err != nil {
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

	user.LastLogin = &now
	responseMap := fiber.Map{
		"message":      fmt.Sprintf("%s has been successfully logged in", user.Username),
		"user":         user,
		"access_token": token,
		"expires_in":   expiresIn,
	}

	if dataValidator.Remember {
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
	returnValidator := new(validators.ReturnValidator)

	if err := context.BodyParser(returnValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(returnValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(auth.UserSession)
	deleteSessionQuery := sqlbuilder.DeleteFrom("open_board_user_session")
	deleteSessionQuery.Where(deleteSessionQuery.Equal("id", session.Id))

	if err := db.Instance.Exec(deleteSessionQuery); err != nil {
		return err
	}

	if returnValidator.ReturnType == "session" {
		context.Status(fiber.StatusOK)
		context.ClearCookie("open_board_session")

		if session.RefreshToken != nil && len(*session.RefreshToken) > 0 {
			context.ClearCookie("open_board_session_remember_me")
		}

		return nil
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: fmt.Sprintf("%s logged out successfully", session.User.Username)}),
	)
}

func LocalRefresh(context *fiber.Ctx) error {
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

	returnValidator := new(validators.ReturnValidator)

	if err := context.BodyParser(returnValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(returnValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	authHeader := context.Get("Authorization")

	if len(authHeader) == 0 {
		authCookie := context.Cookies("open_board_session_remember_me", "")

		if len(authCookie) == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		authHeader = authCookie
	}

	authHeaderSplit := strings.Split(authHeader, " ")

	if len(authHeaderSplit) < 2 {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	session := new(auth.UserSession)
	now := time.Now().Local()
	authToken := authHeaderSplit[1]
	sessionQuery := sqlbuilder.Select(
		"id",
		"date_created",
		"date_updated",
		"expires_on",
		"refresh_expires_on",
		"session_type",
		"access_token",
		"refresh_token",
		"ip_address",
		"user_agent",
		"additional_info",
		"open_board_user.id",
		"open_board_user.date_created",
		"open_board_user.date_updated",
		"open_board_user.last_login",
		"open_board_user.username",
		"open_board_user.email",
		"open_board_user.first_name",
		"open_board_user.last_name",
		"open_board_user.hashed_password",
		"open_board_user.enabled",
		"open_board_user.email_verified",
	).From("open_board_user_session")
	sessionQuery.
		Where(sessionQuery.Equal("refresh_token", authToken)).
		Join("open_board_user", "open_board_user_session.user_id = open_board_user.id")

	if err := db.Instance.One(sessionQuery, session); err != nil {
		return err
	}

	if session == nil {
		return fiber.NewError(fiber.StatusNotFound, "token not found")
	}

	if session.RefreshExpiresOn.After(now) {
		deleteSessionQuery := sqlbuilder.DeleteFrom("open_board_user_session")
		deleteSessionQuery.Where(deleteSessionQuery.Equal("id", session.Id))

		if err := db.Instance.Exec(deleteSessionQuery); err != nil {
			return err
		}

		if returnValidator.ReturnType == "session" {
			context.ClearCookie("open_board_session", "open_board_session_remember_me")
		}

		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	accessToken, err := auth.CreateSessionToken()

	if err != nil {
		return err
	}

	refreshToken, err := auth.CreateSessionToken()

	if err != nil {
		return err
	}

	expiresOn := now.Add(time.Second * time.Duration(expiresIn))
	refreshExpiresOn := now.Add(time.Second * time.Duration(refreshExpiresIn))
	updateSessionQuery := sqlbuilder.Update("open_board_user_session")
	updateSessionQuery.
		Where(updateSessionQuery.Equal("id", session.Id)).
		Set(
			sessionQuery.Equal("date_updated", now),
			sessionQuery.Equal("expires_on", expiresOn),
			sessionQuery.Equal("refresh_expires_on", refreshExpiresOn),
			sessionQuery.Equal("access_token", accessToken),
			sessionQuery.Equal("refresh_token", refreshToken),
		)

	if err := db.Instance.Exec(updateSessionQuery); err != nil {
		return err
	}

	if returnValidator.ReturnType == "session" {
		context.Status(fiber.StatusOK)
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session",
			Value:    accessToken,
			Expires:  expiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   "localhost:8080",
		})
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session_remember_me",
			Value:    refreshToken,
			Expires:  refreshExpiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   "localhost:8080",
		})

		return nil
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, fiber.Map{
			"message":            fmt.Sprintf("%s auth session has been successfully refreshed", session.User.Username),
			"user":               session.User,
			"access_token":       accessToken,
			"refresh_token":      refreshToken,
			"expires_in":         expiresIn,
			"refresh_expires_in": refreshExpiresOn,
		}),
	)
}
