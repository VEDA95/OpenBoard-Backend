package service

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
	"errors"
)

type SettingsService struct {
	generalSettingsRepo *repository.GeneralSettingsRepository
	authSettingsRepo    *repository.AuthSettingsRepository
	emailSettingsRepo   *repository.EmailSettingsRepository
}

func NewSettingsService(
	generalSettingsRepo *repository.GeneralSettingsRepository,
	authSettingsRepo *repository.AuthSettingsRepository,
	emailSettingsRepo *repository.EmailSettingsRepository,
) *SettingsService {
	return &SettingsService{
		generalSettingsRepo: generalSettingsRepo,
		authSettingsRepo:    authSettingsRepo,
		emailSettingsRepo:   emailSettingsRepo,
	}
}

func (settingsService *SettingsService) GetGeneralSettings() (*models.GeneralSettings, error) {
	return settingsService.generalSettingsRepo.Find()
}

func (settingsService *SettingsService) GetAuthSettings() (*models.AuthSettings, error) {
	return settingsService.authSettingsRepo.Find()
}

func (settingsService *SettingsService) GetEmailSettings() (*models.EmailSettings, error) {
	return settingsService.emailSettingsRepo.Find()
}

func (settingsService *SettingsService) UpdateGeneralSettings(data *validators.GeneralSettingsValidator) (*models.GeneralSettings, error) {
	generalSettings, err := settingsService.generalSettingsRepo.Find()
	if err != nil {
		return nil, err
	}

	if generalSettings == nil {
		return nil, errors.New("general settings have not been initialized")
	}

	if data.AppName != nil && *data.AppName != generalSettings.AppName {
		generalSettings.AppName = *data.AppName
	}

	if data.AppURL != nil && *data.AppURL != generalSettings.AppURL {
		generalSettings.AppURL = *data.AppURL
	}

	if data.AppLogo != nil && data.AppLogo != generalSettings.AppLogo {
		if len(*data.AppLogo) == 0 {
			generalSettings.AppLogo = nil
		} else {
			generalSettings.AppLogo = data.AppLogo
		}
	}

	if data.AppFavicon != nil && data.AppFavicon != generalSettings.AppFavicon {
		if len(*data.AppFavicon) == 0 {
			generalSettings.AppFavicon = nil
		} else {
			generalSettings.AppFavicon = data.AppFavicon
		}
	}

	if data.AppDescription != nil && data.AppDescription != generalSettings.AppDescription {
		if len(*data.AppDescription) == 0 {
			generalSettings.AppDescription = nil
		} else {
			generalSettings.AppDescription = data.AppDescription
		}
	}

	if data.ShowAnnouncementBanner != nil && *data.ShowAnnouncementBanner != generalSettings.ShowAnnouncementBanner {
		generalSettings.ShowAnnouncementBanner = *data.ShowAnnouncementBanner
	}

	if data.AnnouncementMessage != nil && data.AnnouncementMessage != generalSettings.AnnouncementMessage {
		if len(*data.AnnouncementMessage) == 0 {
			generalSettings.AnnouncementMessage = nil
		} else {
			generalSettings.AnnouncementMessage = data.AnnouncementMessage
		}
	}

	if data.AnnouncementType != nil && data.AnnouncementType != generalSettings.AnnouncementType {
		if len(*data.AnnouncementType) == 0 {
			generalSettings.AnnouncementType = nil
		} else {
			generalSettings.AnnouncementType = data.AnnouncementType
		}
	}

	if data.DefaultItemsPerPage != nil && *data.DefaultItemsPerPage != generalSettings.DefaultItemsPerPage {
		generalSettings.DefaultItemsPerPage = *data.DefaultItemsPerPage
	}

	if data.DefaultLanguage != nil && *data.DefaultLanguage != generalSettings.DefaultLanguage {
		generalSettings.DefaultLanguage = *data.DefaultLanguage
	}

	if data.DefaultTimezone != nil && *data.DefaultTimezone != generalSettings.DefaultTimezone {
		generalSettings.DefaultTimezone = *data.DefaultTimezone
	}

	if err := settingsService.generalSettingsRepo.Update(generalSettings); err != nil {
		return nil, err
	}

	return generalSettings, nil
}

