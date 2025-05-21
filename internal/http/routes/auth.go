package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/email"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/log"
	genericError "errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-sqlbuilder"
	"github.com/wneessen/go-mail"
	"os"
	"strconv"
	"strings"
	"time"
)

// LocalLogin godoc
//
//		@Description	Handles local login process for users
//		@Summary		Logs in user
//	 	@Tags			authentication
//		@Param 			request body validators.LocalLoginValidator false "Request Data"
//		@Success		201	{object} responses.SuccessResponse[auth.LocalUserLogin]
//		@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//		@Failure		401,404,500 {object} responses.ErrorResponse[responses.GenericMessage]
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

	now := time.Now().Local()
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
		User:        user,
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}

	if dataValidator.Remember {
		refreshToken := queryValues[len(queryValues)-1].(string)
		responseData.RefreshExpiresIn = &refreshExpiresIn
		responseData.RefreshToken = &refreshToken
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.CreateSuccessResponse(
			fiber.StatusCreated,
			fmt.Sprintf("%s has been successfully logged in", user.Username),
			responseData,
		),
	)
}

// LocalLogout godoc
//
//		@Description	Handles local logout process for users
//		@Summary		Logs out user
//	 	@Tags			authentication
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
//		@Summary		Logs out user by session ID
//	 	@Tags			authentication
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
//		@Summary		Refreshes user session
//	 	@Tags			authentication
//		@Param 			request body validators.ReturnValidator false "Request Data"
//		@Success		200	{object} responses.SuccessResponse[auth.LocalUserLogin]
//		@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//		@Failure		401,404,500 {object} responses.ErrorResponse[responses.GenericMessage]
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
	now := time.Now().Local()
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
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			fmt.Sprintf("%s auth session has been successfully refreshed", session.User.Username),
			auth.LocalUserLogin{
				User:             session.User,
				AccessToken:      accessToken,
				RefreshToken:     &refreshToken,
				ExpiresIn:        expiresIn,
				RefreshExpiresIn: &refreshExpiresIn,
			},
		),
	)
}

func LocalUnauthenticatedPasswordTokenIssuer(context *fiber.Ctx) error {
	validatorData := new(validators.ResetPasswordUserLookupValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	user := new(auth.User)
	usersQuery := auth.UsersQuery()

	usersQuery.Where(usersQuery.Equal("email", validatorData.Email))

	if err := db.Instance.One(usersQuery, user); err != nil {
		return err
	}

	if user == nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	token, err := auth.CreateSessionToken()

	if err != nil {
		log.Logger.Err(err).Msg("unable to create password reset token")
	}

	passwordResetQuery := sqlbuilder.InsertInto("open_board_password_reset_token").
		Cols("id", "expires_on", "type", "user_id").
		Values(token, time.Now().Add(time.Minute*15), "form", user.Id)

	if err := db.Instance.Exec(passwordResetQuery); err != nil {
		return err
	}

	go func() {
		if email.MailClient == nil {
			log.Logger.Warn().Msg("email client is nil. Skipping sending email...")
			return
		}

		err := email.MailClient.SendMessage(
			"Open Board Password Reset",
			validatorData.Email,
			mail.TypeTextHTML,
			email.MailTemplateStore.RenderTemplate("login_password_reset", validatorData),
		)

		if err != nil {
			log.Logger.Warn().Err(err).Msg("unable to send email")
		}
	}()

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: "Please check your email to recover your password"},
		),
	)
}

func LocalUnauthenticatedPasswordReset(context *fiber.Ctx) error {
	validatorData := new(validators.ResetPasswordValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	resetToken := new(auth.PasswordResetToken)
	resetTokenQuery := sqlbuilder.Select("id", "date_created", "type", "user_id", "expires_on").From("open_board_password_reset_token")

	resetTokenQuery.Where(resetTokenQuery.Equal("id", validatorData.Token))

	if err := db.Instance.One(resetTokenQuery, resetToken); err != nil {
		return err
	}

	if resetToken == nil {
		return fiber.NewError(fiber.StatusNotFound, "reset token not found")
	}

	now := time.Now().Local()
	deleteResetTokenQuery := sqlbuilder.DeleteFrom("open_board_password_reset_token")

	deleteResetTokenQuery.Where(deleteResetTokenQuery.Equal("id", validatorData.Token))

	if now.After(resetToken.ExpiresOn) {
		if err := db.Instance.Exec(deleteResetTokenQuery); err != nil {
			return err
		}

		return fiber.NewError(fiber.StatusBadRequest, "reset token is expired")
	}

	if resetToken.Type != "form" {
		return fiber.NewError(fiber.StatusBadRequest, "reset token is invalid")
	}

	hashedPassword, err := auth.HashPassword(validatorData.NewPassword)

	if err != nil {
		return err
	}

	transaction, err := db.Instance.Begin()

	if err != nil {
		return err
	}

	updateUserQuery := sqlbuilder.Update("open_board_user")
	deleteSessionQuery := sqlbuilder.DeleteFrom("open_board_user_session")

	updateUserQuery.Where(updateUserQuery.Assign("id", resetToken.UserId)).Set(
		updateUserQuery.Assign("date_updated", now),
		updateUserQuery.Assign("hashed_password", hashedPassword),
	)
	deleteSessionQuery.Where(deleteSessionQuery.Equal("user_id", resetToken.UserId))

	if err := transaction.Exec(updateUserQuery); err != nil {
		return err
	}

	if err := transaction.Exec(deleteSessionQuery); err != nil {
		return err
	}

	if err := transaction.Exec(deleteResetTokenQuery); err != nil {
		return err
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.GenericMessage{Message: "Password reset successfully!"})
}
