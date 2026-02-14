package routes

import (
	"VEDA95/open_board/api/internal/auth"
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/log"
	"VEDA95/open_board/api/internal/service"
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService     *service.AuthService
	userService     *service.UserService
	emailService    *service.EmailService
	settingsService *service.SettingsService
	validator       *validators.Validator
}

func NewAuthHandler(
	authService *service.AuthService,
	userService *service.UserService,
	emailService *service.EmailService,
	settingsService *service.SettingsService,
	validator *validators.Validator,
) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		userService:     userService,
		emailService:    emailService,
		settingsService: settingsService,
		validator:       validator,
	}
}

func (authHandler *AuthHandler) getCookieDomain() string {
	authSettings, err := authHandler.settingsService.GetAuthSettings()

	if err != nil || len(authSettings.CORSDomain) == 0 {
		return "localhost"
	}
	// Extract domain from the CORS URL
	parsedURL, err := url.Parse(authSettings.CORSDomain)
	if err != nil {
		return "localhost"
	}
	return parsedURL.Hostname()
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
			Value:    authData.Session.AccessToken,
			Expires:  authData.Session.ExpiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   authHandler.getCookieDomain(),
		})

		if dataValidator.Remember {
			context.Cookie(&fiber.Cookie{
				Name:     "open_board_session_remember_me",
				Value:    *authData.Session.RefreshToken,
				Expires:  *authData.Session.RefreshExpiresOn,
				HTTPOnly: true,
				Secure:   false,
				Path:     "/",
				Domain:   authHandler.getCookieDomain(),
			})
		}

		return nil
	}

	responseData := auth.LocalUserLogin{
		AccessToken: authData.Session.AccessToken,
		ExpiresIn:   authData.ExpiresIn,
		User:        authData.Session.User,
	}

	if dataValidator.Remember {
		responseData.RefreshExpiresIn = authData.RefreshExpiresIn
		responseData.RefreshToken = authData.Session.RefreshToken
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.CreateSuccessResponse(
			fiber.StatusCreated,
			fmt.Sprintf("%s has been successfully logged in", authData.Session.User.Username),
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
func (authHandler *AuthHandler) LocalLogout(context *fiber.Ctx) error {
	logoutValidator := new(validators.LocalLogoutBodyValidator)

	if err := context.BodyParser(logoutValidator); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(logoutValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	if logoutValidator.All {
		authHandler.authService.LocalLogoutByUserID(session.User.ID)
	} else {
		authHandler.authService.LocalLogout(session.AccessToken)
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
func (authHandler *AuthHandler) LocalLogoutById(context *fiber.Ctx) error {
	paramValidator := new(validators.ParamValidator)

	if err := context.ParamsParser(paramValidator); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(paramValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	logoutValidator := new(validators.ReturnValidator)

	if err := context.BodyParser(logoutValidator); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(logoutValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	if err := authHandler.authService.LocalLogoutByID(paramValidator.Id); err != nil {
		return err
	}

	if (paramValidator.Id == session.ID) && logoutValidator.ReturnType == "session" {
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
func (authHandler *AuthHandler) LocalRefresh(context *fiber.Ctx) error {
	returnValidator := new(validators.ReturnValidator)

	if err := context.BodyParser(returnValidator); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(returnValidator); len(errs) > 0 {
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

	authToken := authHeaderSplit[1]
	authData, err := authHandler.authService.LocalRefresh(authToken)
	if err != nil {
		return err
	}

	if returnValidator.ReturnType == "session" {
		context.Status(fiber.StatusOK)
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session",
			Value:    authData.Session.AccessToken,
			Expires:  authData.Session.ExpiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   authHandler.getCookieDomain(),
		})
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session_remember_me",
			Value:    *authData.Session.RefreshToken,
			Expires:  *authData.Session.RefreshExpiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   authHandler.getCookieDomain(),
		})

		return nil
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			fmt.Sprintf("%s auth session has been successfully refreshed", authData.Session.User.Username),
			auth.LocalUserLogin{
				User:             authData.Session.User,
				AccessToken:      authData.Session.AccessToken,
				RefreshToken:     authData.Session.RefreshToken,
				ExpiresIn:        authData.ExpiresIn,
				RefreshExpiresIn: authData.RefreshExpiresIn,
			},
		),
	)
}

func (authHandler *AuthHandler) LocalUnauthenticatedPasswordTokenIssuer(context *fiber.Ctx) error {
	validatorData := new(validators.ResetPasswordUserLookupValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	passwordResetToken, err := authHandler.authService.IssueForgotPasswordToken(validatorData.Email)
	if err != nil {
		return err
	}

	go func() {
		subject := "password-reset"

		err := authHandler.emailService.SendMessage(
			"Open Board Password Reset",
			passwordResetToken.User.Email,
			auth.PasswordResetEmailVariables{
				Token:    passwordResetToken.Token,
				Email:    passwordResetToken.User.Email,
				Username: passwordResetToken.User.Username,
			},
			&subject,
		)
		if err != nil {
			log.Global.Warn().Err(err).Msg("an error occurred when sending a password reset email")
		}
	}()

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(
			fiber.StatusOK,
			responses.GenericMessage{Message: "Please check your email to reset your password"},
		),
	)
}

func (authHandler *AuthHandler) LocalAuthenticatedPasswordTokenIssuer(context *fiber.Ctx) error {
	validatorData := new(validators.AuthenticatedPasswordResetValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	authSession := context.Locals("auth_session").(models.Session)
	passwordResetToken, err := authHandler.authService.IssueUserPasswordResetToken(authSession.User, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			"Password reset token was created successfully",
			auth.AuthenticatedPasswordResetResponse{
				Token: passwordResetToken.Token,
			},
		),
	)
}

func (authHandler *AuthHandler) LocalUnauthenticatedPasswordReset(context *fiber.Ctx) error {
	validatorData := new(validators.ResetPasswordValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	if err := authHandler.authService.ResetForgottenPassword(validatorData); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.GenericMessage{Message: "Password reset successfully!"})
}

func (authHandler *AuthHandler) LocalAuthenticatedPasswordReset(context *fiber.Ctx) error {
	validatorData := new(validators.ResetPasswordValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(validatorData); errs != nil {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)
	if err := authHandler.authService.ResetUserPassword(session.User, validatorData); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.GenericMessage{Message: "Password Reset successfully!"})
}

// UserInfoGET godoc
//
//	@Description 	Retrieve own user details. IMPORTANT: The actual endpoint is '/auth/@me' (with @ symbol)
//	@Summary 		Retrieve user details
//	@Tags			authentication
//	@Success		200 {object} responses.OkResponse[models.User]
//	@Failure		401,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/auth/me [get]
//	@Produce		json
func (*AuthHandler) UserInfoGET(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, session.User))
}

// UserInfoPATCH godoc
//
//	@Description 	Update own user details. IMPORTANT: The actual endpoint is '/auth/@me' (with @ symbol)
//	@Summary 		Update user details
//	@Tags			authentication
//	@Param			request body validators.UpdateUserValidator false "request Data"
//	@Success		200 {object} responses.SuccessResponse[models.User]
//	@Failure		422 {object} responses.ErrorResponse[validators.ErrorResponseMap]
//	@Failure		401,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/auth/me [patch]
//	@Accept			json
//	@Produce		json
func (authHandler *AuthHandler) UserInfoPATCH(context *fiber.Ctx) error {
	updateUserValidator := new(validators.UpdateUserValidator)

	if err := context.BodyParser(updateUserValidator); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(updateUserValidator); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)
	user, err := authHandler.userService.UpdateUser(session.User.ID, updateUserValidator)
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
func (authHandler *AuthHandler) UserInfoDELETE(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	if err := authHandler.userService.DeleteUser(session.User.ID); err != nil {
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
//	@Success		200 {object} responses.OkCollectionResponse[models.Session]
//	@Failure		401,500 {object} responses.ErrorResponse[responses.GenericMessage]
//	@Router			/auth/me/sessions [get]
//	@Produce		json
func (*AuthHandler) UserSessionsGET(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, session.User.Sessions))
}

func (authHandler *AuthHandler) Register(context *fiber.Ctx) error {
	validatorData := new(validators.RegisterUserValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := authHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := authHandler.authService.RegisterUser(validatorData); err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.OKResponse(
			fiber.StatusCreated,
			responses.GenericMessage{Message: "user created!"},
		),
	)
}
