package validators

type CreateWorkspaceValidator struct {
	Name          string    `json:"name" validate:"required,min=1"`
	UserID        string    `json:"user" validate:"required,uuid"`
	Description   *string   `json:"description,omitempty" validate:"min=1,omitempty"`
	PermissionIDs *[]string `json:"permissions,omitempty" validate:"omitempty"`
}

type UpdateWorkspaceValidator struct {
	Name          *string   `json:"name,omitempty" validate:"min=1,omitempty"`
	UserID        *string   `json:"user,omitempty" validate:"uuid,omitempty"`
	Description   *string   `json:"description,omitempty"`
	PermissionIDs *[]string `json:"permissions,omitempty" validate:"omitempty"`
}

type CreateBoardValidator struct {
	Name        string `json:"name" validate:"required,min=1"`
	IsPublic    *bool  `json:"is_public,omitempty" validate:"bool,omitempty"`
	UserID      string `json:"user" validate:"required,uuid"`
	WorkspaceID string `json:"workspace" validate:"required,uuid"`
}

type UpdateBoardValidator struct {
	Name        *string `json:"name,omitempty" validate:"min=1,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty" validate:"bool,omitempty"`
	UserID      string  `json:"user,omitempty" validate:"uuid,omitempty"`
	WorkspaceID string  `json:"workspace,omitempty" validate:"uuid,omitempty"`
}
