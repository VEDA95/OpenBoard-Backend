package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/datatypes"
)

// WebAuthnUser implements the webauthn.User interface
type WebAuthnUser struct {
	user        *models.User
	credentials []webauthn.Credential
}

func NewWebAuthnUser(user *models.User, credentials []webauthn.Credential) *WebAuthnUser {
	return &WebAuthnUser{
		user:        user,
		credentials: credentials,
	}
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return []byte(u.user.ID)
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.user.Username
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	if u.user.FirstName != nil && u.user.LastName != nil {
		return fmt.Sprintf("%s %s", *u.user.FirstName, *u.user.LastName)
	}
	return u.user.Username
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// WebAuthnService handles WebAuthn registration and authentication ceremonies
type WebAuthnService struct {
	webAuthn         *webauthn.WebAuthn
	multiAuthRepo    *repository.MultiAuthMethodRepository
	userRepo         *repository.UserRepository
	authSettingsRepo *repository.AuthSettingsRepository
	mu               sync.RWMutex
	sessions         map[string]*webauthn.SessionData
}

func NewWebAuthnService(
	multiAuthRepo *repository.MultiAuthMethodRepository,
	userRepo *repository.UserRepository,
	authSettingsRepo *repository.AuthSettingsRepository,
) (*WebAuthnService, error) {
	s := &WebAuthnService{
		multiAuthRepo:    multiAuthRepo,
		userRepo:         userRepo,
		authSettingsRepo: authSettingsRepo,
		sessions:         make(map[string]*webauthn.SessionData),
	}

	if err := s.loadConfig(); err != nil {
		return nil, err
	}

	return s, nil
}

// loadConfig reads WebAuthn configuration from auth settings and initializes the webauthn instance
func (s *WebAuthnService) loadConfig() error {
	authSettings, err := s.authSettingsRepo.Find()
	if err != nil {
		return fmt.Errorf("failed to read auth settings for webauthn: %w", err)
	}

	rpID := authSettings.WebAuthnRPID
	if rpID == "" {
		rpID = "localhost"
	}

	rpDisplayName := authSettings.WebAuthnRPDisplayName
	if rpDisplayName == "" {
		rpDisplayName = "Open Board"
	}

	var rpOrigins []string
	if authSettings.WebAuthnRPOrigins != "" {
		rpOrigins = strings.Split(authSettings.WebAuthnRPOrigins, ",")
	}
	if len(rpOrigins) == 0 {
		// Fall back to CORS domain if no explicit origins are set
		if authSettings.CORSDomain != "" {
			rpOrigins = []string{authSettings.CORSDomain}
		} else {
			rpOrigins = []string{"http://localhost:3000"}
		}
	}

	w, err := webauthn.New(&webauthn.Config{
		RPDisplayName: rpDisplayName,
		RPID:          rpID,
		RPOrigins:     rpOrigins,
	})
	if err != nil {
		return fmt.Errorf("failed to create webauthn instance: %w", err)
	}

	s.webAuthn = w
	return nil
}

// getWebAuthnUser loads a user and their WebAuthn credentials, returning a WebAuthnUser
func (s *WebAuthnService) getWebAuthnUser(userID string) (*WebAuthnUser, error) {
	user, err := s.userRepo.FindByID(userID,
		repository.WithSelect("id", "username", "email", "first_name", "last_name"),
	)
	if err != nil {
		return nil, err
	}

	methods, err := s.multiAuthRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var credentials []webauthn.Credential
	for _, method := range methods {
		if method.Type != MFATypeWebAuthn {
			continue
		}
		var cred webauthn.Credential
		if err := json.Unmarshal(method.Credentials, &cred); err != nil {
			continue
		}
		credentials = append(credentials, cred)
	}

	return NewWebAuthnUser(user, credentials), nil
}

// BeginRegistration starts the WebAuthn registration ceremony
func (s *WebAuthnService) BeginRegistration(userID string) (*protocol.CredentialCreation, error) {
	webAuthnUser, err := s.getWebAuthnUser(userID)
	if err != nil {
		return nil, err
	}

	options, session, err := s.webAuthn.BeginRegistration(webAuthnUser)
	if err != nil {
		return nil, fmt.Errorf("failed to begin registration: %w", err)
	}

	s.mu.Lock()
	s.sessions[userID] = session
	s.mu.Unlock()

	return options, nil
}

// FinishRegistration completes the WebAuthn registration ceremony and stores the credential
func (s *WebAuthnService) FinishRegistration(userID string, name string, response *protocol.ParsedCredentialCreationData) (*models.MultiAuthMethod, error) {
	webAuthnUser, err := s.getWebAuthnUser(userID)
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	session, exists := s.sessions[userID]
	s.mu.RUnlock()
	if !exists {
		return nil, errors.New("no pending registration session")
	}

	credential, err := s.webAuthn.CreateCredential(webAuthnUser, *session, response)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	s.mu.Lock()
	delete(s.sessions, userID)
	s.mu.Unlock()

	// Check if name is already used
	if s.multiAuthRepo.ExistsByUserIDAndName(userID, name) {
		return nil, errors.New("a WebAuthn credential with this name already exists")
	}

	credJSON, err := json.Marshal(credential)
	if err != nil {
		return nil, err
	}

	method := &models.MultiAuthMethod{
		Name:        name,
		Type:        MFATypeWebAuthn,
		UserID:      userID,
		Credentials: datatypes.JSON(credJSON),
	}

	if err := s.multiAuthRepo.Create(method); err != nil {
		return nil, err
	}

	return method, nil
}

// BeginLogin starts the WebAuthn authentication ceremony
func (s *WebAuthnService) BeginLogin(userID string) (*protocol.CredentialAssertion, error) {
	webAuthnUser, err := s.getWebAuthnUser(userID)
	if err != nil {
		return nil, err
	}

	if len(webAuthnUser.credentials) == 0 {
		return nil, errors.New("no WebAuthn credentials registered")
	}

	options, session, err := s.webAuthn.BeginLogin(webAuthnUser)
	if err != nil {
		return nil, fmt.Errorf("failed to begin login: %w", err)
	}

	s.mu.Lock()
	s.sessions[userID] = session
	s.mu.Unlock()

	return options, nil
}

// FinishLogin completes the WebAuthn authentication ceremony
func (s *WebAuthnService) FinishLogin(userID string, response *protocol.ParsedCredentialAssertionData) error {
	webAuthnUser, err := s.getWebAuthnUser(userID)
	if err != nil {
		return err
	}

	s.mu.RLock()
	session, exists := s.sessions[userID]
	s.mu.RUnlock()
	if !exists {
		return errors.New("no pending authentication session")
	}

	credential, err := s.webAuthn.ValidateLogin(webAuthnUser, *session, response)
	if err != nil {
		return fmt.Errorf("failed to validate login: %w", err)
	}

	s.mu.Lock()
	delete(s.sessions, userID)
	s.mu.Unlock()

	// Update the stored credential with new sign count
	return s.updateCredentialAfterLogin(userID, credential)
}

// updateCredentialAfterLogin updates the stored credential's sign count after successful authentication
func (s *WebAuthnService) updateCredentialAfterLogin(userID string, credential *webauthn.Credential) error {
	methods, err := s.multiAuthRepo.FindByUserID(userID)
	if err != nil {
		return err
	}

	for _, method := range methods {
		if method.Type != MFATypeWebAuthn {
			continue
		}

		var storedCred webauthn.Credential
		if err := json.Unmarshal(method.Credentials, &storedCred); err != nil {
			continue
		}

		if string(storedCred.ID) == string(credential.ID) {
			storedCred.Authenticator.SignCount = credential.Authenticator.SignCount
			storedCred.Authenticator.CloneWarning = credential.Authenticator.CloneWarning

			credJSON, err := json.Marshal(storedCred)
			if err != nil {
				return err
			}

			method.Credentials = datatypes.JSON(credJSON)
			return s.multiAuthRepo.Update(method)
		}
	}

	return nil
}

// GetCredentials returns all WebAuthn credentials for a user
func (s *WebAuthnService) GetCredentials(userID string) ([]*models.MultiAuthMethod, error) {
	methods, err := s.multiAuthRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	webAuthnMethods := make([]*models.MultiAuthMethod, 0)
	for _, m := range methods {
		if m.Type == MFATypeWebAuthn {
			webAuthnMethods = append(webAuthnMethods, m)
		}
	}

	return webAuthnMethods, nil
}
