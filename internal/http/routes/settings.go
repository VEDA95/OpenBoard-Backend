package routes

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/log"
	"VEDA95/open_board/api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	settingsService *service.SettingsService
	emailService    *service.EmailService
	validator       *validators.Validator
}

func NewSettingsHandler(
	settingsService *service.SettingsService,
	emailService *service.EmailService,
	validator *validators.Validator,
) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		emailService:    emailService,
		validator:       validator,
	}
}

func (settingsHandler *SettingsHandler) GET(context *fiber.Ctx) error {
	paramsData := new(validators.SettingsTypeValidator)

	if err := context.ParamsParser(paramsData); err != nil {
		return err
	}

	if errs := settingsHandler.validator.Validate(paramsData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if paramsData.Type == "auth" {
		authSettings, err := settingsHandler.settingsService.GetAuthSettings()
		if err != nil {
			return err
		}

		session := context.Locals("auth_session").(models.Session)
		if !session.User.IsSuperuser() && !session.User.IsAuthorized("system:settings") {
			return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(
				fiber.StatusOK,
				fiber.Map{
					"allow_public_registration":          authSettings.AllowPublicRegistration,
					"require_email_verification":         authSettings.RequireEmailVerification,
					"two_factor_authentication":          authSettings.TwoFactorAuthentication,
					"two_factor_authentication_required": authSettings.TwoFactorAuthentication,
					"allow_user_invitations":             authSettings.AllowUserInvitations,
					"invite_only_mode":                   authSettings.InviteOnlyMode,
				},
			))
		}

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, authSettings))
	}

	if paramsData.Type == "email" {
		emailSettings, err := settingsHandler.settingsService.GetEmailSettings()
		if err != nil {
			return err
		}

		session := context.Locals("auth_session").(models.Session)
		if !session.User.IsSuperuser() && !session.User.IsAuthorized("system:settings") {
			return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(
				fiber.StatusOK,
				fiber.Map{
					"email_provider":             emailSettings.EmailProvider,
					"email_enabled":              emailSettings.EmailEnabled,
					"email_from_address":         emailSettings.EmailFromAddress,
					"email_from_name":            emailSettings.EmailFromName,
					"smtp_host":                  emailSettings.SMTPHost,
					"smtp_port":                  emailSettings.SMTPPort,
					"smtp_encryption":            emailSettings.SMTPEncryption,
					"enable_email_notifications": emailSettings.EnableEmailNotifications,
					"send_password_reset_email":  emailSettings.SendPasswordResetEmail,
					"notify_on_card_assigned":    emailSettings.NotifyOnCardAssigned,
					"notify_on_card_comment":     emailSettings.NotifyOnCardComment,
					"notify_on_card_due":         emailSettings.NotifyOnCardDue,
					"notify_on_board_invite":     emailSettings.NotifyOnBoardInvite,
					"notify_on_mention":          emailSettings.NotifyOnMention,
				},
			))
		}

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, emailSettings))
	}

	generalSettings, err := settingsHandler.settingsService.GetGeneralSettings()
	if err != nil {
		return err
	}

	session := context.Locals("auth_session").(models.Session)
	if !session.User.IsSuperuser() && !session.User.IsAuthorized("system:settings") {
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(
			fiber.StatusOK,
			fiber.Map{
				"app_name":                 generalSettings.AppName,
				"app_logo":                 generalSettings.AppLogo,
				"app_favicon":              generalSettings.AppFavicon,
				"app_description":          generalSettings.AppDescription,
				"default_language":         generalSettings.DefaultLanguage,
				"default_timezone":         generalSettings.DefaultTimezone,
				"show_announcement_banner": generalSettings.ShowAnnouncementBanner,
				"announcement_message":     generalSettings.AnnouncementMessage,
				"announcement_type":        generalSettings.AnnouncementType,
				"max_file_size":            generalSettings.MaxFileSize,
			},
		))
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, generalSettings))
}