func (settingsService *SettingsService) UpdateAuthSettings(data *validators.AuthSettingsValidator) (*models.AuthSettings, error) {
	authSettings, err := settingsService.authSettingsRepo.Find()
	if err != nil {
		return nil, err
	}

	if authSettings == nil {
		return nil, errors.New("auth settings have not been initialized")
	}

	if data.AllowPublicRegistration != nil && *data.AllowPublicRegistration != authSettings.AllowPublicRegistration {
		authSettings.AllowPublicRegistration = *data.AllowPublicRegistration
	}

	if data.AllowUserInvitations != nil && *data.AllowUserInvitations != authSettings.AllowUserInvitations {
		authSettings.AllowUserInvitations = *data.AllowUserInvitations
	}

	if data.InviteOnlyMode != nil && *data.InviteOnlyMode != authSettings.InviteOnlyMode {
		authSettings.InviteOnlyMode = *data.InviteOnlyMode
	}

	if data.RequireAdminApproval != nil && *data.RequireAdminApproval != authSettings.RequireAdminApproval {
		authSettings.RequireAdminApproval = *data.RequireAdminApproval
	}

	if data.RequireEmailVerification != nil && *data.RequireEmailVerification != authSettings.RequireEmailVerification {
		authSettings.RequireEmailVerification = *data.RequireEmailVerification
	}

	if data.TwoFactorAuthentication != nil && *data.TwoFactorAuthentication != authSettings.TwoFactorAuthentication {
		authSettings.TwoFactorAuthentication = *data.TwoFactorAuthentication
	}

	if data.TwoFactorRequired != nil && *data.TwoFactorRequired != authSettings.TwoFactorRequired {
		authSettings.TwoFactorRequired = *data.TwoFactorRequired
	}

	if data.EnableOAuth != nil && *data.EnableOAuth != authSettings.EnableOAuth {
		authSettings.EnableOAuth = *data.EnableOAuth
	}

	if data.RegistrationWelcomeEmail != nil && *data.RegistrationWelcomeEmail != authSettings.RegistrationWelcomeEmail {
		authSettings.RegistrationWelcomeEmail = *data.RegistrationWelcomeEmail
	}

	if data.DefaultUserRole != nil && *data.DefaultUserRole != authSettings.DefaultUserRole {
		authSettings.DefaultUserRole = *data.DefaultUserRole
	}

	if data.InvitationExpiry != nil && *data.InvitationExpiry != authSettings.InvitationExpiry {
		authSettings.InvitationExpiry = *data.InvitationExpiry
	}

	if data.SessionTimeout != nil && *data.SessionTimeout != authSettings.SessionTimeout {
		authSettings.SessionTimeout = *data.SessionTimeout
	}

	if data.SessionIdleTimeout != nil && *data.SessionIdleTimeout != authSettings.SessionIdleTimeout {
		authSettings.SessionIdleTimeout = *data.SessionIdleTimeout
	}

	if data.RememberMeDuration != nil && *data.RememberMeDuration != authSettings.RememberMeDuration {
		authSettings.RememberMeDuration = *data.RememberMeDuration
	}

	if data.MaxLoginAttempts != nil && *data.MaxLoginAttempts != authSettings.MaxLoginAttempts {
		authSettings.MaxLoginAttempts = *data.MaxLoginAttempts
	}

	if data.LockoutDuration != nil && *data.LockoutDuration != authSettings.LockoutDuration {
		authSettings.LockoutDuration = *data.LockoutDuration
	}

	if err := settingsService.authSettingsRepo.Update(authSettings); err != nil {
		return nil, err
	}

	return authSettings, nil
}

