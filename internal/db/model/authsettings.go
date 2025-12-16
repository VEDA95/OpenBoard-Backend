package models

import (
	"time"

	"gorm.io/gorm"
)

type AuthSettings struct {
	gorm.Model
	ID                       uint       `gorm:"primaryKey;default:1;constraint:check:check_single_row_auth,id = 1" json:"id"`
	UpdatedAt                *time.Time `json:"updated_at"`
	AllowPublicRegistration  bool       `gorm:"default:true" json:"allow_public_registration"`
	RequireEmailVerification bool       `gorm:"default:true" json:"require_email_verification"`
	RequireAdminApproval     bool       `gorm:"default:false" json:"require_admin_approval"`
	DefaultUserRole          string     `gorm:"type:varchar(50);default:'user'" json:"default_user_role"`
	RegistrationWelcomeEmail bool       `gorm:"default:true" json:"registration_welcome_email"`
	AllowUserInvitations     bool       `gorm:"default:true" json:"allow_user_invitations"`
	InvitationExpiry         int        `gorm:"default:7" json:"invitation_expiry"`
	InviteOnlyMode           bool       `gorm:"default:false" json:"invite_only_mode"`
	SessionTimeout           int        `gorm:"default:10080" json:"session_timeout"`
	SessionIdleTimeout       int        `gorm:"default:60" json:"session_idle_timeout"`
	RefreshTokenExpiry       int        `gorm:"default:30" json:"refresh_token_expiry"`
	RememberMeDuration       int        `gorm:"default:30" json:"remember_me_duration"`
	MaxLoginAttempts         int        `gorm:"default:5" json:"max_login_attempts"`
	LockoutDuration          int        `gorm:"default:30" json:"lockout_duration"`
	TwoFactorAuthentication  bool       `gorm:"default:false" json:"two_factor_authentication"`
	TwoFactorRequired        bool       `gorm:"default:false" json:"two_factor_required"`
	EnableOAuth              bool       `gorm:"default:false" json:"enable_oauth"`
}

func (a *AuthSettings) BeforeCreate(tx *gorm.DB) error {
	a.ID = 1
	return nil
}
