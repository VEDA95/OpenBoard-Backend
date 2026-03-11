package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"VEDA95/open_board/api/internal/auth"
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type ExternalAuthProviderService struct {
	providerRepo     *repository.ExternalAuthProviderRepository
	userRepo         *repository.UserRepository
	sessionRepo      *repository.SessionRepository
	roleRepo         *repository.RoleRepository
	authSettingsRepo *repository.AuthSettingsRepository
}

func NewExternalAuthProviderService(
	providerRepo *repository.ExternalAuthProviderRepository,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	roleRepo *repository.RoleRepository,
	authSettingsRepo *repository.AuthSettingsRepository,
) *ExternalAuthProviderService {
	return &ExternalAuthProviderService{
		providerRepo:     providerRepo,
		userRepo:         userRepo,
		sessionRepo:      sessionRepo,
		roleRepo:         roleRepo,
		authSettingsRepo: authSettingsRepo,
	}
}

// OAuthStateData holds the state for OAuth flow
type OAuthStateData struct {
	State       string
	Provider    string
	RedirectURL string
	CreatedAt   time.Time
}

// OAuthUserInfo represents user info from OAuth provider
type OAuthUserInfo struct {
	ID            string
	Email         string
	EmailVerified bool
	Name          string
	FirstName     string
	LastName      string
	Picture       string
}

// GetProviders returns all external auth providers
func (s *ExternalAuthProviderService) GetProviders() ([]*models.ExternalAuthProvider, error) {
	return s.providerRepo.FindAll(repository.WithPreload("Permissions"))
}

// GetEnabledProviders returns all enabled external auth providers
func (s *ExternalAuthProviderService) GetEnabledProviders() ([]*models.ExternalAuthProvider, error) {
	return s.providerRepo.FindEnabled(repository.WithOmit("ClientSecret"))
}

// GetProviderByID returns a provider by ID
func (s *ExternalAuthProviderService) GetProviderByID(ID string) (*models.ExternalAuthProvider, error) {
	return s.providerRepo.FindByID(ID, repository.WithPreload("Permissions"))
}

// GetProviderByName returns a provider by name
func (s *ExternalAuthProviderService) GetProviderByName(name string) (*models.ExternalAuthProvider, error) {
	return s.providerRepo.FindByName(name, repository.WithPreload("Permissions"))
}

// CreateProvider creates a new external auth provider
func (s *ExternalAuthProviderService) CreateProvider(data *validators.CreateExternalProviderValidator) (*models.ExternalAuthProvider, error) {
	if s.providerRepo.ExistsByName(data.Name) {
		return nil, errors.New("provider with this name already exists")
	}

	provider := &models.ExternalAuthProvider{
		Name:                    data.Name,
		ClientID:                data.ClientID,
		ClientSecret:            data.ClientSecret,
		AuthURL:                 data.AuthURL,
		LoginURL:                data.LoginURL,
		UserInfoURL:             data.UserInfoURL,
		LogoutURL:               data.LogoutURL,
		UsePKCE:                 data.UsePKCE,
		DefaultLoginMethod:      data.DefaultLoginMethod,
		SelfRegistrationEnabled: data.SelfRegistrationEnabled,
	}

	if data.RequiredEmailDomain != nil && len(*data.RequiredEmailDomain) > 0 {
		provider.RequiredEmailDomain = data.RequiredEmailDomain
	}

	if data.PermissionIDs != nil && len(*data.PermissionIDs) > 0 {
		err := s.providerRepo.CreateWithPermissions(
			provider,
			*data.PermissionIDs,
			repository.WithPreload("Permissions"),
		)
		if err != nil {
			return nil, err
		}
		return provider, nil
	}

	err := s.providerRepo.Create(provider, repository.WithPreload("Permissions"))
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// UpdateProvider updates an existing external auth provider
func (s *ExternalAuthProviderService) UpdateProvider(ID string, data *validators.UpdateExternalProviderValidator) (*models.ExternalAuthProvider, error) {
	provider, err := s.providerRepo.FindByID(ID)
	if err != nil {
		return nil, err
	}

	if data.Name != nil && *data.Name != provider.Name {
		if s.providerRepo.ExistsByName(*data.Name) {
			return nil, errors.New("provider with this name already exists")
		}
		provider.Name = *data.Name
	}

	if data.ClientID != nil {
		provider.ClientID = *data.ClientID
	}

	if data.ClientSecret != nil {
		provider.ClientSecret = *data.ClientSecret
	}

	if data.RequiredEmailDomain != nil {
		if len(*data.RequiredEmailDomain) == 0 {
			provider.RequiredEmailDomain = nil
		} else {
			provider.RequiredEmailDomain = data.RequiredEmailDomain
		}
	}

	if data.AuthURL != nil {
		provider.AuthURL = *data.AuthURL
	}

	if data.LoginURL != nil {
		provider.LoginURL = *data.LoginURL
	}

	if data.UserInfoURL != nil {
		provider.UserInfoURL = *data.UserInfoURL
	}

	if data.LogoutURL != nil {
		provider.LogoutURL = *data.LogoutURL
	}

	if data.UsePKCE != nil {
		provider.UsePKCE = *data.UsePKCE
	}

	if data.DefaultLoginMethod != nil {
		provider.DefaultLoginMethod = *data.DefaultLoginMethod
	}

	if data.SelfRegistrationEnabled != nil {
		provider.SelfRegistrationEnabled = *data.SelfRegistrationEnabled
	}

	if data.PermissionIDs != nil {
		err := s.providerRepo.UpdateWithPermissions(
			provider,
			*data.PermissionIDs,
			repository.WithPreload("Permissions"),
		)
		if err != nil {
			return nil, err
		}
		return provider, nil
	}

	err = s.providerRepo.Update(provider, repository.WithPreload("Permissions"))
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// DeleteProvider deletes an external auth provider
func (s *ExternalAuthProviderService) DeleteProvider(ID string) error {
	if !s.providerRepo.Exists(ID) {
		return errors.New("provider does not exist")
	}
	return s.providerRepo.Delete(ID)
}

// GenerateOAuthState generates a random state string for OAuth flow
func (s *ExternalAuthProviderService) GenerateOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetOAuthConfig creates an OAuth2 config for the provider
func (s *ExternalAuthProviderService) GetOAuthConfig(provider *models.ExternalAuthProvider, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthURL,
			TokenURL: provider.LoginURL,
		},
		RedirectURL: redirectURL,
		Scopes:      []string{"openid", "email", "profile"},
	}
}

