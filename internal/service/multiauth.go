package service

import (
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	models "VEDA95/open_board/api/internal/db/model"

	"github.com/pquerna/otp/totp"
	"gorm.io/datatypes"
)

const (
	MFATypeEmailTOTP = "email_totp"
	MFATypeWebAuthn  = "webauthn"
	TOTPCodeLength   = 6
	TOTPValidityMin  = 10
)

type TOTPCredentials struct {
	LastCodeSent time.Time `json:"last_code_sent"`
	CodeHash     string    `json:"code_hash,omitempty"`
	CodeExpiry   time.Time `json:"code_expiry,omitempty"`
}

type WebAuthnCredentials struct {
	CredentialID     string `json:"credential_id"`
	PublicKey        string `json:"public_key"`
	AttestationType  string `json:"attestation_type"`
	AAGUID           string `json:"aaguid"`
	SignCount        uint32 `json:"sign_count"`
	CloneWarning     bool   `json:"clone_warning"`
	Attachment       string `json:"attachment"`
	RegisteredAt     string `json:"registered_at"`
	LastUsedAt       string `json:"last_used_at,omitempty"`
	TransportsString string `json:"transports,omitempty"`
}

// PendingTOTPChallenge represents a pending TOTP verification
type PendingTOTPChallenge struct {
	Code      string
	ExpiresAt time.Time
	UserID    string
	Email     string
}

type MultiAuthService struct {
	multiAuthRepo     *repository.MultiAuthMethodRepository
	userRepo          *repository.UserRepository
	authSettingsRepo  *repository.AuthSettingsRepository
	pendingChallenges map[string]*PendingTOTPChallenge
}

func NewMultiAuthService(
	multiAuthRepo *repository.MultiAuthMethodRepository,
	userRepo *repository.UserRepository,
	authSettingsRepo *repository.AuthSettingsRepository,
) *MultiAuthService {
	return &MultiAuthService{
		multiAuthRepo:     multiAuthRepo,
		userRepo:          userRepo,
		authSettingsRepo:  authSettingsRepo,
		pendingChallenges: make(map[string]*PendingTOTPChallenge),
	}
}

func (s *MultiAuthService) GetUserMFAMethods(userID string) ([]*models.MultiAuthMethod, error) {
	return s.multiAuthRepo.FindByUserID(userID, repository.QueryOptions{
		Omit: []string{"Credentials"},
	})
}

func (s *MultiAuthService) GetMFAMethodByID(ID string) (*models.MultiAuthMethod, error) {
	return s.multiAuthRepo.FindByID(ID, repository.QueryOptions{})
}

func (s *MultiAuthService) HasMFAEnabled(userID string) bool {
	return s.multiAuthRepo.CountByUserID(userID) > 0
}

func (s *MultiAuthService) IsMFARequired() (bool, error) {
	authSettings, err := s.authSettingsRepo.Find()
	if err != nil {
		return false, err
	}
	return authSettings.TwoFactorRequired, nil
}

// IsMFAEnabled checks if MFA is enabled in settings
func (s *MultiAuthService) IsMFAEnabled() (bool, error) {
	authSettings, err := s.authSettingsRepo.Find()
	if err != nil {
		return false, err
	}
	return authSettings.TwoFactorAuthentication, nil
}

// --- Email TOTP Methods ---

// SetupEmailTOTP sets up email-based TOTP for a user
func (s *MultiAuthService) SetupEmailTOTP(userID string, name string) (*models.MultiAuthMethod, error) {
	// Check if user already has email TOTP
	if s.multiAuthRepo.ExistsByUserIDAndType(userID, MFATypeEmailTOTP) {
		return nil, errors.New("email TOTP is already set up for this user")
	}

	credentials := TOTPCredentials{
		LastCodeSent: time.Time{},
	}

	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	method := &models.MultiAuthMethod{
		Name:        name,
		Type:        MFATypeEmailTOTP,
		UserID:      userID,
		Credentials: datatypes.JSON(credentialsJSON),
	}

	if err := s.multiAuthRepo.Create(method, repository.QueryOptions{}); err != nil {
		return nil, err
	}

	return method, nil
}

// GenerateEmailTOTPCode generates a TOTP code and stores it for verification
func (s *MultiAuthService) GenerateEmailTOTPCode(userID string) (*PendingTOTPChallenge, error) {
	// Verify user has email TOTP set up
	_, err := s.multiAuthRepo.FindByUserIDAndType(userID, MFATypeEmailTOTP, repository.QueryOptions{})
	if err != nil {
		return nil, errors.New("email TOTP is not set up for this user")
	}

	user, err := s.userRepo.FindByID(userID, repository.QueryOptions{
		Select: []string{"id", "email"},
	})
	if err != nil {
		return nil, err
	}

	// Generate random 6-digit code
	now := time.Now()
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "open_board",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}

	token, err := totp.GenerateCode(key.Secret(), now)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// Create challenge
	challenge := &PendingTOTPChallenge{
		Code:      token,
		ExpiresAt: time.Now().Add(time.Duration(TOTPValidityMin) * time.Minute),
		UserID:    userID,
		Email:     user.Email,
	}

	// Store challenge (keyed by userID)
	s.pendingChallenges[userID] = challenge

	return challenge, nil
}