func (settingsHandler *SettingsHandler) GETPublic(context *fiber.Ctx) error {
	paramsData := new(validators.SettingsTypeValidator)

	if err := context.ParamsParser(paramsData); err != nil {
		return err
	}

	if errs := settingsHandler.validator.Validate(paramsData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if paramsData.Type == "auth" {
		authSettings, err := settingsHandler.settingsService.GetAuthSettings()
		if err != nil {
			return err
		}

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(
			fiber.StatusOK,
			fiber.Map{
				"allow_public_registration":          authSettings.AllowPublicRegistration,
				"require_email_verification":         authSettings.RequireEmailVerification,
				"two_factor_authentication":          authSettings.TwoFactorAuthentication,
				"two_factor_authentication_required": authSettings.TwoFactorAuthentication,
				"allow_user_invitations":             authSettings.AllowUserInvitations,
				"invite_only_mode":                   authSettings.InviteOnlyMode,
			},
		))
	}

	if paramsData.Type == "email" {
		emailSettings, err := settingsHandler.settingsService.GetEmailSettings()
		if err != nil {
			return err
		}

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(
			fiber.StatusOK,
			fiber.Map{
				"email_enabled":              emailSettings.EmailEnabled,
				"send_password_reset_email":  emailSettings.SendPasswordResetEmail,
				"enable_email_notifications": emailSettings.EnableEmailNotifications,
				"notify_on_card_assigned":    emailSettings.NotifyOnCardAssigned,
				"notify_on_card_comment":     emailSettings.NotifyOnCardComment,
				"notify_on_card_due":         emailSettings.NotifyOnCardDue,
				"notify_on_board_invite":     emailSettings.NotifyOnBoardInvite,
				"notify_on_mention":          emailSettings.NotifyOnMention,
			},
		))
	}

	generalSettings, err := settingsHandler.settingsService.GetGeneralSettings()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(
		fiber.StatusOK,
		fiber.Map{
			"app_name":                 generalSettings.AppName,
			"app_logo":                 generalSettings.AppLogo,
			"app_favicon":              generalSettings.AppFavicon,
			"app_description":          generalSettings.AppDescription,
			"default_language":         generalSettings.DefaultLanguage,
			"default_timezone":         generalSettings.DefaultTimezone,
			"show_announcement_banner": generalSettings.ShowAnnouncementBanner,
			"announcement_message":     generalSettings.AnnouncementMessage,
			"announcement_type":        generalSettings.AnnouncementType,
			"max_file_size":            generalSettings.MaxFileSize,
		},
	))
}

func (settingsHandler *SettingsHandler) PATCH(context *fiber.Ctx) error {
	paramsData := new(validators.SettingsTypeValidator)

	if err := context.ParamsParser(paramsData); err != nil {
		return err
	}

	if errs := settingsHandler.validator.Validate(paramsData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if paramsData.Type == "auth" {
		validatorData := new(validators.AuthSettingsValidator)

		if err := context.BodyParser(validatorData); err != nil {
			return err
		}

		if errs := settingsHandler.validator.Validate(validatorData); len(errs) > 0 {
			return errors.CreateValidationError(errs)
		}

		authSettings, err := settingsHandler.settingsService.UpdateAuthSettings(validatorData)
		if err != nil {
			return err
		}

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, authSettings))
	}

	if paramsData.Type == "email" {
		validatorData := new(validators.EmailSettingsValidator)

		if err := context.BodyParser(validatorData); err != nil {
			return err
		}

		if errs := settingsHandler.validator.Validate(validatorData); len(errs) > 0 {
			return errors.CreateValidationError(errs)
		}

		emailSettings, err := settingsHandler.settingsService.UpdateEmailSettings(validatorData)
		if err != nil {
			return err
		}

		go func() {
			if err := settingsHandler.emailService.InitializeClient(); err != nil {
				log.Global.Warn().Err(err).Msg("an error occurred when initializing the email client")
			}
		}()

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, emailSettings))
	}

	validatorData := new(validators.GeneralSettingsValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := settingsHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	generalSettings, err := settingsHandler.settingsService.UpdateGeneralSettings(validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, generalSettings))
}
