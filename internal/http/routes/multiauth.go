package routes

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MultiAuthHandler struct {
	multiAuthService *service.MultiAuthService
	emailService     *service.EmailService
	validator        *validators.Validator
}

func NewMultiAuthHandler(
	multiAuthService *service.MultiAuthService,
	emailService *service.EmailService,
	validator *validators.Validator,
) *MultiAuthHandler {
	return &MultiAuthHandler{
		multiAuthService: multiAuthService,
		emailService:     emailService,
		validator:        validator,
	}
}

// GetMFAMethods returns all MFA methods for the authenticated user
// @Description Returns all MFA methods for the authenticated user
// @Summary Get MFA methods
// @Tags mfa
// @Success 200 {object} responses.OkCollectionResponse[models.MultiAuthMethod]
// @Failure 401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/methods [get]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) GetMFAMethods(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	methods, err := h.multiAuthService.GetUserMFAMethods(session.User.ID)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, methods))
}

// GetMFAStatus returns the MFA status for the authenticated user
// @Description Returns whether MFA is enabled/required for the user
// @Summary Get MFA status
// @Tags mfa
// @Success 200 {object} responses.OkResponse[fiber.Map]
// @Failure 401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/status [get]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) GetMFAStatus(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	hasEnabled := h.multiAuthService.HasMFAEnabled(session.User.ID)

	required, err := h.multiAuthService.IsMFARequired()
	if err != nil {
		return err
	}

	enabled, err := h.multiAuthService.IsMFAEnabled()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, fiber.Map{
		"has_mfa_enabled": hasEnabled,
		"mfa_required":    required,
		"mfa_enabled":     enabled,
	}))
}

// DeleteMFAMethod deletes an MFA method
// @Description Deletes an MFA method for the authenticated user
// @Summary Delete MFA method
// @Tags mfa
// @Param id path string true "MFA Method ID"
// @Success 200 {object} responses.OkResponse[responses.GenericMessage]
// @Failure 400,401,404,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/methods/{id} [delete]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) DeleteMFAMethod(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := h.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := h.multiAuthService.DeleteMFAMethod(params.Id, session.User.ID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "MFA method has been deleted successfully",
	}))
}

// --- Email TOTP Methods ---

// SetupEmailTOTP sets up email-based TOTP for the authenticated user
// @Description Sets up email-based TOTP for the authenticated user
// @Summary Setup email TOTP
// @Tags mfa
// @Param request body validators.SetupEmailTOTPValidator true "Request Data"
// @Success 201 {object} responses.SuccessResponse[models.MultiAuthMethod]
// @Failure 400,401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/email-totp/setup [post]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) SetupEmailTOTP(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	validatorData := new(validators.SetupEmailTOTPValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := h.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	method, err := h.multiAuthService.SetupEmailTOTP(session.User.ID, validatorData.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.CreateSuccessResponse(
			fiber.StatusCreated,
			"Email TOTP has been set up successfully",
			method,
		),
	)
}

// SendEmailTOTPCode sends a TOTP code via email
// @Description Sends a TOTP code to the user's email
// @Summary Send email TOTP code
// @Tags mfa
// @Success 200 {object} responses.SuccessResponse[responses.GenericMessage]
// @Failure 400,401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/email-totp/send [post]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) SendEmailTOTPCode(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	challenge, err := h.multiAuthService.GenerateEmailTOTPCode(session.User.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Send email with TOTP code
	if h.emailService != nil {
		subject := "Your verification code"
		err := h.emailService.SendMessage(
			challenge.Email,
			"mfa_totp_code", // Template name - to be created later
			map[string]any{
				"code":       challenge.Code,
				"expires_in": "10 minutes",
			},
			&subject,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to send verification code")
		}
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "Verification code has been sent to your email",
	}))
}

// VerifyEmailTOTPCode verifies the email TOTP code
// @Description Verifies the TOTP code sent via email
// @Summary Verify email TOTP code
// @Tags mfa
// @Param request body validators.VerifyEmailTOTPValidator true "Request Data"
// @Success 200 {object} responses.SuccessResponse[responses.GenericMessage]
// @Failure 400,401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/email-totp/verify [post]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) VerifyEmailTOTPCode(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	validatorData := new(validators.VerifyEmailTOTPValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := h.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := h.multiAuthService.VerifyEmailTOTPCode(session.User.ID, validatorData.Code); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "Verification successful",
	}))
}

// --- WebAuthn Methods ---

