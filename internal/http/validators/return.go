package validators

type ReturnValidator struct {
	ReturnType string `json:"type" validate:"required,oneof=token session"`
}
