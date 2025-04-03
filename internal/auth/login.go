package auth

type LocalUserLogin struct {
	Message          string  `json:"message"`
	User             *User   `json:"user"`
	AccessToken      string  `json:"access_token"`
	ExpiresIn        int     `json:"expires_in"`
	RefreshToken     *string `json:"refresh_token,omitempty" extensions:"x-nullable,x-omitempty"`
	RefreshExpiresIn *int    `json:"refresh_expires_in,omitempty" extensions:"x-nullable,x-omitempty"`
}
