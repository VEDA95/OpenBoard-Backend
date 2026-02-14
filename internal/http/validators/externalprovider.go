package validators

type CreateExternalProviderValidator struct {
	Name                    string    `json:"name" validate:"required,min=1,max=255"`
	ClientID                string    `json:"client_id" validate:"required,min=1"`
	ClientSecret            string    `json:"client_secret" validate:"required,min=1"`
	RequiredEmailDomain     *string   `json:"required_email_domain,omitempty" validate:"omitempty,min=1,max=255"`
	AuthURL                 string    `json:"auth_url" validate:"required,url"`
	LoginURL                string    `json:"login_url" validate:"required,url"`
	UserInfoURL             string    `json:"userinfo_url" validate:"required,url"`
	LogoutURL               string    `json:"logout_url,omitempty" validate:"omitempty,url"`
	UsePKCE                 bool      `json:"use_pkce" validate:"omitempty" default:"false"`
	DefaultLoginMethod      bool      `json:"default_login_method" validate:"omitempty" default:"false"`
	SelfRegistrationEnabled bool      `json:"self_registration_enabled" validate:"omitempty" default:"false"`
	PermissionIDs           *[]string `json:"permission_ids,omitempty" validate:"omitempty"`
}

type UpdateExternalProviderValidator struct {
	Name                    *string   `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	ClientID                *string   `json:"client_id,omitempty" validate:"omitempty,min=1"`
	ClientSecret            *string   `json:"client_secret,omitempty" validate:"omitempty,min=1"`
	RequiredEmailDomain     *string   `json:"required_email_domain,omitempty" validate:"omitempty"`
	AuthURL                 *string   `json:"auth_url,omitempty" validate:"omitempty,url"`
	LoginURL                *string   `json:"login_url,omitempty" validate:"omitempty,url"`
	UserInfoURL             *string   `json:"userinfo_url,omitempty" validate:"omitempty,url"`
	LogoutURL               *string   `json:"logout_url,omitempty" validate:"omitempty"`
	UsePKCE                 *bool     `json:"use_pkce,omitempty" validate:"omitempty"`
	DefaultLoginMethod      *bool     `json:"default_login_method,omitempty" validate:"omitempty"`
	SelfRegistrationEnabled *bool     `json:"self_registration_enabled,omitempty" validate:"omitempty"`
	PermissionIDs           *[]string `json:"permission_ids,omitempty" validate:"omitempty"`
}

type OAuthAuthorizeValidator struct {
	ProviderID  string `json:"provider_id" validate:"required,uuid"`
	RedirectURL string `json:"redirect_url" validate:"required,url"`
}

type OAuthCallbackValidator struct {
	ProviderID  string `json:"provider_id" validate:"required,uuid"`
	Code        string `json:"code" validate:"required,min=1"`
	State       string `json:"state" validate:"required,min=1"`
	RedirectURL string `json:"redirect_url" validate:"required,url"`
}

type OAuthTokenExchangeValidator struct {
	Code        string `json:"code" validate:"required,min=1"`
	State       string `json:"state" validate:"required,min=1"`
	RedirectURL string `json:"redirect_url" validate:"required,url"`
	ReturnType  string `json:"type" validate:"required,oneof=token session"`
}
