package validators

type CreateListValidator struct {
	Name     string  `json:"name" validate:"required,min=1"`
	BoardID  string  `json:"board_id" validate:"required,uuid"`
	Color    *string `json:"color,omitempty" validate:"omitempty,len=7"`
	Position *int    `json:"position,omitempty" validate:"omitempty,min=0"`
}

type UpdateListValidator struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Color    *string `json:"color,omitempty" validate:"omitempty,len=7"`
	Position *int    `json:"position,omitempty" validate:"omitempty,min=0"`
}
