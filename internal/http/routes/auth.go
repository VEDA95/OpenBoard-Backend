package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/log"
	genericError "errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-sqlbuilder"
	"os"
	"strconv"
	"strings"
	"time"
)

// LocalLogin godoc
//
//		@Description	Handles local login process for users
//		@Summary		logs in user
//	 	@Tags			auth
//		@Param 			request body validators.LocalLoginValidator false "Request Data"
//		@Success		201	{object} responses.OkResponse[auth.LocalUserLogin]
//		@Failure		401 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Failure		404 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//		@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Router			/auth/login [post]
//		@Accept			json
//		@Produce		json
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

	userQuery := auth.UsersQuery()
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
		now.Add(time.Second * time.Duration(expiresIn)),
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

		queryColumns = append(queryColumns, "refresh_expires_on", "refresh_token")
		queryValues = append(
			queryValues,
			now.Add(time.Second*time.Duration(refreshExpiresIn)),
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
		Set(updateUserQuery.Assign("last_login", now))

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
			Domain:   "localhost",
		})

		if dataValidator.Remember {
			context.Cookie(&fiber.Cookie{
				Name:     "open_board_session_remember_me",
				Value:    queryValues[len(queryValues)-1].(string),
				Expires:  queryValues[len(queryValues)-2].(time.Time),
				HTTPOnly: true,
				Secure:   false,
				Path:     "/",
				Domain:   "localhost",
			})
		}

		return nil
	}

	rows := make([]map[string]interface{}, 0)
	userRolesQuery := auth.UsersRolesQuery()
	userRolesQuery.Where(userRolesQuery.Equal("open_board_user_roles.user_id", user.Id))

	if err := db.Instance.Many(userRolesQuery, &rows); err != nil {
		return err
	}

	auth.AppendRolesToUser(rows, user)

	user.LastLogin = &now
	responseData := auth.LocalUserLogin{
		Message:     fmt.Sprintf("%s has been successfully logged in", user.Username),
		User:        user,
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}

	if dataValidator.Remember {
		responseData.RefreshExpiresIn = &refreshExpiresIn
		responseData.RefreshToken = queryValues[len(queryValues)-1].(*string)
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.OKResponse(fiber.StatusCreated, responseData),
	)
}

// LocalLogout godoc
//
//		@Description	Handles local logout process for users
//		@Summary		logs out user
//	 	@Tags			auth
//		@Param 			request body validators.LocalLogoutBodyValidator false "Request Data"
//		@Success		200	{object} responses.OkResponse[responses.GenericMessage]
//		@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//		@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Router			/auth/logout [post]
//		@Accept			json
//		@Produce		json
func LocalLogout(context *fiber.Ctx) error {
	logoutValidator := new(validators.LocalLogoutBodyValidator)

	if err := context.BodyParser(logoutValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(logoutValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	log.Logger.Debug().Interface("body", logoutValidator).Msg("LocalLogout")

	session := context.Locals("auth_session").(auth.UserSession)
	deleteSessionQuery := sqlbuilder.DeleteFrom("open_board_user_session")

	if logoutValidator.All {
		deleteSessionQuery.Where(deleteSessionQuery.Equal("user_id", session.User.Id))

	} else {
		deleteSessionQuery.Where(deleteSessionQuery.Equal("id", session.Id))
	}

	if err := db.Instance.Exec(deleteSessionQuery); err != nil {
		return err
	}

	if !logoutValidator.All && logoutValidator.ReturnType == "session" {
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

// LocalLogoutById godoc
//
//		@Description	Handles local logout process for users by session ID
//		@Summary		logs out user by session ID
//	 	@Tags			auth
//		@Param			id path string true "Session ID"
//		@Param 			request body validators.ReturnValidator false "Request Data"
//		@Success		200	{object} responses.OkResponse[responses.GenericMessage]
//		@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//		@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Router			/auth/logout/{id} [post]
//		@Accept			json
//		@Produce		json
func LocalLogoutById(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	logoutValidator := new(validators.ReturnValidator)

	if err := context.BodyParser(logoutValidator); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(logoutValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	log.Logger.Debug().Interface("body", logoutValidator).Msg("LocalLogout")

	session := context.Locals("auth_session").(auth.UserSession)
	deleteSessionQuery := sqlbuilder.DeleteFrom("open_board_user_session")
	deleteSessionQuery.Where(
		deleteSessionQuery.Equal("id", paramValidator.Id),
		deleteSessionQuery.Equal("user_id", session.User.Id),
	)

	if err := db.Instance.Exec(deleteSessionQuery); err != nil {
		return err
	}

	if (paramValidator.Id == session.Id) && logoutValidator.ReturnType == "session" {
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

// LocalRefresh godoc
//
//		@Description	Handles local user session renewal
//		@Summary		refreshes user session
//	 	@Tags			auth
//		@Param 			request body validators.ReturnValidator false "Request Data"
//		@Success		200	{object} responses.OkResponse[auth.LocalUserLogin]
//		@Failure		401 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Failure		404 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//		@Failure		500 {object} responses.ErrorResponse[responses.GenericMessage]
//		@Router			/auth/refresh [post]
//		@Accept			json
//		@Produce		json
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

	var session auth.UserSession
	now := time.Now()
	authToken := authHeaderSplit[1]
	sessionQuery := auth.GetSessionQuery()
	sessionQuery.Where(sessionQuery.Equal("refresh_token", authToken))

	if err := db.Instance.One(sessionQuery, &session); err != nil {
		return err
	}

	if len(session.Id) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "token not found")
	}

	if session.RefreshExpiresOn == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	if now.After(*session.RefreshExpiresOn) {
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
			updateSessionQuery.Assign("date_updated", now),
			updateSessionQuery.Assign("expires_on", expiresOn),
			updateSessionQuery.Assign("refresh_expires_on", refreshExpiresOn),
			updateSessionQuery.Assign("access_token", accessToken),
			updateSessionQuery.Assign("refresh_token", refreshToken),
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
			Domain:   "localhost",
		})
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session_remember_me",
			Value:    refreshToken,
			Expires:  refreshExpiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   "localhost",
		})

		return nil
	}

	rows := make([]map[string]interface{}, 0)
	usersRolesQuery := auth.UsersRolesQuery()
	usersRolesQuery.Where(usersRolesQuery.Equal("open_board_user_roles.user_id", session.User.Id))

	if err := db.Instance.Many(usersRolesQuery, &rows); err != nil {
		return err
	}

	auth.AppendRolesToUser(rows, session.User)

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, auth.LocalUserLogin{
			Message:          fmt.Sprintf("%s auth session has been successfully refreshed", session.User.Username),
			User:             session.User,
			AccessToken:      accessToken,
			RefreshToken:     &refreshToken,
			ExpiresIn:        expiresIn,
			RefreshExpiresIn: &refreshExpiresIn,
		}),
	)
}
