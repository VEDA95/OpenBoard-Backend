package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
	"errors"
	"os"
	"strconv"
	"time"
)

type AuthService struct {
	userRepo          *repository.UserRepository
	sessionRepo       *repository.SessionRepository
	passwordResetRepo *repository.PasswordResetRepository
}

func NewSessionService() (*AuthService, error) {
	sessionRepo, err := repository.NewSessionRepository()
	if err != nil {
		return nil, err
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return nil, err
	}

	passwordResetRepo, err := repository.NewPasswordResetRepository()
	if err != nil {
		return nil, err
	}

	return &AuthService{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		passwordResetRepo: passwordResetRepo,
	}, nil
}

func (authService *AuthService) GetUserSessions(ID string) ([]*models.Session, error) {
	sessions, err := authService.sessionRepo.FindByUser(ID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (authService *AuthService) ValidateSession(token string) (*models.Session, error) {
	session, err := authService.sessionRepo.FindByAccessToken(token)
	if err != nil {
		return nil, err
	}

	if !session.IsValid() {
		return nil, errors.New("invalid credentials")
	}

	return session, nil
}

func (authService *AuthService) LocalLogin(data *validators.LocalLoginValidator, userAgent string, IPAdress string) (*models.Session, error) {
	expiresInEnv := os.Getenv("AUTH_SESSION_EXPIRES_IN")
	refreshExpiresInEnv := os.Getenv("AUTH_SESSION_REFRESH_EXPIRES_IN")

	if len(expiresInEnv) == 0 || len(refreshExpiresInEnv) == 0 {
		return nil, errors.New("AUTH_SESSION_EXPIRES_IN and/or AUTH_SESSION_REFRESH_EXPIRES_IN environment variable(s) was not set")
	}

	expiresIn, err := strconv.Atoi(expiresInEnv)
	if err != nil {
		return nil, err
	}

	refreshExpiresIn, err := strconv.Atoi(refreshExpiresInEnv)
	if err != nil {
		return nil, err
	}

	user, err := authService.userRepo.FindByUsername(data.Username)
	if err != nil {
		return nil, err
	}

	if !user.ValidPassword(data.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.Enabled {
		return nil, errors.New("account has been disabled")
	}

	now := time.Now()
	session := &models.Session{
		Type:      "local",
		IPAddress: IPAdress,
		UserAgent: userAgent,
		ExpiresOn: now.Add(time.Second * time.Duration(expiresIn)),
		UserID:    user.ID,
	}
	user.LastLogin = &now

	if data.Remember {
		refreshExpiresOn := now.Add(time.Second * time.Duration(refreshExpiresIn))
		session.RememberMe = true
		session.RefreshExpiresOn = &refreshExpiresOn
	}

	if err := session.GenerateTokens(); err != nil {
		return nil, err
	}

	if err := authService.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	if err := authService.userRepo.Update(user); err != nil {
		return nil, err
	}

	return session, nil
}

func (authService *AuthService) LocalRefresh(token string, remember bool) (*models.Session, error) {
	expiresInEnv := os.Getenv("AUTH_SESSION_EXPIRES_IN")
	refreshExpiresInEnv := os.Getenv("AUTH_SESSION_REFRESH_EXPIRES_IN")

	if len(expiresInEnv) == 0 || len(refreshExpiresInEnv) == 0 {
		return nil, errors.New("AUTH_SESSION_EXPIRES_IN and/or AUTH_SESSION_REFRESH_EXPIRES_IN environment variable(s) was not set")
	}

	expiresIn, err := strconv.Atoi(expiresInEnv)
	if err != nil {
		return nil, err
	}

	refreshExpiresIn, err := strconv.Atoi(refreshExpiresInEnv)
	if err != nil {
		return nil, err
	}

	session, err := authService.sessionRepo.FindByRefreshToken(token)
	if err != nil {
		return nil, err
	}

	if !session.IsRefreshValid() {
		return nil, errors.New("invalid credentials")
	}

	user, err := authService.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errors.New("account has been disabled")
	}

	now := time.Now()
	session.ExpiresOn = now.Add(time.Second * time.Duration(expiresIn))

	if remember {
		refreshExpiresOn := now.Add(time.Second * time.Duration(refreshExpiresIn))
		session.RememberMe = true
		session.RefreshExpiresOn = &refreshExpiresOn
	}

	if err := session.GenerateTokens(); err != nil {
		return nil, err
	}

	if err := authService.sessionRepo.Update(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (authService *AuthService) LocalLogout(token string) error {
	session, err := authService.sessionRepo.FindByAccessToken(token)
	if err != nil {
		return err
	}

	return authService.sessionRepo.Delete(session.ID)
}

func (authService *AuthService) LocalLogoutByID(ID string) error {
	return authService.sessionRepo.Delete(ID)
}

func (authService *AuthService) LocalLogoutByUserID(ID string) error {
	return authService.sessionRepo.DeleteByUserID(ID)
}

func (authService *AuthService) IssueForgotPasswordToken(email string) (*models.PasswordResetToken, error) {
	user, err := authService.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	passwordResetToken := &models.PasswordResetToken{
		Type:      "form",
		UserID:    user.ID,
		ExpiresOn: time.Now().Add(time.Minute * 15),
	}

	if err := passwordResetToken.GenerateToken(user.Email); err != nil {
		return nil, err
	}

	if err := authService.passwordResetRepo.Create(passwordResetToken); err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (authService *AuthService) IssueUserPasswordResetToken(user *models.User, data *validators.AuthenticatedPasswordResetValidator) (*models.PasswordResetToken, error) {
	if !user.ValidPassword(data.Password) {
		return nil, errors.New("invalid credentials")
	}

	passwordResetToken := &models.PasswordResetToken{
		Type:      "auth",
		UserID:    user.ID,
		ExpiresOn: time.Now().Add(time.Minute * 15),
	}

	if err := passwordResetToken.GenerateToken(user.Email); err != nil {
		return nil, err
	}

	if err := authService.passwordResetRepo.Create(passwordResetToken); err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (authService *AuthService) ResetForgottenPassword(data *validators.ResetPasswordValidator) error {
	passwordResetToken, err := authService.passwordResetRepo.FindByToken(data.Token)
	if err != nil {
		return err
	}

	if !passwordResetToken.IsValid() {
		return errors.New("invalid credentials")
	}

	if passwordResetToken.Type != "form" {
		return errors.New("invalid credentials")
	}

	user, err := authService.userRepo.FindByID(passwordResetToken.UserID)
	if err != nil {
		return err
	}

	if err := user.HashPassword(data.ConfirmPassword); err != nil {
		return err
	}

	if err := authService.passwordResetRepo.Delete(passwordResetToken.ID); err != nil {
		return err
	}

	if err := authService.sessionRepo.DeleteByUserID(user.ID); err != nil {
		return err
	}

	now := time.Now()
	user.UpdatedAt = &now

	return authService.userRepo.Update(user)
}

func (authService *AuthService) ResetUserPassword(user *models.User, data *validators.ResetPasswordValidator) error {
	passwordResetToken, err := authService.passwordResetRepo.FindByToken(data.Token)
	if err != nil {
		return err
	}

	if !passwordResetToken.IsValid() {
		return errors.New("invalid credentials")
	}

	if passwordResetToken.Type != "auth" {
		return errors.New("invalid credentials")
	}

	if err := user.HashPassword(data.ConfirmPassword); err != nil {
		return err
	}

	if err := authService.passwordResetRepo.Delete(passwordResetToken.ID); err != nil {
		return err
	}

	if err := authService.sessionRepo.DeleteByUserID(user.ID); err != nil {
		return err
	}

	now := time.Now()
	user.UpdatedAt = &now

	return authService.userRepo.Update(user)
}
