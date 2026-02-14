package routes

import (
	"VEDA95/open_board/api/internal/auth"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"
	"fmt"
	"net/url"

	models "VEDA95/open_board/api/internal/db/model"

	"github.com/gofiber/fiber/v2"
)

var _ = models.ExternalAuthProvider{}

type ExternalProviderHandler struct {
	providerService *service.ExternalAuthProviderService
	settingsService *service.SettingsService
	validator       *validators.Validator
}

func NewExternalProviderHandler(
	providerService *service.ExternalAuthProviderService,
	settingsService *service.SettingsService,
	validator *validators.Validator,
) *ExternalProviderHandler {
	return &ExternalProviderHandler{
		providerService: providerService,
		settingsService: settingsService,
		validator:       validator,
	}
}

func (h *ExternalProviderHandler) getCookieDomain() string {
	authSettings, err := h.settingsService.GetAuthSettings()
	if err != nil || len(authSettings.CORSDomain) == 0 {
		return "localhost"
	}
	parsedURL, err := url.Parse(authSettings.CORSDomain)
	if err != nil {
		return "localhost"
	}
	return parsedURL.Hostname()
}

// GET returns all external auth providers
// @Description Returns all external auth providers
// @Summary Get all providers
// @Tags oauth
// @Success 200 {object} responses.OkCollectionResponse[models.ExternalAuthProvider]
// @Failure 500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/providers [get]
// @Accept json
// @Produce json
func (h *ExternalProviderHandler) GET(context *fiber.Ctx) error {
	providers, err := h.providerService.GetProviders()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, providers))
}

// GETEnabled returns all enabled external auth providers (for public display)
// @Description Returns all enabled external auth providers (for login page)
// @Summary Get enabled providers
// @Tags oauth
// @Success 200 {object} responses.OkCollectionResponse[models.ExternalAuthProvider]
// @Failure 500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/providers/enabled [get]
// @Accept json
// @Produce json
func (h *ExternalProviderHandler) GETEnabled(context *fiber.Ctx) error {
	providers, err := h.providerService.GetEnabledProviders()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, providers))
}

// GETByID returns an external auth provider by ID
// @Description Returns an external auth provider by ID
// @Summary Get provider by ID
// @Tags oauth
// @Param id path string true "Provider ID"
// @Success 200 {object} responses.OkResponse[models.ExternalAuthProvider]
// @Failure 404,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/providers/{id} [get]
// @Accept json
// @Produce json
func (h *ExternalProviderHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := h.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	provider, err := h.providerService.GetProviderByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			fmt.Sprintf("External provider: %s has been created successfully", provider.ID),
			provider,
		),
	)
}

// POST creates a new external auth provider
// @Description Creates a new external auth provider
// @Summary Create provider
// @Tags oauth
// @Param request body validators.CreateExternalProviderValidator true "Request Data"
// @Success 201 {object} responses.SuccessResponse[models.ExternalAuthProvider]
// @Failure 400,422,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/providers [post]
// @Accept json
// @Produce json
func (h *ExternalProviderHandler) POST(context *fiber.Ctx) error {
	validatorData := new(validators.CreateExternalProviderValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := h.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	provider, err := h.providerService.CreateProvider(validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.CreateSuccessResponse(
			fiber.StatusCreated,
			fmt.Sprintf("Provider %s has been created successfully", provider.Name),
			provider,
		),
	)
}

// PATCH updates an external auth provider
// @Description Updates an external auth provider
// @Summary Update provider
// @Tags oauth
// @Param id path string true "Provider ID"
// @Param request body validators.UpdateExternalProviderValidator true "Request Data"
// @Success 200 {object} responses.SuccessResponse[models.ExternalAuthProvider]
// @Failure 400,404,422,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/providers/{id} [patch]
// @Accept json
// @Produce json
func (h *ExternalProviderHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := h.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateExternalProviderValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := h.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	provider, err := h.providerService.UpdateProvider(params.Id, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.CreateSuccessResponse(
			fiber.StatusOK,
			fmt.Sprintf("Provider %s has been updated successfully", provider.Name),
			provider,
		),
	)
}

// DELETE deletes an external auth provider
// @Description Deletes an external auth provider
// @Summary Delete provider
// @Tags oauth
// @Param id path string true "Provider ID"
// @Success 200 {object} responses.SuccessResponse[responses.GenericMessage]
// @Failure 404,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/providers/{id} [delete]
// @Accept json
// @Produce json
func (h *ExternalProviderHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := h.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := h.providerService.DeleteProvider(params.Id); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: fmt.Sprintf("Provider %s has been deleted successfully", params.Id),
	}))
}

