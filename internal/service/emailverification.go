package service

import (
	"errors"
	"time"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
)

type EmailVerificationService struct {
	emailVerificationRepo *repository.EmailVerificationRepository
	userRepo              *repository.UserRepository
	authSettingsRepo      *repository.AuthSettingsRepository
}

func NewEmailVerificationService(
	emailVerificationRepo *repository.EmailVerificationRepository,
	userRepo *repository.UserRepository,
	authSettingsRepo *repository.AuthSettingsRepository,
) *EmailVerificationService {
	return &EmailVerificationService{
		emailVerificationRepo: emailVerificationRepo,
		userRepo:              userRepo,
		authSettingsRepo:      authSettingsRepo,
	}
}

// CreateVerificationToken creates a new email verification token for a user
func (s *EmailVerificationService) CreateVerificationToken(userID string) (*models.EmailVerificationToken, error) {
	// Check if user exists
	user, err := s.userRepo.FindByID(userID,
		repository.WithSelect("id", "email", "email_verified"),
	)
	if err != nil {
		return nil, err
	}

	if user.EmailVerified {
		return nil, errors.New("email is already verified")
	}

	// Delete any existing tokens for this user
	s.emailVerificationRepo.DeleteByUserID(userID)

	// Create new token (valid for 24 hours)
	token := &models.EmailVerificationToken{
		UserID:    userID,
		ExpiresOn: time.Now().Add(24 * time.Hour),
	}

	if err := s.emailVerificationRepo.Create(token, repository.WithPreload("User")); err != nil {
		return nil, err
	}

	return token, nil
}

// VerifyEmail verifies the user's email using the token
func (s *EmailVerificationService) VerifyEmail(tokenID string) error {
	token, err := s.emailVerificationRepo.FindByID(tokenID, repository.WithPreload("User"))
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	// Check if token is expired
	if time.Now().After(token.ExpiresOn) {
		// Delete expired token
		s.emailVerificationRepo.Delete(tokenID)
		return errors.New("verification token has expired")
	}

	// Update user's email_verified status
	token.User.EmailVerified = true
	now := time.Now()
	token.User.UpdatedAt = &now

	if err := s.userRepo.Update(&token.User,
		repository.WithSelect("email_verified", "updated_at"),
	); err != nil {
		return err
	}

	// Delete the used token
	return s.emailVerificationRepo.Delete(tokenID)
}

// ResendVerificationToken creates a new verification token for a user
func (s *EmailVerificationService) ResendVerificationToken(userID string) (*models.EmailVerificationToken, error) {
	return s.CreateVerificationToken(userID)
}

// GetPendingVerification gets the pending verification token for a user
func (s *EmailVerificationService) GetPendingVerification(userID string) (*models.EmailVerificationToken, error) {
	token, err := s.emailVerificationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(token.ExpiresOn) {
		s.emailVerificationRepo.Delete(token.ID)
		return nil, errors.New("verification token has expired")
	}

	return token, nil
}

// IsEmailVerificationRequired checks if email verification is required based on settings
func (s *EmailVerificationService) IsEmailVerificationRequired() (bool, error) {
	authSettings, err := s.authSettingsRepo.Find()
	if err != nil {
		return false, err
	}
	return authSettings.RequireEmailVerification, nil
}

// CleanupExpiredTokens removes all expired verification tokens
func (s *EmailVerificationService) CleanupExpiredTokens() error {
	return s.emailVerificationRepo.DeleteExpired()
}