// VerifyEmailTOTPCode verifies the TOTP code sent via email
func (s *MultiAuthService) VerifyEmailTOTPCode(userID string, code string) error {
	challenge, exists := s.pendingChallenges[userID]
	if !exists {
		return errors.New("no pending verification code")
	}

	// Check expiry
	if time.Now().After(challenge.ExpiresAt) {
		delete(s.pendingChallenges, userID)
		return errors.New("verification code has expired")
	}

	// Check code
	if challenge.Code != code {
		return errors.New("invalid verification code")
	}

	// Remove used challenge
	delete(s.pendingChallenges, userID)

	return nil
}

// --- WebAuthn Methods ---

// SetupWebAuthn initiates WebAuthn registration for a user
func (s *MultiAuthService) SetupWebAuthn(userID string, data *validators.SetupWebAuthnValidator) (*models.MultiAuthMethod, error) {
	// Check if name is already used
	if s.multiAuthRepo.ExistsByUserIDAndName(userID, data.Name) {
		return nil, errors.New("a WebAuthn credential with this name already exists")
	}

	credentials := WebAuthnCredentials{
		CredentialID:    data.CredentialID,
		PublicKey:       data.PublicKey,
		AttestationType: data.AttestationType,
		AAGUID:          data.AAGUID,
		SignCount:       0,
		CloneWarning:    false,
		Attachment:      data.Attachment,
		RegisteredAt:    time.Now().Format(time.RFC3339),
	}

	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	method := &models.MultiAuthMethod{
		Name:        data.Name,
		Type:        MFATypeWebAuthn,
		UserID:      userID,
		Credentials: datatypes.JSON(credentialsJSON),
	}

	if err := s.multiAuthRepo.Create(method, repository.QueryOptions{}); err != nil {
		return nil, err
	}

	return method, nil
}

// GetWebAuthnCredentials returns all WebAuthn credentials for a user
func (s *MultiAuthService) GetWebAuthnCredentials(userID string) ([]*models.MultiAuthMethod, error) {
	methods, err := s.multiAuthRepo.FindByUserID(userID, repository.QueryOptions{})
	if err != nil {
		return nil, err
	}

	// Filter for WebAuthn only
	webAuthnMethods := make([]*models.MultiAuthMethod, 0)
	for _, m := range methods {
		if m.Type == MFATypeWebAuthn {
			webAuthnMethods = append(webAuthnMethods, m)
		}
	}

	return webAuthnMethods, nil
}

// VerifyWebAuthn verifies a WebAuthn assertion
func (s *MultiAuthService) VerifyWebAuthn(userID string, credentialID string, signCount uint32) error {
	// Find the credential
	methods, err := s.GetWebAuthnCredentials(userID)
	if err != nil {
		return err
	}

	for _, method := range methods {
		var creds WebAuthnCredentials
		if err := json.Unmarshal(method.Credentials, &creds); err != nil {
			continue
		}

		if creds.CredentialID == credentialID {
			// Verify sign count to detect cloned authenticators
			if signCount <= creds.SignCount {
				creds.CloneWarning = true
			}

			// Update sign count and last used
			creds.SignCount = signCount
			creds.LastUsedAt = time.Now().Format(time.RFC3339)

			credentialsJSON, err := json.Marshal(creds)
			if err != nil {
				return err
			}

			method.Credentials = datatypes.JSON(credentialsJSON)
			return s.multiAuthRepo.Update(method, repository.QueryOptions{})
		}
	}

	return errors.New("credential not found")
}

// --- Common Methods ---

// DeleteMFAMethod deletes an MFA method
func (s *MultiAuthService) DeleteMFAMethod(ID string, userID string) error {
	method, err := s.multiAuthRepo.FindByID(ID, repository.QueryOptions{})
	if err != nil {
		return err
	}

	// Verify ownership
	if method.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.multiAuthRepo.Delete(ID)
}

// DeleteAllUserMFAMethods deletes all MFA methods for a user
func (s *MultiAuthService) DeleteAllUserMFAMethods(userID string) error {
	return s.multiAuthRepo.DeleteByUserID(userID)
}

// Helper function to generate a random numeric code
func generateTOTPCode(length int) (string, error) {
	const digits = "0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = digits[int(b[i])%len(digits)]
	}

	return string(code), nil
}

// GenerateWebAuthnChallenge generates a random challenge for WebAuthn
func (s *MultiAuthService) GenerateWebAuthnChallenge() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GetUserForWebAuthn returns user data needed for WebAuthn ceremonies
func (s *MultiAuthService) GetUserForWebAuthn(userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID, repository.QueryOptions{
		Select: []string{"id", "username", "email"},
	})
	if err != nil {
		return nil, err
	}

	// Get existing credentials
	credentials, err := s.GetWebAuthnCredentials(userID)
	if err != nil {
		return nil, err
	}

	excludeCredentials := make([]map[string]interface{}, 0)
	for _, cred := range credentials {
		var webAuthnCred WebAuthnCredentials
		if err := json.Unmarshal(cred.Credentials, &webAuthnCred); err != nil {
			continue
		}
		excludeCredentials = append(excludeCredentials, map[string]interface{}{
			"type": "public-key",
			"id":   webAuthnCred.CredentialID,
		})
	}

	return map[string]interface{}{
		"id":                 user.ID,
		"name":               user.Username,
		"displayName":        fmt.Sprintf("%s (%s)", user.Username, user.Email),
		"excludeCredentials": excludeCredentials,
	}, nil
}
