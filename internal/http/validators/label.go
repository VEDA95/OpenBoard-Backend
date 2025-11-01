package validators

type CreateLabelValidator struct {
	Name    string `json:"name" validate:"required,min=1"`
	Color   string `json:"color" validate:"required,len=7"`
	BoardID string `json:"board_id" validate:"required,uuid"`
}

type UpdateLabelValidator struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Color *string `json:"color,omitempty" validate:"omitempty,len=7"`
}