// Authorize initiates the OAuth flow for a provider
// @Description Initiates the OAuth authorization flow
// @Summary Start OAuth flow
// @Tags oauth
// @Param id path string true "Provider ID"
// @Param redirect_url query string true "Redirect URL after OAuth"
// @Success 302 "Redirect to OAuth provider"
// @Failure 400,404,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/authorize/{id} [get]
func (h *ExternalProviderHandler) Authorize(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := h.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	redirectURL := context.Query("redirect_url")
	if len(redirectURL) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "redirect_url is required")
	}

	provider, err := h.providerService.GetProviderByID(params.Id)
	if err != nil {
		return err
	}

	state, err := h.providerService.GenerateOAuthState()
	if err != nil {
		return err
	}

	// Store state in session/cookie for verification
	context.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   600, // 10 minutes
	})

	context.Cookie(&fiber.Cookie{
		Name:     "oauth_redirect",
		Value:    redirectURL,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   600,
	})

	authURL := h.providerService.GetAuthorizationURL(provider, state, redirectURL)

	return context.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// Callback handles the OAuth token exchange (called by frontend after OAuth redirect)
// @Description Exchanges OAuth authorization code for tokens and creates a session
// @Summary OAuth token exchange
// @Tags oauth
// @Param id path string true "Provider ID"
// @Param request body validators.OAuthTokenExchangeValidator true "Request Data"
// @Success 201 {object} responses.SuccessResponse[auth.LocalUserLogin]
// @Failure 400,401,422,500 {object} responses.ErrorResponse[responses.GenericMessage]
// @Router /auth/oauth/callback/{id} [post]
// @Accept json
// @Produce json
func (h *ExternalProviderHandler) Callback(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := h.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.OAuthTokenExchangeValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := h.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	// Verify state from cookie
	storedState := context.Cookies("oauth_state")
	if storedState != validatorData.State {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid state parameter")
	}

	// Clear OAuth state cookies
	context.ClearCookie("oauth_state")
	context.ClearCookie("oauth_redirect")

	provider, err := h.providerService.GetProviderByID(params.Id)
	if err != nil {
		return err
	}

	// Complete OAuth login: exchange code, fetch user info, create/find user, create session
	authData, err := h.providerService.OAuthLogin(
		context.Context(),
		provider,
		validatorData.Code,
		validatorData.RedirectURL,
		context.Get("User-Agent"),
		context.IP(),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Return session cookies or JSON tokens based on return type
	if validatorData.ReturnType == "session" {
		context.Status(fiber.StatusCreated)
		context.Cookie(&fiber.Cookie{
			Name:     "open_board_session",
			Value:    authData.Session.AccessToken,
			Expires:  authData.Session.ExpiresOn,
			HTTPOnly: true,
			Secure:   false,
			Path:     "/",
			Domain:   h.getCookieDomain(),
		})

		if authData.Session.RefreshToken != nil && authData.Session.RefreshExpiresOn != nil {
			context.Cookie(&fiber.Cookie{
				Name:     "open_board_session_remember_me",
				Value:    *authData.Session.RefreshToken,
				Expires:  *authData.Session.RefreshExpiresOn,
				HTTPOnly: true,
				Secure:   false,
				Path:     "/",
				Domain:   h.getCookieDomain(),
			})
		}

		return nil
	}

	// Return JSON tokens
	responseData := auth.LocalUserLogin{
		AccessToken: authData.Session.AccessToken,
		ExpiresIn:   authData.ExpiresIn,
		User:        authData.Session.User,
	}

	if authData.Session.RefreshToken != nil {
		responseData.RefreshToken = authData.Session.RefreshToken
		responseData.RefreshExpiresIn = authData.RefreshExpiresIn
	}

	return responses.JSONResponse(
		context,
		fiber.StatusCreated,
		responses.CreateSuccessResponse(
			fiber.StatusCreated,
			fmt.Sprintf("%s has been successfully logged in via %s", authData.Session.User.Username, provider.Name),
			responseData,
		),
	)
}
