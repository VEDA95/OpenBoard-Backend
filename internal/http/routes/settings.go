package routes

import (
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

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, authSettings))
	}

	if paramsData.Type == "email" {
		emailSettings, err := settingsHandler.settingsService.GetEmailSettings()
		if err != nil {
			return err
		}

		return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, emailSettings))
	}

	generalSettings, err := settingsHandler.settingsService.GetGeneralSettings()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, generalSettings))
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
