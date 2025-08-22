package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/email"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/log"
	"VEDA95/open_board/api/internal/service"
	genericError "errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pquerna/otp/totp"
	"github.com/wneessen/go-mail"
)

type AuthHandler struct {
	authService  *service.AuthService
	userService  *service.UserService
	emailService *service.EmailService
	validator    *validators.Validator
}

func NewAuthHandler(
	authService *service.AuthService,
	userService *service.UserService,
	emailService *service.EmailService,
	validator *validators.Validator,
) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		userService:  userService,
		emailService: emailService,
		validator:    validator,
	}
}

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
func (authHandler *AuthHandler) LocalLogin(context *fiber.Ctx) error {
	dataValidator := new(validators.LocalLoginValidator)

	if err := context.BodyParser(dataValidator); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(dataValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	authData, err := authHandler.authService.LocalLogin(dataValidator, context.Get("User-Agent"), context.IP())
	if err != nil {
		return err
	}

	if dataValidator.ReturnType == "session" {
		context.Status(fiber.StatusCreated)
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session",
			Value:    authData.AccessToken,
			Expires:  authData.ExpiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   "localhost",
		})

		if dataValidator.Remember {
			context.Cookie(&fiber.Cookie{
				Name:     "open_board_session_remember_me",
				Value:    *authData.RefreshToken,
				Expires:  *authData.RefreshExpiresOn,
				HTTPOnly: true,
				Secure:   false,
				Path:     "/",
				Domain:   "localhost",
			})
		}

		return nil
	}

	responseData := auth.LocalUserLogin{
		AccessToken: authData.AccessToken,
		ExpiresIn:   authData.ExpiresOn,
	}

	if dataValidator.Remember {
		responseData.RefreshExpiresIn = authData.RefreshExpiresOn
		responseData.RefreshToken = authData.RefreshToken
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
	if len(user.Id) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	now := time.Now()
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "open_board",
		AccountName: user.Email,
	})
	if err != nil {
		return err
	}

	token, err := totp.GenerateCode(key.Secret(), now)
	if err != nil {
		return err
	}

	passwordResetQuery := sqlbuilder.InsertInto("open_board_password_reset_token").
		Cols("expires_on", "type", "user_id", "token").
		Values(now.Add(time.Minute*15), "form", user.Id, token)

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
			user.Email,
			mail.TypeTextHTML,
			email.MailTemplateStore.RenderTemplate(
				"password-reset",
				auth.PasswordResetEmailVariables{Token: token, Email: user.Email, Username: user.Username},
			),
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

func LocalAuthenticatedPasswordTokenIssuer(context *fiber.Ctx) error {
	validatorData := new(validators.AuthenticatedPasswordResetValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	authSession := context.Locals("auth_session").(auth.UserSession)

	if !auth.CheckPasswordHash(validatorData.Password, authSession.User.HashedPassword) {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	now := time.Now()
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "open_board",
		AccountName: authSession.User.Email,
	})
	if err != nil {
		return err
	}

	token, err := totp.GenerateCode(key.Secret(), now)
	if err != nil {
		return err
	}

	passwordResetQuery := sqlbuilder.InsertInto("open_board_password_reset_token").
		Cols("expires_on", "type", "user_id", "token").
		Values(now.Add(time.Minute*15), "auth", authSession.User.Id, token)

	if err := db.Instance.Exec(passwordResetQuery); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			"Password reset token was created successfully",
			auth.AuthenticatedPasswordResetResponse{
				Token: token,
			},
		),
	)
}

func LocalUnauthenticatedPasswordResetTokenIntrospect(context *fiber.Ctx) error {
	validatorData := new(validators.ResetPasswordTokenIntrospectValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}
	if errs := validators.Instance.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	resetToken := new(auth.PasswordResetToken)
	resetTokenQuery := sqlbuilder.Select("id", "date_created", "type", "user_id", "expires_on", "token").From("open_board_password_reset_token")

	resetTokenQuery.Where(resetTokenQuery.Equal("token", validatorData.Token))

	if err := db.Instance.One(resetTokenQuery, resetToken); err != nil {
		return err
	}
	if resetToken == nil {
		return responses.JSONResponse(context, fiber.StatusBadRequest, fiber.Map{"valid": false})
	}

	now := time.Now().Local()
	deleteResetTokenQuery := sqlbuilder.DeleteFrom("open_board_password_reset_token")

	deleteResetTokenQuery.Where(deleteResetTokenQuery.Equal("token", validatorData.Token))

	if now.After(resetToken.ExpiresOn) {
		if err := db.Instance.Exec(deleteResetTokenQuery); err != nil {
			return err
		}

		return responses.JSONResponse(context, fiber.StatusBadRequest, fiber.Map{"valid": false})
	}

	return responses.JSONResponse(context, fiber.StatusOK, fiber.Map{"valid": true})
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
	resetTokenQuery := sqlbuilder.Select("id", "date_created", "type", "user_id", "expires_on", "token").From("open_board_password_reset_token")

	resetTokenQuery.Where(resetTokenQuery.Equal("token", validatorData.Token))

	if err := db.Instance.One(resetTokenQuery, resetToken); err != nil {
		return err
	}
	if resetToken == nil {
		return fiber.NewError(fiber.StatusNotFound, "reset token not found")
	}

	now := time.Now().Local()
	deleteResetTokenQuery := sqlbuilder.DeleteFrom("open_board_password_reset_token")

	deleteResetTokenQuery.Where(deleteResetTokenQuery.Equal("token", validatorData.Token))

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

	updateUserQuery.Where(updateUserQuery.Equal("id", resetToken.UserId)).Set(
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

func LocalAuthenticatedPasswordReset(context *fiber.Ctx) error {
	validatorData := new(validators.ResetPasswordValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := validators.Instance.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	resetToken := new(auth.PasswordResetToken)
	resetTokenQuery := sqlbuilder.Select("id", "date_created", "type", "user_id", "expires_on", "token").From("open_board_password_reset_token")

	resetTokenQuery.Where(resetTokenQuery.Equal("token", validatorData.Token))

	if err := db.Instance.One(resetTokenQuery, resetToken); err != nil {
		return err
	}
	if resetToken == nil {
		return fiber.NewError(fiber.StatusNotFound, "reset token not found")
	}

	authSession := context.Locals("auth_session").(auth.UserSession)

	if resetToken.UserId != authSession.User.Id {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	now := time.Now().Local()
	deleteResetTokenQuery := sqlbuilder.DeleteFrom("open_board_password_reset_token")

	deleteResetTokenQuery.Where(deleteResetTokenQuery.Equal("token", validatorData.Token))

	if now.After(resetToken.ExpiresOn) {
		if err := db.Instance.Exec(deleteResetTokenQuery); err != nil {
			return err
		}

		return fiber.NewError(fiber.StatusBadRequest, "reset token is expired")
	}

	if resetToken.Type != "auth" {
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

	updateUserQuery.Where(updateUserQuery.Equal("id", resetToken.UserId)).Set(
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

	return responses.JSONResponse(context, fiber.StatusOK, responses.GenericMessage{Message: "Password Reset successfully!"})
}
