package validators

type CreateCheckListItemValidator struct {
	Name     string `json:"name" validate:"required,min=1"`
	CardID   string `json:"card_id" validate:"required,uuid"`
	UserID   string `json:"-" validate:"required,uuid"`
	Position *int   `json:"position,omitempty" validate:"omitempty,min=0"`
}

type UpdateCheckListItemValidator struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=1"`
	IsChecked *bool   `json:"is_checked,omitempty"`
	Position  *int    `json:"position,omitempty" validate:"omitempty,min=0"`
}

type ReorderCheckListItemsValidator struct {
	CardID    string         `json:"-" validate:"required,uuid"`
	Positions map[string]int `json:"positions" validate:"required"`
}
