package validators

type SettingsTypeValidator struct {
	Type string `validate:"required,oneof=general auth email"`
}

type GeneralSettingsValidator struct {
	AppName                *string `json:"app_name,omitempty" validate:"optional"`
	AppURL                 *string `json:"app_url" validate:"optional,min=12"`
	AppLogoID              *string `json:"app_logo" validate:"optional,uuid"`
	AppFaviconID           *string `json:"app_favicon" validate:"optional,uuid"`
	AppDescription         *string `json:"app_description" validate:"optional"`
	ShowAnnouncementBanner *bool   `json:"show_announcement_banner" validate:"optional"`
	AnnouncementMessage    *string `json:"announcement_messange" validate:"optional"`
	AnnouncementType       *string `json:"announcement_type" validate:"optional,oneof=info success warning critical"`
	DefaultLanguage        *string `json:"default_language" validate:"optional"`
	DefaultTimezone        *string `json:"default_timezone" validate:"optional"`
	DefaultItemsPerPage    *int    `json:"default_items_per_page" validate:"optional,gt=0"`
	MaxFileSize            *int    `json:"max_file_size" validate:"optional,gt=0"`
}

type AuthSettingsValidator struct {
	AllowPublicRegistration  *bool   `json:"allow_public_registration" validate:"optional"`
	RequireEmailVerification *bool   `json:"require_email_verification" validate:"optional"`
	RequireAdminApproval     *bool   `json:"require_admin_approval" validate:"optional"`
	DefaultUserRole          *string `json:"default_user_role" validate:"optional,min=1"`
	RegistrationWelcomeEmail *bool   `json:"registration_welcome_email" validate:"optional"`
	AllowUserInvitations     *bool   `json:"allow_user_invitations" validate:"optional"`
	InvitationExpiry         *int    `json:"invitation_expiry" validate:"optional,gt=0"`
	InviteOnlyMode           *bool   `json:"invite_only_mode" validate:"optional"`
	SessionTimeout           *int    `json:"session_timeout" validate:"optional,gt=0"`
	SessionIdleTimeout       *int    `json:"session_idle_timeout" validate:"optional,gt=0"`
	RememberMeDuration       *int    `json:"remember_me_duration" validate:"optional,gt=0"`
	MaxLoginAttempts         *int    `json:"max_login_attempts" validate:"optional,gt=0"`
	LockoutDuration          *int    `json:"lockout_duration" validate:"gt=0"`
	TwoFactorAuthentication  *bool   `json:"two_factor_authentication" validate:"optional"`
	TwoFactorRequired        *bool   `json:"two_factor_required" validate:"optional"`
	EnableOAuth              *bool   `json:"enable_oauth" validate:"optional"`
}

type EmailSettingsValidator struct {
	EmailProvider            *string `json:"email_provider" validate:"optional,oneof=smtp sendgrid mailgun ses postmark"` // smtp, sendgrid, mailgun, ses, postmark
	EmailEnabled             *bool   `json:"email_enabled" validate:"optional"`
	EmailFromAddress         *string `json:"email_from_address" validate:"optional"`
	EmailFromName            *string `json:"email_from_name" validate:"optional"`
	SMTPHost                 *string `json:"smtp_host" validate:"optional"`
	SMTPPort                 *int    `json:"smtp_port" validate:"optional"`
	SMTPUsername             *string `json:"smtp_username" validate:"optional"`
	SMTPPassword             *string `json:"smtp_password" validate:"optional"`
	SMTPEncryption           *string `json:"smtp_encryption" validate:"optional,oneof=tls ssl none"`                   // tls, ssl, none
	SMTPAuthMethod           *string `json:"smtp_auth_method" validate:"optional,oneof=plain login, cram-md5, noauth"` // plain, login, cram-md5
	SMTPVerifySSL            *bool   `json:"smtp_verify_ssl" validate:"optional"`
	SendGridAPIKey           *string `json:"sendgrid_api_key" validate:"optional"`
	MailgunAPIKey            *string `json:"mailgun_api_key" validate:"optional"`
	MailgunDomain            *string `json:"mailgun_domain" validate:"optional"`
	SESAccessKeyID           *string `json:"ses_access_key_id" validate:"optional"`
	SESSecretAccessKey       *string `json:"ses_secret_access_key" validate:"optional"`
	SESRegion                *string `json:"ses_region" validate:"optional"`
	PostmarkServerToken      *string `json:"postmark_server_token" validate:"optional"`
	PostmarkAccountToken     *string `json:"postmark_account_token" validate:"optional"`
	EmailFooterText          *string `json:"email_footer_text" validate:"optional"`
	EnableEmailNotifications *bool   `json:"enable_email_notifications" validate:"optional"`
	SendPasswordResetEmail   *bool   `json:"send_password_reset_email" validate:"optional"`
	NotifyOnCardAssigned     *bool   `json:"notify_on_card_assigned" validate:"optional"`
	NotifyOnCardComment      *bool   `json:"notify_on_card_comment" validate:"optional"`
	NotifyOnCardDue          *bool   `json:"notify_on_card_due" validate:"optional"`
	NotifyOnBoardInvite      *bool   `json:"notify_on_board_invite" validate:"optional"`
	NotifyOnMention          *bool   `json:"notify_on_mention" validate:"optional"`
}
