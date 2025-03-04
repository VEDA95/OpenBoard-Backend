package validators

type ParamValidator struct {
	Id string `validate:"required,uuid"`
}
