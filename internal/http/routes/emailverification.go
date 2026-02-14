package routes

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type EmailVerificationHandler struct {
	emailVerificationService *service.EmailVerificationService
	emailService             *service.EmailService
	validator                *validators.Validator
}

func NewEmailVerificationHandler(
	emailVerificationService *service.EmailVerificationService,
	emailService *service.EmailService,
	validator *validators.Validator,
) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		emailVerificationService: emailVerificationService,
		emailService:             emailService,
		validator:                validator,
	}
}

// SendVerificationEmail sends a verification email to the authenticated user
// @Description Sends a verification email to the authenticated user
// @Summary Send verification email
// @Tags email-verification
// @Success 200 {object} responses.OkResponse[responses.GenericMessage]
// @Failure 400,401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/email/verify/send [post]
// @Accept json
// @Produce json
func (h *EmailVerificationHandler) SendVerificationEmail(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	if session.User.EmailVerified {
		return fiber.NewError(fiber.StatusBadRequest, "email is already verified")
	}

	token, err := h.emailVerificationService.CreateVerificationToken(session.User.ID)
	if err != nil {
		return err
	}

	// Send email with verification link
	if h.emailService != nil {
		subject := "Verify your email address"
		err := h.emailService.SendMessage(
			session.User.Email,
			"email_verification", // Template name - to be created later
			map[string]interface{}{
				"token_id": token.ID,
				"user":     session.User,
			},
			&subject,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to send verification email")
		}
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "Verification email has been sent",
	}))
}

// VerifyEmail verifies the user's email using the token
// @Description Verifies the user's email using the token
// @Summary Verify email
// @Tags email-verification
// @Param id path string true "Verification token ID"
// @Success 200 {object} responses.OkResponse[responses.GenericMessage]
// @Failure 400,404,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/email/verify/{id} [get]
// @Accept json
// @Produce json
func (h *EmailVerificationHandler) VerifyEmail(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := h.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := h.emailVerificationService.VerifyEmail(params.Id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "Email has been verified successfully",
	}))
}

// ResendVerificationEmail resends the verification email
// @Description Resends the verification email to the authenticated user
// @Summary Resend verification email
// @Tags email-verification
// @Success 200 {object} responses.OkResponse[responses.GenericMessage]
// @Failure 400,401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/email/verify/resend [post]
// @Accept json
// @Produce json
func (h *EmailVerificationHandler) ResendVerificationEmail(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	if session.User.EmailVerified {
		return fiber.NewError(fiber.StatusBadRequest, "email is already verified")
	}

	token, err := h.emailVerificationService.ResendVerificationToken(session.User.ID)
	if err != nil {
		return err
	}

	// Send email with verification link
	if h.emailService != nil {
		subject := "Verify your email address"
		err := h.emailService.SendMessage(
			session.User.Email,
			"email_verification", // Template name - to be created later
			map[string]interface{}{
				"token_id": token.ID,
				"user":     session.User,
			},
			&subject,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to send verification email")
		}
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "Verification email has been resent",
	}))
}

// GetVerificationStatus returns the email verification status for the authenticated user
// @Description Returns the email verification status
// @Summary Get verification status
// @Tags email-verification
// @Success 200 {object} responses.OkResponse[fiber.Map]
// @Failure 401,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/email/verify/status [get]
// @Accept json
// @Produce json
func (h *EmailVerificationHandler) GetVerificationStatus(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)

	required, err := h.emailVerificationService.IsEmailVerificationRequired()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, fiber.Map{
		"email_verified": session.User.EmailVerified,
		"required":       required,
	}))
}
