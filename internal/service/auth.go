package service

import (
	"VEDA95/open_board/api/internal/auth"
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo          *repository.UserRepository
	sessionRepo       *repository.SessionRepository
	passwordResetRepo *repository.PasswordResetRepository
	roleRepo          *repository.RoleRepository
	authSettingsRepo  *repository.AuthSettingsRepository
}

func NewAuthService(
	sessionRepo *repository.SessionRepository,
	userRepo *repository.UserRepository,
	passwordResetRepo *repository.PasswordResetRepository,
	roleRepo *repository.RoleRepository,
	authSettingsRepo *repository.AuthSettingsRepository,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		passwordResetRepo: passwordResetRepo,
		roleRepo:          roleRepo,
		authSettingsRepo:  authSettingsRepo,
	}
}

func (authService *AuthService) GetUserSessions(ID string) ([]*models.Session, error) {
	sessions, err := authService.sessionRepo.FindByUser(ID, repository.QueryOptions{
		Preload: []string{"User", "User.Roles", "User.Roles.Permissions"},
		Omit:    []string{"User.Sessions"},
	})
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (authService *AuthService) ValidateSession(token string) (*models.Session, error) {
	session, err := authService.sessionRepo.FindByAccessToken(token, repository.QueryOptions{
		Preload: []string{"User", "User.Roles", "User.Roles.Permissions", "User.Sessions"},
	})
	if err != nil {
		return nil, err
	}

	isValid := session.IsValid()
	isRefreshValid := session.IsRefreshValid()

	if !isValid && !isRefreshValid {
		if err := authService.sessionRepo.Delete(session.ID); err != nil {
			return nil, err
		}

		return nil, errors.New("invalid credentials")
	}

	if !isValid && isRefreshValid {
		return session, errors.New("refresh required")
	}

	return session, nil
}

func (authService *AuthService) LocalLogin(data *validators.LocalLoginValidator, userAgent string, IPAdress string) (*auth.LoginResponse, error) {
	authSettings, err := authService.authSettingsRepo.Find()

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if authSettings == nil {
		return nil, errors.New("unable to fetch auth settings")
	}

	user, err := authService.userRepo.FindByUsername(data.Username, repository.QueryOptions{
		Preload: []string{"Roles", "Roles.Permissions"},
		Omit:    []string{"Sessions"},
	})
	if err != nil {
		return nil, err
	}

	if !user.ValidPassword(data.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.Enabled {
		return nil, errors.New("account has been disabled")
	}

	response := &auth.LoginResponse{ExpiresIn: authSettings.SessionTimeout}
	columns := []string{"id", "type", "ip_address", "user_agent", "expires_on", "user_id", "access_token"}
	now := time.Now()
	session := &models.Session{
		Type:      "local",
		IPAddress: IPAdress,
		UserAgent: userAgent,
		ExpiresOn: now.Add(time.Second * time.Duration(authSettings.SessionTimeout)),
		UserID:    user.ID,
	}
	user.LastLogin = &now

	if data.Remember {
		refreshExpiresOn := now.Add(time.Second * time.Duration(authSettings.RememberMeDuration))
		session.RememberMe = true
		session.RefreshExpiresOn = &refreshExpiresOn
		response.RefreshExpiresIn = &authSettings.RememberMeDuration
		columns = append(columns, "remember_me", "refresh_expires_on", "refresh_token")
	}

	if err := session.GenerateTokens(); err != nil {
		return nil, err
	}

	err2 := authService.sessionRepo.Create(session, repository.QueryOptions{
		Select: columns,
	})
	if err2 != nil {
		return nil, err2
	}

	err3 := authService.userRepo.Update(user, repository.QueryOptions{
		Select: []string{"last_login"},
	})
	if err3 != nil {
		return nil, err3
	}

	session.User = user
	response.Session = session

	return response, nil
}

func (authService *AuthService) LocalRefresh(token string) (*auth.LoginResponse, error) {
	authSettings, err := authService.authSettingsRepo.Find()

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if authSettings == nil {
		return nil, errors.New("unable to fetch auth settings")
	}

	session, err := authService.sessionRepo.FindByRefreshToken(token, repository.QueryOptions{
		Preload: []string{"User", "User.Roles", "User.Roles.Permissions"},
		Omit:    []string{"User.Sessions"},
	})
	if err != nil {
		return nil, err
	}

	if !session.IsRefreshValid() {
		return nil, errors.New("invalid credentials")
	}

	if !session.User.Enabled {
		return nil, errors.New("account has been disabled")
	}

	columns := []string{"expires_on", "access_token"}
	response := &auth.LoginResponse{ExpiresIn: authSettings.SessionTimeout}
	now := time.Now()
	session.ExpiresOn = now.Add(time.Second * time.Duration(authSettings.SessionTimeout))

	if session.RememberMe {
		refreshExpiresOn := now.Add(time.Second * time.Duration(authSettings.RememberMeDuration))
		session.RefreshExpiresOn = &refreshExpiresOn
		response.RefreshExpiresIn = &authSettings.RememberMeDuration
		columns = append(columns, "refresh_expires_on", "refresh_token")
	}

	if err := session.GenerateTokens(); err != nil {
		return nil, err
	}

	err2 := authService.sessionRepo.Update(session, repository.QueryOptions{
		Select: columns,
	})
	if err2 != nil {
		return nil, err2
	}

	response.Session = session

	return response, nil
}

func (authService *AuthService) LocalLogout(token string) error {
	session, err := authService.sessionRepo.FindByAccessToken(token, repository.QueryOptions{
		Select: []string{"id"},
	})
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
	user, err := authService.userRepo.FindByEmail(email, repository.QueryOptions{
		Select: []string{"id", "email"},
	})
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

	err2 := authService.passwordResetRepo.Create(passwordResetToken, repository.QueryOptions{
		Omit: []string{"User"},
	})
	if err2 != nil {
		return nil, err2
	}

	passwordResetToken.User = *user

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

	err := authService.passwordResetRepo.Create(passwordResetToken, repository.QueryOptions{
		Omit: []string{"User"},
	})
	if err != nil {
		return nil, err
	}

	return passwordResetToken, nil
}

func (authService *AuthService) ResetForgottenPassword(data *validators.ResetPasswordValidator) error {
	passwordResetToken, err := authService.passwordResetRepo.FindByToken(data.Token, repository.QueryOptions{
		Preload: []string{"User"},
		Omit:    []string{"User.Roles", "User.Roles.Permissions", "User.Sessions"},
	})
	if err != nil {
		return err
	}

	if !passwordResetToken.IsValid() {
		return errors.New("invalid credentials")
	}

	if passwordResetToken.Type != "form" {
		return errors.New("invalid credentials")
	}

	if err := passwordResetToken.User.HashPassword(data.ConfirmPassword); err != nil {
		return err
	}

	if err := authService.passwordResetRepo.Delete(passwordResetToken.ID); err != nil {
		return err
	}

	if err := authService.sessionRepo.DeleteByUserID(passwordResetToken.User.ID); err != nil {
		return err
	}

	now := time.Now()
	passwordResetToken.User.UpdatedAt = &now

	return authService.userRepo.Update(&passwordResetToken.User, repository.QueryOptions{
		Select: []string{"updated_at", "hashed_password"},
	})
}

func (authService *AuthService) ResetUserPassword(user *models.User, data *validators.ResetPasswordValidator) error {
	passwordResetToken, err := authService.passwordResetRepo.FindByToken(data.Token, repository.QueryOptions{
		Omit: []string{"User", "User.Roles", "User.Roles.Permissions", "User.Sessions"},
	})
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

	return authService.userRepo.Update(user, repository.QueryOptions{
		Select: []string{"updated_at", "hashed_password"},
	})
}

func (authService *AuthService) RegisterUser(data *validators.RegisterUserValidator) error {
	authSettings, err := authService.authSettingsRepo.Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if authSettings == nil {
		return errors.New("unable to fetch auth settings")
	}

	if !authSettings.AllowPublicRegistration {
		return errors.New("user registration is not enabled")
	}

	if authService.userRepo.ExistsByUsernameOrEmail(data.Username, data.Email) {
		return errors.New("user already exists")
	}

	userRole, err := authService.roleRepo.FindByName(authSettings.DefaultUserRole, repository.QueryOptions{
		Select: []string{"id"},
	})
	if err != nil {
		return err
	}

	user := &models.User{
		Username: data.Username,
		Email:    data.Email,
	}

	if err := user.HashPassword(data.Password); err != nil {
		return err
	}

	if data.FirstName != nil && len(*data.FirstName) > 0 {
		user.FirstName = data.FirstName
	}

	if data.LastName != nil && len(*data.LastName) > 0 {
		user.LastName = data.LastName
	}

	return authService.userRepo.CreateWithRoles(
		user,
		[]string{userRole.ID},
		repository.QueryOptions{
			Preload: []string{"Roles", "Roles.Permissions"},
			Omit:    []string{"Sessions"},
		},
		repository.QueryOptions{},
	)
}
