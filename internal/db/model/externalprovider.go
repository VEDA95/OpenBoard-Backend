package models

type ExternalAuthProvider struct {
	Base
	Name                    string        `gorm:"type:varchar(255);not null" json:"name"`
	ClientID                string        `gorm:"type:varchar(255); not null" json:"client_id"`
	ClientSecret            string        `gorm:"type:text;not null" json:"client_secret"`
	RequiredEmailDomain     *string       `gorm:"type:varchar(255)" json:"required_email_domain"`
	AuthURL                 string        `gorm:"type:varchar(255); not null" json:"auth_url"`
	LoginURL                string        `gorm:"type:varchar(255); not null" json:"login_url"`
	UserInfoURL             string        `gorm:"type:varchar(255); not null" json:"userinfo_url"`
	LogoutURL               string        `gorm:"type:varchar(255);" json:"logout_url"`
	UsePKCE                 bool          `gorm:"not null; default:false" json:"use_pkce"`
	DefaultLoginMethod      bool          `gorm:"not null; default:false" json:"default_login_method"`
	SelfRegistrationEnabled bool          `gorm:"not null; default:false" json:"self_registration_enabled"`
	Permissions             []*Permission `gorm:"many2many:external_provider_permissions" json:"permissions"`
}