func (settingsService *SettingsService) UpdateEmailSettings(data *validators.EmailSettingsValidator) (*models.EmailSettings, error) {
	emailSettings, err := settingsService.emailSettingsRepo.Find()
	if err != nil {
		return nil, err
	}

	if emailSettings == nil {
		return nil, errors.New("email settings have not been initialized")
	}

	if data.EmailEnabled != nil && *data.EmailEnabled != emailSettings.EmailEnabled {
		emailSettings.EmailEnabled = *data.EmailEnabled
	}

	if data.EnableEmailNotifications != nil && *data.EnableEmailNotifications != emailSettings.EnableEmailNotifications {
		emailSettings.EnableEmailNotifications = *data.EnableEmailNotifications
	}

	if data.NotifyOnBoardInvite != nil && *data.NotifyOnBoardInvite != emailSettings.NotifyOnBoardInvite {
		emailSettings.NotifyOnBoardInvite = *data.NotifyOnBoardInvite
	}

	if data.NotifyOnCardAssigned != nil && *data.NotifyOnCardAssigned != emailSettings.NotifyOnCardAssigned {
		emailSettings.NotifyOnCardAssigned = *data.NotifyOnCardAssigned
	}

	if data.NotifyOnCardComment != nil && *data.NotifyOnCardComment != emailSettings.NotifyOnCardComment {
		emailSettings.NotifyOnCardComment = *data.NotifyOnCardComment
	}

	if data.NotifyOnCardDue != nil && *data.NotifyOnCardDue != emailSettings.NotifyOnCardDue {
		emailSettings.NotifyOnCardDue = *data.NotifyOnCardDue
	}

	if data.NotifyOnMention != nil && *data.NotifyOnMention != emailSettings.NotifyOnMention {
		emailSettings.NotifyOnMention = *data.NotifyOnMention
	}

	if data.SMTPVerifySSL != nil && *data.SMTPVerifySSL != emailSettings.SMTPVerifySSL {
		emailSettings.SMTPVerifySSL = *data.SMTPVerifySSL
	}

	if data.SendPasswordResetEmail != nil && *data.SendPasswordResetEmail != emailSettings.SendPasswordResetEmail {
		emailSettings.SendPasswordResetEmail = *data.SendPasswordResetEmail
	}

	if data.EmailProvider != nil && *data.EmailProvider != emailSettings.EmailProvider {
		emailSettings.EmailProvider = *data.EmailProvider
	}

	if data.EmailFromAddress != nil && data.EmailFromAddress != emailSettings.EmailFromAddress {
		emailSettings.EmailFromAddress = data.EmailFromAddress
	}

	if data.EmailFromName != nil && data.EmailFromName != emailSettings.EmailFromName {
		emailSettings.EmailFromName = data.EmailFromName
	}

	if data.EmailReplyTo != nil && data.EmailReplyTo != emailSettings.EmailReplyTo {
		emailSettings.EmailReplyTo = data.EmailReplyTo
	}

	if data.SMTPHost != nil && data.SMTPHost != emailSettings.SMTPHost {
		emailSettings.SMTPHost = data.SMTPHost
	}

	if data.SMTPPort != nil && *data.SMTPPort != emailSettings.SMTPPort {
		emailSettings.SMTPPort = *data.SMTPPort
	}

	if data.SMTPUsername != nil && data.SMTPUsername != emailSettings.SMTPUsername {
		emailSettings.SMTPUsername = data.SMTPUsername
	}

	if data.SMTPPassword != nil && data.SMTPPassword != emailSettings.SMTPPassword {
		emailSettings.SMTPPassword = data.SMTPPassword
	}

	if data.SMTPAuthMethod != nil && data.SMTPAuthMethod != emailSettings.SMTPAuthMethod {
		emailSettings.SMTPAuthMethod = data.SMTPAuthMethod
	}

	if data.SMTPEncryption != nil && data.SMTPEncryption != emailSettings.SMTPEncryption {
		emailSettings.SMTPEncryption = data.SMTPEncryption
	}

	if data.SMTPPoolSize != nil && *data.SMTPPoolSize != emailSettings.SMTPPoolSize {
		emailSettings.SMTPPoolSize = *data.SMTPPoolSize
	}

	if data.SMTPTimeout != nil && *data.SMTPTimeout != emailSettings.SMTPTimeout {
		emailSettings.SMTPTimeout = *data.SMTPTimeout
	}

	if data.SendGridAPIKey != nil && data.SendGridAPIKey != emailSettings.SendGridAPIKey {
		emailSettings.SendGridAPIKey = data.SendGridAPIKey
	}

	if data.SendGridWebhookSecret != nil && data.SendGridWebhookSecret != emailSettings.SendGridWebhookSecret {
		emailSettings.SendGridWebhookSecret = data.SendGridWebhookSecret
	}

	if data.MailgunAPIKey != nil && data.MailgunAPIKey != emailSettings.MailgunAPIKey {
		emailSettings.MailgunAPIKey = data.MailgunAPIKey
	}

	if data.MailgunDomain != nil && data.MailgunDomain != emailSettings.MailgunDomain {
		emailSettings.MailgunDomain = data.MailgunDomain
	}

	if data.MailgunRegion != nil && *data.MailgunRegion != emailSettings.MailgunRegion {
		emailSettings.MailgunRegion = *data.MailgunRegion
	}

	if data.SESAccessKeyID != nil && data.SESAccessKeyID != emailSettings.SESAccessKeyID {
		emailSettings.SESAccessKeyID = data.SESAccessKeyID
	}

	if data.SESSecretAccessKey != nil && data.SESSecretAccessKey != emailSettings.SESSecretAccessKey {
		emailSettings.SESSecretAccessKey = data.SESSecretAccessKey
	}

	if data.SESRegion != nil && *data.SESRegion != emailSettings.SESRegion {
		emailSettings.SESRegion = *data.SESRegion
	}

	if data.EmailFooterText != nil && data.EmailFooterText != emailSettings.EmailFooterText {
		emailSettings.EmailFooterText = data.EmailFooterText
	}

	if err := settingsService.emailSettingsRepo.Update(emailSettings); err != nil {
		return nil, err
	}

	return emailSettings, nil
}
