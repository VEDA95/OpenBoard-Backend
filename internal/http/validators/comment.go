package validators

type CreateCommentValidator struct {
	Comment string `json:"comment" validate:"required,min=1"`
	CardID  string `json:"card_id" validate:"required,uuid"`
}

type UpdateCommentValidator struct {
	Comment *string `json:"comment,omitempty" validate:"omitempty,min=1"`
}
