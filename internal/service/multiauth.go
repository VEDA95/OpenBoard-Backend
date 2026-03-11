package service

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"

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
	mu                sync.RWMutex
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
	return s.multiAuthRepo.FindByUserID(userID, repository.WithOmit("Credentials"))
}

func (s *MultiAuthService) GetMFAMethodByID(ID string) (*models.MultiAuthMethod, error) {
	return s.multiAuthRepo.FindByID(ID)
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

	if err := s.multiAuthRepo.Create(method); err != nil {
		return nil, err
	}

	return method, nil
}

// GenerateEmailTOTPCode generates a TOTP code and stores it for verification
func (s *MultiAuthService) GenerateEmailTOTPCode(userID string) (*PendingTOTPChallenge, error) {
	// Verify user has email TOTP set up
	_, err := s.multiAuthRepo.FindByUserIDAndType(userID, MFATypeEmailTOTP)
	if err != nil {
		return nil, errors.New("email TOTP is not set up for this user")
	}

	user, err := s.userRepo.FindByID(userID, repository.WithSelect("id", "email"))
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

	// Create challenge
	challenge := &PendingTOTPChallenge{
		Code:      token,
		ExpiresAt: time.Now().Add(time.Duration(TOTPValidityMin) * time.Minute),
		UserID:    userID,
		Email:     user.Email,
	}

	// Store challenge (keyed by userID)
	s.mu.Lock()
	s.pendingChallenges[userID] = challenge
	s.mu.Unlock()

	return challenge, nil
}

// VerifyEmailTOTPCode verifies the TOTP code sent via email
func (s *MultiAuthService) VerifyEmailTOTPCode(userID string, code string) error {
	s.mu.RLock()
	challenge, exists := s.pendingChallenges[userID]
	s.mu.RUnlock()
	if !exists {
		return errors.New("no pending verification code")
	}

	// Check expiry
	if time.Now().After(challenge.ExpiresAt) {
		s.mu.Lock()
		delete(s.pendingChallenges, userID)
		s.mu.Unlock()
		return errors.New("verification code has expired")
	}

	// Check code
	if challenge.Code != code {
		return errors.New("invalid verification code")
	}

	// Remove used challenge
	s.mu.Lock()
	delete(s.pendingChallenges, userID)
	s.mu.Unlock()

	return nil
}

// --- Common Methods ---

// DeleteMFAMethod deletes an MFA method
func (s *MultiAuthService) DeleteMFAMethod(ID string, userID string) error {
	method, err := s.multiAuthRepo.FindByID(ID)
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
