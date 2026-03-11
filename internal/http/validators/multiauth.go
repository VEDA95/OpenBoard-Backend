package validators

type SetupEmailTOTPValidator struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type VerifyEmailTOTPValidator struct {
	Code string `json:"code" validate:"required,len=6,numeric"`
}

type DeleteMFAMethodValidator struct {
	MethodID string `json:"method_id" validate:"required,uuid"`
}
