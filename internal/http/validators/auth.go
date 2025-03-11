package validators

type LocalLoginValidator struct {
	ReturnType string `json:"type" validate:"required,oneof=token session"`
	Username   string `json:"username" validate:"required,min=1"`
	Password   string `json:"password" validate:"required,min=8"`
	Remember   bool   `json:"remember_me,omitempty" validate:"omitempty" default:"false"`
}

type CreateUserValidator struct {
	Username  string  `json:"username" validate:"required,min=1"`
	Email     string  `json:"email_address" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8,max=32"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1"`
}

type UpdateUserValidator struct {
	Username  *string `json:"username,omitempty" validate:"omitempty,min=1"`
	Email     *string `json:"email_address,omitempty" validate:"omitempty,email"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1"`
}

type ResetPasswordValidator struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=32"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