// ExchangeCode exchanges the authorization code for tokens
func (s *ExternalAuthProviderService) ExchangeCode(ctx context.Context, provider *models.ExternalAuthProvider, code string, redirectURL string) (*oauth2.Token, error) {
	config := s.GetOAuthConfig(provider, redirectURL)
	return config.Exchange(ctx, code)
}

// GetAuthorizationURL returns the OAuth authorization URL
func (s *ExternalAuthProviderService) GetAuthorizationURL(provider *models.ExternalAuthProvider, state string, redirectURL string) string {
	config := s.GetOAuthConfig(provider, redirectURL)
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// FetchUserInfo fetches user information from the OAuth provider
func (s *ExternalAuthProviderService) FetchUserInfo(ctx context.Context, provider *models.ExternalAuthProvider, token *oauth2.Token) (*OAuthUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get(provider.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	// Parse response based on common OAuth provider formats
	var rawInfo map[string]interface{}
	if err := json.Unmarshal(body, &rawInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	userInfo := &OAuthUserInfo{}

	// Extract ID (various field names used by different providers)
	if id, ok := rawInfo["sub"].(string); ok {
		userInfo.ID = id
	} else if id, ok := rawInfo["id"].(string); ok {
		userInfo.ID = id
	} else if id, ok := rawInfo["id"].(float64); ok {
		userInfo.ID = fmt.Sprintf("%.0f", id)
	}

	// Extract email
	if email, ok := rawInfo["email"].(string); ok {
		userInfo.Email = email
	}

	// Extract email verified
	if verified, ok := rawInfo["email_verified"].(bool); ok {
		userInfo.EmailVerified = verified
	}

	// Extract name fields
	if name, ok := rawInfo["name"].(string); ok {
		userInfo.Name = name
	}
	if firstName, ok := rawInfo["given_name"].(string); ok {
		userInfo.FirstName = firstName
	}
	if lastName, ok := rawInfo["family_name"].(string); ok {
		userInfo.LastName = lastName
	}

	// If no first/last name but we have full name, try to split it
	if userInfo.FirstName == "" && userInfo.LastName == "" && userInfo.Name != "" {
		parts := strings.SplitN(userInfo.Name, " ", 2)
		if len(parts) >= 1 {
			userInfo.FirstName = parts[0]
		}
		if len(parts) >= 2 {
			userInfo.LastName = parts[1]
		}
	}

	// Extract picture/avatar
	if picture, ok := rawInfo["picture"].(string); ok {
		userInfo.Picture = picture
	} else if picture, ok := rawInfo["avatar_url"].(string); ok {
		userInfo.Picture = picture
	}

	if userInfo.Email == "" {
		return nil, errors.New("email not provided by OAuth provider")
	}

	return userInfo, nil
}

// OAuthLogin handles the complete OAuth login flow: exchange code, fetch user info, create/find user, create session
func (s *ExternalAuthProviderService) OAuthLogin(
	ctx context.Context,
	provider *models.ExternalAuthProvider,
	code string,
	redirectURL string,
	userAgent string,
	ipAddress string,
) (*auth.LoginResponse, error) {
	// Exchange code for token
	token, err := s.ExchangeCode(ctx, provider, code, redirectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Fetch user info from OAuth provider
	userInfo, err := s.FetchUserInfo(ctx, provider, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Check required email domain if configured
	if provider.RequiredEmailDomain != nil && len(*provider.RequiredEmailDomain) > 0 {
		emailParts := strings.Split(userInfo.Email, "@")
		if len(emailParts) != 2 || emailParts[1] != *provider.RequiredEmailDomain {
			return nil, fmt.Errorf("email domain %s is not allowed for this provider", emailParts[1])
		}
	}

	// Get auth settings
	authSettings, err := s.authSettingsRepo.Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if authSettings == nil {
		return nil, errors.New("unable to fetch auth settings")
	}

	// Find or create user
	user, err := s.userRepo.FindByEmail(userInfo.Email, repository.UserFullLoad...)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// User not found, check if self-registration is enabled
		if !provider.SelfRegistrationEnabled {
			return nil, errors.New("user not found and self-registration is not enabled for this provider")
		}

		// Create new user
		userRole, err := s.roleRepo.FindByName(authSettings.DefaultUserRole, repository.WithSelect("id"))
		if err != nil {
			return nil, fmt.Errorf("failed to get default role: %w", err)
		}

		// Generate a username from email
		username := strings.Split(userInfo.Email, "@")[0]
		// Ensure username is unique
		baseUsername := username
		counter := 1
		for s.userRepo.ExistsByUsername(username) {
			username = fmt.Sprintf("%s%d", baseUsername, counter)
			counter++
		}

		user = &models.User{
			Username:      username,
			Email:         userInfo.Email,
			EmailVerified: userInfo.EmailVerified,
			Enabled:       true,
		}

		if userInfo.FirstName != "" {
			user.FirstName = &userInfo.FirstName
		}
		if userInfo.LastName != "" {
			user.LastName = &userInfo.LastName
		}

		// Generate a random password for OAuth users (they won't use it)
		randomPassword := make([]byte, 32)
		if _, err := rand.Read(randomPassword); err != nil {
			return nil, err
		}
		if err := user.HashPassword(base64.URLEncoding.EncodeToString(randomPassword)); err != nil {
			return nil, err
		}

		err = s.userRepo.CreateWithRoles(
			user,
			[]string{userRole.ID},
			nil,
			repository.UserFullLoad...,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Check if user is enabled
	if !user.Enabled {
		return nil, errors.New("account has been disabled")
	}

	// Create session
	now := time.Now()
	session := &models.Session{
		Type:      "oauth",
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresOn: now.Add(time.Second * time.Duration(authSettings.SessionTimeout)),
		UserID:    user.ID,
	}

	// Enable remember me for OAuth logins by default
	refreshExpiresOn := now.Add(time.Second * time.Duration(authSettings.RememberMeDuration))
	session.RememberMe = true
	session.RefreshExpiresOn = &refreshExpiresOn

	if err := session.GenerateTokens(); err != nil {
		return nil, err
	}

	err = s.sessionRepo.Create(session, repository.WithSelect("id", "type", "ip_address", "user_agent", "expires_on", "user_id", "access_token", "remember_me", "refresh_expires_on", "refresh_token"))
	if err != nil {
		return nil, err
	}

	// Update last login
	user.LastLogin = &now
	if err := s.userRepo.Update(user, repository.WithSelect("last_login")); err != nil {
		return nil, err
	}

	session.User = user

	return &auth.LoginResponse{
		Session:          session,
		ExpiresIn:        authSettings.SessionTimeout,
		RefreshExpiresIn: &authSettings.RememberMeDuration,
	}, nil
}
