package validators

type SetupEmailTOTPValidator struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type VerifyEmailTOTPValidator struct {
	Code string `json:"code" validate:"required,len=6,numeric"`
}

type SetupWebAuthnValidator struct {
	Name            string `json:"name" validate:"required,min=1,max=255"`
	CredentialID    string `json:"credential_id" validate:"required,min=1"`
	PublicKey       string `json:"public_key" validate:"required,min=1"`
	AttestationType string `json:"attestation_type" validate:"required,oneof=none indirect direct enterprise"`
	AAGUID          string `json:"aaguid,omitempty" validate:"omitempty"`
	Attachment      string `json:"attachment,omitempty" validate:"omitempty,oneof=platform cross-platform"`
}

type VerifyWebAuthnValidator struct {
	CredentialID      string `json:"credential_id" validate:"required,min=1"`
	AuthenticatorData string `json:"authenticator_data" validate:"required,min=1"`
	ClientDataJSON    string `json:"client_data_json" validate:"required,min=1"`
	Signature         string `json:"signature" validate:"required,min=1"`
	SignCount         uint32 `json:"sign_count" validate:"omitempty"`
}

type WebAuthnRegistrationOptionsValidator struct {
	AttestationType     string `json:"attestation_type,omitempty" validate:"omitempty,oneof=none indirect direct enterprise"`
	AuthenticatorType   string `json:"authenticator_type,omitempty" validate:"omitempty,oneof=platform cross-platform"`
	ResidentKeyRequired bool   `json:"resident_key_required,omitempty" validate:"omitempty"`
	UserVerification    string `json:"user_verification,omitempty" validate:"omitempty,oneof=required preferred discouraged"`
}

type WebAuthnAuthenticationOptionsValidator struct {
	UserVerification string `json:"user_verification,omitempty" validate:"omitempty,oneof=required preferred discouraged"`
}

type DeleteMFAMethodValidator struct {
	MethodID string `json:"method_id" validate:"required,uuid"`
}