// GetWebAuthnRegistrationOptions returns options for WebAuthn registration
// @Description Returns options for WebAuthn registration ceremony
// @Summary Get WebAuthn registration options
// @Tags mfa
// @Success 200 {object} responses.SuccessResponse[fiber.Map]
// @Failure 401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/webauthn/register/options [get]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) GetWebAuthnRegistrationOptions(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	challenge, err := h.multiAuthService.GenerateWebAuthnChallenge()
	if err != nil {
		return err
	}

	userData, err := h.multiAuthService.GetUserForWebAuthn(session.User.ID)
	if err != nil {
		return err
	}

	// Store challenge in cookie for verification
	context.Cookie(&fiber.Cookie{
		Name:     "webauthn_challenge",
		Value:    challenge,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   300, // 5 minutes
	})

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, fiber.Map{
		"challenge": challenge,
		"rp": fiber.Map{
			"name": "Open Board",
			"id":   context.Hostname(),
		},
		"user":             userData,
		"pubKeyCredParams": getPublicKeyCredentialParams(),
		"authenticatorSelection": map[string]interface{}{
			"userVerification": "preferred",
		},
		"timeout":     60000,
		"attestation": "none",
	}))
}

// SetupWebAuthn completes WebAuthn registration
// @Description Completes WebAuthn registration with credential
// @Summary Setup WebAuthn
// @Tags mfa
// @Param request body validators.SetupWebAuthnValidator true "Request Data"
// @Success 201 {object} responses.SuccessResponse[models.MultiAuthMethod]
// @Failure 400,401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/webauthn/register [post]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) SetupWebAuthn(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	validatorData := new(validators.SetupWebAuthnValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := h.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	// Clear challenge cookie
	context.ClearCookie("webauthn_challenge")

	method, err := h.multiAuthService.SetupWebAuthn(session.User.ID, validatorData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.CreateSuccessResponse(
			fiber.StatusCreated,
			"WebAuthn credential has been registered successfully",
			method,
		),
	)
}

// GetWebAuthnAuthenticationOptions returns options for WebAuthn authentication
// @Description Returns options for WebAuthn authentication ceremony
// @Summary Get WebAuthn authentication options
// @Tags mfa
// @Success 200 {object} responses.SuccessResponse[fiber.Map]
// @Failure 401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/webauthn/authenticate/options [get]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) GetWebAuthnAuthenticationOptions(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	challenge, err := h.multiAuthService.GenerateWebAuthnChallenge()
	if err != nil {
		return err
	}

	credentials, err := h.multiAuthService.GetWebAuthnCredentials(session.User.ID)
	if err != nil {
		return err
	}

	if len(credentials) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no WebAuthn credentials registered")
	}

	// Store challenge in cookie for verification
	context.Cookie(&fiber.Cookie{
		Name:     "webauthn_challenge",
		Value:    challenge,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   300, // 5 minutes
	})

	allowCredentials := make([]fiber.Map, 0)
	for _, cred := range credentials {
		// Parse credentials JSON to get credential ID
		allowCredentials = append(allowCredentials, fiber.Map{
			"type": "public-key",
			"id":   cred.ID, // This should be the credential ID from the JSON
		})
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, fiber.Map{
		"challenge":        challenge,
		"rpId":             context.Hostname(),
		"allowCredentials": allowCredentials,
		"userVerification": "preferred",
		"timeout":          60000,
	}))
}

// VerifyWebAuthn verifies a WebAuthn authentication assertion
// @Description Verifies a WebAuthn authentication assertion
// @Summary Verify WebAuthn
// @Tags mfa
// @Param request body validators.VerifyWebAuthnValidator true "Request Data"
// @Success 200 {object} responses.SuccessResponse[responses.GenericMessage]
// @Failure 400,401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/webauthn/authenticate [post]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) VerifyWebAuthn(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	validatorData := new(validators.VerifyWebAuthnValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := h.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	// Clear challenge cookie
	context.ClearCookie("webauthn_challenge")

	if err := h.multiAuthService.VerifyWebAuthn(session.User.ID, validatorData.CredentialID, validatorData.SignCount); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "WebAuthn authentication successful",
	}))
}

// GetWebAuthnCredentials returns all WebAuthn credentials for the authenticated user
// @Description Returns all WebAuthn credentials for the authenticated user
// @Summary Get WebAuthn credentials
// @Tags mfa
// @Success 200 {object} responses.SuccessCollectionResponse[models.MultiAuthMethod]
// @Failure 401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/mfa/webauthn/credentials [get]
// @Accept json
// @Produce json
func (h *MultiAuthHandler) GetWebAuthnCredentials(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	credentials, err := h.multiAuthService.GetWebAuthnCredentials(session.User.ID)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, credentials))
}

// Helper function to get public key credential parameters
func getPublicKeyCredentialParams() []fiber.Map {
	return []fiber.Map{
		{"type": "public-key", "alg": -7},   // ES256
		{"type": "public-key", "alg": -257}, // RS256
	}
}
