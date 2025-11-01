package validators

import "time"

type CreateCardValidator struct {
	Name          string     `json:"name" validate:"required,min=1"`
	Description   *string    `json:"description,omitempty"`
	ListID        string     `json:"list_id" validate:"required,uuid"`
	Color         *string    `json:"color,omitempty" validate:"omitempty,len=7"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	LabelIDs      *[]string  `json:"label_ids,omitempty"`
	AttachmentIDs *[]string  `json:"attachment_ids,omitempty"`
}

type UpdateCardValidator struct {
	Name          *string    `json:"name,omitempty" validate:"omitempty,min=1"`
	Description   *string    `json:"description,omitempty"`
	Color         *string    `json:"color,omitempty" validate:"omitempty,len=7"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
	LabelIDs      *[]string  `json:"label_ids,omitempty"`
	AttachmentIDs *[]string  `json:"attachment_ids,omitempty"`
}

type MoveCardValidator struct {
	ListID   *string `json:"list_id,omitempty" validate:"omitempty,uuid"`
	Position *int    `json:"position" validate:"required,min=0"`
}

type AddCommentValidator struct {
	Comment string `json:"comment" validate:"required,min=1"`
}

type ReorderListsValidator struct {
	Positions map[string]int `json:"positions" validate:"required"`
}
